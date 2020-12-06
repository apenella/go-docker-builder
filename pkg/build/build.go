package build

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/push"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/apenella/go-docker-builder/pkg/types"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	// DefaultDockerfile is the default filename for Dockerfile
	DefaultDockerfile string = "Dockerfile"
)

// DockerBuilderCmd
type DockerBuildCmd struct {
	// Writer to write the build output
	Writer io.Writer
	// Context manages the build context
	Context context.Context
	// Cli is the docker api client
	Cli *client.Client
	// DockerBuildOptions are the options to build
	DockerBuildOptions *DockerBuildOptions
	// DockerPushOptions are the option to push
	DockerPushOptions *push.DockerPushOptions
	// ExecPrefix defines a prefix to each output lines
	ExecPrefix string
	// Response manages responses from docker client
	Response types.Responser
}

// Run execute the docker build
// https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
func (b *DockerBuildCmd) Run() error {

	var err error
	var contextReader io.Reader

	if b == nil {
		return errors.New("(builder:Run)", "DockerBuilder is nil")
	}

	if b.Writer == nil {
		b.Writer = os.Stdout
	}

	if b.Response == nil {
		b.Response = &response.DefaultResponse{
			Prefix: b.ExecPrefix,
		}
	}

	if b.DockerBuildOptions.ImageName == "" {
		return errors.New("(builder:Run)", "An image name is required to build an image")
	}

	if b.DockerBuildOptions.Tags == nil {
		b.DockerBuildOptions.Tags = []string{b.DockerBuildOptions.ImageName}
	} else {
		b.DockerBuildOptions.Tags = append(b.DockerBuildOptions.Tags, b.DockerBuildOptions.ImageName)
	}

	if b.DockerBuildOptions.Dockerfile == "" {
		b.DockerBuildOptions.Dockerfile = DefaultDockerfile
	}

	contextReader, err = b.DockerBuildOptions.DockerBuildContext.Reader()
	if err != nil {
		return errors.New("(builder:Run)", "Error generating a build context reader", err)
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
		AuthConfigs:    b.DockerBuildOptions.Auth,
	}

	buildResponse, err := b.Cli.ImageBuild(b.Context, contextReader, buildOptions)
	if err != nil {
		return errors.New("(builder:Run)", fmt.Sprintf("Error building '%s'", b.DockerBuildOptions.ImageName), err)
	}
	defer buildResponse.Body.Close()

	err = b.Response.Write(b.Writer, buildResponse.Body)
	if err != nil {
		return errors.New("(builder:Run)", fmt.Sprintf("Error writing build response for '%s'", b.DockerBuildOptions.ImageName), err)
	}

	if b.DockerBuildOptions.PushAfterBuild {
		dockerPush := &push.DockerPushCmd{
			Writer:            b.Writer,
			Cli:               b.Cli,
			Context:           b.Context,
			DockerPushOptions: b.DockerPushOptions,
			ExecPrefix:        b.ExecPrefix,
		}

		err = dockerPush.Run()
		if err != nil {
			return errors.New("(builder:Run)", fmt.Sprintf("Error pushing image '%s'", b.DockerBuildOptions.ImageName), err)
		}
	}

	return nil
}
