package utils

import (
	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CheckError checks for errors
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func SendMessageFromFile(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	templateName string,
) (error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.ParseMode = "html"
	filename := "./templates/" + templateName
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	msg.Text = string(text)
	if len(msg.Text) > 0 {
		bot.Send(msg)
	}
	return nil
}