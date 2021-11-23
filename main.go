package main

import (
	"os"

	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/tray"
	"github.com/hirokimoto/crypto-auto/views"
)

func main() {
	os.Setenv("MAIN_PAIR", "0x7a99822968410431edd1ee75dab78866e31caf39")
	os.Setenv("SWAP_DURATION", "1000")
	os.Setenv("PRICE_ALERT", "0")

	views := views.Get()
	defer views.WaitGroup.Wait()
	systray.Run(tray.OnReady, tray.OnQuit)
}
