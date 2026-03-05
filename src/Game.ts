import { Ball } from "./models/Ball.js";
import { Table } from "./models/Table.js";
import { GameStateManager, GameState } from "./game/GameState.js";
import { GameRules } from "./game/GameRules.js";
import { InputHandler } from "./game/InputHandler.js";
import { PhysicsEngine } from "./physics/PhysicsEngine.js";
import { Renderer, RenderState } from "./rendering/Renderer.js";
import { AudioManager } from "./audio/AudioManager.js";
import { EventBus, EventType } from "./core/EventBus.js";
import type { BallSunkEvent, ScratchEvent, CollisionEvent, EndDragEvent, StateChangeEvent } from "./core/EventBus.js";
import { BALL_COLORS, CANVAS_WIDTH, CANVAS_HEIGHT, GAME_CONFIG } from "./constants.js";

export interface GameCallbacks {
  onSunkCountChange?: (count: number) => void;
  onStateChange?: (state: GameState) => void;
  onReadyToShootChange?: (ready: boolean) => void;
  onWin?: () => void;
}

export class Game {
  private canvas: HTMLCanvasElement;
  private table: Table;
  private balls: Ball[] = [];
  private cueBall!: Ball;
  private stateManager: GameStateManager;
  private gameRules: GameRules;
  private inputHandler: InputHandler;
  private physicsEngine: PhysicsEngine;
  private renderer: Renderer;
  private audioManager: AudioManager;
  private eventBus: EventBus;
  private callbacks: GameCallbacks;

  private shotTimer: number = 0;
  private scratchTimer: number = 0;
  private animationId: number | null = null;
  private lastTime: number = 0;
  private unsubscribeHandlers: (() => void)[] = [];

  constructor(canvas: HTMLCanvasElement, callbacks: GameCallbacks = {}) {
    this.canvas = canvas;
    this.callbacks = callbacks;

    this.canvas.width = CANVAS_WIDTH;
    this.canvas.height = CANVAS_HEIGHT;

    this.eventBus = new EventBus();
    this.table = Table.createDefault();
    this.stateManager = new GameStateManager();
    this.gameRules = new GameRules();
    this.renderer = new Renderer(canvas);
    this.audioManager = new AudioManager();
    this.physicsEngine = new PhysicsEngine(this.table, this.eventBus);
    this.inputHandler = new InputHandler(canvas, this.eventBus);

    this.setupEventListeners();
    this.setupStateListeners();

    this.reset();
  }

  private setupEventListeners(): void {
    this.unsubscribeHandlers.push(
      this.eventBus.on(EventType.BALL_SUNK, (event) => {
        const { ball } = event as BallSunkEvent;
        const result = this.gameRules.handleBallSunk(ball);

        this.audioManager.playPocketSound();
        this.callbacks.onSunkCountChange?.(result.sunkNumbers.length);

        if (result.isWin) {
          this.callbacks.onWin?.();
        }
      })
    );

    this.unsubscribeHandlers.push(
      this.eventBus.on(EventType.SCRATCH, (event) => {
        const { cueBall } = event as ScratchEvent;
        this.gameRules.handleScratch(cueBall);
      })
    );

    this.unsubscribeHandlers.push(
      this.eventBus.on(EventType.COLLISION, (event) => {
        const { impactForce } = event as CollisionEvent;
        this.audioManager.playCollisionSound(impactForce);
      })
    );

    this.unsubscribeHandlers.push(
      this.eventBus.on(EventType.CUSHION_HIT, () => {
        this.audioManager.playCushionSound();
      })
    );

    this.unsubscribeHandlers.push(
      this.eventBus.on(EventType.START_DRAG, () => {
        this.audioManager.ensureContext();
      })
    );

    this.unsubscribeHandlers.push(
      this.eventBus.on(EventType.END_DRAG, (event) => {
        const { power, direction } = event as EndDragEvent;
        
        if (power > 0.2 && this.cueBall.inPlay) {
          this.physicsEngine.applyShot(this.cueBall, power, direction);
          this.stateManager.startShooting();
          this.shotTimer = GAME_CONFIG.shotAnimationDuration;
          this.audioManager.playCueHitSound(power);
        }
      })
    );
  }

  private setupStateListeners(): void {
    this.stateManager.onStateTransition((from, to) => {
      this.eventBus.emit({
        type: EventType.STATE_CHANGE,
        from,
        to
      } as StateChangeEvent);

      this.callbacks.onStateChange?.(to);
      this.callbacks.onReadyToShootChange?.(to === GameState.AIMING);
    });
  }

  reset(): void {
    this.balls = [];
    this.gameRules.reset();
    this.shotTimer = 0;
    this.scratchTimer = 0;
    this.inputHandler.reset();
    this.stateManager.reset();

    this.cueBall = Ball.createCueBall(this.table.headSpot.x, this.table.headSpot.y);
    this.balls.push(this.cueBall);

    const rackPositions = this.table.generateRackPositions();
    for (let i = 1; i <= 7; i++) {
      const pos = rackPositions[i - 1];
      const ball = Ball.createNumberedBall(i, BALL_COLORS[i], pos.x, pos.y);
      this.balls.push(ball);
    }

    this.updateUI();
  }

  start(): void {
    if (this.animationId !== null) return;
    this.lastTime = performance.now();
    this.loop();
  }

  stop(): void {
    if (this.animationId !== null) {
      cancelAnimationFrame(this.animationId);
      this.animationId = null;
    }
  }

  private loop = (): void => {
    const currentTime = performance.now();
    const dt = currentTime - this.lastTime;
    this.lastTime = currentTime;

    this.update(dt);
    this.render();

    this.animationId = requestAnimationFrame(this.loop);
  };

  private update(dt: number): void {
    if (this.stateManager.isAiming && this.cueBall.inPlay) {
      this.inputHandler.updateAimFromCueBall(this.cueBall.position);
    }

    if (this.stateManager.isShooting) {
      this.shotTimer -= dt;
      if (this.shotTimer <= 0) {
        this.stateManager.startBallsMoving();
      }
    }

    if (this.stateManager.isShooting || this.stateManager.isBallsMoving) {
      this.physicsEngine.simulate(this.balls);

      if (this.scratchTimer > 0) {
        this.scratchTimer -= dt;
        if (this.scratchTimer <= 0) {
          this.respawnCueBall();
        }
      }

      if (!this.physicsEngine.areBallsMoving(this.balls) &&
          this.scratchTimer <= 0 &&
          !this.stateManager.isAiming) {
        this.stateManager.returnToAiming();
      }
    }

    this.gameRules.updatePocketFlashes(dt);
    this.updateUI();
  }

  private respawnCueBall(): void {
    const isOccupied = (x: number, y: number): boolean => {
      return this.physicsEngine.isPositionOccupied(x, y, this.cueBall.radius, this.balls, this.cueBall);
    };

    const newPos = this.table.findRespawnPosition(
      this.table.headSpot,
      this.cueBall.radius,
      isOccupied
    );

    if (newPos) {
      this.cueBall.respawn(newPos.x, newPos.y);
    } else {
      this.cueBall.respawn(this.table.headSpot.x, this.table.headSpot.y);
    }
  }

  private render(): void {
    const state: RenderState = {
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

    this.renderer.render(state);
  }

  private updateUI(): void {
    this.callbacks.onSunkCountChange?.(this.gameRules.sunkCount);
    this.callbacks.onStateChange?.(this.stateManager.state);
    this.callbacks.onReadyToShootChange?.(
      this.stateManager.readyToShoot && this.cueBall.inPlay
    );
  }

  getState(): { sunkCount: number; gameState: GameState; readyToShoot: boolean } {
    return {
      sunkCount: this.gameRules.sunkCount,
      gameState: this.stateManager.state,
      readyToShoot: this.stateManager.readyToShoot && this.cueBall.inPlay
    };
  }

  destroy(): void {
    this.stop();
    this.inputHandler.destroy();
    this.audioManager.destroy();
    this.unsubscribeHandlers.forEach(unsubscribe => unsubscribe());
    this.eventBus.clear();
  }
}
