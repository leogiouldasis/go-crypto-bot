package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type Coin struct {
	Id             string  `json:"id"`
	Symbol         string  `json:"symbol"`
	Name           string  `json:"name"`
	MarketCap      float64 `json:"market_cap"`
	Price          float64 `json:"current_price"`
	PercentChange7d float64 `json:"price_change_percentage_7d_in_currency"`
}

type ByPercentChange7d []Coin

func (c ByPercentChange7d) Len() int           { return len(c) }
func (c ByPercentChange7d) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByPercentChange7d) Less(i, j int) bool { return c[i].PercentChange7d < c[j].PercentChange7d }

func main() {
	url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=20&page=1&sparkline=false&price_change_percentage=7d"
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data from CoinGecko API: %s\n", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	var coins []Coin
	err = json.Unmarshal(body, &coins)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %s\n", err)
		return
	}

	sort.Sort(sort.Reverse(ByPercentChange7d(coins)))

	fmt.Printf("Top 20 Cryptocurrencies by Market Cap, Sorted by 7d Change in USD:\n")
	for i, coin := range coins {
		fmt.Printf("%d. %s (%s): $%.2f, Market Cap: $%.0f, 7d Change: %.2f%%\n", i+1, coin.Name, coin.Symbol, coin.Price, coin.MarketCap, coin.PercentChange7d)
	}
}
