package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"

	builder "github.com/apenella/go-docker-builder/pkg/dockerBuilder"
)

func main() {
	var err error
	var dockerCli *client.Client

	image := "ubuntu"
	imageDefinitionPath := filepath.Join("..", "tests", "dockerfiles", image)
	imageName := strings.Join([]string{"gobuild", filepath.Base(imageDefinitionPath)}, "-")

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	dockerBuilderContext := &builder.DockerBuilderContext{
		Path: imageDefinitionPath,
	}

	dockerBuilder := &builder.DockerBuilder{
		ImageName:            imageName,
		Writer:               os.Stdout,
		Context:              context.TODO(),
		DockerBuilderContext: dockerBuilderContext,
		Cli:                  dockerCli,
		Tags:                 []string{imageName},
		Dockerfile:           "Dockerfile",
		ExecPrefix:           imageName,
	}

	err = dockerBuilder.Run()
	if err != nil {
		panic("Error building '" + imageName + "'. " + err.Error())
	}
}
