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
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// go-docker-builder example where is created a ubuntu image
func main() {

	err := buildAndPush(os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func buildAndPush(w io.Writer) error {
	var err error
	var dockerCli *client.Client

	fmt.Println("Yeah! Build and push has started!")

	dockerBuildContext := &contextpath.PathBuildContext{
		Path: filepath.Join(".", "files"),
	}

	registry := "registry.go-docker-builder.test"
	imageName := strings.Join([]string{registry, "dummy-image-layers"}, "/")
	username := "admin"
	password := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("BuildAndPush", "Error on docker client creation", err)
	}

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend("buildAndPush"),
		),
		response.WithWriter(w),
	)

	dockerBuilder := &build.DockerBuildCmd{
		Cli:               dockerCli,
		ImageBuildOptions: &dockertypes.ImageBuildOptions{},
		ImagePushOptions:  &dockertypes.ImagePushOptions{},
		PushAfterBuild:    true,
		RemoveAfterPush:   true,
		ImageName:         imageName,
		ExecPrefix:        imageName,
		Response:          res,
	}

	dockerBuilder.AddTags(strings.Join([]string{imageName, "tag1"}, ":"))

	err = dockerBuilder.AddBuildContext(dockerBuildContext)
	if err != nil {
		return errors.New("BuildAndPush", "Error adding build docker context", err)
	}

	err = dockerBuilder.AddAuth(username, password, registry)
	if err != nil {
		return errors.New("BuildAndPush", "Error adding registry auth", err)
	}

	err = dockerBuilder.Run(context.TODO())
	if err != nil {
		return errors.New("BuildAndPush", fmt.Sprintf("Error building '%s'", imageName), err)
	}

	return nil
}
