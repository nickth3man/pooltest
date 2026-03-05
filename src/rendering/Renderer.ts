/**
 * Composite Renderer - Orchestrates specialized renderers
 * This replaces the monolithic Renderer.ts with a cleaner separation of concerns
 */

import type { PocketFlash } from "../types.js";
import { Ball } from "../models/Ball.js";
import { Table } from "../models/Table.js";
import { Vec2 } from "../models/Vector2.js";
import { GameState } from "../game/GameState.js";
import { BALL_RADIUS, CANVAS_WIDTH, CANVAS_HEIGHT } from "../constants.js";
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

    return ctx.createPattern(patternCanvas, "repeat");
  }

  private drawCueAndGuide(state: RenderState): void {
    const { cueBall, gameState, isDragging, pullDistance, lockedAimDirection } = state;

    if (!cueBall.inPlay || gameState !== GameState.AIMING) return;

    const aimDir = isDragging ? lockedAimDirection : state.aimDirection;
    const VISUAL = {
      guideLength: 220,
      guideStartOffset: 1,
      maxPullDistance: 140,
      cueStickLength: 175,
      cueButtLength: 26
    };

    this.ctx.save();
    this.ctx.setLineDash([7, 7]);
    this.ctx.lineWidth = 2;
    this.ctx.strokeStyle = "rgba(255,255,255,0.5)";
    this.ctx.beginPath();
    this.ctx.moveTo(
      cueBall.x + aimDir.x * (BALL_RADIUS + VISUAL.guideStartOffset),
      cueBall.y + aimDir.y * (BALL_RADIUS + VISUAL.guideStartOffset)
    );
    this.ctx.lineTo(
      cueBall.x + aimDir.x * VISUAL.guideLength,
      cueBall.y + aimDir.y * VISUAL.guideLength
    );
    this.ctx.stroke();
    this.ctx.restore();

    const pull = isDragging ? pullDistance : 0;
    const gap = BALL_RADIUS + 2 + pull * 0.6;

    const tipX = cueBall.x - aimDir.x * gap;
    const tipY = cueBall.y - aimDir.y * gap;
    const buttX = tipX - aimDir.x * VISUAL.cueStickLength;
    const buttY = tipY - aimDir.y * VISUAL.cueStickLength;

    this.ctx.lineCap = "round";
    this.ctx.strokeStyle = "#ecdfbf";
    this.ctx.lineWidth = 6;
    this.ctx.beginPath();
    this.ctx.moveTo(tipX, tipY);
    this.ctx.lineTo(buttX, buttY);
    this.ctx.stroke();

    this.ctx.strokeStyle = "#7b4524";
    this.ctx.lineWidth = 9;
    this.ctx.beginPath();
    this.ctx.moveTo(buttX, buttY);
    this.ctx.lineTo(
      buttX - aimDir.x * VISUAL.cueButtLength,
      buttY - aimDir.y * VISUAL.cueButtLength
    );
    this.ctx.stroke();

    this.ctx.strokeStyle = "#7db8de";
    this.ctx.lineWidth = 3;
    this.ctx.beginPath();
    this.ctx.moveTo(tipX + aimDir.x * 2, tipY + aimDir.y * 2);
    this.ctx.lineTo(tipX - aimDir.x * 2, tipY - aimDir.y * 2);
    this.ctx.stroke();
  }

  resize(): void {
    // Canvas is fixed size
  }
}
