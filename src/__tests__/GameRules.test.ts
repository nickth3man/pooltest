import { describe, it, expect, beforeEach } from 'vitest';
import { GameRules } from '../game/GameRules.js';
import { Ball } from '../models/Ball.js';

describe('GameRules', () => {
  let gameRules: GameRules;

  beforeEach(() => {
    gameRules = new GameRules();
  });

  it('should handle ball sunk', () => {
    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    
    const result = gameRules.handleBallSunk(ball);

    expect(result.sunkNumbers.length).toBe(1);
    expect(result.sunkNumbers).toContain(1);
    expect(result.shouldSpawnCueBall).toBe(false);
    expect(result.isWin).toBe(false);
    expect(result.pocketFlashes).toHaveLength(1);
  });

  it('should detect win condition when all balls sunk', () => {
    for (let i = 1; i <= 7; i++) {
      const ball = Ball.createNumberedBall(i, '#ff0000', 100, 100);
      gameRules.handleBallSunk(ball);
    }

    expect(gameRules.checkWinCondition()).toBe(true);
    expect(gameRules.sunkCount).toBe(7);
  });

  it('should not count duplicate sink events toward win condition', () => {
    for (let i = 0; i < 6; i++) {
      const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
      gameRules.handleBallSunk(ball);
    }

    expect(gameRules.sunkCount).toBe(1);
    expect(gameRules.checkWinCondition()).toBe(false);
  });

  it('should handle scratch', () => {
    const cueBall = Ball.createCueBall(100, 100);
    
    const result = gameRules.handleScratch(cueBall);

    expect(result.shouldSpawnCueBall).toBe(true);
    expect(result.isWin).toBe(false);
    expect(result.pocketFlashes).toHaveLength(1);
  });

  it('should reset state', () => {
    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    gameRules.handleBallSunk(ball);
    
    gameRules.reset();

    expect(gameRules.sunkCount).toBe(0);
    expect(gameRules.getSunkNumbers()).toHaveLength(0);
    expect(gameRules.getPocketFlashes()).toHaveLength(0);
  });

  it('should update pocket flashes', () => {
    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    gameRules.handleBallSunk(ball);

    expect(gameRules.getPocketFlashes()).toHaveLength(1);

    gameRules.updatePocketFlashes(250);

    expect(gameRules.getPocketFlashes()).toHaveLength(0);
  });
});
