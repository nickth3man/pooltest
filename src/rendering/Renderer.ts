/**
 * Composite Renderer - Orchestrates specialized renderers
 * This replaces the monolithic Renderer.ts with a cleaner separation of concerns
 */

import type { PocketFlash } from "../types.js";
import { Ball } from "../models/Ball.js";
import { Table } from "../models/Table.js";
import { Vec2 } from "../models/Vector2.js";
import { GameState } from "../game/GameState.js";
import { BALL_RADIUS, CANVAS_WIDTH, CANVAS_HEIGHT, VISUAL } from "../constants.js";
import { TableRenderer } from "./TableRenderer.js";
import { BallRenderer } from "./BallRenderer.js";
import { UIRenderer } from "./UIRenderer.js";

export interface RenderState {
  balls: Ball[];
  cueBall: Ball;
  table: Table;
  gameState: GameState;
  aimDirection: Vec2;
  lockedAimDirection: Vec2;
  isDragging: boolean;
  pullDistance: number;
  pocketFlashes: PocketFlash[];
  sunkNumbers: number[];
}

export class Renderer {
  private ctx: CanvasRenderingContext2D;
  private tableRenderer: TableRenderer;
  private ballRenderer: BallRenderer;
  private uiRenderer: UIRenderer;
  private feltPattern: CanvasPattern | null = null;

  constructor(canvas: HTMLCanvasElement) {
    const ctx = canvas.getContext("2d");
    if (!ctx) {
      throw new Error("Failed to get 2D rendering context");
    }
    
    this.ctx = ctx;
    this.tableRenderer = new TableRenderer(ctx);
    this.ballRenderer = new BallRenderer(ctx);
    this.uiRenderer = new UIRenderer(ctx);
    this.feltPattern = this.createFeltPattern(ctx);
  }

  render(state: RenderState): void {
    this.clear();
    this.tableRenderer.render(state.table, this.feltPattern);
    this.drawCueAndGuide(state);
    this.ballRenderer.renderBalls(state.balls);
    
    if (state.isDragging) {
      this.uiRenderer.renderPowerMeter(state.pullDistance);
    }
    
    this.uiRenderer.renderSunkRow(state.sunkNumbers);
    this.uiRenderer.renderPocketFlashes(state.pocketFlashes);
  }

  private clear(): void {
    this.ctx.clearRect(0, 0, CANVAS_WIDTH, CANVAS_HEIGHT);
  }

  private createFeltPattern(ctx: CanvasRenderingContext2D): CanvasPattern | null {
    const patternCanvas = document.createElement("canvas");
    patternCanvas.width = 140;
    patternCanvas.height = 140;
    const pctx = patternCanvas.getContext("2d");
    if (!pctx) return null;

    const gradient = pctx.createLinearGradient(0, 0, patternCanvas.width, patternCanvas.height);
    gradient.addColorStop(0, "rgba(255,255,255,0.03)");
    gradient.addColorStop(1, "rgba(0,0,0,0.03)");
    pctx.fillStyle = gradient;
    pctx.fillRect(0, 0, patternCanvas.width, patternCanvas.height);

    for (let i = 0; i < 850; i++) {
      const x = Math.random() * patternCanvas.width;
      const y = Math.random() * patternCanvas.height;
      const alpha = 0.035 + Math.random() * 0.04;
      pctx.fillStyle = `rgba(255,255,255,${alpha})`;
      pctx.fillRect(x, y, 1, 1);
    }

    for (let i = 0; i < 450; i++) {
      const x = Math.random() * patternCanvas.width;
      const y = Math.random() * patternCanvas.height;
      const alpha = 0.02 + Math.random() * 0.03;
      pctx.fillStyle = `rgba(0,0,0,${alpha})`;
      pctx.fillRect(x, y, 1, 1);
    }

    pctx.strokeStyle = "rgba(255,255,255,0.022)";
    pctx.lineWidth = 1;
    for (let i = 0; i < 24; i++) {
      const y = (i / 24) * patternCanvas.height;
      pctx.beginPath();
      pctx.moveTo(0, y);
      pctx.lineTo(patternCanvas.width, y + 8);
      pctx.stroke();
    }

    return ctx.createPattern(patternCanvas, "repeat");
  }

  private drawCueAndGuide(state: RenderState): void {
    const { cueBall, gameState, isDragging, pullDistance, lockedAimDirection } = state;

    if (!cueBall.inPlay || gameState !== GameState.AIMING) return;

    const aimDir = isDragging ? lockedAimDirection : state.aimDirection;
    this.drawAimGuide(cueBall, aimDir);
    this.drawCueStick(cueBall, aimDir, isDragging ? pullDistance : 0);
  }

  private drawAimGuide(cueBall: Ball, aimDir: Vec2): void {
    const startX = cueBall.x + aimDir.x * (BALL_RADIUS + VISUAL.guideStartOffset);
    const startY = cueBall.y + aimDir.y * (BALL_RADIUS + VISUAL.guideStartOffset);
    const endX = cueBall.x + aimDir.x * VISUAL.guideLength;
    const endY = cueBall.y + aimDir.y * VISUAL.guideLength;

    const guideGradient = this.ctx.createLinearGradient(startX, startY, endX, endY);
    guideGradient.addColorStop(0, "rgba(255,255,255,0.62)");
    guideGradient.addColorStop(0.65, "rgba(255,245,198,0.52)");
    guideGradient.addColorStop(1, "rgba(255,255,255,0.14)");

    this.ctx.save();
    this.ctx.setLineDash([7, 7]);
    this.ctx.lineWidth = 2;
    this.ctx.strokeStyle = guideGradient;
    this.ctx.shadowColor = "rgba(255,255,255,0.2)";
    this.ctx.shadowBlur = 6;
    this.ctx.beginPath();
    this.ctx.moveTo(startX, startY);
    this.ctx.lineTo(endX, endY);
    this.ctx.stroke();

    this.ctx.setLineDash([]);
    this.ctx.shadowBlur = 0;
    this.ctx.fillStyle = "rgba(255,250,222,0.78)";
    this.ctx.beginPath();
    this.ctx.arc(endX, endY, 2.2, 0, Math.PI * 2);
    this.ctx.fill();

    this.ctx.restore();
  }

  private drawCueStick(cueBall: Ball, aimDir: Vec2, pull: number): void {
    const gap = BALL_RADIUS + 2 + pull * 0.6;

    const tipX = cueBall.x - aimDir.x * gap;
    const tipY = cueBall.y - aimDir.y * gap;
    const buttX = tipX - aimDir.x * VISUAL.cueStickLength;
    const buttY = tipY - aimDir.y * VISUAL.cueStickLength;

    const shaftGradient = this.ctx.createLinearGradient(tipX, tipY, buttX, buttY);
    shaftGradient.addColorStop(0, "#f5ecd0");
    shaftGradient.addColorStop(0.62, "#e4d6af");
    shaftGradient.addColorStop(1, "#b79f73");

    this.ctx.save();
    this.ctx.lineCap = "round";
    this.ctx.strokeStyle = shaftGradient;
    this.ctx.lineWidth = 6;
    this.ctx.beginPath();
    this.ctx.moveTo(tipX, tipY);
    this.ctx.lineTo(buttX, buttY);
    this.ctx.stroke();

    const buttEndX = buttX - aimDir.x * VISUAL.cueButtLength;
    const buttEndY = buttY - aimDir.y * VISUAL.cueButtLength;
    const buttGradient = this.ctx.createLinearGradient(buttX, buttY, buttEndX, buttEndY);
    buttGradient.addColorStop(0, "#8f562f");
    buttGradient.addColorStop(1, "#5f3218");

    this.ctx.strokeStyle = buttGradient;
    this.ctx.lineWidth = 9;
    this.ctx.beginPath();
    this.ctx.moveTo(buttX, buttY);
    this.ctx.lineTo(buttEndX, buttEndY);
    this.ctx.stroke();

    this.ctx.strokeStyle = "#7db8de";
    this.ctx.lineWidth = 3;
    this.ctx.beginPath();
    this.ctx.moveTo(tipX + aimDir.x * 2, tipY + aimDir.y * 2);
    this.ctx.lineTo(tipX - aimDir.x * 2, tipY - aimDir.y * 2);
    this.ctx.stroke();
    this.ctx.restore();
  }

  resize(): void {
    // Canvas is fixed size
  }
}
