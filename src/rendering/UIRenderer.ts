import type { PocketFlash, RenderContext } from "../types.js";
import { BALL_COLORS, VISUAL, SHOT, TIMING } from "../constants.js";

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

    this.drawRoundedRect(x - 3, y - 3, width + 6, height + 6, 5);
    this.ctx.fillStyle = "rgba(0,0,0,0.28)";
    this.ctx.fill();

    this.drawRoundedRect(x, y, width, height, 4);
    this.ctx.fillStyle = "rgba(12,12,12,0.72)";
    this.ctx.fill();

    const gradient = this.ctx.createLinearGradient(0, y + height, 0, y);
    gradient.addColorStop(0, "#53bf68");
    gradient.addColorStop(0.6, "#e9bd53");
    gradient.addColorStop(1, "#cf4f4f");
    this.ctx.fillStyle = gradient;
    this.ctx.fillRect(x + 1, y + height - fillHeight, width - 2, fillHeight);

    this.ctx.fillStyle = "rgba(255,255,255,0.14)";
    this.ctx.fillRect(x + 1, y + 1, width - 2, 3);

    this.ctx.strokeStyle = "rgba(255,255,255,0.7)";
    this.ctx.lineWidth = 1.2;
    this.drawRoundedRect(x, y, width, height, 4);
    this.ctx.stroke();
  }

  renderSunkRow(sunkNumbers: number[]): void {
    const startX = VISUAL.sunkRowStartX;
    const y = VISUAL.sunkRowY;
    const miniR = VISUAL.miniBallRadius;

    for (let i = 0; i < sunkNumbers.length; i++) {
      const n = sunkNumbers[i];
      const x = startX + i * 20;
      const color = BALL_COLORS[n];

      this.ctx.fillStyle = "rgba(0,0,0,0.26)";
      this.ctx.beginPath();
      this.ctx.arc(x + 0.9, y + 1.6, miniR * 0.98, 0, Math.PI * 2);
      this.ctx.fill();

      this.ctx.fillStyle = color;
      this.ctx.beginPath();
      this.ctx.arc(x, y, miniR, 0, Math.PI * 2);
      this.ctx.fill();

      this.ctx.strokeStyle = "rgba(0,0,0,0.2)";
      this.ctx.lineWidth = 0.8;
      this.ctx.stroke();

      this.ctx.fillStyle = "#fff";
      this.ctx.beginPath();
      this.ctx.arc(x, y, miniR * 0.45, 0, Math.PI * 2);
      this.ctx.fill();

      this.ctx.fillStyle = "#111";
      this.ctx.font = "bold 7px Georgia";
      this.ctx.textAlign = "center";
      this.ctx.textBaseline = "middle";
      this.ctx.fillText(String(n), x, y + 0.2);
    }
  }

  renderPocketFlashes(pocketFlashes: PocketFlash[]): void {
    const maxLife = TIMING.pocketFlashDuration;

    for (const flash of pocketFlashes) {
      const t = flash.life / maxLife;
      this.ctx.strokeStyle = `rgba(255,245,186,${0.55 * t})`;
      this.ctx.lineWidth = 2.5;
      this.ctx.beginPath();
      this.ctx.arc(flash.x, flash.y, 13 + (1 - t) * 16, 0, Math.PI * 2);
      this.ctx.stroke();

      this.ctx.strokeStyle = `rgba(255,255,255,${0.25 * t})`;
      this.ctx.lineWidth = 1.2;
      this.ctx.beginPath();
      this.ctx.arc(flash.x, flash.y, 8 + (1 - t) * 11, 0, Math.PI * 2);
      this.ctx.stroke();
    }
  }

  private drawRoundedRect(x: number, y: number, width: number, height: number, radius: number): void {
    this.ctx.beginPath();
    this.ctx.moveTo(x + radius, y);
    this.ctx.arcTo(x + width, y, x + width, y + height, radius);
    this.ctx.arcTo(x + width, y + height, x, y + height, radius);
    this.ctx.arcTo(x, y + height, x, y, radius);
    this.ctx.arcTo(x, y, x + width, y, radius);
    this.ctx.closePath();
  }
}
