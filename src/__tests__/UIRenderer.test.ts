import { describe, expect, it } from "vitest";
import { UIRenderer } from "../rendering/UIRenderer.js";
import { createMockCanvasContext } from "./helpers/renderTestUtils.js";

describe("UIRenderer", () => {
  it("draws power meter when pull distance is positive", () => {
    const context = createMockCanvasContext();
    const uiRenderer = new UIRenderer(context);

    uiRenderer.renderPowerMeter(42);

    expect(context.fillRect).toHaveBeenCalled();
  });

  it("renders sunk row and pocket flash overlays", () => {
    const context = createMockCanvasContext();
    const uiRenderer = new UIRenderer(context);

    uiRenderer.renderSunkRow([1, 2, 3]);
    uiRenderer.renderPocketFlashes([{ x: 100, y: 110, life: 100 }]);

    expect(context.fillText).toHaveBeenCalled();
    expect(context.arc).toHaveBeenCalled();
  });
});
