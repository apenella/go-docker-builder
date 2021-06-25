package path

import (
	"io"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	"github.com/spf13/afero"
)

// PathBuildContext creates a build context from path
type PathBuildContext struct {
	// Path is context location on the local system
	Path string
}

// Reader return a context reader
func (c *PathBuildContext) Reader() (io.Reader, error) {
	errorContext := "(context::path::Reader)"
	// context, err := os.Open(c.Path)
	// if err != nil {
	// 	return nil, errors.New("(context::path::Reader)", fmt.Sprintf("Error opening '%s'", c.Path), err)
	// }

	// reader, err := tar.Tar(context)
	// if err != nil {
	// 	return nil, errors.New("(context::path::Reader)", fmt.Sprintf("Error archieving '%s'", c.Path), err)
	// }

	// return reader, nil
	fs, err := c.GenerateContextFilesystem()
	if err != nil {
		return nil, errors.New(errorContext, "Error packaging repository files", err)
	}
	return fs.Tar()
}

func (c *PathBuildContext) GenerateContextFilesystem() (*filesystem.ContextFilesystem, error) {
	fs := filesystem.NewContextFilesystem(afero.NewOsFs())
	fs.RootPath = c.Path

	return fs, nil
}
