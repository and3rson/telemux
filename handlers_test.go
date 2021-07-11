package telemux_test

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ExampleNewCommandHandler() {
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	u := tgbotapi.NewUpdate(0)
	updates, _ := bot.GetUpdatesChan(u)
	mux := tm.NewMux()
	mux.AddHandler(tm.NewCommandHandler(
		"add",
		nil,
		func(u *tm.Update) {
			args := u.Context["args"].([]string)
			if len(args) != 2 {
				bot.Send(tgbotapi.NewMessage(
					u.EffectiveChat().ID, "Wrong number of arguments. Example: /add 13 37"),
				)
				return
			}
			a, err1 := strconv.Atoi(args[0])
			b, err2 := strconv.Atoi(args[1])
			if err1 != nil || err2 != nil {
				bot.Send(tgbotapi.NewMessage(
					u.EffectiveChat().ID, "Arguments must be numbers. Example: /add 13 37"),
				)
				return
			}
			bot.Send(tgbotapi.NewMessage(
				u.EffectiveChat().ID, fmt.Sprintf("%d + %d = %d", a, b, a+b),
			))
		},
	))
	for update := range updates {
		mux.Dispatch(bot, update)
	}
}

func TestHandlerConsume(t *testing.T) {
	a, b, c := false, false, false
	h := tm.NewHandler(
		nil,
		func(u *tm.Update) { a = true },
		func(u *tm.Update) { b = true; u.Consume() },
		func(u *tm.Update) { c = true },
	)
	u := &tm.Update{tgbotapi.Update{}, nil, false, nil, nil}
	if !h.Process(u) {
		t.Error("Handler should return true")
	}
	if !a {
		t.Error("First handler should fire")
	}
	if !b {
		t.Error("Second handler should fire")
	}
	if c {
		t.Error("Third handler should not fire")
	}
}

func TestCommandHandler(t *testing.T) {
	h := tm.NewCommandHandler("test", nil, func(u *tm.Update) {})
	u := &tm.Update{}
	u.Update.Message = &tgbotapi.Message{}
	u.Update.Message.Text = "/test foo bar"
	u.Context = make(map[string]interface{})
	u.Bot = &tgbotapi.BotAPI{}
	u.Bot.Self.UserName = "testbot"
	if !h.Process(u) {
		t.Error("Handler should return true")
	}
	args := u.Context["args"].([]string)
	if len(args) != 2 {
		t.Error("There should be 2 args")
	}
	if args[0] != "foo" {
		t.Error("First arg should be 'foo'")
	}
	if args[1] != "bar" {
		t.Error("Second arg should be 'bar'")
	}

	h = tm.NewCommandHandler("foo bar", nil, func(u *tm.Update) {})
	u.Update.Message.Text = "/foo 42"
	if !h.Process(u) {
		t.Error("Handler should process update")
	}
	u.Update.Message.Text = "/bar 42"
	if !h.Process(u) {
		t.Error("Handler should process update")
	}
	u.Update.Message.Text = "/baz 42"
	if h.Process(u) {
		t.Error("Handler should not process update")
	}
}

func TestConversationHandler(t *testing.T) {
	NewUpdate := func(text string) *tm.Update {
		u := &tm.Update{}
		u.Update.Message = &tgbotapi.Message{}
		u.Update.Message.Text = text
		u.Update.Message.From = &tgbotapi.User{}
		u.Update.Message.From.ID = 13
		u.Update.Message.Chat = &tgbotapi.Chat{}
		u.Update.Message.Chat.ID = 37
		u.Context = make(map[string]interface{})
		u.Bot = &tgbotapi.BotAPI{}
		u.Bot.Self.UserName = "testbot"
		return u
	}
	askAgeEntered := false
	p := tm.NewLocalPersistence()
	h := tm.NewConversationHandler(
		"test",
		p,
		map[string][]*tm.Handler{
			"": {
				tm.NewCommandHandler("start", nil, func(u *tm.Update) {
					u.PersistenceContext.SetState("ask_name")
				}),
			},
			"ask_name": {
				tm.NewMessageHandler(tm.HasText(), func(u *tm.Update) {
					data := u.PersistenceContext.GetData()
					data["name"] = u.EffectiveMessage().Text
					u.PersistenceContext.SetData(data)
					u.PersistenceContext.SetState("ask_age")
				}),
			},
			"ask_age:enter": {
				tm.NewHandler(nil, func(u *tm.Update) {
					askAgeEntered = true
				}),
			},
			"ask_age": {
				tm.NewMessageHandler(tm.HasText(), func(u *tm.Update) {
					data := u.PersistenceContext.GetData()
					data["age"] = u.EffectiveMessage().Text
					u.PersistenceContext.SetData(data)
					u.PersistenceContext.SetState("ask_confirm")
				}),
			},
			"ask_confirm": {
				tm.NewCommandHandler("confirm", nil, func(u *tm.Update) {
					u.PersistenceContext.ClearData()
					u.PersistenceContext.SetState("")
				}),
			},
		},
		[]*tm.Handler{
			tm.NewCommandHandler("cancel", nil, func(u *tm.Update) { u.PersistenceContext.SetState(""); u.PersistenceContext.ClearData() }),
		},
	)
	pk := tm.PersistenceKey{"test", 13, 37}
	assert(!h.Process(NewUpdate("just some text")), t, "Random text must be ignored")
	assert(h.Process(NewUpdate("/start")), t, "/start must be processed")
	assert(p.GetState(pk) == "ask_name", t, "State must be ask_name, have", p.GetState(pk))
	assert(!askAgeEntered, t)
	assert(h.Process(NewUpdate("Foobar")), t, "Name must be processed")
	assert(p.GetState(pk) == "ask_age", t, "State must be ask_age, have", p.GetState(pk))
	assert(askAgeEntered, t)
	assert(reflect.DeepEqual(p.GetData(pk), map[string]interface{}{"name": "Foobar"}), t, "Unexpected persistence data")
	assert(h.Process(NewUpdate("18")), t, "Age must be processed")
	assert(p.GetState(pk) == "ask_confirm", t, "State must be ask_confirm, have", p.GetState(pk))
	assert(reflect.DeepEqual(p.GetData(pk), map[string]interface{}{"name": "Foobar", "age": "18"}), t, "Unexpected persistence data")
	assert(!h.Process(NewUpdate("foobar")), t, "Random text must be ignored")
	assert(p.GetState(pk) == "ask_confirm", t, "State must be ask_confirm, have", p.GetState(pk))
	assert(h.Process(NewUpdate("/confirm")), t, "/confirm must be processed")
	assert(p.GetState(pk) == "", t, "State must be empty, have", p.GetState(pk))
	assert(reflect.DeepEqual(p.GetData(pk), map[string]interface{}{}), t, "Persistence data must be empty")

	assert(h.Process(NewUpdate("/start")), t, "/start must be processed")
	assert(h.Process(NewUpdate("OtherUser")), t, "Name must be processed")
	assert(p.GetState(pk) == "ask_age", t, "State must be ask_age, have", p.GetState(pk))
	assert(reflect.DeepEqual(p.GetData(pk), map[string]interface{}{"name": "OtherUser"}), t, "Unexpected persistence data")
	assert(h.Process(NewUpdate("/cancel")), t, "/cancel must be processed")
	assert(p.GetState(pk) == "", t, "State must be empty, have", p.GetState(pk))
	assert(reflect.DeepEqual(p.GetData(pk), map[string]interface{}{}), t, "Persistence data must be empty")
}

func TestConvenienceHandlers(t *testing.T) {
	assert(strings.HasSuffix(
		getFunctionName(tm.NewInlineQueryHandler(".*", nil, func(u *tm.Update) {}).Filter),
		"And.func1",
	), t)
	assert(strings.HasSuffix(
		getFunctionName(tm.NewInlineQueryHandler(".*", tm.Any(), func(u *tm.Update) {}).Filter),
		"And.func1",
	), t)
	update := &tm.Update{
		Update: tgbotapi.Update{
			InlineQuery: &tgbotapi.InlineQuery{
				Query: "foo:bar:42",
			},
		},
		Context: map[string]interface{}{},
	}
	tm.NewInlineQueryHandler(`^foo:(\w+):(\d+)`, nil, func(u *tm.Update) {
	}).Process(update)
	assert(reflect.DeepEqual(update.Context["matches"], []string{"foo:bar:42", "bar", "42"}), t)

	assert(strings.HasSuffix(
		getFunctionName(tm.NewCallbackQueryHandler(".*", nil, func(u *tm.Update) {}).Filter),
		"And.func1",
	), t)
	assert(strings.HasSuffix(
		getFunctionName(tm.NewCallbackQueryHandler(".*", tm.Any(), func(u *tm.Update) {}).Filter),
		"And.func1",
	), t)
	update = &tm.Update{
		Update: tgbotapi.Update{
			CallbackQuery: &tgbotapi.CallbackQuery{
				Data: "foo:bar:42",
			},
		},
		Context: map[string]interface{}{},
	}
	tm.NewCallbackQueryHandler(`^foo:(\w+):(\d+)`, nil, func(u *tm.Update) {
	}).Process(update)
	assert(reflect.DeepEqual(update.Context["matches"], []string{"foo:bar:42", "bar", "42"}), t)

	assert(strings.HasSuffix(
		getFunctionName(tm.NewEditedMessageHandler(nil, func(u *tm.Update) {}).Filter),
		"IsEditedMessage.func1",
	), t)
	assert(strings.HasSuffix(
		getFunctionName(tm.NewEditedMessageHandler(tm.Any(), func(u *tm.Update) {}).Filter),
		"And.func1",
	), t)

	assert(strings.HasSuffix(
		getFunctionName(tm.NewChannelPostHandler(nil, func(u *tm.Update) {}).Filter),
		"IsChannelPost.func1",
	), t)
	assert(strings.HasSuffix(
		getFunctionName(tm.NewChannelPostHandler(tm.Any(), func(u *tm.Update) {}).Filter),
		"And.func1",
	), t)

	assert(strings.HasSuffix(
		getFunctionName(tm.NewEditedChannelPostHandler(nil, func(u *tm.Update) {}).Filter),
		"IsEditedChannelPost.func1",
	), t)
	assert(strings.HasSuffix(
		getFunctionName(tm.NewEditedChannelPostHandler(tm.Any(), func(u *tm.Update) {}).Filter),
		"And.func1",
	), t)
}
