package main

import (
	"context"
	"os"
	"strings"

	"github.com/apenella/go-docker-builder/pkg/push"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/docker/docker/client"
)

func main() {

	var err error
	var dockerCli *client.Client

	registry := "registry"
	namespace := "namespace"
	imageName := strings.Join([]string{registry, namespace, "ubuntu"}, "/")

	dockerCli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic("Error on docker client creation. " + err.Error())
	}

	dockerPushOptions := &push.DockerPushOptions{
		ImageName: strings.Join([]string{registry, namespace, imageName}, "/"),
	}

	user := "myregistryuser"
	pass := "myregistrypass"
	dockerPushOptions.AddUserPasswordRegistryAuth(user, pass, registry)

	response := &response.ResponseHandler{
		Prefix: imageName,
	}

	dockerPusher := &push.DockerPushCmd{
		Writer:            os.Stdout,
		Cli:               dockerCli,
		Context:           context.TODO(),
		DockerPushOptions: dockerPushOptions,
		ExecPrefix:        imageName,
		Response:          response,
	}

	err = dockerPusher.Run()
	if err != nil {
		panic("Error pushing '" + imageName + "'. " + err.Error())
	}
}
