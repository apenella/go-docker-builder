package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/apenella/go-docker-builder/pkg/push"
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

	dockerPusher := &push.DockerPushCmd{
		Writer:     os.Stdout,
		Cli:        dockerCli,
		ImageName:  strings.Join([]string{registry, namespace, imageName}, "/"),
		ExecPrefix: imageName,
	}

	user := "myregistryuser"
	pass := "myregistrypass"
	dockerPusher.AddAuth(user, pass)

	err = dockerPusher.Run(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Error building '%s'. %s", imageName, err.Error()))
	}
}
