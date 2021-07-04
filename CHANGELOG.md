# Changelog

## v1.8.2

- `readme`: **add** info on update type handlers
- `handlers`: **add** missing filter arg to NewCommandHandler

## v1.8.1

- `persistence`: **fix** PersistenceKey error reporting during unmarshal error
- `global`: **add** tests, increase coverage to 98.3%

## v1.8.0

- .travis.yml: add coverage collection for codeclimate
- `readme`: **add** coverage badge
- `filters`: **update** IsCommandMessage & IsAnyCommandMessage to use regex
- `mux`: **make** Processor public for testing convenience
- `global`: **add** unit tests for the most important functionality

## v1.7.0

- `readme`: **update** docs on recovering
- `examples`: **update** recover usage
- `mux`: **rename** Recoverer/SetRecoverer to Recover/SetRecover
- `update`: **add** Context to Update
- `mux`: **add** processor interface
- `mux`: **make** Mux implement processor
- `handlers`: **make** Handler implement processor

## v1.6.0

- `readme`: **drop** async handlers section
- `readme`: **drop** section about update consumption from filters
- `readme`: **add** section about update consumption from handle functions
- `readme`: **replace** "list" with "slice"
- `readme`: **update** examples
- `handlers`: **drop** NewAsyncHandler
- `handlers`: **add** HandleFunc
- `handlers`: **do** not check if u.Consumed by filters in TransitionHandler
- `filters`: **rename** Filter to FilterFunc
- `examples`: **drop** async handler usage from cat_callback
- `mux`: **do** not check if u.Consumed by filters in Mux

## v1.5.3

- `mux`: **fix** varargs for AddHandler & AddMux

## v1.5.2

- `mux`: **allow** varargs for AddHandler & AddMux

## v1.5.1

- `mux`: **add** nested mux support via AddMux
- `mux`: **add** global filters via SetGlobalFilter
- `mux`: **prioritize** update consumption over filter return result
- `handlers`: **prioritize** update consumption over filter return result
- `examples`: **add** nested_mux example

## v1.5.0

- `readme`: **update** docs on conversations
- `examples`: **update** conversation example
- `gormpersistence`: **update** readme
- `handlers`: **set** PersistenceContext during conversation filters/handlers
- `mux`: **set** PersistenceContext to nil by default
- `update`: **add** PersistenceContext to Update
- `persistence`: **update** wording
- `changelog`: **update** mkchangelog.sh and CHANGELOG.md
- `makefile`: **update** TAG variable

## v1.4.5

- `mux`: **pass** stack trace to recoverer
- `readme`: **update** docs
- `examples`: **update** recoverer implementation

## v1.4.4

- `readme`: **fix** errors
- `filters`: **add** todo

## v1.4.3

- `readme`: **add** section on update consumption
- `readme`: **update** TOC
- `mux`: **add** bot arg to Dispatch
- `mux`: **check** if message is consumed
- `filters`: **update** IsCommandMessage filter to check against bot username
- `examples`: **update** Dispatch invocation
- `gormpersistence`: **update** readme
- `update`: **add** Bot & Consumed to Update struct
- `handlers`: **check** if message is consumed

## v1.4.2

- `filters`: **add** /foo@bar format support for IsCommandMessage
- `mux`: **add** todo
- `update`: **add** handlig of missing cases for Update.EffectiveUser

## v1.4.1

- `filters`: **fix** IsAnyCommand improperly matching command prefixes
- `examples`: **fix** spelling
- `gormpersistence`: **add** unit tests
- `gormpersistence`: **change** tabs to spaces in readme
- `global`: **update** makefiles
- `gormpersistence`: **update** readme

## v1.4.0

- `examples`: **remove** redundant os.Exit(1) calls
- `readme`: **add** info about GORMPersistence & update info on default state
- `gormpersistence`: **add** docstrings
- `gormpersistence`: **add** readme
- `persistence`: **add** GORM tags to PersistenceKey
- `gormpersistence`: **models**: add ConversationState & ConversationData models
- `gormpersistence`: **persistence**: add GORMPersistence implementation
- `readme`: **add** TOC
- `global`: **add** mkchangelog.sh & CHANGELOG.md

## v1.3.0

- `readme`: **update** docs on persistence & recoverer
- `handlers`: **use** persistence context
- `mux`: **add** recoverer
- `persistence`: **add** persistence context
- `examples`: **add** error_handling

## v1.2.0

- `readme`: **cleanup**
- `examples`: **update** Dispatch call
- `update`: **embed** Update instead of *Update
- `mux`: **make** Dispatch receive Update by value

## v1.1.5

- `handlers`: **add** ConversationID to PersistenceKey
- `persistence`: **add** FilePersistence
- `examples`: **add** info about FilePersistence

## v1.1.4

- `global`: **add** .travis.yml & Makefile
- `examples`: **change** directory structure
- `handlers`: **allow** nullable filters (defaults to Any())
- `filters`: **rename** content filters to have "Has" prefix

## v1.1.3

- `examples`: **add** members example

## v1.1.2

- `filters`: **fix** bug with IsNewChatMembers

## v1.1.1

- `filters`: **add** IsNewChatMembers & IsLeftChatMember

## v1.1.0


## v1.0.5

- `examples`: **switch** to tm.Update
- `filters`: **use** methods from tm.Update
- `handlers`: **use** methods from tm.Update
- `mux`: **wrap** tgbotapi.Update in tm.Update
- `helpers`: **delete** module
- `update`: **add** Update struct with Effective(User|Chat|Message) methods

## v1.0.4


## v1.0.3

- `filters`: **improve** IsCommand & IsAnyCommand filters
- `readme`: **add** info about filtering edited messages by mistakes

## v1.0.2

- `readme`: **add** "tips and common pitfalls" section
- `examples`: **add** cat_callback example
- `handlers`: **use** GetEffectiveUser & GetEffectiveChat
- `filters`: **add** a lot of new filters
- `helpers`: **add** GetEffectiveUser, GetEffectiveChat & GetEffectiveMessage

## v1.0.1


## v1.0.0


