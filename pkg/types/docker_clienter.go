package types

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
)

type DockerClienter interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options dockertypes.ImageBuildOptions) (dockertypes.ImageBuildResponse, error)
	ImagePull(ctx context.Context, ref string, options dockertypes.ImagePullOptions) (io.ReadCloser, error)
	ImagePush(ctx context.Context, image string, options dockertypes.ImagePushOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, imageID string, options dockertypes.ImageRemoveOptions) ([]dockertypes.ImageDeleteResponseItem, error)
	ImageTag(ctx context.Context, imageID, ref string) error
}
