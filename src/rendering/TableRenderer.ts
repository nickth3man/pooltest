import type { Table } from "../models/Table.js";
import type { RenderContext } from "../types.js";

export class TableRenderer {
  constructor(private ctx: RenderContext) {}

  render(table: Table, feltPattern: CanvasPattern | null): void {
    this.drawRails(table);
    this.drawFelt(table, feltPattern);
    this.drawRailInset(table);
    this.drawPocketShadows(table);
    this.drawPockets(table);
    this.drawDiamonds(table);
    this.drawSpotMarkers(table);
  }

  private drawRails(table: Table): void {
    const railGradient = this.ctx.createLinearGradient(
      table.x, table.y,
      table.x + table.width, table.y + table.height
    );
    railGradient.addColorStop(0, "#6c4934");
    railGradient.addColorStop(0.5, "#4f311f");
    railGradient.addColorStop(1, "#2d1a12");
    this.ctx.fillStyle = railGradient;
    this.roundedRect(table.x, table.y, table.width, table.height, 16);
    this.ctx.fill();

    this.ctx.strokeStyle = "rgba(255,236,190,0.16)";
    this.ctx.lineWidth = 2;
    this.roundedRect(table.x + 3, table.y + 3, table.width - 6, table.height - 6, 13);
    this.ctx.stroke();
  }

  private drawFelt(table: Table, feltPattern: CanvasPattern | null): void {
    const feltGradient = this.ctx.createLinearGradient(
      table.feltX, table.feltY,
      table.feltX + table.feltWidth, table.feltY + table.feltHeight
    );
    feltGradient.addColorStop(0, "#1f7044");
    feltGradient.addColorStop(1, "#0f4f2c");
    this.ctx.fillStyle = feltGradient;
    this.ctx.fillRect(table.feltX, table.feltY, table.feltWidth, table.feltHeight);

    if (feltPattern) {
      this.ctx.save();
      this.ctx.globalAlpha = 0.45;
      this.ctx.fillStyle = feltPattern;
      this.ctx.fillRect(table.feltX, table.feltY, table.feltWidth, table.feltHeight);
      this.ctx.restore();
    }

    const sheen = this.ctx.createLinearGradient(
      table.feltX,
      table.feltY,
      table.feltX,
      table.feltY + table.feltHeight
    );
    sheen.addColorStop(0, "rgba(255,255,255,0.06)");
    sheen.addColorStop(0.35, "rgba(255,255,255,0.02)");
    sheen.addColorStop(1, "rgba(0,0,0,0.12)");
    this.ctx.fillStyle = sheen;
    this.ctx.fillRect(table.feltX, table.feltY, table.feltWidth, table.feltHeight);

    const vignette = this.ctx.createRadialGradient(
      table.feltX + table.feltWidth * 0.5,
      table.feltY + table.feltHeight * 0.5,
      table.feltWidth * 0.22,
      table.feltX + table.feltWidth * 0.5,
      table.feltY + table.feltHeight * 0.5,
      table.feltWidth * 0.62
    );
    vignette.addColorStop(0, "rgba(255,255,255,0.02)");
    vignette.addColorStop(1, "rgba(0,0,0,0.22)");
    this.ctx.fillStyle = vignette;
    this.ctx.fillRect(table.feltX, table.feltY, table.feltWidth, table.feltHeight);

    this.ctx.strokeStyle = "rgba(255,255,255,0.1)";
    this.ctx.lineWidth = 1;
    this.ctx.beginPath();
    this.ctx.moveTo(table.feltX + 8, table.feltY + table.feltHeight * 0.5);
    this.ctx.lineTo(table.feltX + table.feltWidth - 8, table.feltY + table.feltHeight * 0.5);
    this.ctx.stroke();
  }

  private drawRailInset(table: Table): void {
    this.ctx.save();
    this.ctx.strokeStyle = "rgba(0,0,0,0.35)";
    this.ctx.lineWidth = 4;
    this.roundedRect(table.feltX - 1, table.feltY - 1, table.feltWidth + 2, table.feltHeight + 2, 2);
    this.ctx.stroke();
    this.ctx.restore();
  }

  private drawPocketShadows(table: Table): void {
    for (const pocket of table.pockets) {
      const shadow = this.ctx.createRadialGradient(
        pocket.x,
        pocket.y,
        pocket.drawRadius * 0.55,
        pocket.x,
        pocket.y,
        pocket.drawRadius * 1.8
      );
      shadow.addColorStop(0, "rgba(0,0,0,0.36)");
      shadow.addColorStop(1, "rgba(0,0,0,0)");
      this.ctx.fillStyle = shadow;
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, pocket.drawRadius * 1.8, 0, Math.PI * 2);
      this.ctx.fill();
    }
  }

  private drawPockets(table: Table): void {
    for (const pocket of table.pockets) {
      const gradient = this.ctx.createRadialGradient(
        pocket.x - 2, pocket.y - 2, 4,
        pocket.x, pocket.y, pocket.drawRadius
      );
      gradient.addColorStop(0, "#272727");
      gradient.addColorStop(0.5, "#131313");
      gradient.addColorStop(1, "#030303");
      this.ctx.fillStyle = gradient;
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, pocket.drawRadius, 0, Math.PI * 2);
      this.ctx.fill();

      this.ctx.strokeStyle = "rgba(0,0,0,0.45)";
      this.ctx.lineWidth = 3;
      this.ctx.stroke();
    }
  }

  private drawDiamonds(table: Table): void {
    this.ctx.fillStyle = "#dfd0a8";
    const positions = table.getDiamondPositions();
    const diamondSize = 4;

    [...positions.top, ...positions.bottom, ...positions.left, ...positions.right].forEach(pos => {
      this.drawDiamond(pos.x, pos.y, diamondSize);
    });
  }

  private drawSpotMarkers(table: Table): void {
    const spots = [table.headSpot, table.footSpot];
    this.ctx.save();
    this.ctx.strokeStyle = "rgba(241, 227, 186, 0.55)";
    this.ctx.lineWidth = 1.1;

    for (const spot of spots) {
      this.ctx.beginPath();
      this.ctx.arc(spot.x, spot.y, 3.5, 0, Math.PI * 2);
      this.ctx.stroke();
    }

    this.ctx.restore();
  }

  private drawDiamond(x: number, y: number, r: number): void {
    this.ctx.beginPath();
    this.ctx.moveTo(x, y - r);
    this.ctx.lineTo(x + r, y);
    this.ctx.lineTo(x, y + r);
    this.ctx.lineTo(x - r, y);
    this.ctx.closePath();
    this.ctx.fill();
  }

  private roundedRect(x: number, y: number, w: number, h: number, r: number): void {
    this.ctx.beginPath();
    this.ctx.moveTo(x + r, y);
    this.ctx.arcTo(x + w, y, x + w, y + h, r);
    this.ctx.arcTo(x + w, y + h, x, y + h, r);
    this.ctx.arcTo(x, y + h, x, y, r);
    this.ctx.arcTo(x, y, x + w, y, r);
    this.ctx.closePath();
  }
}
