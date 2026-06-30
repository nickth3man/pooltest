#pragma once

extern "C" {
#include "raylib.h"
}

#include <vector>

#include "vec2.h"

struct Segment {
    Vec2 a;
    Vec2 b;
};

struct Pocket {
    Vec2 center;
    float radius = 16.0f;
};

struct Table {
    Rectangle playArea = {50.0f, 80.0f, 700.0f, 350.0f};
    float ballRadius = 10.0f;
    std::vector<Segment> cushions;
    std::vector<Pocket> pockets;

    void build();
};
