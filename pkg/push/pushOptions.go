package push

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	dockertypes "github.com/docker/docker/api/types"
)

// DockerBuilderOptions has an options set to build and image
type DockerPushOptions struct {
	// ImageName is the name of the image
	ImageName string
	// // RegistryAuth is the base64 encoded credentials for the registry
	RegistryAuth *string
}

func (o *DockerPushOptions) AddUserPasswordRegistryAuth(username, password, serverAddress string) error {

	var err error

	if username == "" {
		return errors.New("(push::AddUserPasswordRegistryAuth) Username must be provided")
	}

	if password == "" {
		return errors.New("(push::AddUserPasswordRegistryAuth) Password must be provided")
	}

	authConfig := &dockertypes.AuthConfig{
		Username: username,
		Password: password,
	}

	if serverAddress != "" {
		authConfig.ServerAddress = serverAddress
	}

	o.RegistryAuth, err = encodeAuthConfig(authConfig)
	if err != nil {
		return errors.New("(push::AddUserPasswordRegistryAuth) Error encoding authorization. " + err.Error())
	}

	return nil
}

func encodeAuthConfig(authConfig *dockertypes.AuthConfig) (*string, error) {

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, errors.New("(push::encodeAuthConfig) Error marshaling auth configuration. " + err.Error())
	}
	authString := base64.URLEncoding.EncodeToString(encodedJSON)

	return &authString, nil
}
