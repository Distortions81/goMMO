package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var spritePacks map[string]*spritePack

func initSpritePacks() {
	spritePacks = make(map[string]*spritePack)

	character := spritePack{
		size:    52,
		walking: findItemImage("characters", "player 1", "walk"),
		dead:    findItemImage("characters", "player 1", "dead"),
		attack:  findItemImage("characters", "player 1", "attack"),
		healColors: colorPack{
			frames: []outlineColors{}, numFrames: 3,
		},
	}
	spritePacks["character"] = &character

	zombie := spritePack{
		size:    52,
		walking: findItemImage("creatures", "zombie", "walk"),
		dead:    findItemImage("creatures", "zombie", "dead"),
		attack:  findItemImage("creatures", "zombie", "attack"),
	}
	spritePacks["zombie"] = &zombie
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
