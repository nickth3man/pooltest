# Pool — Go

A real-time 2D **8-ball pool** game written in Go using [Ebitengine v2.9.x](https://ebitengine.org/).

A full two-player game: a racked table of 16 balls, a slingshot cue, ball/rail/pocket physics with spin, and 8-ball rules (group assignment, scratch-with-ball-in-hand, and win/loss on the 8).

## How to play

- **Aim & shoot:** hold the left mouse button, drag *back* from the cue ball, and release. The cue ball fires in the opposite direction with power proportional to how far you pulled. A **ghost ball** marks where you'll first make contact, a faint line shows the object ball's predicted carom, and a side **power meter** ramps green → red with strength.
- **Spin (English):** click the **English dial** in the bottom-left to set the contact point on the cue ball — top/bottom for *follow*/*draw*, left/right for side spin. Right-click anywhere to clear it. Watch the cue ball's spot spin when you apply side english.
- **Groups:** the table is open until the first ball is legally pocketed, which assigns solids/stripes. Pocket one of your own balls to keep shooting.
- **Scratch:** pocketing the cue ball is a foul — the opponent gets ball-in-hand (click inside the cushions to place it).
- **Winning:** clear your group, then pocket the 8-ball. Sinking the 8 early — or scratching on it — loses the game. Press **R** to start a new match.

## Look & feel

Shaded sphere balls with drop shadows and a fixed specular highlight, a felt sheen with rail bevels, diamond sights and table spots, deep pockets with cushion-tip **jaws** that rattle off-angle shots, pocket-drop animations, impact sparks, a small break shake, and synthesized (no asset files) sound for contacts, rails, pockets, and the cue strike.

## Stack

- **Go 1.25+** (Ebitengine 2.9 floor; tested with go1.26.4)
- **Ebitengine v2.9.9** — imported as `github.com/hajimehoshi/ebiten/v2`
- **Direct dependencies:** `golang.org/x/image` (the `gofont/goregular` TTF for ball numbers and the HUD via `text/v2`)
- **Indirect dependencies** (resolved by `go mod tidy`): `golang.org/x/text`, `github.com/go-text/typesetting`, `github.com/rivo/uniseg` (font shaping), `github.com/ebitengine/oto/v3` (audio output for procedurally synthesized SFX), plus `golang.org/x/sync`, `golang.org/x/sys`, `github.com/ebitengine/gomobile`, `github.com/ebitengine/hideconsole`, `github.com/ebitengine/purego`, and `github.com/jezek/xgb`

## Prerequisites

- **Go 1.25+** (Ebitengine 2.9 requires it). Verify with:
  ```bash
  go version
  ```
- **No C compiler needed on Windows** — Ebitengine is pure Go on Windows and works out of the box without CGO / MinGW.
- **macOS** — install the Xcode Command Line Tools (`xcode-select --install`).
- **Linux** — install `gcc` plus the standard graphics and audio development headers (e.g. `libgl1-mesa-dev`, `libxrandr-dev`, `libxinerama-dev`, `libxi-dev`, `libxxf86vm-dev`, `libasound2-dev` on Debian/Ubuntu).

## Dev workflow

Run from the `pool-go/` directory. The [`Makefile`](Makefile) wraps every common workflow — run `make help` for the full list.

| Target | Description |
|--------|-------------|
| `make run` | Build and run the desktop binary |
| `make build` | Build the native desktop binary |
| `make build-cross` | Cross-compile for linux/darwin/windows (amd64 + darwin arm64) |
| `make wasm` | Produce `pool.wasm` (also prints how to copy `wasm_exec.js`) |
| `make wasm-serve` | Instant dev preview at http://localhost:8080 (via `wasmserve`) |
| `make test` | Run all tests |
| `make test-race` | Run tests with the race detector (requires CGO/gcc) |
| `make cover` | Generate `cover.html` with line-level coverage |
| `make vet` | Run `go vet` |
| `make lint` | Run golangci-lint v2 (31 linters) |
| `make fmt` | Format with `gofumpt` + `goimports` |
| `make check` | Full quality gate: vet + lint + test |
| `make mod-tidy` | Tidy `go.mod` / `go.sum` |
| `make clean` | Remove build artifacts |

### Run (desktop)

```bash
make mod-tidy   # first time only — downloads deps
make run
```

This opens an **800×600** resizable window titled **"Pool - Go"** with a racked table ready to break. Close the window to quit.

Manual equivalent:

```bash
go mod tidy
go run .
```

### Build a binary

```bash
make build
```

Produces a single static binary (`pool-go.exe` on Windows, `pool-go` on macOS/Linux). Run it directly:

```bash
./pool-go        # macOS / Linux
./pool-go.exe    # Windows
```

Manual equivalent:

```bash
go build -o pool-go .
```

## Web build (WASM)

Two options.

**Option A — instant in-browser preview** (auto-rebuilds on file changes):

```bash
make wasm-serve
```

Then open <http://localhost:8080> in a modern browser.

Manual equivalent:

```bash
go run github.com/hajimehoshi/wasmserve@latest .
```

**Option B — produce a `.wasm` artifact manually:**

```bash
make wasm
```

Manual equivalent:

```bash
GOOS=js GOARCH=wasm go build -o pool.wasm .
```

Then copy `wasm_exec.js` from your `GOROOT`:

```bash
# from inside pool-go/
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
```

Serve the directory with any static HTTP server (`python -m http.server`, `npx serve`, etc.) and load `index.html` (you'll need a small `index.html` shell that loads `wasm_exec.js` and instantiates the `pool.wasm` module).

## Quality / CI

`make check` runs the full local quality gate: `go vet`, golangci-lint, and all tests. This mirrors [`.github/workflows/ci.yml`](.github/workflows/ci.yml), which runs on every push and pull request — plus a native build smoke test and a WASM build smoke test.

On Windows, use plain `make test` locally; the race detector (`make test-race`) requires CGO/gcc. CI runs the race detector on Linux.

## File layout

```
pool-go/
├── go.mod, go.sum       # module + deps
├── main.go              # entry point: window setup + RunGame
├── Makefile             # all dev workflows (run/build/wasm/test/lint/fmt/check)
├── .golangci.yml        # golangci-lint v2 config (31 linters)
├── .github/workflows/   # CI: vet + lint + test + wasm smoke on push/PR
├── internal/
│   ├── vec/             # 2D vector math
│   ├── ball/            # Ball model + palette
│   ├── table/           # table geometry + pockets
│   ├── rack/            # opening rack
│   ├── physics/         # integration/friction/spin/collisions/pockets/jaws
│   ├── audio/           # procedural SFX (synthesized, no asset files)
│   ├── sprites/         # shaded ball/shadow/specular sprites
│   ├── rules/           # 8-ball rules: groups, fouls, win/loss
│   ├── fx/              # particles, screen shake, pocket-drop animation
│   └── game/            # Game struct, state machine, input, rendering
├── LICENSE              # Apache-2.0
└── README.md
```

## Rules — what's modelled (and simplified)

Implemented: open-table group assignment, continue-on-legal-pocket, scratch → ball-in-hand, and 8-ball win/loss.

Deliberately simplified versus tournament 8-ball (room for future iteration):

- Shots are not "called" — any legally pocketed ball of your group counts.
- The only foul is scratching the cue ball (no rail-contact or wrong-ball-first fouls).
- Clearing your group and sinking the 8 on the *same* stroke counts as a win.

Shipped since the first cut: shaded ball sprites, spin/English with a ghost-ball aim line, pocket jaws, rolling-friction physics, drop/spark/shake juice, and synthesized audio.

Ideas next: called shots, rail-contact and wrong-ball-first fouls, a hosted WASM build, and a simple AI opponent. Keep diffs small and focused — one feature at a time.

## License

Licensed under the [Apache License 2.0](LICENSE).
