// unknown_group is a bot that will say hello/goodbye when members enver/leave group. He will also leave any unknown groups he is invited to. Provide known groups as env var, e. g. KNOWN_GROUPS=-111,-222,333,444
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tm "github.com/and3rson/telemux"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type KnownGroups struct {
	ids []int64
}

func (k KnownGroups) LoadFromEnv() {
	for _, idStr := range strings.Split(os.Getenv("KNOWN_GROUPS"), ",") {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			panic(err)
		}
		k.ids = append(k.ids, id)
	}
}

func (k KnownGroups) IsKnownGroup(id int64) bool {
	for _, otherID := range k.ids {
		if id == otherID {
			return true
		}
	}
	return false
}

var knownGroups KnownGroups

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	knownGroups.LoadFromEnv()

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	mux := tm.NewMux().
		AddHandler(tm.NewHandler(
			tm.IsNewChatMembers(),
			func(u *tm.Update) {
				chat := u.EffectiveChat()
				// Check every new member
				for _, user := range *u.Message.NewChatMembers {
					if user.ID == bot.Self.ID {
						// This is us!
						if !knownGroups.IsKnownGroup(chat.ID) {
							// Group is unknown, leave chat
							bot.LeaveChat(tgbotapi.ChatConfig{ChatID: chat.ID, SuperGroupUsername: ""})
						}
					} else {
						// Greet new member
						bot.Send(tgbotapi.NewMessage(chat.ID, fmt.Sprintf("Hello, %s!", user.UserName)))
					}
				}
			},
		)).
		AddHandler(tm.NewHandler(
			tm.IsLeftChatMember(),
			func(u *tm.Update) {
				// Say goodbye to a member who has just left
				bot.Send(tgbotapi.NewMessage(
					u.Message.Chat.ID,
					fmt.Sprintf("Goodbye, %s!", u.Message.LeftChatMember.UserName),
				))
			},
		))
	for update := range updates {
		mux.Dispatch(&update)
	}
}
