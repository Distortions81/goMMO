package main

import (
	"os"
	"strconv"
	"strings"
)

const indexFileName = "index.dat"

type sectionData struct {
	name     string
	items    []sectionItemData
	numItems uint32
}

type sectionItemData struct {
	name     string
	fileName string
	id       uint32
}

var sections []*sectionData

func readIndex() bool {
	data, err := os.ReadFile(dataDir + indexFileName)
	if err != nil {
		doLog(true, "Unable to read index.dat")
		return false
	}
	doLog(true, "Reading %v", indexFileName)

	lines := strings.Split(string(data), "\n")
	var l int
	var currentSection *sectionData
	for _, line := range lines {

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
				doLog(true, "version header doesn't have two words.")
				return false
			}
			if words[0] != "version" {
				doLog(true, "no version header found.")
				return false
			}
			if words[1] != "1" {
				doLog(true, "index version not supported.")
				return false
			}

			//Reset data
			sections = []*sectionData{}
			currentSection = nil
			if gDevMode {
				doLog(true, "version header found.")
			}
			continue
		}

		//Section
		if strings.HasSuffix(line, ":") {
			sectionName := strings.TrimSuffix(line, ":")
			newSection := &sectionData{name: sectionName}
			sections = append(sections, newSection)
			currentSection = newSection

			if gDevMode {
				doLog(true, "section found: %v", sectionName)
			}
			continue
		}

		//Section data
		if currentSection != nil {
			words := strings.Split(line, ":")
			numWords := len(words)
			if numWords != 3 {
				doLog(true, "Item doesn't have correct number of entries.")
				return false
			}
			idNum, _ := strconv.ParseUint(words[0], 10, 32)
			newItem := sectionItemData{name: words[1], fileName: words[2], id: uint32(idNum)}

			for _, item := range currentSection.items {
				if item.id == newItem.id {
					doLog(true, "Duplicate ID! Section: %v, item1: %v, item2: %v id: %v", currentSection.name, item.name, newItem.name, newItem.id)
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
