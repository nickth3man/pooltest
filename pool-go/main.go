// Command pool-go is a 2D 8-ball pool game built with Ebitengine.
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetTPS(60)
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Pool - Go")
	// Logical size is fixed by Layout; let the window resize and letterbox-scale.
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(screenW/2, screenH/2, -1, -1)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
