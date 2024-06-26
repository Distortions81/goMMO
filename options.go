package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TYPE_BOOL = 0
	TYPE_INT  = 1
	TYPE_TEXT = 2

	settingsFile = "data/settings.json"
)

var (
	optionWindowButtons []image.Rectangle
	settingItems        []settingType
	updateWindowButtons []image.Rectangle
)

type settingType struct {
	ConfigName string
	text       string

	textPosX   int
	textPosY   int
	textBounds image.Rectangle
	rect       image.Rectangle

	Enabled     bool
	wasmExclude bool

	action  func(item int)
	noCheck bool
}

func init() {
	defer reportPanic("options init")
	settingItems = []settingType{
		{ConfigName: "VSYNC", text: "Limit FPS (VSYNC)", action: toggleVsync, Enabled: true},
		{ConfigName: "FULLSCREEN", text: "Full Screen", action: toggleFullscreen},
		{ConfigName: "FAST-SHADOWS", text: "Fast Shadows", action: toggleFastShadow},
		{ConfigName: "MOTION-SMOOTH", text: "Motion Smoothing", action: toggleSmoothing, Enabled: true},
		{ConfigName: "NIGHT-MODE", text: "Disable Shadows", action: toggleNightShadow},
		{ConfigName: "DEBUG-TEXT", text: "Debug info-text", action: toggleInfoLine},
		{ConfigName: "GREEN-SCREEN", text: "Green-Screen", action: toggleGreenScreen},
	}
}

// Load user options settings from disk
func loadOptions() bool {
	defer reportPanic("loadOptions")
	if WASMMode {
		return false
	}

	var tempSettings []settingType

	file, err := os.ReadFile(settingsFile)

	if file != nil && err == nil {

		err := json.Unmarshal([]byte(file), &tempSettings)
		if err != nil {
			doLog(true, "loadOptions: Unmarshal failure")
			doLog(true, err.Error())
			return false
		}
	} else {
		doLog(true, "loadOptions: ReadFile failure")
		return false
	}

	doLog(true, "Settings loaded.")

	for setPos, wSetting := range settingItems {
		for _, fSetting := range tempSettings {
			if wSetting.ConfigName == fSetting.ConfigName {
				if fSetting.Enabled != wSetting.Enabled {
					settingItems[setPos].action(setPos)
				}
			}
		}
	}
	return true
}

// Save user options settings to disk
func saveOptions() {
	defer reportPanic("saveOptions")
	if WASMMode {
		return
	}

	var tempSettings []settingType
	for _, setting := range settingItems {
		if setting.ConfigName != "" {
			tempSettings = append(tempSettings, settingType{ConfigName: setting.ConfigName, Enabled: setting.Enabled})
		}
	}

	tempPath := settingsFile + ".tmp"
	finalPath := settingsFile

	outBuf := new(bytes.Buffer)
	enc := json.NewEncoder(outBuf)
	enc.SetIndent("", "\t")

	if err := enc.Encode(&tempSettings); err != nil {
		doLog(true, "saveOptions: enc.Encode failure")
		return
	}

	os.Mkdir("data", os.ModePerm)
	_, err := os.Create(tempPath)

	if err != nil {
		doLog(true, "saveOptions: os.Create failure")
		return
	}

	err = os.WriteFile(tempPath, outBuf.Bytes(), 0666)

	if err != nil {
		doLog(true, "saveOptions: WriteFile failure")
	}

	err = os.Rename(tempPath, finalPath)

	if err != nil {
		doLog(true, "Couldn't rename settings file.")
		return
	}

	doLog(true, "Settings saved.")
}

// Toggle the debug bottom-screen text
func toggleInfoLine(item int) {
	defer reportPanic("toggleInfoLine")
	if debugLine {
		debugLine = false
		settingItems[item].Enabled = false
	} else {
		debugLine = true
		settingItems[item].Enabled = true
	}
}

// Toggle the debug bottom-screen text
func toggleGreenScreen(item int) {
	defer reportPanic("toggleGreenScreen")
	if greenScreen {
		greenScreen = false
		settingItems[item].Enabled = false
	} else {
		greenScreen = true
		settingItems[item].Enabled = true
	}
}

// Toggle the use of hyper-threading
func toggleNightShadow(item int) {
	defer reportPanic("toggleNightShadow")
	if disableNightShadow {
		disableNightShadow = false
		settingItems[item].Enabled = false
	} else {
		disableNightShadow = true
		settingItems[item].Enabled = true
	}
}

// Toggle full-screen
func toggleFullscreen(item int) {
	defer reportPanic("toggleFullscreen")

	if fullscreen {
		fullscreen = false
		ebiten.SetFullscreen(false)
		settingItems[item].Enabled = false
	} else {
		fullscreen = true
		ebiten.SetFullscreen(true)
		settingItems[item].Enabled = true
	}
	buf := fmt.Sprintf("%v is now %v.",
		settingItems[item].text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorDarkOrange, time.Second*5)
}

// Toggle UI magnification
func toggleFastShadow(item int) {
	defer reportPanic("toggleFastShadow")
	if fastShadow {
		fastShadow = false
		settingItems[item].Enabled = false
	} else {
		fastShadow = true
		settingItems[item].Enabled = true
	}

	//handleResize(int(ScreenWidth), int(ScreenHeight))

	buf := fmt.Sprintf("%v is now %v.",
		settingItems[item].text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorDarkOrange, time.Second*5)
}

// Toggle debug mode
func toggleSmoothing(item int) {
	defer reportPanic("toggleSmoothing")
	if !noSmoothing {
		noSmoothing = true
		settingItems[item].Enabled = false
	} else {
		noSmoothing = false
		settingItems[item].Enabled = true
	}
	buf := fmt.Sprintf("%v is now %v.",
		settingItems[item].text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorDarkOrange, time.Second*5)
}

func toggleVsync(item int) {
	defer reportPanic("toggleVsync")
	if vSync {
		vSync = false
		settingItems[item].Enabled = false
		ebiten.SetVsyncEnabled(false)
	} else {
		vSync = true
		settingItems[item].Enabled = true
		ebiten.SetVsyncEnabled(true)
	}
	buf := fmt.Sprintf("%v is now %v.",
		settingItems[item].text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorDarkOrange, time.Second*5)
}
