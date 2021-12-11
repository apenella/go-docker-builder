# RELEASE NOTES

## [0.6.0]

### Added
- Include `WithImageName` method to set image name attribute on `build` instances
- Include `WithPullParentImage` method to set wether to pull the parent image on `build` instances
- Include `WithSourceImage` method to set source image attriubte on `copy` instances
- Include `WithTargetImage` method to set target image attriubte on `copy` instances
- Include `WithImageName` method to set image name attribute on `push` instances
- Include constructor on git auth methods

### Changed
- **BREAKING CHANGES**: On package `build`, `NewDockerBuildCmd` has changed its signature to `NewDockerBuildCmd(cli types.DockerClienter) *DockerBuildCmd`
- **BREAKING CHANGES**: On package `copy`, `NewDockerImageCopyCmd` has changed its signature to `NewDockerImageCopyCmd(cli types.DockerClienter) *DockerImageCopyCmd`
- **BREAKING CHANGES**: On pacakge `push`, `NewDockerPushCmd` has changed its signature to `NewDockerPushCmd(cli types.DockerClienter) *DockerPushCmd`
