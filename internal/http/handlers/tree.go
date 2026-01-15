package handlers

import (
	"errors"
	"net/http"

	"chainhub-api/internal/repo"

	"github.com/go-chi/chi/v5"
)

type treeResponse struct {
	ID    int64  `json:"id"`
	Username  string `json:"username"`
	Title string `json:"title"`
	Links []linkResponse `json:"links"`
}

type linkResponse struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Position int    `json:"position"`
}

func (h *Handler) GetTreeByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}

	tree, err := repo.GetTreeByUsername(h.DB, username)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			writeError(w, http.StatusNotFound, "tree not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to load tree")
		return
	}

	links, err := repo.ListActiveLinksByTreeID(h.DB, tree.ID)
	if err != nil {
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

	writeJSON(w, http.StatusOK, treeResponse{
		ID:    tree.ID,
		Username:  username,
		Title: tree.Title,
		Links: respLinks,
	})
}
