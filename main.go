package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"runtime"
	"time"

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
var screenWidth int = 1080
var screenHeight int = 1080

var HscreenWidth int
var HscreenHeight int

func main() {
	defer time.Sleep(time.Second * 2) //Wait for log to close

	playerList = make(map[uint32]*playerData)

	StartLog()
	LogDaemon()

	devMode := flag.Bool("dev", false, "dev mode enable")
	flag.Parse()

	/* Temporary for testing */
	if *devMode {
		gDevMode = true
		authSite = "https://127.0.0.1/gs"
		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	/* TODO: use compile flag instead */
	if runtime.GOARCH == "wasm" {
		WASMMode = true
	}

	if !readIndex() {
		return
	}

	/* Set up ebiten and window */
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("goMMO")

	loadTest()

	if err := ebiten.RunGameWithOptions(newGame(), &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryOpenGL}); err != nil {
		return
	}
}

func newGame() *Game {
	updateFonts()
	go connectServer()

	initWindows()
	settingsToggle()

	HscreenWidth = screenWidth / 2
	HscreenHeight = screenHeight / 2

	return &Game{}
}

const maxScreenSize = 1080

/* Ebiten resize handling */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	if outsideWidth > maxScreenSize {
		outsideWidth = maxScreenSize
	}
	if outsideHeight > maxScreenSize {
		outsideHeight = maxScreenSize
	}

	if outsideWidth != screenWidth || outsideHeight != screenHeight {
		screenWidth = outsideWidth
		screenHeight = outsideHeight

		HscreenWidth = outsideWidth / 2
		HscreenHeight = outsideHeight / 2
	}
	return int(outsideWidth), int(outsideHeight)
}
