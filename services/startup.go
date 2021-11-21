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

var autoPrice float64 = 0.0

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
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func trackMainPair() {
	money := accounting.Accounting{Symbol: "$", Precision: 6}
	cc := make(chan string, 1)
	var swaps utils.Swaps
	address := os.Getenv("MAIN_PAIR")
	go utils.SwapsByCounts(cc, 2, address)

	msg := <-cc
	json.Unmarshal([]byte(msg), &swaps)
	n, p, c, d, _, a := SwapsInfo(swaps, 0.1)

	price := money.FormatMoney(p)
	change := money.FormatMoney(c)
	duration := fmt.Sprintf("%.2f hours", d)

	systray.SetTitle(fmt.Sprintf("%s|%f", n, p))
	t := time.Now()
	fmt.Print(".")

	if p != autoPrice {
		message := fmt.Sprintf("%s: %s %s %s", n, price, change, duration)
		Notify("Price changed!", message, "https://kek.tools/", gosxnotifier.Default)
		fmt.Println(t.Format("2006/01/02 15:04:05"), ": ", n, price, change, duration, a)
	}
	autoPrice = p
}
