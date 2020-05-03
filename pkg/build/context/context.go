package context

import "io"

// DockerBuildContexter defines a docker build context entity
type DockerBuildContexter interface {
	Reader() (io.Reader, error)
}
