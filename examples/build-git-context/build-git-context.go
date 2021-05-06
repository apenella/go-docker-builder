package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {

	var err error
	var dockerCli *client.Client

	registry := "registry.go-docker-builder.test"
	imageName := strings.Join([]string{registry, "alpine"}, "/")
	username := "admin"
	password := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	dockerBuilder := &build.DockerBuildCmd{
		Writer:            os.Stdout,
		Cli:               dockerCli,
		ImageName:         imageName,
		ExecPrefix:        imageName,
		ImageBuildOptions: &dockertypes.ImageBuildOptions{},
		ImagePushOptions:  &dockertypes.ImagePushOptions{},
		PushAfterBuild:    true,
		RemoveAfterPush:   true,
	}

	dockerBuilder.AddTags(strings.Join([]string{imageName, "tag1"}, ":"))
	dockerBuildContext := &gitcontext.GitBuildContext{
		Repository: "https://github.com/alpinelinux/docker-alpine.git",
	}
	err = dockerBuilder.AddBuildContext(dockerBuildContext)
	if err != nil {
		panic(fmt.Sprintf("Error adding build docker context. %s", err.Error()))
	}

	err = dockerBuilder.AddAuth(username, password, registry)
	if err != nil {
		panic(fmt.Sprintf("Error adding registry auth. %s", err.Error()))
	}

	err = dockerBuilder.Run(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Error building '%s'. %s", imageName, err.Error()))
	}
}
