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

func TrackPairs() {
	money := accounting.Accounting{Symbol: "$", Precision: 6}
	pairs := []string{"0x7a99822968410431edd1ee75dab78866e31caf39"}
	olds := []float64{0.1}

	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(len(pairs))

			cc := make(chan string, 1)
			var swaps utils.Swaps
			go trackPairs(&wg, pairs, cc, 2)

			ai := 0.1
			msg := <-cc
			json.Unmarshal([]byte(msg), &swaps)
			n, p, c, d, a := SwapsInfo(swaps, ai)

			price := money.FormatMoney(p)
			change := money.FormatMoney(c)
			duration := fmt.Sprintf("%.2f hours", d)

			systray.SetTitle(fmt.Sprintf("%s|%f", n, p))
			systray.SetTooltip("Crypto Auto")
			t := time.Now()
			fmt.Println(t.Format("2006/01/02 15:04:05"), ": ", n, price, change, duration, a)

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

func GetAllPairs() {
	skip := 0
	go func() {
		for {
			target := make(chan string, 1)
			go utils.Post(target, "pairs", 1000, 0, "")
			msg := <-target
			var pairs utils.Pairs
			json.Unmarshal([]byte(msg), &pairs)
			counts := len(pairs.Data.Pairs)
			fmt.Println(".", skip)
			if counts == 0 {
				return
			}
			SaveAllPairs(&pairs)
			skip += 1
		}
	}()
}

func trackPairs(wg *sync.WaitGroup, pairs []string, target chan string, limit int) {
	for _, pair := range pairs {
		defer wg.Done()
		go utils.Post(target, "swaps", limit, 0, pair)
		fmt.Print(".")
	}
	time.Sleep(time.Second * 5)
}
