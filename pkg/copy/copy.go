package copy

import (
	"context"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	auth "github.com/apenella/go-docker-builder/pkg/auth/docker"
	"github.com/apenella/go-docker-builder/pkg/types"
	dockertypes "github.com/docker/docker/api/types"
)

// DockerCopyImageCmd is used to copy images to docker registry. Copy image is understood as tag an existing image and push it to a docker registry
type DockerImageCopyCmd struct {
	// Writer to use to write docker client messges
	Writer io.Writer
	// Cli is the docker client to use
	Cli types.DockerClienter
	// ImagePushOptions from docker sdk
	ImagePullOptions *dockertypes.ImagePullOptions
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockertypes.ImagePushOptions
	// ExecPrefix prefix to include add to each docker client message
	ExecPrefix string
	// SourceImage is the name of the image
	SourceImage string
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
		c.ImagePullOptions = &dockertypes.ImagePullOptions{}
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
		c.ImagePushOptions = &dockertypes.ImagePushOptions{}
	}

	auth, err := auth.GenerateEncodedUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(copy::AddPushAuth)", "Error generating encoded user password auth configuration", err)
	}

	c.ImagePushOptions.RegistryAuth = *auth
	return nil
}

// AddTag add a new copied image tag to tags list
func (c *DockerImageCopyCmd) AddTag(tag ...string) {

	if c.Tags == nil {
		c.Tags = []string{}
	}

	c.Tags = append(c.Tags, tag...)
}

// Run performs the image copy
func (c *DockerImageCopyCmd) Run(ctx context.Context) error {

	// if remote, pull

	// generate tags imageTag

	// prepare DockerPushCmd

	// remove tags and source image

	return nil
}
