# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Support for Flutter SDK provider
- Flutter SDK versions: 3.27.2, 3.27.1, 3.24.5, 3.22.3, 3.19.6, 3.16.9, 3.13.9
- Support for "latest" version alias for Flutter installations
- Flutter downloads from official Google Cloud Storage (flutter_infra_release)
- Support for GraalVM Java distribution provider
- GraalVM versions with simplified Java major version aliases (25, 21, 17)
- GraalVM specific versions: 25.0.2, 21.0.2, 17.0.9

### Changed
- Updated GraalVM download URL format to match official release naming convention
- GraalVM now supports simplified version syntax (e.g., `unosdk install java graalvm 21`)
- Updated CLI examples to include Flutter SDK installation

## [1.0.0] - 2026-02-05

### Added
- Initial release of UnoSDK
- CLI tool for managing multiple SDK installations
- Support for Java SDK providers (Amazon Corretto, OpenJDK)
- Support for Node.js SDK provider
- Support for Python SDK provider
- `install` command to install SDK versions
- `list` command to list available and installed SDKs
- `switch` command to switch between SDK versions
- `uninstall` command to remove installed SDKs
- `update` command to update SDK registry
- Environment setup and PATH management
- SDK registry with caching support
- Download verification with checksums
- Progress indicators for downloads and installations
- Configuration management via `sdks.yaml`
- Cross-platform support (Windows, macOS, Linux)

[1.0.0]: https://github.com/javaquery/unosdk/releases/tag/v1.0.0
