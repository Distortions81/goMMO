package main

// Figure out what option item user clicked
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
			if (WASMMode && !item.wasmExclude) || !WASMMode {
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

func handleLoginWindow(input XYs, window *windowData) bool {
	defer reportPanic("handleLoginWindow")
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

func handleLoginKeys(input []rune, window *windowData, delete, enter bool) bool {

	ilen := len(input)
	window.dirty = true
	if login.selected == SELECTED_LOGIN {
		if enter {
			login.selected = SELECTED_PASSWORD
			return true
		}
		llen := len(login.login)

		if delete {
			if llen > 0 {
				login.login = login.login[:llen-1]
				return true
			}
		}

		if ilen > 0 {
			login.login = append(login.login, input[0])
			return true
		}
	} else if login.selected == SELECTED_PASSWORD {
		if enter {
			login.selected = SELECTED_GO
			return true
		}
		if delete {
			llen := len(login.password)
			if llen > 0 {
				login.login = login.password[:llen-1]
				return true
			}
		}

		if ilen > 0 {
			login.password = append(login.password, input[0])
			return true
		}
	} else if login.selected == SELECTED_GO {
		if enter {
			//login
		}
	}

	return false
}
