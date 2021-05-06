package mock

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type DockerBuildContext struct {
	mock.Mock
}

// NewDockerBuildContext creates a new mock for docker build context
func NewDockerBuildContext() *DockerBuildContext {
	return &DockerBuildContext{}
}

// Reader is mock method to build docker images
func (context *DockerBuildContext) Reader() (io.Reader, error) {
	args := context.Mock.Called()
	return args.Get(0).(io.Reader), args.Error(1)
}
