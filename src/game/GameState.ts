/**
 * GameState manager
 * Handles state transitions and state-dependent logic
 */

import { GameState as GameStateEnum, type GameEventCallback } from "../types.js";

export { GameStateEnum as GameState };

export class GameStateManager {
  private currentState: GameStateEnum;
  private previousState: GameStateEnum | null = null;
  private callbacks: Map<GameStateEnum, Set<GameEventCallback>> = new Map();
  private transitionCallbacks: Set<(from: GameStateEnum, to: GameStateEnum) => void> = new Set();

  constructor(initialState: GameStateEnum = GameStateEnum.AIMING) {
    this.currentState = initialState;
  }

  /** Get current state */
  get state(): GameStateEnum {
    return this.currentState;
  }

  /** Get previous state */
  get previousStateValue(): GameStateEnum | null {
    return this.previousState;
  }

  /** Check if currently aiming */
  get isAiming(): boolean {
    return this.currentState === GameStateEnum.AIMING;
  }

  /** Check if currently shooting */
  get isShooting(): boolean {
    return this.currentState === GameStateEnum.SHOOTING;
  }

  /** Check if balls are moving */
  get isBallsMoving(): boolean {
    return this.currentState === GameStateEnum.BALLS_MOVING;
  }

  /** Check if ready to shoot (in aiming state) */
  get readyToShoot(): boolean {
    return this.isAiming;
  }

  /** Transition to a new state */
  transitionTo(newState: GameStateEnum): void {
    if (newState === this.currentState) return;

    const oldState = this.currentState;
    this.previousState = oldState;
    this.currentState = newState;

    // Notify transition callbacks
    this.transitionCallbacks.forEach(callback => callback(oldState, newState));

    // Notify state-specific callbacks
    const stateCallbacks = this.callbacks.get(newState);
    if (stateCallbacks) {
      stateCallbacks.forEach(callback => callback());
    }
  }

  /** Register a callback for a specific state */
  onEnterState(state: GameStateEnum, callback: GameEventCallback): () => void {
    if (!this.callbacks.has(state)) {
      this.callbacks.set(state, new Set());
    }
    this.callbacks.get(state)!.add(callback);

    // Return unsubscribe function
    return () => {
      this.callbacks.get(state)?.delete(callback);
    };
  }

  /** Register a callback for any state transition */
  onStateTransition(callback: (from: GameStateEnum, to: GameStateEnum) => void): () => void {
    this.transitionCallbacks.add(callback);
    return () => {
      this.transitionCallbacks.delete(callback);
    };
  }

  /** Start shooting animation */
  startShooting(): void {
    this.transitionTo(GameStateEnum.SHOOTING);
  }

  /** Start balls moving state */
  startBallsMoving(): void {
    this.transitionTo(GameStateEnum.BALLS_MOVING);
  }

  /** Return to aiming state */
  returnToAiming(): void {
    this.transitionTo(GameStateEnum.AIMING);
  }

  /** Reset to initial state */
  reset(): void {
    this.previousState = null;
    this.transitionTo(GameStateEnum.AIMING);
  }

  /** Check if can transition to a specific state */
  canTransitionTo(state: GameStateEnum): boolean {
    // Define valid transitions
    const validTransitions: Record<GameStateEnum, GameStateEnum[]> = {
      [GameStateEnum.AIMING]: [GameStateEnum.SHOOTING],
      [GameStateEnum.SHOOTING]: [GameStateEnum.BALLS_MOVING, GameStateEnum.AIMING],
      [GameStateEnum.BALLS_MOVING]: [GameStateEnum.AIMING]
    };

    return validTransitions[this.currentState]?.includes(state) ?? false;
  }

  /** Get state as string for display */
  toString(): string {
    return this.currentState;
  }
}
