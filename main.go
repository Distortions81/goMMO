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
var screenX int = 1080
var screenY int = 1080

var halfScreenX int
var halfScreenY int

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
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("goMMO")

	helpText, _ = getText("help")
	loadSprites()

	initToolbar()
	drawToolbar(false, false, maxItemType)

	if err := ebiten.RunGameWithOptions(newGame(), &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryOpenGL}); err != nil {
		return
	}
}

func newGame() *Game {
	updateFonts()
	go connectServer()

	initWindows()
	//settingsToggle()
	toggleHelp()
	loadOptions()

	halfScreenX = screenX / 2
	halfScreenY = screenY / 2

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

	if outsideWidth != screenX || outsideHeight != screenY {
		screenX = outsideWidth
		screenY = outsideHeight

		halfScreenX = outsideWidth / 2
		halfScreenY = outsideHeight / 2

		//Keep UI windows from going outside the window
		clampUIWindow()
	}
	return int(outsideWidth), int(outsideHeight)
}
