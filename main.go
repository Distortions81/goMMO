package main

import (
	"crypto/tls"
	"net/http"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	transport *http.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	client *http.Client = &http.Client{Transport: transport}
)

const windowStartX = 1500
const windowStartY = 800

const halfWindowStartX = windowStartX / 2
const halfWindowStartY = windowStartY / 2

func main() {
	StartLog()
	LogDaemon()

	/* Temporary for testing */
	authSite = "https://127.0.0.1/gs"
	transport.TLSClientConfig.InsecureSkipVerify = true

	/* TODO: use compile flag instead */
	if runtime.GOARCH == "wasm" {
		//doLog(false, "WASM mode")
		WASMMode = true
	}

	/* Set up ebiten and window */
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(60)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowSize(windowStartX, windowStartY)
	ebiten.SetWindowTitle("goMMO")

	loadTest()

	if err := ebiten.RunGameWithOptions(newGame(), &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryOpenGL}); err != nil {
		return
	}
}

func newGame() *Game {
	go connectServer()

	return &Game{}
}

/* Ebiten resize handling */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(outsideWidth), int(outsideHeight)
}
