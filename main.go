package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
}

func main() {
	/* Set up ebiten and window */
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(512, 512)
	loadTest()

	if err := ebiten.RunGameWithOptions(newGame(), &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryOpenGL}); err != nil {
		return
	}
}

func newGame() *Game {

	/* Initialize the game */
	return &Game{}
}

/* Ebiten: Draw everything */
func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(walkNorth, nil)

}

/* Ebiten resize handling */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(outsideWidth), int(outsideHeight)
}

/* Input interface handler */
func (g *Game) Update() error {

	return nil
}
