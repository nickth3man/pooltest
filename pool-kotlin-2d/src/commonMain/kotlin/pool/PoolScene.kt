package pool

import de.fabmax.kool.KoolContext
import de.fabmax.kool.input.Pointer
import de.fabmax.kool.input.PointerInput
import de.fabmax.kool.math.MutableVec2f
import de.fabmax.kool.math.MutableVec3f
import de.fabmax.kool.math.Vec2f
import de.fabmax.kool.math.Vec3f
import de.fabmax.kool.pipeline.ClearColorFill
import de.fabmax.kool.modules.ksl.KslUnlitShader
import de.fabmax.kool.physics2d.Body
import de.fabmax.kool.physics2d.BodyDef
import de.fabmax.kool.physics2d.BodyType
import de.fabmax.kool.physics2d.Geometry
import de.fabmax.kool.physics2d.Physics2dWorld
import de.fabmax.kool.physics2d.ShapeDef
import de.fabmax.kool.physics2d.SurfaceMaterial
import de.fabmax.kool.physics2d.createBody
import de.fabmax.kool.scene.OrthographicCamera
import de.fabmax.kool.scene.Scene
import de.fabmax.kool.scene.addColorMesh
import de.fabmax.kool.scene.addLineMesh
import de.fabmax.kool.scene.scene
import de.fabmax.kool.util.Color
import de.fabmax.kool.util.InterpolatableSimulation
import kotlin.math.min
import kotlin.math.sqrt

private object PoolTable {
    const val TABLE_W = 12f
    const val TABLE_H = 6f
    const val PLAY_W = 11.2f
    const val PLAY_H = 5.2f
    const val CUSHION = 0.4f
    const val BALL_R = 0.18f
    const val POCKET_R = 0.38f
    const val MAX_SHOT_SPEED = 18f
    const val REST_VELOCITY = 0.08f
    const val POCKET_CHECK_R = POCKET_R * POCKET_R

    val ballMaterial =
        ShapeDef(
            density = 1f,
            material = SurfaceMaterial(friction = 0.12f, restitution = 0.93f),
        )
    val cushionMaterial =
        ShapeDef(
            material = SurfaceMaterial(friction = 0.25f, restitution = 0.78f),
        )

    val pockets: List<Vec2f> =
        buildList {
            val hx = PLAY_W / 2f
            val hy = PLAY_H / 2f
            add(Vec2f(-hx, hy))
            add(Vec2f(0f, hy))
            add(Vec2f(hx, hy))
            add(Vec2f(-hx, -hy))
            add(Vec2f(0f, -hy))
            add(Vec2f(hx, -hy))
        }

    val ballColors =
        listOf(
            Color(0.95f, 0.85f, 0.1f), // 1 yellow
            Color(0.1f, 0.35f, 0.85f), // 2 blue
            Color(0.85f, 0.15f, 0.15f), // 3 red
            Color(0.55f, 0.15f, 0.65f), // 4 purple
            Color(0.95f, 0.45f, 0.05f), // 5 orange
            Color(0.1f, 0.55f, 0.2f), // 6 green
            Color(0.55f, 0.1f, 0.15f), // 7 maroon
            Color(0.08f, 0.08f, 0.08f), // 8 black
            Color(0.95f, 0.85f, 0.1f), // 9
            Color(0.1f, 0.35f, 0.85f), // 10
            Color(0.85f, 0.15f, 0.15f), // 11
            Color(0.55f, 0.15f, 0.65f), // 12
            Color(0.95f, 0.45f, 0.05f), // 13
            Color(0.1f, 0.55f, 0.2f), // 14
            Color(0.95f, 0.95f, 0.95f), // 15 cue ball replacement if needed
        )
}

private data class PoolBall(
    val body: Body,
    val color: Color,
    val mesh: de.fabmax.kool.scene.ColorMesh,
    val isCueBall: Boolean = false,
)

private fun rackPositions(): List<Vec2f> {
    val positions = mutableListOf<Vec2f>()
    val rackX = PoolTable.PLAY_W / 2f - 2.8f * PoolTable.BALL_R
    val rowSpacing = PoolTable.BALL_R * sqrt(3f)
    for (row in 0 until 5) {
        for (col in 0..row) {
            val x = rackX + row * rowSpacing
            val y = (col - row / 2f) * 2f * PoolTable.BALL_R
            positions += Vec2f(x, y)
        }
    }
    return positions
}

private fun Physics2dWorld.addCushion(
    pos: Vec2f,
    halfW: Float,
    halfH: Float,
) {
    val body = createBody(BodyType.Static, pos)
    body.attachShape(Geometry.Box(halfW, halfH), PoolTable.cushionMaterial)
}

private fun Physics2dWorld.buildTable() {
    val hw = PoolTable.PLAY_W / 2f
    val hh = PoolTable.PLAY_H / 2f
    val c = PoolTable.CUSHION / 2f
    val gap = PoolTable.POCKET_R * 1.6f

    // Top and bottom rails split around center and corner pockets.
    addCushion(Vec2f(-(hw + c) / 2f - gap / 4f, hh + c), (hw - gap) / 2f, PoolTable.CUSHION / 2f)
    addCushion(Vec2f((hw + c) / 2f + gap / 4f, hh + c), (hw - gap) / 2f, PoolTable.CUSHION / 2f)
    addCushion(Vec2f(-(hw + c) / 2f - gap / 4f, -hh - c), (hw - gap) / 2f, PoolTable.CUSHION / 2f)
    addCushion(Vec2f((hw + c) / 2f + gap / 4f, -hh - c), (hw - gap) / 2f, PoolTable.CUSHION / 2f)

    // Side rails with corner pocket gaps.
    addCushion(Vec2f(-hw - c, 0f), PoolTable.CUSHION / 2f, hh - gap)
    addCushion(Vec2f(hw + c, 0f), PoolTable.CUSHION / 2f, hh - gap)
}

private fun Physics2dWorld.createBall(
    scene: Scene,
    position: Vec2f,
    color: Color,
    isCueBall: Boolean = false,
): PoolBall {
    val body =
        createBody(
            BodyDef(
                type = BodyType.Dynamic,
                position = position,
                linearDamping = 1.6f,
                angularDamping = 1.8f,
                fixedRotation = false,
                isBullet = true,
            ),
        )
    body.attachShape(Geometry.Circle(PoolTable.BALL_R), PoolTable.ballMaterial)

    val mesh =
        scene.addColorMesh {
            generate {
                circle {
                    radius = PoolTable.BALL_R
                    steps = 24
                }
            }
            shader =
                KslUnlitShader {
                    color { constColor(color) }
                }
            transform.setIdentity().translate(position.x, position.y, PoolTable.BALL_R)
        }

    return PoolBall(body, color, mesh, isCueBall)
}

private fun Scene.pointerToWorld2d(pointer: Pointer): Vec2f? {
    val world = MutableVec3f()
    return if (camera.unProjectScreen(Vec3f(pointer.pos.x, pointer.pos.y, 0f), mainRenderPass.viewport, world)) {
        Vec2f(world.x, world.y)
    } else {
        null
    }
}

private fun Vec2f.distanceSqTo(other: Vec2f): Float {
    val dx = x - other.x
    val dy = y - other.y
    return dx * dx + dy * dy
}

private fun PoolBall.isNearlyStopped(): Boolean = body.isValid && body.linearVelocity.length() < PoolTable.REST_VELOCITY

fun buildPoolScene(ctx: KoolContext) =
    scene("Pool") {
        clearColor = ClearColorFill(Color(0.2f, 0.2f, 0.2f, 1f))
        
        camera =
            OrthographicCamera().apply {
                setCentered(height = PoolTable.TABLE_H + 1f, near = -10f, far = 10f)
            }

        val world = Physics2dWorld(gravity = Vec2f(0f, 0f))
        world.registerHandlers(this)

        val balls = mutableListOf<PoolBall>()
        var cueBall: PoolBall? = null

        world.buildTable()

        val cuePos = Vec2f(-PoolTable.PLAY_W / 2f + 2.5f * PoolTable.BALL_R, 0f)
        cueBall = world.createBall(this, cuePos, Color.WHITE, isCueBall = true).also { balls += it }

        val rack = rackPositions()
        rack.forEachIndexed { i, pos ->
            balls += world.createBall(this, pos, PoolTable.ballColors[i])
        }

        // Table felt
        addColorMesh {
            generate {
                rect { size.set(PoolTable.PLAY_W, PoolTable.PLAY_H) }
            }
            shader =
                KslUnlitShader {
                    color { constColor(Color(0.08f, 0.42f, 0.16f)) }
                }
        }

        // Outer table frame
        addColorMesh {
            transform.translate(0f, 0f, -0.005f)
            generate {
                rect { size.set(PoolTable.TABLE_W, PoolTable.TABLE_H) }
            }
            shader =
                KslUnlitShader {
                    color { constColor(Color(0.35f, 0.18f, 0.08f)) }
                }
        }

        // Cushion visuals
        addColorMesh {
            transform.translate(0f, 0f, 0.002f)
            generate {
                rect { size.set(PoolTable.PLAY_W + PoolTable.CUSHION * 2f, PoolTable.PLAY_H + PoolTable.CUSHION * 2f) }
            }
            shader =
                KslUnlitShader {
                    color { constColor(Color(0.06f, 0.28f, 0.11f)) }
                }
        }

        // Pockets
        PoolTable.pockets.forEach { pocket ->
            addColorMesh {
                transform.translate(pocket.x, pocket.y, -0.003f)
                generate {
                    circle {
                        radius = PoolTable.POCKET_R
                        steps = 32
                    }
                }
                shader =
                    KslUnlitShader {
                        color { constColor(Color(0.02f, 0.02f, 0.02f)) }
                    }
            }
        }

        // Aim / cue line
        val aimLine =
            addLineMesh {
                shader =
                    KslUnlitShader {
                        color { vertexColor() }
                    }
            }

        var aiming = false
        var aimPointerPos = MutableVec2f()

        fun allBallsStopped(): Boolean = balls.all { it.isNearlyStopped() }

        fun pocketBall(ball: PoolBall) {
            if (!ball.body.isValid) return
            world.removeBody(ball.body)
            ball.mesh.isVisible = false
            balls.remove(ball)
            if (ball.isCueBall) {
                cueBall = world.createBall(this, cuePos, Color.WHITE, isCueBall = true).also { balls += it }
            }
        }

        fun checkPockets() {
            val toPocket =
                balls.filter { ball ->
                    ball.body.isValid &&
                        PoolTable.pockets.any { pocket ->
                            ball.body.position.distanceSqTo(pocket) < PoolTable.POCKET_CHECK_R
                        }
                }
            toPocket.forEach { pocketBall(it) }
        }

        fun updateBallMeshes() {
            balls.forEach { ball ->
                if (!ball.body.isValid) {
                    ball.mesh.isVisible = false
                    return@forEach
                }
                val pos = ball.body.position
                ball.mesh.isVisible = true
                ball.mesh.transform.setIdentity().translate(pos.x, pos.y, PoolTable.BALL_R)
            }
        }

        fun updateAimLine() {
            aimLine.clear()
            val cue = cueBall ?: return
            if (!aiming || !cue.body.isValid) return

            val cuePos2d = cue.body.position
            val dx = cuePos2d.x - aimPointerPos.x
            val dy = cuePos2d.y - aimPointerPos.y
            val dist = sqrt(dx * dx + dy * dy)
            if (dist < 0.05f) return

            val power = min(dist * 4f, PoolTable.MAX_SHOT_SPEED) / PoolTable.MAX_SHOT_SPEED
            val lineLen = 0.5f + power * 2.5f
            val nx = dx / dist
            val ny = dy / dist

            val start = Vec3f(cuePos2d.x, cuePos2d.y, PoolTable.BALL_R + 0.02f)
            val end =
                Vec3f(
                    cuePos2d.x + nx * lineLen,
                    cuePos2d.y + ny * lineLen,
                    PoolTable.BALL_R + 0.02f,
                )
            aimLine.addLine(start, end, Color(1f, 0.95f, 0.7f, 0.85f))
        }

        fun handleInput() {
            val ptr = PointerInput.primaryPointer
            if (!ptr.isValid || ptr.isConsumed()) return

            val cue = cueBall ?: return
            if (!cue.body.isValid) return

            val worldPos = pointerToWorld2d(ptr) ?: return

            when {
                ptr.isLeftButtonPressed && allBallsStopped() -> {
                    val dist = cue.body.position.distanceSqTo(worldPos)
                    if (dist < (PoolTable.BALL_R * 3f) * (PoolTable.BALL_R * 3f)) {
                        aiming = true
                        aimPointerPos.set(worldPos.x, worldPos.y)
                    }
                }

                aiming && ptr.isLeftButtonDown -> {
                    aimPointerPos.set(worldPos.x, worldPos.y)
                }

                aiming && ptr.isLeftButtonReleased -> {
                    val cuePos2d = cue.body.position
                    val dx = cuePos2d.x - aimPointerPos.x
                    val dy = cuePos2d.y - aimPointerPos.y
                    val dist = sqrt(dx * dx + dy * dy)
                    if (dist > 0.08f) {
                        val speed = min(dist * 4f, PoolTable.MAX_SHOT_SPEED)
                        cue.body.setLinearVelocity(Vec2f(dx / dist * speed, dy / dist * speed))
                    }
                    aiming = false
                }
            }
        }

        world.simulationListeners +=
            object : InterpolatableSimulation {
                override fun simulateStep(timeStep: Float) {
                    checkPockets()
                }
            }

        onUpdate {
            handleInput()
            updateAimLine()
            updateBallMeshes()
        }
    }
