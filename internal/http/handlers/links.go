package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"chainhub-api/internal/http/middleware"
	"chainhub-api/internal/repo"

	"github.com/go-chi/chi/v5"
)

type createLinkRequest struct {
	TreeID   int64  `json:"tree_id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Position int    `json:"position"`
	IsActive *bool  `json:"is_active"`
}

type updateLinkRequest struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Position int    `json:"position"`
	IsActive bool   `json:"is_active"`
}

type linkListResponse struct {
	Links []linkResponse `json:"links"`
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.TreeID <= 0 {
		writeError(w, http.StatusBadRequest, "tree_id is required")
		return
	}

	title := strings.TrimSpace(req.Title)
	url := strings.TrimSpace(req.URL)
	if title == "" || url == "" {
		writeError(w, http.StatusBadRequest, "title and url are required")
		return
	}

	ownsTree, err := repo.TreeBelongsToUser(h.DB, req.TreeID, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to validate tree")
		return
	}
	if !ownsTree {
		writeError(w, http.StatusNotFound, "tree not found")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	link, err := repo.CreateLink(h.DB, req.TreeID, title, url, req.Position, isActive)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create link")
		return
	}

	writeJSON(w, http.StatusCreated, linkResponse{
		ID:       link.ID,
		Title:    link.Title,
		URL:      link.URL,
		Position: link.Position,
	})
}

func (h *Handler) ListLinks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	username := strings.TrimSpace(r.URL.Query().Get("username"))
	if username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}

	links, err := repo.ListLinksByUsernameAndUser(h.DB, username, userID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, "tree not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load links")
		return
	}

	respLinks := make([]linkResponse, 0, len(links))
	for _, link := range links {
		respLinks = append(respLinks, linkResponse{
			ID:       link.ID,
			Title:    link.Title,
			URL:      link.URL,
			Position: link.Position,
		})
	}

	writeJSON(w, http.StatusOK, linkListResponse{Links: respLinks})
}

func (h *Handler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	linkID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || linkID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid link id")
		return
	}

	var req updateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	title := strings.TrimSpace(req.Title)
	url := strings.TrimSpace(req.URL)
	if title == "" || url == "" {
		writeError(w, http.StatusBadRequest, "title and url are required")
		return
	}

	link, err := repo.UpdateLinkByIDAndUser(h.DB, linkID, userID, title, url, req.Position, req.IsActive)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, "link not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update link")
		return
	}

	writeJSON(w, http.StatusOK, linkResponse{
		ID:       link.ID,
		Title:    link.Title,
		URL:      link.URL,
		Position: link.Position,
	})
}

func (h *Handler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	linkID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || linkID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid link id")
		return
	}

	err = repo.DeleteLinkByIDAndUser(h.DB, linkID, userID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, "link not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete link")
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}
