export class GameLoop {
  private animationId: number | null = null;
  private lastTime: number = 0;

  constructor(private onTick: (dt: number) => void) {}

  private frame = (): void => {
    const currentTime = performance.now();
    const dt = currentTime - this.lastTime;
    this.lastTime = currentTime;

    this.onTick(dt);
    this.animationId = requestAnimationFrame(this.frame);
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
