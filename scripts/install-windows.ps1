# SecretShift Windows Installation Script
# Installs or updates secret-shift to the latest version

param(
    [string]$InstallDir = "$env:ProgramFiles\secret-shift"
)

$ErrorActionPreference = "Stop"

$Repo = "PapaDanielVi/secret-shift"

# Detect architecture
$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
    "AMD64" { $Suffix = "x86_64" }
    "ARM64" { $Suffix = "arm64" }
    default { Write-Host "Unsupported architecture: $Arch"; exit 1 }
}

# Get latest release version
Write-Host "Fetching latest release..."
$LatestVersion = (Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest").tag_name
if (-not $LatestVersion) {
    Write-Host "Failed to fetch latest release version"
    exit 1
}

Write-Host "Installing secret-shift $LatestVersion..."

# Remove existing installation if present
$BinPath = Join-Path $InstallDir "secret-shift.exe"
if (Test-Path $BinPath) {
    Write-Host "Removing existing installation..."
    Remove-Item -Force $BinPath
}

# Create install directory if needed
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# Download tarball
$Url = "https://github.com/$Repo/releases/download/$LatestVersion/secret-shift_Windows_$Suffix.tar.gz"
$TarPath = "$env:TEMP\secret-shift.tar.gz"

Write-Host "Downloading Windows tarball..."
Invoke-WebRequest -Uri $Url -OutFile $TarPath

# Extract tarball (requires tar on Windows 10+)
tar -xzf $TarPath -C $InstallDir
Remove-Item $TarPath

# Add to PATH if not already present
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($CurrentPath -notlike "*$InstallDir*") {
    $NewPath = "$CurrentPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "Machine")
    Write-Host "Added $InstallDir to system PATH. You may need to restart your terminal."
}

Write-Host "secret-shift $LatestVersion installed successfully to $InstallDir!"

# Verify installation
& "$BinPath" version