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

// Recoverer is a function that handles panics.
type Recoverer = func(*Update, error, string)

// ConfigurationError describes internal configuration error.
type ConfigurationError struct {
	s string
}

func (e *ConfigurationError) Error() string {
	return e.s
}

// Mux is a container for handlers.
type Mux struct {
	Targets      []interface{} // Target can be either Handler or Mux
	Recoverer    Recoverer
	GlobalFilter FilterFunc
}

// NewMux creates new Mux.
func NewMux() *Mux {
	return &Mux{}
}

// AddHandler adds handler to Mux.
func (m *Mux) AddHandler(handlers ...*Handler) *Mux {
	for _, handler := range handlers {
		m.Targets = append(m.Targets, handler)
	}
	return m
}

// AddMux adds nested Mux to this Mux.
func (m *Mux) AddMux(others ...*Mux) *Mux {
	for _, other := range others {
		m.Targets = append(m.Targets, other)
	}
	return m
}

// SetGlobalFilter sets a filter to be called for every update before any other filters.
func (m *Mux) SetGlobalFilter(filter FilterFunc) *Mux {
	m.GlobalFilter = filter
	return m
}

// SetRecoverer registers a function to call when a panic occurs.
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
func (m *Mux) Dispatch(bot *tgbotapi.BotAPI, u tgbotapi.Update) bool {
	uu := Update{u, bot, false, nil}

	defer m.tryRecover(&uu)

	if m.GlobalFilter != nil && !m.GlobalFilter(&uu) {
		return false
	}

	for _, target := range m.Targets {
		switch target := target.(type) {
		case *Mux:
			if target.Dispatch(bot, u) {
				return true
			}
		case *Handler:
			if target.Filter(&uu) {
				for i := 0; i < len(target.Handles) && !uu.Consumed; i++ {
					target.Handles[i](&uu)
				}
				return true
			}
		default:
			panic(&ConfigurationError{fmt.Sprintf("%T is not an instance of telemux.Handler or telemux.Mux: %v", target, target)})
		}
	}
	return false
}
