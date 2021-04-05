package mock

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/mock"
)

type DockerClient struct {
	mock.Mock
}

func NewDockerClient() *DockerClient {
	return &DockerClient{}
}

func (client *DockerClient) ImagePush(ctx context.Context, image string, options dockertypes.ImagePushOptions) (io.ReadCloser, error) {
	args := client.Mock.Called(ctx, image, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
