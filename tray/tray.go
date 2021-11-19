package tray

import (
	"fmt"
	"io/ioutil"

	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/services"
	"github.com/hirokimoto/crypto-auto/views"
	"github.com/skratchdot/open-golang/open"
)

func OnReady() {
	systray.SetIcon(getIcon("assets/auto.ico"))

	mETH := systray.AddMenuItem("ETH", "Price of ethereum")
	mBTC := systray.AddMenuItem("BTC", "Price of bitcoin")
	systray.AddSeparator()
	mDashboard := systray.AddMenuItem("Open Dashboard", "Opens a simple HTML Hello, World")
	mKekBrowser := systray.AddMenuItem("KEK in Browser", "Opens Google in a normal browser")
	mDexEmbed := systray.AddMenuItem("DEX in Window", "Opens Google in a custom window")
	mStables := systray.AddMenuItem("Stable tokens", "Find stable tokens")
	mTradables := systray.AddMenuItem("Tradable tokens", "Find tradable tokens")
	mSettings := systray.AddMenuItem("Settings", "Opens Google in a custom window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit example tray application")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	services.TrackPairs()
	c := &services.Tokens{}

	for {
		select {

		case <-mETH.ClickedCh:
			go services.TrackStable(c)
		case <-mBTC.ClickedCh:
		case <-mDashboard.ClickedCh:
			err := views.Get().OpenIndex()
			if err != nil {
				fmt.Println(err)
			}
		case <-mKekBrowser.ClickedCh:
			err := open.Run("https://www.google.com")
			if err != nil {
				fmt.Println(err)
			}
		case <-mDexEmbed.ClickedCh:
			err := views.Get().OpenGoogle()
			if err != nil {
				fmt.Println(err)
			}
		case <-mStables.ClickedCh:
			err := views.Get().OpenStables()
			if err != nil {
				fmt.Println(err)
			}
		case <-mTradables.ClickedCh:
			err := views.Get().OpenTradables()
			if err != nil {
				fmt.Println(err)
			}
		case <-mSettings.ClickedCh:
			err := views.Get().OpenSettings()
			if err != nil {
				fmt.Println(err)
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-sigc:
			systray.Quit()
		}
	}
}

func OnQuit() {
	close(views.Get().Shutdown)
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
