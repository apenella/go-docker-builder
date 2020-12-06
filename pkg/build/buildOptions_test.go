package build

import (
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestAddBuildArgs(t *testing.T) {

	type args struct {
		arg   string
		value string
	}

	tests := []struct {
		name    string
		options *DockerBuildOptions
		args    *args
		err     error
	}{
		{
			name: "Test add argument to nil BuildArgs object",
			options: &DockerBuildOptions{
				ImageName:  "test image",
				Tags:       []string{},
				BuildArgs:  nil,
				Dockerfile: "",
			},
			args: &args{
				arg:   "argument",
				value: "value",
			},
			err: nil,
		},
		{
			name: "Test add an existing argument",
			options: &DockerBuildOptions{
				ImageName: "test image",
				Tags:      []string{},
				BuildArgs: map[string]*string{
					"argument": nil,
				},
				Dockerfile: "",
			},
			args: &args{
				arg:   "argument",
				value: "value",
			},
			err: errors.New("(builder::AddBuildArgs)", "Argument 'argument' already defined"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := test.options.AddBuildArgs(test.args.arg, test.args.value)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				_, exists := test.options.BuildArgs[test.args.arg]
				assert.True(t, exists, "Argument does not exists")
			}
		})
	}
}

func TestAddTags(t *testing.T) {

	type args struct {
		tag string
	}

	tests := []struct {
		name    string
		options *DockerBuildOptions
		args    *args
		err     error
	}{
		{
			name: "Test add argument to nil BuildArgs object",
			options: &DockerBuildOptions{
				ImageName:  "test image",
				Tags:       nil,
				BuildArgs:  nil,
				Dockerfile: "",
			},
			args: &args{
				tag: "argument",
			},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.options.AddTags(test.args.tag)
			it := 0
			found := false
			for it < len(test.options.Tags) && !found {

				if test.options.Tags[it] == test.args.tag {
					found = true
				}
				it++
			}

			assert.True(t, found, "Argument does not exists")
		})
	}
}

func TestAddAuth(t *testing.T) {

	type args struct {
		username string
		password string
		registry string
	}
	tests := []struct {
		name    string
		options *DockerBuildOptions
		args    *args
		err     error
		res     map[string]dockertypes.AuthConfig
	}{
		{
			name: "Test add user-password auth",
			options: &DockerBuildOptions{
				ImageName: "test image",
				Auth:      map[string]dockertypes.AuthConfig{},
			},
			args: &args{
				username: "user",
				password: "AqSwd3Fr",
				registry: "registry",
			},
			err: nil,
			res: map[string]dockertypes.AuthConfig{
				"registry": dockertypes.AuthConfig{
					Username: "user",
					Password: "AqSwd3Fr",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.options.AddAuth(test.args.username, test.args.password, test.args.registry)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, test.options.Auth, "Unexpected auth result")
			}
		})
	}
}
