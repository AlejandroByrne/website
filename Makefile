SHELL := /bin/bash
APP := website
PKG := ./cmd/server
IMAGE := $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPO_NAME)/$(APP):$(TAG)
TAG ?= dev

.PHONY: dev
dev:
	./scripts/dev.sh

.PHONY: build
build:
	templ generate
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css
	go build -o bin/server $(PKG)

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE) .

.PHONY: test
test:
	go vet ./...
	go test ./... -count=1 -run . || true
