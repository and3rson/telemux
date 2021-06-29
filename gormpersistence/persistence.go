package gormpersistence

import (
	tm "github.com/and3rson/telemux"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GORMPersistence struct {
	DB *gorm.DB
}

func (p *GORMPersistence) AutoMigrate() error {
	return p.DB.AutoMigrate(&ConversationState{}, &ConversationData{})
}

func (p *GORMPersistence) GetState(pk tm.PersistenceKey) string {
	var stateRecord ConversationState
	p.DB.Where(pk).Attrs(ConversationState{State: ""}).FirstOrCreate(&stateRecord)
	return stateRecord.State
}

func (p *GORMPersistence) SetState(pk tm.PersistenceKey, state string) {
	p.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&ConversationState{
		PersistenceKey: pk,
		State:          state,
	})
}

func (p *GORMPersistence) GetData(pk tm.PersistenceKey) tm.Data {
	var dataRecord ConversationData
	p.DB.Where(pk).Attrs(ConversationData{Data: datatypes.JSONMap{}}).FirstOrCreate(&dataRecord)
	return dataRecord.Data
}

func (p *GORMPersistence) SetData(pk tm.PersistenceKey, data tm.Data) {
	p.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&ConversationData{
		PersistenceKey: pk,
		Data:           data,
	})
}
