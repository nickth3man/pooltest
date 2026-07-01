# pool-go Makefile.
# Standard targets: run / build / wasm / test / vet / lint / fmt / cover / clean.
# Cross-platform pure-Go; works on Windows (Git Bash), macOS, and Linux.
# `go build` with no -o produces pool-go (Linux/macOS) or pool-go.exe (Windows)
# automatically, so we don't need to detect the extension.

LDFLAGS    = -s -w
MODULE     = github.com/user/pooltest/pool-go

.PHONY: all
all: build

.PHONY: run
run:  ## Build and run the desktop binary
	go run .

.PHONY: build
build:  ## Build the native desktop binary
	go build -ldflags="$(LDFLAGS)" .

.PHONY: build-cross
build-cross:  ## Cross-compile native binaries for linux/darwin/windows
	GOOS=linux   GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o pool-go-linux-amd64   .
	GOOS=darwin  GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o pool-go-darwin-amd64  .
	GOOS=darwin  GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o pool-go-darwin-arm64  .
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o pool-go-windows-amd64.exe .

.PHONY: wasm
wasm:  ## Build the WebAssembly binary (pool.wasm)
	GOOS=js GOARCH=wasm go build -trimpath -ldflags="$(LDFLAGS)" -buildvcs=false -o pool.wasm .
	@echo "WASM built. To serve, also copy wasm_exec.js from your GOROOT:"
	@echo '  cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" .'

.PHONY: wasm-serve
wasm-serve:  ## Run a local server that hot-reloads the WASM build (http://localhost:8080)
	go run github.com/hajimehoshi/wasmserve@latest .

.PHONY: test
test:  ## Run all tests
	go test -count=1 ./...

.PHONY: test-race
test-race:  ## Run all tests with the race detector (requires CGO/gcc)
	go test -race -count=1 ./...

.PHONY: cover
cover:  ## Generate HTML coverage report at cover.html
	go test -coverprofile=cover.out -count=1 ./...
	go tool cover -html=cover.out -o cover.html
	@echo "Open cover.html to view coverage."

.PHONY: vet
vet:  ## Run go vet
	go vet ./...

.PHONY: lint
lint:  ## Run golangci-lint
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...

.PHONY: fmt
fmt:  ## Format all Go source (gofumpt + goimports)
	go run mvdan.cc/gofumpt@latest -w .
	go run golang.org/x/tools/cmd/goimports@latest -w -local $(MODULE) .

.PHONY: mod-tidy
mod-tidy:  ## Tidy go.mod / go.sum
	go mod tidy

.PHONY: clean
clean:  ## Remove build artifacts
	rm -f pool-go pool-go.exe
	rm -f pool-go-linux-amd64 pool-go-darwin-amd64 pool-go-darwin-arm64 pool-go-windows-amd64.exe
	rm -f pool.wasm cover.out cover.html

.PHONY: check
check: vet lint test  ## Run every quality gate (vet + lint + test)

.PHONY: help
help:  ## Show this help
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-15s %s\n", $$1, $$2}' Makefile
