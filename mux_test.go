package telemux_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ExampleNewMux() {
	// This part is a boilerplate from go-telegram-bot-api library.
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	// Create a multiplexer with two handlers: one for command and one for all messages.
	// If a handler cannot handle the update (fails the filter),
	// multiplexer will proceed to the next handler.
	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsCommandMessage("start"),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Hello! Say something. :)"))
			},
		)).
		AddHandler(tm.NewHandler(
			tm.Any(),
			func(u *tm.Update) {
				bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "You said: "+u.Message.Text))
			},
		))
	// Dispatch all telegram updates to multiplexer
	for update := range updates {
		mux.Dispatch(bot, update)
	}
}

func TestMuxDispatch(t *testing.T) {
	NewTGUpdate := func(text string) tgbotapi.Update {
		u := tgbotapi.Update{}
		u.Message = &tgbotapi.Message{}
		u.Message.Text = text
		return u
	}

	stack := []string{}
	mux := tm.NewMux().
		AddHandler(tm.NewMessageHandler(tm.HasRegex("^1"), func(u *tm.Update) { stack = append(stack, "1") })).
		AddMux(
			tm.NewMux().
				SetGlobalFilter(tm.HasRegex("^2")).
				AddHandler(tm.NewMessageHandler(tm.HasRegex("^21"), func(u *tm.Update) { stack = append(stack, "21") })).
				AddHandler(tm.NewMessageHandler(tm.HasRegex("^22"), func(u *tm.Update) { stack = append(stack, "22") })),
		)

	assert(mux.Dispatch(nil, NewTGUpdate("1")), t, "Dispatch 1")
	assert(reflect.DeepEqual(stack, []string{"1"}), t, "Check 1")

	stack = []string{}
	assert(!mux.Dispatch(nil, NewTGUpdate("2")), t, "Dispatch 2")
	assert(mux.Dispatch(nil, NewTGUpdate("21")), t, "Dispatch 21")
	assert(reflect.DeepEqual(stack, []string{"21"}), t, "Check 21")

	stack = []string{}
	assert(mux.Dispatch(nil, NewTGUpdate("22")), t, "Dispatch 22")
	assert(reflect.DeepEqual(stack, []string{"22"}), t, "Check 22")

	stack = []string{}
	assert(!mux.Dispatch(nil, NewTGUpdate("23")), t, "Dispatch 23")
	assert(reflect.DeepEqual(stack, []string{}), t, "Check 22")
	assert(!mux.Dispatch(nil, NewTGUpdate("33")), t, "Dispatch 33")
	assert(reflect.DeepEqual(stack, []string{}), t, "Check 33")
}

func TestMuxRecover(t *testing.T) {
	NewTGUpdate := func(text string) tgbotapi.Update {
		u := tgbotapi.Update{}
		u.Message = &tgbotapi.Message{}
		u.Message.Text = text
		return u
	}

	recovered := struct {
		e error
		s string
	}{}

	mux := tm.NewMux().
		AddHandler(tm.NewMessageHandler(nil, func(u *tm.Update) {
			if u.EffectiveMessage().Text == "panic_string" {
				panic("boom")
			} else if u.EffectiveMessage().Text == "panic_error" {
				panic(errors.New("boom"))
			}
		})).
		SetRecover(func(u *tm.Update, e error, s string) {
			recovered.e = e
			recovered.s = s
		})

	mux.Dispatch(nil, NewTGUpdate("keep_calm"))
	assert(recovered.e == nil, t)
	assert(recovered.s == "", t)
	mux.Dispatch(nil, NewTGUpdate("panic_string"))
	assert(recovered.e.Error() == "boom", t)
	assert(recovered.s != "", t)
	mux.Dispatch(nil, NewTGUpdate("panic_error"))
	assert(recovered.e.Error() == "boom", t)
	assert(recovered.s != "", t)

	mux.SetRecover(nil)
	func() {
		defer func() {
			r := recover()
			if r == nil || r.(error).Error() != "boom" {
				t.Error("Expected unhandled panic")
			}
		}()
		mux.Dispatch(nil, NewTGUpdate("panic_error"))
	}()
}
