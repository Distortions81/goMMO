package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
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
	creatureList = make(map[uint32]*playerData)

	StartLog()
	LogDaemon()

	dMode := flag.Bool("dev", false, "dev mode enable")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		doLog(true, "pprof started")
		pprof.StartCPUProfile(f)
		go func() {
			time.Sleep(time.Minute)
			pprof.StopCPUProfile()
			doLog(true, "pprof complete")
		}()
	}

	// Temporary for testing
	if *dMode {
		devMode = true
		authSite = "https://127.0.0.1/gs"
		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	// TODO: use compile flag instead
	if runtime.GOARCH == "wasm" {
		WASMMode = true
	}

	if !readIndex() {
		return
	}

	// Set up ebiten and window
	if WASMMode {
		vSync = false
		ebiten.SetVsyncEnabled(false)
		settingItems[0].Enabled = false
	} else {
		vSync = true
		ebiten.SetVsyncEnabled(true)
		settingItems[0].Enabled = true
	}
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(640, 360, 8192, 8192)
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("goMMO")

	helpText, _ = getText("help")
	loadSprites()
	initSpritePacks()

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
	loadPlayerModes()

	updateFonts()

	halfScreenX = screenX / 2
	halfScreenY = screenY / 2

	return &Game{}
}

var maxScreenSize = 1080

// Ebiten resize handling
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
		if !smallMode && (screenX < 720 || screenY < 720) {
			smallMode = true
			//uiScale = 0.5
			maxScreenSize = 1080 / 2
			ebiten.SetWindowSize(screenX/2, screenY/2)
		} else if smallMode && (screenX >= 720 || screenY >= 720) {
			smallMode = false
			//uiScale = 0.5
			maxScreenSize = 1080
		}

		clampUIWindows()
	}
	return int(outsideWidth), int(outsideHeight)
}
