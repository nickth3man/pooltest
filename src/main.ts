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

let isCleanedUp = false;

const handleRestart = (): void => {
  if (!isCleanedUp) {
    game.reset();
  }
};

const handleVisibilityChange = (): void => {
  if (document.hidden) {
    // Optional: pause game when tab is hidden
  }
};

const cleanup = (): void => {
  if (isCleanedUp) {
    return;
  }

  isCleanedUp = true;
  restartBtn.removeEventListener("click", handleRestart);
  document.removeEventListener("visibilitychange", handleVisibilityChange);
  window.removeEventListener("beforeunload", cleanup);
  game.destroy();
};

restartBtn.addEventListener("click", handleRestart);
document.addEventListener("visibilitychange", handleVisibilityChange);
window.addEventListener("beforeunload", cleanup, { once: true });

// Start the game
game.start();
