# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-02-18

### Added
- Support for C SDK provider (MinGW-w64)
- Support for C++ SDK provider (MinGW-w64)
- MinGW-w64 GCC versions: 15.2.0, 14.2.0, 14.1.0, 13.2.0, 13.1.0, 12.3.0, 12.2.0, 12.1.0, 11.3.0, 11.2.0, 11.1.0
- Includes GCC (C compiler), G++ (C++ compiler), and related development tools
- Automatic PATH configuration for MinGW bin directory
- Use `unosdk install c mingw <version>` for C or `unosdk install cpp mingw <version>` for C++

## [1.1.0] - 2026-02-06

### Added
- Support for Maven SDK provider (Apache Maven)
- Maven versions: 3.9.9, 3.9.8, 3.9.7, 3.9.6, 3.8.8, 3.8.7, 3.8.6, 3.6.3
- Support for Gradle SDK provider
- Gradle versions: 8.12, 8.11.1, 8.11, 8.10.2, 8.10.1, 8.10, 8.9, 8.8, 8.7, 8.6, 8.5, 8.4, 8.3, 8.2.1, 8.2, 8.1.1, 8.1, 8.0.2, 8.0.1, 8.0, 7.6.4, 7.6.3, 7.6.2, 7.6.1, 7.6
- Support for Go programming language provider
- Go versions: 1.23.5, 1.23.4, 1.23.3, 1.23.2, 1.23.1, 1.23.0, 1.22.x series, 1.21.x series
- Support for Flutter SDK provider
- Flutter SDK versions: 3.27.2, 3.27.1, 3.24.5, 3.22.3, 3.19.6, 3.16.9, 3.13.9
- Support for GraalVM Java distribution provider
- GraalVM versions with simplified Java major version aliases (25, 21, 17)
- GraalVM specific versions: 25.0.2, 21.0.2, 17.0.9

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

[1.2.0]: https://github.com/javaquery/unosdk/releases/tag/v1.2.0
[1.1.0]: https://github.com/javaquery/unosdk/releases/tag/v1.1.0
[1.0.0]: https://github.com/javaquery/unosdk/releases/tag/v1.0.0
