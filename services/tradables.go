package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hirokimoto/crypto-auto/utils"
)

func TrackTradable(t *Tokens) {
	pc := make(chan string, 1)
	for {
		go utils.Post(pc, "pairs", 1000, "")

		msg1 := <-pc
		var pairs utils.Pairs
		json.Unmarshal([]byte(msg1), &pairs)
		counts := len(pairs.Data.Pairs)
		fmt.Println("Counts of Pairs: ", counts)
		if counts > 0 {
			var wg sync.WaitGroup
			wg.Add(counts)
			TradableTokens(&wg, pairs, t)
			wg.Wait()
		}

		time.Sleep(time.Minute * 10)
	}
}

func TradableTokens(wg *sync.WaitGroup, pairs utils.Pairs, t *Tokens) {
	defer wg.Done()

	for index, item := range pairs.Data.Pairs {
		cc := make(chan string, 1)
		go utils.Post(cc, "swaps", 1000, item.Id)
		tradableToken(cc, item.Id, t)
		t.SetProgress(index)
		fmt.Print(".")
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

func tradableToken(pings chan string, id string, t *Tokens) {
	var swaps utils.Swaps
	msg := <-pings
	json.Unmarshal([]byte(msg), &swaps)

	if len(swaps.Data.Swaps) > 0 {
		name, price, change, period, _ := SwapsInfo(swaps, 0.1)

		min, max, _, _, _, _ := minMax(swaps)
		howOld := howMuchOld(swaps)

		if (max-min)/price > 0.1 && period < 24*3 && howOld < 24 && price > 0.0001 {
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
