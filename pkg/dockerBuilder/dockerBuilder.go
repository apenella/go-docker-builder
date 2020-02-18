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

type DockerBuilder struct {
	ImageName            string
	Writer               io.Writer
	Context              context.Context
	DockerBuilderContext *DockerBuilderContext
	Cli                  *client.Client
	Tags                 []string
	BuildArgs            map[string]*string
	Dockerfile           string
	PushImage            bool
	ExecPrefix           string
}

// Run execute the docker build
func (b *DockerBuilder) Run() error {

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
		Dockerfile:     b.Dockerfile,
		Tags:           b.Tags,
		BuildArgs:      b.BuildArgs,
	}

	buildResponse, err := b.Cli.ImageBuild(b.Context, contextReader, options)
	if err != nil {
		return errors.New("(builder:Run) Error building '" + b.ImageName + "'." + err.Error())
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

// AddBuildArgs append new tags to DockerBuilder
func (b *DockerBuilder) AddBuildArgs(arg string, value string) error {

	_, exists := b.BuildArgs[arg]
	if exists {
		return errors.New("(builder::AddBuildArgs) Argument '" + arg + "' already defined")
	}

	b.BuildArgs[arg] = &value
	return nil
}

// AddTags append new tags to DockerBuilder
func (b *DockerBuilder) AddTags(tag string) {
	b.Tags = append(b.Tags, tag)
}
