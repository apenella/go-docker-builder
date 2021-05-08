package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apenella/go-docker-builder/pkg/copy"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	var err error
	var dockerCli *client.Client

	sourceImage := "other-registry.go-docker-builder.test:5000/alpine:3.13"
	targetImage := "registry.go-docker-builder.test/alpine:3.13"

	targetRegistryUsername := "admin"
	targetRegistryPassword := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	copy := &copy.DockerImageCopyCmd{
		Writer:           os.Stdout,
		Cli:              dockerCli,
		ExecPrefix:       "copy alpine:3.13",
		SourceImage:      sourceImage,
		TargetImage:      targetImage,
		ImagePullOptions: &dockertypes.ImagePullOptions{},
		ImagePushOptions: &dockertypes.ImagePushOptions{},
		RemoteSource:     true,
	}

	err = copy.AddPushAuth(targetRegistryUsername, targetRegistryPassword)
	if err != nil {
		panic(fmt.Sprintf("Error adding registry auth. %s", err.Error()))
	}

	err = copy.Run(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Error copying '%s' to '%s'. %s", sourceImage, targetImage, err.Error()))
	}
}
