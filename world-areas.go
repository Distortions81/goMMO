package main

type areaData struct {
	version uint16

	wObjects []wObjectData
	terrain  uint8
	effect   uint8
}

type wObjectData struct {
	id  uint32
	pos XY
}
