package gormpersistence

import (
	tm "github.com/and3rson/telemux/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestPersistence(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	p := GORMPersistence{db}
	p.AutoMigrate()

	pk := tm.PersistenceKey{ConversationID: "a", UserID: 13, ChatID: 37}
	if p.GetState(pk) != "" {
		t.Error("State should be \"\"")
	}
	p.SetState(pk, "new_state")
	if p.GetState(pk) != "new_state" {
		t.Error("State should be \"foo\"")
	}
	if !reflect.DeepEqual(p.GetData(pk), tm.Data{}) {
		t.Error("State should be an empty map")
	}
	p.SetData(pk, tm.Data{"foo": "bar"})
	if !reflect.DeepEqual(p.GetData(pk), tm.Data{"foo": "bar"}) {
		t.Error("State should be [foo:bar]")
	}
}
