package main

import (
	"context"
	"os"
	"strings"

	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
	"github.com/docker/docker/client"
)

func main() {

	var err error
	var dockerCli *client.Client

	registry := "registry"
	namespace := "namespace"
	imageName := strings.Join([]string{registry, namespace, "ubuntu"}, "/")

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	dockerBuildContext := &gitcontext.GitBuildContext{
		Repository: "https://github.com/alpinelinux/docker-alpine.git",
	}

	dockerBuildOptions := &build.DockerBuildOptions{
		ImageName:          imageName,
		Tags:               []string{strings.Join([]string{imageName, "tag1"}, ":")},
		DockerBuildContext: dockerBuildContext,
	}

	dockerBuilder := &build.DockerBuildCmd{
		Writer:             os.Stdout,
		Cli:                dockerCli,
		Context:            context.TODO(),
		DockerBuildOptions: dockerBuildOptions,
		ExecPrefix:         imageName,
	}

	err = dockerBuilder.Run()
	if err != nil {
		panic("Error building '" + imageName + "'. " + err.Error())
	}
}
