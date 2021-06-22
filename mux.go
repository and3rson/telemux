// telemux is a flexible message router add-on for "go-telegram-bot-api".
//
// Make sure to check "go-telegram-bot-api" documentation first:
// https://github.com/go-telegram-bot-api/telegram-bot-api
package telemux

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Mux is a container for handlers.
type Mux struct {
	Handlers []*Handler
}

// Update is an alias for tgbotapi.Update.
type Update = tgbotapi.Update

// NewMux creates new Mux.
func NewMux() *Mux {
	return &Mux{}
}

// AddHandler adds handler to Mux.
func (m *Mux) AddHandler(h *Handler) *Mux {
	m.Handlers = append(m.Handlers, h)
	return m
}

// Dispatch tells Mux to process the update.
func (m *Mux) Dispatch(u *tgbotapi.Update) bool {
	for _, handler := range m.Handlers {
		if handler.Filter(u) {
			handler.Handle(u)
			return true
		}
	}
	return false
}
