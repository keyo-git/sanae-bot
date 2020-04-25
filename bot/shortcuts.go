package bot

import "github.com/bwmarrin/discordgo"

// CmdTrigger ...
type CmdTrigger = *discordgo.MessageCreate

// ReactTrigger ...
type ReactTrigger = *discordgo.MessageReactionAdd

// TagMap ...
type TagMap = map[string][]string
