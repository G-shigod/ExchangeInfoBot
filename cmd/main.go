package main

import (
	binance "ExchangeInfo/internal/BinanceXchange"
	telegrambot "ExchangeInfo/internal/bot"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("failed to download .env: %v", err)
	}

	binanceClient := binance.NewBinanceClient()

	telegramToken := os.Getenv("bot_api")
	if telegramToken == "" {
		log.Fatal("Telegrambot Token is missing")
	}

	telegramBot, err := telegrambot.NewTelegramBot(telegramToken, "http://localhost:8080")
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	go telegramBot.Start()

	http.HandleFunc("/Cryptoticker", binanceClient.BinanceHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
