package handlers

import (
    "encoding/json"
    "net/http"

    "server/internal/services"
)

type FoldersHandler struct {
    folderService *services.FolderService
}

func NewFoldersHandler(folderService *services.FolderService) *FoldersHandler {
    return &FoldersHandler{folderService: folderService}
}

// Handle gestiona la solicitud GET y responde con la estructura de folders.
func (h *FoldersHandler) Handle(w http.ResponseWriter, r *http.Request) {
    folders := h.folderService.GetFolders()
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(folders); err != nil {
        http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
    }
}