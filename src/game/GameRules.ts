/**
 * GameRules - Extracted game rule logic
 * Handles scoring, win conditions, and game-specific rules
 */

import type { Ball } from "../models/Ball.js";
import type { PocketFlash } from "../types.js";
import { GAME_CONFIG } from "../constants.js";

export interface RuleResult {
  pocketFlashes: PocketFlash[];
  sunkNumbers: number[];
  shouldSpawnCueBall: boolean;
  isWin: boolean;
  message?: string;
}

export class GameRules {
  private sunkBalls: Set<number> = new Set();
  private pocketFlashes: PocketFlash[] = [];

  /**
   * Handle a ball being sunk
   * @param ball - The ball that was sunk
   * @returns Rule result with effects and state changes
   */
  handleBallSunk(ball: Ball): RuleResult {
    const pocketFlash: PocketFlash = {
      x: ball.x,
      y: ball.y,
      life: GAME_CONFIG.pocketFlashDuration
    };

    this.pocketFlashes.push(pocketFlash);
    if (!ball.isCue && ball.number >= 1 && ball.number <= 7) {
      this.sunkBalls.add(ball.number);
    }

    const isWin = this.checkWinCondition();

    return {
      pocketFlashes: [...this.pocketFlashes],
      sunkNumbers: Array.from(this.sunkBalls),
      shouldSpawnCueBall: false,
      isWin,
      message: isWin ? "Congratulations! You won!" : undefined
    };
  }

  /**
   * Handle a scratch (cue ball sunk)
   * @param cueBall - The cue ball that was sunk
   * @returns Rule result with effects
   */
  handleScratch(cueBall: Ball): RuleResult {
    const pocketFlash: PocketFlash = {
      x: cueBall.x,
      y: cueBall.y,
      life: GAME_CONFIG.pocketFlashDuration
    };

    this.pocketFlashes.push(pocketFlash);

    return {
      pocketFlashes: [...this.pocketFlashes],
      sunkNumbers: Array.from(this.sunkBalls),
      shouldSpawnCueBall: true,
      isWin: false
    };
  }

  /**
   * Check if the player has won
   * @returns True if all numbered balls (1-7) are sunk
   */
  checkWinCondition(): boolean {
    if (this.sunkBalls.size !== 7) {
      return false;
    }

    for (let number = 1; number <= 7; number++) {
      if (!this.sunkBalls.has(number)) {
        return false;
      }
    }

    return true;
  }

  /**
   * Get current sunk ball count
   */
  get sunkCount(): number {
    return this.sunkBalls.size;
  }

  /**
   * Get sunk ball numbers
   */
  getSunkNumbers(): number[] {
    return Array.from(this.sunkBalls);
  }

  /**
   * Get current pocket flashes
   */
  getPocketFlashes(): PocketFlash[] {
    return [...this.pocketFlashes];
  }

  /**
   * Update pocket flashes (decrement life)
   * @param dt - Delta time in milliseconds
   */
  updatePocketFlashes(dt: number): void {
    for (let i = this.pocketFlashes.length - 1; i >= 0; i--) {
      this.pocketFlashes[i].life -= dt;
      if (this.pocketFlashes[i].life <= 0) {
        this.pocketFlashes.splice(i, 1);
      }
    }
  }

  /**
   * Reset game rules state
   */
  reset(): void {
    this.sunkBalls.clear();
    this.pocketFlashes = [];
  }
}
