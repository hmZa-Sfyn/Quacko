# configure.ps1
# Installs quackolang.exe and retro.exe to C:\Program Files\Quacko\bin and adds it to system PATH

# Check for administrative privileges
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Error "This script requires administrative privileges. Run PowerShell as Administrator."
    exit 1
}

# Define paths
$sourceDir = "."
$sourceExeQuacko = Join-Path $sourceDir "quackolang.exe"
$sourceExeRetro = Join-Path $sourceDir "retro.exe"
$installDir = "C:\Program Files\Quacko\bin"
$installExeQuacko = Join-Path $installDir "quackolang.exe"
$installExeRetro = Join-Path $installDir "retro.exe"

# Check if source executables exist
$missingFiles = @()
if (-not (Test-Path $sourceExeQuacko)) {
    $missingFiles += $sourceExeQuacko
}
if (-not (Test-Path $sourceExeRetro)) {
    $missingFiles += $sourceExeRetro
}
if ($missingFiles.Count -gt 0) {
    Write-Error "The following files were not found: $($missingFiles -join ', '). Ensure they exist in $sourceDir."
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

# Copy executables to installation directory
try {
    if (Test-Path $installExeQuacko) {
        Write-Warning "Overwriting existing $installExeQuacko."
    }
    Copy-Item -Path $sourceExeQuacko -Destination $installExeQuacko -Force
    Write-Host "Copied quackolang.exe to $installExeQuacko"

    if (Test-Path $installExeRetro) {
        Write-Warning "Overwriting existing $installExeRetro."
    }
    Copy-Item -Path $sourceExeRetro -Destination $installExeRetro -Force
    Write-Host "Copied retro.exe to $installExeRetro"
} catch {
    Write-Error "Failed to copy executables to $installDir. Error: $_"
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
$failedInstalls = @()
if (-not (Test-Path $installExeQuacko)) {
    $failedInstalls += "quackolang.exe not found at $installExeQuacko"
}
if (-not (Test-Path $installExeRetro)) {
    $failedInstalls += "retro.exe not found at $installExeRetro"
}
if ($failedInstalls.Count -gt 0) {
    Write-Error "Installation failed: $($failedInstalls -join '; ')"
    exit 1
} else {
    Write-Host "Installation successful. Run 'quackolang' or 'retro' from a new cmd window."
}