package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"os"
)

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

	if message := update.Message.PinnedMessage; message != nil {
		if message.SenderChat.Type == "channel" {
			unpinMessage := tgbotapi.UnpinChatMessageConfig{
				ChatID:    message.Chat.ID,
				MessageID: message.MessageID,
			}
			_, _ = bot.Send(unpinMessage)
		}
	}
}
