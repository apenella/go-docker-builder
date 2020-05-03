package path

import (
	"errors"
	"io"
	"os"

	"github.com/apenella/go-docker-builder/pkg/common/tar"
)

// PathBuildContext creates a build context from path
type PathBuildContext struct {
	// Path is context location on the local system
	Path string
}

// Reader return a context reader
func (c *PathBuildContext) Reader() (io.Reader, error) {

	context, err := os.Open(c.Path)
	if err != nil {
		return nil, errors.New("(context::path::Reader) Error opening '" + c.Path + "'. " + err.Error())
	}

	reader, err := tar.Tar(context)
	if err != nil {
		return nil, errors.New("(context::path::Reader) Error archieving '" + c.Path + "'. " + err.Error())
	}

	return reader, nil
}
