package main

/* Figure out what option item user clicked */
func handleOptions(input XYs, window *windowData) bool {
	defer reportPanic("handleOptions")
	windowsLock.Lock()
	defer windowsLock.Unlock()

	originX := window.position.X
	originY := window.position.Y

	for i, item := range settingItems {
		b := optionWindowButtons[i]
		if PosWithinRect(
			XY{X: uint32(input.X - originX),
				Y: uint32(input.Y - originY)}, b, 1) {
			if (WASMMode && !item.WASMExclude) || !WASMMode {
				item.action(i)
				saveOptions()
				window.dirty = true
				mouseHeld = false

				return true
			}
		}
	}

	return false
}

func handleHelpWindow(input XYs, window *windowData) bool {
	defer reportPanic("handleHelpWindow")
	windowsLock.Lock()
	defer windowsLock.Unlock()

	if !mouseHeld {
		return false
	}

	originX := window.position.X
	originY := window.position.Y

	for i := range updateWindowButtons {
		b := updateWindowButtons[i]
		if PosWithinRect(
			XY{X: uint32(input.X - originX),
				Y: uint32(input.Y - originY)}, b, 1) {
			return true
		}
	}

	return false
}
