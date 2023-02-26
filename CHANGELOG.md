# CHANGELOG

## Undefined
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
- In docker build command, before adding a new tag validate that is normalized 

## [0.7.1]

### Fixed
- Existing labels could not be overwritten

## [0.7.0]

### Added
- Include `WithPullParentImage` method to set wether to pull the parent image on `build` instances
- Include constructor on git auth methods

### Changed
- Use go 1.19
- Use docker compose v2 to start testing stack
- Update testing docker images version

## [0.6.0]

### Added
- Include `WithImageName` method to set image name attribute on `build` instances
- Include `WithSourceImage` method to set source image attriubte on `copy` instances
- Include `WithTargetImage` method to set target image attriubte on `copy` instances
- Include `WithImageName` method to set image name attribute on `push` instances

### Changed
- **BREAKING CHANGES**: On package `build`, `NewDockerBuildCmd` has changed its signature to `NewDockerBuildCmd(cli types.DockerClienter) *DockerBuildCmd`
- **BREAKING CHANGES**: On package `copy`, `NewDockerImageCopyCmd` has changed its signature to `NewDockerImageCopyCmd(cli types.DockerClienter) *DockerImageCopyCmd`
- **BREAKING CHANGES**: On pacakge `push`, `NewDockerPushCmd` has changed its signature to `NewDockerPushCmd(cli types.DockerClienter) *DockerPushCmd`

## [0.5.0]

### Added
- Include constructors on `build`, `push` and `copy` packages
- Include `WithXXX`methods to set attributes on `build`, `push` and `copy` instances
- On `build` package new method to add labels

### Changed
- **BREAKING CHANGES**: On package `copy`, `AddTag` method has been renamed to `AddTags`
- **BREAKING CHANGES**: On pacakge `push`, `AddTag` method has been renamed to `AddTags`

### Fixed
- On `push` packages, tag images defined on `Tags` attribute before push them

## [0.4.0]

### Added
- Image copy package
- New intermediate filesystem to manage the build docker context. Build context filesystems let you to join multiple context before start an image build and tar itself.
- Included examples for new use cases
- Included resources for run examples and tests
- Included functinal test
-

### Changed
- Package for moking has been moved to internal/mock
- Response writer to manage push all pull responses
- git context support to create a build context from a repository subfolder

### Removed
- remove common package
