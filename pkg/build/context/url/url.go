package url

import (
	"io"

	errors "github.com/apenella/go-common-utils/error"
)

// URLBuildContext creates a build context from url
type URLBuildContext struct {
	// URL is a web resource contexts location
	URL string
}

// Reader return a context reader
func (c *URLBuildContext) Reader() (io.Reader, error) {
	return nil, errors.New("(context::url::Tar)", "URL context not available")
}
