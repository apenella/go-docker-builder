package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	auth "github.com/apenella/go-docker-builder/pkg/auth/git/key"
	"github.com/apenella/go-docker-builder/pkg/build"
	gitcontext "github.com/apenella/go-docker-builder/pkg/build/context/git"
	pathcontext "github.com/apenella/go-docker-builder/pkg/build/context/path"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// go-docker-builder example where is created a ubuntu image
func main() {
	var err error
	var dockerCli *client.Client

	registry := "registry.go-docker-builder.test"
	imageName := strings.Join([]string{registry, "go-dummy-app"}, "/")
	registryUsername := "admin"
	registryPassword := "admin"

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

	codeContext := &gitcontext.GitBuildContext{
		Repository: "git@gitserver:/git/repos/go-dummy-app.git",
		// Uncomment the line below in case you want to run the example using the basic auth
		// Repository: "http://gitserver/repos/go-dummy-app.git",
		Auth: authMethod,
	}

	dockerContextInjection := &pathcontext.PathBuildContext{
		Path: filepath.Join(".", "files", "injection"),
	}

	dockerBuilder := &build.DockerBuildCmd{
		Writer:           os.Stdout,
		Cli:              dockerCli,
		ImagePushOptions: &dockertypes.ImagePushOptions{},
		PushAfterBuild:   true,
		RemoveAfterPush:  true,
		ImageName:        imageName,
		ExecPrefix:       imageName,
	}

	err = dockerBuilder.AddAuth(registryUsername, registryPassword, registry)
	if err != nil {
		panic(fmt.Sprintf("Error adding registry auth. %s", err.Error()))
	}

	err = dockerBuilder.AddBuildContext(codeContext, dockerContextInjection)
	if err != nil {
		panic(fmt.Sprintf("Error adding build docker context. %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	err = dockerBuilder.Run(ctx)
	if err != nil {
		panic(fmt.Sprintf("Error building '%s'. %s", imageName, err.Error()))
	}
}
