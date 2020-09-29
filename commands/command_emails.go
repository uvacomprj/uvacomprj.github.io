package commands

import (
	"sort"
	"strings"

	"uva-bot/model"
	"uva-bot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type emailsData struct {
	Docentes []*model.Contato
	Cursos   []*model.Curso
}

type emails struct {
	name    string
	message *tgbotapi.Message
	bot     *tgbotapi.BotAPI
}

func (cmd *emails) Name() string {
	return cmd.name
}

func (cmd *emails) Execute() error {
	contatoRepository := model.GetContatoRepository()
	cursoRepository := model.GetCursoRepository()

	docentes, err := contatoRepository.FindByPredicateFn(func(contato *model.Contato) (bool, error) {
		return !contato.IsAlias() && contato.Tipo == "docente", nil
	})
	if err != nil {
		return err
	}

	sort.Slice(docentes[:], func(i, j int) bool {
		return (*docentes[i].Nome) < (*docentes[j].Nome)
	})

	cursos, err := cursoRepository.FindByPredicateFn(func(curso *model.Curso) (bool, error) {
		return !curso.IsAlias(), nil
	})
	if err != nil {
		return err
	}

	data := emailsData{
		Docentes: docentes,
		Cursos:   cursos,
	}

	template, err := utils.LoadHtmlTemplate("emails")
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

func newEmailsFactory(name string) commandHandlerFactory {
	return func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error) {
		return &emails{
			name:    name,
			message: message,
			bot:     bot,
		}, nil
	}
}
