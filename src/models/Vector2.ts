/**
 * 2D Vector utility class for physics calculations
 * Provides immutable vector operations
 */

import type { Vector2 } from "../types.js";

export class Vec2 implements Vector2 {
  constructor(
    public x: number,
    public y: number
  ) {}

  /** Create a new Vec2 from a Vector2 interface */
  static from(v: Vector2): Vec2 {
    return new Vec2(v.x, v.y);
  }

  /** Zero vector */
  static zero(): Vec2 {
    return new Vec2(0, 0);
  }

  /** Vector with both components set to 1 */
  static one(): Vec2 {
    return new Vec2(1, 1);
  }

  /** Create a unit vector from an angle in radians */
  static fromAngle(angle: number): Vec2 {
    return new Vec2(Math.cos(angle), Math.sin(angle));
  }

  /** Clone this vector */
  clone(): Vec2 {
    return new Vec2(this.x, this.y);
  }

  /** Add another vector */
  add(other: Vector2): Vec2 {
    return new Vec2(this.x + other.x, this.y + other.y);
  }

  /** Subtract another vector */
  sub(other: Vector2): Vec2 {
    return new Vec2(this.x - other.x, this.y - other.y);
  }

  /** Multiply by a scalar */
  mul(scalar: number): Vec2 {
    return new Vec2(this.x * scalar, this.y * scalar);
  }

  /** Divide by a scalar */
  div(scalar: number): Vec2 {
    return new Vec2(this.x / scalar, this.y / scalar);
  }

  /** Negate this vector */
  neg(): Vec2 {
    return new Vec2(-this.x, -this.y);
  }

  /** Dot product with another vector */
  dot(other: Vector2): number {
    return this.x * other.x + this.y * other.y;
  }

  /** Cross product magnitude (2D cross product is a scalar) */
  cross(other: Vector2): number {
    return this.x * other.y - this.y * other.x;
  }

  /** Magnitude squared (faster than magnitude, useful for comparisons) */
  magSq(): number {
    return this.x * this.x + this.y * this.y;
  }

  /** Magnitude (length) of the vector */
  mag(): number {
    return Math.hypot(this.x, this.y);
  }

  /** Distance to another vector */
  distTo(other: Vector2): number {
    const dx = other.x - this.x;
    const dy = other.y - this.y;
    return Math.hypot(dx, dy);
  }

  /** Distance squared to another vector (faster for comparisons) */
  distSqTo(other: Vector2): number {
    const dx = other.x - this.x;
    const dy = other.y - this.y;
    return dx * dx + dy * dy;
  }

  /** Normalize to unit length, return fallback if magnitude is too small */
  normalize(fallback: Vector2 = { x: 1, y: 0 }): Vec2 {
    const m = this.mag();
    if (m < 1e-7) {
      return new Vec2(fallback.x, fallback.y);
    }
    return this.div(m);
  }

  /** Linear interpolation between this and another vector */
  lerp(other: Vector2, t: number): Vec2 {
    return new Vec2(
      this.x + (other.x - this.x) * t,
      this.y + (other.y - this.y) * t
    );
  }

  /** Project this vector onto another vector */
  projectOnto(other: Vector2): Vec2 {
    const otherMagSq = other.x * other.x + other.y * other.y;
    if (otherMagSq < 1e-10) {
      return Vec2.zero();
    }
    const scale = this.dot(other) / otherMagSq;
    return new Vec2(other.x * scale, other.y * scale);
  }

  /** Reflect this vector across a normal */
  reflect(normal: Vector2): Vec2 {
    const n = Vec2.from(normal).normalize();
    const dot = this.dot(n);
    return this.sub(n.mul(2 * dot));
  }

  /** Rotate by an angle in radians */
  rotate(angle: number): Vec2 {
    const cos = Math.cos(angle);
    const sin = Math.sin(angle);
    return new Vec2(
      this.x * cos - this.y * sin,
      this.x * sin + this.y * cos
    );
  }

  /** Check if this vector equals another (within epsilon) */
  equals(other: Vector2, epsilon: number = 1e-10): boolean {
    return Math.abs(this.x - other.x) < epsilon && Math.abs(this.y - other.y) < epsilon;
  }

  /** Convert to plain object */
  toObject(): Vector2 {
    return { x: this.x, y: this.y };
  }

  /** String representation */
  toString(): string {
    return `Vec2(${this.x.toFixed(2)}, ${this.y.toFixed(2)})`;
  }
}

/** Utility functions for Vector2 operations */
export namespace Vector2Utils {
  /** Clamp a value between min and max */
  export function clamp(value: number, min: number, max: number): number {
    return Math.max(min, Math.min(max, value));
  }

  /** Calculate distance between two points */
  export function dist(a: Vector2, b: Vector2): number {
    return Math.hypot(b.x - a.x, b.y - a.y);
  }

  /** Calculate distance squared (faster, for comparisons) */
  export function distSq(a: Vector2, b: Vector2): number {
    const dx = b.x - a.x;
    const dy = b.y - a.y;
    return dx * dx + dy * dy;
  }

  /** Normalize a vector with a fallback */
  export function normalize(v: Vector2, fallback: Vector2 = { x: 1, y: 0 }): Vector2 {
    const m = Math.hypot(v.x, v.y);
    if (m < 1e-7) {
      return { x: fallback.x, y: fallback.y };
    }
    return { x: v.x / m, y: v.y / m };
  }

  /** Linear interpolation between two vectors */
  export function lerp(a: Vector2, b: Vector2, t: number): Vector2 {
    return {
      x: a.x + (b.x - a.x) * t,
      y: a.y + (b.y - a.y) * t
    };
  }
}
