package build

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/apenella/go-docker-builder/pkg/types"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	// DefaultDockerfile is the default filename for Dockerfile
	DefaultDockerfile = "Dockerfile"
)

// DockerBuilderCmd
type DockerBuildCmd struct {
	Writer             io.Writer
	Context            context.Context
	Cli                *client.Client
	DockerBuildContext *DockerBuildContext
	DockerBuildOptions *DockerBuildOptions
	ExecPrefix         string
	Response           types.Responser
}

// Run execute the docker build
// https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
func (b *DockerBuildCmd) Run() error {

	var err error
	var contextReader io.Reader

	if b == nil {
		return errors.New("(builder:Run) DockerBuilder is nil")
	}

	if b.Writer == nil {
		b.Writer = os.Stdout
	}

	if b.Response == nil {
		b.Response = &response.DefaultResponse{
			Prefix: b.ExecPrefix,
		}
	}

	contextReader, err = b.DockerBuildContext.GenerateDockerBuildContext()
	if err != nil {
		return errors.New("(builder:Run) Error generating Docker building context. " + err.Error())
	}

	if b.DockerBuildOptions.ImageName == "" {
		return errors.New("(builder:Run) An image name is required to build an image")
	}

	if b.DockerBuildOptions.Tags == nil {
		b.DockerBuildOptions.Tags = []string{b.DockerBuildOptions.ImageName}
	} else {
		b.DockerBuildOptions.Tags = append(b.DockerBuildOptions.Tags, b.DockerBuildOptions.ImageName)
	}

	if b.DockerBuildOptions.Dockerfile == "" {
		b.DockerBuildOptions.Dockerfile = DefaultDockerfile
	}

	buildOptions := dockertypes.ImageBuildOptions{
		Context:        contextReader,
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Dockerfile:     b.DockerBuildOptions.Dockerfile,
		Tags:           b.DockerBuildOptions.Tags,
		BuildArgs:      b.DockerBuildOptions.BuildArgs,
	}

	buildResponse, err := b.Cli.ImageBuild(b.Context, contextReader, buildOptions)
	if err != nil {
		return errors.New("(builder:Run) Error building '" + b.DockerBuildOptions.ImageName + "'. " + err.Error())
	}
	defer buildResponse.Body.Close()

	err = b.Response.Write(b.Writer, buildResponse.Body)
	if err != nil {
		return errors.New("(builder:Run) " + err.Error())
	}

	return nil
}
