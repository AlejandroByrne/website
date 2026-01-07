# ==========================================
# STAGE 1: The Builder (Compiling the App)
# ==========================================
FROM golang:alpine AS builder

# 1. Install required system tools
# 'curl' is needed to download Tailwind
RUN apk add --no-cache curl git

# 2. Set working directory inside the container
WORKDIR /app

# 3. Install Templ (The Go Code Generator)
RUN go install github.com/a-h/templ/cmd/templ@latest

# 4. Install Tailwind CSS (The Linux Version!)
# Note: We download the LINUX binary because this container runs Linux, even on your Mac.
RUN curl -sL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 -o /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss

# 5. Copy Dependency Files
COPY go.mod go.sum ./
RUN go mod download

# 6. Copy Source Code
COPY . .

# 7. Generate Templ Files
RUN templ generate

# 8. Build CSS
# We run the linux binary we just downloaded
RUN tailwindcss -i ./static/css/input.css -o ./static/css/output.css

# 9. Build the Go Binary
# CGO_ENABLED=0 ensures a static binary that works on any Linux distro
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/server/main.go


# ==========================================
# STAGE 2: The Runner (Production Image)
# ==========================================
FROM alpine:latest

WORKDIR /root/

# 1. Copy the Binary from the Builder Stage
COPY --from=builder /main .

# 2. Copy the Static Assets (CSS, Images) from the Builder Stage
# We need these so the browser can load styles
COPY --from=builder /app/static ./static

# 3. Expose the Port
EXPOSE 8080

# 4. Run the Binary
CMD ["./main"]