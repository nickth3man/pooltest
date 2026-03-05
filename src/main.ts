/**
 * Main entry point for the Billiards game
 * Initializes the game and sets up UI callbacks
 */

import { Game } from "./Game.js";

// Get DOM elements
const canvas = document.getElementById("table") as HTMLCanvasElement;
const sunkCountEl = document.getElementById("sunk-count") as HTMLElement;
const stateTextEl = document.getElementById("state-text") as HTMLElement;
const shootIndicatorEl = document.getElementById("shoot-indicator") as HTMLElement;
const restartBtn = document.getElementById("restart-btn") as HTMLButtonElement;

// Validate required elements
if (!canvas || !sunkCountEl || !stateTextEl || !shootIndicatorEl || !restartBtn) {
  throw new Error("Required DOM elements not found");
}

// Initialize game with UI callbacks
const game = new Game(canvas, {
  onSunkCountChange: (count) => {
    sunkCountEl.textContent = String(count);
  },
  onStateChange: (state) => {
    stateTextEl.textContent = state;
  },
  onReadyToShootChange: (ready) => {
    shootIndicatorEl.style.visibility = ready ? "visible" : "hidden";
  }
});

// Setup restart button
restartBtn.addEventListener("click", () => {
  game.reset();
});

// Start the game
game.start();

// Handle page unload
document.addEventListener("visibilitychange", () => {
  if (document.hidden) {
    // Optional: pause game when tab is hidden
  }
});

// Cleanup on page unload
window.addEventListener("beforeunload", () => {
  game.destroy();
});
