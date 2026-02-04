$ErrorActionPreference = "Stop"

# Change to project root directory (parent of scripts)
$scriptDir = $PSScriptRoot
$projectRoot = Split-Path $scriptDir -Parent
Set-Location $projectRoot

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "  unosdk - Build Script" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# Get version from git or use default
try {
    $version = git describe --tags --always --dirty 2>$null
    if ($LASTEXITCODE -ne 0) {
        $version = "dev"
    }
} catch {
    $version = "dev"
}

$commit = try { git rev-parse --short HEAD 2>$null } catch { "unknown" }
$date = Get-Date -Format "yyyy-MM-dd"

Write-Host "Version: $version" -ForegroundColor Green
Write-Host "Commit:  $commit" -ForegroundColor Green
Write-Host "Date:    $date" -ForegroundColor Green
Write-Host ""

# Build flags
$ldflags = "-s -w -X main.version=$version -X main.commit=$commit -X main.date=$date"

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
    Write-Host "  ✓ Built: bin/unosdk-amd64.exe" -ForegroundColor Green
} else {
    Write-Host "  ✗ Build failed!" -ForegroundColor Red
    exit 1
}

# Build for Windows 386
Write-Host "  → Building for Windows 386..." -ForegroundColor Gray
$env:GOOS = "windows"
$env:GOARCH = "386"
go build -ldflags="$ldflags" -o bin/unosdk-386.exe ./cmd/unosdk

if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Built: bin/unosdk-386.exe" -ForegroundColor Green
} else {
    Write-Host "  ✗ Build failed!" -ForegroundColor Red
}

# Build for Windows ARM64
Write-Host "  → Building for Windows ARM64..." -ForegroundColor Gray
$env:GOOS = "windows"
$env:GOARCH = "arm64"
go build -ldflags="$ldflags" -o bin/unosdk-arm64.exe ./cmd/unosdk

if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Built: bin/unosdk-arm64.exe" -ForegroundColor Green
} else {
    Write-Host "  ✗ Build failed!" -ForegroundColor Red
}

# Create default copy
Copy-Item "bin/unosdk-amd64.exe" "bin/unosdk.exe"

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
Write-Host ""
Write-Host "Run: .\bin\unosdk.exe --help" -ForegroundColor Yellow