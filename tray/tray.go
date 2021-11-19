package tray

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/services"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/hirokimoto/crypto-auto/views"
	"github.com/leekchan/accounting"
	"github.com/skratchdot/open-golang/open"
)

func OnReady() {
	systray.SetIcon(getIcon("assets/auto.ico"))

	mETH := systray.AddMenuItem("ETH", "Price of ethereum")
	mETH.SetIcon(getIcon("assets/eth.ico"))
	mBTC := systray.AddMenuItem("BTC", "Price of bitcoin")
	mBTC.SetIcon(getIcon("assets/btc.ico"))
	mBTC.Disable()
	systray.AddSeparator()
	mStart := systray.AddMenuItem("Start", "Start background tracker to find tradable tokens")
	mStop := systray.AddMenuItem("Stop", "Stop background tracker to find tradable tokens")
	systray.AddSeparator()
	mRefreshPairs := systray.AddMenuItem("Refresh pairs", "Get all available pairs")
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

	money := accounting.Accounting{Symbol: "$", Precision: 2}
	ethc := make(chan string, 1)
	btcc := make(chan string, 1)
	pirc := make(chan int, 1)

	for {
		select {

		case <-mETH.ClickedCh:
			services.TrackETH(ethc)
		case <-mBTC.ClickedCh:
			services.TrackBTC(btcc)
		case <-mStart.ClickedCh:
			services.GetAllPairs(pirc)
		case <-mStop.ClickedCh:
		case <-mRefreshPairs.ClickedCh:
			services.GetAllPairs(pirc)
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
		case <-ethc:
			msg := <-ethc
			var eth utils.Bundles
			json.Unmarshal([]byte(msg), &eth)
			_price, _ := strconv.ParseFloat(eth.Data.Bundles[0].EthPrice, 32)
			price := fmt.Sprintf("$%.2f", _price)
			mETH.SetTitle(price)
			fmt.Println("ETH Price: ", price)
		case <-btcc:
			msg := <-btcc
			var swaps utils.Swaps
			json.Unmarshal([]byte(msg), &swaps)
			_, p, _, _, _ := services.SwapsInfo(swaps, 0.1)
			price := money.FormatMoney(p)
			fmt.Println("BTC Price: ", price)
		case <-pirc:
			msg := <-pirc
			mRefreshPairs.SetTitle(fmt.Sprintf("Refreshing pairs %d...", msg))
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
