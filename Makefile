# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go test -race -vet=off ./...
	go mod verify


# ==================================================================================== #
# BUILD
# ==================================================================================== #
## Build settings and helper (DRYed build command)
LDFLAGS := -s -w -X main.version=9.9.9
OUT ?= ./bin/mods-fixed-pipe

.PHONY: _build
_build:
	go mod verify
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -a -o $(OUT) .

## build: build the cmd application (default: macOS Apple Silicon)
.PHONY: build
build: GOOS=darwin GOARCH=arm64 OUT=./bin/mods-fixed-pipe
build: _build

## build-local: detect host OS/ARCH and build for it (outputs to ./bin/mods-fixed-pipe)
.PHONY: build-local
build-local:
	OS=$$(uname -s); ARCH=$$(uname -m); \
	case $$OS in \
		Darwin) GOOS=darwin ;; \
		Linux) GOOS=linux ;; \
		*) echo "unsupported OS: $$OS"; exit 1 ;; \
	esac; \
	case $$ARCH in \
		x86_64|amd64) GOARCH=amd64 ;; \
		arm64|aarch64) GOARCH=arm64 ;; \
		*) echo "unsupported arch: $$ARCH"; exit 1 ;; \
	esac; \
	echo "Building for $$GOOS/$$GOARCH -> ./bin/mods-fixed-pipe"; \
	$(MAKE) _build GOOS=$$GOOS GOARCH=$$GOARCH OUT=./bin/mods-fixed-pipe

## build-linux: build the cmd application for Linux (amd64)
.PHONY: build-linux
build-linux: GOOS=linux GOARCH=amd64 OUT=./bin/mods-fixed-pipe-linux
build-linux: _build

## build-darwin-amd64: build the cmd application for macOS (Intel)
.PHONY: build-darwin-amd64
build-darwin-amd64: GOOS=darwin GOARCH=amd64 OUT=./bin/mods-fixed-pipe-darwin
build-darwin-amd64: _build

## run: run the cmd application
.PHONY: run
run: build
	./bin/mods-fixed-pipe

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
