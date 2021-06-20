package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	errors "github.com/apenella/go-common-utils/error"
	auth "github.com/apenella/go-docker-builder/pkg/auth/git"
	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/afero"
)

// TODO:
// use any folder inside the repository as context's base

// GitBuildContext defines a build context from a git repository
type GitBuildContext struct {
	// Path identify where is located the context inside the repository
	Path string
	// Repository which will be used as docker build context
	Repository string
	// Reference is the name of the branch to clone. By default is used 'master'
	Reference string
	// Auth
	Auth auth.GitAuther
}

// Reader return a context reader
func (c *GitBuildContext) Reader() (io.Reader, error) {
	errorContext := "(context::git::Reader)"
	fs, err := c.GenerateContextFilesystem()
	if err != nil {
		return nil, errors.New(errorContext, "Error packaging repository files", err)
	}
	return fs.Tar()
}

func (c *GitBuildContext) GenerateContextFilesystem() (*filesystem.ContextFilesystem, error) {
	var err error
	fs := filesystem.NewContextFilesystem(afero.NewMemMapFs())

	// if c.Path != "" {
	// 	fs.RootPath = c.Path
	// }

	gitstorage := memory.NewStorage()
	errorContext := "(context::git::GenerateContextFilesystem)"

	referenceName := plumbing.Master
	if c.Reference != "" {
		referenceName = plumbing.NewBranchReferenceName(c.Reference)
	}

	cloneOption := &git.CloneOptions{
		URL:           c.Repository,
		ReferenceName: referenceName,
		Depth:         1,
		SingleBranch:  true,
	}

	if c.Auth != nil {
		auth, err := c.Auth.Auth()
		if err != nil {
			return nil, errors.New(errorContext, "Error getting authorization method", err)
		}

		cloneOption.Auth = auth
	}

	repo, err := git.Clone(gitstorage, nil, cloneOption)
	if err != nil {
		return nil, errors.New(errorContext, fmt.Sprintf("Error cloning repository '%s'", c.Repository), err)
	}

	referenceHead, err := repo.Head()
	if err != nil {
		return nil, errors.New(errorContext, fmt.Sprintf("Error getting reference '%s' HEAD", referenceName.String()), err)
	}

	commit, err := repo.CommitObject(referenceHead.Hash())
	if err != nil {
		return nil, errors.New(errorContext, fmt.Sprintf("Error getting commit '%s'", referenceHead.Hash().String()), err)
	}

	filesIterator, err := commit.Files()
	if err != nil {
		return nil, errors.New(errorContext, "Error getting files iterator")
	}
	defer filesIterator.Close()

	subPath := c.Path != ""
	subPathRegexp, _ := regexp.Compile(fmt.Sprintf("^%s", c.Path))

	err = filesIterator.ForEach(func(file *object.File) error {

		var buff bytes.Buffer
		var fileContents string
		var err error
		var memfile afero.File
		var mode os.FileMode
		var relativePath string

		subPathMatch := subPathRegexp.MatchString(file.Name)

		// skip files when subpath is defined and its location does not match to subpath base
		if subPath && !subPathMatch {
			return nil
		}

		fileContents, err = file.Contents()
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error achiving '%s' contents", file.Name), err)
		}
		mode, err = file.Mode.ToOSFileMode()
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error converting file mode to '%s' on '%s'", file.Mode.String(), file.Name), err)
		}

		relativePath, err = filepath.Rel(c.Path, file.Name)
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error creating relative path on '%s' from '%s'", file.Name, c.Path), err)
		}

		memfile, err = fs.OpenFile(filepath.Join(fs.RootPath, relativePath), os.O_CREATE, mode)
		if err != nil {
			return errors.New("(Walk)", fmt.Sprintf("Error extracting '%s' from git repository", relativePath), err)
		}

		_, err = buff.WriteString(fileContents)
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error writting '%s' contents to temporal buffer", relativePath), err)
		}

		io.Copy(memfile, io.Reader(&buff))

		buff.Reset()

		return nil
	})
	if err != nil {
		return nil, errors.New(errorContext, "Error packaging repository files", err)
	}

	return fs, nil
}
