package context

import "github.com/apenella/go-docker-builder/pkg/build/context/filesystem"

// DockerBuildContexter defines a docker build context entity
type DockerBuildContexter interface {
	//Reader() (io.Reader, error)
	GenerateContextFilesystem() (*filesystem.ContextFilesystem, error)
}
