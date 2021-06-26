// private_only is a bot that allows you to talk to him only in private chat.
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
	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.And(tm.IsPrivate(), tm.IsCommandMessage("start")),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Psst... Don't tell anyone about our private chat! :)",
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Sorry, I only respond in private chats. Send me a direct message!",
				))
			},
		))
	for update := range updates {
		mux.Dispatch(&update)
	}
}
