package build

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	buildcontext "github.com/apenella/go-docker-builder/pkg/build/context"
	"github.com/apenella/go-docker-builder/test/mock"
	mockclient "github.com/apenella/go-docker-builder/test/mock"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestAddAuth(t *testing.T) {

	type args struct {
		username string
		password string
		registry string
	}
	tests := []struct {
		desc           string
		dockerBuildCmd *DockerBuildCmd
		buildOptions   *dockertypes.ImageBuildOptions
		args           *args
		err            error
		res            map[string]dockertypes.AuthConfig
	}{
		{
			desc:           "Test add user-password auth",
			dockerBuildCmd: &DockerBuildCmd{},
			args: &args{
				username: "user",
				password: "AqSwd3Fr",
				registry: "registry",
			},
			err: nil,
			res: map[string]dockertypes.AuthConfig{
				"registry": {
					Username: "user",
					Password: "AqSwd3Fr",
				},
			},
		},
		{
			desc:           "Test add invalid user-password auth",
			dockerBuildCmd: &DockerBuildCmd{},
			args: &args{
				username: "",
				password: "AqSwd3Fr",
				registry: "registry",
			},
			err: errors.New("(build::AddAuth)", "Error generation user password auth configuration", errors.New(
				"(auth::GenerateUserPasswordAuthConfig)", "Username must be provided")),
			res: map[string]dockertypes.AuthConfig{},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.dockerBuildCmd.AddAuth(test.args.username, test.args.password, test.args.registry)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, test.dockerBuildCmd.ImageBuildOptions.AuthConfigs, "Unexpected auth result")
			}
		})
	}
}

func TestAddBuildArgs(t *testing.T) {

	type args struct {
		arg   string
		value string
	}

	tests := []struct {
		desc           string
		dockerBuildCmd *DockerBuildCmd
		args           *args
		err            error
	}{
		{
			desc:           "Test add argument to nil BuildArgs object",
			dockerBuildCmd: &DockerBuildCmd{},
			args: &args{
				arg:   "argument",
				value: "value",
			},
			err: nil,
		},
		{
			desc: "Test add an existing argument",
			dockerBuildCmd: &DockerBuildCmd{
				ImageBuildOptions: &dockertypes.ImageBuildOptions{
					BuildArgs: map[string]*string{
						"argument": nil,
					},
				},
			},
			args: &args{
				arg:   "argument",
				value: "value",
			},
			err: errors.New("(build::AddBuildArgs)", "Argument 'argument' already defined"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.dockerBuildCmd.AddBuildArgs(test.args.arg, test.args.value)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				_, exists := test.dockerBuildCmd.ImageBuildOptions.BuildArgs[test.args.arg]
				assert.True(t, exists, "Argument does not exists")
			}
		})
	}
}

func TestAddBuildContext(t *testing.T) {
	tests := []struct {
		desc              string
		dockerBuildCmd    *DockerBuildCmd
		buildContext      buildcontext.DockerBuildContexter
		prepareAssertFunc func(*mock.DockerBuildContext)
		err               error
	}{
		{
			desc:           "Testing error when build context is not defined",
			dockerBuildCmd: &DockerBuildCmd{},
			buildContext:   nil,
			err:            errors.New("(build:.AddBuilderContext)", "Docker build context is not defined"),
		},
		{
			desc:           "Testing add docker build context",
			dockerBuildCmd: &DockerBuildCmd{},
			buildContext:   mock.NewDockerBuildContext(),
			prepareAssertFunc: func(mock *mock.DockerBuildContext) {
				mock.On("Reader").Return(ioutil.NopCloser(io.Reader(&bytes.Buffer{})), nil)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				mock := new(mock.DockerBuildContext)
				test.prepareAssertFunc(mock)
				test.buildContext = mock
			}

			err := test.dockerBuildCmd.AddBuildContext(test.buildContext)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.NotNil(t, test.dockerBuildCmd.ImageBuildOptions.Context)
			}
		})
	}
}

func TestAddTags(t *testing.T) {

	type args struct {
		tag string
	}

	tests := []struct {
		desc           string
		dockerBuildCmd *DockerBuildCmd
		args           *args
		res            []string
		err            error
	}{
		{
			desc:           "Test add new tag",
			dockerBuildCmd: &DockerBuildCmd{},
			args: &args{
				tag: "new_tag",
			},
			res: []string{"new_tag"},
			err: nil,
		},
		{
			desc: "Test add new tag with normalized named",
			dockerBuildCmd: &DockerBuildCmd{
				UseNormalizedNamed: true,
			},
			args: &args{
				tag: "new_tag",
			},
			res: []string{"docker.io/library/new_tag"},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			test.dockerBuildCmd.AddTags(test.args.tag)
			assert.Equal(t, test.res, test.dockerBuildCmd.ImageBuildOptions.Tags)
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		desc              string
		dockerBuildCmd    *DockerBuildCmd
		prepareAssertFunc func(context.Context, *mockclient.DockerClient, *DockerBuildCmd)
		assertFunc        func(*mockclient.DockerClient) bool
		err               error
	}{
		{
			desc:           "Testing error when build with not defined DockerBuildCmd",
			dockerBuildCmd: nil,
			err:            errors.New("(build::Run)", "DockerBuildCmd is not defined"),
		},
		{
			desc:           "Testing error when build with not defined ImageBuildOptions",
			dockerBuildCmd: &DockerBuildCmd{},
			err:            errors.New("(build::Run)", "ImageBuildOptions options is not defined"),
		},
		{
			desc: "Testing error when build with not defined ImagePushOptions",
			dockerBuildCmd: &DockerBuildCmd{
				ImageBuildOptions: &dockertypes.ImageBuildOptions{},
				PushAfterBuild:    true,
			},
			err: errors.New("(build::Run)", "ImagePushOptions options is not defined"),
		},
		{
			desc: "Testing error when build with not defined docker build context",
			dockerBuildCmd: &DockerBuildCmd{
				ImageBuildOptions: &dockertypes.ImageBuildOptions{},
			},
			err: errors.New("(build::Run)", "Docker build context is not defined"),
		},

		{
			desc: "Testing build an image",
			dockerBuildCmd: &DockerBuildCmd{
				ImageName: "testing_image",
				ImageBuildOptions: &dockertypes.ImageBuildOptions{
					Context: ioutil.NopCloser(io.Reader(&bytes.Buffer{})),
				},
			},
			err: &errors.Error{},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerBuildCmd) {
				o := dockertypes.ImageBuildOptions{
					Tags:       []string{cmd.ImageName},
					Dockerfile: DefaultDockerfile,
					Context:    ioutil.NopCloser(io.Reader(&bytes.Buffer{})),
				}
				mock.On("ImageBuild", ctx, cmd.ImageBuildOptions.Context, o).Return(
					dockertypes.ImageBuildResponse{
						Body: ioutil.NopCloser(io.Reader(&bytes.Buffer{})),
					}, nil)
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImageBuild", 1)
			},
		},
		{
			desc: "Testing build and push an image",
			dockerBuildCmd: &DockerBuildCmd{
				ImageName: "testing_image",
				ImageBuildOptions: &dockertypes.ImageBuildOptions{
					Context: ioutil.NopCloser(io.Reader(&bytes.Buffer{})),
				},
				ImagePushOptions: &dockertypes.ImagePushOptions{},
				PushAfterBuild:   true,
			},
			err: &errors.Error{},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerBuildCmd) {
				o := dockertypes.ImageBuildOptions{
					Tags:       []string{cmd.ImageName},
					Dockerfile: DefaultDockerfile,
					Context:    ioutil.NopCloser(io.Reader(&bytes.Buffer{})),
				}
				mock.On("ImageBuild", ctx, cmd.ImageBuildOptions.Context, o).Return(
					dockertypes.ImageBuildResponse{
						Body: ioutil.NopCloser(io.Reader(&bytes.Buffer{})),
					}, nil)

				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(ioutil.NopCloser(io.Reader(&bytes.Buffer{})), nil)
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImageBuild", 1) && mock.AssertNumberOfCalls(t, "ImagePush", 1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			mock := new(mockclient.DockerClient)
			ctx := context.TODO()

			if test.prepareAssertFunc != nil {
				test.dockerBuildCmd.Cli = mock
				test.prepareAssertFunc(ctx, mock, test.dockerBuildCmd)
			}

			err := test.dockerBuildCmd.Run(ctx)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.True(t, test.assertFunc(mock))
			}
		})
	}
}
