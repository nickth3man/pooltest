import { describe, expect, it } from "vitest";
import { Table } from "../models/Table.js";
import { TableRenderer } from "../rendering/TableRenderer.js";
import { createMockCanvasContext } from "./helpers/renderTestUtils.js";

describe("TableRenderer", () => {
  it("renders table layers without throwing", () => {
    const context = createMockCanvasContext();
    const tableRenderer = new TableRenderer(context);
    const table = Table.createDefault();

    expect(() => tableRenderer.render(table, null)).not.toThrow();
    expect(context.fillRect).toHaveBeenCalled();
    expect(context.arc).toHaveBeenCalled();
  });
});
