package gormpersistence

import (
	tm "github.com/and3rson/telemux/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GORMPersistence is an implementation of Persistence.
// It stores conversation states & conversation data in database via GORM.
type GORMPersistence struct {
	DB *gorm.DB
}

// AutoMigrate creates tables for ConversationState & ConversationData models
func (p *GORMPersistence) AutoMigrate() error {
	return p.DB.AutoMigrate(&ConversationState{}, &ConversationData{})
}

// GetState reads conversation state from database
func (p *GORMPersistence) GetState(pk tm.PersistenceKey) string {
	var stateRecord ConversationState
	p.DB.Where(pk).Attrs(ConversationState{State: ""}).FirstOrCreate(&stateRecord)
	return stateRecord.State
}

// SetState writes conversation state to database
func (p *GORMPersistence) SetState(pk tm.PersistenceKey, state string) {
	p.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&ConversationState{
		PersistenceKey: pk,
		State:          state,
	})
}

// GetData reads conversation data from database
func (p *GORMPersistence) GetData(pk tm.PersistenceKey) tm.Data {
	var dataRecord ConversationData
	p.DB.Where(pk).Attrs(ConversationData{Data: datatypes.JSONMap{}}).FirstOrCreate(&dataRecord)
	return dataRecord.Data
}

// SetData writes conversation data to database
func (p *GORMPersistence) SetData(pk tm.PersistenceKey, data tm.Data) {
	p.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&ConversationData{
		PersistenceKey: pk,
		Data:           data,
	})
}
