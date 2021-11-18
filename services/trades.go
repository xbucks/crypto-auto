package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hirokimoto/crypto-auto/utils"
)

func StableTokens(wg *sync.WaitGroup, pairs utils.Pairs, t chan string) {
	for _, item := range pairs.Data.Pairs {
		defer wg.Done()
		cc := make(chan string, 1)
		go utils.Post(cc, "swaps", 1000, item.Id)
		fmt.Print(".")
		stableToken(cc, item.Id, t)
	}
}

func TradableTokens(wg *sync.WaitGroup, pairs utils.Pairs, t chan string) {
	defer wg.Done()

	for _, item := range pairs.Data.Pairs {
		cc := make(chan string, 1)
		go utils.Post(cc, "swaps", 1000, item.Id)
		fmt.Print(".")
		tradableToken(cc, item.Id, t)
	}
}

func StoreAndRemovePair(pair string) (err error) {
	if IsExist(pair) {
		err = RemoveOnePair(pair)
		if err == nil {
			// Alert("Removed!", pair)
		}
	} else {
		err = WriteOnePair(pair)
		if err == nil {
			// Alert("Saved!", pair)
		}
	}
	return err
}

func stableToken(pings chan string, id string, t chan string) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		name, price, change, period, _ := SwapsInfo(swaps, 0.1)

		min, max, _, _, _, _ := minMax(swaps)
		howOld := howMuchOld(swaps)

		if (max-min)/price < 0.1 && period > 24 && howOld < 24 {
			ct := Token{
				name:    name,
				address: id,
				price:   fmt.Sprintf("%f", price),
				change:  fmt.Sprintf("%f", change),
				min:     fmt.Sprintf("%f", min),
				max:     fmt.Sprintf("%f", max),
				period:  fmt.Sprintf("%f", period),
			}
			t <- ct.name
		}
	}
}

func tradableToken(pings chan string, id string, t chan string) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		name, price, change, period, _ := SwapsInfo(swaps, 0.1)

		min, max, _, _, _, _ := minMax(swaps)
		howOld := howMuchOld(swaps)

		if (max-min)/price > 0.1 && period < 6 && howOld < 24 {
			ct := Token{
				name:    name,
				address: id,
				price:   fmt.Sprintf("%f", price),
				change:  fmt.Sprintf("%f", change),
				min:     fmt.Sprintf("%f", min),
				max:     fmt.Sprintf("%f", max),
				period:  fmt.Sprintf("%f", period),
			}
			t <- ct.name
		}
	}
}
