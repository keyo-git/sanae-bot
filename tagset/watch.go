package tagset

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/keyo-git/sanae-bot/bot"
	"github.com/keyo-git/sanae-bot/db"
)

func init() {
	bot.GlobalRegistry.RegisterCmd(watch, "!watch", "Add tag(s) to watch list")
	bot.GlobalRegistry.RegisterCmd(hide, "!hide", "Add tag(s) to hidden list")
}

var maxTags = 10

func addTags(
	handle *sql.DB,
	argv []string,
	trigger bot.CmdTrigger,
	watched bool,
) (addedTags []string, err error) {
	for _, tag := range argv {
		tag = strings.ReplaceAll(tag, "_", " ")
		tag = strings.ReplaceAll(tag, "-", " ")

		if e := db.FindTag(handle, trigger.Author.ID, tag); e != nil {
			if e.Watched != watched {
				err := db.SetWatched(handle, trigger.Author.ID, tag, watched)
				if err != nil {
					return nil, err
				}
			} else {
				continue
			}
		} else {
			err := db.InsertTag(handle, trigger.Author.ID, tag, watched)
			if err != nil {
				return nil, err
			}
		}

		addedTags = append(addedTags, tag)
	}

	return
}

func watch(s *bot.Sanae, argv []string, trigger bot.CmdTrigger) error {
	if len(argv) == 0 {
		_, err := s.Sess().ChannelMessageSend(
			trigger.ChannelID,
			"Incorrect number of arguments",
		)
		return err
	}

	var rawmsg string

	ts, err := db.FindAllTags(s.DbHandle(), trigger.Author.ID)
	if err != nil {
		return err
	}

	n := maxTags - (len(ts)) + 1
	if len(argv) < n {
		n = len(argv)
	}
	if n == 0 {
		rawmsg = fmt.Sprintf("Maximum allowed number of tags you can add to tagset is %d", maxTags)
		_, err := s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
		return err
	}
	argv = argv[:n]

	addedTags, err := addTags(s.DbHandle(), argv, trigger, true)
	if err != nil {
		return err
	}
	rawmsg = "Tags "
	for _, tag := range addedTags {
		rawmsg += fmt.Sprintf("`%s` ", tag)
	}
	rawmsg += fmt.Sprintf("have been added to watched list")

	_, err = s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
	return err
}

func hide(s *bot.Sanae, argv []string, trigger bot.CmdTrigger) error {
	if len(argv) == 0 {
		_, err := s.Sess().ChannelMessageSend(
			trigger.ChannelID,
			"Incorrect number of arguments",
		)
		return err
	}

	var rawmsg string

	ts, err := db.FindAllTags(s.DbHandle(), trigger.Author.ID)
	if err != nil {
		return err
	}

	n := maxTags - (len(ts)) + 1
	if len(argv) < n {
		n = len(argv)
	}
	if n == 0 {
		rawmsg = fmt.Sprintf("Maximum allowed number of tags you can add to tagset is %d", maxTags)
		_, err := s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
		return err
	}
	argv = argv[:n]

	addedTags, err := addTags(s.DbHandle(), argv, trigger, false)
	if err != nil {
		return err
	}
	rawmsg = "Tags "
	for _, tag := range addedTags {
		rawmsg += fmt.Sprintf("`%s` ", tag)
	}
	rawmsg += fmt.Sprintf("have been added to hidden list")

	_, err = s.Sess().ChannelMessageSend(trigger.ChannelID, rawmsg)
	return err
}
