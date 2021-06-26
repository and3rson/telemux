package telemux

import (
	"regexp"
	"strings"
)

// Filter checks Update and returns true if the update satisfies this filter.
type Filter func(u *Update) bool

// Any tells handler to process all updates.
func Any() Filter {
	return func(u *Update) bool {
		return true
	}
}

// IsMessage filters updates that look like message (text, photo, location etc.)
func IsMessage() Filter {
	return func(u *Update) bool {
		return u.Message != nil
	}
}

// IsInlineQuery filters updates that are callbacks from inline queries.
func IsInlineQuery() Filter {
	return func(u *Update) bool {
		return u.InlineQuery != nil
	}
}

// IsCallbackQuery filters updates that are callbacks from button presses.
func IsCallbackQuery() Filter {
	return func(u *Update) bool {
		return u.CallbackQuery != nil
	}
}

// IsEditedMessage filters updates that are edits to existing messages.
func IsEditedMessage() Filter {
	return func(u *Update) bool {
		return u.EditedMessage != nil
	}
}

// IsChannelPost filters updates that are channel posts.
func IsChannelPost() Filter {
	return func(u *Update) bool {
		return u.ChannelPost != nil
	}
}

// IsEditedChannelPost filters updates that are edits to existing channel posts.
func IsEditedChannelPost() Filter {
	return func(u *Update) bool {
		return u.EditedChannelPost != nil
	}
}

// HasText filters updates that look like text,
// i. e. have some text and do not start with a slash ("/").
func HasText() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Text != "" && message.Text[0] != '/'
	}
}

// IsAnyCommandMessage filters updates that look like a command,
// i. e. have some text and start with a slash ("/").
// It also filters new message and excludes edited messages, channel posts, callback queries etc.
func IsAnyCommandMessage() Filter {
	return func(u *Update) bool {
		return u.Message != nil && u.Message.Text != "" && u.Message.Text[0] == '/'
	}
}

// IsCommandMessage filters updates that contain a specific command.
// For example, IsCommandMessage("start") will handle a "/start" command.
// This will also allow the user to pass arguments, e. g. "/start foo bar".
// It also filters only new messages (edited messages, channel posts, callback queries etc are all excluded.)
func IsCommandMessage(cmd string) Filter {
	return func(u *Update) bool {
		return u.Message != nil && (u.Message.Text == "/"+cmd || strings.HasPrefix(u.Message.Text, "/"+cmd))
	}
}

// HasRegex filters updates that match a regular expression.
// For example, HasRegex("^/get_(\d+)$") will handle commands like "/get_42".
func HasRegex(pattern string) Filter {
	exp := regexp.MustCompile(pattern)
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && exp.MatchString(message.Text)
	}
}

// HasPhoto filters updates that contain a photo.
func HasPhoto() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Photo != nil
	}
}

// HasVoice filters updates that contain a voice message.
func HasVoice() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Voice != nil
	}
}

// HasAudio filters updates that contain an audio.
func HasAudio() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Audio != nil
	}
}

// HasAnimation filters updates that contain an animation.
func HasAnimation() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Animation != nil
	}
}

// HasDocument filters updates that contain a document.
func HasDocument() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Document != nil
	}
}

// HasSticker filters updates that contain a sticker.
func HasSticker() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Sticker != nil
	}
}

// HasVideo filters updates that contain a video.
func HasVideo() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Video != nil
	}
}

// HasVideoNote filters updates that contain a video note.
func HasVideoNote() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.VideoNote != nil
	}
}

// HasContact filters updates that contain a contact.
func HasContact() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Contact != nil
	}
}

// HasLocation filters updates that contain a location.
func HasLocation() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Location != nil
	}
}

// HasVenue filters updates that contain a venue.
func HasVenue() Filter {
	return func(u *Update) bool {
		message := u.EffectiveMessage()
		return message != nil && message.Venue != nil
	}
}

// IsPrivate filters updates that are sent in private chats.
func IsPrivate() Filter {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsPrivate()
		}
		return false
	}
}

// IsGroup filters updates that are sent in a group. See also IsGroupOrSuperGroup.
func IsGroup() Filter {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsGroup()
		}
		return false
	}
}

// IsSupergroup filters updates that are sent in a superbroup. See also IsGroupOrSuperGroup.
func IsSuperGroup() Filter {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsSuperGroup()
		}
		return false
	}
}

// IsGroupOrSuperGroup filters updates that are sent in both groups and supergroups.
func IsGroupOrSuperGroup() Filter {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsGroup() || chat.IsSuperGroup()
		}
		return false
	}
}

// IsChannel filters updates that are sent in channels.
func IsChannel() Filter {
	return func(u *Update) bool {
		if chat := u.EffectiveChat(); chat != nil {
			return chat.IsChannel()
		}
		return false
	}
}

// IsNewChatMembers filters updates that have users in NewChatMembers property.
func IsNewChatMembers() Filter {
	return func(u *Update) bool {
		if message := u.EffectiveMessage(); message != nil {
			return message.NewChatMembers != nil && len(*message.NewChatMembers) > 0
		}
		return false
	}
}

// IsLeftChatMember filters updates that have user in LeftChatMember property.
func IsLeftChatMember() Filter {
	return func(u *Update) bool {
		if message := u.EffectiveMessage(); message != nil {
			return message.LeftChatMember != nil
		}
		return false
	}
}

// And filters updates that pass ALL of the provided filters.
func And(filters ...Filter) Filter {
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
func Or(filters ...Filter) Filter {
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
func Not(filter Filter) Filter {
	return func(u *Update) bool {
		return !filter(u)
	}
}
