# CHANGELOG

## [0.8.2] - 2024-01-03

### Changed

- Bumps github.com/go-git/go-git/v5 from 5.10.0 to 5.11.0.

## [0.8.1] - 2023-20-12

### Changed

- Bump golang.org/x/crypto from 0.14.0 to 0.17.0

## [0.8.0] - 2023-27-01

### Fixed

- Fixed a bug that caused the library to return an error when the build context contains a symlink.

### Changed

- In the examples, start the client container once and execute actions on it.
- Bump github.com/docker/docker from v20.10+incompatible to v24.0.7+incompatible
- Bump github.com/docker/distribution from v2.8.2+incompatible to v2.8.3+incompatible
- Bump github.com/go-git/go-git/v5 from v5.6.1 to v5.10.0
- Bump github.com/spf13/afero from v1.9.5 to v1.10.0
- Bump github.com/stretchr/testify from v1.8.2 to v1.8.4

## [0.7.8]

### Fixed

- On the build use case, WithPullParentImage is set to true ImageBuildOptions.PullParent
- On the copy use case, set the remove after the push parameter

## [0.7.7]

### Fixed

- Bump github.com/docker/distribution from 2.8.1+incompatible to 2.8.2+incompatible
- Bump github.com/cloudflare/circl from 1.3.2 to 1.3.3 (indirect dependency)

## [0.7.6]

### Fixed

- Bump up github.com/docker/docker to 20.10.24 to fix CVE-2023-28841, CVE-2023-28840, and CVE-2023-28842

## [0.7.5]

### Fixed

- Revert github.com/docker/docker to v20.10.23+incompatible to fix compatibility with Docker API 1.41

## [0.7.4]

### Changed

- Bump up golang.org/x/text to v0.7.0. It fixes [CVE-2022-32149](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-32149)
thub.com/apenella/go-common-utils/error v0.0.0-20221227202648-5452d804e940
- Bump up github.com/apenella/go-common-utils/transformer/string to v0.0.0-20221227202648-5452d804e940
- Bump up github.com/docker/docker to v23.0.1+incompatible
- Bump up github.com/go-git/go-git/v5 to v5.5.2
- Bump up github.com/spf13/afero to v1.9.4
- Bump up github.com/stretchr/testify to v1.8.2
- Bump up github.com/xanzy/ssh-agent to v0.3.3

## [0.7.3]

### Fixed

- Generate auth config when username and password are empty, instead of returning an error

## [0.7.2]

### Fixed

- On the build command, validate that the tag is normalized before adding it to the image name

## [0.7.1]

### Fixed

- Existing labels could not be overwritten

## [0.7.0]

### Added

- Include `WithPullParentImage` method to set whether to pull the parent image on `build` instances
- Include constructor on git auth methods

### Changed

- Use go 1.19
- Use Docker Compose v2 to start testing the stack
- Update testing docker images version

## [0.6.0]

### Added

- Include `WithImageName` method to set the image name attribute on `build` instances
- Include `WithSourceImage` method to set the source image attribute on `copy` instances
- Include `WithTargetImage` method to set the target image attribute on `copy` instances
- Include `WithImageName` method to set the image name attribute on `push` instances

### Changed

- **BREAKING CHANGES**: On package `build`, `NewDockerBuildCmd` has changed its signature to `NewDockerBuildCmd(cli types.DockerClienter) *DockerBuildCmd`
- **BREAKING CHANGES**: On package `copy`, `NewDockerImageCopyCmd` has changed its signature to `NewDockerImageCopyCmd(cli types.DockerClienter) *DockerImageCopyCmd`
- **BREAKING CHANGES**: On package `push`, `NewDockerPushCmd` has changed its signature to `NewDockerPushCmd(cli types.DockerClienter) *DockerPushCmd`

## [0.5.0]

### Added

- Include constructors on `build`, `push` and `copy` packages
- Include `WithXXX` methods to set attributes on the `build`, `push` and `copy` instances
- On the `build` package new method to add labels

### Changed

- **BREAKING CHANGES**: On package `copy`, `AddTag` method has been renamed to `AddTags`
- **BREAKING CHANGES**: On package `push`, `AddTag` method has been renamed to `AddTags`

### Fixed

- On `push` packages, tag images defined on `Tags` attribute before pushing them

## [0.4.0]

### Added

- Image copy package
- New intermediate filesystem to manage the build docker context. Build context filesystems let you join multiple contexts before starting an image build and tar itself.
- Included examples for new use cases
- Included resources for running examples and tests
- Included functional test

### Changed

- The package for mocking has been moved to internal/mock
- Response writer to manage to push all pull responses
- Git context support to create a build context from a repository subfolder

### Removed

- Remove common package
