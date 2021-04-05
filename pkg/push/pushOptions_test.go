package push

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	mockclient "github.com/apenella/go-docker-builder/test/mock"
	"github.com/aws/aws-sdk-go/aws"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	var w bytes.Buffer
	writer := io.Writer(&w)
	reader := ioutil.NopCloser(io.Reader(&bytes.Buffer{}))

	tests := []struct {
		desc          string
		dockerPushCmd *DockerPushCmd
		mock          *mockclient.DockerClient
		pushOptions   dockertypes.ImagePushOptions
		ctx           context.Context
		prepareMock   func(context.Context, *mockclient.DockerClient, *DockerPushCmd, dockertypes.ImagePushOptions)
		expected      struct {
			Ctx     context.Context
			Name    string
			Options dockertypes.ImagePushOptions
		}
		err error
	}{

		{
			desc:          "Testing error when DockerPushCmd is nil",
			ctx:           context.TODO(),
			pushOptions:   dockertypes.ImagePushOptions{},
			dockerPushCmd: nil,
			err:           errors.New("(push::Run)", "DockerPushCmd is nil"),
		},
		{
			desc: "Testing push a single image",
			ctx:  context.TODO(),
			prepareMock: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd, options dockertypes.ImagePushOptions) {
				mock.On("ImagePush", ctx, cmd.DockerPushOptions.ImageName, options).Return(reader, nil)
				cmd.Cli = mock
			},
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				Writer: io.Writer(writer),
				DockerPushOptions: &DockerPushOptions{
					ImageName: "test_image",
				},
				ExecPrefix: "",
			},
			expected: struct {
				Ctx     context.Context
				Name    string
				Options dockertypes.ImagePushOptions
			}{
				Ctx:     context.TODO(),
				Name:    "test_image",
				Options: dockertypes.ImagePushOptions{},
			},
			mock: new(mockclient.DockerClient),
			err:  &errors.Error{},
		},
		{
			desc: "Testing push a single image with auth",
			ctx:  context.TODO(),
			prepareMock: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd, options dockertypes.ImagePushOptions) {
				if cmd.DockerPushOptions.RegistryAuth != nil {
					options.RegistryAuth = *cmd.DockerPushOptions.RegistryAuth
				}
				mock.On("ImagePush", ctx, cmd.DockerPushOptions.ImageName, options).Return(reader, nil)
				cmd.Cli = mock
			},
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				Writer: io.Writer(writer),
				DockerPushOptions: &DockerPushOptions{
					ImageName:    "test_image",
					RegistryAuth: aws.String("auth"),
				},
				ExecPrefix: "",
			},
			expected: struct {
				Ctx     context.Context
				Name    string
				Options dockertypes.ImagePushOptions
			}{
				Ctx:  context.TODO(),
				Name: "test_image",
				Options: dockertypes.ImagePushOptions{
					RegistryAuth: "auth",
				},
			},
			mock: new(mockclient.DockerClient),
			err:  &errors.Error{},
		},
		{
			desc: "Testing push a single image with tags",
			ctx:  context.TODO(),
			prepareMock: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd, options dockertypes.ImagePushOptions) {
				mock.On("ImagePush", ctx, cmd.DockerPushOptions.ImageName, options).Return(reader, nil)
				mock.On("ImagePush", ctx, "tag1", options).Return(reader, nil)
				mock.On("ImagePush", ctx, "tag2", options).Return(reader, nil)
				cmd.Cli = mock
			},
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				Writer: io.Writer(writer),
				DockerPushOptions: &DockerPushOptions{
					ImageName: "test_image",
					Tags:      []string{"tag1", "tag2"},
				},
				ExecPrefix: "",
			},
			expected: struct {
				Ctx     context.Context
				Name    string
				Options dockertypes.ImagePushOptions
			}{
				Ctx:     context.TODO(),
				Name:    "test_image",
				Options: dockertypes.ImagePushOptions{},
			},
			mock: new(mockclient.DockerClient),
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareMock != nil {
				test.prepareMock(test.ctx, test.mock, test.dockerPushCmd, test.pushOptions)
			}

			err := test.dockerPushCmd.Run(test.ctx)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.mock.AssertCalled(t, "ImagePush", test.expected.Ctx, test.expected.Name, test.expected.Options)
			}
		})
	}
}

func TestAddAuth(t *testing.T) {

	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		options *DockerPushOptions
		args    *args
		err     error
		res     string
	}{
		{
			name: "Test add user-password auth",
			options: &DockerPushOptions{
				ImageName:    "test image",
				RegistryAuth: nil,
			},
			args: &args{
				username: "user",
				password: "AqSwd3Fr",
			},
			err: nil,
			res: "eyJ1c2VybmFtZSI6InVzZXIiLCJwYXNzd29yZCI6IkFxU3dkM0ZyIn0=",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.options.AddAuth(test.args.username, test.args.password)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, *test.options.RegistryAuth, "Unexpected auth result")
			}
		})
	}
}
