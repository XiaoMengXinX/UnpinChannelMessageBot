package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/telebot.v4"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable not set")
	}

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	botUsername := bot.Me.Username
	log.Printf("Bot @%s started", botUsername)

	_ = bot.RemoveWebhook()
	err = bot.SetDefaultRights(telebot.Rights{CanPinMessages: true}, false)
	if err != nil {
		log.Fatal("Failed to set default rights.")
	}

	err = bot.SetCommands([]telebot.Command{
		{
			Text:        "start",
			Description: "Start to use the bot",
		},
	})
	if err != nil {
		log.Fatal("Failed to set commands:", err)
	}

	bot.Handle("/start", func(c telebot.Context) error {
		return handleStartCommand(c, botUsername)
	})
	bot.Handle(&telebot.InlineButton{Unique: "check_permissions"}, handleCheckPermissions)
	bot.Handle(telebot.OnText, handleMessage)
	bot.Handle(telebot.OnPoll, handleMessage)
	bot.Handle(telebot.OnMedia, handleMessage)
	bot.Handle(telebot.OnAddedToGroup, func(c telebot.Context) error {
		return handleAddedToGroup(c, botUsername)
	})

	go bot.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Exiting...")
	bot.Stop()
}

func handleStartCommand(c telebot.Context, botUsername string) error {
	if c.Chat().Type != telebot.ChatPrivate {
		return nil
	}

	message := `Hello! I'm a bot for unpinning channel messages.
Click the button below to add me to your group and grant me admin permissions.`

	keyboard := &telebot.ReplyMarkup{}
	btnAddToGroup := keyboard.URL("Add me to a group", "https://t.me/"+botUsername+"?startgroup=true")

	keyboard.Inline(
		keyboard.Row(btnAddToGroup),
	)

	return c.Send(message, keyboard)
}

func handleAddedToGroup(c telebot.Context, botUsername string) error {
	if c.Chat().Type != telebot.ChatGroup && c.Chat().Type != telebot.ChatSuperGroup {
		return nil
	}

	message := `👋 Thank you for adding me to the group!
To work properly, I need admin permissions to unpin messages.
You can click the button to check permission`

	keyboard := &telebot.ReplyMarkup{}
	keyboard.Inline(
		keyboard.Row(keyboard.Data("✅ Check Permission", "check_permissions")),
	)
	return c.Send(message, keyboard)
}

// Handle all messages
func handleMessage(c telebot.Context) error {
	message := c.Message()

	if message.SenderChat != nil && message.SenderChat.Type == telebot.ChatChannel {
		err := c.Bot().Unpin(c.Chat(), message.ID)
		if err != nil {
			log.Printf("Failed to unpin: %v", err)
			return nil
		}
	}
	return nil
}

func handleCheckPermissions(c telebot.Context) error {
	member, err := c.Bot().ChatMemberOf(c.Chat(), c.Message().Sender)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "❌ Failed to check permissions.",
		})
	}

	if member.CanPinMessages {
		return c.Respond(&telebot.CallbackResponse{
			Text: "✅ Permissions set successfully!",
		})
	} else {
		return c.Respond(&telebot.CallbackResponse{
			Text: "❌ I don't have pin message permissions yet.",
		})
	}
}
