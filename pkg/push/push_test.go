package push

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	mockclient "github.com/apenella/go-docker-builder/internal/mock"
	"github.com/apenella/go-docker-builder/pkg/response"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestWithTags(t *testing.T) {
	t.Log("Testing WithTags")
	p := DockerPushCmd{}
	p.WithTags([]string{"tag1", "tag2"})

	assert.Equal(t, []string{"tag1", "tag2"}, p.Tags)
}
func TestWithResponse(t *testing.T) {
	t.Log("Testing WithResponse")
	p := DockerPushCmd{}

	res := response.NewDefaultResponse()

	p.WithResponse(res)

	assert.Equal(t, res, p.Response)
}
func TestWithRemoveAfterPush(t *testing.T) {
	t.Log("Testing WithRemoveAfterPush")
	p := DockerPushCmd{}
	p.WithRemoveAfterPush()

	assert.Equal(t, true, p.RemoveAfterPush)
}
func TestWithUseNormalizedNamed(t *testing.T) {
	t.Log("Testing WithUseNormalizedNamed")
	p := DockerPushCmd{}
	p.WithUseNormalizedNamed()

	assert.Equal(t, true, p.UseNormalizedNamed)
}

func TestAddAuth(t *testing.T) {

	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args *args
		err  error
		res  string
	}{
		{
			name: "Test add user-password auth",
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
			t.Log(test.name)
			cmd := &DockerPushCmd{}
			err := cmd.AddAuth(test.args.username, test.args.password)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, cmd.ImagePushOptions.RegistryAuth, "Unexpected auth result")
			}
		})
	}
}

func TestAddTags(t *testing.T) {

	type args struct {
		tag []string
	}

	tests := []struct {
		desc          string
		dockerPushCmd *DockerPushCmd
		args          *args
		res           []string
		err           error
	}{
		{
			desc:          "Test add new tag",
			dockerPushCmd: &DockerPushCmd{},
			args: &args{
				tag: []string{"new_tag", "other_new_tag"},
			},
			res: []string{"new_tag", "other_new_tag"},
			err: nil,
		},
		{
			desc: "Test add new tag with normalized named",
			dockerPushCmd: &DockerPushCmd{
				UseNormalizedNamed: true,
			},
			args: &args{
				tag: []string{"new_tag", "other_new_tag"},
			},
			res: []string{"docker.io/library/new_tag", "docker.io/library/other_new_tag"},
			err: nil,
		},
		{
			desc: "Test add a tag that already exist",
			dockerPushCmd: &DockerPushCmd{
				Tags: []string{"new_tag"},
			},
			args: &args{
				tag: []string{"new_tag"},
			},
			res: []string{"new_tag"},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			test.dockerPushCmd.AddTags(test.args.tag...)

			assert.Equal(t, test.res, test.dockerPushCmd.Tags)
		})
	}
}

func TestRun(t *testing.T) {

	reader := ioutil.NopCloser(io.Reader(&bytes.Buffer{}))

	tests := []struct {
		desc              string
		dockerPushCmd     *DockerPushCmd
		pushOptions       dockertypes.ImagePushOptions
		prepareAssertFunc func(context.Context, *mockclient.DockerClient, *DockerPushCmd)
		assertFunc        func(*mockclient.DockerClient) bool
		err               error
	}{
		{
			desc:          "Testing error when DockerPushCmd is undefined",
			pushOptions:   dockertypes.ImagePushOptions{},
			dockerPushCmd: nil,
			err:           errors.New("(push::Run)", "DockerPushCmd is undefined"),
		},
		{
			desc:        "Testing error when ImageName is undefined",
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				ImagePushOptions: nil,
			},
			err: errors.New("(push::Run)", "Image name is undefined"),
		},
		{
			desc:        "Testing error when ImagePushOptions is undefined",
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				ImagePushOptions: nil,
				ImageName:        "name",
			},
			err: errors.New("(push::Run)", "Image push options is undefined"),
		},
		{
			desc: "Testing push a single image",
			dockerPushCmd: &DockerPushCmd{
				ImageName:        "test_image",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
			},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd) {
				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImagePush", 1)
			},
			err: &errors.Error{},
		},

		{
			desc: "Testing push a single image with remove after push",
			dockerPushCmd: &DockerPushCmd{
				ImageName:        "test_image",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
				RemoveAfterPush:  true,
			},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd) {
				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(reader, nil)
				mock.On("ImageRemove", ctx, cmd.ImageName, dockertypes.ImageRemoveOptions{
					Force:         true,
					PruneChildren: true,
				}).Return([]dockertypes.ImageDeleteResponseItem{}, nil)
				cmd.Cli = mock
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImagePush", 1) && mock.AssertNumberOfCalls(t, "ImageRemove", 1)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing push a single image with auth",
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd) {
				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				ImageName: "test_image",
				ImagePushOptions: &dockertypes.ImagePushOptions{
					RegistryAuth: "auth",
				},
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImagePush", 1)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing push a single image with tags",
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd) {
				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(reader, nil)
				mock.On("ImagePush", ctx, "tag1", *cmd.ImagePushOptions).Return(reader, nil)
				mock.On("ImagePush", ctx, "tag2", *cmd.ImagePushOptions).Return(reader, nil)
				mock.On("ImageTag", ctx, cmd.ImageName, "tag1").Return(nil)
				mock.On("ImageTag", ctx, cmd.ImageName, "tag2").Return(nil)
				cmd.Cli = mock
			},
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				ImageName:        "test_image",
				Tags:             []string{"tag1", "tag2"},
				ImagePushOptions: &dockertypes.ImagePushOptions{},
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImagePush", 3)
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			mock := new(mockclient.DockerClient)
			ctx := context.TODO()

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(ctx, mock, test.dockerPushCmd)
			}

			err := test.dockerPushCmd.Run(ctx)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.True(t, test.assertFunc(mock))
			}
		})
	}
}
