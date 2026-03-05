import { beforeEach, describe, expect, it, vi } from "vitest";
import { Game } from "../Game.js";
import { EventType } from "../core/EventBus.js";
import type { GameSession } from "../game/GameSession.js";
import { GameState } from "../game/GameState.js";
import { createMockCanvasContext, installCanvasContextMock } from "./helpers/renderTestUtils.js";

describe("Game", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("transitions to SHOOTING when a valid shot is released", () => {
    const context = createMockCanvasContext();
    const contextSpy = installCanvasContextMock(context);

    const canvas = document.createElement("canvas");
    canvas.width = 900;
    canvas.height = 500;
    vi.spyOn(canvas, "getBoundingClientRect").mockReturnValue({
      x: 0,
      y: 0,
      width: 900,
      height: 500,
      top: 0,
      right: 900,
      bottom: 500,
      left: 0,
      toJSON: () => ({})
    });

    const onStateChange = vi.fn();
    const game = new Game(canvas, { onStateChange });

    canvas.dispatchEvent(new MouseEvent("mousedown", { clientX: 300, clientY: 200 }));
    canvas.dispatchEvent(new MouseEvent("mousemove", { clientX: 20, clientY: 200 }));
    canvas.dispatchEvent(new MouseEvent("mouseup", { clientX: 20, clientY: 200 }));

    expect(game.getState().gameState).toBe(GameState.SHOOTING);
    expect(onStateChange).toHaveBeenCalledWith(GameState.SHOOTING);

    game.destroy();
    contextSpy.mockRestore();
  });

  it("cleans up event bus subscriptions on destroy", () => {
    const context = createMockCanvasContext();
    const contextSpy = installCanvasContextMock(context);
    const canvas = document.createElement("canvas");
    canvas.width = 900;
    canvas.height = 500;

    const game = new Game(canvas);
    const session = (game as unknown as { session: GameSession }).session;

    expect(session.eventBus.handlerCount(EventType.BALL_SUNK)).toBeGreaterThan(0);

    game.destroy();

    expect(session.eventBus.handlerCount(EventType.BALL_SUNK)).toBe(0);
    contextSpy.mockRestore();
  });
});
