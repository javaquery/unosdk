---
name: update-sdk-version
description: Update the version of an existing Software Development Kit (SDK) in the project. Use when user wants to "update SDK version" or "upgrade SDK" or "change ${SDK_NAME} version in the project" or "check the latest versions of the SDK(s)".
---

# Update SDK Version Skill

### Overview
Generate code to update the version of an existing Software Development Kit (SDK) in the project. User can specify the SDK they want to update and the new version they want to use, and the skill will provide the necessary code snippets and instructions to perform the update. 

### Step 1. Gather project context for the existing SDK

Before updating the SDK version, gather necessary information about the existing SDK, such as the current version being used, installation instructions for the new version, and any changes or breaking changes in the new version.

- Identify the SDK name and current version being used in the project.
- Research the new version of the SDK on its official documentation or repository to find installation instructions, configuration changes, and any breaking changes or new features introduced in the new version.
- Gather information about the download links for the new version of the SDK, especially for Windows machines, and any specific installation files (e.g., .zip, .exe, .msi) that are available for the new version.
- Note any dependencies or prerequisites required for the new version of the SDK.
- If there are multiple versions of the SDK available, gather information about the latest stable version and a few older versions as well to provide options for the user.
- We will try to use *.zip, *.exe, *.msi files for installation on Windows machines, so prioritize gathering information about these types of installation files for the SDK.
- We must look for the distribution of the SDK that is compatible with Windows machines, as the user will be installing it on a Windows machine.
- If the SDK has a command-line interface (CLI) for installation or management, gather information about the available commands and their usage as well.
- We will try to install SDK without user having to configure the installation process, so gather information about any silent installation options or default configuration settings that can be used during installation to minimize user input and configuration steps.
- Check which version of the SDK we have already setup (./ineternal/providers/*) and update only the latest version only.

### Step 2. Project structure

- ./internal/cli/* - This directory contains the CLI commands for managing SDKs, including installation, switching, and configuration. Use this as a reference for generating code snippets related to SDK management.
- ./internal/installers/* - This directory contains the logic for installing various SDKs. Use this as a reference for generating code snippets related to SDK installation.
- ./internal/config/* - This directory contains configuration management logic, which can be referenced for generating code snippets related to setting environment variables or configuring the SDK.
- ./internal/providers/* - This directory contains logic for handling different SDK providers, which can be referenced for generating code snippets related to supporting multiple versions or providers of the SDK.

### Step 3. Generate code snippets for updating the SDK version

Based on the gathered information, generate code snipppet to add latest SDK version to the project with the download link and installation instructions. Also, provide code snippets that the user can use to switch to the new SDK version on a Windows machine.

### Step 4. Checks

- [] Ensure that the generated code snippets are accurate and up-to-date with the latest version of the SDK.
- [] Verify that the installation commands and configuration settings are correct and will work on a Windows machine
- [] Ensure that the example usage provided is relevant and demonstrates the key features of the new SDK version effectively
- [] Check that the code snippets are formatted correctly and are easy to read and understand.
- [] We must have updated project version and changelog (./CHANGELOG.md) to reflect the addition of the new SDK in the project. Do only this if current version is not released yet. If current version is already released, then we will update the version and changelog in the next release when we have more changes to add along with the new SDK addition.
- [] Run go mod tidy to ensure that any new dependencies added for the SDK are properly included in the project.
- [] Run all tests to ensure that the update of the new SDK does not break any existing functionality in the project.