package services


//Contiene la lógica de negocio para ejecutar búsquedas. A partir de la solicitud 
//(por ejemplo, models.SearchRequest), construye la consulta en JSON, realiza la solicitud HTTP 
//a ZincSearch, maneja la respuesta (incluyendo errores y lectura del cuerpo) y retorna 
//los resultados al manejador.
import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    "server/config"
    "server/internal/models"
)

type SearchService struct {
    config *config.Config
    client *http.Client
}

func NewSearchService(config *config.Config) *SearchService {
    return &SearchService{
        config: config,
        client: &http.Client{},
    }
}

func (s *SearchService) Search(req models.SearchRequest) ([]byte, error) {
    if req.Field == "" {
        req.Field = "body"
    }

    query := models.ZincSearchQuery{
        SearchType: "match",
        Query: models.Query{
            Term:  req.Term,
            Field: req.Field,
        },
        From:       req.From,
        MaxResults: req.Size,
        Source:     []string{"subject", "from", "to", "date", "body", "message_id", "folder"},
    }

    // Convert query to JSON
    jsonQuery, err := json.Marshal(query)

    if err != nil {
        return nil, fmt.Errorf("error marshaling query: %w", err)
    }

    // Log outgoing request
    log.Printf("Sending request to ZincSearch: %s", string(jsonQuery))

    request, err := http.NewRequest("POST", s.config.ZincSearchURL+s.config.EndpointIndex+"/_search", bytes.NewBuffer(jsonQuery))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    request.SetBasicAuth(s.config.ZincSearchUser, s.config.ZincSearchPassword)
    request.Header.Set("Content-Type", "application/json")

    response, err := s.client.Do(request)
    if err != nil {
        return nil, fmt.Errorf("error executing request: %w", err)
    }
    defer response.Body.Close()

    // Check response status
    if response.StatusCode != http.StatusOK {
        bodyBytes, _ := ioutil.ReadAll(response.Body)
        return nil, fmt.Errorf("zinc search error: status=%d body=%s", 
            response.StatusCode, string(bodyBytes))
    }

    // Read response body
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    // Log response
    log.Printf("Received response from ZincSearch: %s", string(bodyBytes))

    return bodyBytes, nil
}