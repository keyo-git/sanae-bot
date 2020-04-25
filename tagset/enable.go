package tagset

import (
	"fmt"

	"github.com/keyo-git/sanae-bot/bot"
	"github.com/keyo-git/sanae-bot/db"
)

func init() {
	bot.GlobalRegistry.RegisterCmd(enable, "!enable", "Enable tag")
	bot.GlobalRegistry.RegisterCmd(disable, "!disable", "Disable tag")
}

func enable(s *bot.Sanae, argv []string, trigger bot.CmdTrigger) error {
	if len(argv) == 0 {
		_, err := s.Sess().ChannelMessageSend(
			trigger.ChannelID,
			"Incorrect number of arguments",
		)
		return err
	}

	var rawmsg string

	r, err := db.EnableTag(s.DbHandle(), trigger.Author.ID, argv[0])
	if err != nil {
		return err
	}
	if n, _ := r.RowsAffected(); n == 0 {
		rawmsg = fmt.Sprintf("Tag `%s` does not exist :c", argv[0])
	} else {
		rawmsg = fmt.Sprintf("Tag `%s` has been enabled", argv[0])
	}

	_, err = s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
	return err
}

func disable(s *bot.Sanae, argv []string, trigger bot.CmdTrigger) error {
	if len(argv) == 0 {
		_, err := s.Sess().ChannelMessageSend(
			trigger.ChannelID,
			"Incorrect number of arguments",
		)
		return err
	}

	var rawmsg string

	r, err := db.DisableTag(s.DbHandle(), trigger.Author.ID, argv[0])
	if err != nil {
		return err
	}
	if n, _ := r.RowsAffected(); n == 0 {
		rawmsg = fmt.Sprintf("Tag `%s` does not exist :c", argv[0])
	} else {
		rawmsg = fmt.Sprintf("Tag `%s` has been disabled", argv[0])
	}

	_, err = s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
	return err
}
