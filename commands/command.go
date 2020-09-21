package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type commandHandlerFactory func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error)

type CommandHandler interface {
	Name() string
	Execute() (error)
}
