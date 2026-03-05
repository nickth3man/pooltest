/**
 * Core type definitions for the Billiards game
 */

/** 2D vector/point interface */
export interface Vector2 {
  x: number;
  y: number;
}

/** Ball interface representing a billiard ball */
export interface BallConfig {
  number: number;
  color: string;
  x: number;
  y: number;
  isCue: boolean;
}

/** Ball state extending config with physics properties */
export interface BallState extends BallConfig {
  vx: number;
  vy: number;
  radius: number;
  inPlay: boolean;
}

/** Pocket configuration */
export interface PocketConfig {
  x: number;
  y: number;
  drawRadius: number;
  catchRadius: number;
}

/** Pocket flash effect for visual feedback */
export interface PocketFlash {
  x: number;
  y: number;
  life: number;
}

/** Table dimensions and layout */
export interface TableDimensions {
  x: number;
  y: number;
  width: number;
  height: number;
  rail: number;
}

/** Felt area (playable surface) */
export interface FeltArea {
  x: number;
  y: number;
  width: number;
  height: number;
}

/** Spot positions on the table */
export interface TableSpots {
  headSpot: Vector2;
  footSpot: Vector2;
}

/** Complete table configuration */
export interface TableConfig extends TableDimensions {
  felt: FeltArea;
  headSpot: Vector2;
  footSpot: Vector2;
}

/** Game state enumeration */
export enum GameState {
  AIMING = "AIMING",
  SHOOTING = "SHOOTING",
  BALLS_MOVING = "BALLS_MOVING"
}

/** Mouse input state */
export interface MouseState {
  x: number;
  y: number;
}

/** Drag input state for shooting */
export interface DragState {
  isDragging: boolean;
  startX: number;
  startY: number;
  pullDistance: number;
}

/** Aim direction and power */
export interface AimState {
  direction: Vector2;
  lockedDirection: Vector2;
}

/** HUD/UI element references */
export interface HUDElements {
  sunkCount: HTMLElement;
  stateText: HTMLElement;
  shootIndicator: HTMLElement;
  restartButton: HTMLButtonElement;
}

/** Canvas rendering context with required methods */
export type RenderContext = CanvasRenderingContext2D;

/** Audio context wrapper */
export interface AudioContextState {
  context: AudioContext | null;
}

/** Physics simulation configuration */
export interface PhysicsConfig {
  friction: number;
  cushionRestitution: number;
  ballRestitution: number;
  minSpeed: number;
  substeps: number;
}

/** Shot configuration */
export interface ShotConfig {
  maxPower: number;
  powerScale: number;
}

/** Game configuration combining all sub-configs */
export interface GameConfig {
  width: number;
  height: number;
  ballRadius: number;
  physics: PhysicsConfig;
  shot: ShotConfig;
  scratchRespawnDelay: number;
  shotAnimationDuration: number;
  pocketFlashDuration: number;
}

/** Event callback types */
export type GameEventCallback = () => void;
export type BallSunkCallback = (ballNumber: number) => void;
export type ScratchCallback = () => void;
