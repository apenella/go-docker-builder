package mock

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/mock"
)

// DockerClient is a docker client mock
type DockerClient struct {
	mock.Mock
}

// NewDockerClient creates a new mock for docker client
func NewDockerClient() *DockerClient {
	return &DockerClient{}
}

// ImageBuild is mock method to build docker images
func (client *DockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options dockertypes.ImageBuildOptions) (dockertypes.ImageBuildResponse, error) {
	args := client.Mock.Called(ctx, buildContext, options)
	return args.Get(0).(dockertypes.ImageBuildResponse), args.Error(1)
}

// ImagePush is a mock method to push docker images to registry
func (client *DockerClient) ImagePush(ctx context.Context, image string, options dockertypes.ImagePushOptions) (io.ReadCloser, error) {
	args := client.Mock.Called(ctx, image, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

// ImageRemove is a mock method to remove docker images locally
func (client *DockerClient) ImageRemove(ctx context.Context, imageID string, options dockertypes.ImageRemoveOptions) ([]dockertypes.ImageDeleteResponseItem, error) {
	args := client.Mock.Called(ctx, imageID, options)

	return args.Get(0).([]dockertypes.ImageDeleteResponseItem), args.Error(1)
}
