import { BALL_COLORS, GAME_CONFIG } from "../constants.js";
import { AudioManager } from "../audio/AudioManager.js";
import { EventBus } from "../core/EventBus.js";
import { Ball } from "../models/Ball.js";
import { Table } from "../models/Table.js";
import type { Vector2 } from "../types.js";
import { PhysicsEngine } from "../physics/PhysicsEngine.js";
import { Renderer, type RenderState } from "../rendering/Renderer.js";
import { InputHandler } from "./InputHandler.js";
import { GameRules, type RuleResult } from "./GameRules.js";
import { GameStateManager, type GameState } from "./GameState.js";
import type { GameStatusSnapshot } from "./GamePresenter.js";

/**
 * GameSession - Core game controller that orchestrates all systems
 * 
 * Responsibilities:
 * - Manages game objects (balls, table)
 * - Coordinates update/render cycle
 * - Handles game rules and state transitions
 * - Bridges physics, rendering, audio, and input
 */
export class GameSession {
  readonly table: Table;
  readonly stateManager: GameStateManager;
  readonly gameRules: GameRules;
  readonly inputHandler: InputHandler;
  readonly physicsEngine: PhysicsEngine;
  readonly renderer: Renderer;
  readonly audioManager: AudioManager;
  readonly eventBus: EventBus;

  private balls: Ball[] = [];
  private cueBall!: Ball;
  private shotTimer: number = 0;
  private scratchTimer: number = 0;

  constructor(canvas: HTMLCanvasElement, eventBus: EventBus = new EventBus()) {
    this.eventBus = eventBus;
    this.table = Table.createDefault();
    this.stateManager = new GameStateManager();
    this.gameRules = new GameRules();
    this.renderer = new Renderer(canvas);
    this.audioManager = new AudioManager();
    this.physicsEngine = new PhysicsEngine(this.table, this.eventBus);
    this.inputHandler = new InputHandler(canvas, this.eventBus);
    this.reset();
  }

  reset(): void {
    this.balls = [];
    this.gameRules.reset();
    this.shotTimer = 0;
    this.scratchTimer = 0;
    this.inputHandler.reset();
    this.stateManager.reset();

    // Place cue ball at head spot (left side of table)
    this.cueBall = Ball.createCueBall(this.table.headSpot.x, this.table.headSpot.y);
    this.balls.push(this.cueBall);

    // Rack numbered balls in triangle formation at foot spot
    const rackPositions = this.table.generateRackPositions();
    for (let i = 1; i <= 7; i++) {
      const position = rackPositions[i - 1];
      const ball = Ball.createNumberedBall(i, BALL_COLORS[i], position.x, position.y);
      this.balls.push(ball);
    }
  }

  update(dt: number): void {
    // Update aim direction based on mouse position relative to cue ball
    if (this.stateManager.isAiming && this.cueBall.inPlay) {
      this.inputHandler.updateAimFromCueBall(this.cueBall.position);
    }

    // Handle shooting state countdown (shows cue stick animation)
    if (this.stateManager.isShooting) {
      this.shotTimer -= dt;
      if (this.shotTimer <= 0) {
        this.stateManager.startBallsMoving();
      }
    }

    // Run physics simulation during shooting or balls moving states
    if (this.stateManager.isShooting || this.stateManager.isBallsMoving) {
      this.physicsEngine.simulate(this.balls);

      // Handle cue ball respawn delay after scratch
      if (this.scratchTimer > 0) {
        this.scratchTimer -= dt;
        if (this.scratchTimer <= 0) {
          this.respawnCueBall();
        }
      }

      // Return to aiming when all balls stop (and no respawn pending)
      if (!this.physicsEngine.areBallsMoving(this.balls) && this.scratchTimer <= 0 && !this.stateManager.isAiming) {
        this.stateManager.returnToAiming();
      }
    }

    // Update visual effects (pocket flashes fade out over time)
    this.gameRules.updatePocketFlashes(dt);
  }

  render(): void {
    this.renderer.render(this.getRenderState());
  }

  handleBallSunk(ball: Ball): RuleResult {
    return this.gameRules.handleBallSunk(ball);
  }

  handleScratch(cueBall: Ball): RuleResult {
    const result = this.gameRules.handleScratch(cueBall);
    if (result.shouldSpawnCueBall) {
      this.scratchTimer = GAME_CONFIG.scratchRespawnDelay;
    }
    return result;
  }

  handleShot(power: number, direction: Vector2): void {
    if (power <= 0.2 || !this.cueBall.inPlay) {
      return;
    }

    this.physicsEngine.applyShot(this.cueBall, power, direction);
    this.stateManager.startShooting();
    this.shotTimer = GAME_CONFIG.shotAnimationDuration;
    this.audioManager.playCueHitSound(power);
  }

  ensureAudioContext(): void {
    this.audioManager.ensureContext();
  }

  playPocketSound(): void {
    this.audioManager.playPocketSound();
  }

  playCollisionSound(impactForce: number): void {
    this.audioManager.playCollisionSound(impactForce);
  }

  playCushionSound(): void {
    this.audioManager.playCushionSound();
  }

  getState(): GameStatusSnapshot {
    return {
      sunkCount: this.gameRules.sunkCount,
      gameState: this.stateManager.state,
      readyToShoot: this.stateManager.readyToShoot && this.cueBall.inPlay
    };
  }

  getGameState(): GameState {
    return this.stateManager.state;
  }

  isReadyToShoot(): boolean {
    return this.stateManager.readyToShoot && this.cueBall.inPlay;
  }

  destroy(): void {
    this.inputHandler.destroy();
    this.audioManager.destroy();
    this.eventBus.clear();
  }

  private getRenderState(): RenderState {
    return {
      balls: this.balls,
      cueBall: this.cueBall,
      table: this.table,
      gameState: this.stateManager.state,
      aimDirection: this.inputHandler.aimDirection,
      lockedAimDirection: this.inputHandler.lockedAimDirection,
      isDragging: this.inputHandler.isDragging,
      pullDistance: this.inputHandler.pullDistance,
      pocketFlashes: this.gameRules.getPocketFlashes(),
      sunkNumbers: this.gameRules.getSunkNumbers()
    };
  }

  private respawnCueBall(): void {
    const isOccupied = (x: number, y: number): boolean => {
      return this.physicsEngine.isPositionOccupied(x, y, this.cueBall.radius, this.balls, this.cueBall);
    };

    const newPosition = this.table.findRespawnPosition(this.table.headSpot, this.cueBall.radius, isOccupied);
    if (newPosition) {
      this.cueBall.respawn(newPosition.x, newPosition.y);
      return;
    }

    this.cueBall.respawn(this.table.headSpot.x, this.table.headSpot.y);
  }
}
