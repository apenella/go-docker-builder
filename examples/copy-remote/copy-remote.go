package main

import (
	"context"
	"fmt"
	"io"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/copy"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/docker/docker/client"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	err := copyRemote(os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func copyRemote(w io.Writer) error {
	var err error
	var dockerCli *client.Client

	sourceImage := "base-registry.go-docker-builder.test:5000/alpine:3.13"
	targetImage := "registry.go-docker-builder.test/alpine:3.13"

	targetRegistryUsername := "admin"
	targetRegistryPassword := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("copyRemote", "Error on docker client creation", err)
	}

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend("copyRemote"),
		),
		response.WithWriter(w),
	)

	copy := copy.NewDockerImageCopyCmd(dockerCli).
		WithSourceImage(sourceImage).
		WithTargetImage(targetImage).
		WithRemoteSource().
		WithRemoveAfterPush().
		WithResponse(res)

	copy.AddTags("registry.go-docker-builder.test/alpine/alpine:three",
		"registry.go-docker-builder.test/alpine/alpine:latest")

	err = copy.AddPushAuth(targetRegistryUsername, targetRegistryPassword)
	if err != nil {
		return errors.New("copyRemote", "Error adding registry auth", err)
	}

	err = copy.Run(context.TODO())
	if err != nil {
		return errors.New("copyRemote", fmt.Sprintf("Error copying '%s' to '%s'", sourceImage, targetImage), err)
	}

	return nil
}
