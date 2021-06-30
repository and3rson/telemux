// error_handling is a bot that handles zero division panic.
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {
				msg := tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello! I divide numbers. For example: `/div 20 4`.\n\nHint:  I can handler errors! Try `/div 42 0`",
				)
				msg.ParseMode = "markdown"
				bot.Send(msg)
			},
		)).
		AddHandler(tm.NewHandler(
			tm.And(tm.IsMessage(), tm.HasRegex(`^/div (\d+) (\d+)$`)),
			func(u *tm.Update) {
				parts := strings.Split(u.Message.Text, " ")
				a, _ := strconv.Atoi(parts[1])
				b, _ := strconv.Atoi(parts[2])
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf("The result is %d", a/b),
				))
			},
		)).
		SetRecoverer(func(u *tm.Update, err error) {
			chat := u.EffectiveChat()
			if chat != nil {
				bot.Send(tgbotapi.NewMessage(
					chat.ID,
					fmt.Sprintf("Oops, an error occurred: %s", err),
				))
				log.Printf("Warning! An error occurred: %s", err)
			}
		})
	for update := range updates {
		mux.Dispatch(bot, update)
	}
}
