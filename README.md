# Pooltest

A multi-language billiards/pool game implemented in **7 different tech stacks** — a polyglot exploration of game development patterns, physics engines, build systems, and rendering pipelines. Each implementation ships a runnable billiards game (or scaffold) targeting the same 2D (and one 3D) pool/snooker concept, allowing direct comparisons across languages, frameworks, and architectures.

| Implementation | Language | Framework / Engine | Platform | Status |
|---------------|----------|-------------------|----------|--------|
| [`pool-go`](./pool-go) | Go | Ebitengine v2.9 | Desktop + WASM | **Full game** — physics, spin, rules, audio, particles |
| [`pool-cpp`](./pool-cpp) | C++17 | raylib 6.0 | Desktop | **Feature-complete** — physics, ball sprites, cue input |
| [`pool-ts`](./pool-ts) | TypeScript | Phaser 4 + Vite 8 | Browser | **Scaffold** — green felt, cue ball, toolchain ready |
| [`pool-python`](./pool-python) | Python 3.11+ | pygame-ce | Desktop | **Scaffold** — window, felt, exit on ESC |
| [`pool-rust`](./pool-rust) | Rust | Bevy 0.18 + Rapier2D | Desktop | **Scaffold** — ECS setup, felt, cue ball circle |
| [`pool-kotlin-2d`](./pool-kotlin-2d) | Kotlin | kool + Box2D (kool-physics-2d) | Desktop + Browser | **Scaffold** — kool scene, felt, single ball |
| [`pool-kotlin-3d`](./pool-kotlin-3d) | Kotlin | kool + NVIDIA PhysX (kool-physics) | Desktop (JVM) | **Scaffold** — 3D orbit camera, green plane, sphere |

## Running

Each implementation lives in its own directory with its own build system, README, and AGENTS.md with full dev-command reference:

```
pool-go/          cd pool-go && go run .
pool-cpp/         cd pool-cpp && cmake -B build && cmake --build build --config Release && .\build\Release\pool.exe
pool-ts/          cd pool-ts && npm install && npm run dev
pool-python/      cd pool-python && uv sync && uv run pool-python
pool-rust/        cd pool-rust && cargo run
pool-kotlin-2d/   cd pool-kotlin-2d && ./gradlew runJvm
pool-kotlin-3d/   cd pool-kotlin-3d && ./gradlew run
```

See each subdirectory's `README.md` for detailed prerequisites, build instructions, and controls.

## Implementations

### pool-go (most complete)

A full two-player 8-ball pool game with shaded ball sprites, spin/english with ghost-ball aim line, pocket jaws with realistic rattle, particle effects, screen shake, synthesized audio, and a turn-based rules engine (group assignment, scratch → ball-in-hand, 8-ball win/loss). Uses Ebitengine's pure-Go rendering — no C compiler required on Windows.

### pool-cpp

Feature-complete 2D 8-ball pool built with raylib 6.0. Includes procedural ball textures with numbers and stripes, elastic ball-ball and ball-cushion physics, rolling friction, pocket detection, mouse aim with power drag, and rack reset.

### pool-ts

Browser-based scaffold using Phaser 4's WebGL2 renderer and Vite 8's dev server with HMR. Intentionally minimal — green felt stage with a single white cue ball centered — to confirm the toolchain end-to-end before layering on physics, rails, pockets, and cue mechanics.

### pool-python

Minimal scaffold using pygame-ce (the actively maintained pygame fork). Ships an 800×600 window with green felt running at 60 FPS and clean ESC exit. Configured with the strictest stable Ruff (lint + format) and ty (type checking) settings.

### pool-rust

Scaffold using Bevy's ECS architecture with Rapier2D physics on deck. Configured with cargo-husky pre-commit hooks, Bevy lint support, and optimized dev profiles. The `dev-dynamic` feature enables fast Bevy iteration via dynamic linking.

### pool-kotlin-2d

Multiplatform scaffold targeting both JVM desktop and browser (WebGL/WebGPU) via the kool engine. Configured with Box2D 3.1.1 physics bindings (kool-physics-2d), ktlint formatting via spotless, and a Gradle version catalog.

### pool-kotlin-3d

3D scaffold using kool's Vulkan-first renderer with NVIDIA PhysX 5.6.1 physics via physx-jni. Features an orbit camera (left-drag orbit, right-drag pan, scroll zoom) and ships natives for all platforms as Maven classifier artifacts.

## Dev tooling

Each sub-project is configured with language-appropriate quality tooling:

| Language | Linter | Formatter | Type Checker | Build System |
|----------|--------|-----------|--------------|--------------|
| Go | `go vet` | `gofmt` | — | `go build` |
| C++ | — | — | — | CMake + MSVC |
| TypeScript | — | — | `tsc --noEmit` | Vite |
| Python | Ruff (ALL) | Ruff | ty | uv + uv_build |
| Rust | Clippy | rustfmt | — | Cargo |
| Kotlin 2D | ktlint (spotless) | ktlint | — | Gradle KMP |
| Kotlin 3D | — | — | — | Gradle KMP |

## Project structure

```
pooltest/
├── README.md
├── .gitignore
├── pool-go/          # Go + Ebitengine (full 8-ball game)
├── pool-cpp/         # C++17 + raylib (feature-complete)
├── pool-ts/          # TypeScript + Phaser 4 + Vite (scaffold)
├── pool-python/      # Python + pygame-ce (scaffold)
├── pool-rust/        # Rust + Bevy + Rapier2D (scaffold)
├── pool-kotlin-2d/   # Kotlin + kool + Box2D (scaffold)
└── pool-kotlin-3d/   # Kotlin + kool + PhysX 3D (scaffold)
```

## License

Each implementation carries its own license as specified in its project metadata.
