// pool-cpp — 2D 8-ball pool with raylib.

extern "C" {
#include "raylib.h"
}

#include "game.h"
#include "render.h"

int main(void) {
    const int screenWidth = 800;
    const int screenHeight = 600;

    InitWindow(screenWidth, screenHeight, "Pool - C++");
    if (!IsWindowReady()) return 1;

    SetTargetFPS(60);

    Game game;
    Renderer renderer;

    game.reset();
    renderer.init();

    while (!WindowShouldClose()) {
        const float dt = GetFrameTime();
        game.update(dt);

        BeginDrawing();
            ClearBackground({20, 20, 28, 255});
            renderer.draw(game);
        EndDrawing();
    }

    renderer.shutdown();
    CloseWindow();
    return 0;
}
