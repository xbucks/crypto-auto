package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hirokimoto/crypto-auto/utils"
)

const TRADES_TARGET string = "/trades.txt"
const ALL_PAIRS string = "/allpairs.txt"

func SaveTradables(tokens *Tokens) {
	path := absolutePath() + TRADES_TARGET
	trades, err := readLines(path)
	if err != nil {
		return
	}
	for _, v := range tokens.data {
		if !isExistedTradables(v.address, trades) {
			trades = append(trades, v.address)
		}
	}
	writeLines(trades, TRADES_TARGET)
	fmt.Println("Saved tradable tokens successfully!")
}

func isExistedTradables(t string, trades []string) bool {
	for _, v := range trades {
		if v == t {
			return true
		}
	}
	return false
}

func SaveAllPairs(p *utils.Pairs) {
	path := absolutePath() + ALL_PAIRS
	pairs, _ := readLines(path)

	for _, v := range p.Data.Pairs {
		if !isExistedPairs(v.Id, pairs) {
			pairs = append(pairs, v.Id)
		}
	}

	writeLines(pairs, path)
}

func ReadAllPairs() ([]string, error) {
	path := absolutePath() + ALL_PAIRS
	pairs, err := readLines(path)
	return pairs, err
}

func isExistedPairs(p string, pairs []string) bool {
	for _, v := range pairs {
		if v == p {
			return true
		}
	}
	return false
}

func WriteOnePair(pair string) error {
	err := writeOnePair(pair)
	return err
}

func RemoveOnePair(pair string) error {
	err := removeOnePair(pair)
	return err
}

func writeOnePair(pair string) error {
	path := absolutePath() + "/pairs.txt"
	pairs, _ := readLines(path)
	pairs = append(pairs, pair)
	err := writeLines(pairs, path)
	return err
}

func removeOnePair(pair string) error {
	path := absolutePath() + "/pairs.txt"
	pairs, _ := readLines(path)
	_pairs := []string{}
	for _, v := range pairs {
		if v != pair {
			_pairs = append(_pairs, v)
		}
	}
	fmt.Println(pairs)
	fmt.Println(_pairs)
	err := writeLines(_pairs, path)
	return err
}

func absolutePath() string {
	// ex, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }
	// exPath := filepath.Dir(ex)
	// return exPath
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func readLines(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		writeLines(lines, path)
		return lines, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
