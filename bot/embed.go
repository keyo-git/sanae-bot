package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/keyo-git/sanae-bot/api"
)

var categoryColors = map[string]int{
	"Doujinshi":  0x9E2720,
	"Manga":      0xDB6C24,
	"Artist CG":  0xD38F1D,
	"Game CG":    0x6A936D,
	"Western":    0xAB9F60,
	"Non-H":      0x5FA9CF,
	"Image Set":  0x325CA2,
	"Cosplay":    0x6A32A2,
	"Asian Porn": 0xA23282,
	"Misc":       0x777777,
}

// EmbedGallery creates an embeded Discord message from gallery metadata
func EmbedGallery(gm *api.GalleryMetadata) (embed *discordgo.MessageEmbed) {
	tagMap := groupTags(gm.Tags)
	galleryDescription := fmt.Sprintf(
		`%s
			**Rating: ** %.2f
			**Uploaded: ** %+v
			**Uploader: **  %s
			**Length: ** %d pages

			`,
		gm.TitleJPN,
		gm.Rating,
		time.Unix(gm.Posted, 0),
		gm.Uploader,
		gm.FileCount,
	)
	galleryDescription += formatTags(tagMap)

	embedImage := &discordgo.MessageEmbedImage{URL: gm.Thumb}
	embed = &discordgo.MessageEmbed{
		URL:         api.BuildURLFromMetadata(gm),
		Title:       gm.Title,
		Description: galleryDescription,
		Image:       embedImage,
		Color:       categoryColors[gm.Category],
	}

	return
}
