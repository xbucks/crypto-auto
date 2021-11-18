package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/hirokimoto/crypto-auto/utils"
)

func TrackPairs(wg *sync.WaitGroup, pairs []string, target chan string) {
	defer wg.Done()

	for _, pair := range pairs {
		go utils.Post(target, "swaps", pair)
		fmt.Print(".")
	}
	time.Sleep(time.Second * 5)
}
