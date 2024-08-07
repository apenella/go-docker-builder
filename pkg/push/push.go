package push

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	auth "github.com/apenella/go-docker-builder/pkg/auth/docker"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/apenella/go-docker-builder/pkg/types"
	"github.com/distribution/reference"
	dockerimagetypes "github.com/docker/docker/api/types/image"
)

// DockerPushCmd is used to push images to docker registry
type DockerPushCmd struct {
	// Cli is the docker client to use
	Cli types.DockerClienter
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockerimagetypes.PushOptions
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

// NewDockerPushCmd return a DockerPushCmd
func NewDockerPushCmd(cli types.DockerClienter) *DockerPushCmd {
	return &DockerPushCmd{
		Cli:              cli,
		ImagePushOptions: &dockerimagetypes.PushOptions{},
	}
}

// WithImageName set to push image automatically after its build
func (p *DockerPushCmd) WithImageName(name string) *DockerPushCmd {
	p.ImageName = name
	return p
}

// WithTags set tags to DockerPushCmd
func (p *DockerPushCmd) WithTags(tags []string) *DockerPushCmd {
	p.Tags = tags
	return p
}

// WithResponse set responser attribute to DockerPushCmd
func (p *DockerPushCmd) WithResponse(res types.Responser) *DockerPushCmd {
	p.Response = res
	return p
}

// WithRemoveAfterPush set to remove source image once the image is pushed
func (p *DockerPushCmd) WithRemoveAfterPush() *DockerPushCmd {
	p.RemoveAfterPush = true
	return p
}

// WithUseNormalizedNamed set to use normalized named to DockerPushCmd
func (p *DockerPushCmd) WithUseNormalizedNamed() *DockerPushCmd {
	p.UseNormalizedNamed = true
	return p
}

// AddAuth append new tags to DockerBuilder
func (p *DockerPushCmd) AddAuth(username, password string) error {

	if p.ImagePushOptions == nil {
		p.ImagePushOptions = &dockerimagetypes.PushOptions{}
	}

	auth, err := auth.GenerateEncodedUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(push::AddAuth)", "Error generating encoded user password auth configuration", err)
	}

	p.ImagePushOptions.RegistryAuth = *auth
	return nil
}

// AddTags append new tags to DockerBuilder
func (p *DockerPushCmd) AddTags(tags ...string) error {
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
				return errors.New("(push::AddTags)", fmt.Sprintf("Error parsing to normalized named on '%s'", tag), err)
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

	if p.ImageName == "" {
		return errors.New("(push::Run)", "Image name is undefined")
	}

	if p.ImagePushOptions == nil {
		return errors.New("(push::Run)", "Image push options is undefined")
	}

	if p.Response == nil {
		p.Response = response.NewDefaultResponse(
			response.WithWriter(os.Stdout),
		)
	}

	p.AddTags(p.ImageName)

	for _, image := range p.Tags {

		if image != p.ImageName {
			err = p.Cli.ImageTag(ctx, p.ImageName, image)
			if err != nil {
				return errors.New("(push::Run)", fmt.Sprintf("Error tagging image '%s' to '%s'", p.ImageName, image), err)
			}
		}

		pushResponse, err = p.Cli.ImagePush(ctx, image, *p.ImagePushOptions)
		if err != nil {
			return errors.New("(push::Run)", fmt.Sprintf("Error pushing image '%s'", image), err)
		}

		err = p.Response.Print(pushResponse)
		if err != nil {
			return errors.New("(push::Run)", fmt.Sprintf("Error writing push response for '%s'", image), err)
		}
	}

	if p.RemoveAfterPush {
		for _, image := range p.Tags {
			deleteResponseItems, err := p.Cli.ImageRemove(ctx, image, dockerimagetypes.RemoveOptions{
				Force:         true,
				PruneChildren: true,
			})
			if err != nil {
				return errors.New("(push::Run)", fmt.Sprintf("Error removing '%s'", image), err)
			}

			for _, item := range deleteResponseItems {

				str := ""
				if item.Deleted != "" {
					str = fmt.Sprintf("deleted: %s %s ", str, strings.TrimSpace(item.Deleted))
				}

				if item.Untagged != "" {
					str = fmt.Sprintf("untagged: %s %s ", str, strings.TrimSpace(item.Untagged))
				}

				if str != "" {
					p.Response.Fwriteln(str)
				}
			}
		}
	}

	return nil
}
