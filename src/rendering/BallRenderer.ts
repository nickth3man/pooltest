import { Ball } from "../models/Ball.js";
import type { RenderContext } from "../types.js";

/**
 * BallRenderer - Handles all ball rendering with realistic shading
 * 
 * Rendering order (painter's algorithm):
 * 1. Contact shadow (ground shadow beneath ball)
 * 2. Drop shadow (floating shadow for depth)
 * 3. Ball body with radial gradient
 * 4. Outline for definition
 * 5. Number circle (only for numbered balls)
 * 6. Specular highlight for 3D effect
 */
export class BallRenderer {
  constructor(private ctx: RenderContext) {}

  renderBalls(balls: Ball[]): void {
    for (const ball of balls) {
      if (!ball.inPlay) continue;  // Skip sunk balls
      this.renderBall(ball);
    }
  }

  private renderBall(ball: Ball): void {
    // Layer rendering from bottom to top for proper depth
    this.drawContactShadow(ball);     // Shadow on the table felt
    this.drawShadow();                // Drop shadow offset
    this.drawBallBody(ball);          // Main ball surface
    this.drawBallOutline(ball);       // Edge definition
    
    if (!ball.isCue) {
      this.drawNumberCircle(ball);    // Only numbered balls get this
    }
    
    this.drawSpecularHighlight(ball); // Glossy reflection
  }

  private drawContactShadow(ball: Ball): void {
    this.ctx.save();
    this.ctx.fillStyle = "rgba(0,0,0,0.18)";
    this.ctx.beginPath();
    this.ctx.ellipse(
      ball.x + 1.3,
      ball.y + ball.radius * 0.72,
      ball.radius * 0.9,
      ball.radius * 0.5,
      0,
      0,
      Math.PI * 2
    );
    this.ctx.fill();
    this.ctx.restore();
  }

  private drawShadow(): void {
    this.ctx.save();
    this.ctx.shadowColor = "rgba(0,0,0,0.42)";
    this.ctx.shadowBlur = 8;
    this.ctx.shadowOffsetX = 2.5;
    this.ctx.shadowOffsetY = 3.5;
  }

  private drawBallBody(ball: Ball): void {
    // Create radial gradient for 3D sphere effect
    // Offset highlight to top-left for light source illusion
    const gradient = this.ctx.createRadialGradient(
      ball.x - ball.radius * 0.38,
      ball.y - ball.radius * 0.36,
      ball.radius * 0.2,
      ball.x,
      ball.y,
      ball.radius
    );
    gradient.addColorStop(0, "rgba(255,255,255,0.46)");  // Highlight
    gradient.addColorStop(0.58, ball.color);              // Base color
    gradient.addColorStop(1, this.darkenHex(ball.color, 0.4));  // Shadow edge
    this.ctx.fillStyle = gradient;
    this.ctx.beginPath();
    this.ctx.arc(ball.x, ball.y, ball.radius, 0, Math.PI * 2);
    this.ctx.fill();
    this.ctx.restore();
  }

  private drawBallOutline(ball: Ball): void {
    this.ctx.strokeStyle = "rgba(0,0,0,0.24)";
    this.ctx.lineWidth = 1.2;
    this.ctx.beginPath();
    this.ctx.arc(ball.x, ball.y, ball.radius, 0, Math.PI * 2);
    this.ctx.stroke();
  }

  private drawNumberCircle(ball: Ball): void {
    this.ctx.fillStyle = "#ffffff";
    this.ctx.beginPath();
    this.ctx.arc(ball.x, ball.y, ball.radius * 0.47, 0, Math.PI * 2);
    this.ctx.fill();

    this.ctx.strokeStyle = "rgba(0,0,0,0.25)";
    this.ctx.lineWidth = 0.9;
    this.ctx.stroke();

    this.ctx.fillStyle = "#1e1e1e";
    this.ctx.font = "bold 10px Georgia";
    this.ctx.textAlign = "center";
    this.ctx.textBaseline = "middle";
    this.ctx.fillText(String(ball.number), ball.x, ball.y + 0.3);
  }

  private drawSpecularHighlight(ball: Ball): void {
    this.ctx.fillStyle = "rgba(255,255,255,0.24)";
    this.ctx.beginPath();
    this.ctx.arc(
      ball.x - ball.radius * 0.33,
      ball.y - ball.radius * 0.38,
      ball.radius * 0.18,
      0,
      Math.PI * 2
    );
    this.ctx.fill();

    this.ctx.fillStyle = "rgba(255,255,255,0.14)";
    this.ctx.beginPath();
    this.ctx.arc(
      ball.x - ball.radius * 0.12,
      ball.y - ball.radius * 0.24,
      ball.radius * 0.09,
      0,
      Math.PI * 2
    );
    this.ctx.fill();
  }

  private darkenHex(hex: string, amount: number): string {
    const value = hex.replace("#", "");
    if (value.length !== 6) return hex;

    const clamp = (n: number): number => Math.max(0, Math.min(255, Math.round(n)));
    const scale = 1 - amount;

    const r = clamp(parseInt(value.slice(0, 2), 16) * scale);
    const g = clamp(parseInt(value.slice(2, 4), 16) * scale);
    const b = clamp(parseInt(value.slice(4, 6), 16) * scale);
    return `rgb(${r}, ${g}, ${b})`;
  }
}
