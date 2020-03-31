package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"

	"github.com/apenella/go-docker-builder/pkg/build"
	"github.com/apenella/go-docker-builder/pkg/push"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	var err error
	var dockerCli *client.Client

	imageDefinitionPath := filepath.Join(".", "files")

	registry := "registry"
	namespace := "namespace"
	imageName := strings.Join([]string{registry, namespace, "alpine"}, "/")
	username := "user"
	password := "pass"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	dockerBuildContext := &build.DockerBuildContext{
		Path: imageDefinitionPath,
	}

	dockerBuildOptions := &build.DockerBuildOptions{
		ImageName:      imageName,
		Dockerfile:     "Dockerfile",
		Tags:           []string{strings.Join([]string{imageName, "tag1"}, ":")},
		PushAfterBuild: true,
	}

	dockerPushOptions := &push.DockerPushOptions{
		ImageName: dockerBuildOptions.ImageName,
		Tags:      dockerBuildOptions.Tags,
	}
	dockerPushOptions.AddUserPasswordRegistryAuth(username, password, registry)

	dockerBuilder := &build.DockerBuildCmd{
		Writer:             os.Stdout,
		Cli:                dockerCli,
		Context:            context.TODO(),
		DockerBuildContext: dockerBuildContext,
		DockerBuildOptions: dockerBuildOptions,
		DockerPushOptions:  dockerPushOptions,
		ExecPrefix:         imageName,
	}

	err = dockerBuilder.Run()
	if err != nil {
		panic("Error building '" + imageName + "'. " + err.Error())
	}
}
