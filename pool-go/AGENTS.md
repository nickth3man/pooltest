# AGENTS.md — pool-go

2D pool/billiards game — Go + Ebitengine (desktop + WebAssembly).

## Stack

- **Language:** Go 1.24+ (Ebitengine 2.9 floor; tested go1.26.4)
- **Engine:** Ebitengine v2.9.9 — module `github.com/hajimehoshi/ebiten/v2`
- **Module:** `github.com/user/pooltest/pool-go`; entry `main.go` (window setup + `RunGame`). The `Game` type (`ebiten.Game`: `Update() error`, `Draw(*ebiten.Image)`, `Layout(int,int) (int,int)`) lives in `game.go`.
- **Source files** (flat `package main`): `game.go` (state machine + input + spin dial), `render.go` (drawing + ghost-ball aim), `physics.go` (integration/rolling friction/spin/collisions/pockets/jaws), `sprites.go` (shaded ball/shadow/specular sprites), `fx.go` (particles/shake/drop animation), `audio.go` (procedural SFX), `rules.go` (8-ball rules), `rack.go`, `ball.go`, `table.go`, `vec.go`, `game_test.go`.
- **Direct dependencies:** `github.com/hajimehoshi/ebiten/v2`, `golang.org/x/image` (the `gofont/goregular` TTF is used for ball numbers and the HUD via `text/v2`).
- **Audio:** uses Ebitengine's `audio` package with PCM synthesized in code at startup (signed 16-bit LE stereo at 44100 Hz) — no sound asset files. This pulls in `github.com/ebitengine/oto/v3` (indirect).
- **Indirect dependencies** (resolved by `go mod tidy`): `golang.org/x/sync`, `golang.org/x/sys`, `golang.org/x/text`, `github.com/go-text/typesetting`, `github.com/rivo/uniseg`, `github.com/ebitengine/gomobile`, `github.com/ebitengine/hideconsole`, `github.com/ebitengine/oto/v3`, `github.com/ebitengine/purego`, `github.com/jezek/xgb`
- **Targets:** native desktop binary (Windows/macOS/Linux), WebAssembly (`GOOS=js GOARCH=wasm`)
- **Windows:** pure-Go — **no C compiler / CGO / MinGW required**
- **API note (v2.9.x):** `vector.FillCircle(dst *ebiten.Image, cx, cy, r float32, clr color.Color, antialias bool)` — the trailing `antialias bool` is required
- **Render loop:** update rate set via `ebiten.SetTPS(60)` (decoupled from draw rate); logical size returned by `Layout`

## Dev commands

Run from the `pool-go/` folder.

### Dependencies
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
- `gofmt -w .` — format; `gofmt -l .` — list unformatted
- `go test ./...` — tests; `go test -v ./...` — verbose; `go test -race ./...` — race detector

### Web (WebAssembly)
- Instant dev preview: `go run github.com/hajimehoshi/wasmserve@latest .` → http://localhost:8080
- Produce a `.wasm`: `GOOS=js GOARCH=wasm go build -o pool.wasm .`, then copy `$(go env GOROOT)/lib/wasm/wasm_exec.js` next to it and serve via any static host

### Cleanup
- `go clean -cache` — clear build cache

No CGO/MinGW needed on Windows. On macOS install Xcode CLT; on Linux install `gcc` + graphics/audio dev headers.
