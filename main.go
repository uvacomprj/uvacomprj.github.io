package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	CheckError(err)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "html"
			switch update.Message.Command() {
			case "readme":
				msg.Text = ReadTextFile("./Messages/readme.html")

			case "emails":
				msg.Text = ReadTextFile("./Messages/emails.html")

			case "grupos":
				msg.Text = ReadTextFile("./Messages/groups.html")

			default:
				msg.Text = ""
			}
			if len(msg.Text) > 0 {
				bot.Send(msg)
			}
		}

	}
}
