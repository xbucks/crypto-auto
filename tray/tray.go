package tray

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

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

	mStart := systray.AddMenuItem("Start", "Start tracking pairs.")
	mStop := systray.AddMenuItem("Stop", "Stop tracking pairs.")
	systray.AddSeparator()
	mDashboard := systray.AddMenuItem("Open Dashboard", "Opens a simple HTML Hello, World")
	mKekBrowser := systray.AddMenuItem("KEK in Browser", "Opens Google in a normal browser")
	mDexEmbed := systray.AddMenuItem("DEX in Window", "Opens Google in a custom window")
	mStables := systray.AddMenuItem("Stable tokens", "Find stable tokens")
	mSettings := systray.AddMenuItem("Settings", "Opens Google in a custom window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit example tray application")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	trackPairs()

	token := make(chan services.Token, 1)

	for {
		select {

		case <-mStart.ClickedCh:
			go trackStable(token)

		case <-mStop.ClickedCh:
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
		case <-mSettings.ClickedCh:
			err := views.Get().OpenSettings()
			if err != nil {
				fmt.Println(err)
			}
		case msg1 := <-token:
			fmt.Println("New token!!!!! ", msg1.Get())
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

func trackPairs() {
	money := accounting.Accounting{Symbol: "$", Precision: 6}
	pairs := []string{"0x7a99822968410431edd1ee75dab78866e31caf39"}
	olds := []float64{0.1}

	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(len(pairs))

			cc := make(chan string, 1)
			var swaps utils.Swaps
			go services.TrackPairs(&wg, pairs, cc)

			ai := 0.1
			msg := <-cc
			json.Unmarshal([]byte(msg), &swaps)
			n, p, c, d, a := services.SwapsInfo(swaps, ai)

			price := money.FormatMoney(p)
			change := money.FormatMoney(c)
			duration := fmt.Sprintf("%.2f hours", d)

			systray.SetTitle(fmt.Sprintf("%s %s", n, price))
			systray.SetTooltip("Crypto Auto")
			fmt.Println(getClockTime("Local"), "---->>>  ", n, change, duration, a)

			if p != olds[0] {
				message := fmt.Sprintf("%s: %s %s %s", n, price, change, duration)
				url := "https://kek.tools/t/0x295b42684f90c77da7ea46336001010f2791ec8c?pair=0x7a99822968410431edd1ee75dab78866e31caf39"
				services.Notify("Price changed!", message, url)
			}
			olds[0] = p

			time.Sleep(1 * time.Second)
			wg.Wait()
		}
	}()
}

func trackStable(t chan services.Token) {
	pc := make(chan string, 1)
	for {
		go utils.Post(pc, "pairs", 1000, "")

		msg1 := <-pc
		var pairs utils.Pairs
		json.Unmarshal([]byte(msg1), &pairs)
		counts := len(pairs.Data.Pairs)
		fmt.Println("Counts of Pairs: ", counts)
		if counts > 0 {
			var wg sync.WaitGroup
			wg.Add(counts)
			services.StableTokens(&wg, pairs, t)
			wg.Wait()
		}

		time.Sleep(time.Minute * 10)
	}
}

func getClockTime(tz string) string {
	t := time.Now()
	utc, _ := time.LoadLocation(tz)

	hour, min, sec := t.In(utc).Clock()
	return itoaTwoDigits(hour) + ":" + itoaTwoDigits(min) + ":" + itoaTwoDigits(sec)
}

func itoaTwoDigits(i int) string {
	b := "0" + strconv.Itoa(i)
	return b[len(b)-2:]
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
