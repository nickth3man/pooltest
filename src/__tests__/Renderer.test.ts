import { beforeEach, describe, expect, it, vi } from "vitest";
import { Renderer } from "../rendering/Renderer.js";
import { BallRenderer } from "../rendering/BallRenderer.js";
import { TableRenderer } from "../rendering/TableRenderer.js";
import { UIRenderer } from "../rendering/UIRenderer.js";
import { createRenderState, createMockCanvasContext, installCanvasContextMock } from "./helpers/renderTestUtils.js";

describe("Renderer", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("delegates rendering to specialized renderer modules", () => {
    const context = createMockCanvasContext();
    const contextSpy = installCanvasContextMock(context);

    const tableRenderSpy = vi.spyOn(TableRenderer.prototype, "render");
    const ballRenderSpy = vi.spyOn(BallRenderer.prototype, "renderBalls");
    const sunkRowSpy = vi.spyOn(UIRenderer.prototype, "renderSunkRow");
    const flashesSpy = vi.spyOn(UIRenderer.prototype, "renderPocketFlashes");

    const canvas = document.createElement("canvas");
    const renderer = new Renderer(canvas);
    renderer.render(createRenderState({ sunkNumbers: [1, 2] }));

    expect(tableRenderSpy).toHaveBeenCalledTimes(1);
    expect(ballRenderSpy).toHaveBeenCalledTimes(1);
    expect(sunkRowSpy).toHaveBeenCalledWith([1, 2]);
    expect(flashesSpy).toHaveBeenCalledTimes(1);

    contextSpy.mockRestore();
  });

  it("renders power meter only while dragging", () => {
    const context = createMockCanvasContext();
    const contextSpy = installCanvasContextMock(context);
    const powerSpy = vi.spyOn(UIRenderer.prototype, "renderPowerMeter");

    const canvas = document.createElement("canvas");
    const renderer = new Renderer(canvas);
    renderer.render(createRenderState({ isDragging: false, pullDistance: 0 }));
    renderer.render(createRenderState({ isDragging: true, pullDistance: 50 }));

    expect(powerSpy).toHaveBeenCalledTimes(1);
    expect(powerSpy).toHaveBeenCalledWith(50);

    contextSpy.mockRestore();
  });
});
