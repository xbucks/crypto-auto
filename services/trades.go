package services

import (
	"encoding/json"
	"fmt"

	"github.com/hirokimoto/crypto-auto/utils"
)

func TradePairs(command <-chan string, progress chan<- int, t *Tokens) {
	pairs, _ := ReadAllPairs()
	t.SetTotal(len(pairs))
	var status = "Play"
	for index, pair := range pairs {
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
				trackPair(pair, index, t)
			}
		}
		progress <- index
	}
}

func trackPair(pair string, index int, t *Tokens) {
	ch := make(chan string, 1)
	go utils.Post(ch, "swaps", 1000, 0, pair)

	msg := <-ch
	var swaps utils.Swaps
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		name, price, change, period, _ := SwapsInfo(swaps, 0.1)

		min, max, _, _, _, _ := minMax(swaps)
		howOld := howMuchOld(swaps)

		var isTradable = (max-min)/price > 0.1 && period < 24*3 && howOld < 24 && price > 0.0001
		var isStable = (max-min)/price < 0.1 && period > 24 && howOld < 24

		target := ""
		if isTradable {
			target = "tradable"
			Notify("Tradable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/")
		}
		if isStable {
			target = "stable"
			Notify("Stable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/")
		}

		if isTradable || isStable {
			ct := &Token{
				target:  target,
				name:    name,
				address: pair,
				price:   fmt.Sprintf("%f", price),
				change:  fmt.Sprintf("%f", change),
				min:     fmt.Sprintf("%f", min),
				max:     fmt.Sprintf("%f", max),
				period:  fmt.Sprintf("%.2f", period),
				swaps:   swaps.Data.Swaps,
			}
			t.Add(ct)
		}
	}
	t.SetProgress(index)
	fmt.Print(index, "|")
}