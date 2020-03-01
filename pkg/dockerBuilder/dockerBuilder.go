package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/apenella/go-docker-builder/pkg/types"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerBuilderCmd
type DockerBuilderCmd struct {
	Writer               io.Writer
	Context              context.Context
	Cli                  *client.Client
	DockerBuilderContext *DockerBuilderContext
	DockerBuilderOptions *DockerBuilderOptions
	ExecPrefix           string
}

// Run execute the docker build
// https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
func (b *DockerBuilderCmd) Run() error {

	var err error
	var contextReader io.Reader

	if b == nil {
		return errors.New("(builder:Run) DockerBuilder is nil")
	}

	if b.Writer == nil {
		b.Writer = os.Stdout
	}

	contextReader, err = b.DockerBuilderContext.GenerateDockerBuilderContext()
	if err != nil {
		return errors.New("(builder:Run) Error generating Docker building context. " + err.Error())
	}

	options := dockertypes.ImageBuildOptions{
		Context:        contextReader,
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Dockerfile:     b.DockerBuilderOptions.Dockerfile,
		Tags:           b.DockerBuilderOptions.Tags,
		BuildArgs:      b.DockerBuilderOptions.BuildArgs,
	}

	buildResponse, err := b.Cli.ImageBuild(b.Context, contextReader, options)
	if err != nil {
		return errors.New("(builder:Run) Error building '" + b.DockerBuilderOptions.ImageName + "'." + err.Error())
	}
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)
	prefix := b.ExecPrefix

	for scanner.Scan() {
		streamMessage := &types.BuildResponseBodyStreamMessage{}
		line := scanner.Bytes()
		err = json.Unmarshal(line, &streamMessage)
		if err != nil {
			return errors.New("(builder:Run) Error unmarshalling line '" + string(line) + "' " + err.Error())
		}

		fmt.Fprintf(b.Writer, "%s \u2500\u2500  %s\n", prefix, streamMessage.String())
	}

	return nil
}
