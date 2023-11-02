package main

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	toolbarCache     *ebiten.Image
	toolbarCacheLock sync.RWMutex
	toolbarMax       int
	toolbarItems           = []toolbarItemData{}
	selectedItemType uint8 = maxItemType

	toolbarHover bool
)

const (
	/* Game data structures */
	/* Subtypes */
	objSubUI   = 0
	objSubGame = 1
	objOverlay = 2

	maxItemType = 255

	/* Toolbar settings */
	toolBarIconSize   = 32
	toolBarSpaceRatio = 4
	tbSelThick        = 2
	halfSelThick      = tbSelThick / 2
)

/* Toolbar list item */
type toolbarItemData struct {
	sType int
	oType *objTypeData
}

/* Object type data, includes image, toolbar action, and update handler */
type objTypeData struct {
	base        string
	name        string
	description string

	image *ebiten.Image

	/* Toolbar Specific */
	excludeWASM bool       //Don't show this object in the toolbar on WASM
	qKey        ebiten.Key //Toolbar quick-key

	/* Function links */
	toolbarAction func() `json:"-"`
}

type subTypeData struct {
	folder string
	list   []*objTypeData
}

/* Toolbar item types, array of array of ObjType */
var subTypes = []subTypeData{
	{
		folder: "ui",
		list:   uiObjs,
	},
}

/* Toolbar actions and images */
var uiObjs = []*objTypeData{
	//Ui Only
	{
		base: "settings",
		name: "Options", toolbarAction: settingsToggle,
		description: "Show game options", qKey: ebiten.KeyF2,
	},
	{
		base: "help",
		name: "Help", toolbarAction: toggleHelp,
		description: "See game controls and help.", qKey: ebiten.KeyF1,
	},
}

/* Make default toolbar list */
func initToolbar() {
	defer reportPanic("InitToolbar")
	toolbarMax = 0
	for subPos, subType := range subTypes {
		for o, oType := range subType.list {
			/* Skips some items for WASM */
			if WASMMode && oType.excludeWASM {
				continue
			}
			toolbarMax++
			toolbarItems = append(toolbarItems, toolbarItemData{sType: subPos, oType: oType})

			subType.list[o].image = getItemImage(subType.folder, oType.base)
		}
	}
}

/* Draw toolbar to an image */
func drawToolbar(click, hover bool, index int) {
	defer reportPanic("drawToolbar")
	iconSize := float32(uiScale * toolBarIconSize)
	spacing := float32(toolBarIconSize / toolBarSpaceRatio)

	toolbarCacheLock.Lock()
	defer toolbarCacheLock.Unlock()

	/* If needed, init image */
	if toolbarCache == nil {
		toolbarCache = ebiten.NewImage(int(iconSize+spacing)*toolbarMax+4, int(iconSize+spacing))
	}
	/* Clear, full with semi-transparent */
	toolbarCache.Clear()
	toolbarCache.Fill(ColorToolTipBG)

	/* Loop through all toolbar items */
	for pos := 0; pos < toolbarMax; pos++ {
		item := toolbarItems[pos]

		/* Get main image */
		img := item.oType.image

		/* Something went wrong, exit */
		if img == nil {
			doLog(false, "FAILURE: %v\n", pos)
			return
		}

		var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}

		op.GeoM.Reset()
		iSize := img.Bounds()

		/* Handle non-square sprites */
		/* Just make toolbar sprites instead */
		var largerDim int
		if iSize.Size().X > largerDim {
			largerDim = iSize.Size().X
		}
		if iSize.Size().Y > largerDim {
			largerDim = iSize.Size().Y
		}

		/* Adjust image to toolbar size */
		op.GeoM.Scale(
			uiScale/(float64(largerDim)/float64(toolBarIconSize)),
			uiScale/(float64(largerDim)/float64(toolBarIconSize)))

		/* Move to correct location in toolbar image */
		op.GeoM.Translate((float64(iconSize+(spacing))*float64(pos))+float64(spacing/2), float64(spacing/2))

		/* hovered/clicked icon highlight */
		if pos == index {
			if click {
				vector.DrawFilledRect(toolbarCache, float32(pos)*(iconSize+spacing),
					0, iconSize+spacing, iconSize+spacing, ColorRed, false)
				toolbarHover = true

				go func() {
					time.Sleep(time.Millisecond * 155)
					drawToolbar(false, false, 0)
				}()
			} else if hover {
				vector.DrawFilledRect(toolbarCache, float32(pos)*(iconSize+spacing),
					0, iconSize+spacing, iconSize+spacing, ColorAqua, false)
				toolbarHover = true
			}

		}

		/* Draw to image */
		toolbarCache.DrawImage(img, op)
	}
}

/* Handle clicks that end up within the toolbar */
func handleToolbar() bool {
	defer reportPanic("handleToolbar")

	iconSize := float32(uiScale * toolBarIconSize)
	spacing := float32(toolBarIconSize / toolBarSpaceRatio)

	tbLength := float32((toolbarMax * int(iconSize+spacing)))

	fmx := float32(MouseX)
	fmy := float32(MouseY)

	/* If the click isn't off the right of the toolbar */
	if fmx <= tbLength {
		/* If the click isn't below the toolbar */
		if fmy <= iconSize {

			tbItem := int(fmx / float32(iconSize+spacing))
			len := len(toolbarItems) - 1
			if tbItem > len {
				tbItem = len
			} else if tbItem < 0 {
				tbItem = 0
			}
			item := toolbarItems[tbItem].oType

			/* Draw item hover */
			drawToolbar(true, false, tbItem)

			/* Actions */
			if item.toolbarAction != nil {
				item.toolbarAction()
				drawToolbar(true, false, tbItem)
			}

			/* Eat this mouse event */
			gMouseHeld = false
			gClickCaptured = true
			return true
		}
	}
	return false
}
