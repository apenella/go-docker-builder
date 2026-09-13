# RELEASE NOTES

## [v0.13.0] (2026-09-12)

### Changed

- Migrate from the deprecated `github.com/docker/docker v28.5.2+incompatible` to `github.com/moby/moby/client v0.6.0` and `github.com/moby/moby/api v1.56.0`
- Bump `github.com/docker/go-connections` from `v0.5.0` to `v0.8.1` (required by the new moby client)
- Bump `github.com/go-git/go-git/v5` from `v5.19.1` to `v5.19.2`
- Bump `github.com/stretchr/testify` from `v1.11.1` to `v1.12.1`
- Bump `dario.cat/mergo` from `v1.0.1` to `v1.0.2`
- Bump `github.com/ProtonMail/go-crypto` from `v1.1.6` to `v1.4.1`
- Bump `github.com/cloudflare/circl` from `v1.6.3` to `v1.6.5`
- Bump `github.com/cyphar/filepath-securejoin` from `v0.6.1` to `v0.7.0`
- Bump `github.com/felixge/httpsnoop` from `v1.0.4` to `v1.1.0`
- Bump `github.com/go-git/go-billy/v5` from `v5.9.0` to `v5.9.1`
- Bump `github.com/go-logr/logr` from `v1.4.3` to `v1.4.4`
- Bump `github.com/kevinburke/ssh_config` from `v1.2.0` to `v1.6.0`
- Bump `github.com/klauspost/cpuid/v2` from `v2.3.0` to `v2.4.0`
- Bump `github.com/sergi/go-diff` from `v1.3.2-0.20230802210424-5b0b94c5c0d3` to `v1.4.0`
- Bump `github.com/skeema/knownhosts` from `v1.3.1` to `v1.3.3`
- Bump `github.com/stretchr/objx` from `v0.5.2` to `v0.5.3`
- Bump `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` from `v0.60.0` to `v0.71.0`
- Bump `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/metric` and `go.opentelemetry.io/otel/trace` from `v1.43.0` to `v1.46.0`
- Bump `golang.org/x/crypto` from `v0.50.0` to `v0.54.0`
- Bump `golang.org/x/net` from `v0.53.0` to `v0.56.0`
- Bump `golang.org/x/sys` from `v0.43.0` to `v0.47.0`
- Bump `golang.org/x/text` from `v0.36.0` to `v0.40.0`
- Add `go.yaml.in/yaml/v3 v3.0.5`

### Breaking changes

- **BREAKING CHANGES**: `pkg/types.DockerClienter` method signatures follow the moby SDK: `ImageBuild` returns `client.ImageBuildResult`, `ImagePull` returns `client.ImagePullResponse`, `ImagePush` returns `client.ImagePushResponse`, `ImageRemove` returns `client.ImageRemoveResult`, and `ImageTag` takes `client.ImageTagOptions` and returns `client.ImageTagResult`
- **BREAKING CHANGES**: `ImageBuildOptions`, `ImagePushOptions`, `ImagePullOptions` and `ImageRemoveOptions` moved from `github.com/docker/docker/api/types*` to `github.com/moby/moby/client`. `AuthConfig` moved from `github.com/docker/docker/api/types/registry` to `github.com/moby/moby/api/types/registry`, and `DeleteResponse` moved to `github.com/moby/moby/api/types/image`
- **BREAKING CHANGES**: `ImageTag` now takes `client.ImageTagOptions{Source, Target}` instead of `(imageID, ref string)`, and `ImageRemove` returns `client.ImageRemoveResult` — iterate over `result.Items` instead of a `[]DeleteResponse` slice
- **BREAKING CHANGES**: consumers must create the client with `github.com/moby/moby/client.NewClientWithOpts(client.FromEnv)` instead of `github.com/docker/docker/client`. A `*docker/docker` client no longer satisfies `types.DockerClienter`, and custom `DockerClienter` implementations or mocks must be updated to the new signatures

## [v0.12.0] (2026-04-17)

### Changed

- Bumps Golang 1.25
- Bump Alpine, Golang and Python version used on examples and testing stacks
