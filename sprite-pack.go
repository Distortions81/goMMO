package main

import "github.com/hajimehoshi/ebiten/v2"

var spritePacks map[string]*spritePack

func initSpritePacks() {
	spritePacks = make(map[string]*spritePack)

	character := spritePack{
		size:    52,
		walking: findItemImage("characters", "player"),
		dead:    findItemImage("characters", "player-dead"),
		attack:  findItemImage("characters", "player-attacj"),

		healing:  findItemImage("characters", "player-heal"),
		healing2: findItemImage("characters", "player-heal2"),

		healingDead:  findItemImage("characters", "player-dead-heal"),
		healingDead2: findItemImage("characters", "player-dead-heal2"),

		healingAttack:  findItemImage("characters", "player-attack-heal"),
		healingAttack2: findItemImage("characters", "player-attack-heal2"),
	}
	spritePacks["character"] = &character
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
}
