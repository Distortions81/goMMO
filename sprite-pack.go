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
				{0, 255, 255, 255},
			},
			outlineWidth: []float64{
				1.25, 1.25, 1.25, 1.25, 1.25,
			},
		},
		{
			colors: []color.RGBA{
				{0, 255, 255, 255},
				{0, 255, 255, 255},
			},
			outlineWidth: []float64{
				1.25, 1.25, 1.25, 1.25, 1.25,
			},
		},
		{
			colors: []color.RGBA{
				{0, 255, 255, 255},
				{0, 255, 255, 255},
				{0, 192, 255, 255},
			},
			outlineWidth: []float64{
				1.25, 1.25, 1.25, 1.25, 1.25,
			},
		},
		{
			colors: []color.RGBA{
				{0, 255, 255, 255},
				{0, 255, 255, 255},
				{0, 192, 255, 255},
				{0, 128, 255, 255},
			},
			outlineWidth: []float64{
				1.25, 1.25, 1.25, 1.25, 1.25,
			},
		},
		{
			colors: []color.RGBA{
				{0, 255, 255, 255},
				{0, 255, 255, 255},
				{0, 192, 255, 255},
				{0, 128, 255, 255},
				{0, 0, 255, 255},
			},
			outlineWidth: []float64{
				1.25, 1.25, 1.25, 1.25, 1.25,
			},
		},
	},
}

func initSpritePacks() {
	spritePacks = make(map[string]*spritePack)

	healAnimation.numFrames = len(healAnimation.frames)

	for _, itemType := range itemTypesList {

		if itemType.name == "characters" ||
			itemType.name == "creatures" {
			doLog(true, "Creating sprite pack: %v", itemType.name)
			for _, item := range itemType.items {
				walkSprite := findItemImage(itemType.name, item.name, "walk")
				deadSprite := findItemImage(itemType.name, item.name, "dead")
				attackSprite := findItemImage(itemType.name, item.name, "attack")

				var healSprite, healDeadSprite, healAttackSprite []*ebiten.Image
				for x := 0; x < healAnimation.numFrames; x++ {
					healSprite = append(healSprite, makeOutlines(walkSprite, healAnimation.frames[x].colors))
					healDeadSprite = append(healDeadSprite, makeOutlines(deadSprite, healAnimation.frames[x].colors))
					healAttackSprite = append(healAttackSprite, makeOutlines(attackSprite, healAnimation.frames[x].colors))
				}

				if walkSprite == nil || deadSprite == nil || attackSprite == nil {
					doLog(true, "Item not found: %v, %v", itemType.name, item.name)
					continue
				}
				newPack := &spritePack{
					walking: walkSprite, dead: deadSprite, attack: attackSprite,
					healing:       healSprite,
					healingDead:   healDeadSprite,
					healingAttack: healAttackSprite,
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
