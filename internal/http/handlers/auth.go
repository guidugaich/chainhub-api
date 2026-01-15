package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"chainhub-api/internal/repo"
	"chainhub-api/internal/services"

	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
	User  userResponse `json:"user"`
	TreeID int64 `json:"tree_id"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	email := strings.TrimSpace(req.Email)
	username := strings.TrimSpace(req.Username)
	if email == "" || req.Password == "" || username == "" {
		writeError(w, http.StatusBadRequest, "email, password, and username are required")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := repo.CreateUser(h.DB, email, username, string(hash))
	if err != nil {
		if errors.Is(err, repo.ErrDuplicate) {
			writeError(w, http.StatusConflict, "email or username already in use")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	treeTitle := fmt.Sprintf("%s's Link Tree", username)
	tree, err := repo.CreateTree(h.DB, user.ID, treeTitle)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create tree")
		return
	}

	token, err := services.GenerateJWT(h.Config.JWTSecret, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	writeJSON(w, http.StatusCreated, authResponse{
		Token: token,
		User: userResponse{
			ID:    user.ID,
			Email: user.Email.String,
			Username:  user.Username.String,
		},
		TreeID: tree.ID,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	identifier := strings.TrimSpace(req.Email)
	if identifier == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email or username and password are required")
		return
	}

	user, err := repo.GetUserByEmailOrUsername(h.DB, identifier)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load user")
		return
	}

	if !user.PasswordHash.Valid || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := services.GenerateJWT(h.Config.JWTSecret, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		Token: token,
		User: userResponse{
			ID:    user.ID,
			Email: user.Email.String,
			Username:  user.Username.String,
		},
	})
}
