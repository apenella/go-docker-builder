package push

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Responser interface {
	Run() error
	SetReader(io.ReadCloser)
	SetWriter(io.Writer)
}

type DockerPushCmd struct {
	Writer            io.Writer
	Context           context.Context
	Cli               *client.Client
	DockerPushOptions *DockerPushOptions
	ExecPrefix        string
	Response          Responser
}

func (p *DockerPushCmd) Run() error {

	if p == nil {
		return errors.New("(pusher:Run) DockerBuilder is nil")
	}

	if p.Writer == nil {
		p.Writer = os.Stdout
	}

	pushOptions := dockertypes.ImagePushOptions{}

	if p.DockerPushOptions.RegistryAuth != nil {
		pushOptions.RegistryAuth = *p.DockerPushOptions.RegistryAuth
	}

	pushResponse, err := p.Cli.ImagePush(p.Context, p.DockerPushOptions.ImageName, pushOptions)
	if err != nil {
		return errors.New("(pusher:Run) Error push '" + p.DockerPushOptions.ImageName + "'. " + err.Error())
	}
	//fmt.Println(pushResponse)
	defer pushResponse.Close()

	p.Response.SetReader(pushResponse)
	p.Response.SetWriter(p.Writer)
	err = p.Response.Run()
	if err != nil {
		return errors.New("(builder:Run) " + err.Error())
	}

	return nil
}

func (p *DockerPushCmd) registryAuthenticationPrivilegedFunc() (string, error) {
	fmt.Println("required authorization")
	return *p.DockerPushOptions.RegistryAuth, nil
}
