package tagset

import (
	"fmt"
	"strings"

	"github.com/keyo-git/sanae-bot/bot"
	"github.com/keyo-git/sanae-bot/db"
)

func init() {
	searchStr := `Remove tag(s) from watched/hidden list`
	bot.GlobalRegistry.RegisterCmd(remove, "!remove", searchStr)
}

func remove(s *bot.Sanae, argv []string, trigger bot.CmdTrigger) error {
	if len(argv) == 0 {
		_, err := s.Sess().ChannelMessageSend(
			trigger.ChannelID,
			"Incorrect number of arguments",
		)
		return err
	}

	var rawmsg string
	var removedTags []string

	for _, tag := range argv {
		tag = strings.ReplaceAll(tag, "_", " ")
		r, err := db.DeleteTag(s.DbHandle(), trigger.Author.ID, tag)
		if err != nil {
			return err
		}
		if n, _ := r.RowsAffected(); n != 0 {
			removedTags = append(removedTags, tag)
		}
	}

	if len(removedTags) == 0 {
		// FIXME: same output if tagset exists but no tags have been matched
		rawmsg = fmt.Sprintf("No tags have been matched")
	} else if len(removedTags) > 0 {
		rawmsg = "Tags "
		for _, removedTag := range removedTags {
			rawmsg += fmt.Sprintf("`%s` ", removedTag)
		}
		rawmsg += fmt.Sprintf(" have been removed")
	}

	_, err := s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
	return err
}
