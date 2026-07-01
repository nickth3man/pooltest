# AGENTS.md — pool-go

2D pool/billiards game — Go + Ebitengine (desktop + WebAssembly).

## Stack

- **Language:** Go 1.25+ (Ebitengine 2.9 floor; tested go1.26.4)
- **Engine:** Ebitengine v2.9.9 — module `github.com/hajimehoshi/ebiten/v2`
- **Module:** `github.com/user/pooltest/pool-go`; entry `main.go` (window setup + `RunGame`). The `Game` type (`ebiten.Game`: `Update() error`, `Draw(*ebiten.Image)`, `Layout(int,int) (int,int)`) lives in `internal/game/game.go`.
- **Layout:** `main.go` (package main, entry only) + `internal/` packages — `vec` (2D vector math), `ball` (Ball model + palette), `table` (geometry + pockets), `rack` (opening rack), `physics` (integration/rolling friction/spin/collisions/pockets/jaws + `On*Impact` hook vars), `audio` (procedural SFX), `sprites` (shaded ball/shadow/specular sprites), `rules` (8-ball rules + group logic), `fx` (particles/shake state), `game` (Game struct, state machine, input, spin dial, rendering, `firstContact`). Each `internal/` package is Go-toolchain-private to this module.
- **Direct dependencies:** `github.com/hajimehoshi/ebiten/v2`, `golang.org/x/image` (the `gofont/goregular` TTF is used for ball numbers and the HUD via `text/v2`).
- **Audio:** uses Ebitengine's `audio` package with PCM synthesized in code at startup (signed 16-bit LE stereo at 44100 Hz) — no sound asset files. This pulls in `github.com/ebitengine/oto/v3` (indirect).
- **Indirect dependencies** (resolved by `go mod tidy`): `golang.org/x/sync`, `golang.org/x/sys`, `golang.org/x/text`, `github.com/go-text/typesetting`, `github.com/rivo/uniseg`, `github.com/ebitengine/gomobile`, `github.com/ebitengine/hideconsole`, `github.com/ebitengine/oto/v3`, `github.com/ebitengine/purego`, `github.com/jezek/xgb`
- **Targets:** native desktop binary (Windows/macOS/Linux), WebAssembly (`GOOS=js GOARCH=wasm`)
- **Windows:** pure-Go — **no C compiler / CGO / MinGW required**
- **API note (v2.9.x):** `vector.FillCircle(dst *ebiten.Image, cx, cy, r float32, clr color.Color, antialias bool)` — the trailing `antialias bool` is required
- **Render loop:** update rate set via `ebiten.SetTPS(60)` (decoupled from draw rate); logical size returned by `Layout`
- **Conventions:** structured logging via `log/slog` (not `log`); randomness via `math/rand/v2` (auto-seeded, not `math/rand`); formatting via `gofumpt` (strict gofmt superset) + `goimports`. Fatal startup errors use `slog.Error` + `os.Exit(1)`.
- **Project config:** `Makefile` (all workflows), `.golangci.yml` v2 (27 linters), `.editorconfig`, `.gitattributes` (LF on `.go`, CRLF on `.bat`), `.github/workflows/ci.yml` (vet+lint+test+wasm smoke on push/PR), `.github/dependabot.yml` (weekly grouped updates), `LICENSE` (Apache-2.0).

## Dev commands

Run from the `pool-go/` folder. The `Makefile` wraps every common workflow — run `make help` for the full list.

### Make targets (preferred)
- `make run` — build + run the desktop binary
- `make build` — native binary (`pool-go.exe` on Windows, `pool-go` elsewhere)
- `make build-cross` — cross-compile for linux/darwin/windows (amd64 + darwin arm64)
- `make wasm` — produce `pool.wasm` (also tells you how to copy `wasm_exec.js`)
- `make wasm-serve` — instant dev preview at http://localhost:8080 (via `wasmserve`)
- `make test` — run all tests; `make test-race` — with the race detector (needs CGO/gcc)
- `make cover` — generate `cover.html` with line-level coverage
- `make vet` — `go vet`
- `make lint` — `golangci-lint` v2 (27 linters across style, bug-finding, modernization, and complexity — see `.golangci.yml` for the full enable list; key ones: revive, gocritic, govet, staticcheck, modernize, perfsprint, errcheck, bodyclose, nilnil, gocyclo, intrange)
- `make fmt` — `gofumpt -w` + `goimports -w`
- `make check` — vet + lint + test in one shot (the full quality gate; what CI runs)
- `make mod-tidy` — tidy go.mod/go.sum
- `make clean` — remove all build artifacts

### Raw go commands (equivalents)
- `go mod tidy` — resolve/add indirect deps to go.mod/go.sum
- `go get -u ./...` — update deps; `go get <pkg>@<version>` — add/pin
- `go list -m all` — list modules; `go mod graph` — dependency graph
- `go clean -modcache` — clear module cache

### Build / run
- `go run .` — build + run
- `go build -o pool-go.exe .` — Windows binary; `go build -o pool-go .` — macOS/Linux
- `go install` — install to GOBIN

### Quality
- `go vet ./...` — static checks
- `gofumpt -w .` — format; `gofumpt -l .` — list unformatted (gofumpt is a strict superset of gofmt)
- `go test ./...` — tests; `go test -v ./...` — verbose; `go test -race ./...` — race detector

### Web (WebAssembly)
- Instant dev preview: `make wasm-serve` (or `go run github.com/hajimehoshi/wasmserve@latest .`) → http://localhost:8080
- Produce a `.wasm`: `make wasm` (or `GOOS=js GOARCH=wasm go build -o pool.wasm .`), then copy `$(go env GOROOT)/lib/wasm/wasm_exec.js` next to it and serve via any static host

### Cleanup
- `make clean` — remove build artifacts; `go clean -cache` — clear build cache

No CGO/MinGW needed on Windows for builds or tests. The race detector (`make test-race`) does require CGO+gcc — use plain `make test` locally on Windows; CI runs the race detector on Linux. On macOS install Xcode CLT; on Linux install `gcc` + graphics/audio dev headers.
