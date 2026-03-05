import { Ball } from "../models/Ball.js";
import type { RenderContext } from "../types.js";

export class BallRenderer {
  constructor(private ctx: RenderContext) {}

  renderBalls(balls: Ball[]): void {
    for (const ball of balls) {
      if (!ball.inPlay) continue;
      this.renderBall(ball);
    }
  }

  private renderBall(ball: Ball): void {
    this.drawShadow();
    this.drawBallBody(ball);
    this.drawBallOutline(ball);
    
    if (!ball.isCue) {
      this.drawNumberCircle(ball);
    }
    
    this.drawSpecularHighlight(ball);
  }

  private drawShadow(): void {
    this.ctx.save();
    this.ctx.shadowColor = "rgba(0,0,0,0.35)";
    this.ctx.shadowBlur = 6;
    this.ctx.shadowOffsetX = 2;
    this.ctx.shadowOffsetY = 3;
  }

  private drawBallBody(ball: Ball): void {
    const gradient = this.ctx.createRadialGradient(
      ball.x - ball.radius * 0.38,
      ball.y - ball.radius * 0.36,
      ball.radius * 0.2,
      ball.x,
      ball.y,
      ball.radius
    );
    gradient.addColorStop(0, "rgba(255,255,255,0.34)");
    gradient.addColorStop(1, ball.color);
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

    this.ctx.fillStyle = "#1e1e1e";
    this.ctx.font = "bold 10px Trebuchet MS";
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
  }
}
