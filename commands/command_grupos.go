package commands

import (
	"uva-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type grupos struct {
	name string
	message *tgbotapi.Message
	bot *tgbotapi.BotAPI
}

func (cmd *grupos) Name() string {
	return cmd.name
}

func (cmd *grupos) Execute() (error) {
	return utils.SendMessageFromFile(cmd.bot, cmd.message, "grupos.html")
}

func newGruposFactory(name string) (commandHandlerFactory) {
	return func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error) {
		return &grupos{
			name: name,
			message: message,
			bot: bot,
		}, nil
	}
}
