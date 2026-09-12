package types

import (
	"context"
	"io"

	"github.com/moby/moby/client"
)

type DockerClienter interface {
	ImageBuild(ctx context.Context, buildContext io.Reader, options client.ImageBuildOptions) (client.ImageBuildResult, error)
	ImagePull(ctx context.Context, ref string, options client.ImagePullOptions) (client.ImagePullResponse, error)
	ImagePush(ctx context.Context, image string, options client.ImagePushOptions) (client.ImagePushResponse, error)
	ImageRemove(ctx context.Context, imageID string, options client.ImageRemoveOptions) (client.ImageRemoveResult, error)
	ImageTag(ctx context.Context, options client.ImageTagOptions) (client.ImageTagResult, error)
}
