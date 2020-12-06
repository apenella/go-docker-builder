package push

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/apenella/go-docker-builder/pkg/types"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerPushCmd struct {
	Writer            io.Writer
	Context           context.Context
	Cli               *client.Client
	DockerPushOptions *DockerPushOptions
	ExecPrefix        string
	Response          types.Responser
}

func (p *DockerPushCmd) Run() error {

	if p == nil {
		return errors.New("(pusher:Run)", "DockerBuilder is nil")
	}

	if p.Writer == nil {
		p.Writer = os.Stdout
	}

	if p.Response == nil {
		p.Response = &response.DefaultResponse{
			Prefix: p.ExecPrefix,
		}
	}

	pushOptions := dockertypes.ImagePushOptions{}

	if p.DockerPushOptions.RegistryAuth != nil {
		pushOptions.RegistryAuth = *p.DockerPushOptions.RegistryAuth
	}

	pushResponse, err := p.Cli.ImagePush(p.Context, p.DockerPushOptions.ImageName, pushOptions)
	if err != nil {
		return errors.New("(pusher:Run)", fmt.Sprintf("Error pushing image '%s'", p.DockerPushOptions.ImageName), err)
	}
	defer pushResponse.Close()

	err = p.Response.Write(p.Writer, pushResponse)
	if err != nil {
		return errors.New("(builder:Run)", fmt.Sprintf("Error writing push response for '%s'", p.DockerPushOptions.ImageName), err)
	}

	for _, tag := range p.DockerPushOptions.Tags {
		pushResponse, err = p.Cli.ImagePush(p.Context, tag, pushOptions)
		if err != nil {
			return errors.New("(pusher:Run)", fmt.Sprintf("Error pushing image '%s'", tag), err)
		}

		err = p.Response.Write(p.Writer, pushResponse)
		if err != nil {
			return errors.New("(builder:Run)", fmt.Sprintf("Error writing push response for '%s'", tag), err)
		}
	}

	return nil
}

func (p *DockerPushCmd) registryAuthenticationPrivilegedFunc() (string, error) {
	fmt.Println("required authorization")
	return *p.DockerPushOptions.RegistryAuth, nil
}
