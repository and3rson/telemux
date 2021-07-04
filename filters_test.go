package telemux_test

import (
	"testing"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestIsCommandMessage(t *testing.T) {
	Check := func(text string, isCommand bool) {
		u := &tm.Update{}
		u.Update.Message = &tgbotapi.Message{}
		u.Update.Message.Text = text
		u.Bot = &tgbotapi.BotAPI{}
		u.Bot.Self.UserName = "testbot"
		actual := tm.IsCommandMessage("foo")(u)
		if actual != isCommand {
			t.Errorf("Testing %s: IsCommandMessage = %v, expected = %v", text, actual, isCommand)
		}
	}

	Check("asd", false)
	Check("/", false)
	Check("/foo", true)
	Check("/foo@", true)
	Check("/foo ", true)
	Check("/foo bar", true)
	Check("/foo@testbot", true)
	Check("/foo@testbot bar", true)
	Check("/foo@ bar", true)
	Check("/foo@nope", false)
	Check("/foo@nope bar", false)
	Check("/bar", false)
	Check("/bar baz", false)
	Check("/bar@testbot baz", false)
}

func TestIsAnyCommandMessage(t *testing.T) {
	Check := func(text string, isCommand bool) {
		u := &tm.Update{}
		u.Update.Message = &tgbotapi.Message{}
		u.Update.Message.Text = text
		u.Bot = &tgbotapi.BotAPI{}
		u.Bot.Self.UserName = "testbot"
		actual := tm.IsAnyCommandMessage()(u)
		if actual != isCommand {
			t.Errorf("Testing %s: IsAnyCommandMessage = %v, expected = %v", text, actual, isCommand)
		}
	}

	Check("asd", false)
	Check("/", false)
	Check("/foo", true)
	Check("/foo@", true)
	Check("/foo ", true)
	Check("/foo bar", true)
	Check("/foo@testbot", true)
	Check("/foo@testbot bar", true)
	Check("/foo@ bar", true)
	Check("/foo@nope", false)
	Check("/foo@nope bar", false)
}
