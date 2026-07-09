#!/usr/bin/env bash
# Build the File Tools macOS app + DMG locally (run on macOS in the repo root):
#   ./scripts/build-macos.sh
#
# Produces a universal (Intel + Apple Silicon) FileTools.app and FileTools.dmg.
# Requires Go and the Xcode command line tools.
set -euo pipefail
cd "$(dirname "$0")/.."

echo "Installing Fyne packaging tool..."
go install fyne.io/tools/cmd/fyne@latest
export PATH="$(go env GOPATH)/bin:$PATH"

echo "Building universal binary..."
export CGO_ENABLED=1
GOARCH=arm64 go build -o ft-arm64 ./cmd/filetools
GOARCH=amd64 CC="clang -arch x86_64" go build -o ft-amd64 ./cmd/filetools
lipo -create -output ft-universal ft-arm64 ft-amd64
lipo -info ft-universal

echo "Packaging FileTools.app..."
rm -rf FileTools.app
fyne package --os darwin --executable ft-universal --name FileTools --icon build/appicon.png --app-id ai.filetools.desktop --release

echo "Building FileTools.dmg..."
rm -rf dmg FileTools.dmg
mkdir -p dmg/root
cp -R FileTools.app dmg/root/
ln -s /Applications dmg/root/Applications
hdiutil create -volname "File Tools" -srcfolder dmg/root -ov -format UDZO FileTools.dmg
rm -rf dmg ft-arm64 ft-amd64 ft-universal

echo "Done: FileTools.app and FileTools.dmg"
