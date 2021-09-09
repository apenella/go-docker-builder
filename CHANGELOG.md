# CHANGELOG

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
