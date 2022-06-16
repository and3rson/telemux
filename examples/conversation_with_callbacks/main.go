package main

import (
	"fmt"
	"log"
	"os"

	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	SearchProductMessage = "Search product by ID"
)

type Product struct {
	Sku   string
	Title string
	Image string
	Price int
}

var storage = map[string]Product{
	"p01": {
		Sku:   "a87xn",
		Title: "My awesome product",
		Image: "https://picsum.photos/id/252/600/400",
		Price: 89,
	},
}

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(SearchProductMessage),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {

				msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
				startKeyboard.OneTimeKeyboard = true
				msg.ReplyMarkup = startKeyboard

				bot.Send(msg)
			},
		)).
		AddHandler(tm.NewConversationHandler(
			"get_product_data",
			tm.NewLocalPersistence(),
			tm.StateMap{
				"": {
					tm.NewHandler(tm.And(tm.IsMessage(), tm.HasRegex("^"+SearchProductMessage)), func(u *tm.Update) {
						msg := tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"Provide a product ID",
						)

						bot.Send(msg)
						u.PersistenceContext.SetState("enter_productId")
					}),
				},
				"enter_productId": {
					tm.NewHandler(tm.HasText(), func(u *tm.Update) {
						bot.Send(tgbotapi.NewChatAction(u.Message.Chat.ID, tgbotapi.ChatTyping))

						productId := u.Message.Text

						msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")

						product, ok := storage[productId]
						if !ok {
							msg.Text = "Product not found"
							bot.Send(msg)
							return
						}

						file := tgbotapi.FileURL(product.Image)

						share := tgbotapi.NewPhoto(u.Message.Chat.ID, file)
						share.Caption = fmt.Sprintf("%s\n\nPrice: %d $\nSKU: %s", product.Title, product.Price, product.Sku)

						share.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Add to cart", "buyProduct "+productId),
							),
						)

						bot.Send(share)
					}),
					tm.NewHandler(tm.And(tm.Not(tm.IsCommandMessage("back")), tm.Not(tm.IsCallbackQuery())), func(u *tm.Update) {
						bot.Send(tgbotapi.NewMessage(
							u.Message.Chat.ID,
							"ID only!",
						))
					}),
				},
			},
			[]*tm.Handler{
				tm.NewHandler(tm.IsCommandMessage("back"), func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")

					msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)

					startKeyboard.OneTimeKeyboard = true
					msg.ReplyMarkup = startKeyboard

					bot.Send(msg)
				}),
				// During the active conversation this callback handler will be invoked first
				tm.NewHandler(tm.IsCallbackQuery(), func(u *tm.Update) {

					loadingMarkup := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Loading...", "load"),
						),
					)

					refreshMarkup := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("In cart", "return"),
						),
					)

					bot.Send(tgbotapi.NewCallback(u.CallbackQuery.ID, "Refreshing..."))
					bot.Send(tgbotapi.NewEditMessageReplyMarkup(
						u.CallbackQuery.Message.Chat.ID,
						u.CallbackQuery.Message.MessageID,
						loadingMarkup,
					))
					bot.Send(tgbotapi.NewChatAction(u.CallbackQuery.Message.Chat.ID, tgbotapi.ChatTyping))

					bot.Send(tgbotapi.NewEditMessageReplyMarkup(
						u.CallbackQuery.Message.Chat.ID,
						u.CallbackQuery.Message.MessageID,
						refreshMarkup,
					))
				}),
			},
		)).
		AddHandler(tm.NewHandler(tm.IsCallbackQuery(), func(u *tm.Update) {
			callback := tgbotapi.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Print(err)
			}
			// Next:
			// Here you can handle any query callbacks even it from closed conversations.
		}))

	for update := range updates {
		mux.Dispatch(bot, update)
	}
}
