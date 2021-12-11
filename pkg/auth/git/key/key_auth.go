package auth

import (
	"fmt"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/go-git/go-git/v5/plumbing/transport"
	ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type KeyAuth struct {
	GitSSHUser string
	PkFile     string
	PkPassword string
}

// NewKeyAuth returns a new KeyAuth
func NewKeyAuth(gitSSHUser, pkFile, pkPassword string) *KeyAuth {
	return &KeyAuth{
		GitSSHUser: gitSSHUser,
		PkFile:     pkFile,
		PkPassword: pkPassword,
	}
}

// Auth returns a transport.AuthMethod created from the KeyAuth
func (a *KeyAuth) Auth() (transport.AuthMethod, error) {

	if a.GitSSHUser == "" {
		a.GitSSHUser = "git"
	}

	key, err := ssh.NewPublicKeysFromFile(a.GitSSHUser, a.PkFile, a.PkPassword)
	if err != nil {
		return nil, errors.New("(auth::SSHAgentAuth::Auth)", fmt.Sprintf("Could not load key from file '%s'", a.PkFile), err)
	}

	return key, nil
}
