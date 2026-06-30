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

- **Go 1.24+**
- **Ebitengine v2.9.9** — imported as `github.com/hajimehoshi/ebiten/v2`
- No direct dependencies beyond Ebitengine. Indirect deps (`golang.org/x/image`, `golang.org/x/sync`) are pulled in automatically by `go mod tidy`.

## Prerequisites

- **Go 1.24+** (Ebitengine 2.9 requires it). Verify with:
  ```bash
  go version
  ```
- **No C compiler needed on Windows** — Ebitengine is pure Go on Windows and works out of the box without CGO / MinGW.
- **macOS** — install the Xcode Command Line Tools (`xcode-select --install`).
- **Linux** — install `gcc` plus the standard graphics and audio development headers (e.g. `libgl1-mesa-dev`, `libxrandr-dev`, `libxinerama-dev`, `libxi-dev`, `libxxf86vm-dev`, `libasound2-dev` on Debian/Ubuntu).

## Run (desktop)

From the `pool-go/` directory:

```bash
go mod tidy   # first time only — downloads Ebitengine and indirect deps
go run .
```

This opens an **800×600** window titled **"Pool - Go"** with a racked table ready to break. Close the window to quit.

## Build a binary

```bash
go build -o pool-go .
```

Produces a single static binary (`pool-go.exe` on Windows, `pool-go` on macOS/Linux). Run it directly:

```bash
./pool-go        # macOS / Linux
./pool-go.exe    # Windows
```

## Web build (WASM)

Two options.

**Option A — instant in-browser preview** (auto-rebuilds on file changes):

```bash
go run github.com/hajimehoshi/wasmserve@latest .
```

Then open <http://localhost:8080> in a modern browser.

**Option B — produce a `.wasm` artifact manually:**

```bash
GOOS=js GOARCH=wasm go build -o pool.wasm .
```

Then copy `wasm_exec.js` from your `GOROOT`:

```bash
# from inside pool-go/
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

Serve the directory with any static HTTP server (`python -m http.server`, `npx serve`, etc.) and load `index.html` (you'll need a small `index.html` shell that loads `wasm_exec.js` and instantiates the `pool.wasm` module).

## File layout

```
pool-go/
├── go.mod        # module declaration + Ebitengine v2.9.9
├── main.go       # entry point: window setup (resizable) + RunGame
├── game.go       # Game struct, state machine, input (aiming/shooting/ball-in-hand), spin dial
├── render.go     # all drawing: table, shaded balls, ghost-ball aim, cue stick, power meter, HUD
├── physics.go    # integration, rolling friction, spin/English, ball/rail collisions, pocket jaws
├── sprites.go    # lazily-built ball/shadow/specular sprites (shaded spheres)
├── fx.go         # particles, screen shake, pocket-drop animation
├── audio.go      # procedurally synthesized sound effects (no asset files)
├── rules.go      # 8-ball rules: groups, fouls, win/loss
├── rack.go       # opening rack + cue-ball placement
├── ball.go       # Ball model (position, spin, roll) and colors
├── table.go      # table geometry and pocket positions
├── vec.go        # 2D vector math
├── game_test.go  # tests for rack, collisions, rails, pockets, friction, spin, jaws, aim
├── .gitignore
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
