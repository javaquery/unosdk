# UnoSDK

![UnoSDK Logo](https://github.com/javaquery/unosdk/raw/master/unosdk%20list.png)

**UnoSDK** is a native CLI tool for Windows that installs and manages multiple software development kits (SDKs) from various providers. Think of it as **SDKMAN for Windows**, built as a single `.exe` that runs directly in PowerShell or Command Prompt. Say goodbye to manual downloads, extractions, and environment variable configuration.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/javaquery/unosdk?style=flat&logo=github)](https://github.com/javaquery/unosdk/releases/latest)
[![CI](https://github.com/javaquery/unosdk/actions/workflows/ci.yml/badge.svg)](https://github.com/javaquery/unosdk/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/javaquery/unosdk/blob/master/LICENSE)

## Table of Contents

- [Why UnoSDK?](#why-unosdk)
- [UnoSDK vs SDKMAN on Windows](#unosdk-vs-sdkman-on-windows)
- [Features](#features)
- [Supported SDKs](#supported-sdks)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

## Why UnoSDK?

If you've used [SDKMAN!](https://sdkman.io/) on Linux or macOS and wished for the same experience on Windows, **UnoSDK** is for you. It gives you a native Windows workflow for managing multiple SDK versions, without bash, WSL, or manual PATH management.

## UnoSDK vs SDKMAN on Windows

SDKMAN is an excellent tool, but it is a set of bash scripts. On Windows it only runs inside a Unix-like layer such as WSL, Git Bash, Cygwin, or MSYS2. UnoSDK is a compiled Go binary built specifically for Windows.

### Benefits at a glance

| | **UnoSDK** | **SDKMAN** on Windows |
| --- | --- | --- |
| **Runs in** | PowerShell and Command Prompt (native) | Needs WSL, Git Bash, Cygwin, or MSYS2 |
| **Installation** | One PowerShell command, or drop in a single `.exe` | Requires a bash environment plus `curl`, `zip`, and `unzip` first |
| **Environment variables** | Sets Windows `PATH` / `JAVA_HOME` etc. for you | Works inside the bash shell session; not wired into native Windows apps |
| **Works with Windows tools** | IDEs, Windows Terminal, PowerShell, cmd, CI agents, and services all see the same SDKs | SDKs installed in WSL live in the Linux filesystem and are not visible to Windows-native tools |
| **SDKs available** | Java, Node.js, Python, Flutter, Maven, Gradle, Go, C, C++ | Strong JVM ecosystem (Java, Maven, Gradle, Kotlin, Scala, and more), but no Node.js, Python, or Go |
| **Windows binaries** | Downloads Windows builds (e.g. MinGW-w64 for C/C++) | In WSL, downloads Linux builds, which cannot run natively on Windows |
| **Footprint** | Single static binary, no runtime dependencies | Shell scripts plus the Unix tooling they depend on |

### Why this matters

- **No Linux layer required.** You don't need to install or maintain WSL or Git Bash just to switch Java versions.
- **One tool for your whole stack.** Manage your JDK, build tools, Node.js, Python, Go, Flutter, and a C/C++ toolchain together instead of juggling several installers.
- **Native paths and variables.** SDKs are installed under `%USERPROFILE%\.unosdk\` (or any path you choose) and registered in the Windows environment, so Visual Studio Code, IntelliJ IDEA, Android Studio, and build servers can use them without extra setup.
- **Fast, safe installs.** Parallel downloads with progress bars and checksum verification.
- **Easy to update.** Re-run the install command to upgrade UnoSDK itself.

### When SDKMAN might still be a better fit

- You already work primarily inside WSL or a Linux/macOS environment. SDKMAN is the established choice there.
- You need a JVM-ecosystem tool that UnoSDK doesn't support yet (for example Kotlin or Scala).
- You want SDKMAN's wide range of Java vendor distributions (UnoSDK currently offers Amazon Corretto, OpenJDK, and GraalVM).

> UnoSDK and SDKMAN can coexist. UnoSDK manages Windows-native SDKs, while SDKMAN inside WSL manages Linux-side ones.

## Features

- 🚀 **Multi-SDK Support**: Manage Java, Node.js, Python, Flutter, Maven, Gradle, Go, C, and C++ from a single tool
- 🪟 **Windows Native**: A single `.exe` that runs in PowerShell and Command Prompt, with no WSL or bash needed
- 🔄 **Version Switching**: Switch between installed SDK versions with one command
- 📦 **Multiple Providers**: Support for various distribution providers
  - Java: Amazon Corretto, OpenJDK, GraalVM
  - Node.js: Official Node.js distributions
  - Python: Official Python distributions
  - Flutter: Official Flutter SDK
  - Maven: Apache Maven build tool
  - Gradle: Gradle build automation tool
  - Go: Official Go programming language
  - C / C++: MinGW-w64 (GCC/G++ toolchain)
- 🔧 **Automatic Environment Setup**: Configures PATH and environment variables for you
- 📋 **Registry Management**: Keeps track of all installed SDKs
- ⚡ **Fast Downloads**: Parallel downloads with progress tracking
- 🛡️ **Verification**: Checksum verification to ensure download integrity

## Supported SDKs

| SDK Type | Providers                         | Description                                   |
| -------- | --------------------------------- | --------------------------------------------- |
| Java     | Amazon Corretto, OpenJDK, GraalVM | Java Development Kit                          |
| Node.js  | nodejs                            | JavaScript runtime environment                |
| Python   | python                            | Python programming language                   |
| Flutter  | flutter                           | Flutter SDK for mobile, web, and desktop apps |
| Maven    | apache                            | Apache Maven build automation tool            |
| Gradle   | gradle                            | Gradle build automation tool                  |
| Go       | golang                            | Go programming language                       |
| C        | mingw                             | MinGW-w64 GCC toolchain                       |
| C++      | mingw                             | MinGW-w64 GCC/G++ toolchain                   |

## Installation

### Prerequisites

- Windows 10 or later
- PowerShell 5.1 or later

### Quick Installation (Recommended)

Open PowerShell and run:

```powershell
irm https://raw.githubusercontent.com/javaquery/unosdk/refs/heads/master/scripts/install.ps1 | iex
```

This will automatically:

- Download the latest release from GitHub
- Install to `%LOCALAPPDATA%\unosdk`
- Add unosdk to your PATH
- Replace any existing installation

**To update UnoSDK:** run the same command again. The script detects the existing installation and replaces it with the latest version.

### Manual Installation

1. Go to the [releases page](https://github.com/javaquery/unosdk/releases)
2. Download the latest `unosdk.exe` for Windows
3. Move it to a permanent location (e.g., `C:\Program Files\unosdk\`)
4. Add that directory to your PATH:

```powershell
$path = [Environment]::GetEnvironmentVariable('Path', 'User')
$newPath = $path + ';C:\Program Files\unosdk'
[Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
```

5. Open a new terminal and verify:

```powershell
unosdk version
```

### Quick Start

```powershell
# List available SDKs
unosdk list

# Install Java
unosdk install java amazoncorretto 21

# Install Node.js
unosdk install node nodejs latest
```

## Usage

### Basic Commands

```powershell
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

```powershell
# Java
unosdk install java amazoncorretto 21
unosdk install java graalvm 23.1.2

# Node.js
unosdk install node nodejs latest

# Python
unosdk install python python 3.11

# Flutter
unosdk install flutter flutter latest
unosdk install flutter flutter 3.27.2

# Maven
unosdk install maven apache 3.9.9
unosdk install maven apache 3.8.8

# Gradle
unosdk install gradle gradle 8.12
unosdk install gradle gradle 8.10

# Go
unosdk install go golang 1.23.5
unosdk install go golang 1.22.10

# C++ (MinGW-w64)
unosdk install cpp mingw 15.2.0
unosdk install cpp mingw 14.2.0

# C (MinGW-w64)
unosdk install c mingw 15.2.0
```

**Install options:**

```powershell
# Install to a custom path
unosdk install java openjdk 17 --path C:\SDKs\java

# Skip environment setup
unosdk install java amazoncorretto 21 --skip-env

# Set as default version
unosdk install java openjdk 21 --set-default
```

### Switch Between Versions

```powershell
unosdk switch java openjdk 21
unosdk switch node nodejs 20
unosdk switch gradle gradle 8.12
unosdk switch go golang 1.23.5
unosdk switch cpp mingw 15.2.0
unosdk switch c mingw 15.2.0
```

### Uninstall SDKs

```powershell
# Uninstall a specific version
unosdk uninstall java amazoncorretto 21

# Force uninstall (skip confirmation)
unosdk uninstall java openjdk 17 --force
```

### Update SDK Registry

```powershell
# Refresh the list of available SDKs
unosdk update
```

## Configuration

UnoSDK manages its configuration and tracks installed SDKs automatically. Data is stored in:

```
%USERPROFILE%\.unosdk\
├── config.yaml          # User configuration
├── registry.json        # Installed SDKs registry
├── cache/               # Cached SDK metadata
└── sdks/                # Installed SDKs
```

By default, SDKs are installed under `%USERPROFILE%\.unosdk\`:

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
├── maven\
│   └── 3.9.9\
├── gradle\
│   └── 8.12\
├── go\
│   └── golang\
│       └── 1.23.5\
├── c\
│   └── mingw\
│       └── 15.2.0\
│           └── mingw64\  # bin/ (gcc), include/, lib/, etc.
└── cpp\
    └── mingw\
        └── 15.2.0\
            └── mingw64\  # bin/ (g++, gcc), include/, lib/, etc.
```

For example, Amazon Corretto 11 is installed at `C:\Users\<username>\.unosdk\java\amazoncorretto\11`. Use `--path` to choose a different location.

## Troubleshooting

### Command Not Found

- Make sure the directory containing `unosdk.exe` is in your PATH
- Open a new terminal window after PATH changes

### Permission Denied

Run PowerShell or Command Prompt as Administrator when:

- Installing SDKs (to set environment variables)
- Switching between SDK versions
- First-time setup

### SDK Not Working After Install

1. Verify the SDK is installed: `unosdk list --installed`
2. Check that environment variables are set correctly
3. Open a new terminal to refresh environment variables
4. Re-select the version: `unosdk switch <sdk-type> <provider> <version>`

## FAQ

**Q: How is UnoSDK different from SDKMAN?**
A: SDKMAN is a bash-based tool designed for Linux and macOS; on Windows it requires WSL, Git Bash, Cygwin, or MSYS2. UnoSDK is a native Windows executable that works in PowerShell and Command Prompt, sets Windows environment variables directly, and also covers Node.js, Python, Go, Flutter, and C/C++. See [UnoSDK vs SDKMAN on Windows](#unosdk-vs-sdkman-on-windows).

**Q: Do I need WSL or Git Bash?**
A: No. UnoSDK is a standalone Windows binary.

**Q: Do I need to configure environment variables manually?**
A: No. UnoSDK configures PATH and other required variables automatically.

**Q: Can I install multiple versions of the same SDK?**
A: Yes. Install as many as you like and switch with `unosdk switch`.

**Q: Where are SDKs installed?**
A: By default in `%USERPROFILE%\.unosdk\` (e.g., `C:\Users\<username>\.unosdk\java\amazoncorretto\11`). Use `--path` to customize.

**Q: Is an internet connection required?**
A: Only for downloading SDKs. Once installed, SDKs work offline.

**Q: Can I use this alongside other SDK managers (including SDKMAN in WSL)?**
A: Yes, but watch for PATH conflicts. UnoSDK manages its own installations independently.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

- **Issues**: Report bugs on [GitHub Issues](https://github.com/javaquery/unosdk/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/javaquery/unosdk/discussions)

## Acknowledgments

Special thanks to all the SDK providers for making their distributions available.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/javaquery/unosdk/blob/master/LICENSE) file for details.

---

## For Contributors

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

Edit `pkg/version/version.go`:

```go
const Version = "1.2.0"  // Change this line
```

Then build and release:

```powershell
.\scripts\build.ps1
git commit -am "bump version to 1.2.0"
git tag v1.2.0
git push origin master --tags
```

### Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [zap](https://github.com/uber-go/zap) - Structured logging
- [progressbar](https://github.com/schollz/progressbar) - Terminal progress bars
- [grab](https://github.com/cavaliergopher/grab) - File downloading

---

**Made with ❤️ for the developer community**