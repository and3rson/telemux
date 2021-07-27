package telemux_test

import (
	"fmt"
	"testing"

	tm "github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	Check("/foox", false)
	Check("/fo", false)
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

func TestUpdateTypeFilters(t *testing.T) {
	u := &tm.Update{}
	assert(!tm.IsInlineQuery()(u), t)
	u.InlineQuery = &tgbotapi.InlineQuery{}
	assert(tm.IsInlineQuery()(u), t)

	u = &tm.Update{}
	assert(!tm.IsCallbackQuery()(u), t)
	u.CallbackQuery = &tgbotapi.CallbackQuery{}
	assert(tm.IsCallbackQuery()(u), t)

	u = &tm.Update{}
	assert(!tm.IsEditedMessage()(u), t)
	u.EditedMessage = &tgbotapi.Message{}
	assert(tm.IsEditedMessage()(u), t)

	u = &tm.Update{}
	assert(!tm.IsChannelPost()(u), t)
	u.ChannelPost = &tgbotapi.Message{}
	assert(tm.IsChannelPost()(u), t)

	u = &tm.Update{}
	assert(!tm.IsEditedChannelPost()(u), t)
	u.EditedChannelPost = &tgbotapi.Message{}
	assert(tm.IsEditedChannelPost()(u), t)
}

func TestContentFilters(t *testing.T) {
	u := &tm.Update{}
	assert(!tm.HasText()(u), t)
	u.Message = &tgbotapi.Message{Text: "asd"}
	assert(tm.HasText()(u), t)

	u = &tm.Update{}
	assert(!tm.HasPhoto()(u), t)
	u.Message = &tgbotapi.Message{Photo: []tgbotapi.PhotoSize{}}
	assert(tm.HasPhoto()(u), t)

	u = &tm.Update{}
	assert(!tm.HasVoice()(u), t)
	u.Message = &tgbotapi.Message{Voice: &tgbotapi.Voice{}}
	assert(tm.HasVoice()(u), t)

	u = &tm.Update{}
	assert(!tm.HasAudio()(u), t)
	u.Message = &tgbotapi.Message{Audio: &tgbotapi.Audio{}}
	assert(tm.HasAudio()(u), t)

	u = &tm.Update{}
	assert(!tm.HasAnimation()(u), t)
	u.Message = &tgbotapi.Message{Animation: &tgbotapi.Animation{}}
	assert(tm.HasAnimation()(u), t)

	u = &tm.Update{}
	assert(!tm.HasDocument()(u), t)
	u.Message = &tgbotapi.Message{Document: &tgbotapi.Document{}}
	assert(tm.HasDocument()(u), t)

	u = &tm.Update{}
	assert(!tm.HasSticker()(u), t)
	u.Message = &tgbotapi.Message{Sticker: &tgbotapi.Sticker{}}
	assert(tm.HasSticker()(u), t)

	u = &tm.Update{}
	assert(!tm.HasVideo()(u), t)
	u.Message = &tgbotapi.Message{Video: &tgbotapi.Video{}}
	assert(tm.HasVideo()(u), t)

	u = &tm.Update{}
	assert(!tm.HasVideoNote()(u), t)
	u.Message = &tgbotapi.Message{VideoNote: &tgbotapi.VideoNote{}}
	assert(tm.HasVideoNote()(u), t)

	u = &tm.Update{}
	assert(!tm.HasContact()(u), t)
	u.Message = &tgbotapi.Message{Contact: &tgbotapi.Contact{}}
	assert(tm.HasContact()(u), t)

	u = &tm.Update{}
	assert(!tm.HasLocation()(u), t)
	u.Message = &tgbotapi.Message{Location: &tgbotapi.Location{}}
	assert(tm.HasLocation()(u), t)

	u = &tm.Update{}
	assert(!tm.HasVenue()(u), t)
	u.Message = &tgbotapi.Message{Venue: &tgbotapi.Venue{}}
	assert(tm.HasVenue()(u), t)
}

func TestUpdateChatType(t *testing.T) {
	u := &tm.Update{}
	u.Message = &tgbotapi.Message{}

	assert(!tm.IsPrivate()(u), t)
	assert(!tm.IsGroup()(u), t)
	assert(!tm.IsSuperGroup()(u), t)
	assert(!tm.IsGroupOrSuperGroup()(u), t)
	assert(!tm.IsChannel()(u), t)

	u.Message.Chat = &tgbotapi.Chat{}

	assert(!tm.IsPrivate()(u), t)
	assert(!tm.IsGroup()(u), t)
	assert(!tm.IsSuperGroup()(u), t)
	assert(!tm.IsGroupOrSuperGroup()(u), t)
	assert(!tm.IsChannel()(u), t)

	u.Message.Chat.Type = "private"
	assert(tm.IsPrivate()(u), t)
	u.Message.Chat.Type = "group"
	assert(tm.IsGroup()(u), t)
	assert(tm.IsGroupOrSuperGroup()(u), t)
	u.Message.Chat.Type = "supergroup"
	assert(tm.IsSuperGroup()(u), t)
	assert(tm.IsGroupOrSuperGroup()(u), t)
	u.Message.Chat.Type = "channel"
	assert(tm.IsChannel()(u), t)
}

func TestUpdateMembers(t *testing.T) {
	u := &tm.Update{}
	assert(!tm.IsNewChatMembers()(u), t)
	u.Message = &tgbotapi.Message{}
	assert(!tm.IsNewChatMembers()(u), t)
	u.Message.NewChatMembers = []tgbotapi.User{{}}
	assert(tm.IsNewChatMembers()(u), t)

	u = &tm.Update{}
	assert(!tm.IsLeftChatMember()(u), t)
	u.Message = &tgbotapi.Message{}
	assert(!tm.IsLeftChatMember()(u), t)
	u.Message.LeftChatMember = &tgbotapi.User{}
	assert(tm.IsLeftChatMember()(u), t)
}

func TestCombinationFilters(t *testing.T) {
	u := &tm.Update{}
	for _, test := range []struct {
		a bool
		b bool
		r bool
	}{
		{false, false, false},
		{false, true, false},
		{true, false, false},
		{true, true, true},
	} {
		actual := tm.And(func(u *tm.Update) bool { return test.a }, func(u *tm.Update) bool { return test.b })(u)
		assert(actual == test.r, t, fmt.Sprintf("And(%v, %v) should be %v, got %v", test.a, test.b, test.r, actual))
	}
	for _, test := range []struct {
		a bool
		b bool
		r bool
	}{
		{false, false, false},
		{false, true, true},
		{true, false, true},
		{true, true, true},
	} {
		actual := tm.Or(func(u *tm.Update) bool { return test.a }, func(u *tm.Update) bool { return test.b })(u)
		assert(actual == test.r, t, fmt.Sprintf("Or(%v, %v) should be %v, got %v", test.a, test.b, test.r, actual))
	}
	assert(!tm.Not(tm.Any())(u), t)
	assert(tm.Not(tm.Not(tm.Any()))(u), t)
}
