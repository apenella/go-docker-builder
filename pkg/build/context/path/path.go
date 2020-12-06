package path

import (
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
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
		return nil, errors.New("(context::path::Reader)", fmt.Sprintf("Error opening '%s'", c.Path), err)
	}

	reader, err := tar.Tar(context)
	if err != nil {
		return nil, errors.New("(context::path::Reader)", fmt.Sprintf("Error archieving '%s'", c.Path), err)
	}

	return reader, nil
}
