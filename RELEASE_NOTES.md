# RELEASE NOTES

## [v0.13.0] (2026-09-12)

### Changed

- Migrate from the deprecated `github.com/docker/docker v28.5.2+incompatible` to `github.com/moby/moby/client v0.6.0` and `github.com/moby/moby/api v1.56.0`
- Bump `github.com/docker/go-connections` from `v0.5.0` to `v0.7.0` (required by the new moby client)

### Breaking changes

- **BREAKING CHANGES**: `pkg/types.DockerClienter` method signatures follow the moby SDK: `ImageBuild` returns `client.ImageBuildResult`, `ImagePull` returns `client.ImagePullResponse`, `ImagePush` returns `client.ImagePushResponse`, `ImageRemove` returns `client.ImageRemoveResult`, and `ImageTag` takes `client.ImageTagOptions` and returns `client.ImageTagResult`
- **BREAKING CHANGES**: `ImageBuildOptions`, `ImagePushOptions`, `ImagePullOptions` and `ImageRemoveOptions` moved from `github.com/docker/docker/api/types*` to `github.com/moby/moby/client`. `AuthConfig` moved from `github.com/docker/docker/api/types/registry` to `github.com/moby/moby/api/types/registry`, and `DeleteResponse` moved to `github.com/moby/moby/api/types/image`
- **BREAKING CHANGES**: `ImageTag` now takes `client.ImageTagOptions{Source, Target}` instead of `(imageID, ref string)`, and `ImageRemove` returns `client.ImageRemoveResult` — iterate over `result.Items` instead of a `[]DeleteResponse` slice
- **BREAKING CHANGES**: consumers must create the client with `github.com/moby/moby/client.NewClientWithOpts(client.FromEnv)` instead of `github.com/docker/docker/client`. A `*docker/docker` client no longer satisfies `types.DockerClienter`, and custom `DockerClienter` implementations or mocks must be updated to the new signatures

## [v0.12.0] (2026-04-17)

### Changed

- Bumps Golang 1.25
- Bump Alpine, Golang and Python version used on examples and testing stacks
