package models

type CryptoTicker struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	Volume    string `json:"volume"`
}

type CryptoTickerRequest struct {
	Ticker string `json:"ticker"`
}
