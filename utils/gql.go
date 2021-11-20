package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func request(query map[string]string, target chan string) {
	jsonQuery, _ := json.Marshal(query)
	request, _ := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", bytes.NewBuffer(jsonQuery))
	client := &http.Client{Timeout: time.Second * 50}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		target <- ""
		return
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	target <- string(data)
}

func Post(target chan string, to string, limit int, skip int, id string) {
	query := Query(to, limit, skip, id)
	request(query, target)
}

func SwapsByDays(target chan string, limit int, id string) {
	var results Swaps
	var temp Swaps
	i := 0
	now := time.Now()
	for {
		ch := make(chan string)
		go Post(ch, "swaps", 1000, i*1000, id)

		msg := <-ch
		var swaps Swaps
		json.Unmarshal([]byte(msg), &swaps)
		endIndex := len(swaps.Data.Swaps) - 1

		lastInt, _ := strconv.ParseInt(swaps.Data.Swaps[endIndex].Timestamp, 10, 64)
		lastTime := time.Unix(lastInt, 0)
		period := now.Sub(lastTime)

		if period.Hours() >= 24*float64(limit) {
			temp = swaps
			break
		} else {
			results.Data.Swaps = append(results.Data.Swaps, swaps.Data.Swaps...)
		}
		i += 1
	}

	for i = 0; i < len(temp.Data.Swaps); i++ {
		lastInt, _ := strconv.ParseInt(temp.Data.Swaps[i].Timestamp, 10, 64)
		lastTime := time.Unix(lastInt, 0)
		period := now.Sub(lastTime)
		if period.Hours() > 24*float64(limit) {
			break
		} else {
			results.Data.Swaps = append(results.Data.Swaps, temp.Data.Swaps[i])
		}
	}

	tg, _ := json.Marshal(results)
	target <- string(tg)
}

func SwapsByCounts(target chan string, limit int, id string) {
	var results Swaps
	length := limit/1000 + 1
	var wg sync.WaitGroup
	wg.Add(length)

	for i := 0; i < length; i++ {
		ch := make(chan string)
		counts := 1000
		skip := i * 1000
		if limit < 1000 {
			counts = limit
		}
		if (i+1)*1000 > limit {
			counts = limit - i*counts
		}

		go Post(ch, "swaps", counts, skip, id)

		msg := <-ch
		var swaps Swaps
		json.Unmarshal([]byte(msg), &swaps)
		results.Data.Swaps = append(results.Data.Swaps, swaps.Data.Swaps...)

		wg.Done()
	}

	wg.Wait()

	tg, _ := json.Marshal(results)
	target <- string(tg)
}
