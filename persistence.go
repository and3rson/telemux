package telemux

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Data is an alias for map[string]interface{}.
type Data = map[string]interface{}

// ConversationPersistence interface tells conversation where to store & how to retrieve the current state of the conversation,
// i. e. which "step" the given user is currently at.
type ConversationPersistence interface {
	// GetState & SetState tell conversation handlers how to retrieve & set conversation state.
	GetState(pk PersistenceKey) string
	SetState(pk PersistenceKey, state string)
	// GetConversationData & SetConversationData allow conversation handlers to store intermediate data.
	GetData(pk PersistenceKey) Data
	SetData(pk PersistenceKey, data Data)
}

// PersistenceContext allows handler to get/set conversation data & change conversation state.
type PersistenceContext struct {
	Persistence ConversationPersistence
	PK          PersistenceKey
}

// GetData returns data of current conversation.
func (c *PersistenceContext) GetData() Data {
	return c.Persistence.GetData(c.PK)
}

// SetData updates data of current conversation.
func (c *PersistenceContext) SetData(data Data) {
	c.Persistence.SetData(c.PK, data)
}

// ClearData clears data of current conversation.
func (c *PersistenceContext) ClearData() {
	c.Persistence.SetData(c.PK, make(map[string]interface{}))
}

// SetState changes state of current conversation.
func (c *PersistenceContext) SetState(state string) {
	c.Persistence.SetState(c.PK, state)
}

// PersistenceKey contains user & chat IDs. It is used to identify conversations with different users in different chats.
type PersistenceKey struct {
	ConversationID string `gorm:"primaryKey;autoIncrement:false"`
	UserID         int    `gorm:"primaryKey;autoIncrement:false"`
	ChatID         int64  `gorm:"primaryKey;autoIncrement:false"`
}

// String returns a string in form "USER_ID:CHAT_ID".
func (k PersistenceKey) String() string {
	return fmt.Sprintf("%s:%d:%d", k.ConversationID, k.UserID, k.ChatID)
}

// MarshalText marshals persistence for use in map keys
func (k PersistenceKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// UnmarshalText unmarshals persistence from "CONV:USER:CHAT" string
func (k *PersistenceKey) UnmarshalText(b []byte) error {
	parts := strings.Split(string(b), ":")
	chatIDstr := parts[len(parts)-1]
	userIDstr := parts[len(parts)-2]
	k.ConversationID = strings.Join(parts[:len(parts)-2], ":")

	chatID, err := strconv.ParseInt(chatIDstr, 10, 64)
	if err != nil {
		panic(err)
	}
	k.ChatID = chatID
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		panic(err)
	}
	k.UserID = userID
	return nil
}

// LocalPersistence is an implementation of Persistence.
// It stores conversation states & conversation data in memory.
//
// All data in this implementation of persistence is lost if an application is restarted.
// If you want to store the data permanently you will need to implement your own Persistence
// which will use redis, database or something else to store states & conversation data.
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

// GetState returns conversation state from memory
func (p *LocalPersistence) GetState(pk PersistenceKey) string {
	state, ok := p.States[pk]
	if !ok {
		return ""
	}
	return state
}

// SetState stores conversation state in memory
func (p *LocalPersistence) SetState(pk PersistenceKey, state string) {
	p.States[pk] = state
}

// GetData returns conversation data from memory
func (p *LocalPersistence) GetData(pk PersistenceKey) Data {
	data, ok := p.Data[pk]
	if !ok {
		p.Data[pk] = make(Data)
		return p.Data[pk]
	}
	return data
}

// SetData stores conversation data in memory
func (p *LocalPersistence) SetData(pk PersistenceKey, data Data) {
	p.Data[pk] = data
}

// FilePersistence is an implementation of Persistence.
// It stores conversation states & conversation data in file.
type FilePersistence struct {
	mutex    *sync.Mutex
	Filename string
}

// NewFilePersistence creates new instance of FilePersistence.
func NewFilePersistence(filename string) *FilePersistence {
	return &FilePersistence{
		&sync.Mutex{},
		filename,
	}
}

type filePersistenceContent struct {
	States map[PersistenceKey]string `json:"states"`
	Data   map[PersistenceKey]Data   `json:"data"`
}

func (p *FilePersistence) readContent() *filePersistenceContent {
	if _, err := os.Stat(p.Filename); err != nil {
		if os.IsNotExist(err) {
			err := ioutil.WriteFile(p.Filename, []byte("{}"), 0644)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	data, err := ioutil.ReadFile(p.Filename)
	if err != nil {
		panic(err)
	}
	content := filePersistenceContent{
		make(map[PersistenceKey]string),
		make(map[PersistenceKey]Data),
	}
	json.Unmarshal(data, &content)
	return &content
}

func (p *FilePersistence) writeContent(content *filePersistenceContent) {
	data, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(p.Filename, data, 0644)
	if err != nil {
		panic(err)
	}
}

// GetState reads conversation state from file
func (p *FilePersistence) GetState(pk PersistenceKey) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	content := p.readContent()
	state, ok := content.States[pk]
	if !ok {
		return ""
	}
	return state
}

// SetState writes conversation state to file
func (p *FilePersistence) SetState(pk PersistenceKey, state string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	content := p.readContent()
	content.States[pk] = state
	p.writeContent(content)
}

// GetData reads conversation data from file
func (p *FilePersistence) GetData(pk PersistenceKey) Data {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	content := p.readContent()
	data, ok := content.Data[pk]
	if !ok {
		return make(Data)
	}
	return data
}

// SetData writes conversation data to file
func (p *FilePersistence) SetData(pk PersistenceKey, data Data) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	content := p.readContent()
	content.Data[pk] = data
	p.writeContent(content)
}
