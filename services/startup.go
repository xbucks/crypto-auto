package services

import (
	"encoding/json"
	"fmt"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/leekchan/accounting"
)

var PAIRS = []WatchPair{
	{"0x3dd49f67e9d5bc4c5e6634b3f70bfd9dc1b6bd74", 7.0, 7.3},   // SAND
	{"0x7a99822968410431edd1ee75dab78866e31caf39", 0.42, 0.5},  // XI
	{"0x0d0d65e7a7db277d3e0f5e1676325e75f3340455", 1.3, 1.5},   // MTA
	{"0x11b1f53204d03e5529f09eb3091939e4fd8c9cf3", 4.5, 5.3},   // MANA
	{"0x22527f92f43dc8bea6387ce40b87ebaa21f51df3", 1.5, 2.0},   // NUM
	{"0xc88ac988a655b91b70def427c8778b4d43f2048d", 6.7, 8.0},   // DERC
	{"0xccb63225a7b19dcf66717e4d40c9a72b39331d61", 8.0, 11.0},  // MC
	{"0xc0a6bb3d31bb63033176edba7c48542d6b4e406d", 5.0, 8.0},   // RNDR
	{"0xc8ca3c0f011fe42c48258ecbbf5d94c51f141c17", 2.0, 2.5},   // CGG
	{"0x4d3138931437dcc356ca511ac812e14ba8199fd6", 0.16, 0.22}, // BONDLY
}
var oldPrices = map[string]float64{}

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
				trackSubPairs()
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func trackMainPair() {
	mainPair := PAIRS[0]
	trackOnePair(mainPair, "main")
}

func trackSubPairs() {
	for i := 1; i < len(PAIRS); i++ {
		pair := PAIRS[i]
		trackOnePair(pair, "sub")
	}
}

func trackOnePair(pair WatchPair, target string) {
	money := accounting.Accounting{Symbol: "$", Precision: 6}
	cc := make(chan string, 1)
	var swaps utils.Swaps
	go utils.SwapsByCounts(cc, 2, pair.address)

	msg := <-cc
	json.Unmarshal([]byte(msg), &swaps)
	n, p, c, d, _, _ := SwapsInfo(swaps, 0.1)

	price := money.FormatMoney(p)
	change := money.FormatMoney(c)
	duration := fmt.Sprintf("%.2f hours", d)

	fmt.Print(".")

	if p != oldPrices[pair.address] {
		t := time.Now()
		message := fmt.Sprintf("%s: %s %s %s", n, price, change, duration)
		title := "Price changed up!"
		if c < 0 {
			title = "Price changed down!"
		}
		link := fmt.Sprintf("https://www.dextools.io/app/ether/pair-explorer/%s", pair.address)
		var sound gosxnotifier.Sound
		if target == "main" {
			systray.SetTitle(fmt.Sprintf("%s|%f", n, p))
			sound = gosxnotifier.Sosumi
		} else {
			sound = gosxnotifier.Morse
		}

		if p < pair.min {
			title = fmt.Sprintf("Warning Low! Watch %s", n)
			sound = gosxnotifier.Default
		}
		if p > pair.max {
			title = fmt.Sprintf("Warning High! Watch %s", n)
			sound = gosxnotifier.Default
		}
		Notify(title, message, link, sound)
		fmt.Println(".")
		fmt.Println(t.Format("2006/01/02 15:04:05"), ": ", n, price, change, duration)
		fmt.Println(".")
	}
	oldPrices[pair.address] = p
}
