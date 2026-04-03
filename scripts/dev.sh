#!/usr/bin/env bash
set -euo pipefail

# Ensure tailwindcss binary exists (platform-aware)
BIN_DIR="./bin"
TW_BIN="$BIN_DIR/tailwindcss"

mkdir -p "$BIN_DIR"
if [[ ! -x "$TW_BIN" ]]; then
  echo "Downloading TailwindCSS binary..."
  OS=$(uname -s)
  ARCH=$(uname -m)
  case "$OS" in
    Linux)  TW_OS=linux ;;
    Darwin) TW_OS=macos ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
  esac
  case "$ARCH" in
    x86_64|amd64) TW_ARCH=x64 ;;
    arm64|aarch64) TW_ARCH=arm64 ;;
    *) echo "Unsupported arch: $ARCH"; exit 1 ;;
  esac
  URL="https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-${TW_OS}-${TW_ARCH}"
  curl -fsSL "$URL" -o "$TW_BIN"
  chmod +x "$TW_BIN"
fi

# Ensure templ exists
if ! command -v templ >/dev/null 2>&1; then
  echo "Installing templ..."
  go install github.com/a-h/templ/cmd/templ@latest
fi

# Start watchers
"$TW_BIN" -i ./static/css/input.css -o ./static/css/output.css --watch &
TW_PID=$!

templ generate --watch &
TPL_PID=$!

# Run server
sleep 1
trap 'kill $TW_PID $TPL_PID' EXIT
exec go run cmd/server/main.go
