import type { Table } from "../models/Table.js";
import type { RenderContext } from "../types.js";

/**
 * TableRenderer - Renders the billiards table with all visual elements
 * 
 * Drawing order (back to front):
 * 1. Ambient shadow (table float effect)
 * 2. Wooden rails (outer frame)
 * 3. Rail details (grain, bevels)
 * 4. Green felt surface with texture pattern
 * 5. Rail inset border
 * 6. Pocket shadows for depth
 * 7. Pocket holes with beveled edges
 * 8. Diamond markers
 * 9. Head/foot spot markers
 */
export class TableRenderer {
  private feltPattern: CanvasPattern | null = null;

  constructor(private ctx: RenderContext) {}

  render(table: Table, feltPattern: CanvasPattern | null): void {
    this.drawAmbientShadow(table);      // Floating effect
    this.drawRails(table);              // Wooden outer frame
    this.drawRailDetails(table);        // Grain and bevels
    this.drawFelt(table, feltPattern);  // Green playing surface
    this.drawRailInset(table);          // Border between rail and felt
    this.drawPocketShadows(table);      // Depth shadows around pockets
    this.drawPockets(table);            // The pocket holes with depth
    this.drawDiamonds(table);           // Aiming markers
    this.drawSpotMarkers(table);        // Head and foot spot circles
  }

  /**
   * Creates a subtle noise pattern for felt texture
   * Cached for performance - only created once
   */
  createFeltPattern(ctx: CanvasRenderingContext2D): CanvasPattern | null {
    if (this.feltPattern) return this.feltPattern;

    const patternCanvas = document.createElement("canvas");
    patternCanvas.width = 4;
    patternCanvas.height = 4;
    const pCtx = patternCanvas.getContext("2d");
    if (!pCtx) return null;

    // Base: transparent
    pCtx.clearRect(0, 0, 4, 4);

    // Add subtle noise pixels
    pCtx.fillStyle = "rgba(20, 60, 35, 0.04)";
    pCtx.fillRect(0, 0, 1, 1);
    pCtx.fillRect(2, 1, 1, 1);
    pCtx.fillRect(1, 3, 1, 1);
    pCtx.fillRect(3, 2, 1, 1);

    pCtx.fillStyle = "rgba(40, 90, 55, 0.03)";
    pCtx.fillRect(1, 0, 1, 1);
    pCtx.fillRect(3, 3, 1, 1);
    pCtx.fillRect(0, 2, 1, 1);

    try {
      this.feltPattern = ctx.createPattern(patternCanvas, "repeat");
    } catch {
      return null;
    }

    return this.feltPattern;
  }

  private drawAmbientShadow(table: Table): void {
    const spread = 28;
    const gradient = this.ctx.createLinearGradient(
      table.x - spread,
      table.y - spread,
      table.x + table.width + spread,
      table.y + table.height + spread
    );
    gradient.addColorStop(0, "rgba(0, 0, 0, 0)");
    gradient.addColorStop(0.4, "rgba(0, 0, 0, 0.03)");
    gradient.addColorStop(0.5, "rgba(0, 0, 0, 0.08)");
    gradient.addColorStop(0.6, "rgba(0, 0, 0, 0.03)");
    gradient.addColorStop(1, "rgba(0, 0, 0, 0)");

    this.ctx.fillStyle = gradient;
    this.ctx.fillRect(
      table.x - spread,
      table.y - spread,
      table.width + spread * 2,
      table.height + spread * 2
    );
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

  private drawRailDetails(table: Table): void {
    this.ctx.save();
    
    // Subtle inner bevel highlight on top edge
    const inset = 6;
    this.ctx.strokeStyle = "rgba(255,248,230,0.08)";
    this.ctx.lineWidth = 1;
    this.roundedRect(
      table.x + inset,
      table.y + inset,
      table.width - inset * 2,
      table.height - inset * 2,
      12
    );
    this.ctx.stroke();

    // Subtle wood grain lines - horizontal variation
    this.ctx.strokeStyle = "rgba(30,20,12,0.06)";
    this.ctx.lineWidth = 1;
    
    const grainY1 = table.y + table.height * 0.25;
    this.ctx.beginPath();
    this.ctx.moveTo(table.x + 20, grainY1);
    this.ctx.lineTo(table.x + table.width - 20, grainY1);
    this.ctx.stroke();

    const grainY2 = table.y + table.height * 0.7;
    this.ctx.beginPath();
    this.ctx.moveTo(table.x + 25, grainY2);
    this.ctx.lineTo(table.x + table.width - 30, grainY2);
    this.ctx.stroke();

    // Subtle vertical grain
    const grainX = table.x + table.width * 0.6;
    this.ctx.beginPath();
    this.ctx.moveTo(grainX, table.y + 15);
    this.ctx.lineTo(grainX, table.y + table.height - 15);
    this.ctx.stroke();

    // Corner brass accent glow
    const cornerRadius = 40;
    const cornerGradient1 = this.ctx.createRadialGradient(
      table.x + cornerRadius, table.y + cornerRadius, 0,
      table.x + cornerRadius, table.y + cornerRadius, cornerRadius
    );
    cornerGradient1.addColorStop(0, "rgba(180,140,80,0.12)");
    cornerGradient1.addColorStop(1, "rgba(180,140,80,0)");
    this.ctx.fillStyle = cornerGradient1;
    this.ctx.fillRect(table.x, table.y, cornerRadius * 2, cornerRadius * 2);

    const cornerGradient2 = this.ctx.createRadialGradient(
      table.x + table.width - cornerRadius, table.y + table.height - cornerRadius, 0,
      table.x + table.width - cornerRadius, table.y + table.height - cornerRadius, cornerRadius
    );
    cornerGradient2.addColorStop(0, "rgba(180,140,80,0.1)");
    cornerGradient2.addColorStop(1, "rgba(180,140,80,0)");
    this.ctx.fillStyle = cornerGradient2;
    this.ctx.fillRect(table.x + table.width - cornerRadius * 2, table.y + table.height - cornerRadius * 2, cornerRadius * 2, cornerRadius * 2);

    this.ctx.restore();
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
      const r = pocket.drawRadius;

      this.ctx.save();
      
      // Outer bevel highlight (top-left light source)
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, r + 2, 0, Math.PI * 2);
      this.ctx.strokeStyle = "rgba(255,248,220,0.15)";
      this.ctx.lineWidth = 2;
      this.ctx.stroke();

      // Outer bevel shadow (bottom-right)
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, r + 2, 0, Math.PI * 2);
      this.ctx.strokeStyle = "rgba(0,0,0,0.25)";
      this.ctx.lineWidth = 2;
      this.ctx.stroke();

      // Main pocket cavity - multi-stop gradient for depth
      const cavityGradient = this.ctx.createRadialGradient(
        pocket.x - r * 0.3, pocket.y - r * 0.3, 0,
        pocket.x, pocket.y, r
      );
      cavityGradient.addColorStop(0, "#1a1a1a");
      cavityGradient.addColorStop(0.4, "#0d0d0d");
      cavityGradient.addColorStop(0.7, "#050505");
      cavityGradient.addColorStop(1, "#020202");
      this.ctx.fillStyle = cavityGradient;
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, r, 0, Math.PI * 2);
      this.ctx.fill();

      // Inner shadow (simulates pocket depth)
      const innerShadow = this.ctx.createRadialGradient(
        pocket.x, pocket.y, r * 0.3,
        pocket.x, pocket.y, r
      );
      innerShadow.addColorStop(0, "rgba(0,0,0,0)");
      innerShadow.addColorStop(0.6, "rgba(0,0,0,0.15)");
      innerShadow.addColorStop(1, "rgba(0,0,0,0.35)");
      this.ctx.fillStyle = innerShadow;
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, r, 0, Math.PI * 2);
      this.ctx.fill();

      // Brass rim reflection (subtle highlight on top)
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, r - 1, -Math.PI * 0.7, -Math.PI * 0.3);
      this.ctx.strokeStyle = "rgba(180,150,100,0.18)";
      this.ctx.lineWidth = 1.5;
      this.ctx.stroke();

      // Outer ring stroke
      this.ctx.beginPath();
      this.ctx.arc(pocket.x, pocket.y, r, 0, Math.PI * 2);
      this.ctx.strokeStyle = "rgba(0,0,0,0.5)";
      this.ctx.lineWidth = 2;
      this.ctx.stroke();

      this.ctx.restore();
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
