package telemux

import "fmt"

// Data is an alias for map[string]interface{}.
type Data = map[string]interface{}

// Persistence interface tells conversation where to store & how to retrieve the current state of the conversation,
// i. e. which "step" the given user is currently at.
type ConversationPersistence interface {
	// GetState & SetState tell conversation handlers how to retrieve & set conversation state.
	GetState(conversationID string, pk PersistenceKey) string
	SetState(conversationID string, pk PersistenceKey, state string)
	// GetConversationData & SetConversationData allow conversation transition handlers to store intermediate data.
	GetConversationData(conversationID string, pk PersistenceKey) Data
	SetConversationData(conversationID string, pk PersistenceKey, data Data)
}

// PersistenceKey contains user & chat IDs. It is used to identify conversations with different users in different chats.
type PersistenceKey struct {
	UserID int
	ChatID int64
}

// String returns a string in form "USER_ID:CHAT_ID".
func (k PersistenceKey) String() string {
	return fmt.Sprintf("%d:%d", k.UserID, k.ChatID)
}

// LocalPersistence is an implementation of Persistence.
// It stores conversation states & conversation data in memory.
//
// All data in this implementation of persistence is lost if an application is restarted.
// If you want to store the data permanently you will need to implement your own Persistence
// which will use files, redis, database or something else to store states & conversation data.
type LocalPersistence struct {
	States map[PersistenceKey]string
	Data   map[PersistenceKey]Data
}

// NewLocalPersistence creates new instance of LocalPersistence.
func NewLocalPersistence() *LocalPersistence {
	return &LocalPersistence{
		make(map[PersistenceKey]string),
		make(map[PersistenceKey]Data),
	}
}

func (p *LocalPersistence) GetState(conversationID string, pk PersistenceKey) string {
	state, ok := p.States[pk]
	if !ok {
		return ""
	}
	return state
}

func (p *LocalPersistence) SetState(conversationID string, pk PersistenceKey, state string) {
	p.States[pk] = state
}

func (p *LocalPersistence) GetConversationData(conversationID string, pk PersistenceKey) Data {
	data, ok := p.Data[pk]
	if !ok {
		p.Data[pk] = make(Data)
		return p.Data[pk]
	}
	return data
}

func (p *LocalPersistence) SetConversationData(conversationID string, pk PersistenceKey, data Data) {
	p.Data[pk] = data
}
