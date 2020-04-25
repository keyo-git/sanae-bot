package viewer

import (
	"github.com/bwmarrin/discordgo"
	"github.com/keyo-git/sanae-bot/api"
	"github.com/keyo-git/sanae-bot/bot"
)

var viewersRegistry []viewer

const (
	pageLen    = 25
	maxViewers = 20

	emojiLeft  = "🅰"
	emojiRight = "🇧"
)

func init() {
	bot.GlobalRegistry.RegisterReaction(flipLeft, emojiLeft)
	bot.GlobalRegistry.RegisterReaction(flipRight, emojiRight)
}

type fetcher interface {
	fetch(int) ([]*api.GalleryMetadata, error)
	lastIdx() (int, error)
}

//type Fetcher = func(map[string]string, int, int) ([]*api.GalleryMetadata, error)

type viewer struct {
	channelID string
	messageID string

	idx  int
	page int

	lastIdx       int
	fetcherParams map[string]string
	fetcher       fetcher

	gms []*api.GalleryMetadata
}

func (v *viewer) MessageID() string {
	return v.messageID
}

func (v *viewer) flipDecorator(
	sess *discordgo.Session,
	trigger bot.ReactTrigger,
	f func(*viewer) error,
) (err error) {
	if err = sess.MessageReactionRemove(
		trigger.ChannelID,
		trigger.MessageID,
		trigger.Emoji.Name,
		trigger.UserID,
	); err != nil {
		return
	}

	idx := v.idx
	if err = f(v); err != nil || idx == v.idx {
		return err
	}

	gm := v.gms[v.idx%pageLen]
	err = v.update(sess, gm)

	return
}

func (v *viewer) update(
	sess *discordgo.Session,
	gm *api.GalleryMetadata,
) (err error) {
	embed := bot.EmbedGallery(gm)
	msgEdit := &discordgo.MessageEdit{
		Channel: v.channelID,
		ID:      v.messageID,
		Embed:   embed,
	}
	_, err = sess.ChannelMessageEditComplex(msgEdit)

	return
}

func decIdx(v *viewer) (err error) {
	if v.idx == 0 {
		return
	}
	if (v.idx-1)/pageLen != v.page {
		v.page--
		if v.gms, err = v.fetcher.fetch(v.page); err != nil {
			return
		}
	}
	v.idx--

	return
}

func incIdx(v *viewer) (err error) {
	if v.idx == v.lastIdx {
		return
	}
	if (v.idx+1)/pageLen != v.page {
		v.page++
		if v.gms, err = v.fetcher.fetch(v.page); err != nil {
			return
		}
	}
	v.idx++

	return
}

func createGalleriesViewer(
	sess *discordgo.Session,
	trigger bot.CmdTrigger,
	fetcher fetcher,
) error {
	gms, err := fetcher.fetch(0)
	if err != nil {
		return err
	}

	if len(gms) == 0 {
		_, err := sess.ChannelMessageSend(
			trigger.ChannelID,
			"No galleris found :c")
		return err
	}

	gm := gms[0]

	embed := bot.EmbedGallery(gm)
	msgSend := &discordgo.MessageSend{
		Embed: embed,
	}

	response, err := sess.ChannelMessageSendComplex(
		trigger.ChannelID,
		msgSend,
	)
	if err != nil {
		return err
	}

	channelID := response.ChannelID
	responseID := response.ID

	err = sess.MessageReactionAdd(channelID, responseID, emojiLeft)
	if err != nil {
		return err
	}
	err = sess.MessageReactionAdd(channelID, responseID, emojiRight)
	if err != nil {
		return err
	}

	lastIdx, err := fetcher.lastIdx()

	v := viewer{
		channelID: channelID,
		messageID: responseID,
		lastIdx:   lastIdx,
		fetcher:   fetcher,
		gms:       gms,
	}

	if len(viewersRegistry) > maxViewers {
		viewersRegistry = append(viewersRegistry[1:], v)
	} else {
		viewersRegistry = append(viewersRegistry, v)
	}

	return nil
}

func flipLeft(s *bot.Sanae, r bot.ReactTrigger) (err error) {
	for i, v := range viewersRegistry {
		if v.messageID == r.MessageID {
			err = viewersRegistry[i].flipDecorator(s.Sess(), r, decIdx)
			break
		}
	}

	return
}

func flipRight(s *bot.Sanae, r bot.ReactTrigger) (err error) {
	for i, v := range viewersRegistry {
		if v.messageID == r.MessageID {
			err = viewersRegistry[i].flipDecorator(s.Sess(), r, incIdx)
			break
		}
	}

	return
}
