package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	auth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	// Uncomment the line below in case you want to run the example using the basic auth
	//auth "github.com/apenella/go-docker-builder/pkg/auth/git/basic"
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

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
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

	dockerBuilder := &build.DockerBuildCmd{
		Writer: os.Stdout,
		Cli:    dockerCli,
		ImageBuildOptions: &dockertypes.ImageBuildOptions{
			Dockerfile: "Dockerfile.custom",
		},
		ImageName:  imageName,
		ExecPrefix: imageName,
	}
	dockerBuilder.AddTags(strings.Join([]string{imageName, "custom"}, ":"))
	dockerBuildContext := &gitcontext.GitBuildContext{
		Repository: "git@gitserver:/git/repos/go-docker-builder-alpine.git",
		// Uncomment the line below in case you want to run the example using the basic auth
		// Repository: "http://gitserver/repos/go-docker-builder-alpine.git",
		Auth: authMethod,
	}

	err = dockerBuilder.AddBuildContext(dockerBuildContext)
	if err != nil {
		panic(fmt.Sprintf("Error adding build docker context. %s", err.Error()))
	}

	err = dockerBuilder.Run(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Error building '%s'. %s", imageName, err.Error()))
	}
}
