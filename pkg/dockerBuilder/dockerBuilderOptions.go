package builder

import "errors"

// DockerBuilderOptions has an options set to build and image
type DockerBuilderOptions struct {
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the image tags
	Tags []string
	// BuildArgs ia a list of arguments to set during the building
	BuildArgs map[string]*string
	// Dockerfile is the file name for dockerfile file
	Dockerfile string
	// PushImage set to push or not an image
	PushImage bool
}

// AddBuildArgs append new tags to DockerBuilder
func (o *DockerBuilderOptions) AddBuildArgs(arg string, value string) error {

	_, exists := o.BuildArgs[arg]
	if exists {
		return errors.New("(builder::AddBuildArgs) Argument '" + arg + "' already defined")
	}

	o.BuildArgs[arg] = &value
	return nil
}

// AddTags append new tags to DockerBuilder
func (o *DockerBuilderOptions) AddTags(tag string) {
	o.Tags = append(o.Tags, tag)
}
