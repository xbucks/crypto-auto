package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/leekchan/accounting"
)

var autoPrice float64 = 0.0

func Startup(command <-chan string) {
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
	go utils.Post(cc, "swaps", 2, 0, "0x7a99822968410431edd1ee75dab78866e31caf39")

	msg := <-cc
	json.Unmarshal([]byte(msg), &swaps)
	n, p, c, d, a := SwapsInfo(swaps, 0.1)

	price := money.FormatMoney(p)
	change := money.FormatMoney(c)
	duration := fmt.Sprintf("%.2f hours", d)

	systray.SetTitle(fmt.Sprintf("%s|%f", n, p))
	t := time.Now()
	fmt.Println(t.Format("2006/01/02 15:04:05"), ": ", n, price, change, duration, a)

	if p != autoPrice {
		message := fmt.Sprintf("%s: %s %s %s", n, price, change, duration)
		Notify("Price changed!", message, "https://kek.tools/")
	}
	autoPrice = p
}
