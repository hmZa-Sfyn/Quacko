# configure.ps1
# Installs quackolang.exe to C:\Program Files\Quacko\bin and adds it to system PATH

# Check for administrative privileges
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Error "This script requires administrative privileges. Run PowerShell as Administrator."
    exit 1
}

# Define paths
$sourceExe = ".\quackolang.exe"  # Relative to D:\Quacko\windows_bin
$installDir = "C:\Program Files\Quacko\bin"
$installExe = Join-Path $installDir "quackolang.exe"

# Check if quackolang.exe exists
if (-not (Test-Path $sourceExe)) {
    Write-Error "quackolang.exe not found at $sourceExe. Ensure it exists in the windows_bin directory."
    exit 1
}

# Create installation directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    try {
        New-Item -Path $installDir -ItemType Directory -Force | Out-Null
        Write-Host "Created directory: $installDir"
    } catch {
        Write-Error "Failed to create directory $installDir. Error: $_"
        exit 1
    }
}

# Copy quackolang.exe to installation directory
try {
    if (Test-Path $installExe) {
        Write-Warning "Overwriting existing $installExe."
    }
    Copy-Item -Path $sourceExe -Destination $installExe -Force
    Write-Host "Copied quackolang.exe to $installExe"
} catch {
    Write-Error "Failed to copy $sourceExe to $installExe. Error: $_"
    exit 1
}

# Add installation directory to system PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::Machine)
if ($currentPath -notlike "*$installDir*") {
    try {
        $updatedPath = "$currentPath;$installDir"
        [Environment]::SetEnvironmentVariable("Path", $updatedPath, [EnvironmentVariableTarget]::Machine)
        Write-Host "Added $installDir to system PATH"
    } catch {
        Write-Error "Failed to update system PATH. Error: $_"
        exit 1
    }
} else {
    Write-Host "$installDir is already in system PATH"
}

# Verify installation
if (Test-Path $installExe) {
    Write-Host "Installation successful. Run 'quackolang' from a new cmd window."
} else {
    Write-Error "Installation failed. quackolang.exe not found at $installExe."
    exit 1
}
