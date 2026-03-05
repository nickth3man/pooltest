import type { PocketFlash, RenderContext } from "../types.js";
import { BALL_COLORS, VISUAL, SHOT, TIMING } from "../constants.js";

/**
 * UIRenderer - Handles all UI overlay elements
 * 
 * Includes:
 * - Power meter (shows shot power during drag)
 * - Sunk ball display (mini balls at top)
 * - Pocket flash effects (celebration animation)
 */
export class UIRenderer {
  constructor(private ctx: RenderContext) {}

  renderPowerMeter(pullDistance: number): void {
    if (pullDistance <= 0) return;  // Hide when not dragging

    // Destructure power meter dimensions for readability
    const { x, y, width, height } = {
      x: VISUAL.powerMeterX,
      y: VISUAL.powerMeterY,
      width: VISUAL.powerMeterWidth,
      height: VISUAL.powerMeterHeight
    };

    // Convert pull distance to power ratio (0.0 to 1.0+)
    const power = (pullDistance * SHOT.powerScale) / SHOT.maxPower;
    const fillHeight = height * Math.min(power, 1);  // Cap at full height

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

      // Spring physics bounce effect
      const springT = Math.pow(1 - t, 2) * Math.sin((1 - t) * Math.PI * 3);
      const bounceScale = 1 + springT * 0.4;

      // Outer glow with spring bounce
      const outerRadius = (13 + (1 - t) * 16) * bounceScale;
      this.ctx.strokeStyle = `rgba(255,245,186,${0.55 * t})`;
      this.ctx.lineWidth = 2.5 + springT * 1.5;
      this.ctx.shadowColor = "rgba(255,245,186,0.6)";
      this.ctx.shadowBlur = 8 + springT * 6;
      this.ctx.beginPath();
      this.ctx.arc(flash.x, flash.y, outerRadius, 0, Math.PI * 2);
      this.ctx.stroke();

      // Inner ring with spring bounce
      const innerRadius = (8 + (1 - t) * 11) * bounceScale;
      this.ctx.strokeStyle = `rgba(255,255,255,${0.25 * t})`;
      this.ctx.lineWidth = 1.2 + springT * 0.8;
      this.ctx.shadowBlur = 4 + springT * 3;
      this.ctx.beginPath();
      this.ctx.arc(flash.x, flash.y, innerRadius, 0, Math.PI * 2);
      this.ctx.stroke();

      // Sparkle particles for celebration
      if (t > 0.3) {
        const sparkleCount = 6;
        const sparkleT = (t - 0.3) / 0.7;
        for (let i = 0; i < sparkleCount; i++) {
          const angle = (i / sparkleCount) * Math.PI * 2 + (1 - t) * 2;
          const dist = outerRadius * (0.6 + sparkleT * 0.8);
          const sx = flash.x + Math.cos(angle) * dist;
          const sy = flash.y + Math.sin(angle) * dist;
          const sparkleSize = (2 + springT * 1.5) * sparkleT;

          this.ctx.fillStyle = `rgba(255,250,220,${0.8 * sparkleT})`;
          this.ctx.beginPath();
          this.ctx.arc(sx, sy, sparkleSize, 0, Math.PI * 2);
          this.ctx.fill();
        }
      }

      this.ctx.shadowBlur = 0;
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
