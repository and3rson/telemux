# gormpersistence

Support for GORM as persistence backend.

## Installation

```sh
go get github.com/and3rson/telemux/gormpersistence
```

## Example usage

```go
package main

import (
    "fmt"
    "log"
    "os"

    tm "github.com/and3rson/telemux"
    "github.com/and3rson/telemux/gormpersistence"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
    db, _ := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
    // Create GORMPersistence
    p = gormpersistence.GORMPersistence{db}

    // Create required tables
    p.AutoMigrate()
    // ...alias for:
    // db.AutoMigrate(&gormpersistence.ConversationState{}, &gormpersistence.ConversationData{})

    bot, _ := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
    bot.Debug = true
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, _ := bot.GetUpdatesChan(u)

    mux := tm.NewMux().
        AddHandler(tm.NewConversationHandler(
            "upload_photo_dialog",
            p, // We provide GORMPersistence as persistence backend for our conversation handler
            map[string][]*tm.TransitionHandler{
                "": {
                    // Handlers for initial state, i. e. for conversation activation
                    // ...
                },
                "state1": {
                    // Handlers for state1
                    // ...
                },
                "state2": {
                    // Handlers for state2
                    // ...
                }
            },
            []*tm.TransitionHandler{
                // Default handlers
            },
        ))

    for update := range updates {
        mux.Dispatch(update)
    }
}
```
