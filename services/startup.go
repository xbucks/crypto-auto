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

func Startup() {
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan string)
	go func() {
		money := accounting.Accounting{Symbol: "$", Precision: 6}
		olds := []float64{0.1}

		for {
			pair, ok := <-ch
			if !ok {
				println("done")
				wg.Done()
				return
			}

			cc := make(chan string, 1)
			var swaps utils.Swaps
			go utils.Post(cc, "swaps", 2, 0, pair)
			msg := <-cc
			ai := 0.1
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
			println(pair)
		}
	}()

	ch <- "0x7a99822968410431edd1ee75dab78866e31caf39"
	close(ch)

	wg.Wait()
}
