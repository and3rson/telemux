// echo is a bot that repeats whatever you tell him.
package main

import (
	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello! I'm a simple bot who repeats everything you say. :)",
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.HasText(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"You said: "+u.Message.Text,
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.Any(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Uh-oh, I can't repeat that!",
				))
			},
		))

	for update := range updates {
		mux.Dispatch(bot, update)
	}
}
