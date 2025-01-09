package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Estructura que modela la solicitud para la búsqueda en el índice.
// Permite configurar términos de búsqueda, paginación y el campo a consultar.
type SearchRequest struct {
	Term  string `json:"term"`  // Término a buscar.
	From  int    `json:"from"`  // Offset para paginación.
	Size  int    `json:"size"`  // Número máximo de resultados.
	Field string `json:"field"` // Campo en el que se realizará la búsqueda.
}

func main() {
	// Crea un enrutador usando el framework chi.
	r := chi.NewRouter()

	// Middlewares que añaden funcionalidades adicionales al servidor.
	r.Use(middleware.Logger)    // Registra información sobre cada solicitud HTTP.
	r.Use(middleware.Recoverer) // Captura y maneja errores inesperados (panics).
	r.Use(enableCors)           // Habilita CORS para solicitudes desde otros orígenes.

	// Define rutas para la API.
	r.Route("/api", func(r chi.Router) {
		r.Post("/search", searchHandler) // Ruta para manejar solicitudes de búsqueda.
	})

	// Inicia el servidor en el puerto 8080.
	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r)) // Detiene la ejecución si ocurre un error al iniciar el servidor.
}

// Maneja las solicitudes de búsqueda y las traduce en consultas al motor de búsqueda.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody SearchRequest
	// Decodifica el cuerpo de la solicitud JSON en una estructura SearchRequest.
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Si no se especifica un campo para la búsqueda, se usa 'body' como predeterminado.
	field := reqBody.Field
	if field == "" {
		field = "body"
	}

	// Construye una consulta en formato JSON para el motor de búsqueda.
	query := fmt.Sprintf(`{
        "search_type": "match",
        "query": {
            "term": "%s",        
            "field": "%s"
        },
        "from": %d,                
        "max_results": %d,         
        "_source": ["subject", "from", "to", "date", "body", "message_id", "folder"] 
    }`, reqBody.Term, field, reqBody.From, reqBody.Size)

	// Crea una nueva solicitud HTTP para enviar la consulta al motor de búsqueda.
	req, err := http.NewRequest("POST", "http://localhost:4080/api/emails/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err) // Detiene el servidor si ocurre un error al crear la solicitud.
	}

	// Obtener las credenciales desde las variables de entorno.
	username := os.Getenv("ZINC_FIRST_ADMIN_USER")
	password := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Las credenciales no están configuradas en las variables de entorno.")
	}
	req.SetBasicAuth(username, password)            // Configura autenticación básica.
	req.Header.Set("Content-Type", "application/json") // Establece el tipo de contenido.

	// Envía la solicitud al motor de búsqueda.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err) // Detiene el servidor si ocurre un error al enviar la solicitud.
	}
	defer resp.Body.Close() // Cierra el cuerpo de la respuesta para liberar recursos.

	// Lee el cuerpo de la respuesta y maneja posibles errores.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err) // Detiene el servidor si ocurre un error al leer la respuesta.
	}

	// Configura la respuesta HTTP con el tipo de contenido JSON y escribe los resultados.
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Middleware que habilita Cross-Origin Resource Sharing (CORS).
// Permite solicitudes desde cualquier origen, útil en desarrollo o cuando el cliente está en otro dominio.
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configura los encabezados CORS.
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permite acceso desde cualquier origen.
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Maneja las solicitudes OPTIONS (preflight) devolviendo un estado 200 sin procesar más.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r) // Pasa la solicitud al siguiente middleware o handler.
	})
}
