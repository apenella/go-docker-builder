package builder

import (
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DockerBuilderContext defines the building context
type DockerBuilderContext struct {
	// Path defines the path location when the contxext is placed on a local folder
	Path string `yaml:"path"`
	// URL defines the url when the contxext is located remotely and published via HTTP
	URL string `yaml:"status"`
	// Git defines git reposiotry url when the contxext is located remote repository
	Git string `yaml:"git"`
}

// GenerateDockerBuilderContext return io Reader for a given DockerBuilderContext
//   Build context precedence is:
//		1) Path
//		2) URL
//		3) Git
func (c *DockerBuilderContext) GenerateDockerBuilderContext() (io.Reader, error) {

	if c.Path != "" {
		return c.Tar()
	}

	if c.URL != "" {
		return nil, errors.New("(dockerBuilder::GenerateDockerBuilderContext) URL build context not already defined")
	}

	if c.Git != "" {
		return nil, errors.New("(dockerBuilder::GenerateDockerBuilderContext) Git build context not already defined")
	}

	return nil, errors.New("(dockerBuilder::GenerateDockerBuilderContext) No build context defined")
}

// Tar return a tarball io.Reader which contains the whole directory structure
//   This method has been inspeared thanks to https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func (c *DockerBuilderContext) Tar() (io.Reader, error) {

	var err error
	var tarBuff bytes.Buffer

	// ensure the src actually exists before trying to tar it
	_, err = os.Stat(c.Path)
	if err != nil {
		return nil, errors.New("(dockerBuilder::Tar) '" + c.Path + "' stat. " + err.Error())
	}

	tw := tar.NewWriter(&tarBuff)
	defer tw.Close()

	err = filepath.Walk(c.Path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return errors.New("(dockerBuilder::Tar::Walk) Error at the beginning of the walk. " + err.Error())
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.New("(dockerBuilder::Tar::Walk) Error creating '" + file + "' header. " + err.Error())
		}
		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, c.Path, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return errors.New("(dockerBuilder::Tar::Walk) Error writing '" + file + "' header. " + err.Error())
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return errors.New("(dockerBuilder::Tar::Walk) Error copying '" + file + "' into tar. " + err.Error())
		}
		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
	if err != nil {
		return nil, errors.New("(dockerBuilder::Tar) Error explorint '" + c.Path + "'. " + err.Error())
	}

	return bytes.NewReader(tarBuff.Bytes()), nil
}
