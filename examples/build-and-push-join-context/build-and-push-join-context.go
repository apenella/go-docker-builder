package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	auth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
	pathcontext "github.com/apenella/go-docker-builder/pkg/build/context/path"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/docker/docker/client"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	err := buildAndPushJoinContext(os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func buildAndPushJoinContext(w io.Writer) error {

	var err error
	var dockerCli *client.Client

	registry := "registry.go-docker-builder.test"
	imageName := strings.Join([]string{registry, "go-dummy-app"}, "/")
	registryUsername := "admin"
	registryPassword := "admin"

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return errors.New("buildAndPushJoinContext", "Error on docker client creation", err)
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

	codeContext := &gitcontext.GitBuildContext{
		Repository: "git@gitserver:/git/repos/go-dummy-app.git",
		// Uncomment the line below in case you want to run the example using the basic auth
		// Repository: "http://gitserver/repos/go-dummy-app.git",
		Auth: authMethod,
	}

	dockerContextInjection := &pathcontext.PathBuildContext{
		Path: filepath.Join(".", "files", "injection"),
	}

	res := response.NewDefaultResponse(
		response.WithTransformers(
			transformer.Prepend("buildAndPushJoinContext"),
		),
		response.WithWriter(w),
	)

	dockerBuilder := build.NewDockerBuildCmd(dockerCli, imageName).
		WithPushAfterBuild().
		WithRemoveAfterPush().
		WithResponse(res)

	err = dockerBuilder.AddAuth(registryUsername, registryPassword, registry)
	if err != nil {
		return errors.New("buildAndPushJoinContext", "Error adding registry auth", err)
	}

	err = dockerBuilder.AddBuildContext(codeContext, dockerContextInjection)
	if err != nil {
		return errors.New("buildAndPushJoinContext", "Error adding build docker context", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	err = dockerBuilder.Run(ctx)
	if err != nil {
		return errors.New("buildAndPushJoinContext", fmt.Sprintf("Error building '%s'", imageName), err)
	}

	return nil
}
