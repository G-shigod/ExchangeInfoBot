package binance

import (
	"ExchangeInfo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type BinanceClient struct {
	client  *http.Client
	baseURL string
}

func NewBinanceClient() *BinanceClient {
	return &BinanceClient{
		client:  &http.Client{},
		baseURL: "https://api.binance.com",
	}
}

func (bc *BinanceClient) GetTicker(request models.CryptoTickerRequest) (*models.CryptoTicker, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/24hr?symbol=%s", bc.baseURL, request.Ticker)
	resp, err := bc.client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to make request to Binance API: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API called 200 error: %d", resp.StatusCode)
	}

	var ticker models.CryptoTicker
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		return nil, fmt.Errorf("failed to decode data from binance API: %v", err)
	}
	return &ticker, nil
}
