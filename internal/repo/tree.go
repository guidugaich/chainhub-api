package repo

import (
	"database/sql"

	"chainhub-api/internal/models"
)

func GetTreeByUsername(db *sql.DB, username string) (models.Tree, error) {
	var tree models.Tree
	err := db.QueryRow(
		`SELECT t.id, t.user_id, t.title, t.created_at
		 FROM trees t
		 JOIN users u ON t.user_id = u.id
		 WHERE u.username = $1`,
		username,
	).Scan(&tree.ID, &tree.UserID, &tree.Title, &tree.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Tree{}, ErrNotFound
		}
		return models.Tree{}, err
	}
	return tree, nil
}

func TreeBelongsToUser(db *sql.DB, treeID, userID int64) (bool, error) {
	var exists bool
	err := db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM trees WHERE id = $1 AND user_id = $2
		)`,
		treeID,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func CreateTree(db *sql.DB, userID int64, title string) (models.Tree, error) {
	var tree models.Tree
	err := db.QueryRow(
		`INSERT INTO trees (user_id, title)
		 VALUES ($1, $2)
		 RETURNING id, user_id, title, created_at`,
		userID,
		title,
	).Scan(&tree.ID, &tree.UserID, &tree.Title, &tree.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Tree{}, ErrDuplicate
		}
		return models.Tree{}, err
	}
	return tree, nil
}
