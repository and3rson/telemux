// Package telemux is a flexible message router add-on for "go-telegram-bot-api".
//
// Make sure to check "go-telegram-bot-api" documentation first:
// https://github.com/go-telegram-bot-api/telegram-bot-api
package telemux

import (
	"fmt"
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Recoverer is a function that handles panics which happen during Dispatch.
type Recoverer = func(*Update, error, string)

type processor interface {
	process(u *Update) bool
}

// Mux contains handlers, nested multiplexers and global filter.
type Mux struct {
	processors   []processor // Contains instances of Mux & Handler
	Recoverer    Recoverer
	GlobalFilter FilterFunc
}

// NewMux creates new multiplexer.
func NewMux() *Mux {
	return &Mux{}
}

// AddHandler adds one or more handlers to multiplexer.
// This function returns the receiver for convenient chaining.
func (m *Mux) AddHandler(handlers ...*Handler) *Mux {
	for _, handler := range handlers {
		m.processors = append(m.processors, handler)
	}
	return m
}

// AddMux adds one or more nested multiplexers to this multiplexer.
// This function returns the receiver for convenient chaining.
func (m *Mux) AddMux(others ...*Mux) *Mux {
	for _, other := range others {
		m.processors = append(m.processors, other)
	}
	return m
}

// SetGlobalFilter sets a filter to be called for every update before any other filters.
// This function returns the receiver for convenient chaining.
func (m *Mux) SetGlobalFilter(filter FilterFunc) *Mux {
	m.GlobalFilter = filter
	return m
}

// SetRecoverer registers a function to call when a panic occurs.
// This function returns the receiver for convenient chaining.
func (m *Mux) SetRecoverer(recoverer Recoverer) *Mux {
	m.Recoverer = recoverer
	return m
}

func (m *Mux) tryRecover(u *Update) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", err)
		}
		if m.Recoverer != nil {
			m.Recoverer(u, err, string(debug.Stack()))
		} else {
			panic(err)
		}
	}
}

// Dispatch tells Mux to process the update.
// Returns true if the update was processed by one of the handlers.
func (m *Mux) Dispatch(bot *tgbotapi.BotAPI, u tgbotapi.Update) bool {
	return m.process(&Update{u, bot, false, nil})
}

func (m *Mux) process(u *Update) bool {
	defer m.tryRecover(u)

	if m.GlobalFilter != nil && !m.GlobalFilter(u) {
		return false
	}

	for _, processor := range m.processors {
		if processor.process(u) {
			return true
		}
	}
	return false
}
