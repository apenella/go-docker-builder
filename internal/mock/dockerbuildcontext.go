package mock

import (
	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	"github.com/stretchr/testify/mock"
)

type DockerBuildContext struct {
	mock.Mock
}

// NewDockerBuildContext creates a new mock for docker build context
func NewDockerBuildContext() *DockerBuildContext {
	return &DockerBuildContext{}
}

func (context *DockerBuildContext) GenerateContextFilesystem() (*filesystem.ContextFilesystem, error) {
	args := context.Mock.Called()
	return args.Get(0).(*filesystem.ContextFilesystem), args.Error(1)
}
