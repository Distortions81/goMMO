package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const fpx = 96.0

var (
	fontDPI float64 = fpx
	uiScale         = 1.0
	magnify         = true

	/* Fonts */
	bootFont  font.Face
	bootFontH int

	toolTipFont  font.Face
	toolTipFontH int

	monoFont  font.Face
	monoFontH int

	logoFont  font.Face
	logoFontH int

	introFont  font.Face
	introFontH int

	generalFont  font.Face
	generalFontH int

	largeGeneralFont  font.Face
	largeGeneralFontH int

	objectFont  font.Face
	objectFontH int
)

func updateFonts() {
	defer reportPanic("updateFonts")

	newVal := fpx * uiScale
	if newVal < 1 {
		newVal = 1
	}
	fontDPI = newVal

	var mono, tt *opentype.Font
	var err error

	fontData := getFont("Exo2-Regular.ttf")
	collection, err := opentype.ParseCollection(fontData)
	if err != nil {
		log.Fatal(err)
	}

	tt, err = collection.Font(0)
	if err != nil {
		log.Fatal(err)
	}

	/* Mono font */
	fontData = getFont("Hack-Regular.ttf")
	collection, err = opentype.ParseCollection(fontData)
	if err != nil {
		log.Fatal(err)
	}

	mono, err = collection.Font(0)
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * Font DPI
	 * Changes how large the font is for a given point value
	 */

	/* General font */
	generalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    10,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	generalFontH = getFontHeight(generalFont)

	/* Large general font */
	largeGeneralFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	largeGeneralFontH = getFontHeight(largeGeneralFont)

	/* Missing texture font */
	objectFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    6,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	objectFontH = getFontHeight(objectFont)

	/* Tooltip font */
	toolTipFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	toolTipFontH = getFontHeight(toolTipFont)

	/* Mono font */
	monoFont, err = opentype.NewFace(mono, &opentype.FaceOptions{
		Size:    8,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	monoFontH = getFontHeight(monoFont)

}

const sizingText = "!@#$%^&*()_+-=[]{}|;':,.<>?`~qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

func getFontHeight(font font.Face) int {
	defer reportPanic("getFontHeight")
	tRect := text.BoundString(font, sizingText)
	return tRect.Dy()
}
