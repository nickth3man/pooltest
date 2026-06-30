# Pool - Kotlin 2D

Real-time 2D pool/billiards scaffold built with [kool](https://github.com/kool-engine/kool) (Vulkan/WebGPU/OpenGL engine) and the bundled `kool-physics-2d` module (Box2D 3.1.1 bindings). This is a hello-world scaffold: a single kool scene, an orthographic top-down camera, a green felt table, and one white ball. It launches on desktop (JVM) and in the browser (JS). Physics integration lands in the next iteration.

## Stack

| Component            | Version / Detail                                                                |
|----------------------|---------------------------------------------------------------------------------|
| Language             | Kotlin 2.2.21 (Kotlin Multiplatform)                                            |
| Engine               | kool 0.19.0 (`de.fabmax.kool:kool-core`)                                        |
| Physics              | kool-physics-2d 0.19.0 (`de.fabmax.kool:kool-physics-2d`, Box2D 3.1.1)         |
| Targets              | JVM desktop (Windows / macOS / Linux) + browser (WebGL / WebGPU)                |
| Build tool           | Gradle 8.10.2                                                                   |

## Prerequisites

- **JDK 21+** — verify with `java -version`. On Windows, install Adoptium Temurin 21 or run `winget install Microsoft.OpenJDK.21`.
- **Gradle is not required system-wide** — the wrapper is provided, but the wrapper script (`gradlew`, `gradlew.bat`) and the wrapper JAR must be generated once (see below).
- A **Vulkan-** or **OpenGL 4.5-capable GPU** for desktop runs.

## Generate the Gradle wrapper (first time only)

From `pool-kotlin-2d/`, if you have a local Gradle 8.x installed:

```sh
gradle wrapper --gradle-version 8.10.2
```

This creates `gradlew`, `gradlew.bat`, and `gradle/wrapper/gradle-wrapper.jar`. After that, all subsequent commands use the wrapper and no system Gradle is needed.

## Run desktop

macOS / Linux:

```sh
./gradlew runJvm
```

Windows:

```sh
gradlew.bat runJvm
```

A window titled **"Pool - Kotlin 2D"** opens, showing a green felt table and a single white ball.

## Run browser

```sh
./gradlew jsBrowserDevelopmentRun
```

Then open the printed URL (default <http://localhost:8080>) in Chrome for WebGPU. Firefox / Safari fall back to WebGL2.

## Production web build

```sh
./gradlew jsBrowserProductionWebpack
```

Output is written to `dist/js/`.

## What's next

- Wire up `kool-physics-2d`: `Physics2dWorld()`, `world.registerHandlers(scene)`, static `Geometry.Box` walls + dynamic `Geometry.Circle` ball.
- Cue stick aim & power via pointer input.
- Full 15-ball rack.
- Score / fouls / turn UI.
- For a verified `Physics2dWorld` usage pattern, see the [kool MixerDemo](https://github.com/kool-engine/kool/blob/main/kool-demo/src/commonMain/kotlin/de/fabmax/kool/demo/physics/box2d/mixer/MixerDemo.kt).

## Troubleshooting

- **JNI natives not unpacked on first run** — run `./gradlew jvmJar` once, then re-run `runJvm`.
- **Black window on Windows Intel iGPU** — Vulkan init failed. `useOpenGlFallback = true` in `KoolConfigJvm` auto-falls back, or set `renderBackend = RenderBackendGl` explicitly.
- **`JAVA_HOME` errors** — set `JAVA_HOME` to your JDK 21 path (e.g. `setx JAVA_HOME "C:\Program Files\Eclipse Adoptium\jdk-21"`), reopen the terminal, and retry.
