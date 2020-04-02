package build

import (
	"errors"

	"github.com/apenella/go-docker-builder/pkg/auth"

	dockertypes "github.com/docker/docker/api/types"
)

// DockerBuildOptions has an options set to build and image
type DockerBuildOptions struct {
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the image tags
	Tags []string
	// BuildArgs ia a list of arguments to set during the building
	BuildArgs map[string]*string
	// Dockerfile is the file name for dockerfile file
	Dockerfile string
	// PushAfterBuild push image to registry after building
	PushAfterBuild bool
	// Auth
	Auth map[string]dockertypes.AuthConfig
}

// AddBuildArgs append new tags to DockerBuilder
func (o *DockerBuildOptions) AddBuildArgs(arg string, value string) error {

	if o.BuildArgs == nil {
		o.BuildArgs = map[string]*string{}
	}

	_, exists := o.BuildArgs[arg]
	if exists {
		return errors.New("(builder::AddBuildArgs) Argument '" + arg + "' already defined")
	}

	o.BuildArgs[arg] = &value
	return nil
}

// AddTags append new tags to DockerBuilder
func (o *DockerBuildOptions) AddTags(tag string) {

	if o.Tags == nil {
		o.Tags = []string{}
	}

	o.Tags = append(o.Tags, tag)
}

// AddAuth append new tags to DockerBuilder
func (o *DockerBuildOptions) AddAuth(username, password, registry string) error {

	if o.Auth == nil {
		o.Auth = map[string]dockertypes.AuthConfig{}
	}

	authConfig, err := auth.GenerateUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(build::AddAuth) " + err.Error())
	}

	o.Auth[registry] = *authConfig
	return nil
}
