package main

import (
	"github.com/ben-turner/explosive-transistor2/controllers"
	"log"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32 = syscall.MustLoadDLL("User32.dll")
)

type Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func AddWatches(devs map[string]controllers.Controller) {
	go func() {
		var previous bool
		alreadySet := false
		for {
			time.Sleep(time.Second * 2)
			if t, err := isMainWindowFullscreen(); err == nil {
				if t && !alreadySet {
					p, err := devs["huebridge"].Get(controllers.GroupId(1))
					if err != nil {
						log.Println(err)
						return
					}
					previous = p.On
					if err := devs["huebridge"].Set(controllers.GroupId(1), &controllers.State{
						On: false,
					}); err != nil {
						log.Println(err.Error())
						return
					}
					alreadySet = true
				} else if !t && alreadySet {
					if err := devs["huebridge"].Set(controllers.GroupId(1), &controllers.State{
						On: previous,
					}); err != nil {
						log.Println(err.Error())
						return
					}
					alreadySet = false
				}
			} else {
				log.Println(err.Error())
				return
			}
		}
	}()
}

func isMainWindowFullscreen() (bool, error) {
	topWindow, _, err := user32.MustFindProc("GetForegroundWindow").Call()
	if topWindow == 0 {
		return false, err
	}

	desktop, _, err := user32.MustFindProc("GetDesktopWindow").Call()
	if desktop == 0 {
		log.Println(err)
	} else if desktop == topWindow {
		log.Println("Desktop is focused", topWindow)
		return false, nil
	}

	GetWindowRect := user32.MustFindProc("GetWindowRect")
	windowRect := Rect{}
	a, _, err := GetWindowRect.Call(topWindow, uintptr(unsafe.Pointer(&windowRect)))
	if a == 0 {
		return false, err
	}

	desktopRect := Rect{}
	a, _, err = GetWindowRect.Call(desktop, uintptr(unsafe.Pointer(&desktopRect)))
	if a == 0 {
		return false, err
	}

	return windowRect == desktopRect, nil
}
