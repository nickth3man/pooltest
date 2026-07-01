import Phaser from 'phaser';

class BootScene extends Phaser.Scene {
  constructor() {
    super('BootScene');
  }

  preload(): void {
    // Nothing to load yet — no external assets in the scaffold.
  }

  create(): void {
    const { width, height } = this.scale;

    this.cameras.main.setBackgroundColor(0x0f5a2e);

    const ballRadius = 24;
    this.add.circle(width / 2, height / 2, ballRadius, 0xffffff);

    this.add
      .circle(width / 2, height / 2, ballRadius, 0x000000, 0.15)
      .setStrokeStyle(2, 0x000000, 0.35);

    this.add
      .text(width / 2, height - 32, 'pool-ts — Phaser 4 scaffold', {
        fontFamily: 'system-ui, sans-serif',
        fontSize: '16px',
        color: '#e8f3ec',
      })
      .setOrigin(0.5);
  }

  update(_time: number, _delta: number): void {
    // Reserved for per-frame physics loop (next phase).
  }
}

const config: Phaser.Types.Core.GameConfig = {
  type: Phaser.AUTO,
  width: 1024,
  height: 768,
  parent: 'game-container',
  backgroundColor: '#0b0f0c',
  scene: [BootScene],
};

new Phaser.Game(config);
