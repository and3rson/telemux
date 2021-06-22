// echo is a bot that repeats whatever you tell him.
package main

import (
	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	mux := tm.CreateMux().
		AddHandler(tm.CreateHandler(
			tm.IsCommand("start"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello! I'm a simple bot who repeats everything you say. :)",
				))
			},
		)).
		AddHandler(tm.CreateHandler(
			tm.IsText(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"You said: "+u.Message.Text,
				))
			},
		)).
		AddHandler(tm.CreateHandler(
			tm.Any(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Uh-oh, I can't repeat that!",
				))
			},
		))
	for update := range updates {
		mux.Dispatch(&update)
	}
}
