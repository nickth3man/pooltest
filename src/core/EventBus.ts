/**
 * EventBus - Central event system for decoupled communication
 * Replaces manual callback wiring with a publish-subscribe pattern
 */

import type { Ball } from "../models/Ball.js";
import type { Vector2 } from "../types.js";
import { GameState } from "../game/GameState.js";

/** Event types supported by the system */
export enum EventType {
  // Physics events
  BALL_SUNK = "ballSunk",
  SCRATCH = "scratch",
  COLLISION = "collision",
  CUSHION_HIT = "cushionHit",
  
  // Input events
  START_DRAG = "startDrag",
  DRAG = "drag",
  END_DRAG = "endDrag",
  AIM_CHANGE = "aimChange",
  
  // Game state events
  STATE_CHANGE = "stateChange",
  BALLS_STOPPED = "ballsStopped",
  
  // UI events
  SUNK_COUNT_CHANGE = "sunkCountChange",
  READY_TO_SHOOT_CHANGE = "readyToShootChange"
}

/** Event payloads */
export interface BallSunkEvent {
  type: EventType.BALL_SUNK;
  ball: Ball;
}

export interface ScratchEvent {
  type: EventType.SCRATCH;
  cueBall: Ball;
}

export interface CollisionEvent {
  type: EventType.COLLISION;
  ballA: Ball;
  ballB: Ball;
  impactForce: number;
}

export interface CushionHitEvent {
  type: EventType.CUSHION_HIT;
  ball: Ball;
}

export interface StartDragEvent {
  type: EventType.START_DRAG;
  position: Vector2;
}

export interface DragEvent {
  type: EventType.DRAG;
  pullDistance: number;
}

export interface EndDragEvent {
  type: EventType.END_DRAG;
  power: number;
  direction: Vector2;
}

export interface AimChangeEvent {
  type: EventType.AIM_CHANGE;
  direction: Vector2;
}

export interface StateChangeEvent {
  type: EventType.STATE_CHANGE;
  from: GameState;
  to: GameState;
}

export interface BallsStoppedEvent {
  type: EventType.BALLS_STOPPED;
}

export interface SunkCountChangeEvent {
  type: EventType.SUNK_COUNT_CHANGE;
  count: number;
}

export interface ReadyToShootChangeEvent {
  type: EventType.READY_TO_SHOOT_CHANGE;
  ready: boolean;
}

export type GameEvent =
  | BallSunkEvent
  | ScratchEvent
  | CollisionEvent
  | CushionHitEvent
  | StartDragEvent
  | DragEvent
  | EndDragEvent
  | AimChangeEvent
  | StateChangeEvent
  | BallsStoppedEvent
  | SunkCountChangeEvent
  | ReadyToShootChangeEvent;

export type EventHandler = (event: GameEvent) => void;

/**
 * EventBus - Central event management system
 * Provides type-safe event emission and subscription
 */
export class EventBus {
  private handlers: Map<EventType, Set<EventHandler>> = new Map();

  /**
   * Subscribe to an event type
   * @param type - Event type to subscribe to
   * @param handler - Handler function to call when event is emitted
   * @returns Unsubscribe function
   */
  on(type: EventType, handler: EventHandler): () => void {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, new Set());
    }
    
    const handlers = this.handlers.get(type)!;
    handlers.add(handler);
    
    // Return unsubscribe function
    return () => {
      handlers.delete(handler);
      if (handlers.size === 0) {
        this.handlers.delete(type);
      }
    };
  }

  /**
   * Subscribe to an event type (once - auto-unsubscribes after first emission)
   * @param type - Event type to subscribe to
   * @param handler - Handler function to call when event is emitted
   * @returns Unsubscribe function
   */
  once(type: EventType, handler: EventHandler): () => void {
    const wrappedHandler = (event: GameEvent) => {
      handler(event);
      unsubscribe();
    };
    
    const unsubscribe = this.on(type, wrappedHandler);
    return unsubscribe;
  }

  /**
   * Emit an event to all subscribers
   * @param event - Event to emit
   */
  emit(event: GameEvent): void {
    const handlers = this.handlers.get(event.type);
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(event);
        } catch (error) {
          console.error(`Error in event handler for ${event.type}:`, error);
        }
      });
    }
  }

  /**
   * Remove all handlers for a specific event type
   * @param type - Event type to clear
   */
  off(type: EventType): void {
    this.handlers.delete(type);
  }

  /**
   * Remove all handlers for all event types
   */
  clear(): void {
    this.handlers.clear();
  }

  /**
   * Get the number of handlers for a specific event type
   * @param type - Event type to check
   * @returns Number of handlers
   */
  handlerCount(type: EventType): number {
    return this.handlers.get(type)?.size ?? 0;
  }

  /**
   * Check if there are any handlers for a specific event type
   * @param type - Event type to check
   * @returns True if handlers exist
   */
  hasHandlers(type: EventType): boolean {
    return (this.handlers.get(type)?.size ?? 0) > 0;
  }
}
