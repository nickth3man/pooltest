// Command pool-go is a 2D 8-ball pool game built with Ebitengine.
package main

import (
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/user/pooltest/pool-go/internal/game"
	"github.com/user/pooltest/pool-go/internal/table"
)

func main() {
	ebiten.SetTPS(60)
	ebiten.SetWindowSize(table.ScreenW, table.ScreenH)
	ebiten.SetWindowTitle("Pool - Go")
	// Logical size is fixed by Layout; let the window resize and letterbox-scale.
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(table.ScreenW/2, table.ScreenH/2, -1, -1)

	if err := ebiten.RunGame(game.NewGame()); err != nil {
		slog.Error("game exited with error", "err", err)
		os.Exit(1)
	}
}
