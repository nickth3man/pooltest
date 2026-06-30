#pragma once

#include <cmath>

// Tuned for 2D play. Ball-ball values inspired by tailuge/billiards (Alciatore
// throw, Han friction). Cushion values are a 2D simplification of published
// billiards research (Mathavan / Stronge).

namespace PhysicsConfig {

constexpr float kBallRestitution = 0.86f;
constexpr float kCushionRestitution = 0.77f;

// Alciatore dynamic ball-ball friction: mu(v) = a + b * exp(-c * v)
constexpr float kThrowFrictionA = 0.01f;
constexpr float kThrowFrictionB = 0.108f;
constexpr float kThrowFrictionC = 1.088f;

// Han-style table friction (screen units per second squared).
constexpr float kRollingDecel = 42.0f;
constexpr float kSlidingDamping = 1.35f;

constexpr float kCushionFriction = 0.18f;
constexpr float kStopSpeed = 6.0f;

constexpr int kMinSubSteps = 10;
constexpr int kMaxSubSteps = 32;
constexpr int kPositionIterations = 4;

constexpr float kPocketCaptureFactor = 0.35f;
constexpr float kPocketSuckStrength = 280.0f;

inline float dynamicBallFriction(float relativeTangentialSpeed) {
    return kThrowFrictionA +
           kThrowFrictionB *
               std::exp(-kThrowFrictionC * relativeTangentialSpeed);
}

}  // namespace PhysicsConfig
