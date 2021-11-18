package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/leekchan/accounting"
)

func TracPairs() {
	money := accounting.Accounting{Symbol: "$", Precision: 6}
	pairs := []string{"0x7a99822968410431edd1ee75dab78866e31caf39"}
	olds := []float64{0.1}

	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(len(pairs))

			cc := make(chan string, 1)
			var swaps utils.Swaps
			go TrackPairs(&wg, pairs, cc)

			ai := 0.1
			msg := <-cc
			json.Unmarshal([]byte(msg), &swaps)
			n, p, c, d, a := SwapsInfo(swaps, ai)

			price := money.FormatMoney(p)
			change := money.FormatMoney(c)
			duration := fmt.Sprintf("%.2f hours", d)

			systray.SetTitle(fmt.Sprintf("%s %s", n, price))
			systray.SetTooltip("Crypto Auto")
			fmt.Println("---->>>  ", n, change, duration, a)

			if p != olds[0] {
				message := fmt.Sprintf("%s: %s %s %s", n, price, change, duration)
				url := "https://kek.tools/t/0x295b42684f90c77da7ea46336001010f2791ec8c?pair=0x7a99822968410431edd1ee75dab78866e31caf39"
				Notify("Price changed!", message, url)
			}
			olds[0] = p

			time.Sleep(1 * time.Second)
			wg.Wait()
		}
	}()
}

func TrackStable(t *Tokens) {
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
			StableTokens(&wg, pairs, t)
			wg.Wait()
		}

		time.Sleep(time.Minute * 10)
	}
}

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

func StableTokens(wg *sync.WaitGroup, pairs utils.Pairs, t *Tokens) {
	for _, item := range pairs.Data.Pairs {
		defer wg.Done()
		cc := make(chan string, 1)
		go utils.Post(cc, "swaps", 1000, item.Id)
		fmt.Print(".")
		stableToken(cc, item.Id, t)
	}
}

func TradableTokens(wg *sync.WaitGroup, pairs utils.Pairs, t *Tokens) {
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

func stableToken(pings chan string, id string, t *Tokens) {
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
			t.Add(ct)
			fmt.Println("New token!!!!!   ", ct.name)
		}
	}
}

func tradableToken(pings chan string, id string, t *Tokens) {
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
			t.Add(ct)
			fmt.Println("New token!!!!!   ", ct.name)
		}
	}
}
