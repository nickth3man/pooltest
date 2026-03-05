import { vi } from "vitest";
import { GameState } from "../../game/GameState.js";
import { Ball } from "../../models/Ball.js";
import { Table } from "../../models/Table.js";
import { Vec2 } from "../../models/Vector2.js";
import type { RenderState } from "../../rendering/Renderer.js";

export function createMockCanvasContext(): CanvasRenderingContext2D {
  const gradient = {
    addColorStop: vi.fn()
  } as unknown as CanvasGradient;

  const context = {
    canvas: document.createElement("canvas"),
    clearRect: vi.fn(),
    createLinearGradient: vi.fn(() => gradient),
    createRadialGradient: vi.fn(() => gradient),
    createPattern: vi.fn(() => ({} as CanvasPattern)),
    fillRect: vi.fn(),
    beginPath: vi.fn(),
    moveTo: vi.fn(),
    lineTo: vi.fn(),
    stroke: vi.fn(),
    arc: vi.fn(),
    fill: vi.fn(),
    save: vi.fn(),
    restore: vi.fn(),
    setLineDash: vi.fn(),
    ellipse: vi.fn(),
    arcTo: vi.fn(),
    closePath: vi.fn(),
    fillText: vi.fn()
  } satisfies Partial<CanvasRenderingContext2D>;

  return context as unknown as CanvasRenderingContext2D;
}

export function installCanvasContextMock(context: CanvasRenderingContext2D): ReturnType<typeof vi.spyOn> {
  return vi
    .spyOn(HTMLCanvasElement.prototype, "getContext")
    .mockImplementation(() => context);
}

export function createRenderState(overrides: Partial<RenderState> = {}): RenderState {
  const table = Table.createDefault();
  const cueBall = Ball.createCueBall(table.headSpot.x, table.headSpot.y);

  return {
    balls: [cueBall],
    cueBall,
    table,
    gameState: GameState.AIMING,
    aimDirection: new Vec2(1, 0),
    lockedAimDirection: new Vec2(1, 0),
    isDragging: false,
    pullDistance: 0,
    pocketFlashes: [],
    sunkNumbers: [],
    ...overrides
  };
}
