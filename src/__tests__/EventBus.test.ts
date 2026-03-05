import { describe, it, expect, vi } from 'vitest';
import { EventBus, EventType } from '../core/EventBus.js';
import type { BallSunkEvent, CollisionEvent } from '../core/EventBus.js';
import { Ball } from '../models/Ball.js';

describe('EventBus', () => {
  it('should subscribe and emit events', () => {
    const eventBus = new EventBus();
    const handler = vi.fn();

    eventBus.on(EventType.BALL_SUNK, handler);
    
    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    const event: BallSunkEvent = {
      type: EventType.BALL_SUNK,
      ball
    };
    
    eventBus.emit(event);

    expect(handler).toHaveBeenCalledWith(event);
  });

  it('should support multiple handlers for same event', () => {
    const eventBus = new EventBus();
    const handler1 = vi.fn();
    const handler2 = vi.fn();

    eventBus.on(EventType.COLLISION, handler1);
    eventBus.on(EventType.COLLISION, handler2);

    const ballA = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    const ballB = Ball.createNumberedBall(2, '#00ff00', 110, 100);
    const event: CollisionEvent = {
      type: EventType.COLLISION,
      ballA,
      ballB,
      impactForce: 5
    };

    eventBus.emit(event);

    expect(handler1).toHaveBeenCalledWith(event);
    expect(handler2).toHaveBeenCalledWith(event);
  });

  it('should unsubscribe handlers', () => {
    const eventBus = new EventBus();
    const handler = vi.fn();

    const unsubscribe = eventBus.on(EventType.BALL_SUNK, handler);
    unsubscribe();

    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    eventBus.emit({
      type: EventType.BALL_SUNK,
      ball
    } as BallSunkEvent);

    expect(handler).not.toHaveBeenCalled();
  });

  it('should handle once subscriptions', () => {
    const eventBus = new EventBus();
    const handler = vi.fn();

    eventBus.once(EventType.BALL_SUNK, handler);

    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    const event: BallSunkEvent = {
      type: EventType.BALL_SUNK,
      ball
    };

    eventBus.emit(event);
    eventBus.emit(event);

    expect(handler).toHaveBeenCalledTimes(1);
  });

  it('should track handler count', () => {
    const eventBus = new EventBus();
    
    expect(eventBus.handlerCount(EventType.BALL_SUNK)).toBe(0);
    
    const unsubscribe = eventBus.on(EventType.BALL_SUNK, () => {});
    expect(eventBus.handlerCount(EventType.BALL_SUNK)).toBe(1);
    
    unsubscribe();
    expect(eventBus.handlerCount(EventType.BALL_SUNK)).toBe(0);
  });

  it('should clear all handlers', () => {
    const eventBus = new EventBus();
    const handler = vi.fn();

    eventBus.on(EventType.BALL_SUNK, handler);
    eventBus.clear();

    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    eventBus.emit({
      type: EventType.BALL_SUNK,
      ball
    } as BallSunkEvent);

    expect(handler).not.toHaveBeenCalled();
  });

  it('should report handler errors to configured error hook', () => {
    const eventBus = new EventBus();
    const errorHook = vi.fn();
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    eventBus.setErrorHandler(errorHook);
    eventBus.on(EventType.BALL_SUNK, () => {
      throw new Error('handler failed');
    });

    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    const result = eventBus.emit({
      type: EventType.BALL_SUNK,
      ball
    });

    expect(result.delivered).toBe(0);
    expect(result.errors).toBe(1);
    expect(errorHook).toHaveBeenCalledTimes(1);
    expect(consoleSpy).toHaveBeenCalledTimes(1);

    consoleSpy.mockRestore();
  });

  it('should fail fast when handler error throwing is enabled', () => {
    const eventBus = new EventBus({ throwOnHandlerError: true });
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    eventBus.on(EventType.BALL_SUNK, () => {
      throw new Error('boom');
    });

    const ball = Ball.createNumberedBall(1, '#ff0000', 100, 100);
    expect(() => eventBus.emit({ type: EventType.BALL_SUNK, ball })).toThrow('boom');

    consoleSpy.mockRestore();
  });
});
