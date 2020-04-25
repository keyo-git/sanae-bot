package db

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
)

// InsertGallery inserts gallery into history table
func InsertGallery(db *sql.DB,
	msg *discordgo.Message,
	rawurl string,
) (galleryID int, err error) {
	channelID := msg.ChannelID
	messageID := msg.ID
	date, _ := msg.Timestamp.Parse()

	query := `INSERT INTO gallery_history (channel_id, message_id, date, url)
						VALUES ($1, $2, $3, $4)
						RETURNING id`

	err = db.QueryRow(
		query,
		channelID,
		messageID,
		date,
		rawurl,
	).Scan(&galleryID)

	return
}

// InsertMatchedTags inserts matched tags into matched tags table
func InsertMatchedTags(
	db *sql.DB,
	matchedTags map[string][]string,
	galleryID int,
) error {
	query := `INSERT INTO matched_tags (user_id, tags, gallery_id)
						VALUES ($1, $2, $3)`
	for userID, tags := range matchedTags {
		_, err := db.Exec(query, userID, pq.Array(tags), galleryID)
		if err != nil {
			return err
		}
	}

	return nil
}
