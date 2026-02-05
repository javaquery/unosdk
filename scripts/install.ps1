# UnoSDK Installation Script for Windows
# This script downloads the latest unosdk release from GitHub and installs it

param(
    [string]$InstallPath = "$env:LOCALAPPDATA\unosdk",
    [switch]$Force
)

$ErrorActionPreference = "Stop"

# GitHub repository details
$GH_REPO = "javaquery/unosdk"
$GH_API_URL = "https://api.github.com/repos/$GH_REPO/releases/latest"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  UnoSDK Installation Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

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
        Write-Host "✓ Added to PATH: $Directory" -ForegroundColor Green
        
        # Update current session PATH
        $env:Path = "$env:Path;$Directory"
    } else {
        Write-Host "✓ Already in PATH: $Directory" -ForegroundColor Yellow
    }
}

# Check if unosdk is already installed
$existingInstall = $false
if (Test-Path "$InstallPath\unosdk.exe") {
    $existingInstall = $true
    Write-Host "⚠ Found existing installation at: $InstallPath" -ForegroundColor Yellow
    
    if (-not $Force) {
        $response = Read-Host "Do you want to replace it with the latest version? (Y/N)"
        if ($response -ne "Y" -and $response -ne "y") {
            Write-Host "Installation cancelled." -ForegroundColor Red
            exit 0
        }
    }
    Write-Host "→ Replacing existing installation..." -ForegroundColor Cyan
}

Write-Host "→ Fetching latest release information..." -ForegroundColor Cyan

try {
    # Fetch latest release information
    $release = Invoke-RestMethod -Uri $GH_API_URL -Headers @{ "User-Agent" = "unosdk-installer" }
    $version = $release.tag_name
    Write-Host "✓ Latest version: $version" -ForegroundColor Green
    
    # Find Windows binary asset
    $asset = $release.assets | Where-Object { $_.name -match "unosdk.*windows.*\.exe$|unosdk\.exe$" } | Select-Object -First 1
    
    if (-not $asset) {
        Write-Host "✗ Could not find Windows binary in release assets" -ForegroundColor Red
        Write-Host "Available assets:" -ForegroundColor Yellow
        $release.assets | ForEach-Object { Write-Host "  - $($_.name)" -ForegroundColor Yellow }
        exit 1
    }
    
    $downloadUrl = $asset.browser_download_url
    $fileName = $asset.name
    
    Write-Host "→ Downloading: $fileName" -ForegroundColor Cyan
    Write-Host "  URL: $downloadUrl" -ForegroundColor Gray
    
    # Create installation directory
    if (-not (Test-Path $InstallPath)) {
        New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
        Write-Host "✓ Created installation directory: $InstallPath" -ForegroundColor Green
    }
    
    # Download the binary
    $tempFile = Join-Path $env:TEMP "unosdk_download.exe"
    $progressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
    $progressPreference = 'Continue'
    
    Write-Host "✓ Download completed" -ForegroundColor Green
    
    # Stop any running unosdk processes
    $runningProcesses = Get-Process -Name "unosdk" -ErrorAction SilentlyContinue
    if ($runningProcesses) {
        Write-Host "→ Stopping running unosdk processes..." -ForegroundColor Cyan
        $runningProcesses | Stop-Process -Force
        Start-Sleep -Seconds 1
    }
    
    # Move downloaded file to installation directory
    $targetPath = Join-Path $InstallPath "unosdk.exe"
    Move-Item -Path $tempFile -Destination $targetPath -Force
    
    Write-Host "✓ Installed to: $targetPath" -ForegroundColor Green
    
    # Add to PATH
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
        Write-Host "✓ Existing installation has been replaced" -ForegroundColor Green
    }
    
    Write-Host "You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Try it out:" -ForegroundColor Cyan
    Write-Host "  unosdk --version" -ForegroundColor White
    Write-Host "  unosdk --help" -ForegroundColor White
    Write-Host ""
    
} catch {
    Write-Host "✗ Installation failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
