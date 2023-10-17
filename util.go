package main

import "image"

func dirToCharOffset(dir int) int {
	switch dir {
	case DIR_S:
		return charSpriteSize * 0
	case DIR_SE:
		return charSpriteSize * 1
	case DIR_E:
		return charSpriteSize * 2
	case DIR_NE:
		return charSpriteSize * 3
	case DIR_N:
		return charSpriteSize * 4
	case DIR_NW:
		return charSpriteSize * 5
	case DIR_W:
		return charSpriteSize * 6
	case DIR_SW:
		return charSpriteSize * 7
	}
	return 0
}

func getCharFrame(dir int) image.Image {

	dirOff := dirToCharOffset(dir)

	rect := image.Rectangle{}
	rect.Min.X = (walkframe * charSpriteSize)
	rect.Max.X = (walkframe * charSpriteSize) + charSpriteSize
	rect.Min.Y = dirOff
	rect.Max.Y = charSpriteSize + dirOff

	return testChar.SubImage(rect)

}
