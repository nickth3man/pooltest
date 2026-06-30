# AGENTS.md — pool-cpp

2D pool/billiards game — C++17 + raylib 6.0, built with CMake.

## Stack

- **Language:** C++17 (`CMAKE_CXX_STANDARD 17`, `CXX_EXTENSIONS OFF`)
- **Graphics/audio/input:** raylib 6.0 (released 2026-04-23)
  - NOT pre-installed — fetched by CMake `FetchContent` from `https://github.com/raysan5/raylib/archive/refs/tags/6.0.tar.gz` (first configure needs network)
  - Built as a static library using its bundled GLFW + miniaudio
  - Backend: `PLATFORM_DESKTOP`, `GRAPHICS_API_OPENGL_33`
- **Build system:** CMake >=3.15 (tested 4.x)
- **Compiler:** MSVC 19.50 (VS 18 BuildTools) on Windows; Visual Studio generator auto-selected (multi-config: Debug/Release/MinSizeRel/RelWithDebInfo)
- **Target:** executable `pool` from `src/main.cpp`
- **Links:** `raylib`, `Threads::Threads`; on macOS also IOKit/Cocoa/OpenGL frameworks
- **raylib is a C library** — from C++, wrap `#include "raylib.h"` in `extern "C" { ... }`
- **`CMAKE_EXPORT_COMPILE_COMMANDS ON`** — `compile_commands.json` generated for clangd/IDEs

## Dev commands

Run from the `pool-cpp/` folder.

### Configure / build / run
- `cmake -S . -B build` — configure (first run downloads + builds raylib 6.0, ~1–2 min)
- `cmake --build build --config Release` — build Release (or `Debug`)
- `./build/Release/pool.exe` — run (Windows); `./build/pool` (single-config generators)
- Generator override: `cmake -S . -B build -G "Visual Studio 18 2026" -A x64`

### Reconfigure / clean
- Reconfigure: re-run `cmake -S . -B build`
- Clean build: `rm -rf build && cmake -S . -B build`
- `cat build/compile_commands.json` — inspection (generated for tooling)

### Dependencies
- No system raylib install — `FetchContent` manages it. To force a re-fetch, delete `build/_deps/raylib-src`.

### Prerequisites
- CMake on PATH (>=3.15)
- Visual Studio Build Tools with the "Desktop development with C++" workload (MSVC `cl.exe` + Windows SDK). On first configure CMake locates the compiler automatically via the VS generator.

Note: editor IntelliSense may flag `raylib.h` "not found" until the first CMake configure downloads it — expected.
