package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var spritePacks map[string]*spritePack

func initSpritePacks() {
	spritePacks = make(map[string]*spritePack)

	for _, itemType := range itemTypesList {

		if itemType.name == "characters" ||
			itemType.name == "creatures" {
			doLog(true, "Creating sprite pack: %v", itemType.name)
			for _, item := range itemType.items {
				walkSprite := findItemImage(itemType.name, item.name, "walk")
				deadSprite := findItemImage(itemType.name, item.name, "dead")
				attackSprite := findItemImage(itemType.name, item.name, "attack")
				if walkSprite == nil || deadSprite == nil || attackSprite == nil {
					doLog(true, "Item not found: %v, %v", itemType.name, item.name)
					continue
				}
				newPack := &spritePack{walking: walkSprite, dead: deadSprite, attack: attackSprite}
				spritePacks[item.name] = newPack
			}
		}
	}
}

type spritePack struct {
	size int

	walking *ebiten.Image
	dead    *ebiten.Image
	attack  *ebiten.Image

	healing  *ebiten.Image
	healing2 *ebiten.Image

	healingDead  *ebiten.Image
	healingDead2 *ebiten.Image

	healingAttack  *ebiten.Image
	healingAttack2 *ebiten.Image

	healColors colorPack
}
