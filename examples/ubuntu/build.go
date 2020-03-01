package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"

	builder "github.com/apenella/go-docker-builder/pkg/dockerBuilder"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	var err error
	var dockerCli *client.Client

	image := "ubuntu"
	imageDefinitionPath := filepath.Join(".", "build")
	imageName := strings.Join([]string{"gobuild", image}, "-")

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	dockerBuilderContext := &builder.DockerBuilderContext{
		Path: imageDefinitionPath,
	}

	dockerBuilderOptions := &builder.DockerBuilderOptions{
		ImageName:  imageName,
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
	}

	dockerBuilder := &builder.DockerBuilder{
		Writer:               os.Stdout,
		Cli:                  dockerCli,
		Context:              context.TODO(),
		DockerBuilderContext: dockerBuilderContext,
		DockerBuilderOptions: dockerBuilderOptions,
		ExecPrefix:           imageName,
	}

	err = dockerBuilder.Run()
	if err != nil {
		panic("Error building '" + imageName + "'. " + err.Error())
	}
}
