# UnoSDK Installation Script for Windows
# This script downloads the latest unosdk release from GitHub and installs it.
#
# Layout after install:
#   <InstallPath>\unosdk.cmd        <- wrapper for cmd.exe (this folder is on PATH)
#   <InstallPath>\core\unosdk.exe   <- real binary
# plus an 'unosdk' function in your PowerShell profile(s).
# The wrappers let `unosdk switch ...` take effect in the current terminal
# without restarting it.

param(
    [string]$InstallPath = "$env:LOCALAPPDATA\unosdk",
    [switch]$Force
)

$ErrorActionPreference = "Stop"

# Function to pause and wait for user input before closing
function Wait-ForKeyPress {
    param([string]$Message = "Press any key to continue...")
    Write-Host ""
    Write-Host $Message -ForegroundColor Gray
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
}

# Function to convert short 8.3 paths to full paths
function Get-LongPath {
    param([string]$Path)
    
    try {
        # If path doesn't exist, try to create parent directories
        $parent = Split-Path $Path -Parent
        if ($parent -and -not (Test-Path -LiteralPath $parent)) {
            New-Item -ItemType Directory -Path $parent -Force -ErrorAction SilentlyContinue | Out-Null
        }
        
        # Get the full path
        $fullPath = [System.IO.Path]::GetFullPath($Path)
        
        # Convert short names to long names if the path exists
        if (Test-Path -LiteralPath $fullPath) {
            $item = Get-Item -LiteralPath $fullPath -Force
            return $item.FullName
        }
        
        return $fullPath
    } catch {
        return $Path
    }
}

# GitHub repository details
$GH_REPO = "javaquery/unosdk"
$GH_API_URL = "https://api.github.com/repos/$GH_REPO/releases/latest"

# Function to check PowerShell execution policy
function Test-ExecutionPolicy {
    $executionPolicy = Get-ExecutionPolicy -Scope CurrentUser
    $machinePolicy = Get-ExecutionPolicy -Scope LocalMachine
    
    Write-Host "[*] Checking PowerShell execution policy..." -ForegroundColor Cyan
    Write-Host "  Current User: $executionPolicy" -ForegroundColor Gray
    Write-Host "  Local Machine: $machinePolicy" -ForegroundColor Gray
    Write-Host ""
    
    if ($executionPolicy -eq "Restricted" -or ($executionPolicy -eq "Undefined" -and $machinePolicy -eq "Restricted")) {
        Write-Host "[!] WARNING: PowerShell script execution is restricted!" -ForegroundColor Red
        Write-Host ""
        Write-Host "Your system's execution policy prevents running scripts." -ForegroundColor Yellow
        Write-Host "This also stops your PowerShell profile (and the 'unosdk' function) from loading." -ForegroundColor Yellow
        Write-Host "To fix this, run PowerShell as Administrator and execute:" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "  Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser" -ForegroundColor White
        Write-Host ""
        Write-Host "Or, for this session only, run:" -ForegroundColor Yellow
        Write-Host "  powershell -ExecutionPolicy Bypass -File install.ps1" -ForegroundColor White
        Write-Host ""
        
        $response = Read-Host "Do you want to continue anyway? (Y/N)"
        if ($response -ne "Y" -and $response -ne "y") {
            Write-Host "Installation cancelled." -ForegroundColor Red
            Wait-ForKeyPress
            exit 1
        }
    } elseif ($executionPolicy -eq "AllSigned") {
        Write-Host "[!] Note: Only signed scripts can run with current policy ($executionPolicy)" -ForegroundColor Yellow
        Write-Host "If this script is not signed, you may encounter errors." -ForegroundColor Yellow
        Write-Host ""
    } else {
        Write-Host "[OK] Execution policy is adequate ($executionPolicy)" -ForegroundColor Green
        Write-Host ""
    }
}

# Function to check if running as administrator
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Function to add directory to PATH
function Add-ToPath {
    param([string]$Directory)
    
    $currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    
    if ($currentPath -notlike "*$Directory*") {
        $newPath = if ($currentPath) { "$currentPath;$Directory" } else { $Directory }
        [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        Write-Host "[OK] Added to PATH: $Directory" -ForegroundColor Green
        
        # Update current session PATH
        $env:Path = "$env:Path;$Directory"
    } else {
        Write-Host "[OK] Already in PATH: $Directory" -ForegroundColor Yellow
    }
}

# Writes unosdk.cmd (wrapper for cmd.exe) into the install directory
function Install-CmdWrapper {
    param([string]$Directory)

    $wrapper = @'
@echo off
rem unosdk shell wrapper for cmd.exe
rem Runs the real binary (core\unosdk.exe) and then applies any environment
rem changes it requested to THIS shell, so no terminal restart is needed.
rem
rem From another .bat file, call it with:  call unosdk switch java openjdk 21

set "UNOSDK_ENV_FILE=%TEMP%\unosdk_%RANDOM%%RANDOM%.cmd"
"%~dp0core\unosdk.exe" %*
set "UNOSDK_RC=%ERRORLEVEL%"
if exist "%UNOSDK_ENV_FILE%" (
    call "%UNOSDK_ENV_FILE%"
    del "%UNOSDK_ENV_FILE%" >nul 2>&1
)
set "UNOSDK_ENV_FILE=" & set "UNOSDK_RC=" & exit /b %UNOSDK_RC%
'@

    $target = Join-Path $Directory "unosdk.cmd"
    Set-Content -LiteralPath $target -Value $wrapper -Encoding ASCII
    Write-Host "[OK] Created cmd wrapper: $target" -ForegroundColor Green
}

# Adds (or refreshes) the 'unosdk' function in the PowerShell profiles
function Install-PowerShellFunction {
    param([string]$CoreExe)

    $function = @'
# >>> unosdk >>>
function unosdk {
    $rc = 0
    $env:UNOSDK_ENV_FILE = Join-Path $env:TEMP "unosdk_$([guid]::NewGuid().ToString('N')).ps1"
    try {
        & '__CORE_EXE__' @args
        $rc = $LASTEXITCODE
        if (Test-Path $env:UNOSDK_ENV_FILE) { . $env:UNOSDK_ENV_FILE }
    } finally {
        Remove-Item $env:UNOSDK_ENV_FILE -ErrorAction SilentlyContinue
        Remove-Item Env:UNOSDK_ENV_FILE -ErrorAction SilentlyContinue
    }
    $global:LASTEXITCODE = $rc
}
# <<< unosdk <<<
'@
    $function = $function.Replace("__CORE_EXE__", $CoreExe.Replace("'", "''"))

    # Windows PowerShell 5.1 and PowerShell 7+ use different profile folders; update both
    $docs = [Environment]::GetFolderPath('MyDocuments')
    $profiles = @(
        (Join-Path $docs "WindowsPowerShell\profile.ps1"),
        (Join-Path $docs "PowerShell\profile.ps1")
    )

    foreach ($p in $profiles) {
        try {
            $dir = Split-Path $p -Parent
            if (-not (Test-Path -LiteralPath $dir)) {
                New-Item -ItemType Directory -Path $dir -Force | Out-Null
            }

            $existing = ""
            if (Test-Path -LiteralPath $p) {
                $raw = Get-Content -LiteralPath $p -Raw
                if ($null -ne $raw) { $existing = $raw }
            }

            # Replace any previous unosdk block, then append the fresh one
            $existing = [regex]::Replace($existing, '(?s)# >>> unosdk >>>.*?# <<< unosdk <<<\r?\n?', '')
            $content = $existing.TrimEnd() + "`r`n`r`n" + $function
            Set-Content -LiteralPath $p -Value $content.TrimStart() -Encoding UTF8
            Write-Host "[OK] Updated PowerShell profile: $p" -ForegroundColor Green
        } catch {
            Write-Host "[!] Could not update profile $p : $($_.Exception.Message)" -ForegroundColor Yellow
        }
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  UnoSDK Installation Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check execution policy before proceeding
Test-ExecutionPolicy

# Normalize installation path to handle 8.3 short names
$InstallPath = Get-LongPath $InstallPath
Write-Host "[*] Installation directory: $InstallPath" -ForegroundColor Gray
Write-Host ""

# Layout: wrapper in root, real binary in core\
$coreDir    = Join-Path $InstallPath "core"
$targetPath = Join-Path $coreDir "unosdk.exe"
$legacyExe  = Join-Path $InstallPath "unosdk.exe"

# Check if unosdk is already installed (new layout or legacy root-level exe)
$existingInstall = $false
if ((Test-Path -LiteralPath $targetPath) -or (Test-Path -LiteralPath $legacyExe)) {
    $existingInstall = $true
    Write-Host "[!] Found existing installation at: $InstallPath" -ForegroundColor Yellow
    
    if (-not $Force) {
        $response = Read-Host "Do you want to replace it with the latest version? (Y/N)"
        if ($response -ne "Y" -and $response -ne "y") {
            Write-Host "Installation cancelled." -ForegroundColor Red
            Wait-ForKeyPress
            exit 0
        }
    }
    Write-Host "[*] Replacing existing installation..." -ForegroundColor Cyan
}

Write-Host "[*] Fetching latest release information..." -ForegroundColor Cyan

try {
    # Fetch latest release information
    $release = Invoke-RestMethod -Uri $GH_API_URL -Headers @{ "User-Agent" = "unosdk-installer" }
    $version = $release.tag_name
    Write-Host "[OK] Latest version: $version" -ForegroundColor Green
    
    # Find Windows binary asset
    $asset = $release.assets | Where-Object { $_.name -match "unosdk.*windows.*\.exe$|unosdk\.exe$" } | Select-Object -First 1
    
    if (-not $asset) {
        Write-Host "[ERROR] Could not find Windows binary in release assets" -ForegroundColor Red
        Write-Host "Available assets:" -ForegroundColor Yellow
        $release.assets | ForEach-Object { Write-Host "  - $($_.name)" -ForegroundColor Yellow }
        Wait-ForKeyPress
        exit 1
    }
    
    $downloadUrl = $asset.browser_download_url
    $fileName = $asset.name
    
    Write-Host "[*] Downloading: $fileName" -ForegroundColor Cyan
    Write-Host "  URL: $downloadUrl" -ForegroundColor Gray
    
    # Create installation directories
    if (-not (Test-Path -LiteralPath $InstallPath)) {
        New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
        Write-Host "[OK] Created installation directory: $InstallPath" -ForegroundColor Green
    }
    if (-not (Test-Path -LiteralPath $coreDir)) {
        New-Item -ItemType Directory -Path $coreDir -Force | Out-Null
    }
    
    # Download the binary
    $tempDir = Get-LongPath $env:TEMP
    $tempFile = Join-Path $tempDir "unosdk_download.exe"
    
    # Ensure temp directory exists
    if (-not (Test-Path -LiteralPath $tempDir)) {
        Write-Host "[ERROR] Temp directory does not exist: $tempDir" -ForegroundColor Red
        Write-Host "[ERROR] Please check your TEMP environment variable" -ForegroundColor Yellow
        Wait-ForKeyPress "Press any key to exit..."
        exit 1
    }
    
    # Remove old temp file if exists
    if (Test-Path -LiteralPath $tempFile) {
        Remove-Item -LiteralPath $tempFile -Force -ErrorAction SilentlyContinue
    }
    
    Write-Host "[*] Downloading to: $tempFile" -ForegroundColor Gray
    $progressPreference = 'SilentlyContinue'
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
    } catch {
        Write-Host "[ERROR] Download failed: $($_.Exception.Message)" -ForegroundColor Red
        Wait-ForKeyPress "Press any key to exit..."
        exit 1
    } finally {
        $progressPreference = 'Continue'
    }
    
    # Verify the download was successful
    if (-not (Test-Path -LiteralPath $tempFile)) {
        Write-Host "[ERROR] Downloaded file not found at: $tempFile" -ForegroundColor Red
        Write-Host "[ERROR] The download may have failed silently." -ForegroundColor Red
        Wait-ForKeyPress "Press any key to exit..."
        exit 1
    }
    
    $fileSize = (Get-Item -LiteralPath $tempFile).Length
    if ($fileSize -eq 0) {
        Write-Host "[ERROR] Downloaded file is empty (0 bytes)" -ForegroundColor Red
        Wait-ForKeyPress "Press any key to exit..."
        exit 1
    }
    
    Write-Host "[OK] Download completed ($([math]::Round($fileSize/1MB, 2)) MB)" -ForegroundColor Green
    
    # Stop any running unosdk processes
    $runningProcesses = Get-Process -Name "unosdk" -ErrorAction SilentlyContinue
    if ($runningProcesses) {
        Write-Host "[*] Stopping running unosdk processes..." -ForegroundColor Cyan
        $runningProcesses | Stop-Process -Force
        Start-Sleep -Seconds 1
    }
    
    # Move downloaded file to core\ inside the installation directory
    Write-Host "[*] Installing to: $targetPath" -ForegroundColor Gray
    
    # Verify source file exists before moving
    if (-not (Test-Path -LiteralPath $tempFile)) {
        Write-Host "[ERROR] Source file disappeared: $tempFile" -ForegroundColor Red
        Wait-ForKeyPress "Press any key to exit..."
        exit 1
    }
    
    try {
        Move-Item -LiteralPath $tempFile -Destination $targetPath -Force -ErrorAction Stop
    } catch {
        Write-Host "[ERROR] Failed to move file to installation directory" -ForegroundColor Red
        Write-Host "[ERROR] Details: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "[ERROR] Source: $tempFile" -ForegroundColor Yellow
        Write-Host "[ERROR] Target: $targetPath" -ForegroundColor Yellow
        Write-Host "[ERROR] You may need to run as Administrator or check folder permissions" -ForegroundColor Yellow
        Wait-ForKeyPress "Press any key to exit..."
        exit 1
    }
    
    # Verify installation succeeded
    if (-not (Test-Path -LiteralPath $targetPath)) {
        Write-Host "[ERROR] Installation verification failed - file not found at: $targetPath" -ForegroundColor Red
        Wait-ForKeyPress "Press any key to exit..."
        exit 1
    }
    
    Write-Host "[OK] Installed to: $targetPath" -ForegroundColor Green
    
    # Remove legacy root-level exe so the wrapper wins (existing users migrating).
    # In cmd.exe, unosdk.exe beats unosdk.cmd when both are in the same folder.
    if (Test-Path -LiteralPath $legacyExe) {
        Remove-Item -LiteralPath $legacyExe -Force -ErrorAction SilentlyContinue
        Write-Host "[OK] Removed legacy binary: $legacyExe" -ForegroundColor Green
    }
    
    # Wrappers that let switch/install/uninstall apply to the current terminal
    Install-CmdWrapper -Directory $InstallPath
    Install-PowerShellFunction -CoreExe $targetPath
    
    # Add the root folder (contains unosdk.cmd) to PATH. core\ is NOT added.
    Add-ToPath -Directory $InstallPath
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  Installation Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Installed version: $version" -ForegroundColor Cyan
    Write-Host "Installation path: $InstallPath" -ForegroundColor Cyan
    Write-Host ""
    
    if ($existingInstall) {
        Write-Host "[OK] Existing installation has been replaced" -ForegroundColor Green
    }
    
    Write-Host "Open a NEW terminal once so the PATH and the 'unosdk' wrapper load." -ForegroundColor Yellow
    Write-Host "After that, 'unosdk switch ...' applies immediately with no further restarts." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Try it out:" -ForegroundColor Cyan
    Write-Host "  unosdk version" -ForegroundColor White
    Write-Host "  unosdk --help" -ForegroundColor White
    Write-Host ""
    
} catch {
    Write-Host "[ERROR] Installation failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "[ERROR] Stack trace: $($_.ScriptStackTrace)" -ForegroundColor Red
    Wait-ForKeyPress "Press any key to exit..."
    exit 1
}