import { describe, expect, it, vi } from "vitest";
import { InputHandler } from "../game/InputHandler.js";
import { EventBus } from "../core/EventBus.js";

describe("InputHandler", () => {
  it("registers and unregisters the same listener references", () => {
    const canvas = document.createElement("canvas");
    canvas.width = 900;
    canvas.height = 500;

    const addSpy = vi.spyOn(canvas, "addEventListener");
    const removeSpy = vi.spyOn(canvas, "removeEventListener");

    const handler = new InputHandler(canvas, new EventBus());

    const addedByType = new Map<string, Set<EventListenerOrEventListenerObject>>();
    for (const call of addSpy.mock.calls) {
      const type = call[0];
      const listener = call[1];
      if (!addedByType.has(type)) {
        addedByType.set(type, new Set());
      }
      addedByType.get(type)?.add(listener);
    }

    handler.destroy();

    for (const call of removeSpy.mock.calls) {
      const type = call[0];
      const listener = call[1];
      expect(addedByType.get(type)?.has(listener)).toBe(true);
    }

    expect(removeSpy).toHaveBeenCalledTimes(addSpy.mock.calls.length);
  });
});
