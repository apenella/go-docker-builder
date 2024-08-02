package auth

import (
	"testing"

	dockerregistrytypes "github.com/docker/docker/api/types/registry"
	"github.com/stretchr/testify/assert"
)

func TestEncodeAuthConfig(t *testing.T) {
	type args struct {
		authConfig *dockerregistrytypes.AuthConfig
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
