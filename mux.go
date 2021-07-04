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

// RecoverFunc handles panics which happen during Dispatch.
type RecoverFunc = func(*Update, error, string)

// Processor is either Handler or Mux.
type Processor interface {
	Process(u *Update) bool
}

// Mux contains handlers, nested multiplexers and global filter.
type Mux struct {
	Processors   []Processor // Contains instances of Mux & Handler
	Recover      RecoverFunc
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
		m.Processors = append(m.Processors, handler)
	}
	return m
}

// AddMux adds one or more nested multiplexers to this multiplexer.
// This function returns the receiver for convenient chaining.
func (m *Mux) AddMux(others ...*Mux) *Mux {
	for _, other := range others {
		m.Processors = append(m.Processors, other)
	}
	return m
}

// SetGlobalFilter sets a filter to be called for every update before any other filters.
// This function returns the receiver for convenient chaining.
func (m *Mux) SetGlobalFilter(filter FilterFunc) *Mux {
	m.GlobalFilter = filter
	return m
}

// SetRecover registers a function to call when a panic occurs.
// This function returns the receiver for convenient chaining.
func (m *Mux) SetRecover(recover RecoverFunc) *Mux {
	m.Recover = recover
	return m
}

func (m *Mux) tryRecover(u *Update) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		if m.Recover != nil {
			m.Recover(u, err, string(debug.Stack()))
		} else {
			panic(err)
		}
	}
}

// Dispatch tells Mux to process the update.
// Returns true if the update was processed by one of the handlers.
func (m *Mux) Dispatch(bot *tgbotapi.BotAPI, u tgbotapi.Update) bool {
	return m.Process(&Update{u, bot, false, nil, make(map[string]interface{})})
}

// Process runs mux with provided update.
func (m *Mux) Process(u *Update) bool {
	defer m.tryRecover(u)

	if m.GlobalFilter != nil && !m.GlobalFilter(u) {
		return false
	}

	for _, Processor := range m.Processors {
		if Processor.Process(u) {
			return true
		}
	}
	return false
}
