package viewer

import (
	"strings"

	"github.com/keyo-git/sanae-bot/api"
	"github.com/keyo-git/sanae-bot/bot"
)

func init() {
	helpStr := `Search galleries on ExHentai`
	bot.GlobalRegistry.RegisterCmd(search, "!search", helpStr)
}

type searchFetcher struct {
	exAPI *api.ExHentaiAPI

	query string
}

func (f *searchFetcher) fetch(
	page int,
) (gms []*api.GalleryMetadata, err error) {
	return f.exAPI.SearchGalleries(f.query, page)
}

func (f *searchFetcher) lastIdx() (lastIdx int, err error) {
	lastIdx, err = f.exAPI.NumSearchGalleries(f.query)
	lastIdx--
	return
}

func search(s *bot.Sanae, argv []string, trigger bot.CmdTrigger) error {
	searchQuery := strings.Join(argv, " ")
	fetcher := &searchFetcher{s.ExAPI(), searchQuery}

	return createGalleriesViewer(s.Sess(), trigger, fetcher)
}
