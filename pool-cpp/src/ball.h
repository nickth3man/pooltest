#pragma once

extern "C" {
#include "raylib.h"
}

#include <vector>

#include "vec2.h"

struct Ball {
    int id = 0;
    Vec2 pos;
    Vec2 vel;
    float radius = 10.0f;
    Color color = WHITE;
    bool active = true;
    bool isStripe = false;
    bool isCue = false;
};

constexpr int kBallCount = 16;
constexpr int kCueBallId = 0;

inline Color ballColor(int id) {
    switch (id) {
        case 0:  return {245, 245, 245, 255};
        case 1:  return {255, 215, 0, 255};
        case 2:  return {0, 80, 200, 255};
        case 3:  return {200, 30, 30, 255};
        case 4:  return {120, 30, 160, 255};
        case 5:  return {255, 120, 0, 255};
        case 6:  return {0, 140, 60, 255};
        case 7:  return {140, 20, 40, 255};
        case 8:  return {20, 20, 20, 255};
        case 9:  return {255, 215, 0, 255};
        case 10: return {0, 80, 200, 255};
        case 11: return {200, 30, 30, 255};
        case 12: return {120, 30, 160, 255};
        case 13: return {255, 120, 0, 255};
        case 14: return {0, 140, 60, 255};
        case 15: return {140, 20, 40, 255};
        default: return WHITE;
    }
}

void rackBalls(std::vector<Ball>& balls, float rackX, float rackY, float spacing);
