package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

type SearchRequest struct {
    Term  string `json:"term"`
    From  int    `json:"from"`
    Size  int    `json:"size"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    var reqBody SearchRequest
    err := json.NewDecoder(r.Body).Decode(&reqBody)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    query := fmt.Sprintf(`{
        "search_type": "match",
        "query": {
            "term": "%s",
            "field": "body",
            "start_time": "2021-12-25T15:08:48.777Z",
            "end_time": "2025-12-27T15:08:48.777Z"
        },
        "from": %d,
        "max_results": %d,
        "_source": ["subject", "from", "to", "date", "body", "message_id"]
    }`, reqBody.Term, reqBody.From, reqBody.Size)

    req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_search", strings.NewReader(query))
    if err != nil {
        log.Fatal(err)
    }
    req.SetBasicAuth("admin", "admin123")
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(body)
}

// Middleware para habilitar CORS
func enableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")  // Permitir desde cualquier origen (ajustar para producción)
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Manejo de preflight (OPTIONS)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/search", searchHandler)

    // Aplicar middleware de CORS
    handler := enableCors(mux)
    
    log.Println("Servidor escuchando en http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
