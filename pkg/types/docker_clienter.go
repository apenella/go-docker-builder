package types

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
)

type DockerClienter interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options dockertypes.ImageBuildOptions) (dockertypes.ImageBuildResponse, error)
	ImagePush(ctx context.Context, image string, options dockertypes.ImagePushOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, imageID string, options dockertypes.ImageRemoveOptions) ([]dockertypes.ImageDeleteResponseItem, error)
}
