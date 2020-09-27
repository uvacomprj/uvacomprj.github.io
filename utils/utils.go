package utils

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode"

	html "html/template"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//CheckError checks for errors
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

var templatesFolderPrefix = "./templates/"

func SendMessageFromFile(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	templateName string,
) error {
	filename := templatesFolderPrefix + templateName
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return SendHtmlMessage(bot, message, string(text))
}

func SendHtmlMessage(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	content string,
) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.ParseMode = "html"
	msg.Text = content
	if len(msg.Text) > 0 {
		_, err := bot.Send(msg)
		return err
	}
	return nil
}

func LoadHtmlTemplate(filename string) (*html.Template, error) {
	template, err := html.ParseFiles(templatesFolderPrefix + filename + ".html")
	if err != nil {
		return nil, err
	}
	files, err := filepath.Glob(templatesFolderPrefix + "common/*.html")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		template, err = template.ParseGlob(file)
		if err != nil {
			return nil, err
		}
	}
	return template, nil
}

func NormalizeFnid(input string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, input)
	if err != nil {
		return "", err
	}
	return strings.ToLower(result), nil
}
