package gormpersistence

import (
	tm "github.com/and3rson/telemux/v2"
	datatypes "gorm.io/datatypes"
)

// ConversationState is a model that contains conversation states for users.
type ConversationState struct {
	tm.PersistenceKey
	State string `gorm:"not null"`
}

// ConversationData is a model that contains conversation data for users.
type ConversationData struct {
	tm.PersistenceKey
	Data datatypes.JSONMap `gorm:"not null"`
}
