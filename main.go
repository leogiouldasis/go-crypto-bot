package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	client := binance.NewClient(apiKey, secretKey)

	// Check account balance
	balances, err := getBalances(client)
	if err != nil {
		log.Fatalf("Error retrieving balances: %v", err)
	}

	for _, balance := range balances {
		fmt.Printf("%s: %s\n", balance.Asset, balance.Free)
	}

	top10CryptosByMarketCap, err := getTopCryptosByMarketCap(10)
	if err != nil {
		log.Fatalf("Error retrieving top 10 cryptocurrencies by market cap: %v", err)
	}

	top10Symbols := make(map[string]bool)
	for _, coin := range top10CryptosByMarketCap {
		top10Symbols[coin.Symbol] = true
	}

	// Get the top 10 cryptocurrencies by percentage move against Bitcoin
	topCryptos, err := getTopCryptosByPercentageMove(client, 10, top10Symbols)
	if err != nil {
		log.Fatalf("Error retrieving top cryptocurrencies: %v", err)
	}

	fmt.Println("Top 10 cryptocurrencies by percentage move against Bitcoin:")
	for _, crypto := range top10Cryptos {
		fmt.Printf("%s: %.2f%%\n", crypto.Symbol, crypto.PercentageMove)
	}
}

func getBalances(client *binance.Client) ([]binance.Balance, error) {
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	return account.Balances, nil
}

type CryptoMove struct {
	Symbol         string
	PercentageMove float64
}

func getTopCryptosByPercentageMove(client *binance.Client, top int, top10Symbols map[string]bool) ([]CryptoMove, error) {
	ticker24h, err := client.NewListPriceChangeStatsService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	var moves []CryptoMove

	for _, ticker := range ticker24h {
		baseSymbol := strings.TrimSuffix(ticker.Symbol, "BTC")
		if top10Symbols[baseSymbol] {
			priceChangePercent, err := strconv.ParseFloat(ticker.PriceChangePercent, 64)
			if err != nil {
				return nil, err
			}
			moves = append(moves, CryptoMove{Symbol: ticker.Symbol, PercentageMove: priceChangePercent})
		}
	}

	sort.SliceStable(moves, func(i, j int) bool {
		return moves[i].PercentageMove > moves[j].PercentageMove
	})

	if top > len(moves) {
		top = len(moves)
	}

	return moves[:top], nil
}


func getTopCryptosByMarketCap(count int) ([]gecko.Coin, error) {
	geckoClient := gecko.NewClient(nil)
	coins, err := geckoClient.CoinsMarkets("usd", gecko.WithPerPage(count), gecko.WithOrder(gecko.MarketCapDesc))
	if err != nil {
		return nil, err
	}
	return coins, nil
}




