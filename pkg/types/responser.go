package types

import "io"

type Responser interface {
	Write(io.Writer, io.ReadCloser) error
}
