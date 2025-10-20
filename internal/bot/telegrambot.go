package telegrambot

import (
	"ExchangeInfo/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	telegrambotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Bot    *telegrambotapi.BotAPI
	ApiURL string
	Client *http.Client
}

func NewTelegramBot(token, apiURL string) (*TelegramBot, error) {
	bot, err := telegrambotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %v", err)
	}
	return &TelegramBot{
		Bot:    bot,
		ApiURL: apiURL,
		Client: &http.Client{},
	}, nil
}

func (tb *TelegramBot) Start() {
	u := telegrambotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := telegrambotapi.NewMessage(update.Message.Chat.ID, "")
		command := update.Message.Command()
		args := strings.TrimSpace(update.Message.CommandArguments())

		if command == "ticker" {
			if args == "" {
				msg.Text = "Please provide a ticker symbol(e.g., /ticker BTCUSDT )"
				tb.Bot.Send(msg)
				continue
			}
			//prepare request for /Cryptoticker endpoint
			request := models.CryptoTickerRequest{Ticker: args}
			reqBody, err := json.Marshal(request)
			if err != nil {
				msg.Text = "Internal error: failed to prepare request"
				tb.Bot.Send(msg)
				continue
			}
			//отправляет POST запрос
			resp, err := tb.Client.Post(tb.ApiURL+"/Cryptoticker", "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				msg.Text = fmt.Sprintf("Error connecting to API: %v", err)
				tb.Bot.Send(msg)
				continue
			}
			defer resp.Body.Close()

			//обрабатывает апишку
			if resp.StatusCode != http.StatusOK {
				msg.Text = fmt.Sprintf("Error handling API: %v", resp.StatusCode)
				tb.Bot.Send(msg)
				continue
			}

			var ticker models.CryptoTicker
			if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
				msg.Text = "Error parsing API response"
				tb.Bot.Send(msg)
				continue
			}
			msg.Text = fmt.Sprintf("Ticker: %s\nLastPrice: %s\nVolume: %s", ticker.Symbol, ticker.LastPrice, ticker.Volume)
			tb.Bot.Send(msg)
		} else {
			msg.Text = "Unknown command. Use /ticker <SYMBOL> (e.g., /ticker BTCUSDT)"
			tb.Bot.Send(msg)
		}
	}
}
