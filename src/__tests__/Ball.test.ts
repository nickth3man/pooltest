import { describe, it, expect } from 'vitest';
import { Ball } from '../models/Ball.js';

describe('Ball', () => {
  it('should create cue ball', () => {
    const ball = Ball.createCueBall(100, 200);

    expect(ball.isCue).toBe(true);
    expect(ball.number).toBe(0);
    expect(ball.color).toBe('#ffffff');
    expect(ball.x).toBe(100);
    expect(ball.y).toBe(200);
    expect(ball.inPlay).toBe(true);
  });

  it('should create numbered ball', () => {
    const ball = Ball.createNumberedBall(5, '#ff0000', 150, 250);

    expect(ball.isCue).toBe(false);
    expect(ball.number).toBe(5);
    expect(ball.color).toBe('#ff0000');
    expect(ball.x).toBe(150);
    expect(ball.y).toBe(250);
  });

  it('should calculate speed', () => {
    const ball = Ball.createCueBall(0, 0);
    ball.vx = 3;
    ball.vy = 4;

    expect(ball.speed).toBe(5);
    expect(ball.speedSq).toBe(25);
  });

  it('should apply impulse', () => {
    const ball = Ball.createCueBall(0, 0);
    ball.applyImpulse({ x: 5, y: 10 });

    expect(ball.vx).toBe(5);
    expect(ball.vy).toBe(10);
  });

  it('should stop ball', () => {
    const ball = Ball.createCueBall(0, 0);
    ball.vx = 10;
    ball.vy = 20;
    
    ball.stop();

    expect(ball.vx).toBe(0);
    expect(ball.vy).toBe(0);
    expect(ball.isMoving).toBe(false);
  });

  it('should sink ball', () => {
    const ball = Ball.createCueBall(0, 0);
    ball.vx = 10;
    ball.vy = 20;

    ball.sink();

    expect(ball.inPlay).toBe(false);
    expect(ball.vx).toBe(0);
    expect(ball.vy).toBe(0);
  });

  it('should respawn ball', () => {
    const ball = Ball.createCueBall(0, 0);
    ball.sink();

    ball.respawn(100, 200);

    expect(ball.inPlay).toBe(true);
    expect(ball.x).toBe(100);
    expect(ball.y).toBe(200);
  });

  it('should detect overlap', () => {
    const ball1 = Ball.createCueBall(0, 0);
    const ball2 = Ball.createNumberedBall(1, '#ff0000', 20, 0);

    expect(ball1.overlaps(ball2)).toBe(true);

    const ball3 = Ball.createNumberedBall(2, '#00ff00', 100, 100);
    expect(ball1.overlaps(ball3)).toBe(false);
  });

  it('should clone ball', () => {
    const ball = Ball.createCueBall(50, 60);
    ball.vx = 5;
    ball.vy = 10;

    const cloned = ball.clone();

    expect(cloned.x).toBe(ball.x);
    expect(cloned.y).toBe(ball.y);
    expect(cloned.vx).toBe(ball.vx);
    expect(cloned.vy).toBe(ball.vy);
    expect(cloned).not.toBe(ball);
  });
});
