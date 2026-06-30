# Pool — Kotlin 3D

A 3D pool/billiards scaffold in Kotlin, built on [kool](https://github.com/kool-engine/kool) (a Vulkan/WebGPU/OpenGL-first game engine) and the [kool-physics](https://github.com/kool-engine/kool) module, which wraps NVIDIA PhysX 5.6.1 via [physx-jni](https://github.com/fabmax/physx-jni).

This is a **hello-world scaffold only**. It opens a 1280×800 window titled "Pool - Kotlin 3D" with a perspective orbit camera, a green ground plane, one white sphere, and a single directional light. There is no physics, no gameplay, and no balls yet — those land in the next phase.

## Stack

| Component | Coordinates | Version | Role |
|-----------|-------------|---------|------|
| kool-core | `de.fabmax.kool:kool-core` | 0.19.0 | Vulkan-first renderer with OpenGL fallback |
| kool-physics | `de.fabmax.kool:kool-physics` | 0.19.0 | Multiplatform physics facade over PhysX |
| physx-jni | `de.fabmax:physx-jni` | 2.7.2 | JNI bindings to NVIDIA PhysX 5.6.1 |
| Kotlin | `org.jetbrains.kotlin:kotlin-multiplatform` | 2.2.21 | Matches the official kool-basic-template |
| JVM toolchain | — | 17 | kool's minimum supported JDK |

## Prerequisites

- **JDK 17+** (21 recommended). Verify with `java -version`.
- **Gradle** is provided by the wrapper — no system install needed. You must generate the wrapper once (see below).
- A **Vulkan-capable GPU + driver**. kool auto-falls back to OpenGL (`useOpenGlFallback = true` in `KoolConfigJvm`), so an OpenGL 4.1+ context also works.
- **No manual native library install.** physx-jni ships natives as Maven classifier artifacts (`natives-windows`, `natives-linux`, `natives-macos`, `natives-macos-arm64`); `NativeLib.load()` extracts them to `%TEMP%\de.fabmax.physx-jni\2.7.2\` on Windows at runtime.

## Generate the Gradle wrapper (first time only)

From `pool-kotlin-3d/`, with a local Gradle 8.x available:

```bash
gradle wrapper --gradle-version 8.10.2
```

This creates `gradlew`, `gradlew.bat`, and `gradle/wrapper/gradle-wrapper.jar`. After that, all subsequent commands use the wrapper — your system Gradle is not needed again.

## Run

```bash
./gradlew run          # macOS, Linux, or Git Bash on Windows
gradlew.bat run        # Windows cmd or PowerShell
```

`./gradlew runJvm` is provided as a canonical alias matching the kool-basic-template convention.

## First run

The first run downloads roughly 250 MB of dependencies (kool + LWJGL + physx-jni natives) and may take a few minutes on a cold cache. When it finishes, a 1280×800 window opens showing a green ground plane and a single white sphere. Subsequent runs are much faster thanks to the Gradle dependency cache.

## Camera controls

- **Left-drag** — orbit
- **Right-drag** — pan
- **Wheel** — zoom
- **Esc** — close

## What's next

- **Table geometry** — 6 pockets, cushion walls, and a felt PBR material.
- **15 balls in a rack** — per-ball color and number textures.
- **Cue stick** — a cylinder mesh.
- **PhysX rigid bodies** — `RigidDynamic` + `PxSphereGeometry(0.0285f)` (a real 57 mm ball) driven from `onUpdate`; mesh transforms are synced back to PhysX each frame.
- **Cue strike impulse** — apply a linear impulse along the cue's forward axis on click.
- **Materials** — `PxMaterial` with high friction for cloth and near-elastic restitution for the phenolic-resin balls.
- **Game logic** — turns, fouls, shot clock.

## Windows notes

- Prefer `./gradlew run` from **Git Bash** or `gradlew.bat run` from **cmd / PowerShell**.
- physx-jni extracts natives to `%TEMP%\de.fabmax.physx-jni\2.7.2\`. If you see `UnsatisfiedLinkError: no PhysX in java.library.path`, delete that folder and rerun to force a clean re-extract.
- If Windows SmartScreen blocks `gradlew.bat`, choose **More info → Run anyway**.

## Project layout

```
pool-kotlin-3d/
├── settings.gradle.kts
├── gradle.properties
├── build.gradle.kts
├── gradle/
│   ├── libs.versions.toml
│   └── wrapper/
│       └── gradle-wrapper.properties
├── gradlew                       # generated
├── gradlew.bat                   # generated
├── src/
│   └── jvmMain/
│       └── kotlin/
│           └── Main.kt
└── README.md
```
