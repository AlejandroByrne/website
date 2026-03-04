# ==========================================
# STAGE 1: The Builder (Compiling the App)
# ==========================================
FROM golang:1.25-alpine AS builder

# 1. Install required system tools
RUN apk add --no-cache curl git ca-certificates

# 2. Set working directory inside the container
WORKDIR /app

# 3. Install Templ (The Go Code Generator)
RUN GOBIN=/usr/local/bin go install github.com/a-h/templ/cmd/templ@latest

# 4. Install Tailwind CSS (Linux binary)
RUN curl -sSL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 -o /usr/local/bin/tailwindcss \
  && chmod +x /usr/local/bin/tailwindcss

# 5. Copy Dependency Files (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# 6. Copy Source Code
COPY . .

# 7. Generate Templ Files
RUN templ generate

# 8. Build CSS
RUN tailwindcss -i ./static/css/input.css -o ./static/css/output.css

# 9. Build the Go Binary
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -trimpath -ldflags="-s -w" -o /app/bin/server ./cmd/server

# ==========================================
# STAGE 2: The Runner (Production Image)
# ==========================================
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy binary and static assets
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/static /app/static

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app/server"]