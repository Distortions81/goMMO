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

var topSection, topItem uint8

const (
	OBJDATA_NONE CMD = iota
	OBJDAT_INFO
	OBJDAT_SPRITES
)

type IID struct {
	Section uint8
	Num     uint8
	UID     uint32
}

type sectionData struct {
	id    uint8
	name  string
	items [assetArraySize]*sectionItemData
}

type sectionItemData struct {
	name string

	id       IID
	OnGround bool
	SizeW    uint16
	SizeH    uint16

	sprites []spriteData
}

type spriteData struct {
	name  string
	image *ebiten.Image
}

var itemTypesList [assetArraySize]*sectionData
var currentSection *sectionData

func readIndex() bool {

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

		if numWords != 2 {
			doLog(true, "Section header invalid: %v words, not 2", numWords)
			return false
		}

		secID, _ := strconv.ParseUint(words[0], 10, 8)
		newSection := &sectionData{name: words[1], id: uint8(secID)}
		if secID > uint64(topSection) {
			topSection = uint8(secID)
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

	return false
}

func readObjects(section string) bool {
	dirs, err := efs.ReadDir(gfxDir + section)
	if err != nil {
		doLog(true, "Unable to read directory: %v", section)
	}

	for _, item := range dirs {
		if item.IsDir() {
			readObject(section + "/" + item.Name())
		}
	}
	return false
}

func readObject(name string) bool {
	if currentSection == nil {
		doLog(true, "ReadObject: No valid current section?")
		return false
	}

	filePath := fmt.Sprintf("%v%v/object.dat", gfxDir, name)
	data, err := efs.ReadFile(filePath)
	if err != nil {
		doLog(true, "Unable to read %v", filePath)
		return false
	}
	doLog(true, "Reading %v", filePath)

	var xs, ys, id uint64
	var spriteName, spriteFile string

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
		if area == OBJDATA_NONE {
			if words[0] == "info" {
				area = OBJDAT_INFO
			} else if words[0] == "sprites" {
				area = OBJDAT_SPRITES
			}
		} else if area == OBJDAT_INFO {
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
			spriteName = words[0]
			spriteFile = words[1]
		}
	}

	itemID := IID{Section: currentSection.id, Num: uint8(id)}
	newItem := &sectionItemData{
		name:  spriteName,
		SizeW: uint16(ys), SizeH: uint16(xs),
		id: itemID,
	}
	currentSection.items[newItem.id.Num] = newItem

	if devMode {
		doLog(true, "item found: %v (%v)", newItem.name, newItem.id)
	}

	return true
}
