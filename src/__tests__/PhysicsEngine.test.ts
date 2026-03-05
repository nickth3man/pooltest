import { describe, it, expect, vi } from 'vitest';
import { PhysicsEngine } from '../physics/PhysicsEngine.js';
import { Table } from '../models/Table.js';
import { Ball } from '../models/Ball.js';
import { EventBus, EventType } from '../core/EventBus.js';
import type { CollisionEvent } from '../core/EventBus.js';

describe('PhysicsEngine', () => {
  it('should detect when balls are moving', () => {
    const eventBus = new EventBus();
    const table = Table.createDefault();
    const physics = new PhysicsEngine(table, eventBus);
    
    const ball = Ball.createCueBall(100, 100);
    ball.vx = 5;
    
    expect(physics.areBallsMoving([ball])).toBe(true);
    
    ball.stop();
    expect(physics.areBallsMoving([ball])).toBe(false);
  });

  it('should emit collision event', () => {
    const eventBus = new EventBus();
    const collisionHandler = vi.fn();
    eventBus.on(EventType.COLLISION, collisionHandler);
    
    const table = Table.createDefault();
    const physics = new PhysicsEngine(table, eventBus);
    
    const ballA = Ball.createCueBall(100, 100);
    const ballB = Ball.createNumberedBall(1, '#ff0000', 111, 100);
    
    ballA.vx = 10;
    ballB.vx = -10;
    
    physics.simulate([ballA, ballB]);

    expect(collisionHandler).toHaveBeenCalled();
    const event = collisionHandler.mock.calls[0][0] as CollisionEvent;
    expect(event.type).toBe(EventType.COLLISION);
    expect(event.impactForce).toBeGreaterThan(0);
  });

  it('should check position occupancy', () => {
    const eventBus = new EventBus();
    const table = Table.createDefault();
    const physics = new PhysicsEngine(table, eventBus);
    
    const ball = Ball.createCueBall(100, 100);
    
    expect(physics.isPositionOccupied(100, 100, 10, [ball])).toBe(true);
    expect(physics.isPositionOccupied(500, 500, 10, [ball])).toBe(false);
  });

  it('should apply shot to cue ball', () => {
    const eventBus = new EventBus();
    const table = Table.createDefault();
    const physics = new PhysicsEngine(table, eventBus);
    
    const cueBall = Ball.createCueBall(100, 100);
    
    physics.applyShot(cueBall, 10, { x: 1, y: 0 });
    
    expect(cueBall.vx).toBe(10);
    expect(cueBall.vy).toBe(0);
  });

  it('should calculate total energy', () => {
    const eventBus = new EventBus();
    const table = Table.createDefault();
    const physics = new PhysicsEngine(table, eventBus);
    
    const ball1 = Ball.createCueBall(0, 0);
    ball1.vx = 3;
    ball1.vy = 4;
    
    const ball2 = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    ball2.vx = 6;
    ball2.vy = 8;
    
    const energy = physics.getTotalEnergy([ball1, ball2]);
    
    expect(energy).toBe(0.5 * (25 + 100));
  });
});
