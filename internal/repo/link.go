package repo

import (
	"database/sql"

	"chainhub-api/internal/models"
)

func ListActiveLinksByTreeID(db *sql.DB, treeID int64) ([]models.Link, error) {
	rows, err := db.Query(
		`SELECT id, tree_id, title, url, position, is_active, created_at
		 FROM links
		 WHERE tree_id = $1 AND is_active = TRUE
		 ORDER BY position ASC, id ASC`,
		treeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var link models.Link
		if err := rows.Scan(&link.ID, &link.TreeID, &link.Title, &link.URL, &link.Position, &link.IsActive, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

func CreateLink(db *sql.DB, treeID int64, title, url string, position int, isActive bool) (models.Link, error) {
	var link models.Link
	err := db.QueryRow(
		`INSERT INTO links (tree_id, title, url, position, is_active)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, tree_id, title, url, position, is_active, created_at`,
		treeID,
		title,
		url,
		position,
		isActive,
	).Scan(&link.ID, &link.TreeID, &link.Title, &link.URL, &link.Position, &link.IsActive, &link.CreatedAt)
	if err != nil {
		return models.Link{}, err
	}
	return link, nil
}

func ListLinksByTreeIDAndUser(db *sql.DB, treeID, userID int64) ([]models.Link, error) {
	rows, err := db.Query(
		`SELECT l.id, l.tree_id, l.title, l.url, l.position, l.is_active, l.created_at
		 FROM links l
		 JOIN trees t ON l.tree_id = t.id
		 WHERE t.user_id = $1 AND t.id = $2
		 ORDER BY l.position ASC, l.id ASC`,
		userID,
		treeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var link models.Link
		if err := rows.Scan(&link.ID, &link.TreeID, &link.Title, &link.URL, &link.Position, &link.IsActive, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

func ListLinksByUsernameAndUser(db *sql.DB, username string, userID int64) ([]models.Link, error) {
	rows, err := db.Query(
		`SELECT l.id, l.tree_id, l.title, l.url, l.position, l.is_active, l.created_at
		 FROM trees t
		 JOIN users u ON t.user_id = u.id
		 LEFT JOIN links l ON l.tree_id = t.id
		 WHERE t.user_id = $1 AND u.username = $2
		 ORDER BY l.position ASC, l.id ASC`,
		userID,
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	found := false
	for rows.Next() {
		found = true

		var linkID sql.NullInt64
		var treeID sql.NullInt64
		var title sql.NullString
		var url sql.NullString
		var position sql.NullInt64
		var isActive sql.NullBool
		var createdAt sql.NullTime

		if err := rows.Scan(&linkID, &treeID, &title, &url, &position, &isActive, &createdAt); err != nil {
			return nil, err
		}

		if linkID.Valid {
			links = append(links, models.Link{
				ID:        linkID.Int64,
				TreeID:    treeID.Int64,
				Title:     title.String,
				URL:       url.String,
				Position:  int(position.Int64),
				IsActive:  isActive.Bool,
				CreatedAt: createdAt.Time,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNotFound
	}
	return links, nil
}

func UpdateLinkByIDAndUser(db *sql.DB, linkID, userID int64, title, url string, position int, isActive bool) (models.Link, error) {
	var link models.Link
	err := db.QueryRow(
		`UPDATE links l
		 SET title = $1, url = $2, position = $3, is_active = $4
		 FROM trees t
		 WHERE l.tree_id = t.id AND l.id = $5 AND t.user_id = $6
		 RETURNING l.id, l.tree_id, l.title, l.url, l.position, l.is_active, l.created_at`,
		title,
		url,
		position,
		isActive,
		linkID,
		userID,
	).Scan(&link.ID, &link.TreeID, &link.Title, &link.URL, &link.Position, &link.IsActive, &link.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Link{}, ErrNotFound
		}
		return models.Link{}, err
	}
	return link, nil
}

func DeleteLinkByIDAndUser(db *sql.DB, linkID, userID int64) error {
	result, err := db.Exec(
		`DELETE FROM links l
		 USING trees t
		 WHERE l.tree_id = t.id AND l.id = $1 AND t.user_id = $2`,
		linkID,
		userID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
