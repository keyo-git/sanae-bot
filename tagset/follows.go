package tagset

import (
	"fmt"

	"github.com/keyo-git/sanae-bot/bot"
	"github.com/keyo-git/sanae-bot/db"
)

func init() {
	searchStr := `Print watched/hidden list`
	bot.GlobalRegistry.RegisterCmd(myfollows, "!myfollows", searchStr)
}

func myfollows(s *bot.Sanae, argv []string, trigger bot.CmdTrigger) error {
	var rawmsg string

	tagset, err := db.FindAllTags(s.DbHandle(), trigger.Author.ID)
	if err != nil {
		return err
	}

	for _, e := range tagset {
		var status string
		var watched string
		if e.Watched {
			watched = "watched"
		} else {
			watched = "hidden"
		}
		if e.Enabled {
			status = "✔️"
		} else {
			status = "✖"
		}
		rawmsg += fmt.Sprintf("`%s` `%s`%s\n", e.Tag, watched, status)
	}
	if rawmsg == "" {
		rawmsg = "You don't follow any tags"
	}

	_, err = s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
	return err
}
