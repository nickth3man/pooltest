#include "game.h"

#include "physics.h"

#include <algorithm>
#include <cmath>

void Game::reset() {
    table.ballRadius = kBallRadius;
    table.build();
    balls.clear();
    balls.resize(kBallCount);

    const float rackX = table.playArea.x + table.playArea.width * 0.72f;
    const float rackY = table.playArea.y + table.playArea.height * 0.5f;
    rackBalls(balls, rackX, rackY, kBallSpacing);

    ballsMoving = false;
    aiming = false;
    charging = false;
    aimPower = 0.0f;
    aimDir = {1.0f, 0.0f};
    chargeStartMouse = {0.0f, 0.0f};
}

Ball* Game::cueBall() {
    for (Ball& ball : balls) {
        if (ball.isCue && ball.active) return &ball;
    }
    return nullptr;
}

const Ball* Game::cueBall() const {
    for (const Ball& ball : balls) {
        if (ball.isCue && ball.active) return &ball;
    }
    return nullptr;
}

bool Game::canShoot() const {
    return !ballsMoving && cueBall() != nullptr;
}

int Game::activeBallCount() const {
    int count = 0;
    for (const Ball& ball : balls) {
        if (ball.active) ++count;
    }
    return count;
}

void Game::handleInput() {
    ballsMoving = !Physics::allBallsStopped(balls);

    if (IsKeyPressed(KEY_R)) {
        reset();
        return;
    }

    Ball* cue = cueBall();
    if (!cue) {
        aiming = false;
        charging = false;
        return;
    }

    if (!canShoot()) {
        aiming = false;
        charging = false;
        return;
    }

    const Vec2 mouse = {static_cast<float>(GetMouseX()), static_cast<float>(GetMouseY())};
    Vec2 toMouse = mouse - cue->pos;

    aiming = true;

    // Free aim: only update direction when not charging
    if (!charging) {
        if (toMouse.lengthSq() > 1.0f) {
            aimDir = toMouse.normalized();
        }
    }

    if (IsMouseButtonPressed(MOUSE_BUTTON_LEFT)) {
        charging = true;
        aimPower = 0.0f;
        chargeStartMouse = mouse;
        // aimDir is now locked for the duration of the charge
    }

    if (charging && IsMouseButtonDown(MOUSE_BUTTON_LEFT)) {
        // Pull distance measured backward along the *locked* aim direction relative to initial click
        Vec2 dragVec = mouse - chargeStartMouse;
        const float pullback = -(dragVec.dot(aimDir));
        aimPower = std::clamp(pullback / 160.0f, 0.0f, 1.0f);
    }

    if (charging && IsMouseButtonReleased(MOUSE_BUTTON_LEFT)) {
        if (aimPower >= 0.04f) {
            cue->vel = aimDir * (aimPower * kMaxShotSpeed);
            ballsMoving = true;
        }
        charging = false;
        aiming = false;
        aimPower = 0.0f;
    }
}

void Game::update(float dt) {
    handleInput();
    Physics::step(balls, table, dt);
    ballsMoving = !Physics::allBallsStopped(balls);
}
