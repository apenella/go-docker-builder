package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/docker/docker/client"
)

func main() {
	err := buildGitContext(os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func buildGitContext(w io.Writer) error {
	var err error
	var dockerCli *client.Client

	registry := "registry.go-docker-builder.test"
	imageName := strings.Join([]string{registry, "alpine"}, "/")
	username := "admin"
	password := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("BuildGitContext", "Error on docker client creation", err)
	}

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend("buildGitContext"),
		),
		response.WithWriter(w),
	)

	dockerBuilder := build.NewDockerBuildCmd(dockerCli, imageName).
		WithPushAfterBuild().
		WithRemoveAfterPush().
		WithResponse(res)

	// dockerBuilder := &build.DockerBuildCmd{
	// 	Cli:               dockerCli,
	// 	ImageName:         imageName,
	// 	ImageBuildOptions: &dockertypes.ImageBuildOptions{},
	// 	ImagePushOptions:  &dockertypes.ImagePushOptions{},
	// 	PushAfterBuild:    true,
	// 	RemoveAfterPush:   true,
	// 	Response:          res,
	// }

	dockerBuilder.AddTags(strings.Join([]string{imageName, "a-tag"}, ":"),
		strings.Join([]string{imageName, "b-tag"}, ":"),
		strings.Join([]string{imageName, "z-tag"}, ":"))

	dockerBuildContext := &gitcontext.GitBuildContext{
		Repository: "https://github.com/alpinelinux/docker-alpine.git",
		Reference:  "v3.13",
		Path:       "x86_64",
	}
	err = dockerBuilder.AddBuildContext(dockerBuildContext)
	if err != nil {
		return errors.New("BuildGitContext", "Error adding build docker context", err)
	}

	err = dockerBuilder.AddAuth(username, password, registry)
	if err != nil {
		return errors.New("BuildGitContext", "Error adding registry auth", err)
	}

	err = dockerBuilder.Run(context.TODO())
	if err != nil {
		return errors.New("BuildGitContext", fmt.Sprintf("Error building '%s'", imageName), err)
	}

	return nil
}
