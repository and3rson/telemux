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

// Persistence interface tells conversation where to store & how to retrieve the current state of the conversation,
// i. e. which "step" the given user is currently at.
type ConversationPersistence interface {
	// GetState & SetState tell conversation handlers how to retrieve & set conversation state.
	GetState(pk PersistenceKey) string
	SetState(pk PersistenceKey, state string)
	// GetConversationData & SetConversationData allow conversation transition handlers to store intermediate data.
	GetConversationData(pk PersistenceKey) Data
	SetConversationData(pk PersistenceKey, data Data)
}

// PersistenceKey contains user & chat IDs. It is used to identify conversations with different users in different chats.
type PersistenceKey struct {
	ConversationID string
	UserID         int
	ChatID         int64
}

// String returns a string in form "USER_ID:CHAT_ID".
func (k PersistenceKey) String() string {
	return fmt.Sprintf("%s:%d:%d", k.ConversationID, k.UserID, k.ChatID)
}

func (k PersistenceKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

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

func (p *LocalPersistence) GetState(pk PersistenceKey) string {
	state, ok := p.States[pk]
	if !ok {
		return ""
	}
	return state
}

func (p *LocalPersistence) SetState(pk PersistenceKey, state string) {
	p.States[pk] = state
}

func (p *LocalPersistence) GetConversationData(pk PersistenceKey) Data {
	data, ok := p.Data[pk]
	if !ok {
		p.Data[pk] = make(Data)
		return p.Data[pk]
	}
	return data
}

func (p *LocalPersistence) SetConversationData(pk PersistenceKey, data Data) {
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

func (p *FilePersistence) SetState(pk PersistenceKey, state string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	content := p.readContent()
	content.States[pk] = state
	p.writeContent(content)
}

func (p *FilePersistence) GetConversationData(pk PersistenceKey) Data {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	content := p.readContent()
	data, ok := content.Data[pk]
	if !ok {
		return make(Data)
	}
	return data
}

func (p *FilePersistence) SetConversationData(pk PersistenceKey, data Data) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	content := p.readContent()
	content.Data[pk] = data
	p.writeContent(content)
}
