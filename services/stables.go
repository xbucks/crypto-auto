package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hirokimoto/crypto-auto/utils"
)

func StableTokens(wg *sync.WaitGroup, pairs utils.Pairs, t *Tokens) {
	for index, item := range pairs.Data.Pairs {
		defer wg.Done()
		cc := make(chan string, 1)
		go utils.Post(cc, "swaps", 1000, 0, item.Id)
		stableToken(cc, item.Id, t)
		t.SetProgress(index)
		fmt.Print(".")
	}
}

func stableToken(pings chan string, id string, t *Tokens) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		name, price, change, period, _ := SwapsInfo(swaps, 0.1)

		min, max, _, _, _, _ := minMax(swaps)
		howOld := howMuchOld(swaps)

		if (max-min)/price < 0.1 && period > 24 && howOld < 24 {
			ct := &Token{
				name:    name,
				address: id,
				price:   fmt.Sprintf("%f", price),
				change:  fmt.Sprintf("%f", change),
				min:     fmt.Sprintf("%f", min),
				max:     fmt.Sprintf("%f", max),
				period:  fmt.Sprintf("%.2f", period),
			}
			t.Add(ct)
			fmt.Println("New token!!!!!   ", ct.name)
		}
	}
}
