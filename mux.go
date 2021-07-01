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
	GlobalFilter Filter
}

// NewMux creates new Mux.
func NewMux() *Mux {
	return &Mux{}
}

// AddHandler adds handler to Mux.
func (m *Mux) AddHandler(h ...*Handler) *Mux {
	m.Targets = append(m.Targets, h)
	return m
}

// AddMux adds nested Mux to this Mux.
func (m *Mux) AddMux(other ...*Mux) *Mux {
	m.Targets = append(m.Targets, other)
	return m
}

// SetGlobalFilter sets a filter to be called for every update before any other filters.
func (m *Mux) SetGlobalFilter(filter Filter) *Mux {
	m.GlobalFilter = filter
	return m
}

// SetRecoverer registers a function to call when a panic occurs.
func (m *Mux) SetRecoverer(recoverer Recoverer) *Mux {
	m.Recoverer = recoverer
	return m
}

// Dispatch tells Mux to process the update.
func (m *Mux) Dispatch(bot *tgbotapi.BotAPI, u tgbotapi.Update) bool {
	uu := Update{u, bot, false, nil}

	defer func() {
		if err, ok := recover().(error); ok {
			if m.Recoverer != nil {
				m.Recoverer(&uu, error(err), string(debug.Stack()))
			} else {
				panic(err)
			}
		}
		// TODO: what if err is string?
	}()

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
			accepted := target.Filter(&uu)
			if uu.Consumed {
				return true
			} else if accepted {
				target.Handle(&uu)
				return true
			}
		default:
			panic(&ConfigurationError{fmt.Sprintf("%T is not an instance of telemux.Handler or telemux.Mux: %V", target, target)})
		}
	}
	return false
}
