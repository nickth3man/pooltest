#pragma once

#include <vector>

#include "ball.h"
#include "table.h"
#include "vec2.h"

namespace Physics {

bool allBallsStopped(const std::vector<Ball>& balls);
void step(std::vector<Ball>& balls, const Table& table, float dt);

struct RayHit {
    bool hit = false;
    float distance = 0.0f;
    int ballId = -1;
    Vec2 point;
    Vec2 normal;
};

RayHit raycastFirstBall(const Vec2& origin,
                        const Vec2& dir,
                        float maxDistance,
                        const std::vector<Ball>& balls,
                        float cueRadius,
                        int ignoreId = kCueBallId);

}  // namespace Physics
