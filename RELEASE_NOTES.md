# RELEASE NOTES

## [0.4.0]

### Added
- Image copy package
- New intermediate filesystem to manage the build docker context. Build context filesystems let you to join multiple context before start an image build and tar itself.
- Included examples for new use cases
- Included functinal test

### Changed
- Response writer to manage push all pull responses
- git context support to create a build context from a repository subfolder

### Removed
- remove common package
