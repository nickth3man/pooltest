#include "physics.h"

#include "constants.h"

#include <algorithm>
#include <cmath>

namespace Physics {
namespace {

using namespace PhysicsConfig;

Vec2 closestPointOnSegment(const Vec2& p, const Segment& seg) {
    const Vec2 ab = seg.b - seg.a;
    const float abLenSq = ab.lengthSq();
    if (abLenSq < 1e-6f) return seg.a;

    float t = (p - seg.a).dot(ab) / abLenSq;
    t = std::clamp(t, 0.0f, 1.0f);
    return seg.a + ab * t;
}

Vec2 segmentTangent(const Segment& seg) {
    const Vec2 ab = seg.b - seg.a;
    return ab.normalized();
}

float maxBallSpeed(const std::vector<Ball>& balls) {
    float maxSpeedSq = 0.0f;
    for (const Ball& ball : balls) {
        if (!ball.active) continue;
        maxSpeedSq = std::max(maxSpeedSq, ball.vel.lengthSq());
    }
    return std::sqrt(maxSpeedSq);
}

int adaptiveSubSteps(const std::vector<Ball>& balls, float ballRadius) {
    const float speed = maxBallSpeed(balls);
    const float travel = speed / static_cast<float>(kMinSubSteps);
    if (travel <= ballRadius * 0.45f) return kMinSubSteps;

    const int needed = static_cast<int>(std::ceil(speed / (ballRadius * 0.45f)));
    return std::clamp(needed, kMinSubSteps, kMaxSubSteps);
}

void separateOverlap(Ball& a, Ball& b) {
    Vec2 delta = b.pos - a.pos;
    const float minDist = a.radius + b.radius;
    const float distSq = delta.lengthSq();
    if (distSq >= minDist * minDist || distSq < 1e-8f) return;

    const float dist = std::sqrt(distSq);
    Vec2 normal = delta / dist;
    const float overlap = minDist - dist;
    a.pos -= normal * (overlap * 0.5f);
    b.pos += normal * (overlap * 0.5f);
}

void resolveBallBall(Ball& a, Ball& b) {
    Vec2 delta = b.pos - a.pos;
    const float minDist = a.radius + b.radius;
    const float distSq = delta.lengthSq();
    if (distSq >= minDist * minDist || distSq < 1e-8f) return;

    const float dist = std::sqrt(distSq);
    Vec2 normal = delta / dist;
    Vec2 tangent = {-normal.y, normal.x};

    const Vec2 relVel = a.vel - b.vel;
    const float velAlongNormal = relVel.dot(normal);
    if (velAlongNormal > 0.0f) return;

    const float normalImpulse = -(1.0f + kBallRestitution) * velAlongNormal * 0.5f;
    a.vel += normal * normalImpulse;
    b.vel -= normal * normalImpulse;

    const float relTangential = relVel.dot(tangent);
    const float relTangentialSpeed = std::fabs(relTangential);
    if (relTangentialSpeed > 1e-6f) {
        const float mu = dynamicBallFriction(relTangentialSpeed);
        const float maxTangentialImpulse = mu * std::fabs(normalImpulse);
        const float stickImpulse = relTangentialSpeed * 0.5f;
        const float tangentialImpulse =
            -std::copysign(std::min(maxTangentialImpulse, stickImpulse), relTangential);

        a.vel += tangent * tangentialImpulse;
        b.vel -= tangent * tangentialImpulse;
    }

    separateOverlap(a, b);
}

void resolveBallCushion(Ball& ball, const Segment& seg) {
    const Vec2 closest = closestPointOnSegment(ball.pos, seg);
    Vec2 delta = ball.pos - closest;
    const float distSq = delta.lengthSq();
    if (distSq >= ball.radius * ball.radius || distSq < 1e-8f) return;

    const float dist = std::sqrt(distSq);
    Vec2 normal = delta / dist;

    const float velAlongNormal = ball.vel.dot(normal);
    if (velAlongNormal < 0.0f) {
        ball.vel -= normal * ((1.0f + kCushionRestitution) * velAlongNormal);
    }

    const Vec2 tangent = segmentTangent(seg);
    const float velAlongTangent = ball.vel.dot(tangent);
    ball.vel -= tangent * (velAlongTangent * kCushionFriction);

    const float overlap = ball.radius - dist;
    ball.pos += normal * overlap;
}

void applyPockets(Ball& ball, const Table& table, float dt) {
    for (const Pocket& pocket : table.pockets) {
        const Vec2 toPocket = pocket.center - ball.pos;
        const float dist = toPocket.length();
        const float captureDist = pocket.radius - ball.radius * kPocketCaptureFactor;

        if (dist <= captureDist) {
            ball.active = false;
            ball.vel = {0.0f, 0.0f};
            return;
        }

        const float mouthDist = pocket.radius + ball.radius * 0.15f;
        if (dist < mouthDist) {
            const Vec2 pull = toPocket.normalized();
            ball.vel += pull * (kPocketSuckStrength * dt * (1.0f - dist / mouthDist));
        }
    }
}

void applyFriction(Ball& ball, float dt) {
    const float speed = ball.vel.length();
    if (speed <= kStopSpeed) {
        ball.vel = {0.0f, 0.0f};
        return;
    }

    const float decel = kRollingDecel + kSlidingDamping * speed;
    const float newSpeed = std::max(0.0f, speed - decel * dt);
    ball.vel *= newSpeed / speed;

    if (ball.vel.length() < kStopSpeed) {
        ball.vel = {0.0f, 0.0f};
    }
}

void integrate(Ball& ball, float dt) {
    ball.pos += ball.vel * dt;
}

void resolveCollisions(std::vector<Ball>& balls, const Table& table) {
    for (int iter = 0; iter < kPositionIterations; ++iter) {
        for (size_t i = 0; i < balls.size(); ++i) {
            if (!balls[i].active) continue;
            for (size_t j = i + 1; j < balls.size(); ++j) {
                if (!balls[j].active) continue;
                resolveBallBall(balls[i], balls[j]);
            }
        }

        for (Ball& ball : balls) {
            if (!ball.active) continue;
            for (const Segment& seg : table.cushions) {
                resolveBallCushion(ball, seg);
            }
        }
    }
}

float ballRadius(const std::vector<Ball>& balls) {
    for (const Ball& ball : balls) {
        if (ball.active) return ball.radius;
    }
    return 10.0f;
}

}  // namespace

bool allBallsStopped(const std::vector<Ball>& balls) {
    for (const Ball& ball : balls) {
        if (!ball.active) continue;
        if (ball.vel.lengthSq() > kStopSpeed * kStopSpeed) return false;
    }
    return true;
}

RayHit raycastFirstBall(const Vec2& origin,
                        const Vec2& dir,
                        float maxDistance,
                        const std::vector<Ball>& balls,
                        int ignoreId) {
    RayHit best;
    best.distance = maxDistance;

    for (const Ball& ball : balls) {
        if (!ball.active || ball.id == ignoreId) continue;

        const Vec2 oc = origin - ball.pos;
        const float b = oc.dot(dir);
        const float c = oc.lengthSq() - ball.radius * ball.radius;
        const float discriminant = b * b - c;
        if (discriminant < 0.0f) continue;

        const float sqrtD = std::sqrt(discriminant);
        float t = -b - sqrtD;
        if (t < 0.0f) t = -b + sqrtD;
        if (t < 0.0f || t > best.distance) continue;

        best.hit = true;
        best.distance = t;
        best.ballId = ball.id;
        best.point = origin + dir * t;
        best.normal = (best.point - ball.pos).normalized();
    }

    return best;
}

void step(std::vector<Ball>& balls, const Table& table, float dt) {
    const float radius = ballRadius(balls);
    const int subSteps = adaptiveSubSteps(balls, radius);
    const float subDt = dt / static_cast<float>(subSteps);

    for (int stepIndex = 0; stepIndex < subSteps; ++stepIndex) {
        for (Ball& ball : balls) {
            if (!ball.active) continue;
            integrate(ball, subDt);
        }

        resolveCollisions(balls, table);

        for (Ball& ball : balls) {
            if (!ball.active) continue;
            applyPockets(ball, table, subDt);
        }
    }

    for (Ball& ball : balls) {
        if (!ball.active) continue;
        applyFriction(ball, dt);
    }
}

}  // namespace Physics
