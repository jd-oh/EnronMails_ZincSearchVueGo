package handlers

//Se encarga de gestionar las solicitudes HTTP de búsqueda. Decodifica el JSON entrante en una 
//estructura de solicitud, valida y establece un campo predeterminado (por ejemplo, "body" en caso de 
//ausencia), y luego invoca el servicio de búsqueda para obtener los resultados.
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