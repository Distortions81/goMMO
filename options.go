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
	Text       string `json:"-"`

	TextPosX   int             `json:"-"`
	TextPosY   int             `json:"-"`
	TextBounds image.Rectangle `json:"-"`
	Rect       image.Rectangle `json:"-"`

	Enabled     bool
	WASMExclude bool

	action  func(item int) `json:"-"`
	NoCheck bool           `json:"-"`
}

func init() {
	defer reportPanic("options init")
	settingItems = []settingType{
		{ConfigName: "VSYNC", Text: "Limit FPS (VSYNC)", action: toggleVsync, Enabled: true, WASMExclude: true},
		{ConfigName: "FULLSCREEN", Text: "Full Screen", action: toggleFullscreen},
		{ConfigName: "FAST-SHADOWS", Text: "Fast Shadows", action: toggleFastShadow},
		{ConfigName: "MOTION-SMOOTH", Text: "Motion Smoothing", action: toggleSmoothing, Enabled: true},
		{ConfigName: "NIGHT-MODE", Text: "Disable Shadows", action: toggleNightShadow},
		{ConfigName: "DEBUG-TEXT", Text: "Debug info-text", action: toggleInfoLine},
	}
}

/* Load user options settings from disk */
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

/* Save user options settings to disk */
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

/* Toggle the debug bottom-screen text */
func toggleInfoLine(item int) {
	defer reportPanic("toggleInfoLine")
	if infoLine {
		infoLine = false
		settingItems[item].Enabled = false
	} else {
		infoLine = true
		settingItems[item].Enabled = true
	}
}

/* Toggle the use of hyper-threading */
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

/* Toggle units */
func toggleUnits(item int) {
	defer reportPanic("toggleUnits")
	if usUnits {
		usUnits = false
		settingItems[item].Enabled = false
	} else {
		usUnits = true
		settingItems[item].Enabled = true
	}
}

/* Toggle full-screen */
func toggleFullscreen(item int) {
	defer reportPanic("toggleFullscreen")
	if ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
		settingItems[item].Enabled = false
	} else {
		ebiten.SetFullscreen(true)
		settingItems[item].Enabled = true
	}
	buf := fmt.Sprintf("%v is now %v.",
		settingItems[item].Text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorOrange, time.Second*5)
}

/* Toggle UI magnification */
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
		settingItems[item].Text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorOrange, time.Second*5)
}

/* Toggle debug mode */
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
		settingItems[item].Text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorOrange, time.Second*5)
}

/* Toggle autosave */
func toggleAutosave(item int) {
	defer reportPanic("toggleAutosave")
	if autoSave {
		autoSave = false
		settingItems[item].Enabled = false
	} else {
		autoSave = true
		settingItems[item].Enabled = true
	}
	buf := fmt.Sprintf("%v is now %v.",
		settingItems[item].Text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorOrange, time.Second*5)
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
		settingItems[item].Text,
		BoolToOnOff(settingItems[item].Enabled))
	chatDetailed(buf, ColorOrange, time.Second*5)
}
