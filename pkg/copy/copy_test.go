package copy

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

	c := DockerImageCopyCmd{}
	c.WithTags([]string{"tag1", "tag2"})
	assert.Equal(t, []string{"tag1", "tag2"}, c.Tags)
}
func TestWithRemoteSource(t *testing.T) {
	t.Log("Testing WithRemoteSource")

	c := DockerImageCopyCmd{}
	c.WithRemoteSource()
	assert.Equal(t, true, c.RemoteSource)
}
func TestWithRemoveAfterPush(t *testing.T) {
	t.Log("Testing WithRemoveAfterPush")

	c := DockerImageCopyCmd{}
	c.WithRemoveAfterPush()
	assert.Equal(t, true, c.RemoveAfterPush)
}
func TestWithResponse(t *testing.T) {
	t.Log("Testing WithResponse")

	res := response.NewDefaultResponse()

	c := DockerImageCopyCmd{}
	c.WithResponse(res)
	assert.Equal(t, res, c.Response)
}
func TestWithUseNormalizedNamed(t *testing.T) {
	t.Log("Testing WithUseNormalizedNamed")

	c := DockerImageCopyCmd{}
	c.WithUseNormalizedNamed()
	assert.Equal(t, true, c.UseNormalizedNamed)
}

func TestAddPullAuth(t *testing.T) {

	type args struct {
		username string
		password string
	}
	tests := []struct {
		desc string
		args *args
		err  error
		res  string
	}{
		{
			desc: "Test add user-password auth",
			args: &args{
				username: "user",
				password: "AqSwd3Fr",
			},
			err: nil,
			res: "eyJ1c2VybmFtZSI6InVzZXIiLCJwYXNzd29yZCI6IkFxU3dkM0ZyIn0=",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			cmd := &DockerImageCopyCmd{}
			err := cmd.AddPullAuth(test.args.username, test.args.password)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, cmd.ImagePullOptions.RegistryAuth, "Unexpected auth result")
			}
		})
	}
}

func TestAddPushAuth(t *testing.T) {

	type args struct {
		username string
		password string
	}
	tests := []struct {
		desc string
		args *args
		err  error
		res  string
	}{
		{
			desc: "Test add user-password auth",
			args: &args{
				username: "user",
				password: "AqSwd3Fr",
			},
			err: nil,
			res: "eyJ1c2VybmFtZSI6InVzZXIiLCJwYXNzd29yZCI6IkFxU3dkM0ZyIn0=",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			cmd := &DockerImageCopyCmd{}
			err := cmd.AddPushAuth(test.args.username, test.args.password)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, cmd.ImagePushOptions.RegistryAuth, "Unexpected auth result")
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
		desc string
		args *args
		err  error
		res  string
	}{
		{
			desc: "Test add user-password auth",
			args: &args{
				username: "user",
				password: "AqSwd3Fr",
			},
			err: nil,
			res: "eyJ1c2VybmFtZSI6InVzZXIiLCJwYXNzd29yZCI6IkFxU3dkM0ZyIn0=",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			cmd := &DockerImageCopyCmd{}
			err := cmd.AddAuth(test.args.username, test.args.password)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, cmd.ImagePushOptions.RegistryAuth, "Unexpected auth result")
				assert.Equal(t, test.res, cmd.ImagePullOptions.RegistryAuth, "Unexpected auth result")
			}
		})
	}
}

func TestAddTags(t *testing.T) {

	tests := []struct {
		desc string
		tags []string
		res  []string
	}{
		{
			desc: "Testing add tag",
			tags: []string{"tag1"},
			res:  []string{"tag1"},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			cmd := &DockerImageCopyCmd{}
			cmd.AddTags(test.tags...)

			assert.Equal(t, len(test.res), len(cmd.Tags), "Unexpected auth result")
		})
	}
}

func TestRun(t *testing.T) {

	//	var w bytes.Buffer
	//	writer := io.Writer(&w)
	reader := ioutil.NopCloser(io.Reader(&bytes.Buffer{}))

	tests := []struct {
		desc               string
		dockerImageCopyCmd *DockerImageCopyCmd
		prepareAssertFunc  func(context.Context, *mockclient.DockerClient, *DockerImageCopyCmd)
		assertFunc         func(*mockclient.DockerClient) bool
		err                error
	}{
		{
			desc:               "Testing error when DockerImageCopyCmd is undefined",
			dockerImageCopyCmd: nil,
			err:                errors.New("(copy::Run)", "DockerImageCopyCmd is undefined"),
		},
		{
			desc:               "Testing error when source image is undefined",
			dockerImageCopyCmd: &DockerImageCopyCmd{},
			err:                errors.New("(copy::Run)", "Source image must be defined"),
		},
		{
			desc: "Testing error when target image is undefined",
			dockerImageCopyCmd: &DockerImageCopyCmd{
				SourceImage: "source",
			},
			err: errors.New("(copy::Run)", "Target image must be defined"),
		},
		{
			desc: "Testing error when image push options is undefined",
			dockerImageCopyCmd: &DockerImageCopyCmd{
				SourceImage: "source",
				TargetImage: "target",
			},
			err: errors.New("(copy::Run)", "Image push options is undefined"),
		},
		{
			desc: "Testing local image copy",
			dockerImageCopyCmd: &DockerImageCopyCmd{
				SourceImage:      "source",
				TargetImage:      "target",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
			},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerImageCopyCmd) {
				mock.On("ImageTag", ctx, cmd.SourceImage, cmd.TargetImage).Return(nil)
				mock.On("ImagePush", ctx, cmd.TargetImage, *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImageTag", 1)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing error when tagging source image",
			dockerImageCopyCmd: &DockerImageCopyCmd{
				SourceImage:      "source",
				TargetImage:      "target",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
			},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerImageCopyCmd) {
				mock.On("ImageTag", ctx, cmd.SourceImage, cmd.TargetImage).Return(errors.New("(test)", "Error tagging"))
				//		mock.On("ImagePush", ctx, cmd.TargetImage, *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImageTag", 1)
			},
			err: errors.New("(copy::Run)", "Error tagging image 'source' to 'target'",
				errors.New("(test)", "Error tagging")),
		},
		{
			desc: "Testing pushing target image",
			dockerImageCopyCmd: &DockerImageCopyCmd{
				SourceImage:      "source",
				TargetImage:      "target",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
			},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerImageCopyCmd) {
				mock.On("ImageTag", ctx, cmd.SourceImage, cmd.TargetImage).Return(nil)
				mock.On("ImagePush", ctx, cmd.TargetImage, *cmd.ImagePushOptions).Return(reader, errors.New("(test)", "Error pushing image"))
				cmd.Cli = mock
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImageTag", 1)
			},
			err: errors.New("(copy::Run)", "Error pushing image",
				errors.New("(push::Run)", "Error pushing image 'target'",
					errors.New("(test)", "Error pushing image"))),
		},
		{
			desc: "Testing error copying remote image without pull options",
			dockerImageCopyCmd: &DockerImageCopyCmd{
				SourceImage:      "source",
				TargetImage:      "target",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
				RemoteSource:     true,
			},
			err: errors.New("(copy::Run)", "Image pull options is undefined"),
		},
		{
			desc: "Testing remote image copy",
			dockerImageCopyCmd: &DockerImageCopyCmd{
				SourceImage:      "source",
				TargetImage:      "target",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
				ImagePullOptions: &dockertypes.ImagePullOptions{},
				RemoteSource:     true,
			},
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerImageCopyCmd) {
				mock.On("ImagePull", ctx, cmd.SourceImage, *cmd.ImagePullOptions).Return(reader, nil)
				mock.On("ImageTag", ctx, cmd.SourceImage, cmd.TargetImage).Return(nil)
				mock.On("ImagePush", ctx, cmd.TargetImage, *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImageTag", 1)
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
				test.prepareAssertFunc(ctx, mock, test.dockerImageCopyCmd)
			}

			err := test.dockerImageCopyCmd.Run(ctx)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.True(t, test.assertFunc(mock))
			}
		})
	}
}
