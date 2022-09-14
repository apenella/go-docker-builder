package auth

import (
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

// func TestAddUserPasswordRegistryAuth(t *testing.T) {

// 	type args struct {
// 		username      string
// 		password      string
// 		serverAddress string
// 	}
// 	tests := []struct {
// 		name    string
// 		options *DockerPushOptions
// 		args    *args
// 		err     error
// 		res     string
// 	}{
// 		{
// 			name: "Test add user-password auth",
// 			options: &DockerPushOptions{
// 				ImageName:    "test image",
// 				RegistryAuth: nil,
// 			},
// 			args: &args{
// 				username:      "user",
// 				password:      "AqSwd3Fr",
// 				serverAddress: "localhost",
// 			},
// 			err: nil,
// 			res: "eyJ1c2VybmFtZSI6InVzZXIiLCJwYXNzd29yZCI6IkFxU3dkM0ZyIiwic2VydmVyYWRkcmVzcyI6ImxvY2FsaG9zdCJ9",
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			err := test.options.AddUserPasswordRegistryAuth(test.args.username, test.args.password, test.args.serverAddress)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				assert.Equal(t, test.res, *test.options.RegistryAuth, "Unexpected auth result")
// 			}
// 		})
// 	}
// }

func TestEncodeAuthConfig(t *testing.T) {
	type args struct {
		authConfig *dockertypes.AuthConfig
	}
	tests := []struct {
		name string
		args args
		err  error
		res  string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			auth, err := encodeAuthConfig(test.args.authConfig)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.res, *auth, "Unexpected auth result")
			}
		})
	}
}
