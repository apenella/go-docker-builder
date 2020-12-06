package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apenella/go-docker-builder/pkg/build"
	contextpath "github.com/apenella/go-docker-builder/pkg/build/context/path"
	"github.com/apenella/go-docker-builder/pkg/push"
	"github.com/docker/docker/client"
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

	dockerBuildContext := &contextpath.PathBuildContext{
		Path: imageDefinitionPath,
	}

	dockerBuildOptions := &build.DockerBuildOptions{
		ImageName:          imageName,
		Dockerfile:         "Dockerfile",
		Tags:               []string{strings.Join([]string{imageName, "tag1"}, ":")},
		PushAfterBuild:     true,
		DockerBuildContext: dockerBuildContext,
	}

	dockerPushOptions := &push.DockerPushOptions{
		ImageName: dockerBuildOptions.ImageName,
		Tags:      dockerBuildOptions.Tags,
	}
	dockerPushOptions.AddAuth(username, password)
	err = dockerBuildOptions.AddAuth(username, password, registry)
	if err != nil {
		fmt.Println("Error including auth to build options: " + err.Error())
	}

	dockerBuilder := &build.DockerBuildCmd{
		Writer:             os.Stdout,
		Cli:                dockerCli,
		Context:            context.TODO(),
		DockerBuildOptions: dockerBuildOptions,
		DockerPushOptions:  dockerPushOptions,
		ExecPrefix:         imageName,
	}

	err = dockerBuilder.Run()
	if err != nil {
		panic("Error building '" + imageName + "'. " + err.Error())
	}
}
