package telemux

// Handler defines a function that will handle updates that pass the filtering.
type Handler struct {
	Filter Filter
	Handle func(u *Update)
}

// NewHandler creates a new handler.
func NewHandler(filter Filter, handle func(u *Update)) *Handler {
	if filter == nil {
		filter = Any()
	}
	return &Handler{filter, handle}
}

// NewAsyncHandler creates a new handler which will be executed in a goroutine.
func NewAsyncHandler(filter Filter, handle func(u *Update)) *Handler {
	return NewHandler(filter, func(u *Update) {
		go handle(u)
	})
}

// NewConversationHandler creates a conversation handler.
//
// "conversationID" distinguishes this conversation from the others. The main goal of this identifier is to allow persistence to keep track of different conversation states independently without mixing them together.
//
// "persistence" defines where to store conversation state & intermediate inputs from the user. Without persistence, a conversation would not be able to "remember" what "step" the user is at.
//
// "states" define what handlers to use in which state. States are usually strings like "upload_photo", "send_confirmation", "wait_for_text" and describe the "step" the user is currently at.
// It's recommended to have an empty string (`""`) as an initial state (i. e. if the conversation has not started yet or has already finished.)
// For each state you can provide a list of at least one Handler. If none of the handlers can handle the update, the default handlers are attempted (see below).
// In order to switch to a different state your Handler must call `u.PersistenceContext.SetState("STATE_NAME") ` replacing STATE_NAME with the name of the state you want to switch into.
// Conversation data can be accessed with `u.PersistenceContext.GetData()` and updated with `u.PersistenceContext.SetData(newData)`.
//
// "defaults" are "appended" to every state. They are useful to handle commands such as "/cancel" or to display some default message.
func NewConversationHandler(
	conversationID string,
	persistence ConversationPersistence,
	states map[string][]*Handler,
	defaults []*Handler,
) *Handler {
	return &Handler{
		func(u *Update) bool {
			user, chat := u.EffectiveUser(), u.EffectiveChat()
			pk := PersistenceKey{conversationID, user.ID, chat.ID}
			candidates := states[persistence.GetState(pk)]
			if len(defaults) > 0 {
				candidates = append(candidates, defaults...)
			}
			u.PersistenceContext = &PersistenceContext{
				Persistence: persistence,
				PK:          pk,
			}
			defer func() { u.PersistenceContext = nil }()
			for _, handler := range candidates {
				accepted := handler.Filter(u)
				if accepted {
					return true
				} else if u.Consumed {
					return true
				}
			}
			return false
		},
		func(u *Update) {
			user, chat := u.EffectiveUser(), u.EffectiveChat()
			pk := PersistenceKey{conversationID, user.ID, chat.ID}
			candidates := states[persistence.GetState(pk)]
			if len(defaults) > 0 {
				candidates = append(candidates, defaults...)
			}
			u.PersistenceContext = &PersistenceContext{
				Persistence: persistence,
				PK:          pk,
			}
			defer func() { u.PersistenceContext = nil }()
			for _, handler := range candidates {
				if handler.Filter(u) {
					handler.Handle(u)
					return
				}
			}
		},
	}
}
