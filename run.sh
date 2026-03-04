#!/usr/bin/env bash
set -euo pipefail

# Prefer unified dev script
if [[ -x "./scripts/dev.sh" ]]; then
  exec ./scripts/dev.sh
fi

echo "scripts/dev.sh not found; falling back to legacy dev flow"

trap "kill 0" EXIT

# Try system tailwindcss
if command -v tailwindcss >/dev/null 2>&1; then
  tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch &
else
  echo "tailwindcss not found. Install it or run scripts/dev.sh"
fi

templ generate --watch &

sleep 1
exec go run cmd/server/main.go