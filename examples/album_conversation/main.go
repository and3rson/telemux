// album_conversation is a bot that allows users to upload & share photos.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tm "github.com/and3rson/telemux/v2"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Photo describes a submitted photo
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
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var photos []Photo

	mux := tm.NewMux().
		AddHandler(tm.NewConversationHandler(
			"upload_photo_dialog",
			tm.NewLocalPersistence(), // we could also use `tm.NewFilePersistence("db.json")` or `&gormpersistence.GORMPersistence(db)` to keep data across bot restarts
			tm.StateMap{
				"": {
					tm.NewHandler(tm.IsCommandMessage("add"), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Please send me your photo.",
						))
						u.PersistenceContext.SetState("upload_photo")
					}),
				},
				"upload_photo": {
					tm.NewHandler(tm.HasPhoto(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						data["photoID"] = u.Message.Photo[0].FileID
						u.PersistenceContext.SetData(data)
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Please enter photo description.",
						))
						u.PersistenceContext.SetState("enter_description")
					}),
					tm.NewHandler(tm.Not(tm.IsCommandMessage("cancel")), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Sorry, I only accept photos. Please try again!",
						))
					}),
				},
				"enter_description": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						data["photoDescription"] = u.Message.Text
						u.PersistenceContext.SetData(data)
						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Are you sure you want to save this photo?")
						msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
							tgbotapi.NewKeyboardButtonRow(
								tgbotapi.NewKeyboardButton("Yes"),
								tgbotapi.NewKeyboardButton("No"),
							),
						)
						bot.Send(msg)
						u.PersistenceContext.SetState("confirm_submission")
					}),
					tm.NewHandler(tm.Not(tm.IsCommandMessage("cancel")), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Sorry, I did not understand that. Please enter some text!",
						))
					}),
				},
				"confirm_submission": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						data := u.PersistenceContext.GetData()
						var msg tgbotapi.MessageConfig
						if u.Message.Text == "Yes" {
							lastID++
							photos = append(photos, Photo{
								lastID,
								data["photoID"].(string),
								data["photoDescription"].(string),
							})
							msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Photo submitted! Type /list to list all photos.")
						} else {
							msg = tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled.")
						}
						u.PersistenceContext.ClearData()
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)
						u.PersistenceContext.SetState("")
					}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("cancel"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
					bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Cancelled."))
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
				for i, photo := range photos {
					if fmt.Sprint(photo.ID) == photoID {
						match = &photos[i]
					}
				}
				if match == nil {
					bot.Send(tgbotapi.NewMessage(
						u.Message.Chat.ID,
						"Photo not found!",
					))
				} else {
					share := tgbotapi.NewPhoto(u.Message.Chat.ID, tgbotapi.FileID(match.FileID))
					share.Caption = fmt.Sprintf("Description: %s", match.Description)
					_, err := bot.Send(share)
					if err != nil {
						log.Printf("main: send file: %s", err)
					}
				}
			},
		)).
		AddHandler(tm.NewMessageHandler(
			nil,
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
		mux.Dispatch(bot, update)
	}
}
