package push

import (
	"errors"

	auth "github.com/apenella/go-docker-builder/pkg/auth/docker"
)

// DockerBuilderOptions has an options set to build and image
type DockerPushOptions struct {
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the images to push
	Tags []string
	// RegistryAuth is the base64 encoded credentials for the registry
	RegistryAuth *string
}

// AddAuth append new tags to DockerBuilder
func (o *DockerPushOptions) AddAuth(username, password string) error {

	auth, err := auth.GenerateEncodedUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(push::AddAuth) " + err.Error())
	}

	o.RegistryAuth = auth
	return nil
}
