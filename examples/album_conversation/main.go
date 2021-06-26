// album_conversation is a bot that allows users to upload & share photos.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	var photos []Photo
	mux := tm.NewMux().
		AddHandler(tm.NewConversationHandler(
			"upload_photo_dialog",
			tm.NewLocalPersistence(), // we could also use `tm.NewFilePersistence("db.json"),` to keep data across bot restarts
			map[string][]*tm.TransitionHandler{
				"": {
					tm.NewTransitionHandler(tm.IsCommandMessage("add"), func(u *tm.Update, data tm.Data) string {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Please send me your photo.",
						))
						return "upload_photo"
					}),
				},
				"upload_photo": {
					tm.NewTransitionHandler(tm.HasPhoto(), func(u *tm.Update, data tm.Data) string {
						data["photoID"] = (*u.Message.Photo)[0].FileID
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Please enter photo description.",
						))
						return "enter_description"
					}),
					tm.NewTransitionHandler(tm.Not(tm.IsAnyCommandMessage()), func(u *tm.Update, data tm.Data) string {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Sorry, I only accept photos. Please try again!",
						))
						return "upload_photo"
					}),
				},
				"enter_description": {
					tm.NewTransitionHandler(tm.HasText(), func(u *tm.Update, data tm.Data) string {
						data["photoDescription"] = u.Message.Text
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Are you sure you want to save this photo?")
						msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
							tgbotapi.NewKeyboardButtonRow(
								tgbotapi.NewKeyboardButton("Yes"),
								tgbotapi.NewKeyboardButton("No"),
							),
						)
						bot.Send(msg)
						return "confirm_submission"
					}),
				},
				"confirm_submission": {
					tm.NewTransitionHandler(tm.HasText(), func(u *tm.Update, data tm.Data) string {
						var msg tgbotapi.MessageConfig
						if u.Message.Text == "Yes" {
							lastID += 1
							photos = append(photos, Photo{
								lastID,
								data["photoID"].(string),
								data["photoDescription"].(string),
							})
							msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Photo submitted! Type /list to list all photos.")
						} else {
							msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled.")
						}
						for k := range data {
							delete(data, k)
						}
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)
						return ""
					}),
				},
			},
			[]*tm.TransitionHandler{
				tm.NewTransitionHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update, data tm.Data) string {
					for k := range data {
						delete(data, k)
					}
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled."))
					return ""
				}),
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("list"),
			func(u *tm.Update) {
				var lines []string
				for _, photo := range photos {
					lines = append(lines, fmt.Sprintf("- %s (/view_%d)", photo.Description, photo.ID))
				}
				if len(lines) == 0 {
					lines = append(lines, "No photos yet.")
				}
				message := tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Photos:\n"+strings.Join(lines, "\n"),
				)
				message.ReplyToMessageID = u.Message.MessageID
				bot.Send(message)
			},
		)).
		AddHandler(tm.NewHandler(
			tm.HasRegex(`^/view_(\d+)$`),
			func(u *tm.Update) {
				photoID := strings.Split(u.Message.Text, "_")[1]
				var match *Photo
				for _, photo := range photos {
					if fmt.Sprint(photo.ID) == photoID {
						match = &photo
					}
				}
				if match == nil {
					bot.Send(tgbotapi.NewMessage(
						u.Message.Chat.ID,
						"Photo not found!",
					))
				} else {
					share := tgbotapi.NewPhotoShare(u.Message.Chat.ID, match.FileID)
					share.Caption = fmt.Sprintf("Description: %s", match.Description)
					bot.Send(share)
				}
			},
		)).
		AddHandler(tm.NewHandler(
			tm.Any(),
			func(u *tm.Update) {
				message := tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello! I'm a gallery bot.\n\nI allow users to upload & share their photos!\n\nAvailable commands:\n/add - add photo\n/list - list photos",
				)
				message.ReplyToMessageID = u.Message.MessageID
				bot.Send(message)
			},
		))
	for update := range updates {
		mux.Dispatch(update)
	}
}
