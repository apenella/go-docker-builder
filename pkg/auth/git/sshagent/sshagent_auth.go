package auth

import (
	"errors"

	"github.com/go-git/go-git/v5/plumbing/transport"
	ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	sshagent "github.com/xanzy/ssh-agent"
)

type SSHAgentAuth struct {
	GitSSHUser string
}

func (a *SSHAgentAuth) Auth() (transport.AuthMethod, error) {

	if a.GitSSHUser == "" {
		a.GitSSHUser = "git"
	}

	agent, err := ssh.NewSSHAgentAuth(a.GitSSHUser)
	if err != nil || !sshagent.Available() {
		return nil, errors.New("(auth::SSHAgentAuth::Auth) Agent is not available.\n  " + err.Error())
	}

	return agent, nil
}
