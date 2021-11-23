package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hirokimoto/crypto-auto/utils"
)

func SwapsInfo(swaps utils.Swaps, ps float64) (name string, price float64, change float64, duration float64, average float64, alert bool) {
	name = tokenName(swaps)
	price, change = priceChanges(swaps)
	_, _, duration = periodOfSwaps(swaps)
	alert = priceAlert(swaps, ps)
	average = averageOfSwaps(swaps.Data.Swaps)
	return name, price, change, duration, average, alert
}

func SwapInfo(swap utils.Swap) (price float64, target string, amount string, amount1 string, amount2 string) {
	price, target = priceOfSwap(swap)
	_amount, _ := strconv.ParseFloat(swap.AmountUSD, 32)
	amount = fmt.Sprintf("$%.2f", _amount)

	_amount1 := 0.1
	if swap.Amount0In == "0" {
		_amount1, _ = strconv.ParseFloat(swap.Amount0Out, 32)
	} else {
		_amount1, _ = strconv.ParseFloat(swap.Amount0In, 32)
	}
	amount1 = fmt.Sprintf("%.4f %s", _amount1, swap.Pair.Token0.Symbol)

	_amount2 := 0.1
	if swap.Amount1In == "0" {
		_amount2, _ = strconv.ParseFloat(swap.Amount1Out, 32)
	} else {
		_amount2, _ = strconv.ParseFloat(swap.Amount1In, 32)
	}
	amount2 = fmt.Sprintf("%.4f %s", _amount2, swap.Pair.Token1.Symbol)

	return price, target, amount, amount1, amount2
}

func tokenName(swaps utils.Swaps) (name string) {
	if swaps.Data.Swaps != nil && len(swaps.Data.Swaps) > 0 {
		if swaps.Data.Swaps[0].Pair.Token0.Symbol == "WETH" {
			name = swaps.Data.Swaps[0].Pair.Token1.Name
		} else {
			name = swaps.Data.Swaps[0].Pair.Token0.Name
		}

	}
	return name
}

func priceChanges(swaps utils.Swaps) (price float64, change float64) {
	if swaps.Data.Swaps != nil && len(swaps.Data.Swaps) > 0 {
		price, _ = priceOfSwap(swaps.Data.Swaps[0])
		last, _ := priceOfSwap(swaps.Data.Swaps[len(swaps.Data.Swaps)-1])
		change = price - last
	}
	return price, change
}

func priceOfSwap(swap utils.Swap) (price float64, target string) {
	amountUSD, _ := strconv.ParseFloat(swap.AmountUSD, 32)
	amountToken, _ := strconv.ParseFloat(swap.Amount0Out, 32)

	if swap.Pair.Token0.Symbol == "WETH" {
		if swap.Amount0In != "0" && swap.Amount1Out != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount1Out, 32)
			target = "BUY"
		} else if swap.Amount0Out != "0" && swap.Amount1In != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount1In, 32)
			target = "SELL"
		}
	} else {
		if swap.Amount0Out != "0" && swap.Amount1In != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount0Out, 32)
			target = "BUY"
		} else if swap.Amount0In != "0" && swap.Amount1Out != "0" {
			amountToken, _ = strconv.ParseFloat(swap.Amount0In, 32)
			target = "SELL"
		}
	}

	price = amountUSD / amountToken
	return price, target
}

func periodOfSwaps(swaps utils.Swaps) (first time.Time, last time.Time, period float64) {
	var duration float64
	if swaps.Data.Swaps != nil && len(swaps.Data.Swaps) > 0 {
		_first, _ := strconv.ParseInt(swaps.Data.Swaps[0].Timestamp, 10, 64)
		_last, _ := strconv.ParseInt(swaps.Data.Swaps[len(swaps.Data.Swaps)-1].Timestamp, 10, 64)
		first = time.Unix(_first, 0)
		last = time.Unix(_last, 0)
		_period := first.Sub(last)
		duration = _period.Hours()
	}
	return first, last, duration
}

func priceAlert(swaps utils.Swaps, change float64) (state bool) {
	if swaps.Data.Swaps != nil && len(swaps.Data.Swaps) > 1 {
		first, _ := priceOfSwap(swaps.Data.Swaps[0])
		second, _ := priceOfSwap(swaps.Data.Swaps[1])
		state = (first-second)/second > change
	}
	return state
}

func minMax(swaps utils.Swaps) (
	min float64,
	max float64,
	minTarget string,
	maxTarget string,
	minTime time.Time,
	maxTime time.Time,
) {
	min = 0
	max = 0
	var _min int64
	var _max int64
	for _, item := range swaps.Data.Swaps {
		price, target := priceOfSwap(item)
		minTarget = target
		maxTarget = target
		if min == 0 || max == 0 {
			min = price
			max = price
		}
		if price < min {
			min = price
			_min, _ = strconv.ParseInt(item.Timestamp, 10, 64)
		}
		if price > max {
			max = price
			_max, _ = strconv.ParseInt(item.Timestamp, 10, 64)
		}
	}
	minTime = time.Unix(_min, 0)
	maxTime = time.Unix(_max, 0)
	return min, max, minTarget, maxTarget, minTime, maxTime
}

func howMuchOld(swaps utils.Swaps) float64 {
	latest, _ := strconv.ParseInt(swaps.Data.Swaps[0].Timestamp, 10, 64)
	end := time.Unix(latest, 0)
	now := time.Now()
	period := now.Sub(end)
	return period.Hours()
}

func averageOfSwaps(swaps []utils.Swap) float64 {
	sum := 0.0
	for _, item := range swaps {
		price, _ := priceOfSwap(item)
		sum += price
	}
	return sum / float64(len(swaps))
}

func checkupOfSwaps(swaps utils.Swaps) bool {
	avg := 0.0
	duration := 20
	checkUp := 0
	checkDown := 0
	amount := int(len(swaps.Data.Swaps) / duration)
	for i := 0; i < duration; i++ {
		start := i * amount
		end := (i + 1) * amount
		if i == duration-1 {
			end = len(swaps.Data.Swaps)
		}
		var temp []utils.Swap
		for j := start; j < end; j++ {
			temp = append(temp, swaps.Data.Swaps[j])
		}
		cavg := averageOfSwaps(temp)
		if cavg > avg {
			checkUp += 1
			avg = cavg
		} else {
			checkDown += 1
		}
	}
	return checkUp > 2*checkDown
}

func checkdownOfSwaps(swaps utils.Swaps) bool {
	avg := 10000000.0
	duration := 20
	checkUp := 0
	checkDown := 0
	amount := int(len(swaps.Data.Swaps) / duration)
	for i := 0; i < duration; i++ {
		start := i * amount
		end := (i + 1) * amount
		if i == duration-1 {
			end = len(swaps.Data.Swaps)
		}
		var temp []utils.Swap
		for j := start; j < end; j++ {
			temp = append(temp, swaps.Data.Swaps[j])
		}
		cavg := averageOfSwaps(temp)
		if cavg > avg {
			checkUp += 1
			avg = cavg
		} else {
			checkDown += 1
		}
	}
	return checkUp < 2*checkDown
}
