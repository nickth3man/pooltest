import { describe, expect, it } from "vitest";
import {
  CANVAS_WIDTH,
  TABLE_HEIGHT,
  TABLE_WIDTH,
  TABLE_X,
  TABLE_Y,
} from "../constants.js";

describe("table layout constants", () => {
  it("centers the pool table horizontally in the canvas", () => {
    expect(TABLE_X).toBe((CANVAS_WIDTH - TABLE_WIDTH) / 2);
  });

  it("keeps pool table dimensions unchanged", () => {
    expect(TABLE_WIDTH).toBe(784);
    expect(TABLE_HEIGHT).toBe(432);
    expect(TABLE_Y).toBe(34);
  });
});
