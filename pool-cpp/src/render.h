#pragma once

extern "C" {
#include "raylib.h"
}

#include <vector>

#include "ball.h"
#include "game.h"
#include "table.h"

struct BallSprites {
    Texture2D textures[kBallCount] = {};
    int size = 0;

    void create(int diameter);
    void destroy();
    void draw(const Ball& ball) const;
};

struct Renderer {
    BallSprites sprites;

    void init();
    void shutdown();
    void draw(const Game& game) const;
};
