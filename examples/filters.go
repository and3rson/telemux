// filters is a bot that can receive photos, text messages & geolocations.
package main

import (
	"fmt"
	"log"
	"os"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Photo struct {
	ID          int
	FileID      string
	Description string
}

var lastID = 0

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
			tm.IsCommand("start"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello! Send me a text message, photo or geolocation.",
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsText(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"You sent me a text message: "+u.Message.Text,
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsPhoto(),
			func(u *tm.Update) {
				photo := (*u.Message.Photo)[0]
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf("You sent me a photo of size %d x %d", photo.Width, photo.Height),
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsLocation(),
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
		mux.Dispatch(&update)
	}
}
