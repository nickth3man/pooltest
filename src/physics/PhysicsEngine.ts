/**
 * PhysicsEngine class
 * Handles all physics simulation including movement, collisions, and pocket detection
 */

import type { Vector2 } from "../types.js";
import { Ball } from "../models/Ball.js";
import { Table } from "../models/Table.js";
import { PHYSICS, BALL_RADIUS } from "../constants.js";
import { EventBus, EventType } from "../core/EventBus.js";



export class PhysicsEngine {
  private table: Table;
  private eventBus: EventBus;
  private stepFriction: number;
  private lastSimulationCollisionCount: number = 0;

  constructor(table: Table, eventBus: EventBus) {
    this.table = table;
    this.eventBus = eventBus;
    // Pre-calculate per-step friction for performance
    this.stepFriction = Math.pow(PHYSICS.friction, 1 / PHYSICS.substeps);
  }

  /** Simulate physics for one frame */
  simulate(balls: Ball[]): void {
    this.lastSimulationCollisionCount = 0;

    for (let step = 0; step < PHYSICS.substeps; step++) {
      // Update ball positions
      this.updatePositions(balls);

      // Handle cushion collisions
      this.handleCushions(balls);

      // Resolve ball-to-ball collisions
      this.resolveCollisions(balls);

      // Detect pocket sinks
      this.detectPockets(balls);
    }
  }

  /** Update ball positions based on velocity */
  private updatePositions(balls: Ball[]): void {
    const dt = 1 / PHYSICS.substeps;

    for (const ball of balls) {
      if (!ball.inPlay) continue;

      // Update position
      ball.x += ball.vx * dt;
      ball.y += ball.vy * dt;

      // Apply friction
      ball.vx *= this.stepFriction;
      ball.vy *= this.stepFriction;

      // Stop if below minimum speed
      if (ball.speed < PHYSICS.minSpeed) {
        ball.stop();
      }
    }
  }

  /** Handle ball collisions with table cushions */
  private handleCushions(balls: Ball[]): void {
    const minX = this.table.getMinX(BALL_RADIUS);
    const maxX = this.table.getMaxX(BALL_RADIUS);
    const minY = this.table.getMinY(BALL_RADIUS);
    const maxY = this.table.getMaxY(BALL_RADIUS);

    for (const ball of balls) {
      if (!ball.inPlay) continue;

      let hitCushion = false;

      // Left cushion
      if (ball.x < minX) {
        ball.x = minX;
        ball.vx = Math.abs(ball.vx) * PHYSICS.cushionRestitution;
        hitCushion = true;
      }
      // Right cushion
      else if (ball.x > maxX) {
        ball.x = maxX;
        ball.vx = -Math.abs(ball.vx) * PHYSICS.cushionRestitution;
        hitCushion = true;
      }

      // Top cushion
      if (ball.y < minY) {
        ball.y = minY;
        ball.vy = Math.abs(ball.vy) * PHYSICS.cushionRestitution;
        hitCushion = true;
      }
      // Bottom cushion
      else if (ball.y > maxY) {
        ball.y = maxY;
        ball.vy = -Math.abs(ball.vy) * PHYSICS.cushionRestitution;
        hitCushion = true;
      }

      if (hitCushion) {
        this.eventBus.emit({
          type: EventType.CUSHION_HIT,
          ball
        });
      }
    }
  }

  /** Resolve ball-to-ball collisions */
  private resolveCollisions(balls: Ball[]): void {
    for (let i = 0; i < balls.length; i++) {
      const ballA = balls[i];
      if (!ballA.inPlay) continue;

      for (let j = i + 1; j < balls.length; j++) {
        const ballB = balls[j];
        if (!ballB.inPlay) continue;

        if (this.resolveBallCollision(ballA, ballB)) {
          this.lastSimulationCollisionCount++;
        }
      }
    }
  }

  /** Resolve collision between two balls */
  private resolveBallCollision(a: Ball, b: Ball): boolean {
    const dx = b.x - a.x;
    const dy = b.y - a.y;
    const distSq = dx * dx + dy * dy;
    const minDist = a.radius + b.radius;

    // No collision
    if (distSq > minDist * minDist) return false;

    let nx: number;
    let ny: number;
    let dist: number;

    if (distSq === 0) {
      const seed = (a.number + 1) * 92821 + (b.number + 1) * 68917;
      const angle = (seed % 360) * (Math.PI / 180);
      nx = Math.cos(angle);
      ny = Math.sin(angle);
      dist = 0;
    } else {
      dist = Math.sqrt(distSq);
      nx = dx / dist;
      ny = dy / dist;
    }

    const overlap = minDist - dist;

    // Separate balls to prevent sticking
    a.x -= nx * overlap * 0.5;
    a.y -= ny * overlap * 0.5;
    b.x += nx * overlap * 0.5;
    b.y += ny * overlap * 0.5;

    // Calculate relative velocity
    const rvx = b.vx - a.vx;
    const rvy = b.vy - a.vy;
    const velAlongNormal = rvx * nx + rvy * ny;

    // Do not resolve if velocities are separating
    if (velAlongNormal > 0) return true;

    // Calculate impulse with restitution
    const impulse = -(1 + PHYSICS.ballRestitution) * velAlongNormal / 2;
    const impactForce = Math.abs(velAlongNormal);

    // Apply impulse
    a.vx -= impulse * nx;
    a.vy -= impulse * ny;
    b.vx += impulse * nx;
    b.vy += impulse * ny;

    if (impactForce > 0) {
      this.eventBus.emit({
        type: EventType.COLLISION,
        ballA: a,
        ballB: b,
        impactForce
      });
    }

    return true;
  }

  /** Detect and handle balls entering pockets */
  private detectPockets(balls: Ball[]): void {
    for (const ball of balls) {
      if (!ball.inPlay) continue

      for (const pocket of this.table.pockets) {
        if (pocket.canCatchBall(ball.x, ball.y, ball.radius)) {
          ball.sink()

          if (ball.isCue) {
            this.eventBus.emit({
              type: EventType.SCRATCH,
              cueBall: ball
            });
          } else {
            this.eventBus.emit({
              type: EventType.BALL_SUNK,
              ball
            });
          }
          break
        }
      }
    }
  }

  /** Check if any balls are still moving */
  areBallsMoving(balls: Ball[]): boolean {
    for (const ball of balls) {
      if (!ball.inPlay) continue;
      if (ball.speed > 0) return true;
    }
    return false;
  }

  /** Check if a position overlaps with any in-play ball */
  isPositionOccupied(x: number, y: number, radius: number, balls: Ball[], excludeBall?: Ball): boolean {
    for (const ball of balls) {
      if (!ball.inPlay || ball === excludeBall) continue;

      const dx = ball.x - x;
      const dy = ball.y - y;
      const minDist = radius + ball.radius + 1; // +1 for safety margin

      if (dx * dx + dy * dy < minDist * minDist) {
        return true;
      }
    }
    return false;
  }

  /** Find a valid position near the given position */
  findValidPosition(
    baseX: number,
    baseY: number,
    radius: number,
    balls: Ball[],
    excludeBall?: Ball
  ): Vector2 | null {
    const maxTries = 40;
    const step = radius * 1.6;

    for (let tries = 0; tries < maxTries; tries++) {
      if (!this.isPositionOccupied(baseX, baseY, radius, balls, excludeBall)) {
        return { x: baseX, y: baseY };
      }
      baseX += step;
    }

    return null;
  }

  /** Apply a shot impulse to the cue ball */
  applyShot(cueBall: Ball, power: number, direction: Vector2): void {
    cueBall.vx = direction.x * power;
    cueBall.vy = direction.y * power;
  }

  /** Get total kinetic energy of all balls (for debugging) */
  getTotalEnergy(balls: Ball[]): number {
    let energy = 0;
    for (const ball of balls) {
      if (!ball.inPlay) continue;
      energy += 0.5 * ball.speedSq;
    }
    return energy;
  }

  getLastSimulationCollisionCount(): number {
    return this.lastSimulationCollisionCount;
  }

  /** Reset the physics engine */
  reset(): void {
    this.lastSimulationCollisionCount = 0;
    this.stepFriction = Math.pow(PHYSICS.friction, 1 / PHYSICS.substeps);
  }
}
