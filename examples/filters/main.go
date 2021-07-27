// filters is a bot that can receive photos, text messages & geolocations.
package main

import (
	"fmt"
	"log"
	"os"

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
	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello! Send me a text message, photo or geolocation.",
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.HasText(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"You sent me a text message: "+u.Message.Text,
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.HasPhoto(),
			func(u *tm.Update) {
				photo := u.Message.Photo[0]
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf("You sent me a photo of size %d x %d", photo.Width, photo.Height),
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.HasLocation(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf("You sent me a geolocation: %f;%f", u.Message.Location.Latitude, u.Message.Location.Longitude),
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.Any(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Sorry, I only accept text messages, photos & geolocations. :(",
				))
			},
		))
	for update := range updates {
		mux.Dispatch(bot, update)
	}
}
