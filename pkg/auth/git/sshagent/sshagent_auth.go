package auth

import (
	errors "github.com/apenella/go-common-utils/error"
	"github.com/go-git/go-git/v5/plumbing/transport"
	ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	sshagent "github.com/xanzy/ssh-agent"
)

type SSHAgentAuth struct {
	GitSSHUser string
}

// NewSSHAgentAuth returns a new SSHAgentAuth
func NewSSHAgentAuth(gitSSHUser string) *SSHAgentAuth {
	return &SSHAgentAuth{GitSSHUser: gitSSHUser}
}

// Auth returns a new transport.AuthMethod created from the SSHAgentAuth
func (a *SSHAgentAuth) Auth() (transport.AuthMethod, error) {

	if a.GitSSHUser == "" {
		a.GitSSHUser = "git"
	}

	agent, err := ssh.NewSSHAgentAuth(a.GitSSHUser)
	if err != nil || !sshagent.Available() {
		return nil, errors.New("(auth::SSHAgentAuth::Auth)", "Agent is not available", err)
	}

	return agent, nil
}
