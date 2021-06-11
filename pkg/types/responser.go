package types

import "io"

const (
	LayerMessagePrefix string = "\u2023"
)

// Responser interface to write responses
type Responser interface {
	Print(io.ReadCloser) error
	Fwriteln(interface{})
}
