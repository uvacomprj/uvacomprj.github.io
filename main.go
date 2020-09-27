package main

import (
	"log"
	"os"

	"uva-bot/commands"
	"uva-bot/model"
	"uva-bot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	model.Init()
	commands.Init()

	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	utils.CheckError(err)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	utils.CheckError(err)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			// Para adicionar novos comandos, crie um novo arquivo contendo uma implementação da
			// interface CommandHandler, no pacote "commands". Em seguida, cadastre um factory
			// method para o comando criado no método "commands.Init" (arquivo "commands/commands.go")
			command, err := commands.NewCommandHandler(bot, update.Message)

			if err != nil {
				log.Printf("[%s] %s", update.Message.From.UserName, err.Error())
			} else {
				log.Printf("[%s] Executing %s", update.Message.From.UserName, command.Name())
				execErr := command.Execute()

				if execErr != nil {
					log.Println(execErr)
				}
			}
		}
	}
}
