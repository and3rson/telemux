# telemux
Flexible message router add-on for [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) library.

[![GitHub tag](https://img.shields.io/github/tag/and3rson/telemux.svg?sort=semver)](https://GitHub.com/and3rson/telemux/tags/) [![Go Reference](https://pkg.go.dev/badge/github.com/and3rson/telemux.svg)](https://pkg.go.dev/github.com/and3rson/telemux) [![Build Status](https://travis-ci.com/and3rson/telemux.svg?branch=main)](https://travis-ci.com/and3rson/telemux) [![Maintainability](https://api.codeclimate.com/v1/badges/63d82ddd3151594c3765/maintainability)](https://codeclimate.com/github/and3rson/telemux/maintainability) [![Go Report Card](https://goreportcard.com/badge/github.com/and3rson/telemux)](https://goreportcard.com/report/github.com/and3rson/telemux) [![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable)

![Screenshot](./sample_screenshot.png)

<!-- TOC generated with https://luciopaiva.com/markdown-toc/ -->
# Table of contents

- [Motivation](#motivation)
- [Features](#features)
- [Minimal example](#minimal-example)
- [Documentation](#documentation)
- [Changelog](#changelog)
- [Terminology](#terminology)
  - [Mux](#mux)
  - [Handlers & filters](#handlers--filters)
    - [Combining filters](#combining-filters)
    - [Asynchronous handlers](#asynchronous-handlers)
    - [Consuming update from filters](#consuming-update-from-filters)
  - [Conversations & persistence](#conversations--persistence)
  - [Error handling](#error-handling)
- [Tips & common pitfalls](#tips--common-pitfalls)
  - [tgbotapi.Update vs tm.Update confusion](#tgbotapiupdate-vs-tmupdate-confusion)
  - [Getting user/chat/message object from update](#getting-userchatmessage-object-from-update)
  - [Properly filtering updates](#properly-filtering-updates)

# Motivation

This library serves as an addition to the [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) library.
I strongly recommend you to take a look at it since telemux is mostly an extension to it.

Patterns such as handlers, persistence & filters were inspired by a wonderful [python-telegram-bot](https://github.com/python-telegram-bot/python-telegram-bot) library.

This project is in early beta stage. Contributions are welcome! Feel free to submit an issue if you have any questions, suggestions or simply want to help.

# Features

- Extension for [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) library, meaning you'll still use all of its features
- Designed with statelessness in mind
- Extensible handler configuration inspired by [python-telegram-bot](https://github.com/python-telegram-bot/python-telegram-bot) library
- Conversations (aka Dialogs) based on finite-state machines (see [./examples/album_conversation/main.go](./examples/album_conversation/main.go))
- Pluggable persistence for conversations. E. g. you can use database to store the states & intermediate values of conversations (see [./examples/album_conversation/main.go](./examples/album_conversation/main.go) and [./persistence.go](./persistence.go))
- Support for GORM as a persistence backend via ![gormpersistence](./gormpersistence) module
- Flexible handler filtering. E. g. `And(Or(HasText(), HasPhoto()), IsPrivate())` will only accept direct messages containing photo or text (see [./filters.go](./filters.go))

# Minimal example

```go
package main

import (
    tm "github.com/and3rson/telemux"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
    "log"
    "os"
)

func main() {
    // This part is a boilerplate from go-telegram-bot-api library.
    bot, _ := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
    bot.Debug = true
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates, _ := bot.GetUpdatesChan(u)

    // Create a multiplexer with two handlers: one for command and one for all messages.
    // If a handler cannot handle the update (fails the filter),
    // multiplexer will proceed to the next handler.
    mux := tm.NewMux().
        AddHandler(tm.NewHandler(
            tm.IsCommandMessage("start"),
            func(u *tm.Update) {
                bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Hello! Say something. :)"))
            },
        )).
        AddHandler(tm.NewHandler(
            tm.Any(),
            func(u *tm.Update) {
                bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "You said: "+u.Message.Text))
            },
        ))
    // Dispatch all telegram updates to multiplexer
    for update := range updates {
        mux.Dispatch(bot, update)
    }
}
```

# Documentation

The documentation is available [here](https://pkg.go.dev/github.com/and3rson/telemux).

Examples are available [here](./examples/).

# Changelog

Changelog is available [here](./CHANGELOG.md).

# Terminology

## Mux
Mux (multiplexer) is a "router" for instances of `tgbotapi.Update`.

It allows you to register handlers and will take care to choose an appropriate handler based on the incoming update.

In order to work, you must dispatch messages (that come from go-telegram-bot-api channel):

```go
mux := tm.NewMux()
// ...
// add handlers to mux here
// ...
updates, _ := bot.GetUpdatesChan(u)
for update := range updates {
    mux.Dispatch(bot, update)
}
```

## Handlers & filters

Handler consists of filter and handle-function.

Handler's filter decides whether this handler can handle the incoming update.
If so, handle-function is called. Otherwise multiplexer will proceed to the next handler.

Filters are divided in two groups: content filters (starting with "Has", such as `HasPhoto()`, `HasAudio()`, `HasSticker()` etc)
and update type filters (starting with "Is", such as `IsEditedMessage()`, `IsInlineQuery()` or `IsGroupOrSuperGroup()`).

There is also a special filter `Any()` which makes handler accept all updates.

### Combining filters

Filters can be chained using `And`, `Or`, and `Not` meta-filters. For example:

```go
mux := tm.NewMux()

// Add handler that accepts photos sent to the bot in a private chat:
mux.AddHandler(And(tm.IsPrivate(), tm.HasPhoto()), func(u *tm.Update) { /* ... */ })

// Add handler that accepts photos and text messages:
mux.AddHandler(Or(tm.HasText(), tm.HasPhoto()), func(u *tm.Update) { /* ... */ })

// Since filters are plain functions, you can easily implement them yourself.
// Below we add handler that allows onle a specific user to call "/restart" command:
mux.AddHandler(tm.NewHandler(
    tm.And(tm.IsCommandMessage("restart"), func(u *tm.Update) bool {
        return u.Message.From.ID == 3442691337
    }),
    func(u *tm.Update) { /* ... */ },
))
```

### Asynchronous handlers

If you want your handler to be executed in a goroutine, use `tm.NewAsyncHandler`. It's similar to wrapping a handler function in an anonymous goroutine invocation:

```go
mux := tm.NewMux()
// Option 1: using NewAsyncHandler
mux.AddHandler(tm.NewAsyncHandler(
    tm.IsCommandMessage("do_work"),
    func(u *tm.Update) {
        // Slow operation
    },
))
// Option 2: wrapping manually
mus.AddHandler(tm.NewHandler(
    tm.IsCommandMessage("do_work"),
    func(u *tm.Update) {
        go func() {
            // Slow operation
        }
    },
))
```

### Consuming update from filters

Although main purpose of filters is to decide whether a handler can process an update, there are often situations
when a filter needs to mark update as "consumed" (i. e. "processed") and prevent its further processing *without actually invoking the handler*.
In this case filters can call `u.Consume()` on Update and return `false`. This will prevent handler from executing and also
prevent `Mux` from going further down the handler chain. Here's an example:

```go
mux.AddHandler(tm.NewHandler(
    tm.IsCommandMessage("do_work"),
    func(u *tm.Update) {
        if u.EffectiveUser().ID != 3442691337 { // Boilerplate code that will be copy-pasted way too much
            u.Bot.Send(tgbotapi.Message(u.EffectiveChat().ID, "You are not allowed to ask me to work!"))
            return
        }
        if !u.EffectiveChat().IsPrivate() { // Another boilerplate code
            u.Bot.Send(tgbotapi.Message(u.EffectiveChat().ID, "I do not accept commands in group chats. Send me a PM."))
            return
        }
        // Do actual work
    },
))
```

To avoid repeating boilerplate checks like `if user is not "3442691337" then send error and stop", you can "consume" the update from within a filter.
The above code can be rewritten as follows:

```go
// CheckAdmin is a reusable filter that not only checks for user's ID but marks update as processed as well
func CheckAdmin(u *tm.Update) {
    if u.EffectiveUser().ID != 3442691337 { // Boilerplate code that will be copy-pasted way too much
        u.Bot.Send(tgbotapi.Message(u.EffectiveChat().ID, "You are not allowed to ask me to work!"))
        return false
    }
    return true
}

// CheckPrivate is a reusable filter that not only checks for private chat but marks update as processed as well
func CheckPrivate(u *tm.Update) {
    if !u.EffectiveChat().IsPrivate() { // Boilerplate code that will be copy-pasted way too much
        u.Bot.Send(tgbotapi.Message(u.EffectiveChat().ID, "I do not accept commands in group chats. Send me a PM."))
        return false
    }
    return true
}

// ...

mux.AddHandler(tm.NewHandler(
    And(tm.IsCommandMessage("do_work"), CheckAdmin, CheckPrivate),
    func(u *tm.Update) {
        // Do actual work
    },
))
```


## Conversations & persistence

Conversations are handlers on steroids based on the finite-state machine pattern.

They allow you to have complex dialog interactions with different handlers.

Persistence interface tells conversation where to store & how to retrieve the current state of the conversation, i. e. which "step" the given user is currently at.

To create a ConversationHandler you need to provide the following:

- `conversationID string` - identifier that distinguishes this conversation from the others.

    The main goal of this identifier is to allow persistence to keep track of different conversation states independently without mixing them together.

- `persistence Persistence` - defines where to store conversation state & intermediate inputs from the user.

    Without persistence, a conversation would not be able to "remember" what "step" the user is at.

    Persistence is also useful when you want to collect some data from the user step-by-step).

    Two convenient implementations of `Persistence` are available out of the box: `LocalPersistence` & `FilePersistence`.

    Telemux also supports GORM persistence. If you use GORM, you can store conversation states & data in your database by using `GORMPersistence` from a ![gormpersistence](./gormpersistence) module.

- `states map[string][]*TransitionHandler` - defines which TransitionHandlers to use in what state.

    States are usually strings like "upload_photo", "send_confirmation", "wait_for_text" and describe the "step" the user is currently at.
    Empty string (`""`) shoulb be used as an initial/final state (i. e. if the conversation has not started yet or has already finished.)

    For each state you can provide a list of at least one TransitionHandler. If none of the handlers can handle the update, the default handlers are attempted (see below).

    In order to switch to a different state your TransitionHandler must call `ctx.SetState("STATE_NAME") ` replacing STATE_NAME with the name of the state you want to switch into.

    Conversation data can be accessed with `ctx.GetData()` and updated with `ctx.SetData(newData)`.


- `defaults []*TransitionHandler` - these handlers are "appended" to every state.

    Useful to handle commands such as "/cancel" or to display some default message.

See [./examples/album_conversation/main.go](./examples/album_conversation/main.go) for a conversation example.

## Error handling

By default, panics in handlers are propagated all the way to the top (`Dispatch` method).

In order to intercept all panics in your handlers globally and handle them gracefully, register your function using `SetRecoverer`:

```go
mux := tm.NewMux()
# ...
mux.SetRecoverer(func(u *tm.Update, err error) {
    fmt.Printf("An error occured: %s", err)
})
```

# Tips & common pitfalls

## tgbotapi.Update vs tm.Update confusion

Since `Update` struct from go-telegram-bot-api already provides most of the functionality, telemux implements its own `Update` struct
which embeds the `Update` from go-telegram-bot-api. Main reason for this is to add some extra convenient methods and include Bot instance
with every update.

## Getting user/chat/message object from update

When having handlers for wide filters (e. g. `Or(And(HasText(), IsEditedMessage()), IsInlineQuery())`) you may often fall
in situations when you need to check for multiple user/chat/message attributes. In such situations sender's data may
be in one of few places depending on which update has arrived: `u.Message.From`, `u.EditedMessage.From`, or `u.InlineQuery.From`.
Similar issue applies to fetching actual chat info or message object from an update.

In such cases it's highly recommended to use functions such as `EffectiveChat()` (see the [update](./update.go) module for more info):

```go
// Bad:
fmt.Println(u.Message.Chat.ID) // u.Message may be nil

// Better, but not so DRY:
chatId int64
if u.Message != nil {
    chatId = u.Message.Chat.ID
} else if u.EditedMessage != nil {
    chatId = u.EditedMessage.Chat.ID
} else if u.CallbackQuery != nil {
    chatId = u.CallbackQuery.Chat.ID
} // And so on... Duh.
fmt.Println(chatId)

// Best:
chat := u.EffectiveChat()
if chat != nil {
    fmt.Println(chat.ID)
}
```

## Properly filtering updates

Keep in mind that using content filters such as `HasText()`, `HasPhoto()`, `HasLocation()`, `HasVoice()` etc does not guarantee
that the `Update` describes an actual new message. In fact, an `Update` also happens when a user edits his message!
Thus your handler will be executed even if a user just edited one of his messages.

To avoid situations like these, make sure to use filters such as `IsMessage()`, `IsEditedMessage()`, `IsCallbackQuery()` etc
in conjunction with content filters. For example:

```go
tm.NewHandler(HasText(), func(u *tm.Update) { /* ... */ }) // Will handle new messages, updated messages, channel posts & channel post edits which contain text
tm.NewHandler(And(IsMessage(), HasText()), func(u *tm.Update) { /* ... */ }) // Will handle new messages that contain text
tm.NewHandler(And(IsEditedMessage(), HasText()), func(u *tm.Update) { /* ... */ }) // Will handle edited that which contain text
```

The only exceptions are `IsCommandMessage("...")` and `IsAnyCommandMessage()` filters. Since it does not make sense to react to edited messages that contain
commands, this filter also checks if the update designates a new message and not an edited message, inline query, callback query etc.
This means you can safely use `IsCommandMessage("my_command")` without joining it with the `IsMessage()` filter:

```go
IsCommandMessage("my_command") // OK: IsCommand() already checks for IsMessage()
And(IsCommandMessage("start"), IsMessage()) // IsMessage() is unnecessary
And(IsCommandMessage("start"), Not(IsEditedMessage())) // Not(IsEditedMessage()) is unnecessary
```
