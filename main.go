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
	ebiten.SetWindowSize(512, 512)
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
