package pool

import de.fabmax.kool.KoolApplication
import de.fabmax.kool.KoolConfigJs
import de.fabmax.kool.NativeAssetLoader
import de.fabmax.kool.physics2d.Physics2d
import de.fabmax.kool.pipeline.backend.gl.RenderBackendGl

/** Browser (JS) entry point. */
fun main() =
    KoolApplication(
        config =
            KoolConfigJs(
                canvasName = "glCanvas",
                defaultAssetLoader = NativeAssetLoader("."),
                renderBackend = RenderBackendGl,
                isGlobalKeyEventGrabbing = true,
                deviceScaleLimit = 1.5,
                useWebGlFallback = true,
            ),
    ) {
        Physics2d.loadPhysics2d()
        ctx.scenes += buildPoolScene(ctx)
    }
