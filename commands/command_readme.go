package commands

import (
	"uva-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type readme struct {
	name string
	message *tgbotapi.Message
	bot *tgbotapi.BotAPI
}

func (cmd *readme) Name() string {
	return cmd.name
}

func (cmd *readme) Execute() (error) {
	return utils.SendMessageFromFile(cmd.bot, cmd.message, "readme.html")
}

func newReadmeFactory(name string) (commandHandlerFactory) {
	return func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error) {
		return &readme{
			name: name,
			message: message,
			bot: bot,
		}, nil
	}
}
