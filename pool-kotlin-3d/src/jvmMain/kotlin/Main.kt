import de.fabmax.kool.KoolApplication
import de.fabmax.kool.KoolConfigJvm
import de.fabmax.kool.math.Vec2i
import de.fabmax.kool.math.Vec3f
import de.fabmax.kool.math.deg
import de.fabmax.kool.modules.ksl.KslPbrShader
import de.fabmax.kool.scene.Scene
import de.fabmax.kool.scene.addColorMesh
import de.fabmax.kool.scene.defaultOrbitCamera
import de.fabmax.kool.scene.scene
import de.fabmax.kool.util.Color

/**
 * Pool — Kotlin 3D (scaffold).
 *
 * Renders a perspective orbit camera, a green ground plane, one white sphere (cue ball),
 * and a single directional light. Physics (PhysX rigid bodies via kool-physics) lands next.
 *
 * Run: ./gradlew run   (or ./gradlew runJvm)
 */
fun main() = KoolApplication(
    config = KoolConfigJvm(
        windowTitle = "Pool - Kotlin 3D",
        windowSize = Vec2i(1280, 800),
        useOpenGlFallback = true,
    )
) {
    ctx.scenes += scene("main") {
        setupMainScene()
    }
}

private fun Scene.setupMainScene() {
    // Perspective orbit camera (LMB orbit, RMB pan, wheel zoom).
    defaultOrbitCamera()

    lighting.singleDirectionalLight {
        setup(Vec3f(-0.4f, -1f, -0.6f))
        setColor(Color.WHITE, 5f)
    }

    // Green ground plane (5m x 3m), rotated from XY to XZ (Y-up).
    addColorMesh("ground") {
        generate {
            rect {
                size.set(5f, 3f)
            }
        }
        transform.rotate((-90f).deg, Vec3f.X_AXIS)
        shader = KslPbrShader {
            color { constColor(Color(0.10f, 0.45f, 0.18f)) }
            metallic(0f)
            roughness(0.9f)
        }
    }

    // White cue ball (radius 0.1m) sitting on the ground.
    addColorMesh("cue-ball") {
        generate {
            uvSphere {
                radius = 0.1f
                steps = 24
            }
        }
        transform.translate(0f, 0.1f, 0f)
        shader = KslPbrShader {
            color { constColor(Color.WHITE) }
            metallic(0f)
            roughness(0.25f)
        }
    }
}
