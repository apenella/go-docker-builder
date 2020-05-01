package auth

import "github.com/go-git/go-git/v5/plumbing/transport"

type GitAuther interface {
	Auth() (transport.AuthMethod, error)
}
