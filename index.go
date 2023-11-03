package main

import (
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const indexFileName = "index.dat"
const assetArraySize = 255

type IID struct {
	section uint8
	num     uint8
}

type sectionData struct {
	id    uint8
	name  string
	items [assetArraySize]*sectionItemData
}

type sectionItemData struct {
	name     string
	fileName string
	itemType uint32
	id       IID
	OnGround bool
	SizeW    uint16
	SizeH    uint16

	image *ebiten.Image
}

var itemTypesList [assetArraySize]*sectionData

func readIndex() bool {

	data, err := efs.ReadFile(dataDir + indexFileName)
	if err != nil {
		doLog(true, "Unable to read %v", indexFileName)
		return false
	}
	doLog(true, "Reading %v", indexFileName)

	lines := strings.Split(string(data), "\n")
	var l int
	var currentSection *sectionData
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
			if !strings.EqualFold("2", words[1]) {
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

		//Section
		if strings.HasSuffix(line, ":") {
			sName := strings.TrimSuffix(line, ":")
			words := strings.Split(sName, ":")
			numWords := len(words)

			if numWords != 2 {
				doLog(true, "Section header invalid: %v words, not 2", numWords)
				return false
			}

			secID, _ := strconv.ParseUint(words[0], 10, 8)
			newSection := &sectionData{name: words[1], id: uint8(secID)}

			itemTypesList[newSection.id] = newSection
			currentSection = newSection

			if devMode {
				doLog(false, "")
				doLog(true, "section found: (%v) %v", newSection.id, newSection.name)
			}
			continue
		}

		//Item data
		if currentSection != nil {
			words := strings.Split(line, ":")
			numWords := len(words)
			if numWords < 2 {
				doLog(true, "Item doesn't have correct number of entries on line %v.", lnum)
				return false
			}
			itemID, _ := strconv.ParseUint(words[0], 10, 8)
			newItem := &sectionItemData{
				name: words[1], fileName: words[2],
				id: IID{section: currentSection.id, num: uint8(itemID)}}
			if numWords == 6 {
				if words[3] == "true" {
					newItem.OnGround = true
				}
				sizeW, _ := strconv.ParseUint(words[4], 10, 16)
				newItem.SizeW = uint16(sizeW)
				sizeH, _ := strconv.ParseUint(words[5], 10, 16)
				newItem.SizeH = uint16(sizeH)
			}
			currentSection.items[newItem.id.num] = newItem

			if devMode {
				doLog(true, "item found: %v:%v", newItem.id, newItem.name)
			}
			continue
		}

	}

	return true
}
