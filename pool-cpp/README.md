# pool-cpp

2D 8-ball pool written in C++ with [raylib 6.0](https://www.raylib.com/) and CMake.

## Features

- **Table** — green felt, wooden rails, 6 pockets, cushion segments with ball reflection
- **16 balls** — cue ball + 15 object balls in a standard 8-ball rack (solids, stripes, 8-ball)
- **Physics** — equal-mass elastic ball-ball collisions, cushion bounces, rolling friction
- **Cue stick** — aim with the mouse, hold **LMB** and drag to set power, release to shoot
- **Sprites** — procedural ball textures with numbers, stripes, and gloss

## Controls

| Input | Action |
|-------|--------|
| Mouse move | Aim from cue ball |
| Hold LMB + drag | Set shot power |
| Release LMB | Shoot |
| R | Reset rack |

## Stack

- **C++17**
- **raylib 6.0** — fetched automatically by CMake `FetchContent`
- **CMake >= 3.15**

## Prerequisites

| Tool | Notes |
|------|-------|
| CMake >= 3.15 | During install, tick **"Add CMake to PATH"**. |
| MSVC: Visual Studio Build Tools 2019+ | Install the **"Desktop development with C++"** workload. |

## Configure & build

```bash
cmake -B build
cmake --build build --config Release
```

## Run

```bash
.\build\Release\pool.exe
```

On the first configure, CMake downloads the raylib 6.0 source tarball and builds it as part of the project (~1–2 minutes). Subsequent builds are incremental.

## Project layout

```
src/
  main.cpp      — entry point, game loop
  game.cpp      — input, rack reset, shot logic
  physics.cpp   — collisions, friction, pocketing
  table.cpp     — cushions and pocket geometry
  ball.cpp      — ball setup and 8-ball rack
  render.cpp    — table drawing, cue stick, ball sprites
  vec2.h        — 2D vector math
```

## Windows notes

- The Visual Studio generator is used by default when VS is installed, so the binary lands at `build/Release/pool.exe`.
- raylib builds as a static library on MSVC — no extra system libraries required.
