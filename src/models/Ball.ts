/**
 * Ball class representing a billiard ball
 * Encapsulates physics state and rendering properties
 */

import type { BallConfig, BallState, Vector2 } from "../types.js";
import { BALL_RADIUS } from "../constants.js";
import { Vec2 } from "./Vector2.js";

export class Ball implements BallState {
  public vx: number;
  public vy: number;
  public inPlay: boolean;
  public readonly radius: number;

  constructor(
    public readonly number: number,
    public readonly color: string,
    public x: number,
    public y: number,
    public readonly isCue: boolean
  ) {
    this.vx = 0;
    this.vy = 0;
    this.inPlay = true;
    this.radius = BALL_RADIUS;
  }

  /** Create a ball from a configuration object */
  static fromConfig(config: BallConfig): Ball {
    return new Ball(
      config.number,
      config.color,
      config.x,
      config.y,
      config.isCue
    );
  }

  /** Create the cue ball at a specific position */
  static createCueBall(x: number, y: number): Ball {
    return new Ball(0, "#ffffff", x, y, true);
  }

  /** Create a numbered ball */
  static createNumberedBall(number: number, color: string, x: number, y: number): Ball {
    return new Ball(number, color, x, y, false);
  }

  /** Get current position as a Vec2 */
  get position(): Vec2 {
    return new Vec2(this.x, this.y);
  }

  /** Set position from a Vector2 */
  set position(pos: Vector2) {
    this.x = pos.x;
    this.y = pos.y;
  }

  /** Get current velocity as a Vec2 */
  get velocity(): Vec2 {
    return new Vec2(this.vx, this.vy);
  }

  /** Set velocity from a Vector2 */
  set velocity(vel: Vector2) {
    this.vx = vel.x;
    this.vy = vel.y;
  }

  /** Get speed (magnitude of velocity) */
  get speed(): number {
    return Math.hypot(this.vx, this.vy);
  }

  /** Get speed squared (faster for comparisons) */
  get speedSq(): number {
    return this.vx * this.vx + this.vy * this.vy;
  }

  /** Check if the ball is moving */
  get isMoving(): boolean {
    return this.speed > 0;
  }

  /** Stop the ball (zero velocity) */
  stop(): void {
    this.vx = 0;
    this.vy = 0;
  }

  /** Apply an impulse (change in velocity) */
  applyImpulse(impulse: Vector2): void {
    this.vx += impulse.x;
    this.vy += impulse.y;
  }

  /** Apply a force over time (scales by delta time) */
  applyForce(force: Vector2, dt: number): void {
    this.vx += force.x * dt;
    this.vy += force.y * dt;
  }

  /** Update position based on velocity */
  updatePosition(dt: number): void {
    this.x += this.vx * dt;
    this.y += this.vy * dt;
  }

  /** Apply friction to slow down */
  applyFriction(friction: number): void {
    this.vx *= friction;
    this.vy *= friction;
  }

  /** Check if ball overlaps with another ball */
  overlaps(other: Ball): boolean {
    if (!this.inPlay || !other.inPlay) return false;
    if (other === this) return false;
    const minDist = this.radius + other.radius;
    const dx = other.x - this.x;
    const dy = other.y - this.y;
    return dx * dx + dy * dy < minDist * minDist;
  }

  /** Check if ball is within a certain distance of a point */
  isNearPoint(point: Vector2, distance: number): boolean {
    const dx = point.x - this.x;
    const dy = point.y - this.y;
    return dx * dx + dy * dy <= distance * distance;
  }

  /** Distance to another ball */
  distanceTo(other: Ball): number {
    return Math.hypot(other.x - this.x, other.y - this.y);
  }

  /** Distance squared to another ball (faster for comparisons) */
  distanceSqTo(other: Ball): number {
    const dx = other.x - this.x;
    const dy = other.y - this.y;
    return dx * dx + dy * dy;
  }

  /** Sink the ball (remove from play) */
  sink(): void {
    this.inPlay = false;
    this.stop();
  }

  /** Respawn the ball at a new position */
  respawn(x: number, y: number): void {
    this.x = x;
    this.y = y;
    this.stop();
    this.inPlay = true;
  }

  /** Clone the ball's state */
  clone(): Ball {
    const cloned = new Ball(this.number, this.color, this.x, this.y, this.isCue);
    cloned.vx = this.vx;
    cloned.vy = this.vy;
    cloned.inPlay = this.inPlay;
    return cloned;
  }

  /** Serialize to a plain object */
  toJSON(): BallState {
    return {
      number: this.number,
      color: this.color,
      x: this.x,
      y: this.y,
      vx: this.vx,
      vy: this.vy,
      radius: this.radius,
      inPlay: this.inPlay,
      isCue: this.isCue
    };
  }
}
