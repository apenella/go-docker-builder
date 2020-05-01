package auth

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type BasicAuth struct {
	Username string
	Password string
}

func (a *BasicAuth) Auth() transport.AuthMethod {
	return &http.BasicAuth{
		Username: a.Username,
		Password: a.Password,
	}
}
