#include "ball.h"

#include <cmath>

namespace {

void placeBall(std::vector<Ball>& balls, int id, float x, float y, float radius) {
    Ball& ball = balls[id];
    ball.id = id;
    ball.pos = {x, y};
    ball.vel = {0.0f, 0.0f};
    ball.radius = radius;
    ball.color = ballColor(id);
    ball.active = true;
    ball.isStripe = id >= 9;
    ball.isCue = id == kCueBallId;
}

}  // namespace

void rackBalls(std::vector<Ball>& balls, float rackX, float rackY, float spacing) {
    if (balls.size() < static_cast<size_t>(kBallCount)) {
        balls.resize(kBallCount);
    }

    const float radius = spacing * 0.5f;
    const float rowH = spacing * std::sqrt(3.0f) * 0.5f;

    // 8-ball rack: 1 at apex, 8 in center, alternating solids/stripes.
    const int rows[5][5] = {
        {1, 0, 0, 0, 0},
        {9, 2, 0, 0, 0},
        {3, 8, 10, 0, 0},
        {11, 4, 12, 5, 0},
        {14, 6, 15, 7, 13},
    };

    for (int row = 0; row < 5; ++row) {
        const int count = row + 1;
        const float rowY = rackY + row * rowH;
        const float rowWidth = (count - 1) * spacing;
        const float startX = rackX - rowWidth * 0.5f;

        for (int col = 0; col < count; ++col) {
            const int id = rows[row][col];
            if (id == 0) continue;

            const float x = startX + col * spacing;
            placeBall(balls, id, x, rowY, radius);
        }
    }

    placeBall(balls, kCueBallId, rackX - spacing * 6.0f, rackY, radius);
}
