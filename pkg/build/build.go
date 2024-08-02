package build

import (
	"context"
	"fmt"
	"os"

	errors "github.com/apenella/go-common-utils/error"
	auth "github.com/apenella/go-docker-builder/pkg/auth/docker"
	buildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/pkg/build/context/filesystem"
	"github.com/apenella/go-docker-builder/pkg/push"
	"github.com/apenella/go-docker-builder/pkg/response"
	"github.com/apenella/go-docker-builder/pkg/types"
	"github.com/distribution/reference"
	dockertypes "github.com/docker/docker/api/types"
	dockerimagetypes "github.com/docker/docker/api/types/image"
	dockerregistrytypes "github.com/docker/docker/api/types/registry"
	"github.com/spf13/afero"
)

const (
	// DefaultDockerfile is the default filename for Dockerfile
	DefaultDockerfile string = "Dockerfile"
)

// DockerBuilderCmd
type DockerBuildCmd struct {
	// Cli is the docker api client
	Cli types.DockerClienter
	// ImageName is the name of the image
	ImageName string
	// ImageBuildOptions from docker sdk
	ImageBuildOptions *dockertypes.ImageBuildOptions
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockerimagetypes.PushOptions
	// PullParentImage if true pull parent image
	PullParentImage bool
	// PushAfterBuild when is true images are automatically pushed to registry after build
	PushAfterBuild bool
	// Response manages responses from docker client
	Response types.Responser
	// UseNormalizedNamed when is true tags are transformed to a fully qualified reference
	UseNormalizedNamed bool
	// RemoveAfterPush when is true images are removed from local after push
	RemoveAfterPush bool
}

// NewDockerBuildCmd return a DockerBuildCmd
func NewDockerBuildCmd(cli types.DockerClienter) *DockerBuildCmd {
	return &DockerBuildCmd{
		Cli:               cli,
		ImageBuildOptions: &dockertypes.ImageBuildOptions{},
		ImagePushOptions:  &dockerimagetypes.PushOptions{},
	}
}

// WithDockerfile set responser attribute to DockerBuildCmd
func (b *DockerBuildCmd) WithDockerfile(dockerfile string) *DockerBuildCmd {
	if b.ImageBuildOptions == nil {
		b.ImageBuildOptions = &dockertypes.ImageBuildOptions{}
	}

	b.ImageBuildOptions.Dockerfile = dockerfile

	return b
}

// WithImageName set to push image automatically after its build
func (b *DockerBuildCmd) WithImageName(name string) *DockerBuildCmd {
	b.ImageName = name
	return b
}

// WithPullParentImage set to pull parent image
func (b *DockerBuildCmd) WithPullParentImage() *DockerBuildCmd {
	if b.ImageBuildOptions == nil {
		b.ImageBuildOptions = &dockertypes.ImageBuildOptions{}
	}

	b.ImageBuildOptions.PullParent = true
	return b
}

// WithPushAfterBuild set to push image automatically after its build
func (b *DockerBuildCmd) WithPushAfterBuild() *DockerBuildCmd {
	b.PushAfterBuild = true
	return b
}

// WithResponse set responser attribute to DockerBuildCmd
func (b *DockerBuildCmd) WithResponse(res types.Responser) *DockerBuildCmd {
	b.Response = res
	return b
}

// WithUseNormalizedNamed set to use normalized named to DockerBuildCmd
func (b *DockerBuildCmd) WithUseNormalizedNamed() *DockerBuildCmd {
	b.UseNormalizedNamed = true
	return b
}

// WithRemoveAfterPush set to remove source image once the image is pushed
func (b *DockerBuildCmd) WithRemoveAfterPush() *DockerBuildCmd {
	b.RemoveAfterPush = true
	return b
}

// AddAuth append new tags to DockerBuilder
func (b *DockerBuildCmd) AddAuth(username, password, registry string) error {

	if b.ImageBuildOptions == nil {
		b.ImageBuildOptions = &dockertypes.ImageBuildOptions{}
	}

	if b.ImageBuildOptions.AuthConfigs == nil {
		b.ImageBuildOptions.AuthConfigs = map[string]dockerregistrytypes.AuthConfig{}
	}

	authConfig, err := auth.GenerateUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(build::AddAuth)", "Error generation user password auth configuration", err)
	}

	b.ImageBuildOptions.AuthConfigs[registry] = *authConfig
	return nil
}

// AddPushAuth append new tags to DockerBuilder
func (b *DockerBuildCmd) AddPushAuth(username, password string) error {

	if b.ImagePushOptions == nil {
		b.ImagePushOptions = &dockerimagetypes.PushOptions{}
	}

	auth, err := auth.GenerateEncodedUserPasswordAuthConfig(username, password)
	if err != nil {
		return errors.New("(build::AddPushAuth)", "Error generating encoded user password auth configuration", err)
	}

	b.ImagePushOptions.RegistryAuth = *auth
	return nil
}

// AddBuildArgs append new tags to DockerBuilder. Returns an error when adding an existing argument
func (b *DockerBuildCmd) AddBuildArgs(arg string, value string) error {

	if b.ImageBuildOptions == nil {
		b.ImageBuildOptions = &dockertypes.ImageBuildOptions{}
	}

	if b.ImageBuildOptions.BuildArgs == nil {
		b.ImageBuildOptions.BuildArgs = map[string]*string{}
	}

	_, exists := b.ImageBuildOptions.BuildArgs[arg]
	if exists {
		return errors.New("(build::AddBuildArgs)", fmt.Sprintf("Argument '%s' already defined", arg))
	}

	b.ImageBuildOptions.BuildArgs[arg] = &value
	return nil
}

// AddBuildContext include the docker build context. It supports to use several context which are merged before to start the image build
func (b *DockerBuildCmd) AddBuildContext(dockercontexts ...buildcontext.DockerBuildContexter) error {
	var err error
	errorContext := "(build::AddBuilderContext)"
	dockercontext := filesystem.NewContextFilesystem(afero.NewMemMapFs())

	if b.ImageBuildOptions == nil {
		b.ImageBuildOptions = &dockertypes.ImageBuildOptions{}
	}

	for _, dc := range dockercontexts {
		var cfs *filesystem.ContextFilesystem

		if dc == nil {
			return errors.New(errorContext, "Docker build context is not defined")
		}

		cfs, err = dc.GenerateContextFilesystem()
		if err != nil {
			return errors.New(errorContext, "Error generationg context filesystem", err)
		}

		dockercontext, err = filesystem.Join(true, dockercontext, cfs)
		if err != nil {
			return errors.New(errorContext, "Error joining docker context", err)
		}
	}

	dockerBuildContextReader, err := dockercontext.Tar()
	if err != nil {
		return errors.New(errorContext, "Error creating docker build context reader", err)
	}

	b.ImageBuildOptions.Context = dockerBuildContextReader

	return nil
}

// AddLabel append new tags to DockerBuilder. Returns an error when adding an existing label
func (b *DockerBuildCmd) AddLabel(label string, value string) error {

	if b.ImageBuildOptions == nil {
		b.ImageBuildOptions = &dockertypes.ImageBuildOptions{}
	}

	if b.ImageBuildOptions.Labels == nil {
		b.ImageBuildOptions.Labels = map[string]string{}
	}

	_, exists := b.ImageBuildOptions.Labels[label]
	if exists {
		return errors.New("(build::AddLabel)", fmt.Sprintf("Label '%s' already defined", label))
	}

	b.ImageBuildOptions.Labels[label] = value

	return nil
}

// AddTags append new tags to DockerBuilder
func (b *DockerBuildCmd) AddTags(tags ...string) error {

	if b.ImageBuildOptions == nil {
		b.ImageBuildOptions = &dockertypes.ImageBuildOptions{}
	}

	if b.ImageBuildOptions.Tags == nil {
		b.ImageBuildOptions.Tags = []string{}
	}

	for _, tag := range tags {
		if b.UseNormalizedNamed {
			normalizedTag, err := reference.ParseNormalizedNamed(tag)
			if err != nil {
				return errors.New("(build::AddTags)", fmt.Sprintf("Tag '%s' could not be normalized", tag), err)
			}
			tag = normalizedTag.String()
		}
		b.ImageBuildOptions.Tags = append(b.ImageBuildOptions.Tags, tag)
	}

	return nil
}

// Run execute the docker build
// https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
func (b *DockerBuildCmd) Run(ctx context.Context) error {

	var err error

	if b == nil {
		return errors.New("(build::Run)", "DockerBuildCmd is not defined")
	}

	if b.ImageBuildOptions == nil {
		return errors.New("(build::Run)", "ImageBuildOptions options is not defined")
	}

	if b.PushAfterBuild && b.ImagePushOptions == nil {
		return errors.New("(build::Run)", "ImagePushOptions options is not defined")
	}

	if b.ImageBuildOptions.Context == nil {
		return errors.New("(build::Run)", "Docker build context is not defined")
	}

	if b.Response == nil {
		b.Response = response.NewDefaultResponse(
			response.WithWriter(os.Stdout),
		)
	}

	if b.ImageName == "" {
		return errors.New("(build::Run)", "An image name is required to build an image")
	}
	b.AddTags(b.ImageName)

	if b.ImageBuildOptions.Dockerfile == "" {
		b.ImageBuildOptions.Dockerfile = DefaultDockerfile
	}

	buildResponse, err := b.Cli.ImageBuild(ctx, b.ImageBuildOptions.Context, *b.ImageBuildOptions)
	if err != nil {
		return errors.New("(build::Run)", fmt.Sprintf("Error building '%s'", b.ImageName), err)
	}
	defer buildResponse.Body.Close()

	err = b.Response.Print(buildResponse.Body)
	if err != nil {
		return errors.New("(build::Run)", fmt.Sprintf("Error writing build response for '%s'", b.ImageName), err)
	}

	if b.PushAfterBuild {

		dockerPush := &push.DockerPushCmd{
			Cli:              b.Cli,
			ImageName:        b.ImageName,
			ImagePushOptions: b.ImagePushOptions,
			Tags:             b.ImageBuildOptions.Tags,
			Response:         b.Response,
		}

		if b.RemoveAfterPush {
			dockerPush.RemoveAfterPush = b.RemoveAfterPush
		}

		// in case that build auth is set configure it on push image
		if b.ImageBuildOptions.AuthConfigs != nil {
			named, _ := reference.ParseNormalizedNamed(b.ImageName)
			registryHost := reference.Domain(named)

			registryHostAuth, ok := b.ImageBuildOptions.AuthConfigs[registryHost]

			if ok {
				if dockerPush.ImagePushOptions.RegistryAuth == "" {
					auth.GenerateEncodedUserPasswordAuthConfig(registryHostAuth.Username, registryHostAuth.Password)
					if registryHostAuth.Auth != "" {
						dockerPush.ImagePushOptions.RegistryAuth = registryHostAuth.Auth
					} else {
						registryHostBuildAuth, err := auth.GenerateEncodedUserPasswordAuthConfig(registryHostAuth.Username, registryHostAuth.Password)
						if err != nil {
							return errors.New("(build::Run)", "Error encoding username and password", err)
						}
						dockerPush.ImagePushOptions.RegistryAuth = *registryHostBuildAuth
					}
				} else {
					dockerPush.ImagePushOptions.PrivilegeFunc = generateDefaultImagePushOptionsPrivilegeFunc(registryHostAuth.Auth)
				}
			}
		}

		err = dockerPush.Run(ctx)
		if err != nil {
			return errors.New("(build::Run)", fmt.Sprintf("Error pushing image '%s'", b.ImageName), err)
		}
	}

	return nil
}

func generateDefaultImagePushOptionsPrivilegeFunc(auth string) dockertypes.RequestPrivilegeFunc {
	return func(context.Context) (string, error) {
		return auth, nil
	}
}
