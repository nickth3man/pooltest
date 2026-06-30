import org.gradle.internal.os.OperatingSystem
import org.jetbrains.kotlin.gradle.ExperimentalKotlinGradlePluginApi

plugins {
    kotlin("multiplatform") version "2.2.21"
}

group = "pool.kotlin3d"
version = "0.1.0"

kotlin {
    jvm {
        @OptIn(ExperimentalKotlinGradlePluginApi::class)
        binaries {
            executable {
                mainClass.set("MainKt")
                applicationDefaultJvmArgs = buildList {
                    add("--add-opens=java.base/java.lang=ALL-UNNAMED")
                    add("--enable-native-access=ALL-UNNAMED")
                    if (OperatingSystem.current().isMacOsX) {
                        add("-XstartOnFirstThread")
                    }
                }
            }
        }
    }
    jvmToolchain(21)

    sourceSets {
        commonMain.dependencies {
            implementation(libs.kool.core)
            // kool-physics is on the classpath so the next phase can add PhysX rigid bodies.
            implementation(libs.kool.physics)
        }
        val jvmMain by getting {
            dependencies {
                implementation(libs.physx.jni)
                // physx-jni platform natives — auto-extracted from the jar at runtime by NativeLib.load()
                runtimeOnly("de.fabmax:physx-jni:2.7.2:natives-windows")
                runtimeOnly("de.fabmax:physx-jni:2.7.2:natives-linux")
                runtimeOnly("de.fabmax:physx-jni:2.7.2:natives-macos")
                runtimeOnly("de.fabmax:physx-jni:2.7.2:natives-macos-arm64")
            }
        }
    }
}

tasks.register("run") {
    group = "application"
    dependsOn("runJvm")
}
