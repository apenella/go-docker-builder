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
	"github.com/docker/distribution/reference"
	dockertypes "github.com/docker/docker/api/types"
)

// DockerPushCmd is used to push images to docker registry
type DockerPushCmd struct {
	// Writer to use to write docker client messges
	Writer io.Writer
	// Cli is the docker client to use
	Cli types.DockerClienter
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockertypes.ImagePushOptions
	// ExecPrefix prefix to include add to each docker client message
	ExecPrefix string
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the images to push
	Tags []string
	// Response manages the docker client output
	Response types.Responser
	// UseNormalizedNamed when is true tags are transformed to a fully qualified reference
	UseNormalizedNamed bool
	// RemoveAfterPush when is true the image from local is removed after push
	RemoveAfterPush bool
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

// AddTag append new tags to DockerBuilder
func (p *DockerPushCmd) AddTag(tags ...string) error {
	var err error
	var named reference.Named

	if p.Tags == nil {
		p.Tags = []string{}
	}

	for _, tag := range tags {
		exists := false

		if p.UseNormalizedNamed {
			named, err = reference.ParseNormalizedNamed(tag)
			if err != nil {
				return errors.New("(push::AddTag)", fmt.Sprintf("Error parsing to normalized named on '%s'", tag), err)
			}
			tag = named.String()
		}

		for _, t := range p.Tags {
			if t == tag {
				exists = true
			}
		}

		if !exists {
			p.Tags = append(p.Tags, tag)
		}
	}

	return nil
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

	p.AddTag(p.ImageName)

	for _, image := range p.Tags {
		pushResponse, err = p.Cli.ImagePush(ctx, image, *p.ImagePushOptions)
		if err != nil {
			return errors.New("(push::Run)", fmt.Sprintf("Error pushing image '%s'", image), err)
		}

		err = p.Response.Write(p.Writer, pushResponse)
		if err != nil {
			return errors.New("(push::Run)", fmt.Sprintf("Error writing push response for '%s'", image), err)
		}

		if p.RemoveAfterPush {
			deleteResponseItems, err := p.Cli.ImageRemove(ctx, image, dockertypes.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			})
			if err != nil {
				return errors.New("(push::Run)", fmt.Sprintf("Error removing '%s'", image), err)
			}

			// TODO
			fmt.Println(deleteResponseItems)
		}
	}

	return nil
}
