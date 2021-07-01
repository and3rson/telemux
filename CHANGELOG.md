# Changelog

## v1.6.0

- `readme`: _drop_ async handlers section
- `readme`: _drop_ section about update consumption from filters
- `readme`: _add_ section about update consumption from handle functions
- `readme`: _replace_ "list" with "slice"
- `readme`: _update_ examples
- `handlers`: _drop_ NewAsyncHandler
- `handlers`: _add_ HandleFunc
- `handlers`: _do_ not check if u.Consumed by filters in TransitionHandler
- `filters`: _rename_ Filter to FilterFunc
- `examples`: _drop_ async handler usage from cat_callback
- `mux`: _do_ not check if u.Consumed by filters in Mux

## v1.5.3

- `mux`: _fix_ varargs for AddHandler & AddMux

## v1.5.2

- `mux`: _allow_ varargs for AddHandler & AddMux

## v1.5.1

- `mux`: _add_ nested mux support via AddMux
- `mux`: _add_ global filters via SetGlobalFilter
- `mux`: _prioritize_ update consumption over filter return result
- `handlers`: _prioritize_ update consumption over filter return result
- `examples`: _add_ nested_mux example

## v1.5.0

- `readme`: _update_ docs on conversations
- `examples`: _update_ conversation example
- `gormpersistence`: _update_ readme
- `handlers`: _set_ PersistenceContext during conversation filters/handlers
- `mux`: _set_ PersistenceContext to nil by default
- `update`: _add_ PersistenceContext to Update
- `persistence`: _update_ wording
- `changelog`: _update_ mkchangelog.sh and CHANGELOG.md
- `makefile`: _update_ TAG variable

## v1.4.5

- `mux`: _pass_ stack trace to recoverer
- `readme`: _update_ docs
- `examples`: _update_ recoverer implementation

## v1.4.4

- `readme`: _fix_ errors
- `filters`: _add_ todo

## v1.4.3

- `readme`: _add_ section on update consumption
- `readme`: _update_ TOC
- `mux`: _add_ bot arg to Dispatch
- `mux`: _check_ if message is consumed
- `filters`: _update_ IsCommandMessage filter to check against bot username
- `examples`: _update_ Dispatch invocation
- `gormpersistence`: _update_ readme
- `update`: _add_ Bot & Consumed to Update struct
- `handlers`: _check_ if message is consumed

## v1.4.2

- `filters`: _add_ /foo@bar format support for IsCommandMessage
- `mux`: _add_ todo
- `update`: _add_ handlig of missing cases for Update.EffectiveUser

## v1.4.1

- `filters`: _fix_ IsAnyCommand improperly matching command prefixes
- `examples`: _fix_ spelling
- `gormpersistence`: _add_ unit tests
- `gormpersistence`: _change_ tabs to spaces in readme
- `global`: _update_ makefiles
- `gormpersistence`: _update_ readme

## v1.4.0

- `examples`: _remove_ redundant os.Exit(1) calls
- `readme`: _add_ info about GORMPersistence & update info on default state
- `gormpersistence`: _add_ docstrings
- `gormpersistence`: _add_ readme
- `persistence`: _add_ GORM tags to PersistenceKey
- `gormpersistence`: _models_: add ConversationState & ConversationData models
- `gormpersistence`: _persistence_: add GORMPersistence implementation
- `readme`: _add_ TOC
- `global`: _add_ mkchangelog.sh & CHANGELOG.md

## v1.3.0

- `readme`: _update_ docs on persistence & recoverer
- `handlers`: _use_ persistence context
- `mux`: _add_ recoverer
- `persistence`: _add_ persistence context
- `examples`: _add_ error_handling

## v1.2.0

- `readme`: _cleanup_
- `examples`: _update_ Dispatch call
- `update`: _embed_ Update instead of *Update
- `mux`: _make_ Dispatch receive Update by value

## v1.1.5

- `handlers`: _add_ ConversationID to PersistenceKey
- `persistence`: _add_ FilePersistence
- `examples`: _add_ info about FilePersistence

## v1.1.4

- `global`: _add_ .travis.yml & Makefile
- `examples`: _change_ directory structure
- `handlers`: _allow_ nullable filters (defaults to Any())
- `filters`: _rename_ content filters to have "Has" prefix

## v1.1.3

- `examples`: _add_ members example

## v1.1.2

- `filters`: _fix_ bug with IsNewChatMembers

## v1.1.1

- `filters`: _add_ IsNewChatMembers & IsLeftChatMember

## v1.1.0


## v1.0.5

- `examples`: _switch_ to tm.Update
- `filters`: _use_ methods from tm.Update
- `handlers`: _use_ methods from tm.Update
- `mux`: _wrap_ tgbotapi.Update in tm.Update
- `helpers`: _delete_ module
- `update`: _add_ Update struct with Effective(User|Chat|Message) methods

## v1.0.4


## v1.0.3

- `filters`: _improve_ IsCommand & IsAnyCommand filters
- `readme`: _add_ info about filtering edited messages by mistakes

## v1.0.2

- `readme`: _add_ "tips and common pitfalls" section
- `examples`: _add_ cat_callback example
- `handlers`: _use_ GetEffectiveUser & GetEffectiveChat
- `filters`: _add_ a lot of new filters
- `helpers`: _add_ GetEffectiveUser, GetEffectiveChat & GetEffectiveMessage

## v1.0.1


## v1.0.0


