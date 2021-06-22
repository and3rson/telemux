package telemux

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Handler defines a function that will handle updates that pass the filtering.
type Handler struct {
	Filter Filter
	Handle func(u *Update)
}

// TransitionHandler is similar to Handler but requires Handle function to return a string.
// This string is used to change conversation state in conversation handlers.
type TransitionHandler struct {
	Filter Filter
	Handle func(u *tgbotapi.Update, data Data) string
}

// NewHandler creates a new handler.
func NewHandler(filter Filter, handle func(u *Update)) *Handler {
	return &Handler{filter, handle}
}

// NewTransitionHandler creates a new transition handler for a conversation handler.
func NewTransitionHandler(filter Filter, handle func(u *Update, data Data) string) *TransitionHandler {
	return &TransitionHandler{filter, handle}
}

// NewConversationHandler creates a conversation handler.
//
// "conversationID" distinguishes this conversation from the others. The main goal of this identifier is to allow persistence to keep track of different conversation states independently without mixing them together.
//
// "persistence" defines where to store conversation state & intermediate inputs from the user. Without persistence, a conversation would not be able to "remember" what "step" the user is at.
//
// "states" define the TransitionHandlers to use in what state. States are usually strings like "upload_photo", "send_confirmation", "wait_for_text" and describe the "step" the user is currently at. It's recommended to have an empty string (`""`) as an initial state (i. e. if the conversation has not started yet or has already finished.) For each state you can provide a list of at least one TransitionHandler. If none of the handlers can handle the update, the default handlers are attempted (see below).
// In order to switch to a different state your TransitionHandler must return a string that contains the name of the state you want to switch into.
//
// "defaults" are "appended" to every state. They are useful to handle commands such as "/cancel" or to display some default message.
func NewConversationHandler(
	conversationID string,
	persistence ConversationPersistence,
	states map[string][]*TransitionHandler,
	defaults []*TransitionHandler,
) *Handler {
	return &Handler{
		func(u *Update) bool {
			user, chat := GetEffectiveUser(u), GetEffectiveChat(u)
			pk := PersistenceKey{user.ID, chat.ID}
			candidates := states[persistence.GetState(conversationID, pk)]
			if len(defaults) > 0 {
				candidates = append(candidates, defaults...)
			}
			for _, handler := range candidates {
				if handler.Filter(u) {
					return true
				}
			}
			return false
		},
		func(u *tgbotapi.Update) {
			user, chat := GetEffectiveUser(u), GetEffectiveChat(u)
			pk := PersistenceKey{user.ID, chat.ID}
			candidates := states[persistence.GetState(conversationID, pk)]
			if len(defaults) > 0 {
				candidates = append(candidates, defaults...)
			}
			for _, handler := range candidates {
				if handler.Filter(u) {
					data := persistence.GetConversationData(conversationID, pk)
					persistence.SetState(conversationID, pk, handler.Handle(u, data))
					persistence.SetConversationData(conversationID, pk, data)
					return
				}
			}
		},
	}
}
