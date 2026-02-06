# UnoSDK

<p align="center">
  <img src="logo.svg" alt="UnoSDK Logo" width="200"/>
</p>

**UnoSDK** is a powerful CLI tool for Windows that simplifies the installation and management of multiple software development kits (SDKs) from various providers. Think of it as **SDKMAN for Windows** - bringing the same ease of SDK management to Windows machines. Say goodbye to manual downloads, extractions, and environment variable configurations.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/javaquery/unosdk?style=flat&logo=github)](https://github.com/javaquery/unosdk/releases/latest)
[![CI](https://github.com/javaquery/unosdk/actions/workflows/ci.yml/badge.svg)](https://github.com/javaquery/unosdk/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Why UnoSDK?

If you've used [SDKMAN!](https://sdkman.io/) on Linux or macOS and wished for something similar on Windows, **UnoSDK** is your answer. It provides a native Windows experience for managing multiple SDK versions without the complexity of manual installation and PATH management.

## Features

- 🚀 **Multi-SDK Support**: Manage Java, Node.js, Python, and Flutter installations from a single tool
- 🔄 **Version Switching**: Easily switch between different SDK versions
- 📦 **Multiple Providers**: Support for various distribution providers
  - Java: Amazon Corretto, OpenJDK, GraalVM
  - Node.js: Official Node.js distributions
  - Python: Official Python distributions
  - Flutter: Official Flutter SDK
  - Maven: Apache Maven build tool
- 🔧 **Automatic Environment Setup**: Automatically configures PATH and environment variables
- 📋 **Registry Management**: Keeps track of all installed SDKs
- ⚡ **Fast Downloads**: Parallel downloads with progress tracking
- 🛡️ **Verification**: Ensures download integrity with checksum verification

## Supported SDKs

| SDK Type | Providers | Description |
|----------|-----------|-------------|
| Java | Amazon Corretto, OpenJDK, GraalVM | Java Development Kit |
| Node.js | nodejs | JavaScript runtime environment |
| Python | python | Python programming language |
| Flutter | flutter | Flutter SDK for mobile, web, and desktop apps |
| Maven | apache | Apache Maven build automation tool |

## Installation

### Prerequisites

- Windows OS (Windows 10 or later)
- PowerShell 5.1 or later

### Quick Installation

**Automated Installation (Recommended):**

Open PowerShell and run:

```powershell
irm https://raw.githubusercontent.com/javaquery/unosdk/refs/heads/master/scripts/install.ps1 | iex
```

This will automatically:
- Download the latest release from GitHub
- Install to `%LOCALAPPDATA%\unosdk`
- Add unosdk to your PATH
- Replace existing installation if present

**To reinstall/update unosdk:**

Simply run the same command again. The script will detect the existing installation and replace it with the latest version.

### Manual Installation

1. Go to the [releases page](https://github.com/javaquery/unosdk/releases)
2. Download the latest `unosdk.exe` binary for Windows
3. Move the binary to a permanent location (e.g., `C:\Program Files\unosdk\`)
4. Add the directory to your system PATH:

```powershell
# Open PowerShell and run:
$path = [Environment]::GetEnvironmentVariable('Path', 'User')
$newPath = $path + ';C:\Program Files\unosdk'
[Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
```

5. Verify installation:

```powershell
unosdk version
```

### Quick Start

After installation, you can immediately start using UnoSDK:

```bash
# List available SDKs
unosdk list

# Install Java
unosdk install java amazoncorretto 21

# Install Node.js
unosdk install node nodejs latest
```

## Usage

### Basic Commands

```bash
# Display help
unosdk --help

# Show version
unosdk version

# List all available providers and versions
unosdk list

# List installed SDKs
unosdk list --installed
```

### Install SDKs

```bash
# Install Amazon Corretto Java 21
unosdk install java amazoncorretto 21

# Install GraalVM
unosdk install java graalvm 23.1.2

# Install latest Node.js
unosdk install node nodejs latest

# Install specific Python version
unosdk install python python 3.11

# Install latest Flutter SDK
unosdk install flutter flutter latest

# Install specific Flutter version
unosdk install flutter flutter 3.27.2

# Install Apache Maven
unosdk install maven apache 3.9.9

# Install specific Maven version
unosdk install maven apache 3.8.8

# Install with custom path
unosdk install java openjdk 17 --path C:\SDKs\java

# Skip environment setup
unosdk install java amazoncorretto 21 --skip-env

# Set as default version
unosdk install java openjdk 21 --set-default
```

### Switch Between Versions

```bash
# Switch to a different Java version
unosdk switch java openjdk 21

# Switch to a different Node.js version
unosdk switch node nodejs 20
```

### Uninstall SDKs

```bash
# Uninstall specific version
unosdk uninstall java amazoncorretto 21

# Force uninstall (skip confirmation)
unosdk uninstall java openjdk 17 --force
```

### Update SDK Registry

```bash
# Update the list of available SDKs
unosdk update
```

## Configuration

UnoSDK automatically manages configuration and keeps track of installed SDKs. All data is stored in:

```
%USERPROFILE%\.unosdk\
├── config.yaml          # User configuration
├── registry.json        # Installed SDKs registry
├── cache/               # Cached SDK metadata
└── sdks/                # Installed SDKs
```

By default, SDKs are installed to `%USERPROFILE%\.unosdk\` directory:
```
C:\Users\<username>\.unosdk\
├── java\
│   ├── amazoncorretto\
│   │   ├── 11\
│   │   ├── 17\
│   │   └── 21\
│   └── openjdk\
│       └── 21\
├── node\
│   └── nodejs\
│       └── 20\
├── python\
│   └── python\
│       └── 3.11\
└── maven\
    └── 3.9.9\
```

For example, Java Amazon Corretto 11 would be installed at:
```
C:\Users\<username>\.unosdk\java\amazoncorretto\11
```

You can customize the installation path using the `--path` flag when installing SDKs.

## Troubleshooting

### Command Not Found

If you get "command not found" after installation, ensure:
- The directory containing `unosdk.exe` is in your PATH
- You've opened a new terminal window after modifying PATH

### Permission Denied

Run PowerShell or Command Prompt as Administrator when:
- Installing SDKs (to set environment variables)
- Switching between SDK versions
- First-time setup

### SDK Not Working After Install

1. Verify the SDK is installed: `unosdk list --installed`
2. Check environment variables are set correctly
3. Open a new terminal to refresh environment variables
4. Try switching to the SDK version: `unosdk switch <sdk-type> <provider> <version>`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [x] Add support for GraalVM
- [ ] Add support for more SDK providers (Adoptium, etc.)
- [ ] Cross-platform support (Linux, macOS)
- [ ] GUI interface
- [ ] Automatic version detection from project files
- [ ] SDK cleanup and maintenance commands
- [ ] Integration with common build tools

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/javaquery/unosdk).

## Acknowledgments

Special thanks to all the SDK providers for making their distributions available.FAQ

**Q: Do I need to manually configure environment variables?**  
A: No, UnoSDK automatically configures PATH and other necessary environment variables.

**Q: Can I install multiple versions of the same SDK?**  
A: Yes, you can install multiple versions and switch between them using `unosdk switch`.

**Q: Where are the SDKs installed?**  
A: By default in `%USERPROFILE%\.unosdk\` (e.g., `C:\Users\<username>\.unosdk\java\amazoncorretto\11`), but you can specify a custom path with `--path`.

**Q: Is internet connection required?**  
A: Yes, for downloading SDKs. After installation, SDKs work offline.

**Q: Can I use this alongside other SDK managers?**  
A: Yes, but be aware of potential PATH conflicts. UnoSDK manages its own installations independently.

## Support

- **Issues**: Report bugs on [GitHub Issues](https://github.com/javaquery/unosdk/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/javaquery/unosdk/discussions)
- **Documentation**: Visit the [GitHub repository](https://github.com/javaquery/unosdk)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Add support for more SDK providers (GraalVM, Adoptium, etc.)
- [ ] Cross-platform support (Linux, macOS)
- [ ] GUI interface
- [ ] Automatic version detection from project files
- [ ] SDK cleanup and maintenance commands
- [ ] Integration with common build tools

---

## For Contributors

Interested in contributing to UnoSDK? Check out our development guide.

### Building from Source

```powershell
# Clone the repository
git clone https://github.com/javaquery/unosdk.git
cd unosdk

# Build the project (requires Go 1.21+)
.\scripts\build.ps1

# Run tests
go test ./...
```

### Version Management

To update the version, simply edit `pkg/version/version.go`:

```go
const Version = "1.2.0"  // Change this line
```

Then build and release:

```powershell
.\scripts\build.ps1
git commit -am "bump version to 1.2.0"
git tag v1.2.0
git push origin main --tags
```

### Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [zap](https://github.com/uber-go/zap) - Structured logging
- [progressbar](https://github.com/schollz/progressbar) - Terminal progress bars
- [grab](https://github.com/cavaliergopher/grab) - File downloading

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

**Made with ❤️ for the developer community**