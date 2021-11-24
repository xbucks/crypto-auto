package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hirokimoto/crypto-auto/utils"
)

func GetAllPairs(target chan int) {
	skip := 0

	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(1)
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
			defer wg.Done()
		}
	}()
	target <- 111
}
