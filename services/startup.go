package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/leekchan/accounting"
)

var oldPrices = map[string]float64{}

func Startup(command <-chan string, alert float64) {
	var status = "Play"
	for {
		select {
		case cmd := <-command:
			fmt.Println(cmd)
			switch cmd {
			case "Stop":
				return
			case "Pause":
				status = "Pause"
			default:
				status = "Play"
			}
		default:
			if status == "Play" {
				trackMainPair()
				trackSubPairs()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func trackMainPair() {
	address := os.Getenv("MAIN_PAIR")
	trackOnePair(address, "main")
}

func trackSubPairs() {
	pairs := []string{
		"0x22527f92f43dc8bea6387ce40b87ebaa21f51df3",
		"0x684b00a5773679f88598a19976fbeb25a68e9a5f",
		"0xc88ac988a655b91b70def427c8778b4d43f2048d"}
	for _, v := range pairs {
		trackOnePair(v, "sub")
	}
}

func trackOnePair(address string, target string) {
	money := accounting.Accounting{Symbol: "$", Precision: 6}
	cc := make(chan string, 1)
	var swaps utils.Swaps
	go utils.SwapsByCounts(cc, 2, address)

	msg := <-cc
	json.Unmarshal([]byte(msg), &swaps)
	n, p, c, d, _, a := SwapsInfo(swaps, 0.1)

	price := money.FormatMoney(p)
	change := money.FormatMoney(c)
	duration := fmt.Sprintf("%.2f hours", d)

	fmt.Print(".")

	if p != oldPrices[address] {
		t := time.Now()
		message := fmt.Sprintf("%s: %s %s %s", n, price, change, duration)
		title := "Price changed up!"
		if c < 0 {
			title = "Price changed down!"
		}
		link := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", address)
		if target == "main" {
			systray.SetTitle(fmt.Sprintf("%s|%f", n, p))
			Notify(title, message, link, gosxnotifier.Default)
		} else {
			Notify(title, message, link, gosxnotifier.Glass)
		}
		fmt.Println(".")
		fmt.Println(t.Format("2006/01/02 15:04:05"), ": ", n, price, change, duration, a)
		fmt.Println(".")
	}
	oldPrices[address] = p
}
