package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("826697982:AAGMxWYofeRZ5pZl9Wb_pcPBjbmZ3xETDO0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "html"
			switch update.Message.Command() {
			case "readme":
				msg.Text =
					`<strong>IMPORTANTE — informações úteis</strong>:

				<strong>README</strong> (um pouco defasado no momento): https://telegra.ph/README-12-19
				<strong>Google Drive</strong>: https://drive.google.com/drive/folders/0B5eVrMuAr6lbUzhya09SQ05fckU
				<strong>Grupo no LinkedIn</strong>: https://www.linkedin.com/groups/8655283
				<strong>URL "amigável" deste próprio grupo de Telegram</strong>: http://bit.ly/chat-uvacc`

			case "emails":
				msg.Text =
					`<strong>Adriana Aparício Sicsú Ayres do Nascimento</strong> adriana.nascimento@uva.br
			<strong>Ana Maria</strong> ana.vianna@uva.br
			<strong>André Lucio de Oliveira</strong> andre.oliveira@uva.br
			<strong>Alfredo Boente</strong> professor@boente.eti.br
			<strong>Camilla Lobo Paulino</strong> profcamilla@uol.com.br
			<strong>Carlos Alberto Alves Lemos</strong> caalemos@ig.com.br
			<strong>Carlos Augusto Sicsú Ayres Nascimento</strong> caugusto.sicsu@uva.br
			<strong>Carlos Frederico Motta Vasconcelos</strong> cfmotta_2001@yahoo.com
			<strong>Claudio Jose Marques de Souza</strong> cjms70@gmail.com
			<strong>Douglas Ericsson</strong> ericsonmarc@gmail.com
			<strong>Edgar Gurgel</strong> edgar@uva.br edgar@cos.ufrj.br
			<strong>Elias Restum Antônio</strong> eliasra@globo.com
			<strong>Eliane Xavier Cavalcanti</strong> eliane@cnen.gov.br
			<strong>Jobson Luiz</strong> jobson.silva@uva.br
			<strong>Matheus Bandini</strong> matheusbandini@gmail.com
			<strong>Miguel Figueiredo</strong> miguel.figueiredo@uva.br miguel.azf@gmail.com
			<strong>Roberto Luís Miranda Pereira de Castro</strong> roberto.castro@uva.br
			<strong>Rossandro Ramos</strong> prof.rossandro@gmail.com
			<strong>Thiago Gabriel</strong> thiago.gabriel@uva.br
			<strong>Vincenzo</strong> vdroberto@gmail.com`

			case "grupos":
				msg.Text =
					`<strong>AOO</strong> (Análise Orientada a Objetos): https://t.me/joinchat/AvwFpUiH4s8M2STQjUsYzQ
			
			<strong>Aplicações na Internet</strong>: https://t.me/joinchat/AvwFpU0WeWw_BBK-4dQEjg
			
			<strong>Arquitetura de Computadores</strong>: https://t.me/joinchat/AvwFpVYEYrxirrsx71620w
			
			<strong>Bancos de Dados (BD1 e BD2)</strong>: https://t.me/joinchat/AvwFpRfAyL0qxsi2WJ-o5g
			
			<strong>Computação Gráfica</strong> (CG): https://t.me/joinchat/AvwFpU34q8ljFiGfdyWo_A
			
			<strong>Conceitos de Linguagens de Programação</strong> (CLP): https://t.me/joinchat/AvwFpUYM3nYX7-V4fT3SzA
			
			<strong>CVGA</strong> (Cálculo Vetorial e Geometria Analítica): https://t.me/joinchat/AvwFpRbTkP9Iw8uUspZt9Q
			
			<strong>Engenharia de Software</strong>: https://t.me/joinchat/Bj0snhDQSrpKjYdNZ5Pfag
			
			<strong>Estatística</strong>: https://t.me/joinchat/AvwFpQlAzWxR2V5_R1f7cw
			
			<strong>Estruturas de Dados</strong>: https://t.me/joinchat/AvwFpVf0Ol7c-W8nBhZsCQ
			
			<strong>Filosofia</strong>: https://t.me/joinchat/AvwFpRH9KyBcd2GWRsDpiA
			
			<strong>Gerência de Projetos</strong>: https://t.me/joinchat/AvwFpUd6M0WqZgEovlB-2Q
			
			<strong>Gestão de Negócios e Ética</strong> (GNE): https://t.me/joinchat/AvwFpQtTraQcFjG9cnXY_A
			
			<strong>Governança Corporativa em TI</strong> (GCTI): https://t.me/joinchat/AvwFpRHa6v6nvD0TFiBMeA
			
			<strong>Inteligência Computacional</strong> (IC): https://t.me/joinchat/Chd0XUMeZXte7mOeonuNEA
			
			<strong>IHC</strong> (Interação Humano-Computador): https://t.me/joinchat/AvwFpVQuKJodFobcETEIVQ
			
			<strong>LFAC</strong> (Linguagens Formais, Autômatos e Computabilidade): https://t.me/joinchat/AvwFpQs0e8zA0JF2qd9A-A
			
			<strong>Matemática Discreta</strong> (MD1/MD2): https://t.me/joinchat/AvwFpRXgnrzjQBbBBd-dMA
			
			<strong>Modelagem Computacional</strong>: https://t.me/joinchat/AvwFpQqtUFAMxO-7XbvsZA
			
			<strong>Monografia</strong>: https://t.me/joinchat/CibjPA2Yk3T7im9Ih1xY8g
			
			<strong>Programação para Dispositivos Móveis</strong> (PDM): https://t.me/joinchat/AvwFpQrEkwOj1Ox223DYXw
			
			<strong>Programação para Jogos</strong>: https://t.me/joinchat/AvwFpQ8GUKpIogRAnc84ng
			
			<strong>Programação Paralela</strong>: https://t.me/joinchat/AvwFpQyMAEOz71Vvy3BS_Q
			
			<strong>PSOO</strong> (Projeto de Software Orientado a Objetos): https://t.me/joinchat/AvwFpQ492QDEhLVUHxiDXQ
			
			<strong>Projeto em Automação</strong>: https://t.me/joinchat/GM7jHBK7Vd7LwRHMUqFaxA
			
			<strong>Projeto e Análise de Algoritmos</strong>: https://t.me/joinchat/AvwFpQvtn2jkgAb6JJQo_g 
			
			<strong>Redes de Computadores</strong> (I e II): https://t.me/joinchat/AvwFpRR8Zp2hP72AGyodBg
			
			<strong>Sistemas Digitais</strong>: https://t.me/joinchat/AvwFpVlKX-0p-zM10oaGEw
			
			<strong>Sistemas Distribuídos</strong>: https://t.me/joinchat/AvwFpUv_XeRt0UAlgaLP8Q
			
			<strong>Sistemas Embarcados</strong>: https://t.me/joinchat/GM7jHBBWYEmKjhNwp1_MJQ
			
			<strong>Sistemas em Tempo Real</strong> (STR, RTOS): https://t.me/joinchat/AvwFpU2GGx76aPOexKzKuA
			
			<strong>Teoria dos Números</strong>: https://t.me/joinchat/Chd0XQ2u0uIh0HiFnb11Mw
			
			<strong>Tópicos Especiais em Compiladores</strong>: https://t.me/joinchat/AvwFpQwuw5Z4WLsz1pWcvQ`

			default:
				msg.Text = "Não conheço este comando!"
			}
			bot.Send(msg)
		}

	}
}
