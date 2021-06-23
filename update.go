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
	if u.Message != nil {
		return u.Message.Chat
	} else if u.EditedMessage != nil {
		return u.EditedMessage.Chat
	} else if u.ChannelPost != nil {
		return u.ChannelPost.Chat
	} else if u.EditedChannelPost != nil {
		return u.EditedChannelPost.Chat
	} else if u.CallbackQuery != nil && u.CallbackQuery.Message != nil {
		return u.CallbackQuery.Message.Chat
	}
	return nil
}

// GetEffectiveMessage retrieves message object from update.
func (u *Update) EffectiveMessage() *tgbotapi.Message {
	if u.Message != nil {
		return u.Message
	} else if u.EditedMessage != nil {
		return u.EditedMessage
	} else if u.ChannelPost != nil {
		return u.ChannelPost
	} else if u.EditedChannelPost != nil {
		return u.EditedChannelPost
	} else if u.CallbackQuery != nil {
		return u.CallbackQuery.Message
	}
	return nil
}
