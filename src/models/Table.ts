/**
 * Table class representing the billiards table
 * Manages dimensions, felt area, spots, and pocket positions
 */

import type { TableConfig, Vector2 } from "../types.js";
import { Pocket } from "./Pocket.js";
import {
  TABLE_X,
  TABLE_Y,
  TABLE_WIDTH,
  TABLE_HEIGHT,
  TABLE_RAIL,
  BALL_DIAMETER,
  RACK_ROWS
} from "../constants.js";
import { Vec2 } from "./Vector2.js";

export class Table {
  public readonly feltX: number;
  public readonly feltY: number;
  public readonly feltWidth: number;
  public readonly feltHeight: number;
  public readonly headSpot: Vec2;
  public readonly footSpot: Vec2;
  public readonly pockets: Pocket[];

  constructor(
    public readonly x: number = TABLE_X,
    public readonly y: number = TABLE_Y,
    public readonly width: number = TABLE_WIDTH,
    public readonly height: number = TABLE_HEIGHT,
    public readonly rail: number = TABLE_RAIL
  ) {
    // Calculate felt area (playable surface inside rails)
    this.feltX = this.x + this.rail;
    this.feltY = this.y + this.rail;
    this.feltWidth = this.width - this.rail * 2;
    this.feltHeight = this.height - this.rail * 2;

    // Calculate spot positions
    this.headSpot = new Vec2(
      this.feltX + this.feltWidth * 0.25,
      this.feltY + this.feltHeight * 0.5
    );

    this.footSpot = new Vec2(
      this.feltX + this.feltWidth * 0.75,
      this.feltY + this.feltHeight * 0.5
    );

    // Create pockets
    this.pockets = this.createPockets();
  }

  /** Create the default table configuration */
  static createDefault(): Table {
    return new Table();
  }

  /** Create table from a configuration object */
  static fromConfig(config: TableConfig): Table {
    return new Table(
      config.x,
      config.y,
      config.width,
      config.height,
      config.rail
    );
  }

  /** Create the 6 pockets on the table */
  private createPockets(): Pocket[] {
    return [
      // Top pockets
      new Pocket(this.feltX, this.feltY, 25, 17),
      new Pocket(this.feltX + this.feltWidth / 2, this.feltY, 22, 15),
      new Pocket(this.feltX + this.feltWidth, this.feltY, 25, 17),
      // Bottom pockets
      new Pocket(this.feltX, this.feltY + this.feltHeight, 25, 17),
      new Pocket(this.feltX + this.feltWidth / 2, this.feltY + this.feltHeight, 22, 15),
      new Pocket(this.feltX + this.feltWidth, this.feltY + this.feltHeight, 25, 17)
    ];
  }

  /** Get felt area bounds */
  get feltBounds(): {
    left: number;
    right: number;
    top: number;
    bottom: number;
  } {
    return {
      left: this.feltX,
      right: this.feltX + this.feltWidth,
      top: this.feltY,
      bottom: this.feltY + this.feltHeight
    };
  }

  /** Get the minimum X coordinate for a ball (accounting for radius) */
  getMinX(ballRadius: number): number {
    return this.feltX + ballRadius;
  }

  /** Get the maximum X coordinate for a ball (accounting for radius) */
  getMaxX(ballRadius: number): number {
    return this.feltX + this.feltWidth - ballRadius;
  }

  /** Get the minimum Y coordinate for a ball (accounting for radius) */
  getMinY(ballRadius: number): number {
    return this.feltY + ballRadius;
  }

  /** Get the maximum Y coordinate for a ball (accounting for radius) */
  getMaxY(ballRadius: number): number {
    return this.feltY + this.feltHeight - ballRadius;
  }

  /** Check if a position is within the felt area (for ball placement) */
  isWithinFelt(x: number, y: number, radius: number = 0): boolean {
    return (
      x >= this.feltX + radius &&
      x <= this.feltX + this.feltWidth - radius &&
      y >= this.feltY + radius &&
      y <= this.feltY + this.feltHeight - radius
    );
  }

  /** Generate rack positions for the numbered balls at the foot spot */
  generateRackPositions(): Vector2[] {
    const rowSpacing = Math.sqrt(3) * BALL_DIAMETER / 2;
    const positions: Vector2[] = [];

    for (let row = 0; row < RACK_ROWS.length; row++) {
      const count = RACK_ROWS[row];
      for (let i = 0; i < count; i++) {
        positions.push({
          x: this.footSpot.x + row * rowSpacing,
          y: this.footSpot.y + (i - (count - 1) / 2) * BALL_DIAMETER
        });
      }
    }

    return positions;
  }

  /** Find a valid position for respawning the cue ball */
  findRespawnPosition(
    basePosition: Vector2,
    ballRadius: number,
    isOccupied: (x: number, y: number) => boolean
  ): Vector2 | null {
    const maxTries = 40;
    let x = basePosition.x;
    let y = basePosition.y;

    for (let tries = 0; tries < maxTries; tries++) {
      if (!isOccupied(x, y)) {
        return { x, y };
      }

      // Try next position to the right
      x += ballRadius * 1.6;

      // If we've gone too far right, move down and reset x
      if (x > this.feltX + this.feltWidth * 0.4) {
        x = basePosition.x;
        y += ballRadius * 1.5;
      }

      // If we've gone too far down, reset to base
      if (y > this.feltY + this.feltHeight - ballRadius * 2) {
        y = basePosition.y;
      }
    }

    return null;
  }

  /** Get diamond marker positions */
  getDiamondPositions(): {
    top: Vector2[];
    bottom: Vector2[];
    left: Vector2[];
    right: Vector2[];
  } {
    const topY = this.y + this.rail * 0.47;
    const bottomY = this.y + this.height - this.rail * 0.47;
    const leftX = this.x + this.rail * 0.47;
    const rightX = this.x + this.width - this.rail * 0.47;

    const top: Vector2[] = [];
    const bottom: Vector2[] = [];
    const left: Vector2[] = [];
    const right: Vector2[] = [];

    // Horizontal diamonds
    for (let i = 1; i <= 3; i++) {
      const x = this.feltX + (this.feltWidth * i) / 4;
      top.push({ x, y: topY });
      bottom.push({ x, y: bottomY });
    }

    // Vertical diamonds
    for (let i = 1; i <= 2; i++) {
      const y = this.feltY + (this.feltHeight * i) / 3;
      left.push({ x: leftX, y });
      right.push({ x: rightX, y });
    }

    return { top, bottom, left, right };
  }

  /** Convert to plain configuration object */
  toConfig(): TableConfig {
    return {
      x: this.x,
      y: this.y,
      width: this.width,
      height: this.height,
      rail: this.rail,
      felt: {
        x: this.feltX,
        y: this.feltY,
        width: this.feltWidth,
        height: this.feltHeight
      },
      headSpot: this.headSpot.toObject(),
      footSpot: this.footSpot.toObject()
    };
  }
}
