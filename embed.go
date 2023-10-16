package main

import (
	"embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	//go:embed data
	f embed.FS

	walkNorth *ebiten.Image
	walkEast  *ebiten.Image
	walkSouth *ebiten.Image
	walkWest  *ebiten.Image
)

func loadTest() {
	var err error
	walkNorth, err = getSpriteImage(testCharDir+"Walk Up.png", false)
	if err != nil {
		log.Fatalln(err)
	}
	walkEast, err = getSpriteImage(testCharDir+"Walk Right.png", false)
	if err != nil {
		log.Fatalln(err)
	}
	walkSouth, err = getSpriteImage(testCharDir+"Walk Down.png", false)
	if err != nil {
		log.Fatalln(err)
	}
	walkWest, err = getSpriteImage(testCharDir+"Walk Left.png", false)
	if err != nil {
		log.Fatalln(err)
	}
}

const cLoadEmbedSprites = true

func getSpriteImage(name string, unmanaged bool) (*ebiten.Image, error) {

	if cLoadEmbedSprites {
		gpng, err := f.Open(gfxDir + name)
		if err != nil {
			//DoLog(true, "GetSpriteImage: Embedded: %v", err)
			return nil, err
		}

		m, _, err := image.Decode(gpng)
		if err != nil {
			doLog(true, "GetSpriteImage: Embedded: %v", err)
			return nil, err
		}
		var img *ebiten.Image
		if unmanaged {
			img = ebiten.NewImageFromImageWithOptions(m, &ebiten.NewImageFromImageOptions{Unmanaged: true})
		} else {
			img = ebiten.NewImageFromImage(m)
		}
		return img, nil

	} else {
		img, _, err := ebitenutil.NewImageFromFile(dataDir + gfxDir + name)
		if err != nil {
			doLog(true, "GetSpriteImage: File: %v", err)
		}
		return img, err
	}
}
