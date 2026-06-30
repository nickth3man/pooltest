package pool

import de.fabmax.kool.KoolApplication
import de.fabmax.kool.KoolConfigJvm

/** Desktop (JVM) entry point. */
fun main() =
    KoolApplication(
        config =
            KoolConfigJvm(
                windowTitle = "Pool - Kotlin 2D",
                useOpenGlFallback = true,
            ),
    ) {
        ctx.scenes += buildPoolScene(ctx)
    }
