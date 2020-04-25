package bot

import (
	"log"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/keyo-git/sanae-bot/api"
	"github.com/keyo-git/sanae-bot/db"
)

func isExLink(rawurl string) bool {
	galleryURL, err := url.Parse(rawurl)
	if err != nil {
		return false
	}

	host := galleryURL.Hostname()
	return host == "exhentai.org" || host == "e-hentai.org"
}

func getUsersToNotify(
	s *Sanae,
	gm *api.GalleryMetadata,
	guildID string,
) (usersToNotify []*discordgo.User, matchedTags TagMap) {
	matchedTags = make(TagMap)

	users, err := db.FindUsers(s.DbHandle())
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	for _, userID := range users {
		// Skip users not in server
		if mem, _ := s.Sess().GuildMember(guildID, userID); mem == nil {
			continue
		}

		ts, err := db.FindAllTags(s.DbHandle(), userID)
		if err != nil {
			return nil, nil
		}

		for _, e := range ts {
			if !e.Enabled {
				continue
			}

			n := -1
			if noCategory(e.Tag) {
				for i := range gm.Tags {
					if strings.Contains(gm.Tags[i], e.Tag) {
						n = i
						break
					}
				}
			} else {
				for i := range gm.Tags {
					if gm.Tags[i] == e.Tag {
						n = i
						break
					}
				}
			}
			if n != -1 {
				if e.Watched {
					matchedTags[userID] = append(matchedTags[userID], e.Tag)
				} else {
					delete(matchedTags, userID)
					break
				}
			}
		}

		if _, ok := matchedTags[userID]; ok {
			user, err := s.Sess().User(userID)
			if err != nil {
				return nil, nil
			}
			usersToNotify = append(usersToNotify, user)
			break
		}
	}

	return
}

func ex(s *Sanae, argv []string, trigger *discordgo.MessageCreate) error {
	for _, rawurl := range argv {
		if !isExLink(rawurl) {
			continue
		}

		gm, err := s.ExAPI().RequestGalleryMetadata(rawurl)
		if err != nil {
			log.Println(err)
			continue
		}

		var notifications string
		usersToNotify, matchedTags := getUsersToNotify(s, gm, trigger.GuildID)
		for _, user := range usersToNotify {
			notifications += user.Mention()
		}

		embed := EmbedGallery(gm)

		msgSend := &discordgo.MessageSend{
			Content: notifications,
			Embed:   embed,
		}

		galleryID, err := db.InsertGallery(s.DbHandle(), trigger.Message, rawurl)
		if err != nil {
			return err
		}

		err = db.InsertMatchedTags(s.DbHandle(), matchedTags, galleryID)
		if err != nil {
			return err
		}

		_, err = s.Sess().ChannelMessageSendComplex(trigger.ChannelID, msgSend)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
