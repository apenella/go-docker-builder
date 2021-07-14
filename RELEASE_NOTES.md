# RELEASE NOTES

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