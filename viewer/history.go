package viewer

import (
	"database/sql"

	"github.com/keyo-git/sanae-bot/api"
	"github.com/keyo-git/sanae-bot/bot"
)

func init() {
	searchStr := `Show all galleries posted in this channel`
	bot.GlobalRegistry.RegisterCmd(history, "!history", searchStr)
	//bot.CmdRegistry = append(bot.CmdRegistry, bot.CmdEntry{"!whistory", whistory, "aaa"})
}

type historyFetcher struct {
	handle *sql.DB
	exAPI  *api.ExHentaiAPI

	fetchQuery string
	idxQuery   string
	id         string
}

func newChannelHistoryFetcher(
	handle *sql.DB,
	exAPI *api.ExHentaiAPI,
	channelID string,
) *historyFetcher {
	fetchQuery := `SELECT url FROM gallery_history
                 WHERE channel_id=$1
                 ORDER BY id DESC LIMIT $2 OFFSET $3`
	idxQuery := `SELECT COUNT(*) FROM gallery_history WHERE channel_id=$1`

	return &historyFetcher{
		handle:     handle,
		exAPI:      exAPI,
		fetchQuery: fetchQuery,
		idxQuery:   idxQuery,
		id:         channelID,
	}
}

/*func NewWatchedHistoryFetcher(
	handle *sql.DB,
	exAPI *api.ExHentaiAPI,
	userID string,
	pageLen int,
) *HistoryFetcher {
	fetchQuery := `SELECT url FROM gallery_history gh
                 INNER JOIN matched_tags mt
								 ON gh.id = mt.gallery_id
								 WHERE mt.user_id=$1
								 ORDER BY mt.gallery_id DESC LIMIT $2 OFFSET $3`
	idxQuery := `SELECT COUNT(*) FROM gallery_history gh
							 INNER JOIN matched_tags mt
							 ON gh.id = mt.gallery_id
							 WHERE mt.user_id=$1`

	return &HistoryFetcher{
		handle:     handle,
		exAPI:      exAPI,
		fetchQuery: fetchQuery,
		idxQuery:   idxQuery,
		pageLen:    pageLen,
	}
}*/

func (f *historyFetcher) fetch(
	page int,
) (gms []*api.GalleryMetadata, err error) {

	rows, err := f.handle.Query(f.fetchQuery, f.id, pageLen, page*pageLen)
	if err != nil {
		return
	}

	for rows.Next() {
		var rawurl string
		err = rows.Scan(&rawurl)
		if err != nil {
			return nil, err
		}

		gm, err := f.exAPI.RequestGalleryMetadata(rawurl)
		if err != nil {
			return nil, err
		}
		gms = append(gms, gm)
	}

	return
}

func (f *historyFetcher) lastIdx() (lastIdx int, err error) {
	err = f.handle.QueryRow(f.idxQuery, f.id).Scan(&lastIdx)
	lastIdx--
	return
}

func history(s *bot.Sanae, _ []string, trigger bot.CmdTrigger) error {
	fetcher := newChannelHistoryFetcher(
		s.DbHandle(),
		s.ExAPI(),
		trigger.ChannelID,
	)

	return createGalleriesViewer(s.Sess(), trigger, fetcher)
}

/*func whistory(
	s *bot.Sanae,
	_ []string,
	trigger *discordgo.MessageCreate,
) error {
	fetcher := NewWatchedHistoryFetcher(
		s.DbHandle(),
		s.ExAPI(),
		trigger.Author.ID,
		25,
	)

	return createGalleriesViewer(s, trigger, 25, fetcher)
}*/
