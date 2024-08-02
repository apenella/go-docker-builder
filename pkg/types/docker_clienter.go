package types

import (
	"context"
	"io"

	dockertypes "github.com/docker/docker/api/types"
	dockerimagetypes "github.com/docker/docker/api/types/image"
)

type DockerClienter interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options dockertypes.ImageBuildOptions) (dockertypes.ImageBuildResponse, error)
	ImagePull(ctx context.Context, ref string, options dockerimagetypes.PullOptions) (io.ReadCloser, error)
	ImagePush(ctx context.Context, image string, options dockerimagetypes.PushOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, imageID string, options dockerimagetypes.RemoveOptions) ([]dockerimagetypes.DeleteResponse, error)
	ImageTag(ctx context.Context, imageID, ref string) error
}
