package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	errors "github.com/apenella/go-common-utils/error"
	dockertypes "github.com/docker/docker/api/types"
)

// GenerateUserPasswordAuthConfig return an AuthConfig to identify to docker registry using base64 auth credentials
func GenerateAuthConfig(username, password string) (*dockertypes.AuthConfig, error) {

	if username == "" {
		return nil, errors.New("(auth::GenerateUserPasswordAuthConfig)", "Username must be provided")
	}

	if password == "" {
		return nil, errors.New("(auth::GenerateUserPasswordAuthConfig)", "Password must be provided")
	}

	authConfig := &dockertypes.AuthConfig{
		Auth: base64.URLEncoding.EncodeToString([]byte(strings.Join([]string{username, password}, ":"))),
	}

	return authConfig, nil
}

// GenerateUserPasswordAuthConfig return an AuthConfig to identify to docker registry using user-password credentials
func GenerateUserPasswordAuthConfig(username, password string) (*dockertypes.AuthConfig, error) {

	if username == "" {
		return nil, errors.New("(auth::GenerateUserPasswordAuthConfig)", "Username must be provided")
	}

	if password == "" {
		return nil, errors.New("(auth::GenerateUserPasswordAuthConfig)", "Password must be provided")
	}

	authConfig := &dockertypes.AuthConfig{
		Username: username,
		Password: password,
	}

	return authConfig, nil
}

// GenerateEncodedUserPasswordAuthConfig return a pointer to an encoded docker registry auth string
func GenerateEncodedUserPasswordAuthConfig(username, password string) (*string, error) {

	var err error
	var auth *dockertypes.AuthConfig
	var encodedAuth *string

	auth, err = GenerateUserPasswordAuthConfig(username, password)
	if err != nil {
		return nil, errors.New("(auth::GenerateEncodedUserPasswordAuthConfig)", "Error generation user password auth configuration", err)
	}

	encodedAuth, err = encodeAuthConfig(auth)
	if err != nil {
		return nil, errors.New("(auth::GenerateEncodedUserPasswordAuthConfig)", "Error encoding auth configuration", err)
	}

	return encodedAuth, nil
}

func encodeAuthConfig(authConfig *dockertypes.AuthConfig) (*string, error) {

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return nil, errors.New("(auth::encodeAuthConfig)", "Error marshaling auth configuration", err)
	}
	authString := base64.URLEncoding.EncodeToString(encodedJSON)

	return &authString, nil
}
