import type { PocketFlash, RenderContext } from "../types.js";
import { BALL_COLORS, VISUAL, SHOT } from "../constants.js";

export class UIRenderer {
  constructor(private ctx: RenderContext) {}

  renderPowerMeter(pullDistance: number): void {
    if (pullDistance <= 0) return;

    const { x, y, width, height } = {
      x: VISUAL.powerMeterX,
      y: VISUAL.powerMeterY,
      width: VISUAL.powerMeterWidth,
      height: VISUAL.powerMeterHeight
    };

    const power = (pullDistance * SHOT.powerScale) / SHOT.maxPower;
    const fillHeight = height * Math.min(power, 1);

    this.ctx.fillStyle = "rgba(12,12,12,0.6)";
    this.ctx.fillRect(x, y, width, height);

    const gradient = this.ctx.createLinearGradient(0, y + height, 0, y);
    gradient.addColorStop(0, "#73c95b");
    gradient.addColorStop(0.6, "#f1c24f");
    gradient.addColorStop(1, "#cf4f4f");
    this.ctx.fillStyle = gradient;
    this.ctx.fillRect(x, y + height - fillHeight, width, fillHeight);

    this.ctx.strokeStyle = "rgba(255,255,255,0.6)";
    this.ctx.lineWidth = 1.2;
    this.ctx.strokeRect(x, y, width, height);
  }

  renderSunkRow(sunkNumbers: number[]): void {
    const startX = VISUAL.sunkRowStartX;
    const y = VISUAL.sunkRowY;
    const miniR = VISUAL.miniBallRadius;

    for (let i = 0; i < sunkNumbers.length; i++) {
      const n = sunkNumbers[i];
      const x = startX + i * 20;
      const color = BALL_COLORS[n];

      this.ctx.fillStyle = color;
      this.ctx.beginPath();
      this.ctx.arc(x, y, miniR, 0, Math.PI * 2);
      this.ctx.fill();

      this.ctx.fillStyle = "#fff";
      this.ctx.beginPath();
      this.ctx.arc(x, y, miniR * 0.45, 0, Math.PI * 2);
      this.ctx.fill();

      this.ctx.fillStyle = "#111";
      this.ctx.font = "bold 7px Trebuchet MS";
      this.ctx.textAlign = "center";
      this.ctx.textBaseline = "middle";
      this.ctx.fillText(String(n), x, y + 0.2);
    }
  }

  renderPocketFlashes(pocketFlashes: PocketFlash[]): void {
    const maxLife = 220;

    for (const flash of pocketFlashes) {
      const t = flash.life / maxLife;
      this.ctx.strokeStyle = `rgba(255,245,186,${0.55 * t})`;
      this.ctx.lineWidth = 2.5;
      this.ctx.beginPath();
      this.ctx.arc(flash.x, flash.y, 13 + (1 - t) * 16, 0, Math.PI * 2);
      this.ctx.stroke();
    }
  }
}
