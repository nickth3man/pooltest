# Pool Games Collection

A multi-language billiards/pool game implemented in **7 different tech stacks**. Each implementation now lives in its own repository so every game has an isolated dev setup, build system, and release cycle.

| Implementation | Language | Framework / Engine | Platform | Status | Repo |
|---------------|----------|-------------------|----------|--------|------|
| pool-go | Go | Ebitengine v2.9 | Desktop + WASM | **Full game** | [`YOUR_USERNAME/pool-go`](https://github.com/YOUR_USERNAME/pool-go) |
| pool-cpp | C++17 | raylib 6.0 | Desktop | **Feature-complete** | [`YOUR_USERNAME/pool-cpp`](https://github.com/YOUR_USERNAME/pool-cpp) |
| pool-ts | TypeScript | Phaser 4 + Vite 8 | Browser | **Scaffold** | [`YOUR_USERNAME/pool-ts`](https://github.com/YOUR_USERNAME/pool-ts) |
| pool-python | Python 3.11+ | pygame-ce | Desktop | **Scaffold** | [`YOUR_USERNAME/pool-python`](https://github.com/YOUR_USERNAME/pool-python) |
| pool-rust | Rust | Bevy 0.18 + Rapier2D | Desktop | **Scaffold** | [`YOUR_USERNAME/pool-rust`](https://github.com/YOUR_USERNAME/pool-rust) |
| pool-kotlin-2d | Kotlin | kool + Box2D (kool-physics-2d) | Desktop + Browser | **Scaffold** | [`YOUR_USERNAME/pool-kotlin-2d`](https://github.com/YOUR_USERNAME/pool-kotlin-2d) |
| pool-kotlin-3d | Kotlin | kool + NVIDIA PhysX (kool-physics) | Desktop (JVM) | **Scaffold** | [`YOUR_USERNAME/pool-kotlin-3d`](https://github.com/YOUR_USERNAME/pool-kotlin-3d) |

Replace `YOUR_USERNAME` in the links above with your actual GitHub username or organization.

## Local copies

Each game has been moved to a sibling directory next to this meta-repo:

```
C:\Users\nicolas\Documents\GitHub\
├── pooltest/           # this repo — index + links
├── pool-go/
├── pool-cpp/
├── pool-ts/
├── pool-python/
├── pool-rust/
├── pool-kotlin-2d/
└── pool-kotlin-3d/
```

## Running

Each repo contains its own `README.md` and `AGENTS.md` with full prerequisites and dev commands. The short version:

```bash
# pool-go
cd ../pool-go && go run .

# pool-cpp
cd ../pool-cpp && cmake -B build && cmake --build build --config Release && .\build\Release\pool.exe

# pool-ts
cd ../pool-ts && npm install && npm run dev

# pool-python
cd ../pool-python && uv sync && uv run pool-python

# pool-rust
cd ../pool-rust && cargo run

# pool-kotlin-2d
cd ../pool-kotlin-2d && ./gradlew runJvm

# pool-kotlin-3d
cd ../pool-kotlin-3d && ./gradlew run
```

## License

Each implementation carries its own license as specified in its project metadata.
