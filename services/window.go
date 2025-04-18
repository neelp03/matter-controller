package services

import (
	"fmt"
	"sync"
)

var WindowOpen = false
var WindowEventFlag = false
var WindowMu sync.Mutex

func OpenWindow() error {
	WindowMu.Lock()
	defer WindowMu.Unlock()

	if WindowOpen {
		fmt.Println("Window is already open")
		return nil
	}

	OpenWindowMotor();
	fmt.Println("Opening window...")
	var err error // may return an error

	if err == nil {
		WindowOpen = true
		WindowEventFlag = !WindowEventFlag
	}

	return err
}

func CloseWindow() error {
	WindowMu.Lock()
	defer WindowMu.Unlock()

	if !WindowOpen {
		fmt.Println("Window is already closed")
		return nil
	}

	CloseWindowMotor();
	fmt.Println("Closing window...")
	var err error // may return an error

	if err == nil {
		WindowOpen = false
		WindowEventFlag = !WindowEventFlag
	}

	return nil
}
