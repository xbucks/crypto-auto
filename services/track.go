package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/crypto-auto/utils"
	"github.com/leekchan/accounting"
)

func TrackPairs(wg *sync.WaitGroup, pairs []string) {
	defer wg.Done()

	money := accounting.Accounting{Symbol: "$", Precision: 6}
	for _, pair := range pairs {
		var swaps utils.Swaps
		cc := make(chan string, 1)
		ai := 0.1

		go utils.Post(cc, "swaps", pair)
		fmt.Print(".")

		msg := <-cc
		json.Unmarshal([]byte(msg), &swaps)
		n, p, c, d, a := SwapsInfo(swaps, ai)

		price := money.FormatMoney(p)
		change := money.FormatMoney(c)
		duration := fmt.Sprintf("%.2f hours", d)
	}
	time.Sleep(time.Second * 5)
}
