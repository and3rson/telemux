package gormpersistence

import (
	tm "github.com/and3rson/telemux"
	datatypes "gorm.io/datatypes"
)

type ConversationState struct {
	tm.PersistenceKey
	State string `gorm:"not null"`
}

type ConversationData struct {
	tm.PersistenceKey
	Data datatypes.JSONMap `gorm:"not null"`
}
