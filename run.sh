#!/bin/bash

# Kill background processes when script exits
trap "kill 0" EXIT

# 1. Run Tailwind (Background)
./tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch &

# 2. Run Templ (Background)
templ generate --watch &

# 3. Run Go Server (Foreground)
# Wait a second for templ/css to generate first
sleep 1
go run cmd/server/main.go