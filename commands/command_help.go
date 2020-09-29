package commands

import (
	"uva-bot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type help struct {
	name    string
	message *tgbotapi.Message
	bot     *tgbotapi.BotAPI
}

func (cmd *help) Name() string {
	return cmd.name
}

func (cmd *help) Execute() error {
	return utils.SendMessageFromFile(cmd.bot, cmd.message, "help.html")
}

func newHelpFactory(name string) commandHandlerFactory {
	return func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error) {
		return &help{
			name:    name,
			message: message,
			bot:     bot,
		}, nil
	}
}
