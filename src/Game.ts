import { CANVAS_HEIGHT, CANVAS_WIDTH } from "./constants.js";
import { EventType } from "./core/EventBus.js";
import { GameLoop } from "./game/GameLoop.js";
import { GamePresenter, type GameCallbacks } from "./game/GamePresenter.js";
import { GameSession } from "./game/GameSession.js";
import type { GameState } from "./game/GameState.js";

export type { GameCallbacks } from "./game/GamePresenter.js";

export class Game {
  private session: GameSession;
  private presenter: GamePresenter;
  private loop: GameLoop;
  private unsubscribeHandlers: Array<() => void> = [];

  constructor(canvas: HTMLCanvasElement, callbacks: GameCallbacks = {}) {
    canvas.width = CANVAS_WIDTH;
    canvas.height = CANVAS_HEIGHT;

    this.presenter = new GamePresenter(callbacks);
    this.session = new GameSession(canvas);
    this.loop = new GameLoop((dt) => {
      this.session.update(dt);
      this.session.render();
      this.presenter.present(this.session.getState());
    });

    this.setupEventListeners();
    this.setupStateListeners();
    this.presenter.present(this.session.getState());
  }

  reset(): void {
    this.session.reset();
    this.presenter.present(this.session.getState());
  }

  start(): void {
    this.loop.start();
  }

  stop(): void {
    this.loop.stop();
  }

  getState(): { sunkCount: number; gameState: GameState; readyToShoot: boolean } {
    return this.session.getState();
  }

  destroy(): void {
    this.stop();
    this.unsubscribeHandlers.forEach((unsubscribe) => unsubscribe());
    this.unsubscribeHandlers = [];
    this.session.destroy();
  }

  private setupEventListeners(): void {
    this.unsubscribeHandlers.push(
      this.session.eventBus.on(EventType.BALL_SUNK, ({ ball }) => {
        const result = this.session.handleBallSunk(ball);
        this.session.playPocketSound();
        this.presenter.presentSunkCount(result.sunkNumbers.length);

        if (result.isWin) {
          this.presenter.presentWin();
        }
      })
    );

    this.unsubscribeHandlers.push(
      this.session.eventBus.on(EventType.SCRATCH, ({ cueBall }) => {
        this.session.handleScratch(cueBall);
      })
    );

    this.unsubscribeHandlers.push(
      this.session.eventBus.on(EventType.COLLISION, ({ impactForce }) => {
        this.session.playCollisionSound(impactForce);
      })
    );

    this.unsubscribeHandlers.push(
      this.session.eventBus.on(EventType.CUSHION_HIT, () => {
        this.session.playCushionSound();
      })
    );

    this.unsubscribeHandlers.push(
      this.session.eventBus.on(EventType.START_DRAG, () => {
        this.session.ensureAudioContext();
      })
    );

    this.unsubscribeHandlers.push(
      this.session.eventBus.on(EventType.END_DRAG, ({ power, direction }) => {
        this.session.handleShot(power, direction);
      })
    );
  }

  private setupStateListeners(): void {
    this.unsubscribeHandlers.push(
      this.session.stateManager.onStateTransition((from, to) => {
        this.session.eventBus.emit({
          type: EventType.STATE_CHANGE,
          from,
          to
        });

        this.presenter.presentState(this.session.getGameState(), this.session.isReadyToShoot());
      })
    );
  }
}
