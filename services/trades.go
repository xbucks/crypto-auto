package services

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"

	regression "github.com/gaillard/go-online-linear-regression/v1"
	"github.com/hirokimoto/crypto-auto/utils"
)

func AnalyzePairs(command <-chan string, progress chan<- int, t *Tokens) {
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
				trackPair(pair, index, t, progress)
			}
		}
	}
}

func trackPair(pair string, index int, t *Tokens, progress chan<- int) {
	duration, _ := strconv.Atoi(os.Getenv("SWAP_DURATION"))
	var wg sync.WaitGroup
	wg.Add(1)

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

		// Filter our some tokens which is in the active trading in recent3 days.
		if howOld < 3*24 && price > 0.0001 {
			slope, _, _ := testRegression(swaps)
			var isGoingUp = slope > 0
			var isGoingDown = slope < 0
			// var isGoingUp = checkupOfSwaps(swaps)
			// var isGoingDown = checkdownOfSwaps(swaps)
			var isStable = math.Abs((average-price)/price) < 0.1
			var isUnStable = math.Abs((average-price)/price) > 0.1

			target := ""
			updown := ""
			if isUnStable {
				target = "unstable"
				// Notify("Unstable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/", gosxnotifier.Blow)
				fmt.Println("Unstable token ", name, price, average, change, period)
			}
			if isStable {
				target = "stable"
				// Notify("Stable token!", fmt.Sprintf("%s %f %f", name, price, change), "https://kek.tools/", gosxnotifier.Blow)
				fmt.Println("Stable token ", name, price, average, change, period)
			}
			if isGoingUp {
				updown = "up"
				fmt.Println("Trending up token ", name, price, average, change, period)
			}
			if isGoingDown {
				updown = "down"
				fmt.Println("Trending down token ", name, price, average, change, period)
			}

			if isUnStable || isStable || isGoingUp || isGoingDown {
				ct := &Token{
					target:  target,
					updown:  updown,
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
	}
	t.SetProgress(index)
	fmt.Print(index, "|")

	defer wg.Done()
	progress <- index
}

func testRegression(swaps utils.Swaps) (float64, float64, float64) {
	r := regression.New(7)

	for i := 0; i < len(swaps.Data.Swaps); i++ {
		swap := swaps.Data.Swaps[i]
		price, _ := priceOfSwap(swap)
		r.Add(float64(i), price)
	}

	slope, intercept, stdError := r.CalculateWithStdError()
	return slope, intercept, stdError
}
