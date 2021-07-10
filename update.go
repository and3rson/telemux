package telemux

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Update wraps tgbotapi.Update and stores some additional data.
type Update struct {
	tgbotapi.Update
	Bot                *tgbotapi.BotAPI
	Consumed           bool
	PersistenceContext *PersistenceContext
	Context            Map
}

// Consume marks update as processed. Used by handler functions to interrupt further processing of the update.
func (u *Update) Consume() {
	u.Consumed = true
}

// EffectiveUser retrieves user object from update.
func (u *Update) EffectiveUser() *tgbotapi.User {
	if u.Message != nil {
		return u.Message.From
	} else if u.EditedMessage != nil {
		return u.EditedMessage.From
	} else if u.ChannelPost != nil {
		return u.ChannelPost.From
	} else if u.EditedChannelPost != nil {
		return u.EditedChannelPost.From
	} else if u.InlineQuery != nil {
		return u.InlineQuery.From
	} else if u.ChosenInlineResult != nil {
		return u.ChosenInlineResult.From
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.From
	} else if u.ShippingQuery != nil {
		return u.ShippingQuery.From
	} else if u.PreCheckoutQuery != nil {
		return u.PreCheckoutQuery.From
	} // TODO: Polls not yet supported by go-telegram-bot-api?
	log.Println("Sender not found in update object! This is possibly a bug.")
	return nil
}

// EffectiveChat retrieves chat object from update.
func (u *Update) EffectiveChat() *tgbotapi.Chat {
	message := u.EffectiveMessage()
	if message != nil {
		return message.Chat
	}
	return nil
}

// EffectiveMessage retrieves message object from update.
func (u *Update) EffectiveMessage() *tgbotapi.Message {
	candidates := []*tgbotapi.Message{u.Message, u.EditedMessage, u.ChannelPost, u.EditedChannelPost}
	for _, message := range candidates {
		if message != nil {
			return message
		}
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.Message
	}
	return nil
}

// Fields returns some metadata of this update. Useful for passing this directly into logrus.WithFields() or other loggers.
func (u *Update) Fields() Map {
	chatID := ""
	chatname := ""
	userID := ""
	username := ""
	chat := u.EffectiveChat()
	user := u.EffectiveUser()
	if chat != nil {
		chatID = fmt.Sprint(chat.ID)
		chatname = chat.Title
	}
	if user != nil {
		userID = fmt.Sprint(user.ID)
		username = user.String()
	}
	return Map{"chatID": chatID, "chatname": chatname, "userID": userID, "username": username}
}
