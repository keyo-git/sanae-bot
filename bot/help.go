package bot

import (
	"fmt"
)

func init() {
	GlobalRegistry.RegisterCmd(help, "!sanae", "Print this message")
}

func pretty(cs []cmdEntry, level int) (prettyMsg string) {
	var pad string
	for i := 0; i < level; i++ {
		pad += "\t"
	}

	var maxLen int
	for _, c := range cs {
		if len(c.Trigger) > maxLen {
			maxLen = len(c.Trigger)
		}
	}

	for _, c := range cs {
		var blanks string
		for i := 0; i < maxLen-len(c.Trigger); i++ {
			blanks += " "
		}
		prettyMsg += fmt.Sprintf("%s%s%s\t%s\n", pad, c.Trigger, blanks, c.Help)
	}

	return
}

func help(s *Sanae, _ []string, trigger CmdTrigger) (err error) {
	helpMsg := "```"
	helpMsg += pretty(GlobalRegistry.cmdRegistry, 0)
	helpMsg += "```"
	_, err = s.dg.ChannelMessageSend(trigger.ChannelID, helpMsg)
	return
}
