package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const indexFileName = "index.dat"

type sectionData struct {
	name     string
	filePath string
	items    []sectionItemData
}

type sectionItemData struct {
	name     string
	fileName string
	id       uint32
	image    *ebiten.Image
}

var itemTypesList map[string]*sectionData

func readIndex() bool {

	data, err := os.ReadFile(dataDir + indexFileName)
	if err != nil {
		doLog(true, "Unable to read %v", indexFileName)
		return false
	}
	doLog(true, "Reading %v", indexFileName)

	lines := strings.Split(string(data), "\n")
	var l int
	var currentSection *sectionData
	for lnum, line := range lines {

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
			itemTypesList = map[string]*sectionData{}

			if gDevMode {
				doLog(true, "version header found.")
			}
			continue
		}

		//Section
		if strings.HasSuffix(line, ":") {
			sectionName := strings.TrimSuffix(line, ":")
			newSection := &sectionData{name: sectionName}
			itemTypesList[newSection.name] = newSection
			currentSection = newSection

			if gDevMode {
				doLog(false, "")
				doLog(true, "section found: %v", sectionName)
			}
			continue
		}

		//Section data
		if currentSection != nil {
			words := strings.Split(line, ":")
			numWords := len(words)
			if numWords != 3 {
				doLog(true, "Item doesn't have correct number of entries on line %v.", lnum)
				return false
			}
			idNum, _ := strconv.ParseUint(words[0], 10, 32)
			newItem := sectionItemData{name: words[1], fileName: words[2], id: uint32(idNum)}

			for _, item := range currentSection.items {
				if item.id == newItem.id {
					doLog(true, "Duplicate ID! Section: %v, item1: %v, item2: %v id: %v, line: %v", currentSection.name, item.name, newItem.name, newItem.id, lnum)
					return false
				}
			}
			currentSection.items = append(currentSection.items, newItem)

			if gDevMode {
				doLog(true, "item found: %v:%v", newItem.id, newItem.name)
			}
			continue
		}
	}

	return true
}
