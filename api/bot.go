package api

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"os"
)

type Response struct {
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	Method    string `json:"method"`
}

func BotHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bot := &tgbotapi.BotAPI{
		Token:  os.Getenv("BOT_TOKEN"),
		Client: &http.Client{},
		Buffer: 100,
	}
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)

	update, err := bot.HandleUpdate(r)
	if err != nil {
		return
	}

	var response []byte

	if update.Message.PinnedMessage != nil {
		message := update.Message.PinnedMessage
		fmt.Printf("%+v", update.Message)
		if message.SenderChat.Type == "channel" {
			data := Response{
				Method:    "sendMessage",
				ChatID:    message.Chat.ID,
				MessageID: message.MessageID,
			}
			response, _ = json.Marshal(data)
		}
	}
	w.Header().Add("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, string(response))
}
