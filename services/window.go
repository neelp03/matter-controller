package services

import "sync"

var WindowOpen = false
var WindowEventFlag = false
var WindowMu sync.Mutex

func GetWindowStatus() (boolean, boolean) {
	windowMu.Lock()
	defer windowMu.Unlock()
	if windowEventFlag {
		windowEventFlag = false
		return windowOpen, true
	}
	return windowOpen, false
}
