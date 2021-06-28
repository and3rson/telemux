## v1.3.0

- readme: update docs on persistence & recoverer
- handlers: use persistence context
- mux: add recoverer
- persistence: add persistence context
- examples: add error_handling

## v1.2.0

- readme: cleanup
- examples: update Dispatch call
- update: embed Update instead of *Update
- mux: make Dispatch receive Update by value

## v1.1.5

- handlers: add ConversationID to PersistenceKey
- persistence: add FilePersistence
- examples: add info about FilePersistence

## v1.1.4

- global: add .travis.yml & Makefile
- examples: change directory structure
- handlers: allow nullable filters (defaults to Any())
- filters: rename content filters to have "Has" prefix

## v1.1.3

- examples: add members example

## v1.1.2

- filters: fix bug with IsNewChatMembers

## v1.1.1

- filters: add IsNewChatMembers & IsLeftChatMember

## v1.1.0


## v1.0.5

- examples: switch to tm.Update
- filters: use methods from tm.Update
- handlers: use methods from tm.Update
- mux: wrap tgbotapi.Update in tm.Update
- helpers: delete module
- update: add Update struct with Effective(User|Chat|Message) methods

## v1.0.4


## v1.0.3

- filters: improve IsCommand & IsAnyCommand filters
- readme: add info about filtering edited messages by mistakes

## v1.0.2

- readme: add "tips and common pitfalls" section
- examples: add cat_callback example
- handlers: use GetEffectiveUser & GetEffectiveChat
- filters: add a lot of new filters
- helpers: add GetEffectiveUser, GetEffectiveChat & GetEffectiveMessage

## v1.0.1


## v1.0.0


