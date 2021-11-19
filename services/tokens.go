package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/hirokimoto/crypto-auto/utils"
)

func TrackTokens(target chan string) {
	tokens := []string{"0x74b23882a30290451a17c44f4f05243b6b58c76d"}

	go func() {
		for {
			var wg sync.WaitGroup
			wg.Add(len(tokens))

			go trackTokens(&wg, tokens, target)

			time.Sleep(1 * time.Second)
			wg.Wait()
		}
	}()
}

func trackTokens(wg *sync.WaitGroup, tokens []string, target chan string) {
	go func() {
		for _, pair := range tokens {
			defer wg.Done()
			go utils.Post(target, "tokens", 10, pair)
			fmt.Print(".")
			time.Sleep(time.Second * 5)
		}
	}()
}

func TrackETH(target chan string) {
	go func() {
		for {
			go utils.Post(target, "bundles", 10, "")
			time.Sleep(time.Second * 5)
		}
	}()
}
