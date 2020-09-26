package commands

import (
	"uva-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type emails struct {
	name string
	message *tgbotapi.Message
	bot *tgbotapi.BotAPI
}

func (cmd *emails) Name() string {
	return cmd.name
}

func (cmd *emails) Execute() (error) {
	return utils.SendMessageFromFile(cmd.bot, cmd.message, "emails.html")
}

func newEmailsFactory(name string) (commandHandlerFactory) {
	return func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error) {
		return &emails{
			name: name,
			message: message,
			bot: bot,
		}, nil
	}
}
