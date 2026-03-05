/**
 * Pocket class representing a table pocket
 * Handles pocket position and detection radius
 */

import type { PocketConfig, Vector2 } from "../types.js";

export class Pocket implements PocketConfig {
  constructor(
    public readonly x: number,
    public readonly y: number,
    public readonly drawRadius: number,
    public readonly catchRadius: number
  ) {}

  /** Create a pocket from a configuration object */
  static fromConfig(config: PocketConfig): Pocket {
    return new Pocket(
      config.x,
      config.y,
      config.drawRadius,
      config.catchRadius
    );
  }

  /** Get pocket position as a Vector2 */
  get position(): Vector2 {
    return { x: this.x, y: this.y };
  }

  /** Check if a point is within the catch radius */
  containsPoint(point: Vector2): boolean {
    const dx = point.x - this.x;
    const dy = point.y - this.y;
    return dx * dx + dy * dy <= this.catchRadius * this.catchRadius;
  }

  /** Check if a ball center is close enough to be caught */
  canCatchBall(ballX: number, ballY: number, ballRadius: number): boolean {
    void ballRadius;
    const dx = ballX - this.x;
    const dy = ballY - this.y;
    // Ball is caught if its center enters the catch radius
    return dx * dx + dy * dy <= this.catchRadius * this.catchRadius;
  }

  /** Get distance from pocket to a point */
  distanceTo(point: Vector2): number {
    return Math.hypot(point.x - this.x, point.y - this.y);
  }

  /** Clone the pocket */
  clone(): Pocket {
    return new Pocket(this.x, this.y, this.drawRadius, this.catchRadius);
  }
}
