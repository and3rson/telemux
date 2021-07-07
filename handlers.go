package telemux

import (
	"regexp"
	"strings"
)

// HandleFunc processes update.
type HandleFunc func(u *Update)

// Handler defines a function that will handle updates that pass the filtering.
type Handler struct {
	Filter  FilterFunc
	Handles []HandleFunc
}

// Process runs handler with provided Update.
func (h *Handler) Process(u *Update) bool {
	if h.Filter(u) {
		for i := 0; i < len(h.Handles) && !u.Consumed; i++ {
			h.Handles[i](u)
		}
		return true
	}
	return false
}

// NewHandler creates a new generic handler.
func NewHandler(filter FilterFunc, handles ...HandleFunc) *Handler {
	if filter == nil {
		filter = Any()
	}
	return &Handler{filter, handles}
}

// NewMessageHandler creates a handler for updates that contain message.
func NewMessageHandler(filter FilterFunc, handles ...HandleFunc) *Handler {
	newFilter := IsMessage()
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	return NewHandler(newFilter, handles...)
}

// NewCommandHandler is an extension for NewMessageHandler that creates a handler for updates that contain message with command.
// It also populates u.Context["args"] with a slice of strings.
//
// For example, when invoked as `/somecmd foo bar 1337`, u.Context["args"] will be set to []string{"foo", "bar", "1337"}
func NewCommandHandler(command string, filter FilterFunc, handles ...HandleFunc) *Handler {
	handles = append([]HandleFunc{
		func(u *Update) {
			u.Context["args"] = strings.Split(u.Message.Text, " ")[1:]
		},
	}, handles...)
	newFilter := IsCommandMessage(command)
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	return NewMessageHandler(
		newFilter,
		handles...,
	)
}

// NewInlineQueryHandler creates a handler for updates that contain inline query which matches the pattern as regexp.
func NewInlineQueryHandler(pattern string, filter FilterFunc, handles ...HandleFunc) *Handler {
	exp := regexp.MustCompile(pattern)
	newFilter := And(IsInlineQuery(), func(u *Update) bool {
		return exp.Match([]byte(u.InlineQuery.Query))
	})
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	handles = append([]HandleFunc{
		func(u *Update) {
			u.Context["exp"] = exp
			u.Context["matches"] = exp.FindStringSubmatch(u.InlineQuery.Query)
		},
	}, handles...)
	return NewHandler(newFilter, handles...)
}

// NewCallbackQueryHandler creates a handler for updates that contain callback query which matches the pattern as regexp.
func NewCallbackQueryHandler(pattern string, filter FilterFunc, handles ...HandleFunc) *Handler {
	exp := regexp.MustCompile(pattern)
	newFilter := And(IsCallbackQuery(), func(u *Update) bool {
		return exp.Match([]byte(u.CallbackQuery.Data))
	})
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	handles = append([]HandleFunc{
		func(u *Update) {
			u.Context["exp"] = exp
			u.Context["matches"] = exp.FindStringSubmatch(u.CallbackQuery.Data)
		},
	}, handles...)
	return NewHandler(newFilter, handles...)
}

// NewEditedMessageHandler creates a handler for updates that contain edited message.
func NewEditedMessageHandler(filter FilterFunc, handles ...HandleFunc) *Handler {
	newFilter := IsEditedMessage()
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	return NewHandler(newFilter, handles...)
}

// NewChannelPostHandler creates a handler for updates that contain channel post.
func NewChannelPostHandler(filter FilterFunc, handles ...HandleFunc) *Handler {
	newFilter := IsChannelPost()
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	return NewHandler(newFilter, handles...)
}

// NewEditedChannelPostHandler creates a handler for updates that contain edited channel post.
func NewEditedChannelPostHandler(filter FilterFunc, handles ...HandleFunc) *Handler {
	newFilter := IsEditedChannelPost()
	if filter != nil {
		newFilter = And(newFilter, filter)
	}
	return NewHandler(newFilter, handles...)
}

// NewConversationHandler creates a conversation handler.
//
// "conversationID" distinguishes this conversation from the others. The main goal of this identifier is to allow persistence to keep track of different conversation states independently without mixing them together.
//
// "persistence" defines where to store conversation state & intermediate inputs from the user. Without persistence, a conversation would not be able to "remember" what "step" the user is at.
//
// "states" define what handlers to use in which state. States are usually strings like "upload_photo", "send_confirmation", "wait_for_text" and describe the "step" the user is currently at.
// Empty string (`""`) should be used as an initial/final state (i. e. if the conversation has not started yet or has already finished.)
// For each state you must provide a slice with at least one Handler. If none of the handlers can handle the update, the default handlers are attempted (see below).
// In order to switch to a different state your Handler must call `u.PersistenceContext.SetState("STATE_NAME") ` replacing STATE_NAME with the name of the state you want to switch into.
// Conversation data can be accessed with `u.PersistenceContext.GetData()` and updated with `u.PersistenceContext.SetData(newData)`.
//
// "defaults" are "appended" to every state except default state (`""`). They are useful to handle commands such as "/cancel" or to display some default message.
func NewConversationHandler(
	conversationID string,
	persistence ConversationPersistence,
	states map[string][]*Handler,
	defaults []*Handler,
) *Handler {
	var handler *Handler
	handler = &Handler{
		func(u *Update) bool {
			user, chat := u.EffectiveUser(), u.EffectiveChat()
			pk := PersistenceKey{conversationID, user.ID, chat.ID}
			state := persistence.GetState(pk)
			candidates := states[state]
			if state != "" {
				candidates = append(candidates, defaults...)
			}
			u.PersistenceContext = &PersistenceContext{
				Persistence: persistence,
				PK:          pk,
			}
			defer func() { u.PersistenceContext = nil }()
			for _, handler := range candidates {
				if handler.Filter(u) {
					return true
				}
			}
			return false
		},
		[]HandleFunc{func(u *Update) {
			user, chat := u.EffectiveUser(), u.EffectiveChat()
			pk := PersistenceKey{conversationID, user.ID, chat.ID}
			state := persistence.GetState(pk)
			candidates := states[state]
			if state != "" {
				candidates = append(candidates, defaults...)
			}
			if u.PersistenceContext == nil {
				u.PersistenceContext = &PersistenceContext{
					Persistence: persistence,
					PK:          pk,
				}
				defer func() { u.PersistenceContext = nil }()
			}
			defer func() {
				if u.PersistenceContext.NewState != nil {
					// TODO: Add docs for :enter hook
					if handlers, ok := states[*u.PersistenceContext.NewState+":enter"]; ok {
						for _, handler := range handlers {
							handler.Process(u)
						}
					}
				}
			}()
			for _, handler := range candidates {
				if handler.Filter(u) {
					for i := 0; i < len(handler.Handles) && !u.Consumed; i++ {
						handler.Handles[i](u)
					}
					return
				}
			}
		}},
	}
	return handler
}
