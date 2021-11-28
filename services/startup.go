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
	{"0x3dd49f67e9d5bc4c5e6634b3f70bfd9dc1b6bd74", 5.0, 8.0},       // SAND
	{"0x7a99822968410431edd1ee75dab78866e31caf39", 0.3, 0.5},       // XI
	{"0x0d0d65e7a7db277d3e0f5e1676325e75f3340455", 1.3, 1.5},       // MTA
	{"0x11b1f53204d03e5529f09eb3091939e4fd8c9cf3", 4.0, 5.3},       // MANA
	{"0x22527f92f43dc8bea6387ce40b87ebaa21f51df3", 2.0, 2.5},       // NUM
	{"0x0529bf56c9448eafe144c151402bc11c0ff47c4c", 0.4, 0.5},       // EPK
	{"0x98d677887af8a699be38ef6276f4cd84aca29d74", 0.0002, 0.0003}, // GM
	{"0x6ada49aeccf6e556bb7a35ef0119cc8ca795294a", 0.8, 0.9},       // WOO
	{"0x05be6820730b30086d6355c44c424230aaff41fb", 0.4, 0.45},      // VEMP
	{"0x9ff68f61ca5eb0c6606dc517a9d44001e564bb66", 1.15, 1.3},      // BOTTO
	{"0xe6f19dab7d43317344282f803f8e8d240708174a", 0.8, 0.85},      // KEEP
	{"0xd1b32d65ae5add070fb54bab3f28033c5a72c849", 0.02, 0.03},     // CENT
	{"0x27fd0857f0ef224097001e87e61026e39e1b04d1", 0.45, 0.52},     // RLY
	{"0xb6909b960dbbe7392d405429eb2b3649752b4838", 1.5, 1.8},       // BAT
	{"0xf5e875b9f457f2dd8112bd68999eb72befb17b03", 0.25, 0.32},     // $ADS
	{"0x700fc86c46299cf2a8fd86edadae3f57014351b0", 0.0045, 0.0052}, // RACA
	{"0x0f5a2eb364d8b722cba4e1e30e2cf57b6f515b2a", 0.4, 0.43},      // TVK
	{"0x470e8de2ebaef52014a47cb5e6af86884947f08c", 0.65, 0.71},     // FOX
	{"0xc5be99a02c6857f9eac67bbce58df5572498f40c", 1.10, 1.18},     // AMP
	{"0x4a7d4be868e0b811ea804faf0d3a325c3a29a9ad", 0.55, 0.65},     // REQ
	{"0x452c60e1e3ae0965cd27db1c7b3a525d197ca0aa", 0.025, 0.032},   // VADER
	{"0x1e9ed2a6ae58f49b3f847eb9f301849c4a20b7e3", 7.0, 7.5},       // GSWAP
	{"0x4214290310264a27b0ba8cff02b4c592d0234aa1", 0.25, 0.28},     // RFOX
	{"0x80d972d2a62ba71814f4e08bd27f95e5d81d02a9", 3.6, 4.0},       // STOS
	{"0x80b4d4e9d88d9f78198c56c5a27f3bacb9a685c5", 0.55, 0.65},     // TRU
	{"0xb8b7c440c36e31686bf1e1bdca76a52e730190fc", 8.0, 9.0},       // NGL
	{"0xd3e5ca0afeae61a24ff7a9219067e51f4bfdd8d9", 0.12, 0.18},     // UMAD
	{"0xc88ac988a655b91b70def427c8778b4d43f2048d", 6.7, 8.0},       // DERC
	{"0xccb63225a7b19dcf66717e4d40c9a72b39331d61", 8.0, 11.0},      // MC
	{"0xc0a6bb3d31bb63033176edba7c48542d6b4e406d", 5.0, 8.0},       // RNDR
	{"0xc8ca3c0f011fe42c48258ecbbf5d94c51f141c17", 2.0, 2.5},       // CGG
	{"0x4d3138931437dcc356ca511ac812e14ba8199fd6", 0.16, 0.22},     // BONDLY
	{"0x0dbd5d63e04aadee8641b04829d125e3943c6b19", 3.3, 3.8},       // $ICONS
	{"0x42d52847be255eacee8c3f96b3b223c0b3cc0438", 2.2, 2.4},       // UOS

	{"0xb3d994978d2bc50d2ce74c45fcd923e7c9c06730", 0.12, 0.14}, // NTX
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
		title := "Priced Up!"
		if c < 0 {
			title = "Priced Down!"
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
