#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

VER="$(tr -d '[:space:]' < VERSION)"
if [ -z "$VER" ]; then
  echo "VERSION file is empty" >&2
  exit 1
fi
echo "version: $VER"

mkdir -p dist/windows dist/linux cmd/setup/payload

echo "running tests"
go test ./...
go vet ./...
GOOS=windows GOARCH=amd64 go vet ./...

echo "sync web assets"
cp web/index.html web/app.css web/app.js cmd/panel/web/

echo "linux panel (preview)"
go build -trimpath \
  -ldflags "-s -w -X main.appVersion=$VER" \
  -o dist/linux/StellarisConsolePanel ./cmd/panel

echo "windows panel exe"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
  -ldflags "-s -w -H windowsgui -X main.appVersion=$VER" \
  -o dist/windows/StellarisConsolePanel.exe ./cmd/panel

echo "stage installer payload"
cp dist/windows/StellarisConsolePanel.exe cmd/setup/payload/StellarisConsolePanel.exe
cp README.md cmd/setup/payload/README.md

echo "windows setup exe"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
  -ldflags "-s -w -X main.appVersion=$VER" \
  -o dist/windows/Setup.exe ./cmd/setup

echo "done"
ls -lh dist/windows dist/linux
echo
echo "安装：在 Win10 双击 dist/windows/Setup.exe（会弹出 UAC）"
echo "静默：Setup.exe /S"
echo "卸载：开始菜单「卸载 群星控制台面板」或 设置-应用"
