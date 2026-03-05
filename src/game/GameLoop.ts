/**
 * GameLoop - Manages the requestAnimationFrame game loop
 * 
 * Uses high-resolution timestamp for smooth, frame-rate independent updates.
 * dt (delta time) is in milliseconds since last frame.
 */
export class GameLoop {
  private animationId: number | null = null;
  private lastTime: number = 0;

  constructor(private onTick: (dt: number) => void) {}

  private frame = (): void => {
    const currentTime = performance.now();
    const dt = currentTime - this.lastTime;  // Milliseconds since last frame
    this.lastTime = currentTime;

    this.onTick(dt);  // Delegate to game update logic
    this.animationId = requestAnimationFrame(this.frame);  // Next frame
  };

  start(): void {
    if (this.animationId !== null) {
      return;
    }

    this.lastTime = performance.now();
    this.animationId = requestAnimationFrame(this.frame);
  }

  stop(): void {
    if (this.animationId === null) {
      return;
    }

    cancelAnimationFrame(this.animationId);
    this.animationId = null;
  }

  get isRunning(): boolean {
    return this.animationId !== null;
  }
}
