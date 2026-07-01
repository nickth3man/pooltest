#pragma once

extern "C" {
#include "raylib.h"
}

#include <vector>

#include "ball.h"
#include "table.h"

struct Game {
    Table table;
    std::vector<Ball> balls;

    bool ballsMoving = false;
    bool aiming = false;
    bool charging = false;
    Vec2 aimDir = {1.0f, 0.0f};
    Vec2 chargeStartMouse = {0.0f, 0.0f};
    float aimPower = 0.0f;

    static constexpr float kMaxShotSpeed = 900.0f;
    static constexpr float kBallRadius = 10.0f;
    static constexpr float kBallSpacing = 21.5f;

    void reset();
    void update(float dt);
    void handleInput();
    Ball* cueBall();
    const Ball* cueBall() const;
    bool canShoot() const;
    int activeBallCount() const;
};
