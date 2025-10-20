package binance

import (
	"ExchangeInfo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *BinanceClient) BinanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed for this request", http.StatusMethodNotAllowed)
	}

	var request models.CryptoTickerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Ticker == "" {
		http.Error(w, "ticker required", http.StatusBadRequest)
	}

	ticker, err := c.GetTicker(request)
	if err != nil {
		http.Error(w, fmt.Sprintf("error during fetching ticker: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.CryptoTicker{
		Symbol:    ticker.Symbol,
		LastPrice: ticker.LastPrice,
		Volume:    ticker.Volume,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
