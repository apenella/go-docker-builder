package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apenella/go-docker-builder/pkg/build"
	contextpath "github.com/apenella/go-docker-builder/pkg/build/context/path"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	var err error
	var dockerCli *client.Client

	imageDefinitionPath := filepath.Join(".", "files")

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
		ImageBuildOptions: &dockertypes.ImageBuildOptions{},
		ImagePushOptions:  &dockertypes.ImagePushOptions{},
		PushAfterBuild:    true,
		RemoveAfterPush:   true,
		ImageName:         imageName,
		ExecPrefix:        imageName,
	}

	dockerBuilder.AddTags(strings.Join([]string{imageName, "tag1"}, ":"))
	dockerBuildContext := &contextpath.PathBuildContext{
		Path: imageDefinitionPath,
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
