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

// IsText filters updates that look like text,
// i. e. have some text and do not start with a slash ("/").
func IsText() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Text != "" && message.Text[0] != '/'
	}
}

// IsAnyCommand filters updates that look like a command,
// i. e. have some text and start with a slash ("/").
func IsAnyCommand() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Text != "" && message.Text[0] == '/'
	}
}

// IsCommand filters updates that contain a specific command.
// For example, IsCommand("start") will handle a "/start" command.
// This will also allow the user to pass arguments, e. g. "/start foo bar".
func IsCommand(cmd string) Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && (message.Text == "/"+cmd || strings.HasPrefix(message.Text, "/"+cmd))
	}
}

// IsRegex filters updates that match a regular expression.
// For example, IsRegex("^/get_(\d+)$") will handle commands like "/get_42".
func IsRegex(pattern string) Filter {
	exp := regexp.MustCompile(pattern)
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && exp.MatchString(message.Text)
	}
}

// IsPhoto filters updates that contain a photo.
func IsPhoto() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Photo != nil
	}
}

// IsVoice filters updates that contain a voice message.
func IsVoice() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Voice != nil
	}
}

// IsAudio filters updates that contain an audio.
func IsAudio() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Audio != nil
	}
}

// IsAnimation filters updates that contain an animation.
func IsAnimation() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Animation != nil
	}
}

// IsDocument filters updates that contain a document.
func IsDocument() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Document != nil
	}
}

// IsSticker filters updates that contain a sticker.
func Sticker() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Sticker != nil
	}
}

// IsVideo filters updates that contain a video.
func IsVideo() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Video != nil
	}
}

// IsVideoNote filters updates that contain a video note.
func IsVideoNote() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.VideoNote != nil
	}
}

// IsContact filters updates that contain a contact.
func IsContact() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Contact != nil
	}
}

// IsLocation filters updates that contain a location.
func IsLocation() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Location != nil
	}
}

// IsVenue filters updates that contain a venue.
func IsVenue() Filter {
	return func(u *Update) bool {
		message := GetEffectiveMessage(u)
		return message != nil && message.Venue != nil
	}
}

// IsPrivate filters updates that are sent in private chats.
func IsPrivate() Filter {
	return func(u *Update) bool {
		if chat := GetEffectiveChat(u); chat != nil {
			return chat.IsPrivate()
		}
		return false
	}
}

// IsGroup filters updates that are sent in a group. See also IsGroupOrSuperGroup.
func IsGroup() Filter {
	return func(u *Update) bool {
		if chat := GetEffectiveChat(u); chat != nil {
			return chat.IsGroup()
		}
		return false
	}
}

// IsSupergroup filters updates that are sent in a superbroup. See also IsGroupOrSuperGroup.
func IsSuperGroup() Filter {
	return func(u *Update) bool {
		if chat := GetEffectiveChat(u); chat != nil {
			return chat.IsSuperGroup()
		}
		return false
	}
}

// IsGroupOrSuperGroup filters updates that are sent in both groups and supergroups.
func IsGroupOrSuperGroup() Filter {
	return func(u *Update) bool {
		if chat := GetEffectiveChat(u); chat != nil {
			return chat.IsGroup() || chat.IsSuperGroup()
		}
		return false
	}
}

// IsChannel filters updates that are sent in channels.
func IsChannel() Filter {
	return func(u *Update) bool {
		if chat := GetEffectiveChat(u); chat != nil {
			return chat.IsChannel()
		}
		return false
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

// EditedChannelPost filters updates that are edits to existing channel posts.
func IsEditedChannelPost() Filter {
	return func(u *Update) bool {
		return u.EditedChannelPost != nil
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
