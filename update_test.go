package telemux_test

import (
	"testing"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestEffectiveUser(t *testing.T) {
	u := tm.Update{}
	assert(u.EffectiveUser() == nil, t)

	u.Update.PreCheckoutQuery = &tgbotapi.PreCheckoutQuery{}
	u.Update.PreCheckoutQuery.From = &tgbotapi.User{ID: 1}
	assert(u.EffectiveUser().ID == 1, t)

	u.Update.ShippingQuery = &tgbotapi.ShippingQuery{}
	u.Update.ShippingQuery.From = &tgbotapi.User{ID: 2}
	assert(u.EffectiveUser().ID == 2, t)

	u.Update.CallbackQuery = &tgbotapi.CallbackQuery{}
	u.Update.CallbackQuery.From = &tgbotapi.User{ID: 3}
	assert(u.EffectiveUser().ID == 3, t)

	u.Update.ChosenInlineResult = &tgbotapi.ChosenInlineResult{}
	u.Update.ChosenInlineResult.From = &tgbotapi.User{ID: 4}
	assert(u.EffectiveUser().ID == 4, t)

	u.Update.InlineQuery = &tgbotapi.InlineQuery{}
	u.Update.InlineQuery.From = &tgbotapi.User{ID: 5}
	assert(u.EffectiveUser().ID == 5, t)

	u.Update.EditedChannelPost = &tgbotapi.Message{}
	u.Update.EditedChannelPost.From = &tgbotapi.User{ID: 6}
	assert(u.EffectiveUser().ID == 6, t)

	u.Update.ChannelPost = &tgbotapi.Message{}
	u.Update.ChannelPost.From = &tgbotapi.User{ID: 7}
	assert(u.EffectiveUser().ID == 7, t)

	u.Update.EditedMessage = &tgbotapi.Message{}
	u.Update.EditedMessage.From = &tgbotapi.User{ID: 8}
	assert(u.EffectiveUser().ID == 8, t)

	u.Update.Message = &tgbotapi.Message{}
	u.Update.Message.From = &tgbotapi.User{ID: 9}
	assert(u.EffectiveUser().ID == 9, t)
}

func TestEffectiveChat(t *testing.T) {
	u := tm.Update{}
	assert(u.EffectiveChat() == nil, t)

	u.Update.Message = &tgbotapi.Message{}
	u.Update.Message.Chat = &tgbotapi.Chat{ID: 42}
	assert(u.EffectiveChat().ID == 42, t)
}

func TestEffectiveMessage(t *testing.T) {
	u := tm.Update{}
	assert(u.EffectiveMessage() == nil, t)

	u.Update.CallbackQuery = &tgbotapi.CallbackQuery{}
	u.Update.CallbackQuery.Message = &tgbotapi.Message{Text: "Foo"}
	assert(u.EffectiveMessage().Text == "Foo", t)
}
