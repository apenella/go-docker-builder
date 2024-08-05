# RELEASE NOTES

## [0.9.0] - 2024-08-05

- Upgraded to Go 1.22
- Bumps github.com/docker/docker from 24.0.9+incompatible to v27.1.1+incompatible. With this change, the Go client Docker Engine API types packages must be updated to the latest version.
- Bumps github.com/go-git/go-git/v5 from v5.11.0 to v5.12.0
- Bumps github.com/spf13/afero from v1.10.0 to v1.11.0
- Bumps golang.org/x/net from 0.19.0 to 0.23.0.
- In the release.yaml workflow, update the action actions/checkout from v3 to v4
- In the release.yaml workflow, update the action actions/setup-go from v3 to v5, and update the Go version from 1.19 to 1.22
- In the release.yaml workflow, update the action goreleaser/goreleaser-action from v6 to v6
- Replace github.com/docker/distribution from v2.8.3+incompatible to github.com/distribution/reference v0.6.0
- Update Docker-dind image used for testing from docker:28.0-dind to docker:27.1-dind
