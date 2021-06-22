// echo is a bot that repeats whatever you tell him.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CatInfo struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func GetRandomCatURL() (string, error) {
	resp, err := http.Get("https://api.thecatapi.com/v1/images/search?mime_types=jpg")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	var catInfos []CatInfo
	if err := json.Unmarshal(data, &catInfos); err != nil {
		return "", err
	}
	return catInfos[0].URL, nil
}

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

	loading_markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Loading...", "refresh"),
		),
	)
	refresh_markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Another image", "refresh"),
		),
	)

	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsCommand("start"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					"Hello! Type /cat to display a picture of a random cat.",
				))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsCommand("cat"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewChatAction(u.Message.Chat.ID, tgbotapi.ChatTyping))
				url, err := GetRandomCatURL()
				if err != nil {
					bot.Send(tgbotapi.NewMessage(
						u.Message.Chat.ID,
						fmt.Sprintf("Oops, an error occured: %s", err),
					))
					return
				}
				message := tgbotapi.NewMessage(
					u.Message.Chat.ID,
					url,
				)
				message.ReplyMarkup = refresh_markup
				bot.Send(message)
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsCallbackQuery(),
			func(u *tgbotapi.Update) {
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Refreshing..."))
				bot.Send(tgbotapi.NewEditMessageReplyMarkup(
					u.CallbackQuery.Message.Chat.ID,
					u.CallbackQuery.Message.MessageID,
					loading_markup,
				))
				bot.Send(tgbotapi.NewChatAction(u.CallbackQuery.Message.Chat.ID, tgbotapi.ChatTyping))
				url, err := GetRandomCatURL()
				if err != nil {
					bot.Send(tgbotapi.NewMessage(
						u.Message.Chat.ID,
						fmt.Sprintf("Oops, an error occured: %s", err),
					))
				}
				edit := tgbotapi.NewEditMessageText(
					u.CallbackQuery.Message.Chat.ID,
					u.CallbackQuery.Message.MessageID,
					url,
				)
				edit.ReplyMarkup = &refresh_markup
				bot.Send(edit)
			},
		))
	for update := range updates {
		mux.Dispatch(&update)
	}
}
