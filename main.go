package main

import (
	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/tray"
	"github.com/hirokimoto/crypto-auto/views"
)

func main() {
	views := views.Get()
	defer views.WaitGroup.Wait()
	systray.Run(tray.OnReady, tray.OnQuit)
}
