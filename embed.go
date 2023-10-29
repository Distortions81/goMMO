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

	spritelist []*sectionItemData
	numSprites uint32
)

func loadTest() {

	for typeid, typeData := range itemTypesList {
		for itemid, itemData := range typeData.items {
			imageData, err := getSpriteImage(typeData.name+"/"+itemData.fileName, false)
			if err != nil {
				log.Fatalln(err)
			}
			itemTypesList[typeid].items[itemid].image = imageData
			spritelist = append(spritelist, itemTypesList[typeid].items[itemid])
			numSprites++
		}
	}

	testGrass = getItemImage("terrain", "grass-1")
	testChar = getItemImage("characters", "player")
	testlight = getItemImage("effects", "light")
}

func getItemImage(typeName string, itemName string) *ebiten.Image {

	var typeID uint32
	for itemid, item := range itemTypesList {
		if item.name == typeName {
			typeID = itemid
			break
		}
	}
	iType := itemTypesList[typeID]
	if iType == nil {
		doLog(true, "Item type not found: %v", typeName)
		return nil
	}

	var itemID uint32
	for itemid, item := range iType.items {
		if item.name == itemName {
			itemID = itemid
			break
		}
	}
	item := iType.items[itemID]
	if item == nil {
		doLog(true, "Item not found: %v : %v", typeName, itemName)
		return nil
	}
	if item.image == nil {
		doLog(true, "Image not found: %v : %v", typeName, itemName)
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
