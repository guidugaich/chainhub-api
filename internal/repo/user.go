package repo

import (
	"database/sql"

	"chainhub-api/internal/models"

	"github.com/lib/pq"
)

func CreateUser(db *sql.DB, email, username, passwordHash string) (models.User, error) {
	var user models.User
	err := db.QueryRow(
		`INSERT INTO users (email, username, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, username, password_hash, wallet_address, created_at`,
		email,
		username,
		passwordHash,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.WalletAddress, &user.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return models.User{}, ErrDuplicate
		}
		return models.User{}, err
	}
	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
	var user models.User
	err := db.QueryRow(
		`SELECT id, email, username, password_hash, wallet_address, created_at
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.WalletAddress, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func GetUserByEmailOrUsername(db *sql.DB, identifier string) (models.User, error) {
	var user models.User
	err := db.QueryRow(
		`SELECT id, email, username, password_hash, wallet_address, created_at
		 FROM users
		 WHERE email = $1 OR username = $1`,
		identifier,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.WalletAddress, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func isUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == "23505"
	}
	return false
}
