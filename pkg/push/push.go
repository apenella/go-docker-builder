package push

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	auth "github.com/apenella/go-docker-builder/pkg/auth/docker"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/apenella/go-docker-builder/pkg/types"
	dockertypes "github.com/docker/docker/api/types"
)

type PusherClient interface {
	ImagePush(ctx context.Context, image string, options dockertypes.ImagePushOptions) (io.ReadCloser, error)
}

// DockerPushCmd is used to push images to docker registry
type DockerPushCmd struct {
	// Writer to use to write docker client messges
	Writer io.Writer
	// Cli is the docker client to use
	Cli PusherClient
	// ImagePushOptions docker sdk push options
	ImagePushOptions *dockertypes.ImagePushOptions
	// ExecPrefix prefix to include add to each docker client message
	ExecPrefix string
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the images to push
	Tags []string
	// Response manages the docker client output
	Response types.Responser
}

// Run performs the push action
func (p *DockerPushCmd) Run(ctx context.Context) error {

	var err error
	var pushResponse io.ReadCloser

	if p == nil {
		return errors.New("(push::Run)", "DockerPushCmd is undefined")
	}

	if p.ImagePushOptions == nil {
		return errors.New("(push::Run)", "Image push options is undefined")
	}

	if p.Writer == nil {
		p.Writer = os.Stdout
	}

	if p.Response == nil {
		p.Response = &response.DefaultResponse{
			Prefix: p.ExecPrefix,
		}
	}

	images := []string{p.ImageName}
	if len(p.Tags) > 0 {
		images = append(images, p.Tags...)
	}

	for _, image := range images {
		pushResponse, err = p.Cli.ImagePush(ctx, image, *p.ImagePushOptions)
		if err != nil {
			return errors.New("(push::Run)", fmt.Sprintf("Error pushing image '%s'", image), err)
		}

		err = p.Response.Write(p.Writer, pushResponse)
		if err != nil {
			return errors.New("(push::Run)", fmt.Sprintf("Error writing push response for '%s'", image), err)
		}
	}

	return nil
}

// AddAuth append new tags to DockerBuilder
func (p *DockerPushCmd) AddAuth(username, password string) error {

	if p.ImagePushOptions == nil {
		p.ImagePushOptions = &dockertypes.ImagePushOptions{}
	}

	auth, err := auth.GenerateEncodedUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(push::AddAuth)", "Error generating encoded user password auth configuration", err)
	}

	p.ImagePushOptions.RegistryAuth = *auth
	return nil
}
