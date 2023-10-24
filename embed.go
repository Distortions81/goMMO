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
	efs embed.FS

	testChar  *ebiten.Image
	testGrass *ebiten.Image
	testlight *ebiten.Image
)

func loadTest() {

	for typeName, itemType := range itemTypesList {
		for itemName, item := range itemType.items {
			imageData, err := getSpriteImage(itemType.name+"/"+item.fileName, false)
			if err != nil {
				log.Fatalln(err)
			}
			itemTypesList[typeName].items[itemName].image = imageData
		}
	}

	testGrass = getItemImage("terrain", "grass-1")
	testChar = getItemImage("characters", "player")
	testlight = getItemImage("effects", "light")
}

func getItemImage(itemType string, name string) *ebiten.Image {
	iType := itemTypesList[itemType]
	if iType == nil {
		doLog(true, "Item type not found: %v", itemType)
		return nil
	}
	item := iType.items[name]
	if item == nil {
		doLog(true, "Item not found: %v", name)
		return nil
	}
	if item.image == nil {
		doLog(true, "Item has no image: %v", name)
		return nil
	}

	return item.image
}

func getFont(name string) []byte {
	data, err := efs.ReadFile(gfxDir + "fonts/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return data

}

func getSpriteImage(name string, unmanaged bool) (*ebiten.Image, error) {

	if cLoadEmbedSprites {
		gpng, err := efs.Open(gfxDir + name)
		if err != nil {
			doLog(true, "GetSpriteImage: Embedded: %v", err)
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
