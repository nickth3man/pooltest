/**
 * InputHandler class
 * Manages mouse/touch input for aiming and shooting
 */

import type { Vector2, DragState } from "../types.js";
import { Vec2, Vector2Utils } from "../models/Vector2.js";
import { SHOT, VISUAL } from "../constants.js";
import { EventBus, EventType } from "../core/EventBus.js";

interface AimStateInternal {
  direction: Vec2;
  lockedDirection: Vec2;
}

export class InputHandler {
  private readonly touchEventOptions: AddEventListenerOptions = {
    passive: false,
  };
  private mouse: Vec2;
  private drag: DragState;
  private aim: AimStateInternal;
  private eventBus: EventBus;
  private canvas: HTMLCanvasElement;
  private isEnabled: boolean = true;

  constructor(canvas: HTMLCanvasElement, eventBus: EventBus) {
    this.canvas = canvas;
    this.eventBus = eventBus;
    this.mouse = new Vec2(0, 0);
    this.drag = {
      isDragging: false,
      startX: 0,
      startY: 0,
      pullDistance: 0,
    };
    this.aim = {
      direction: new Vec2(1, 0),
      lockedDirection: new Vec2(1, 0),
    };

    this.bindEvents();
  }

  private bindEvents(): void {
    // Mouse events for desktop
    this.canvas.addEventListener("mousemove", this.handleMouseMove);
    this.canvas.addEventListener("mousedown", this.handleMouseDown);
    this.canvas.addEventListener("mouseup", this.handleMouseUp);
    this.canvas.addEventListener("mouseleave", this.handleMouseUp);

    // Touch events for mobile devices
    this.canvas.addEventListener(
      "touchmove",
      this.handleTouchMove,
      this.touchEventOptions,
    );
    this.canvas.addEventListener(
      "touchstart",
      this.handleTouchStart,
      this.touchEventOptions,
    );
    this.canvas.addEventListener(
      "touchend",
      this.handleMouseUp,
      this.touchEventOptions,
    );
    this.canvas.addEventListener(
      "touchcancel",
      this.handleMouseUp,
      this.touchEventOptions,
    );
  }

  private getMousePosition(clientX: number, clientY: number): Vec2 {
    const rect = this.canvas.getBoundingClientRect();
    // Handle CSS scaling: canvas internal size may differ from display size
    const scaleX = this.canvas.width / rect.width;
    const scaleY = this.canvas.height / rect.height;

    return new Vec2(
      (clientX - rect.left) * scaleX,
      (clientY - rect.top) * scaleY,
    );
  }

  private handleMouseMove = (event: MouseEvent): void => {
    if (!this.isEnabled) return;

    this.mouse = this.getMousePosition(event.clientX, event.clientY);

    if (this.drag.isDragging) {
      this.updateDrag();
    }
  };

  private handleTouchMove = (event: TouchEvent): void => {
    if (!this.isEnabled) return;
    event.preventDefault();

    const touch = event.touches[0];
    this.mouse = this.getMousePosition(touch.clientX, touch.clientY);

    if (this.drag.isDragging) {
      this.updateDrag();
    }
  };

  private handleMouseDown = (event: MouseEvent): void => {
    if (!this.isEnabled) return;

    this.mouse = this.getMousePosition(event.clientX, event.clientY);
    this.startDrag();
  };

  private handleTouchStart = (event: TouchEvent): void => {
    if (!this.isEnabled) return;
    event.preventDefault();

    const touch = event.touches[0];
    this.mouse = this.getMousePosition(touch.clientX, touch.clientY);
    this.startDrag();
  };

  private startDrag(): void {
    this.drag.isDragging = true;
    this.drag.startX = this.mouse.x;
    this.drag.startY = this.mouse.y;
    this.drag.pullDistance = 0;
    this.aim.lockedDirection = this.aim.direction.clone(); // Lock aim direction when drag starts

    // Notify game that drag has begun (used to unlock audio context)
    this.eventBus.emit({
      type: EventType.START_DRAG,
      position: this.mouse.toObject(),
    });
  }

  private handleMouseUp = (): void => {
    if (!this.drag.isDragging) return;

    const power = this.calculatePower();
    this.drag.isDragging = false;
    this.drag.pullDistance = 0;

    if (power > 0.2) {
      this.eventBus.emit({
        type: EventType.END_DRAG,
        power,
        direction: this.aim.lockedDirection.toObject(),
      });
    }
  };

  updateAimFromCueBall(cueBallPosition: Vector2): void {
    const direction = new Vec2(
      this.mouse.x - cueBallPosition.x,
      this.mouse.y - cueBallPosition.y,
    ).normalize({ x: 1, y: 0 });

    this.aim.direction = direction;
    this.eventBus.emit({
      type: EventType.AIM_CHANGE,
      direction: direction.toObject(),
    });
  }

  private updateDrag(): void {
    const dragDelta = new Vec2(
      this.mouse.x - this.drag.startX,
      this.mouse.y - this.drag.startY,
    );

    const projection = -dragDelta.dot(this.aim.lockedDirection);
    this.drag.pullDistance = Vector2Utils.clamp(
      projection,
      0,
      VISUAL.maxPullDistance,
    );

    this.eventBus.emit({
      type: EventType.DRAG,
      pullDistance: this.drag.pullDistance,
    });
  }

  private calculatePower(): number {
    return Vector2Utils.clamp(
      this.drag.pullDistance * SHOT.powerScale,
      0,
      SHOT.maxPower,
    );
  }

  get aimDirection(): Vec2 {
    return this.aim.direction;
  }

  get lockedAimDirection(): Vec2 {
    return this.aim.lockedDirection;
  }

  get mousePosition(): Vec2 {
    return this.mouse;
  }

  get isDragging(): boolean {
    return this.drag.isDragging;
  }

  get pullDistance(): number {
    return this.drag.pullDistance;
  }

  enable(): void {
    this.isEnabled = true;
  }

  disable(): void {
    this.isEnabled = false;
    if (this.drag.isDragging) {
      this.handleMouseUp();
    }
  }

  reset(): void {
    this.drag.isDragging = false;
    this.drag.pullDistance = 0;
    this.aim.direction = new Vec2(1, 0);
    this.aim.lockedDirection = new Vec2(1, 0);
  }

  destroy(): void {
    this.canvas.removeEventListener("mousemove", this.handleMouseMove);
    this.canvas.removeEventListener("mousedown", this.handleMouseDown);
    this.canvas.removeEventListener("mouseup", this.handleMouseUp);
    this.canvas.removeEventListener("mouseleave", this.handleMouseUp);
    this.canvas.removeEventListener(
      "touchmove",
      this.handleTouchMove,
      this.touchEventOptions,
    );
    this.canvas.removeEventListener(
      "touchstart",
      this.handleTouchStart,
      this.touchEventOptions,
    );
    this.canvas.removeEventListener(
      "touchend",
      this.handleMouseUp,
      this.touchEventOptions,
    );
    this.canvas.removeEventListener(
      "touchcancel",
      this.handleMouseUp,
      this.touchEventOptions,
    );
  }
}
