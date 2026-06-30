# AGENTS.md — pool-kotlin-3d

3D pool/billiards — kool engine + kool-physics (NVIDIA PhysX 5.6.1 via physx-jni), Kotlin Multiplatform (JVM desktop).

## Stack

- **Language:** Kotlin 2.2.21 (Kotlin Multiplatform plugin `org.jetbrains.kotlin.multiplatform`)
- **Target:** `jvm` only (the 3D/PhysX path is JVM-native)
- **Dependencies (Maven Central):**
  - `de.fabmax.kool:kool-core:0.19.0` — Vulkan/WebGPU/OpenGL engine
  - `de.fabmax.kool:kool-physics:0.19.0` — multiplatform physics facade
  - `de.fabmax:physx-jni:2.7.2` — NVIDIA PhysX 5.6.1 Java bindings
  - Native libraries (declared as `runtimeOnly` classifier artifacts): `de.fabmax:physx-jni:2.7.2:natives-windows`, `:natives-linux`, `:natives-macos`, `:natives-macos-arm64` — auto-extracted from the jar to `%TEMP%\de.fabmax.physx-jni\2.7.2\` on Windows at runtime (no manual native install)
- **Version catalog:** `gradle/libs.versions.toml` — aliases `kool-core`, `kool-physics`, `physx-jni`
- **Build tool:** Gradle 8.10.2 (wrapper)
- **JDK toolchain:** 17+ (21 recommended; Eclipse Temurin 21 tested)
- **JVM executable binary:** `jvm { binaries { executable { mainClass = "MainKt"; applicationDefaultJvmArgs = ["--add-opens=java.base/java.lang=ALL-UNNAMED", "--enable-native-access=ALL-UNNAMED"] } } }` — auto-creates the `runJvm` task; a `run` alias (`dependsOn("runJvm")`) is also registered
- **KMP constraints:** do NOT apply the Gradle `application` or `java` plugin (incompatible with KMP). `runJvm` is auto-created — do NOT `tasks.register("runJvm")` (duplicate-task error). Registering a `run` alias is fine.
- **Scene attach:** inside `KoolApplication { ... }`, use `ctx.scenes += scene("name") { ... }`
- **Render:** Vulkan default + OpenGL fallback (`KoolConfigJvm.useOpenGlFallback = true`); GPU needs Vulkan or OpenGL 4.5
- **First run:** downloads kool + LWJGL + physx-jni natives (~250 MB)
- **Layout:** `src/jvmMain/kotlin/Main.kt`

## Dev commands

Run from the `pool-kotlin-3d/` folder. Requires `java` (JDK 17+/21) on PATH.

### Wrapper (first time only, if `gradlew` is missing)
- `gradle wrapper --gradle-version 8.10.2`

### Run
- `./gradlew run` — desktop run via alias (Windows: `.\gradlew.bat run`)
- `./gradlew runJvm` — equivalent, the auto-created task

### Compile / build
- `./gradlew compileKotlinJvm` — compile JVM target (no GUI launch)
- `./gradlew build` — assemble + check

### Quality / deps / cleanup
- `./gradlew test`
- `./gradlew dependencies` — dependency graph; `--refresh-dependencies` — force re-resolve
- `./gradlew clean` — remove `build/`

### Environment / troubleshooting
- `export JAVA_HOME=/path/to/jdk-21`; verify `java -version`
- On `UnsatisfiedLinkError: no PhysX in java.library.path`, delete `%TEMP%\de.fabmax.physx-jni\2.7.2\` and rerun (forces native re-extraction)
