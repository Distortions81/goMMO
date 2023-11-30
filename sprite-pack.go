package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var spritePacks map[string]*spritePack

var healAnimation = colorPack{
	frames: []outlineColors{
		{
			colors: []color.RGBA{
				{R: 1, G: 0, B: 0, A: 1},
				{R: 0, G: 1, B: 0, A: 1},
				{R: 0, G: 0, B: 1, A: 1},
			},
			outlineWidth: []int{2},
		},
		{
			colors: []color.RGBA{
				{R: 0, G: 0, B: 1, A: 1},
				{R: 0, G: 1, B: 0, A: 1},
				{R: 1, G: 0, B: 0, A: 1},
			},
			outlineWidth: []int{2},
		},
	},
	numFrames: 2,
}

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

				healSprite := makeOutlines(walkSprite, healAnimation.frames[0].colors)
				healDeadSprite := makeOutlines(deadSprite, healAnimation.frames[0].colors)
				healAttackSprite := makeOutlines(attackSprite, healAnimation.frames[0].colors)

				if walkSprite == nil || deadSprite == nil || attackSprite == nil {
					doLog(true, "Item not found: %v, %v", itemType.name, item.name)
					continue
				}
				newPack := &spritePack{
					walking: walkSprite, dead: deadSprite, attack: attackSprite,
					healing:       []*ebiten.Image{healSprite},
					healingDead:   []*ebiten.Image{healDeadSprite},
					healingAttack: []*ebiten.Image{healAttackSprite},
					sizeH:         int(item.sizeH), sizeW: int(item.sizeW)}
				spritePacks[item.name] = newPack
			}
		}
	}
}

type spritePack struct {
	sizeH, sizeW int

	walking *ebiten.Image
	dead    *ebiten.Image
	attack  *ebiten.Image

	healing       []*ebiten.Image
	healingDead   []*ebiten.Image
	healingAttack []*ebiten.Image

	healColors colorPack
}
