package db

import "database/sql"

// TagsetEntry ...
type TagsetEntry struct {
	Tag     string
	Watched bool
	Enabled bool
}

// Tagset ...
type Tagset []*TagsetEntry

// FindAllTags retrieves all tags for given user
func FindAllTags(db *sql.DB, userID string) (ts Tagset, err error) {
	query := `SELECT tag, watched, enabled FROM tags WHERE user_id=$1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return
	}

	for rows.Next() {
		t := &TagsetEntry{}
		err = rows.Scan(&t.Tag, &t.Watched, &t.Enabled)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}

	return
}

// FindTag retrieves tag entry from database given user and tag
func FindTag(db *sql.DB, userID, tag string) *TagsetEntry {
	t := &TagsetEntry{}
	query := `SELECT tag, watched, enabled
						FROM tags WHERE user_id=$1 and tag=$2`
	row := db.QueryRow(query, userID, tag)
	if err := row.Scan(&t.Tag, &t.Watched, &t.Enabled); err != nil {
		return nil
	}
	return t
}

// InsertTag inserts new tag entry into database
func InsertTag(db *sql.DB, userID, tag string, watched bool) (err error) {
	query := `INSERT INTO tags (user_id, tag, watched) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, userID, tag, watched)
	return
}

// DeleteTag removes tag entry from database
func DeleteTag(db *sql.DB, userID, tag string) (sql.Result, error) {
	query := `DELETE FROM tags WHERE user_id=$1 and tag=$2`
	return db.Exec(query, userID, tag)
}

// EnableTag enables tag
func EnableTag(db *sql.DB, userID, tag string) (sql.Result, error) {
	query := `UPDATE tags SET enabled=$1 WHERE user_id=$2 and tag=$3`
	return db.Exec(query, true, userID, tag)
}

// DisableTag disables tag
func DisableTag(db *sql.DB, userID, tag string) (sql.Result, error) {
	query := `UPDATE tags SET enabled=$1 WHERE user_id=$2 and tag=$3`
	return db.Exec(query, false, userID, tag)
}

// SetWatched sets watch field to given value
func SetWatched(db *sql.DB, userID, tag string, watched bool) (err error) {
	query := `UPDATE tags SET watched=$1 WHERE user_id=$2 and tag=$3`
	_, err = db.Exec(query, watched, userID, tag)
	return
}

// FindUsers retrieves all users who have registered tags
func FindUsers(db *sql.DB) (users []string, err error) {
	rows, err := db.Query("SELECT DISTINCT user_id FROM tags")
	if err != nil {
		return
	}

	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	return
}
