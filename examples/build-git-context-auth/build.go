package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	auth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {

	var err error
	var dockerCli *client.Client

	registry := "registry"
	namespace := "namespace"
	imageName := strings.Join([]string{registry, namespace, "myimage"}, "/")

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	// authenticate to git server using a key
	authMethod := &auth.KeyAuth{
		PkFile:     "mykeyfile",
		PkPassword: "mypass",
	}
	//
	// Other auth methods
	//
	// - basic auth authentication ("github.com/apenella/go-docker-builder/pkg/auth/git/basic"):
	// 		authMethod := &auth.BasicAuth{
	// 			Username: "aleix.penella",
	// 			Password: "mypass",
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
	dockerBuilder.AddTags(strings.Join([]string{imageName, "tag1"}, ":"))
	dockerBuildContext := &gitcontext.GitBuildContext{
		Repository: "git@myprinvategitserver:base/alpine.git",
		Auth:       authMethod,
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
