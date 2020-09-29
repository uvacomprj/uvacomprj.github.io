package commands

import (
	"log"
	"math"
	"regexp"
	"sort"
	"strings"
	"uva-bot/model"
	"uva-bot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type email struct {
	name    string
	message *tgbotapi.Message
	bot     *tgbotapi.BotAPI
}

type emailResult struct {
	MensagemErro  string
	Query         string
	Contatos      []*model.Contato
	ContatosFuzzy []*model.Contato
}

func (cmd *email) Name() string {
	return cmd.name
}

func (cmd *email) Execute() error {
	contatoRepository := model.GetContatoRepository()
	query := cmd.message.CommandArguments()

	var result emailResult

	if len(query) < 3 {
		result.MensagemErro = "Por favor, informe um nome para pesquisa com no mínimo 3 (três) caracteres."

	} else if len(query) > 50 {
		result.MensagemErro = "O nome para pesquisa não pode exceder 50 (cinquenta) caracteres. Por favor, tente novamente."

	} else {
		regexSpaces := regexp.MustCompile(`\s+`)
		normalizedQuery, err := utils.NormalizeFnid(query)
		querySplits := regexSpaces.Split(normalizedQuery, -1)

		if err != nil {
			return err
		}

		contatos, err := contatoRepository.FindByPredicateFn(func(contato *model.Contato) (bool, error) {
			if contato.IsAlias() {
				return false, nil
			}

			normalizedName, err := utils.NormalizeFnid(*contato.Nome)

			if err != nil {
				return false, err
			}

			nameSplits := regexSpaces.Split(normalizedName, -1)
			equivalent := false

			for _, queryWord := range querySplits {
				for _, nameWord := range nameSplits {
					equivalent = strings.Contains(nameWord, queryWord)
					if equivalent {
						break
					}
				}
				if !equivalent {
					break
				}
			}

			return equivalent, nil
		})
		if err != nil {
			return err
		}

		sort.Slice(contatos[:], func(i, j int) bool {
			return (*contatos[i].Nome) < (*contatos[j].Nome)
		})

		contatosFuzzy, err := contatoRepository.FindByPredicateFn(func(contato *model.Contato) (bool, error) {
			if contato.IsAlias() {
				return false, nil
			}

			normalizedName, err := utils.NormalizeFnid(*contato.Nome)

			if err != nil {
				return false, err
			}

			nameSplits := regexSpaces.Split(normalizedName, -1)
			equivalent := false

			for _, queryWord := range querySplits {
				for _, nameWord := range nameSplits {
					equivalent = fuzzy.Match(queryWord, nameWord)
					if equivalent {
						break
					}
				}
				if !equivalent {
					break
				}
			}

			if !equivalent {
				for _, queryWord := range querySplits {
					maxRank := int(math.Ceil(float64(len(queryWord)) / 2.0))
					for _, nameWord := range nameSplits {
						rank := fuzzy.LevenshteinDistance(queryWord, nameWord)
						log.Println(rank, queryWord, nameWord)
						equivalent = rank >= 0 && rank < maxRank
						if equivalent {
							break
						}
					}
					if !equivalent {
						break
					}
				}
			}

			return equivalent, nil
		})
		if err != nil {
			return err
		}

		sort.Slice(contatosFuzzy[:], func(i, j int) bool {
			return (*contatosFuzzy[i].Nome) < (*contatosFuzzy[j].Nome)
		})

		result.Query = query
		result.Contatos = contatos
		result.ContatosFuzzy = contatosFuzzy
	}

	template, err := utils.LoadHtmlTemplate("email")
	if err != nil {
		return err
	}

	var builder strings.Builder
	err = template.Execute(&builder, result)
	if err != nil {
		return err
	}

	return utils.SendHtmlMessage(cmd.bot, cmd.message, builder.String())
}

func newEmailFactory(name string) commandHandlerFactory {
	return func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (CommandHandler, error) {
		return &email{
			name:    name,
			message: message,
			bot:     bot,
		}, nil
	}
}
