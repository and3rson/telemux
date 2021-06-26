package telemux

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Update wraps tgbotapi.Update.
// It provides some convenient functions such as GetEffectiveUser.
type Update struct {
	*tgbotapi.Update
}

// GetEffectiveUser retrieves user object from update.
func (u *Update) EffectiveUser() *tgbotapi.User {
	if u.Message != nil {
		return u.Message.From
	} else if u.EditedMessage != nil {
		return u.EditedMessage.From
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
	return nil
}

// GetEffectiveChat retrieves chat object from update.
func (u *Update) EffectiveChat() *tgbotapi.Chat {
	message := u.EffectiveMessage()
	if message != nil {
		return message.Chat
	}
	return nil
}

// GetEffectiveMessage retrieves message object from update.
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
