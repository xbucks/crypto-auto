package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/utils"
	"github.com/leekchan/accounting"
)

func TrackPairs() {
	go func() {
		money := accounting.Accounting{Symbol: "$", Precision: 6}
		olds := []float64{0.1}
		for {
			cc := make(chan string, 1)
			var swaps utils.Swaps
			go utils.Post(cc, "swaps", 2, 0, "0x7a99822968410431edd1ee75dab78866e31caf39")

			msg := <-cc
			ai := 0.1
			json.Unmarshal([]byte(msg), &swaps)
			n, p, c, d, _, a := SwapsInfo(swaps, ai)

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
				Notify("Price changed!", message, url, gosxnotifier.Default)
			}
			olds[0] = p

			time.Sleep(1 * time.Second)
		}
	}()
}

func GetAllPairs(target chan int) {
	skip := 0
	var v sync.WaitGroup
	v.Add(1)
	go func(wg sync.WaitGroup) {
		defer wg.Done()

		for {
			cc := make(chan string, 1)
			go utils.Post(cc, "pairs", 1000, 1000*skip, "")
			msg := <-cc
			var pairs utils.Pairs
			json.Unmarshal([]byte(msg), &pairs)
			counts := len(pairs.Data.Pairs)
			fmt.Println(skip, ": ", counts)
			if counts == 0 {
				target <- 111
				return
			}
			SaveAllPairs(&pairs)
			skip += 1
			target <- skip
			time.Sleep(time.Millisecond * 200)
		}
	}(v)
	target <- 111
}

func trackPairs(wg *sync.WaitGroup, pairs []string, target chan string, limit int) {
	for _, pair := range pairs {
		defer wg.Done()
		go utils.Post(target, "swaps", limit, 0, pair)
		fmt.Print(".")
	}
	time.Sleep(time.Second * 5)
}
