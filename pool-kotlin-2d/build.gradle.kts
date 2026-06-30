import org.gradle.internal.os.OperatingSystem
import org.jetbrains.kotlin.gradle.ExperimentalKotlinGradlePluginApi
import org.jetbrains.kotlin.gradle.targets.js.dsl.ExperimentalDistributionDsl
import org.jetbrains.kotlin.gradle.targets.js.webpack.KotlinWebpackConfig

plugins {
    kotlin("multiplatform") version "2.2.21"
    id("com.diffplug.spotless") version "6.25.0"
}

kotlin {
    jvm {
        @OptIn(ExperimentalKotlinGradlePluginApi::class)
        binaries {
            executable {
                mainClass.set("pool.LinuxLauncherKt")
                applicationDefaultJvmArgs =
                    buildList {
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

    js {
        binaries.executable()
        browser {
            @OptIn(ExperimentalDistributionDsl::class)
            distribution {
                outputDirectory.set(File("$rootDir/dist/js"))
            }
            commonWebpackConfig {
                mode = KotlinWebpackConfig.Mode.DEVELOPMENT
                cssSupport { enabled.set(true) }
                // Enable WebAssembly support for kool-physics-2d
                configDirectory = File("$rootDir/webpack.config.d")
            }
        }
        @OptIn(ExperimentalKotlinGradlePluginApi::class)
        compilerOptions {
            target.set("es2015")
        }
    }

    sourceSets {
        val commonMain by getting {
            dependencies {
                implementation(libs.kool.core)
                implementation(libs.kool.physics2d)
            }
        }
        val jvmMain by getting
        val jsMain by getting
    }
}

val clean by tasks.getting(Task::class) {
    doLast {
        delete("$rootDir/dist")
    }
}

spotless {
    kotlin {
        target("src/**/*.kt")
        ktlint()
    }
    kotlinGradle {
        target("*.gradle.kts")
        ktlint()
    }
}
