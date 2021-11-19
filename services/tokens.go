package services

import (
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
	go func() {
		for {
			go utils.Post(target, "swaps", 2, 0, "0xec454eda10accdd66209c57af8c12924556f3abd")
			time.Sleep(1 * time.Second)
		}
	}()
}
