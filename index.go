package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	assetArraySize = 255
	indexPath      = gfxDir + "index.dat"
)

const (
	OBJDATA_NONE CMD = iota
	OBJDAT_INFO
	OBJDAT_SPRITES
)

type IID struct {
	section uint8
	num     uint8
	sprite  uint8
}

type sectionData struct {
	id       uint8
	name     string
	items    map[uint8]*sectionItemData
	onGround bool
}

type sectionItemData struct {
	name string

	id       IID
	onGround bool
	sizeW    uint16
	sizeH    uint16

	sprites map[uint8]*spriteData
}

type spriteData struct {
	name     string
	filepath string
	id       uint8
	image    *ebiten.Image
}

var itemTypesList map[uint8]*sectionData
var currentSection *sectionData

func readIndex() bool {
	itemTypesList = make(map[uint8]*sectionData)

	data, err := efs.ReadFile(indexPath)
	if err != nil {
		doLog(true, "Unable to read %v", indexPath)
		return false
	}
	doLog(true, "Reading %v", indexPath)

	lines := strings.Split(string(data), "\n")
	var l int

	for ln, line := range lines {
		lnum := ln + 1

		//Ignore comments and blank lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		l++

		//Version header
		if l == 1 {
			words := strings.Split(line, " ")
			numWords := len(words)
			if numWords != 2 {
				doLog(true, "Version header doesn't have two words, line: %v", lnum)
				return false
			}
			if !strings.EqualFold("version", words[0]) {
				doLog(true, "No version header found, line: %v: '%v %v'", lnum, words[0], words[1])
				return false
			}
			if !strings.EqualFold("3", words[1]) {
				doLog(true, "Index version not supported, line: %v", lnum)
				return false
			}

			//Reset data
			currentSection = nil

			if devMode {
				doLog(true, "version header found.")
			}
			continue
		}

		sName := strings.TrimSuffix(line, ":")
		words := strings.Split(sName, ":")
		numWords := len(words)

		if numWords < 2 {
			doLog(true, "Section header invalid: %v words, less than 2.", numWords)
			return false
		}

		secID, _ := strconv.ParseUint(words[0], 10, 8)
		iMap := make(map[uint8]*sectionItemData)
		newSection := &sectionData{name: words[1], id: uint8(secID), items: iMap}

		if numWords > 3 {
			if strings.EqualFold(words[3], "onground") {
				newSection.onGround = true
			}
		}

		itemTypesList[newSection.id] = newSection
		currentSection = newSection

		if devMode {
			doLog(false, "")
			doLog(true, "section found: (%v) %v", newSection.id, newSection.name)
		}

		readObjects(words[1])
		continue

	}

	return true
}

func readObjects(section string) bool {
	dirs, err := efs.ReadDir(gfxDir + section)
	if err != nil {
		doLog(true, "Unable to read directory: %v", section)
	}

	for _, item := range dirs {
		if item.IsDir() {
			readObject(section, item.Name())
		}
	}
	return true
}

func readObject(folder, name string) bool {
	if currentSection == nil {
		doLog(true, "ReadObject: No valid current section?")
		return false
	}

	filePath := fmt.Sprintf("%v%v/%v/object.dat", gfxDir, folder, name)
	data, err := efs.ReadFile(filePath)
	if err != nil {
		doLog(true, "Unable to read %v", filePath)
		return false
	}
	doLog(true, "Reading %v", filePath)

	var xs, ys, id uint64
	var lid uint8
	var sdata []spriteData

	lines := strings.Split(string(data), "\n")

	area := OBJDATA_NONE

	var l int
	for ln, line := range lines {
		lnum := ln + 1

		//Ignore comments and blank lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		l++

		//Version header
		if l == 1 {
			words := strings.Split(line, " ")
			numWords := len(words)
			if numWords != 2 {
				doLog(true, "Version header doesn't have two words, line: %v", lnum)
				return false
			}
			if !strings.EqualFold("version", words[0]) {
				doLog(true, "No version header found, line: %v: '%v %v'", lnum, words[0], words[1])
				return false
			}
			if !strings.EqualFold("3", words[1]) {
				doLog(true, "Index version not supported, line: %v", lnum)
				return false
			}

			if devMode {
				doLog(true, "version header found.")
			}
			continue
		}

		words := strings.Split(line, ":")
		numWords := len(words)

		if numWords == 0 {
			continue
		}

		if words[0] == "info" {
			area = OBJDAT_INFO
			continue
		} else if words[0] == "sprites" {
			area = OBJDAT_SPRITES
			continue
		}
		if area == OBJDAT_INFO {
			if words[0] == "size" {
				dims := strings.Split(words[1], ",")
				numDims := len(dims)
				if numDims == 2 {
					xs, _ = strconv.ParseUint(dims[0], 10, 32)
					ys, _ = strconv.ParseUint(dims[1], 10, 32)
				} else if numDims == 1 {
					xs, _ = strconv.ParseUint(dims[0], 10, 32)
					ys = xs
				} else {
					doLog(true, "Invalid number of size dimensions: %v, line: %v", numDims, lnum)
				}
			} else if words[0] == "id" {
				id, _ = strconv.ParseUint(words[1], 10, 32)
			}
		} else if area == OBJDAT_SPRITES {
			sdata = append(sdata, spriteData{name: words[0], filepath: words[1], id: lid})
			lid++
			doLog(true, "sprite added: %v: %v", words[0], words[1])
		}
	}

	itemID := IID{section: currentSection.id, num: uint8(id)}
	newItem := &sectionItemData{
		name:  name,
		sizeW: uint16(ys), sizeH: uint16(xs),
		id:      itemID,
		sprites: make(map[uint8]*spriteData),
	}
	for s, sprite := range sdata {
		newItem.sprites[sprite.id] = &sdata[s]
	}
	currentSection.items[newItem.id.num] = newItem

	if devMode {
		doLog(true, "object added: %v", newItem.name)
	}

	return true
}
