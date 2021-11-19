package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hirokimoto/crypto-auto/utils"
)

func TradablePairs(command <-chan string, t *Tokens) {
	pairs, _ := ReadAllPairs()
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
		time.Sleep(1 * time.Second)
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

		if (max-min)/price > 0.1 && period < 24*3 && howOld < 24 && price > 0.0001 {
			fmt.Println("Tradable token !!!!!   ", name, price, change, period)
			ct := &Token{
				name:    name,
				address: pair,
				price:   fmt.Sprintf("%f", price),
				change:  fmt.Sprintf("%f", change),
				min:     fmt.Sprintf("%f", min),
				max:     fmt.Sprintf("%f", max),
				period:  fmt.Sprintf("%.2f", period),
			}
			t.Add(ct)
		}
	}
	fmt.Print(index, "|")
}
