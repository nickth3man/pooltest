import { GameState } from "./GameState.js";

export interface GameCallbacks {
  onSunkCountChange?: (count: number) => void;
  onStateChange?: (state: GameState) => void;
  onReadyToShootChange?: (ready: boolean) => void;
  onWin?: () => void;
}

export interface GameStatusSnapshot {
  sunkCount: number;
  gameState: GameState;
  readyToShoot: boolean;
}

export class GamePresenter {
  constructor(private callbacks: GameCallbacks = {}) {}

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
