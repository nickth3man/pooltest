import { GameState } from "./GameState.js";

/** Callbacks for UI updates - all optional for flexibility */
export interface GameCallbacks {
  onSunkCountChange?: (count: number) => void;
  onStateChange?: (state: GameState) => void;
  onReadyToShootChange?: (ready: boolean) => void;
  onWin?: () => void;
}

/** Snapshot of game status for UI presentation */
export interface GameStatusSnapshot {
  sunkCount: number;
  gameState: GameState;
  readyToShoot: boolean;
}

/**
 * GamePresenter - Bridges game state to UI callbacks
 * 
 * Follows the Presenter pattern: transforms internal game state
 * into UI-friendly updates. Decouples game logic from DOM manipulation.
 */
export class GamePresenter {
  constructor(private callbacks: GameCallbacks = {}) {}

  /** Present full snapshot - called each frame */
  present(snapshot: GameStatusSnapshot): void {
    this.callbacks.onSunkCountChange?.(snapshot.sunkCount);
    this.callbacks.onStateChange?.(snapshot.gameState);
    this.callbacks.onReadyToShootChange?.(snapshot.readyToShoot);
  }

  presentState(state: GameState, readyToShoot: boolean): void {
    this.callbacks.onStateChange?.(state);
    this.callbacks.onReadyToShootChange?.(readyToShoot);
  }

  presentSunkCount(count: number): void {
    this.callbacks.onSunkCountChange?.(count);
  }

  presentWin(): void {
    this.callbacks.onWin?.();
  }
}
