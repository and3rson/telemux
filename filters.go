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
		return u.Message.Text != "" && u.Message.Text[0] != '/'
	}
}

// IsAnyCommand filters updates that look like a command,
// i. e. have some text and start with a slash ("/").
func IsAnyCommand() Filter {
	return func(u *Update) bool {
		return u.Message.Text != "" && u.Message.Text[0] == '/'
	}
}

// IsCommand filters updates that contain a specific command.
// For example, IsCommand("start") will handle a "/start" command.
// This will also allow the user to pass arguments, e. g. "/start foo bar".
func IsCommand(cmd string) Filter {
	return func(u *Update) bool {
		return u.Message.Text == "/"+cmd || strings.HasPrefix(u.Message.Text, "/"+cmd)
	}
}

// IsRegex filters updates that match a regular expression.
// For example, IsRegex("^/get_(\d+)$") will handle commands like "/get_42".
func IsRegex(pattern string) Filter {
	exp := regexp.MustCompile(pattern)
	return func(u *Update) bool {
		return exp.MatchString(u.Message.Text)
	}
}

// IsPhoto filters updates that contain a photo.
func IsPhoto() Filter {
	return func(u *Update) bool {
		return u.Message.Photo != nil
	}
}

// IsLocation filters updates that contain a location.
func IsLocation() Filter {
	return func(u *Update) bool {
		return u.Message.Location != nil
	}
}

// IsPrivate filters updates that are sent in private chats.
func IsPrivate() Filter {
	return func(u *Update) bool {
		return u.Message.Chat.IsPrivate()
	}
}

// IsGroup filters updates that are sent in a group. See also IsGroupOrSuperGroup.
func IsGroup() Filter {
	return func(u *Update) bool {
		return u.Message.Chat.IsGroup()
	}
}

// IsSupergroup filters updates that are sent in a superbroup. See also IsGroupOrSuperGroup.
func IsSuperGroup() Filter {
	return func(u *Update) bool {
		return u.Message.Chat.IsSuperGroup()
	}
}

// IsGroupOrSuperGroup filters updates that are sent in both groups and supergroups.
func IsGroupOrSuperGroup() Filter {
	return func(u *Update) bool {
		return u.Message.Chat.IsGroup() || u.Message.Chat.IsSuperGroup()
	}
}

// IsChannel filters updates that are sent in channels.
func IsChannel() Filter {
	return func(u *Update) bool {
		return u.Message.Chat.IsChannel()
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
