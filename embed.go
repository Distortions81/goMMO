package main

import (
	"embed"
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

	testGrass,
	testlight,
	splashScreen,

	checkOn,
	checkOff,
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

	doLog(true, "Loading sprites.")

	var x, y uint32
	for x = 0; x <= uint32(topSection); x++ {
		if itemTypesList[x] == nil {
			continue
		}
		typeData := itemTypesList[x]

		for y = 0; y <= uint32(topItem); y++ {
			if typeData.items[y] == nil {
				continue
			}
			for s, sprite := range typeData.items[y].sprites {
				doLog(true, "loading '%v:%v'", typeData.items[y].name, sprite.filepath)
				imageData, err := loadSprite(typeData.items[y].name+"/"+sprite.filepath, false)
				if err != nil {
					doLog(true, "loadSprites: %v", err.Error())
					return
				}
				typeData.items[y].sprites[s].image = imageData
			}
		}
	}

	testGrass = findItemImage("ground", "grass", "grass")

	testlight = findItemImage("effects", "light", "light")

	splashScreen = findItemImage("ui", "login", "login")
	checkOn = findItemImage("ui", "check box", "on")
	checkOff = findItemImage("ui", "check box", "off")
	closeBox = findItemImage("ui", "close", "close")
}

func findItemImage(typeName string, itemName string, spriteName string) *ebiten.Image {

	var typeID uint8
	for _, item := range itemTypesList {
		if item == nil {
			continue
		}
		if item.name == typeName {
			typeID = item.id
			break
		}
	}
	iType := itemTypesList[typeID]
	if iType == nil {
		doLog(true, "findItemImage: Type not found: %v:%v:%v", typeName, itemName, spriteName)
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
		doLog(true, "findItemImage: Item not found: %v:%v:%v", typeName, itemName, spriteName)
		return nil
	}

	for _, sprite := range item.sprites {
		if sprite.name == spriteName {
			if sprite.image == nil {
				doLog(true, "findItemImage: Nil match: %v:%v:%v", typeName, itemName, spriteName)
				return nil
			}
			doLog(true, "findItemImage: Matched: %v:%v:%v", typeName, itemName, spriteName)
			return sprite.image
		}
	}

	doLog(true, "findItemImage: Sprite not found: %v:%v:%v", typeName, itemName, spriteName)
	return nil
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
