package git

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// TODO:
// use any folder inside the repository as context's base

// GitBuildContext defines a build context from a git repository
type GitBuildContext struct {
	// Repository which will be used as docker build context
	Repository string
	// Reference is the name of the branch to clone. By default is used 'master'
	Reference string
	// Dockerfile is the dockerfile placement inside the repository
	Dockerfile string
	// Auth
	// TODO
}

// Reader return a context reader
func (c *GitBuildContext) Reader() (io.Reader, error) {

	var tarBuff bytes.Buffer

	//aeromemfs := afero.NewMemMapFs()
	gitstorage := memory.NewStorage()

	referenceName := plumbing.Master
	if c.Reference != "" {
		referenceName = plumbing.NewBranchReferenceName(c.Reference)
	}

	cloneOption := &git.CloneOptions{
		URL:           c.Repository,
		ReferenceName: referenceName,
		Depth:         1,
	}

	repo, err := git.Clone(gitstorage, nil, cloneOption)
	if err != nil {
		return nil, errors.New("(context::git::Reader) Error cloning '" + c.Repository + "'." + err.Error())
	}

	referenceHead, err := repo.Head()
	if err != nil {
		return nil, errors.New("(context::git::Reader) Error getting reference '" + referenceName.String() + "' HEAD." + err.Error())
	}

	commit, err := repo.CommitObject(referenceHead.Hash())
	if err != nil {
		return nil, errors.New("(context::git::Reader) Error getting commit '" + referenceHead.Hash().String() + "'. " + err.Error())
	}

	filesIterator, err := commit.Files()
	defer filesIterator.Close()
	if err != nil {
		return nil, errors.New("(context::git::Reader) Error getting files iterator")
	}

	tw := tar.NewWriter(&tarBuff)
	defer tw.Close()

	err = filesIterator.ForEach(func(file *object.File) error {

		var buff bytes.Buffer

		fileContents, err := file.Contents()
		if err != nil {
			return errors.New("(context::git::Reader) Error achiving '" + file.Name + "' contents. " + err.Error())
		}

		_, err = buff.WriteString(fileContents)
		if err != nil {
			return errors.New("(context::git::Reader) Error writting '" + file.Name + "' contents to temporal buffer. " + err.Error())
		}

		header := &tar.Header{
			Name: file.Name,
			Size: file.Size,
			Mode: int64(file.Mode),
		}

		if err := tw.WriteHeader(header); err != nil {
			return errors.New("(context::git::Reader) Error writing '" + file.Name + "' header. " + err.Error())
		}
		// write file content into tar writer
		fmt.Fprint(tw, buff.String())

		buff.Reset()

		return nil
	})
	if err != nil {
		return nil, errors.New("(context::git::Reader) Error packaging repository files. " + err.Error())
	}

	return bytes.NewReader(tarBuff.Bytes()), nil
}
