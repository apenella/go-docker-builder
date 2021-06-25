package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/build"
	contextpath "github.com/apenella/go-docker-builder/pkg/build/context/path"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/docker/docker/client"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	err := buildPathContext(os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func buildPathContext(w io.Writer) error {
	var err error
	var dockerCli *client.Client

	imageDefinitionPath := filepath.Join(".", "files")
	registry := "registry.go-docker-builder.test"
	imageName := strings.Join([]string{registry, "alpine"}, "/")
	username := "admin"
	password := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("buildPathContext", "Error on docker client creation", err)
	}

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend("buildPathContext"),
		),
		response.WithWriter(w),
	)

	dockerBuilder := &build.DockerBuildCmd{
		Cli:       dockerCli,
		ImageName: imageName,
		Response:  res,
	}

	dockerBuilder.AddTags(strings.Join([]string{imageName, "tag1"}, ":"))
	dockerBuildContext := &contextpath.PathBuildContext{
		Path: imageDefinitionPath,
	}
	err = dockerBuilder.AddBuildContext(dockerBuildContext)
	if err != nil {
		return errors.New("buildPathContext", "Error adding build docker context", err)
	}

	err = dockerBuilder.AddAuth(username, password, registry)
	if err != nil {
		return errors.New("buildPathContext", "Error adding registry auth", err)
	}

	err = dockerBuilder.Run(context.TODO())
	if err != nil {
		return errors.New("buildPathContext", fmt.Sprintf("Error building '%s'", imageName), err)
	}

	return nil
}
