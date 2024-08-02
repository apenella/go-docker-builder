package copy

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	auth "github.com/apenella/go-docker-builder/pkg/auth/docker"
	"github.com/apenella/go-docker-builder/pkg/push"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/apenella/go-docker-builder/pkg/types"
	dockerimagetypes "github.com/docker/docker/api/types/image"
)

// DockerCopyImageCmd is used to copy images to docker registry. Copy image is understood as tag an existing image and push it to a docker registry
type DockerImageCopyCmd struct {
	// Cli is the docker client to use
	Cli types.DockerClienter
	// ImagePushOptions from docker sdk
	ImagePullOptions *dockerimagetypes.PullOptions
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockerimagetypes.PushOptions
	// SourceImage is the name of the image to be copied
	SourceImage string
	// TargetImage is the name of the copied image
	TargetImage string
	// Tags is a copied images tags list
	Tags []string
	// UseNormalizedNamed when is true tags are transformed to a fully qualified reference
	UseNormalizedNamed bool
	// RemoteSource when is true the source image is pulled from registry before push it to its destination
	RemoteSource bool
	// RemoveAfterPush when is true the image from local is removed after push
	RemoveAfterPush bool
	// Response manages the docker client output
	Response types.Responser
}

// NewDockerImageCopyCmd return a DockerImageCopyCmd
func NewDockerImageCopyCmd(cli types.DockerClienter) *DockerImageCopyCmd {
	return &DockerImageCopyCmd{
		Cli:              cli,
		ImagePullOptions: &dockerimagetypes.PullOptions{},
		ImagePushOptions: &dockerimagetypes.PushOptions{},
	}
}

// WithSourceImage set tags to DockerImageCopyCmd
func (c *DockerImageCopyCmd) WithSourceImage(source string) *DockerImageCopyCmd {
	c.SourceImage = source
	return c
}

// WithTags set tags to DockerImageCopyCmd
func (c *DockerImageCopyCmd) WithTags(tags []string) *DockerImageCopyCmd {
	c.Tags = tags
	return c
}

// WithTargetImage set tags to DockerImageCopyCmd
func (c *DockerImageCopyCmd) WithTargetImage(target string) *DockerImageCopyCmd {
	c.TargetImage = target
	return c
}

// WithRemoteSource set to use remote source image to DockerImageCopyCmd
func (c *DockerImageCopyCmd) WithRemoteSource() *DockerImageCopyCmd {
	c.RemoteSource = true
	return c
}

// WithRemoveAfterPush set to remove source image once the image is pushed
func (c *DockerImageCopyCmd) WithRemoveAfterPush() *DockerImageCopyCmd {
	c.RemoveAfterPush = true
	return c
}

// WithResponse set responser attribute to DockerImageCopyCmd
func (c *DockerImageCopyCmd) WithResponse(res types.Responser) *DockerImageCopyCmd {
	c.Response = res
	return c
}

// WithUseNormalizedNamed set to use normalized named to DockerImageCopyCmd
func (c *DockerImageCopyCmd) WithUseNormalizedNamed() *DockerImageCopyCmd {
	c.UseNormalizedNamed = true
	return c
}

// AddAuth adds the same auth to image pull options and image push options
func (c *DockerImageCopyCmd) AddAuth(username, password string) error {
	var err error

	err = c.AddPullAuth(username, password)
	if err != nil {
		return errors.New("(copy::AddAuth)", "Error adding authorization to pull the source image", err)
	}

	err = c.AddPushAuth(username, password)
	if err != nil {
		return errors.New("(copy::AddAuth)", "Error adding authorization to push the copied image", err)
	}

	return nil
}

// AddPullAuth adds auth to pull the source image from remote location
func (c *DockerImageCopyCmd) AddPullAuth(username, password string) error {

	if c.ImagePullOptions == nil {
		c.ImagePullOptions = &dockerimagetypes.PullOptions{}
	}

	auth, err := auth.GenerateEncodedUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(copy::AddPullAuth)", "Error generating encoded user password auth configuration", err)
	}

	c.ImagePullOptions.RegistryAuth = *auth
	return nil
}

// AddPushAuth adds auth to push the tagged image to its destination
func (c *DockerImageCopyCmd) AddPushAuth(username, password string) error {

	if c.ImagePushOptions == nil {
		c.ImagePushOptions = &dockerimagetypes.PushOptions{}
	}

	auth, err := auth.GenerateEncodedUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(copy::AddPushAuth)", "Error generating encoded user password auth configuration", err)
	}

	c.ImagePushOptions.RegistryAuth = *auth
	return nil
}

// AddTag add a new copied image tag to tags list
func (c *DockerImageCopyCmd) AddTags(tag ...string) {

	if c.Tags == nil {
		c.Tags = []string{}
	}

	c.Tags = append(c.Tags, tag...)
}

// Run performs the image copy
func (c *DockerImageCopyCmd) Run(ctx context.Context) error {
	var err error
	var pullResponse io.ReadCloser

	if c == nil {
		return errors.New("(copy::Run)", "DockerImageCopyCmd is undefined")
	}

	if c.SourceImage == "" {
		return errors.New("(copy::Run)", "Source image must be defined")
	}

	if c.TargetImage == "" {
		return errors.New("(copy::Run)", "Target image must be defined")
	}

	if c.ImagePushOptions == nil {
		return errors.New("(copy::Run)", "Image push options is undefined")
	}

	if c.Response == nil {
		c.Response = response.NewDefaultResponse(
			response.WithWriter(os.Stdout),
		)
	}

	// if remote, pull
	if c.RemoteSource {
		if c.ImagePullOptions == nil {
			return errors.New("(copy::Run)", "Image pull options is undefined")
		}

		pullResponse, err = c.Cli.ImagePull(ctx, c.SourceImage, *c.ImagePullOptions)
		if err != nil {
			return errors.New("(copy::Run)", fmt.Sprintf("Error pull image '%s", c.SourceImage), err)
		}

		err = c.Response.Print(pullResponse)
		if err != nil {
			return errors.New("(copy::Run)", fmt.Sprintf("Error writing push response for '%s'", c.SourceImage), err)
		}
	}

	err = c.Cli.ImageTag(ctx, c.SourceImage, c.TargetImage)
	if err != nil {
		return errors.New("(copy::Run)", fmt.Sprintf("Error tagging image '%s' to '%s'", c.SourceImage, c.TargetImage), err)
	}

	push := &push.DockerPushCmd{
		Cli:                c.Cli,
		Response:           c.Response,
		ImageName:          c.TargetImage,
		Tags:               c.Tags,
		ImagePushOptions:   c.ImagePushOptions,
		UseNormalizedNamed: c.UseNormalizedNamed,
		RemoveAfterPush:    c.RemoveAfterPush,
	}

	err = push.Run(ctx)
	if err != nil {
		return errors.New("(copy::Run)", "Error pushing image", err)
	}

	return nil
}
