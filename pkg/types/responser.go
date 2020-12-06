package types

import "io"

const (
	separator string = "\u2023"
)

type Responser interface {
	Write(io.Writer, io.ReadCloser) error
}
