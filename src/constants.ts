/**
 * Game constants - extracted from the original monolithic file
 * These are immutable constants that define the game behavior
 */

import type { PhysicsConfig, ShotConfig, GameConfig } from "./types.js";

/** Canvas dimensions */
export const CANVAS_WIDTH = 1100;
export const CANVAS_HEIGHT = 500;

/** Ball geometry */
export const BALL_RADIUS = 11;
export const BALL_DIAMETER = BALL_RADIUS * 2;

/** Table layout constants */
export const TABLE_X = 58;
export const TABLE_Y = 34;
export const TABLE_WIDTH = 784;
export const TABLE_HEIGHT = 432;
export const TABLE_RAIL = 38;

/** Physics constants */
export const PHYSICS: PhysicsConfig = {
  friction: 0.985,
  cushionRestitution: 0.7,
  ballRestitution: 0.95,
  minSpeed: 0.05,
  substeps: 8,
};

/** Shot/power constants */
export const SHOT: ShotConfig = {
  maxPower: 22,
  powerScale: 0.157,
};

/** Timing constants (in milliseconds) */
export const TIMING = {
  scratchRespawnDelay: 850,
  shotAnimationDuration: 90,
  pocketFlashDuration: 220,
} as const;

/** UI/Visual constants */
export const VISUAL = {
  guideLength: 220,
  guideStartOffset: 1,
  maxPullDistance: 200,
  cueStickLength: 175,
  cueButtLength: 26,
  powerMeterX: 1060, // WIDTH - 40
  powerMeterY: 120,
  powerMeterHeight: 240,
  powerMeterWidth: 14,
  sunkRowStartX: 18,
  sunkRowY: 20,
  miniBallRadius: 7,
} as const;

/** Complete game configuration */
export const GAME_CONFIG: GameConfig = {
  width: CANVAS_WIDTH,
  height: CANVAS_HEIGHT,
  ballRadius: BALL_RADIUS,
  physics: PHYSICS,
  shot: SHOT,
  scratchRespawnDelay: TIMING.scratchRespawnDelay,
  shotAnimationDuration: TIMING.shotAnimationDuration,
  pocketFlashDuration: TIMING.pocketFlashDuration,
};

/** Ball colors by number (1-7) */
export const BALL_COLORS: Record<number, string> = {
  1: "#f4d73b", // Yellow
  2: "#3060d8", // Blue
  3: "#d83434", // Red
  4: "#6f3eb8", // Purple
  5: "#e98927", // Orange
  6: "#1f944e", // Green
  7: "#6b2330", // Maroon
};

/** Rack configuration */
export const RACK_ROWS = [1, 2, 4] as const;

/** Diamond marker positions */
export const DIAMOND_COUNT_HORIZONTAL = 3;
export const DIAMOND_COUNT_VERTICAL = 2;
