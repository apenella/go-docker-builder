package main

import (
	"context"
	"fmt"
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

	dockerPush := &push.DockerPushCmd{
		Cli:       dockerCli,
		ImageName: strings.Join([]string{registry, namespace, imageName}, "/"),
	}

	user := "myregistryuser"
	pass := "myregistrypass"
	dockerPush.AddAuth(user, pass)

	err = dockerPush.Run(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Error building '%s'. %s", imageName, err.Error()))
	}
}
