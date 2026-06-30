"""Pool Python entry point.

Opens an 800x600 window, fills it with green felt, and runs a 60 FPS loop.
Exits cleanly on ESC or window close. Physics, balls, and cushions are
intentionally NOT implemented yet — this is the runnable scaffold.
"""

from __future__ import annotations

import sys

import pygame

# --- Configuration -----------------------------------------------------------

# Authentic billiard-cloth green. Adjust to taste.
FELT_COLOR: tuple[int, int, int] = (15, 92, 50)

WINDOW_SIZE: tuple[int, int] = (800, 600)
WINDOW_TITLE: str = "Pool - Python"
TARGET_FPS: int = 60


# --- Game loop ---------------------------------------------------------------


def main() -> int:
    """Run the game loop. Returns the process exit code."""
    pygame.init()
    try:
        screen = pygame.display.set_mode(WINDOW_SIZE)
        pygame.display.set_caption(WINDOW_TITLE)
        clock = pygame.time.Clock()

        running = True
        while running:
            # --- Event polling ---
            for event in pygame.event.get():
                if event.type == pygame.QUIT or (
                    event.type == pygame.KEYDOWN and event.key == pygame.K_ESCAPE
                ):
                    running = False

            # --- Update (physics, AI, etc. go here later) ---

            # --- Render ---
            screen.fill(FELT_COLOR)

            pygame.display.flip()
            clock.tick(TARGET_FPS)

        return 0
    finally:
        pygame.quit()


if __name__ == "__main__":
    sys.exit(main())
