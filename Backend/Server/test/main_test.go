package main

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"
    "testing"

    "server/config"
    "server/internal/handlers"
    "server/internal/models"
    "server/internal/services"
)

// MockTransport es una estructura que implementa http.RoundTripper para mockear respuestas.
type MockTransport struct {
    mockResponse *http.Response
    mockError    error
}

// RoundTrip implementa la interfaz http.RoundTripper.
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.mockResponse, m.mockError
}

func newTestSearchHandler(t *testing.T) http.HandlerFunc {
    // Usamos una configuración dummy o de test
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatalf("Error en LoadConfig: %v", err)
    }
    // Crea el servicio de búsqueda y el handler
    searchService := services.NewSearchService(cfg)
    h := handlers.NewSearchHandler(searchService)
    return h.Handle
}

func TestSearchHandler(t *testing.T) {
    os.Setenv("ZINC_FIRST_ADMIN_USER", "testuser")
    os.Setenv("ZINC_FIRST_ADMIN_PASSWORD", "testpass")
    defer os.Unsetenv("ZINC_FIRST_ADMIN_USER")
    defer os.Unsetenv("ZINC_FIRST_ADMIN_PASSWORD")

    mockZincResponse := `{
        "hits": {
            "hits": [
                { "_source": { "subject": "Test Email 1", "from": "sender@example.com", "to": "recipient@example.com", "date": "2023-01-01", "body": "Hello World", "message_id": "1", "folder": "inbox" } },
                { "_source": { "subject": "Test Email 2", "from": "sender2@example.com", "to": "recipient2@example.com", "date": "2023-01-02", "body": "Hello Again", "message_id": "2", "folder": "inbox" } }
            ],
            "total": { "value": 2 }
        }
    }`

    mockTransport := &MockTransport{
        mockResponse: &http.Response{
            StatusCode: 200,
            Body:       ioutil.NopCloser(bytes.NewBufferString(mockZincResponse)),
            Header:     make(http.Header),
        },
        mockError: nil,
    }
    originalTransport := http.DefaultTransport
    http.DefaultTransport = mockTransport
    defer func() { http.DefaultTransport = originalTransport }()

    searchReq := models.SearchRequest{
        Term:  "Hello",
        From:  0,
        Size:  5,
        Field: "body",
    }
    reqBody, err := json.Marshal(searchReq)
    if err != nil {
        t.Fatalf("Error al marshallizar la solicitud: %v", err)
    }

    req := httptest.NewRequest("POST", "/api/search", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()

    handlerFunc := newTestSearchHandler(t)
    handlerFunc(recorder, req)

    res := recorder.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        t.Errorf("Esperado status code %d, obtenido %d", http.StatusOK, res.StatusCode)
    }

    bodyBytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.Fatalf("Error al leer el cuerpo de la respuesta: %v", err)
    }
    trimmedBody := strings.TrimSpace(string(bodyBytes))
    trimmedExpected := strings.TrimSpace(mockZincResponse)
    if trimmedBody != trimmedExpected {
        t.Errorf("Cuerpo de respuesta no coincide.\nEsperado: %s\nObtenido: %s", trimmedExpected, trimmedBody)
    }
}

func TestSearchHandler_InvalidPayload(t *testing.T) {
    os.Setenv("ZINC_FIRST_ADMIN_USER", "testuser")
    os.Setenv("ZINC_FIRST_ADMIN_PASSWORD", "testpass")
    defer os.Unsetenv("ZINC_FIRST_ADMIN_USER")
    defer os.Unsetenv("ZINC_FIRST_ADMIN_PASSWORD")

    req := httptest.NewRequest("POST", "/api/search", strings.NewReader("invalid json"))
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()

    handlerFunc := newTestSearchHandler(t)
    handlerFunc(recorder, req)

    res := recorder.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusBadRequest {
        t.Errorf("Esperado status code %d, obtenido %d", http.StatusBadRequest, res.StatusCode)
    }

    bodyBytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.Fatalf("Error al leer el cuerpo de la respuesta: %v", err)
    }
    expectedMessage := "Invalid request payload\n"
    if string(bodyBytes) != expectedMessage {
        t.Errorf("Mensaje de error no coincide.\nEsperado: %s\nObtenido: %s", expectedMessage, string(bodyBytes))
    }
}

func TestSearchHandler_DefaultField(t *testing.T) {
    os.Setenv("ZINC_FIRST_ADMIN_USER", "testuser")
    os.Setenv("ZINC_FIRST_ADMIN_PASSWORD", "testpass")
    defer os.Unsetenv("ZINC_FIRST_ADMIN_USER")
    defer os.Unsetenv("ZINC_FIRST_ADMIN_PASSWORD")

    mockZincResponse := `{
        "hits": {
            "hits": [
                { "_source": { "subject": "Test Email Default 1", "from": "sender@example.com", "to": "recipient@example.com", "date": "2023-01-01", "body": "Default Body", "message_id": "1", "folder": "inbox" } }
            ],
            "total": { "value": 1 }
        }
    }`

    mockTransport := &MockTransport{
        mockResponse: &http.Response{
            StatusCode: 200,
            Body:       ioutil.NopCloser(bytes.NewBufferString(mockZincResponse)),
            Header:     make(http.Header),
        },
        mockError: nil,
    }
    originalTransport := http.DefaultTransport
    http.DefaultTransport = mockTransport
    defer func() { http.DefaultTransport = originalTransport }()

    searchReq := models.SearchRequest{
        Term:  "Default",
        From:  0,
        Size:  5,
        Field: "", // Campo vacío para probar el valor por defecto
    }
    reqBody, err := json.Marshal(searchReq)
    if err != nil {
        t.Fatalf("Error al marshallizar la solicitud: %v", err)
    }

    req := httptest.NewRequest("POST", "/api/search", bytes.NewBuffer(reqBody))
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()

    handlerFunc := newTestSearchHandler(t)
    handlerFunc(recorder, req)

    res := recorder.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        t.Errorf("Esperado status code %d, obtenido %d", http.StatusOK, res.StatusCode)
    }

    bodyBytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.Fatalf("Error al leer el cuerpo de la respuesta: %v", err)
    }
    trimmedBody := strings.TrimSpace(string(bodyBytes))
    trimmedExpected := strings.TrimSpace(mockZincResponse)
    if trimmedBody != trimmedExpected {
        t.Errorf("Cuerpo de respuesta no coincide.\nEsperado: %s\nObtenido: %s", trimmedExpected, trimmedBody)
    }
}

// Nota: Para TestSearchHandler_ZincSearchError se recomienda refactorizar el handler
// para que retorne errores en lugar de llamar a log.Fatal y así poder testearlo.
// Se deja comentado o pendiente de refactorización.