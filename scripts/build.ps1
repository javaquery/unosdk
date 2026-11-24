$ErrorActionPreference = "Stop"

# Change to project root directory (parent of scripts)
$scriptDir = $PSScriptRoot
$projectRoot = Split-Path $scriptDir -Parent
Set-Location $projectRoot

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "  unosdk - Build Script" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# Get version from version.go file
$versionFile = "pkg/version/version.go"
$versionContent = Get-Content $versionFile -Raw
if ($versionContent -match 'Version\s*=\s*"([^"]+)"') {
    $version = $matches[1]
} else {
    $version = "dev"
}

$commit = try { git rev-parse --short HEAD 2>$null } catch { "unknown" }
$date = Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"

Write-Host "Version: $version" -ForegroundColor Green
Write-Host "Commit:  $commit" -ForegroundColor Green
Write-Host "Date:    $date" -ForegroundColor Green
Write-Host ""

# Clear any existing GOOS/GOARCH environment variables to ensure tests run on native architecture
Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue

# Run tests before building
Write-Host "Running tests..." -ForegroundColor Yellow
go test ./...

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "=====================================" -ForegroundColor Red
    Write-Host "  Tests Failed!" -ForegroundColor Red
    Write-Host "=====================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Build aborted. Please fix failing tests before building." -ForegroundColor Red
    exit 1
}

Write-Host "  All tests passed!" -ForegroundColor Green
Write-Host ""

# Build flags - inject into version package variables
$ldflags = "-s -w -X github.com/javaquery/unosdk/pkg/version.GitCommit=$commit -X github.com/javaquery/unosdk/pkg/version.BuildDate=$date"

# Create bin directory if it doesn't exist
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

Write-Host "Building unosdk..." -ForegroundColor Yellow

# Build for Windows AMD64
Write-Host "  → Building for Windows AMD64..." -ForegroundColor Gray
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags="$ldflags" -o bin/unosdk-amd64.exe ./cmd/unosdk

if ($LASTEXITCODE -eq 0) {
    Write-Host "  Built: bin/unosdk-amd64.exe" -ForegroundColor Green
} else {
    Write-Host "  Build failed!" -ForegroundColor Red
    exit 1
}

# Build for Windows 386
Write-Host "  → Building for Windows 386..." -ForegroundColor Gray
$env:GOOS = "windows"
$env:GOARCH = "386"
go build -ldflags="$ldflags" -o bin/unosdk-386.exe ./cmd/unosdk

if ($LASTEXITCODE -eq 0) {
    Write-Host "  Built: bin/unosdk-386.exe" -ForegroundColor Green
} else {
    Write-Host "  Build failed!" -ForegroundColor Red
}

# Build for Windows ARM64
Write-Host "  → Building for Windows ARM64..." -ForegroundColor Gray
$env:GOOS = "windows"
$env:GOARCH = "arm64"
go build -ldflags="$ldflags" -o bin/unosdk-arm64.exe ./cmd/unosdk

if ($LASTEXITCODE -eq 0) {
    Write-Host "  Built: bin/unosdk-arm64.exe" -ForegroundColor Green
} else {
    Write-Host "  Build failed!" -ForegroundColor Red
}

# Create default copy
Copy-Item "bin/unosdk-amd64.exe" "bin/unosdk.exe"

# Create a local test layout that mirrors the installed one, so restart-free
# switching can be tested without running install.ps1:
#   bin\dev\unosdk.cmd
#   bin\dev\core\unosdk.exe
$devDir = "bin/dev"
$devCoreDir = "$devDir/core"
$wrapperSource = "scripts/unosdk.cmd"
$devLayoutReady = $false

if (Test-Path $wrapperSource) {
    New-Item -ItemType Directory -Path $devCoreDir -Force | Out-Null
    Copy-Item "bin/unosdk-amd64.exe" "$devCoreDir/unosdk.exe" -Force
    Copy-Item $wrapperSource "$devDir/unosdk.cmd" -Force
    $devLayoutReady = $true
    Write-Host "  Created wrapper test layout: bin/dev" -ForegroundColor Green
} else {
    Write-Host "  Skipped wrapper test layout ($wrapperSource not found)" -ForegroundColor Yellow
}

# Clean up environment variables
Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Build Complete!" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Binaries available in ./bin/" -ForegroundColor White
Write-Host "  - unosdk.exe (default, amd64)" -ForegroundColor Gray
Write-Host "  - unosdk-amd64.exe" -ForegroundColor Gray
Write-Host "  - unosdk-386.exe" -ForegroundColor Gray
Write-Host "  - unosdk-arm64.exe" -ForegroundColor Gray
if ($devLayoutReady) {
    Write-Host "  - dev\unosdk.cmd + dev\core\unosdk.exe (wrapper layout for testing)" -ForegroundColor Gray
}
Write-Host ""
Write-Host "Run: .\bin\unosdk.exe --help" -ForegroundColor Yellow
if ($devLayoutReady) {
    Write-Host "Test restart-free switching in cmd:  .\bin\dev\unosdk.cmd switch java openjdk 21" -ForegroundColor Yellow
}