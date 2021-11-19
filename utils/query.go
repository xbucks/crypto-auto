package utils

import (
	"fmt"
)

type Bundles struct {
	Data struct {
		Bundles []struct {
			EthPrice string `json:"ethPrice"`
		} `json:"bundles"`
	} `json:"data"`
}

type Tokens struct {
	Data struct {
		Tokens []struct {
			Id             string `json:"id"`
			Name           string `json:"name"`
			Symbol         string `json:"symbol"`
			DerivedETH     string `json:"derivedETH"`
			TotalLiquidity string `json:"totalLiquidity"`
		} `json:"tokens"`
	} `json:"data"`
}

type Swap struct {
	Amount0In  string `json:"amount0In"`
	Amount0Out string `json:"amount0Out"`
	Amount1In  string `json:"amount1In"`
	Amount1Out string `json:"amount1Out"`
	AmountUSD  string `json:"amountUSD"`
	Id         string `json:"id"`
	Pair       struct {
		Token0 struct {
			Symbol     string `json:"symbol"`
			Name       string `json:"name"`
			DerivedETH string `json:"derivedETH"`
		} `json:"token0"`
		Token1 struct {
			Symbol     string `json:"symbol"`
			Name       string `json:"name"`
			DerivedETH string `json:"derivedETH"`
		} `json:"token1"`
	} `json:"pair"`
	Timestamp string `json:"timestamp"`
	To        string `json:"to"`
}

type Swaps struct {
	Data struct {
		Swaps []Swap
	}
}

type Pairs struct {
	Data struct {
		Pairs []struct {
			Id     string `json:"id"`
			Token0 struct {
				Symbol string `json:"symbol"`
			} `json:"token0"`
			Token1 struct {
				Symbol string `json:"symbol"`
			} `json:"token1"`
			Token0Price string `json:"token0Price"`
			Token1Price string `json:"token1Price"`
		}
	}
}

func Query(target string, limit int, id string) map[string]string {
	var query map[string]string
	switch target {
	case "bundles":
		query = map[string]string{
			"query": `
				query bundles {
					bundles(where: { id: "1" }) {
						ethPrice
					}
				}
			`,
		}
	case "tokens":
		sub := fmt.Sprintf(`
			query tokens {
				tokens(where: { id: "%s" }) {
					id
					name
					symbol
					derivedETH
					totalLiquidity
				}
			}
		`, id)
		query = map[string]string{"query": sub}
	case "swaps":
		sub := fmt.Sprintf(`
			query swaps {
				swaps(first: %d, orderBy: timestamp, orderDirection: desc, where:
					{ pair: "%s" }
				) {
					pair {
						token0 {
							symbol
							name
       						derivedETH
						}
						token1 {
							symbol
							name
       						derivedETH
						}
					}
					amount0In
					amount0Out
					amount1In
					amount1Out
					amountUSD
					to
					timestamp
					id
				}
			}
		`, limit, id)
		query = map[string]string{"query": sub}
	case "pairs":
		query = map[string]string{
			"query": `
				query pairs {
					pairs(first: 1000, orderBy: reserveUSD, orderDirection: desc) {
						id,
						token0 {
							symbol
						},
						token1 {
							symbol
						},
						token0Price,
						token1Price,
					}
				}
			`,
		}
	default:
	}
	return query
}
