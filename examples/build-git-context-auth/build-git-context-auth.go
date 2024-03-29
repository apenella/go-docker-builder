package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	auth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	"github.com/apenella/go-docker-builder/pkg/response"

	// Uncomment the line below in case you want to run the example using the basic auth
	//auth "github.com/apenella/go-docker-builder/pkg/auth/git/basic"
	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
	"github.com/docker/docker/client"
)

func main() {

	err := buildGitContextAuth(os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func buildGitContextAuth(w io.Writer) error {
	var err error
	var dockerCli *client.Client

	registry := "registry.go-docker-builder.test"
	imageName := strings.Join([]string{registry, "alpine"}, "/")
	username := "admin"
	password := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("buildGitContextAuth", "Error on docker client creation", err)
	}

	// authenticate to git server using a key
	authMethod := &auth.KeyAuth{
		PkFile:     "/root/.ssh/id_rsa",
		PkPassword: "password",
	}

	//
	// Other auth methods
	//
	// - basic auth authentication ("github.com/apenella/go-docker-builder/pkg/auth/git/basic"):
	// 		authMethod := &auth.BasicAuth{
	// 			Username: "admin",
	// 			Password: "admin",
	// 		}
	//
	// - sshagent authentication ("github.com/apenella/go-docker-builder/pkg/auth/git/basic"):
	// 		authMethod := &auth.SSHAgentAuth{}
	//

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend("buildGitContextAuth"),
		),
		response.WithWriter(w),
	)

	dockerBuilder := build.NewDockerBuildCmd(dockerCli).
		WithImageName(imageName).
		WithDockerfile("Dockerfile.custom").
		WithPushAfterBuild().
		WithRemoveAfterPush().
		WithResponse(res)

	dockerBuilder.AddTags(strings.Join([]string{imageName, "custom"}, ":"))
	dockerBuildContext := &gitcontext.GitBuildContext{
		Repository: "git@gitserver:/git/repos/go-docker-builder-alpine.git",
		// Uncomment the line below in case you want to run the example using the basic auth
		// Repository: "http://gitserver/repos/go-docker-builder-alpine.git",
		Auth: authMethod,
	}

	err = dockerBuilder.AddBuildContext(dockerBuildContext)
	if err != nil {
		return errors.New("buildGitContextAuth", "Error adding build docker context", err)
	}

	err = dockerBuilder.AddAuth(username, password, registry)
	if err != nil {
		return errors.New("buildGitContextAuth", "Error adding registry auth", err)
	}

	err = dockerBuilder.Run(context.TODO())
	if err != nil {
		return errors.New("buildGitContextAuth", fmt.Sprintf("Error building '%s'", imageName), err)
	}

	return nil
}
