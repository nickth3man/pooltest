# AGENTS.md — pool-kotlin-2d

2D pool/billiards — kool engine (Kotlin Multiplatform: JVM desktop + JS browser) + kool-physics-2d (Box2D 3.1.1).

## Stack

- **Language:** Kotlin 2.2.21 (Kotlin Multiplatform plugin `org.jetbrains.kotlin.multiplatform`)
- **Targets:** `jvm` (desktop) + `js` (browser, WebGL/WebGPU)
- **Dependencies (Maven Central):**
  - `de.fabmax.kool:kool-core:0.19.0` — Vulkan/WebGPU/OpenGL engine
  - `de.fabmax.kool:kool-physics-2d:0.19.0` — Box2D 3.1.1 JNI/WASM bindings
- **Version catalog:** `gradle/libs.versions.toml` — aliases `kool-core`, `kool-physics2d` (a catalog alias segment must NOT start with a digit — that's why it's `kool-physics2d`, not `kool-physics-2d`; the artifact coordinate keeps its hyphen)
- **Build tool:** Gradle 8.10.2 (wrapper); `org.gradle.configuration-cache=false`
- **JDK toolchain:** 21 (Eclipse Temurin 21 tested)
- **JVM executable binary:** `jvm { binaries { executable { mainClass = "pool.LinuxLauncherKt"; applicationDefaultJvmArgs = ["--add-opens=java.base/java.lang=ALL-UNNAMED", "--enable-native-access=ALL-UNNAMED"] } } }` — auto-creates the `runJvm` task
- **KMP constraints:** do NOT apply the Gradle `application` or `java` plugin (incompatible with KMP). `runJvm` is auto-created — do NOT `tasks.register("runJvm")`.
- **Scene attach:** inside `KoolApplication { ... }`, use `ctx.scenes += scene("name") { ... }`
- **Layout:** `src/commonMain/kotlin/pool/PoolScene.kt`, `src/jvmMain/kotlin/pool/LinuxLauncher.kt`, `src/jsMain/kotlin/pool/JsLauncher.kt` + `src/jsMain/resources/index.html`

## Dev commands

Run from the `pool-kotlin-2d/` folder. Requires `java` (JDK 21) on PATH.

### Wrapper (first time only, if `gradlew` is missing)
- `gradle wrapper --gradle-version 8.10.2`

### Run
- `./gradlew runJvm` — desktop run (Windows: `.\gradlew.bat runJvm`)
- `./gradlew jsBrowserDevelopmentRun` — browser dev server (~http://localhost:8080; Chrome for WebGPU)

### Compile / build
- `./gradlew compileKotlinJvm` — compile JVM target (no GUI launch)
- `./gradlew compileKotlinJs` / `compileDevelopmentExecutableKotlinJs` — compile JS
- `./gradlew build` — assemble + check all targets
- `./gradlew jsBrowserProductionWebpack` — production web bundle → `dist/js/`

### Quality / deps / cleanup
- `./gradlew spotlessCheck` — check code formatting/linting (ktlint)
- `./gradlew spotlessApply` — auto-format and lint-fix code
- Note: If `.editorconfig` is modified, run with `--no-daemon` (e.g., `./gradlew spotlessApply --no-daemon`) to bypass Gradle daemon caching.
- `./gradlew test` (all targets); `jvmTest`; `jsTest`
- `./gradlew dependencies` — print dependency graph; `--refresh-dependencies` — force re-resolve
- `./gradlew clean` — remove `build/`
- Stale config: prefix any task with `--no-configuration-cache`

### Environment
- `export JAVA_HOME=/path/to/jdk-21` (or Windows equivalent); verify with `java -version`
