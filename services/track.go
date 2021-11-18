package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/hirokimoto/crypto-auto/utils"
)

func TrackPairs(wg *sync.WaitGroup, pairs []string, target chan string) {
	for _, pair := range pairs {
		defer wg.Done()
		go utils.Post(target, "swaps", 10, pair)
		fmt.Print(".")
	}
	time.Sleep(time.Second * 5)
}
