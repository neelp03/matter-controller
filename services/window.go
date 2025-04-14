package services

import "sync"

var WindowOpen = false
var WindowEventFlag = false
var WindowMu sync.Mutex

func GetWindowStatus() (bool, bool) {
	WindowMu.Lock()
	defer WindowMu.Unlock()
	if WindowEventFlag {
		WindowEventFlag = false
		return WindowOpen, true
	}
	return WindowOpen, false
}
