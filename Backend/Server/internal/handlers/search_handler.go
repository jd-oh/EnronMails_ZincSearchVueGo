package handlers

import (
    "encoding/json"
    "net/http"

    "server/internal/models"
    "server/internal/services"
)

type SearchHandler struct {
    searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
    return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
    var req models.SearchRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    result, err := h.searchService.Search(req)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(result)
}