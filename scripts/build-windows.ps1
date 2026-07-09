# Build the File Tools Windows installer locally (run on Windows in the repo root):
#   powershell -ExecutionPolicy Bypass -File scripts\build-windows.ps1
#
# Requirements: Go, a C compiler (MinGW gcc), and NSIS.
# The script installs the Fyne tool and NSIS (via Chocolatey) if missing.

$ErrorActionPreference = "Stop"

if (-not (Get-Command gcc -ErrorAction SilentlyContinue)) {
    Write-Error "A C compiler (gcc) is required. Install MinGW-w64, e.g. 'choco install mingw -y'."
}

Write-Host "Installing Fyne packaging tool..."
go install fyne.io/fyne/v2/cmd/fyne@latest
$env:PATH = "$(go env GOPATH)\bin;$env:PATH"

Write-Host "Packaging FileTools.exe..."
$env:CGO_ENABLED = "1"
$icon = Join-Path (Get-Location) "build\appicon.png"
fyne package --os windows --src ./cmd/filetools --icon "$icon" --name FileTools --app-id ai.filetools.desktop --release
$exe = Get-ChildItem -Recurse -Filter FileTools.exe | Select-Object -First 1
if (-not $exe) { throw "FileTools.exe was not produced" }
$dest = Join-Path (Get-Location) "FileTools.exe"
if ($exe.FullName -ne $dest) { Move-Item -Force $exe.FullName $dest }

$makensis = "C:\Program Files (x86)\NSIS\makensis.exe"
if (-not (Test-Path $makensis)) {
    Write-Host "Installing NSIS..."
    choco install nsis -y
}

Write-Host "Building installer..."
& $makensis build\installer.nsi
if (-not (Test-Path "build\FileToolsSetup.exe")) { throw "installer was not produced" }

Write-Host "Done. Installer at build\FileToolsSetup.exe"
