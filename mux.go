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
	uu := Update{u}
	for _, handler := range m.Handlers {
		if handler.Filter(&uu) {
			handler.Handle(&uu)
			return true
		}
	}
	return false
}
