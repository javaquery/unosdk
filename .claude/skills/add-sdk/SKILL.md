---
name: add-sdk
description: Add a new Software Development Kit (SDK) to the project. Use when user wants to "add a new SDK" or "integrate a new SDK" or "add ${SDK_NAME} in the project".
---

# Add SDK Skill

### Overview

Generate code to add a new Software Development Kit (SDK) to the project. User can specify the SDK they want to add, and the skill will provide the necessary code snippets and instructions to integrate it into the project. This can include installation commands, configuration settings, and example usage.

### Step 1. Gather project context for an existing SDK

Before adding a new SDK, gather necessary information about it, such as installation instructions, configuration requirements, and example usage.

- Identify the SDK name provided by the user.
- Research the SDK on its official documentation or repository to find installation binaries, commands, configuration settings, and example usage.
- There may be multiple versions of the SDK available, so ensure to gather information about the latest stable version or the version specified by the user and few older versions as well.
- There may be different providers or distributions of the SDK (e.g., OpenJDK vs Oracle JDK for Java), so gather information about the most popular or widely used ones.
- Note any dependencies or prerequisites required for the SDK.
- We will try to use *.zip, *.exe, *.msi files for installation on Windows machines, so prioritize gathering information about these types of installation files for the SDK.
- We must look for the distribution of the SDK that is compatible with Windows machines, as the user will be installing it on a Windows machine.
- If the SDK has a command-line interface (CLI) for installation or management, gather information about the available commands and their usage as well.
- We will try to install SDK without user having to configure the installation process, so gather information about any silent installation options or default configuration settings that can be used during installation to minimize user input and configuration steps.

### Step 2. Project structure

- ./internal/cli/* - This directory contains the CLI commands for managing SDKs, including installation, switching, and configuration. Use this as a reference for generating code snippets related to SDK management.
- ./internal/installers/* - This directory contains the logic for installing various SDKs. Use this as a reference for generating code snippets related to SDK installation.
- ./internal/config/* - This directory contains configuration management logic, which can be referenced for generating code snippets related to setting environment variables or configuring the SDK.
- ./internal/providers/* - This directory contains logic for handling different SDK providers, which can be referenced for generating code snippets related to supporting multiple versions or providers of the SDK.

### Step 3. Generate code snippets for adding the SDK

Based on the gathered information, generate code snippets that the user can use to add the SDK to a project and later use it on a Windows machine to install, switch, or configure the SDK.

- Provide installation commands (e.g., unosdk install <SDK_NAME>, unosdk switch <SDK_NAME>, unosdk uninstall <SDK_NAME>). 
- Set path variables or environment settings in the system. (e.g JAVA_HOME for Java SDKs). Use `./internal/cli/env_setup.go` as a reference for setting environment variables.
- Include example usage of the SDK in a sample code snippet to demonstrate how to use it in a project.
- If there are multiple versions or providers of the SDK, provide code snippets for each to give the user options to choose from.
- Ensure that the generated code snippets are clear, concise, and easy to follow for users of varying technical expertise.

### Step 4. Checks

- [] Ensure that the generated code snippets are accurate and up-to-date with the latest version of the SDK.
- [] Verify that the installation commands and configuration settings are correct and will work on a Windows machine.
- [] Ensure that the example usage provided is relevant and demonstrates the key features of the SDK effectively
- [] Check that the code snippets are formatted correctly and are easy to read and understand.
- [] We must have updated project version and changelog (./CHANGELOG.md) to reflect the addition of the new SDK in the project. Do only this if current version is not released yet. If current version is already released, then we will update the version and changelog in the next release when we have more changes to add along with the new SDK addition.