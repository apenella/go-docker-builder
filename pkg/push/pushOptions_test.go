package push

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
