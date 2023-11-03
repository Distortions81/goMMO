package main

import (
	"embed"
	"fmt"
	"image"
	"io"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	//go:embed data
	efs embed.FS

	playerSprite *ebiten.Image
	testGrass    *ebiten.Image
	testlight    *ebiten.Image
	splashScreen *ebiten.Image

	checkOn  *ebiten.Image
	checkOff *ebiten.Image
	closeBox *ebiten.Image
)

// Read text files
func getText(name string) (string, error) {
	file, err := efs.Open(txtDir + name + ".txt")
	if err != nil {
		doLog(true, "GetText: %v", err)
		return "GetText: File: " + name + " not found in embed.", err
	}

	txt, err := io.ReadAll(file)
	if err != nil {
		doLog(true, "GetText: %v", err)
		return "Error: Failed read: " + name, err
	}

	if len(txt) > 0 {
		doLog(true, "GetText: %v", name)
		return strings.ReplaceAll(string(txt), "\r", ""), nil
	} else {
		return "Error: length 0!", err
	}

}

// Load sprites
func loadSprites() {

	var x, y uint32
	for x = 0; x < uint32(topSection); x++ {
		if itemTypesList[x] == nil {
			continue
		}
		typeData := itemTypesList[x]

		for y = 0; y < uint32(topItem); y++ {
			if typeData.items[y] == nil {
				continue
			}
			itemData := typeData.items[y]

			doLog(true, "loading %v:%v", typeData.name, itemData.fileName)
			imageData, err := loadSprite(typeData.name+"/"+itemData.fileName, false)
			if err != nil {
				log.Fatalln(err)
			}
			itemTypesList[x].items[y].image = imageData
		}
	}

	testGrass = findItemImage("ground", "grass-1")
	playerSprite = findItemImage("characters", "player")
	testlight = findItemImage("effects", "light")
	splashScreen = findItemImage("ui", "login")

	checkOn = findItemImage("ui", "check on")
	checkOff = findItemImage("ui", "check off")
	closeBox = findItemImage("ui", "close box")
}

func findItemImage(typeName string, itemName string) *ebiten.Image {

	var typeID uint8
	for _, item := range itemTypesList {
		if item == nil {
			continue
		}
		fmt.Println(item.name)
		if item.name == typeName {
			typeID = item.id
			continue
		}
	}
	iType := itemTypesList[typeID]
	if iType == nil {
		doLog(true, "Item type not found: %v", typeName)
		return nil
	}

	var itemID uint8
	for _, item := range iType.items {
		if item == nil {
			continue
		}
		if item.name == itemName {
			itemID = item.id.num
			continue
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

// Read fonts
func getFont(name string) []byte {
	data, err := efs.ReadFile(gfxDir + "fonts/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return data

}

// Load sprites
func loadSprite(name string, unmanaged bool) (*ebiten.Image, error) {

	if loadEmbedSprites {
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
