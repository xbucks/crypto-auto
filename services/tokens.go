package services

import (
	"sync"
	"time"

	"github.com/hirokimoto/crypto-auto/utils"
)

func TrackETH(target chan string) {
	go func() {
		for {
			go utils.Post(target, "bundles", 10, 0, "")
			time.Sleep(time.Second * 5)
		}
	}()
}

func TrackBTC(target chan string) {
	pairs := []string{"0xec454eda10accdd66209c57af8c12924556f3abd"}

	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(len(pairs))
			go trackPairs(&wg, pairs, target, 2)
			wg.Wait()
			time.Sleep(1 * time.Second)
		}
	}()
}
