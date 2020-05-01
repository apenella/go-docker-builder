package main

import (
	"context"
	"os"
	"strings"

	auth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
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
	dockerBuildContext := &gitcontext.GitBuildContext{
		Repository: "git@myprinvategitserver:base/alpine.git",
		Auth:       authMethod,
	}

	dockerBuildOptions := &build.DockerBuildOptions{
		ImageName:          imageName,
		Tags:               []string{strings.Join([]string{imageName, "tag1"}, ":")},
		DockerBuildContext: dockerBuildContext,
		Dockerfile:         "build/docker/ms-suppliers-analytics/Dockerfile",
	}

	dockerBuilder := &build.DockerBuildCmd{
		Writer:             os.Stdout,
		Cli:                dockerCli,
		Context:            context.TODO(),
		DockerBuildOptions: dockerBuildOptions,
		ExecPrefix:         imageName,
	}

	err = dockerBuilder.Run()
	if err != nil {
		panic("Error building '" + imageName + "'. " + err.Error())
	}
}
