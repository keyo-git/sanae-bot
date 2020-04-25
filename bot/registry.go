package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// GlobalRegistry ...
var GlobalRegistry registry

type cmd = func(*Sanae, []string, *discordgo.MessageCreate) error
type emojiCmd = func(*Sanae, *discordgo.MessageReactionAdd) error

type cmdEntry struct {
	Cmd     cmd
	Trigger string
	Help    string
}

type reactionEntry struct {
	Cmd   emojiCmd
	Emoji string
}

type registry struct {
	cmdRegistry      []cmdEntry
	reactionRegistry []reactionEntry
}

func (r *registry) RegisterCmd(cmd cmd, name string, help string) {
	for _, c := range r.cmdRegistry {
		if c.Trigger == name {
			log.Printf("Command with name %s already registred\n", name)
			return
		}
	}

	r.cmdRegistry = append(r.cmdRegistry, cmdEntry{cmd, name, help})
	log.Printf("Registered command %s\n", name)
}

func (r *registry) RegisterReaction(cmd emojiCmd, emoji string) {
	for _, c := range r.reactionRegistry {
		if c.Emoji == emoji {
			log.Printf("Emoji %s already registred\n", emoji)
			return
		}
	}

	r.reactionRegistry = append(r.reactionRegistry, reactionEntry{cmd, emoji})
	log.Printf("Registered emoji %s\n", emoji)
}
