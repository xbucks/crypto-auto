package services

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/hirokimoto/crypto-auto/utils"
)

func AnalyzePairs(command <-chan string, progress chan<- int, duration int, t *Tokens) {
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
				trackPair(pair, index, duration, t)
			}
		}
		progress <- index
	}
}

func trackPair(pair string, index int, duration int, t *Tokens) {
	ch := make(chan string, 1)
	if duration > 100 {
		go utils.SwapsByCounts(ch, duration, pair)
	} else {
		go utils.SwapsByDays(ch, duration, pair)
	}

	msg := <-ch
	var swaps utils.Swaps
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		name, price, change, period, average, _ := SwapsInfo(swaps, 0.1)

		min, max, _, _, _, _ := minMax(swaps)
		howOld := howMuchOld(swaps)

		inPeriodStable := true
		inPeriodUnStable := true
		if duration > 100 {
			inPeriodStable = period > 24 && howOld < 24
			inPeriodUnStable = period < 3*24 && howOld < 24
		} else {
			inPeriodStable = true
			inPeriodUnStable = true
		}

		var isGoingUp = checkupOfSwaps(swaps)
		var isGoingDown = checkdownOfSwaps(swaps)
		var isStable = math.Abs((average-price)/price) < 0.1 && inPeriodStable
		var isUnStable = math.Abs((average-price)/price) > 0.1 && inPeriodUnStable && price > 0.0001

		target := ""
		if isUnStable {
			target = "tradable"
			Notify("Tradable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/")
			fmt.Println("Tradable token ", name, price, average, change, period)
		}
		if isStable {
			target = "stable"
			Notify("Stable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/")
			fmt.Println("Stable token ", name, price, average, change, period)
		}
		if isGoingUp {
			target = "up"
			fmt.Println("Trending up token ", name, price, average, change, period)
		}
		if isGoingDown {
			target = "down"
			fmt.Println("Trending down token ", name, price, average, change, period)
		}

		if isUnStable || isStable || isGoingUp || isGoingDown {
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
