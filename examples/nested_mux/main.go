// nested_mux is a bot that demonstrates nested Mux usage
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	// Register two sub-Mux instances. First will handle updates in private chats, second - in group chats.
	mux := tm.NewMux().
		AddMux(tm.NewMux().
			SetGlobalFilter(tm.IsPrivate()).
			AddHandler(tm.NewHandler(tm.IsCommandMessage("start"), func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello!\n\nCommands in private chat:\n- /start - Show info\n- /version - Print my version\n- /cheer - Send a happy message\n\nCommands in group chats:\n- /time - Tell current time",
				))
			})).
			AddHandler(tm.NewHandler(tm.IsCommandMessage("version"), func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"My version is 0.0.0-alpha",
				))
			})).
			AddHandler(tm.NewHandler(tm.IsCommandMessage("cheer"), func(u *tm.Update) {
				bot.Send(tgbotapi.NewSticker(u.EffectiveChat().ID, tgbotapi.FileID("CAACAgIAAxkBAAECg_1g3b2j0AHBrbm0zPxlkWGDxoYq7QACsQADwPsIAAED7avN0x5kmSAE")))
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"PRT HRD!",
				))
			})),
		).
		AddMux(tm.NewMux().
			SetGlobalFilter(tm.IsGroupOrSuperGroup()).
			AddHandler(tm.NewHandler(tm.IsCommandMessage("time"), func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf("The time is %s", time.Now()),
				))
			})),
		).
		AddHandler(tm.NewHandler(tm.Any(), func(u *tm.Update) {
			bot.Send(tgbotapi.NewMessage(
				u.Message.Chat.ID,
				"Sorry, I can't do that.",
			))
		}))
	for update := range updates {
		mux.Dispatch(bot, update)
	}
}
