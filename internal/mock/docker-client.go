package mock

import (
	"context"
	"io"

	"github.com/moby/moby/client"
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
func (c *DockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options client.ImageBuildOptions) (client.ImageBuildResult, error) {
	args := c.Mock.Called(ctx, buildContext, options)
	return args.Get(0).(client.ImageBuildResult), args.Error(1)
}

// ImagePull is a mock method to pull docker images from registry
func (c *DockerClient) ImagePull(ctx context.Context, ref string, options client.ImagePullOptions) (client.ImagePullResponse, error) {
	args := c.Mock.Called(ctx, ref, options)
	return args.Get(0).(client.ImagePullResponse), args.Error(1)
}

// ImagePush is a mock method to push docker images to registry
func (c *DockerClient) ImagePush(ctx context.Context, image string, options client.ImagePushOptions) (client.ImagePushResponse, error) {
	args := c.Mock.Called(ctx, image, options)
	return args.Get(0).(client.ImagePushResponse), args.Error(1)
}

// ImageRemove is a mock method to remove docker images locally
func (c *DockerClient) ImageRemove(ctx context.Context, imageID string, options client.ImageRemoveOptions) (client.ImageRemoveResult, error) {
	args := c.Mock.Called(ctx, imageID, options)

	return args.Get(0).(client.ImageRemoveResult), args.Error(1)
}

// ImageTag is a mock method to tag docker images
func (c *DockerClient) ImageTag(ctx context.Context, options client.ImageTagOptions) (client.ImageTagResult, error) {
	args := c.Mock.Called(ctx, options)

	return args.Get(0).(client.ImageTagResult), args.Error(1)
}
