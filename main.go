package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"flag"
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
var screenWidth int = 512
var screenHeight int = 512

var HscreenWidth int
var HscreenHeight int

func main() {
	playerList = make(map[uint32]*playerData)

	StartLog()
	LogDaemon()

	devMode := flag.Bool("dev", false, "dev mode enable")
	flag.Parse()

	/* Temporary for testing */
	if *devMode {
		authSite = "https://127.0.0.1/gs"
		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	/* TODO: use compile flag instead */
	if runtime.GOARCH == "wasm" {
		//doLog(false, "WASM mode")
		WASMMode = true
	}

	/* Set up ebiten and window */
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(120)
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
	go connectServer()
	updateFonts()

	HscreenWidth = screenWidth / 2
	HscreenHeight = screenHeight / 2

	return &Game{}
}

const maxScreenSize = 1024

/* Ebiten resize handling */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != screenWidth || outsideHeight != screenHeight {
		if outsideWidth > maxScreenSize {
			outsideWidth = maxScreenSize
		}
		if outsideHeight > maxScreenSize {
			outsideHeight = maxScreenSize
		}

		screenWidth = outsideWidth
		screenHeight = outsideHeight

		HscreenWidth = outsideWidth / 2
		HscreenHeight = outsideHeight / 2

		var buf []byte
		outbuf := bytes.NewBuffer(buf)

		var outHeight uint16 = uint16(screenHeight)
		var outWidth uint16 = uint16(screenWidth)
		binary.Write(outbuf, binary.LittleEndian, &outHeight)
		binary.Write(outbuf, binary.LittleEndian, &outWidth)

		sendCommand(CMD_SCREENSIZE, outbuf.Bytes())
	}
	return int(outsideWidth), int(outsideHeight)
}
