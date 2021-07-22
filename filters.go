package telemux

import (
	"regexp"
)

// FilterFunc is used to check if this update should be processed by handler.
type FilterFunc func(u *Update) bool

var commandRegex = regexp.MustCompile("^/([0-9a-zA-Z_]+)(@[0-9a-zA-Z_]{3,})?")

// Any tells handler to process all updates.
func Any() FilterFunc {
	return func(u *Update) bool {
		return true
	}
}

// IsMessage filters updates that look like message (text, photo, location etc.)
func IsMessage() FilterFunc {
	return func(u *Update) bool {
		return u.Message != nil
	}
}

// IsInlineQuery filters updates that are callbacks from inline queries.
func IsInlineQuery() FilterFunc {
	return func(u *Update) bool {
		return u.InlineQuery != nil
	}
}

// IsCallbackQuery filters updates that are callbacks from button presses.
func IsCallbackQuery() FilterFunc {
	return func(u *Update) bool {
		return u.CallbackQuery != nil
	}
}

// IsEditedMessage filters updates that are edits to existing messages.
func IsEditedMessage() FilterFunc {
	return func(u *Update) bool {
		return u.EditedMessage != nil
	}
}

// IsChannelPost filters updates that are channel posts.
func IsChannelPost() FilterFunc {
	return func(u *Update) bool {
		return u.ChannelPost != nil
	}
}

// IsEditedChannelPost filters updates that are edits to existing channel posts.
func IsEditedChannelPost() FilterFunc {
	return func(u *Update) bool {
		return u.EditedChannelPost != nil
	}
}

// HasText filters updates that look like text,
// i. e. have some text and do not start with a slash ("/").
func HasText() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Text != "" && message.Text[0] != '/'
	}
}

// IsAnyCommandMessage filters updates that contain a message and look like a command,
// i. e. have some text and start with a slash ("/").
// If command contains bot username, it is also checked.
func IsAnyCommandMessage() FilterFunc {
	return And(IsMessage(), func(u *Update) bool {
		matches := commandRegex.FindStringSubmatch(u.Message.Text)
		if len(matches) == 0 {
			return false
		}
		botName := matches[2]
		if botName != "" && botName != "@"+u.Bot.Self.UserName {
			return false
		}
		return true
	})
}

// IsCommandMessage filters updates that contain a specific command.
// For example, IsCommandMessage("start") will handle a "/start" command.
// This will also allow the user to pass arguments, e. g. "/start foo bar".
// Commands in format "/start@bot_name" and "/start@bot_name foo bar" are also supported.
// If command contains bot username, it is also checked.
func IsCommandMessage(cmd string) FilterFunc {
	return And(IsAnyCommandMessage(), func(u *Update) bool {
		matches := commandRegex.FindStringSubmatch(u.Message.Text)
		actualCmd := matches[1]
		return actualCmd == cmd
	})
}

// HasRegex filters updates that match a regular expression.
// For example, HasRegex("^/get_(\d+)$") will handle commands like "/get_42".
func HasRegex(pattern string) FilterFunc {
	exp := regexp.MustCompile(pattern)
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && exp.MatchString(message.Text)
	}
}

// HasPhoto filters updates that contain a photo.
func HasPhoto() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Photo != nil
	}
}

// HasVoice filters updates that contain a voice message.
func HasVoice() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Voice != nil
	}
}

// HasAudio filters updates that contain an audio.
func HasAudio() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Audio != nil
	}
}

// HasAnimation filters updates that contain an animation.
func HasAnimation() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Animation != nil
	}
}

// HasDocument filters updates that contain a document.
func HasDocument() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Document != nil
	}
}

// HasSticker filters updates that contain a sticker.
func HasSticker() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Sticker != nil
	}
}

// HasVideo filters updates that contain a video.
func HasVideo() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Video != nil
	}
}

// HasVideoNote filters updates that contain a video note.
func HasVideoNote() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.VideoNote != nil
	}
}

// HasContact filters updates that contain a contact.
func HasContact() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Contact != nil
	}
}

// HasLocation filters updates that contain a location.
func HasLocation() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Location != nil
	}
}

// HasVenue filters updates that contain a venue.
func HasVenue() FilterFunc {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Venue != nil
	}
}

// IsPrivate filters updates that are sent in private chats.
func IsPrivate() FilterFunc {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsPrivate()
		}
		return false
	}
}

// IsGroup filters updates that are sent in a group. See also IsGroupOrSuperGroup.
func IsGroup() FilterFunc {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsGroup()
		}
		return false
	}
}

// IsSuperGroup filters updates that are sent in a superbroup. See also IsGroupOrSuperGroup.
func IsSuperGroup() FilterFunc {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsSuperGroup()
		}
		return false
	}
}

// IsGroupOrSuperGroup filters updates that are sent in both groups and supergroups.
func IsGroupOrSuperGroup() FilterFunc {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsGroup() || chat.IsSuperGroup()
		}
		return false
	}
}

// IsChannel filters updates that are sent in channels.
func IsChannel() FilterFunc {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsChannel()
		}
		return false
	}
}

// IsNewChatMembers filters updates that have users in NewChatMembers property.
func IsNewChatMembers() FilterFunc {
	return func(u *Update) bool {
		if message := u.EffectiveMessage(); message != nil {
			return message.NewChatMembers != nil && len(*message.NewChatMembers) > 0
		}
		return false
	}
}

// IsLeftChatMember filters updates that have user in LeftChatMember property.
func IsLeftChatMember() FilterFunc {
	return func(u *Update) bool {
		if message := u.EffectiveMessage(); message != nil {
			return message.LeftChatMember != nil
		}
		return false
	}
}

// And filters updates that pass ALL of the provided filters.
func And(filters ...FilterFunc) FilterFunc {
	return func(u *Update) bool {
		for _, filter := range filters {
			if !filter(u) {
				return false
			}
		}
		return true
	}
}

// Or filters updates that pass ANY of the provided filters.
func Or(filters ...FilterFunc) FilterFunc {
	return func(u *Update) bool {
		for _, filter := range filters {
			if filter(u) {
				return true
			}
		}
		return false
	}
}

// Not filters updates that do not pass the provided filter.
func Not(filter FilterFunc) FilterFunc {
	return func(u *Update) bool {
		return !filter(u)
	}
}
