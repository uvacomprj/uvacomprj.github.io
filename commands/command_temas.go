package commands

import (
	"sort"
	"strings"
	"uva-bot/model"
	"uva-bot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type temasResult struct {
	TemaPai             *model.Tema
	TemasRaiz           []*model.Tema
	TemasFilhos         []*model.Tema
	ContagemTemasFilhos map[string]int64
}

type temas struct {
	name    string
	message *tgbotapi.Message
	bot     *tgbotapi.BotAPI
}

func (cmd *temas) Name() string {
	return cmd.name
}

func (cmd *temas) Execute() error {
	args := strings.TrimSpace(cmd.message.CommandArguments())

	data := temasResult{
		ContagemTemasFilhos: map[string]int64{},
	}

	temaRepository := model.GetTemaRepository()

	if len(args) == 0 {
		temasRaiz, err := temaRepository.FindTemasRaiz()

		if err != nil {
			return err
		}

		sort.Slice(temasRaiz[:], func(i, j int) bool {
			return (*temasRaiz[i].Nome) < (*temasRaiz[j].Nome)
		})

		for _, temaRaiz := range temasRaiz {
			numeroFilhos, err := temaRepository.CountTemasFilhos(temaRaiz)

			if err != nil {
				return err
			}

			data.ContagemTemasFilhos[temaRaiz.Fnid] = numeroFilhos
		}

		data.TemasRaiz = temasRaiz

	} else {
		temaPai, err := temaRepository.FindByFnid(args)

		if err != nil {
			return err
		}

		temasFilhos, err := temaRepository.FindTemasFilhos(temaPai)

		if err != nil {
			return err
		}

		sort.Slice(temasFilhos[:], func(i, j int) bool {
			return (*temasFilhos[i].Nome) < (*temasFilhos[j].Nome)
		})

		data.TemaPai = temaPai
		data.TemasFilhos = temasFilhos
	}

	template, err := utils.LoadHtmlTemplate("temas")
	if err != nil {
		return err
	}

	var builder strings.Builder
	err = template.Execute(&builder, data)
	if err != nil {
		return err
	}

	return utils.SendHtmlMessage(cmd.bot, cmd.message, builder.String())
}

func newTemasFactory(name string) commandHandlerFactory {
	return func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error) {
		return &temas{
			name:    name,
			message: message,
			bot:     bot,
		}, nil
	}
}
