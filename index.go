package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const assetArraySize = 255

var topSection, topItem uint8

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
	name     string
	fileName string
	id       IID
	OnGround bool
	SizeW    uint16
	SizeH    uint16

	image *ebiten.Image
}

var itemTypesList [assetArraySize]*sectionData

func readDir(path string) bool {

	//ReadDir doesn't like leading slashes
	cleanPath := strings.TrimSuffix(path, "/")

	dirs, err := efs.ReadDir(cleanPath)
	if err != nil {
		doLog(true, "Unable to read directory: %v (%v)", cleanPath, err.Error())
		return false
	}
	for _, item := range dirs {
		buf := fmt.Sprintf("%v/%v", cleanPath, item.Name())

		if item.IsDir() {
			readDir(buf)
		} else if strings.EqualFold(item.Name(), "object.dat") {
			if !readObject(buf) {
				return false
			}
		}
	}

	return true
}

func readObject(filepath string) bool {
	data, err := efs.ReadFile(filepath)
	if err != nil {
		doLog(true, "Unable to read %v", filepath)
		return false
	}
	doLog(true, "Reading %v", filepath)

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
			if secID > uint64(topSection) {
				topSection = uint8(secID)
			}

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
				id: IID{Section: currentSection.id, Num: uint8(itemID)}}
			if itemID > uint64(topItem) {
				topItem = uint8(itemID)
			}
			if numWords == 5 {
				sizeW, _ := strconv.ParseUint(words[3], 10, 16)
				newItem.SizeW = uint16(sizeW)
				sizeH, _ := strconv.ParseUint(words[4], 10, 16)
				newItem.SizeH = uint16(sizeH)
			}
			currentSection.items[newItem.id.Num] = newItem

			if devMode {
				doLog(true, "item found: %v:%v", newItem.id, newItem.name)
			}
			continue
		}

	}

	return true
}
