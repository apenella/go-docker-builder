package auth

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type BasicAuth struct {
	Username string
	Password string
}

// NewBasicAuth returns a new BasicAuth with the given username and password
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

// Auth return a transport.AuthMethod created from BasicAuth
func (a *BasicAuth) Auth() (transport.AuthMethod, error) {
	return &http.BasicAuth{
		Username: a.Username,
		Password: a.Password,
	}, nil
}
