package main

import (
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const indexFileName = "index.dat"

type sectionData struct {
	id    uint32
	name  string
	items map[uint32]*sectionItemData
}

type sectionItemData struct {
	name     string
	fileName string
	itemType uint32
	id       uint32
	OnGround bool
	SizeW    uint16
	SizeH    uint16

	image *ebiten.Image
}

var itemTypesList map[uint32]*sectionData

func readIndex() bool {

	var sectionID uint32
	var itemID uint32

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
			if !strings.EqualFold("1", words[1]) {
				doLog(true, "Index version not supported, line: %v", lnum)
				return false
			}

			//Reset data
			currentSection = nil
			itemTypesList = map[uint32]*sectionData{}

			if devMode {
				doLog(true, "version header found.")
			}
			continue
		}

		//Section
		if strings.HasSuffix(line, ":") {
			itemID = 0
			sName := strings.TrimSuffix(line, ":")
			newSection := &sectionData{name: sName, id: uint32(sectionID)}
			sectionID++

			itemTypesList[newSection.id] = newSection
			currentSection = newSection

			if itemTypesList[newSection.id].items == nil {
				itemTypesList[newSection.id].items = make(map[uint32]*sectionItemData)
			}

			if devMode {
				doLog(false, "")
				doLog(true, "section found: (%v) %v", newSection.id, newSection.name)
			}
			continue
		}

		//Section data
		if currentSection != nil {
			words := strings.Split(line, ":")
			numWords := len(words)
			if numWords < 2 {
				doLog(true, "Item doesn't have correct number of entries on line %v.", lnum)
				return false
			}
			newItem := &sectionItemData{
				name: words[0], fileName: words[1],
				id: uint32(itemID), itemType: currentSection.id}
			if numWords == 5 {
				if words[2] == "true" {
					newItem.OnGround = true
				}
				sizeW, _ := strconv.ParseUint(words[3], 10, 16)
				newItem.SizeW = uint16(sizeW)
				sizeH, _ := strconv.ParseUint(words[4], 10, 16)
				newItem.SizeH = uint16(sizeH)
			}
			itemID++
			currentSection.items[newItem.id] = newItem

			if devMode {
				doLog(true, "item found: %v:%v", newItem.id, newItem.name)
			}
			continue
		}

	}

	return true
}
