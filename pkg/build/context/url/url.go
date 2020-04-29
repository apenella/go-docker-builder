package url

import (
	"errors"
	"io"
)

func Tar(url string) (io.Reader, error) {
	return nil, errors.New("(context::url::Tar) URL context not available")
}
