package main

import "time"

func changeGameMode(newMode MODE, delay time.Duration) {

	/* Skip if the same */
	if newMode == gameMode {
		return
	}

	if newMode == MODE_CONNECTED {
		//Send INIT

	}

	time.Sleep(delay)
	gameMode = newMode
}
