package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
	"log"
)

// Estructura Email
// Esta estructura refleja el formato de los correos electrónicos que vamos a procesar e indexar.
// Cada campo tiene su respectiva etiqueta JSON para facilitar la serialización y deserialización.
type Email struct {
	MessageID string `json:"message_id"`
	Date      string `json:"date"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Folder    string `json:"folder"`
}

func main() {
	// Configuración del archivo de log para registrar el progreso del procesamiento.
	logFile, err := os.OpenFile("Logs/process_bulk.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creando archivo de log: %v\n", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Iniciando procesamiento con bulk index")

	// Iniciar perfilado de CPU para analizar el rendimiento de la aplicación.
	f, err := os.Create("Profiles/cpu_profile_bulk.prof")
	if err != nil {
		log.Fatalf("Error creando archivo de perfil: %v", err)
	}
	defer f.Close()

	// Comienza a recolectar información sobre el uso de la CPU
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Configuración para el procesamiento
	folderPath := "enron_mail_20110402" // Cambia esto a la ruta de tu carpeta
	indexName := "emails"               // El nombre del índice en ZincSearch
	batchSize := 1000                   // Tamaño del lote de documentos a enviar por cada solicitud
	numWorkers := 16                    // Número de workers concurrentes para procesar archivos

	start := time.Now()
	log.Printf("Procesando carpeta: %s\n", folderPath)

	// Llama a la función que procesa la carpeta de manera concurrente
	err = processFolderConcurrent(folderPath, indexName, numWorkers, batchSize)
	if err != nil {
		log.Printf("Error procesando carpeta: %v\n", err)
	}

	// Registra el tiempo de duración del procesamiento
	duration := time.Since(start)
	log.Printf("Procesamiento completado en %s\n", duration)
}

// processFolderConcurrent procesa los archivos en la carpeta de manera concurrente utilizando workers.
// La carpeta de archivos se recorre con filepath.Walk y cada archivo es enviado a los workers.
func processFolderConcurrent(folderPath, indexName string, numWorkers, batchSize int) error {
	// Canal para transmitir las rutas de los archivos a los workers
	files := make(chan string, 10000) // Canal con buffer grande para evitar bloqueos
	wg := &sync.WaitGroup{}

	// Iniciar los workers concurrentes
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go workerBulk(files, indexName, wg, batchSize)
	}

	// Recorrer los archivos en la carpeta y enviarlos al canal
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accediendo a %s: %v", path, err)
			return nil // Continuar procesando otros archivos
		}
		if !info.IsDir() {
			files <- path // Enviar ruta de archivo al canal
		}
		return nil
	})

	// Cerrar el canal cuando se haya procesado todo
	close(files)
	wg.Wait() // Esperar a que todos los workers terminen
	return err
}

// workerBulk es una función que ejecuta el procesamiento de archivos en paralelo y
// envía los lotes de documentos a la API de ZincSearch.
func workerBulk(files <-chan string, indexName string, wg *sync.WaitGroup, batchSize int) {
	defer wg.Done() // Asegura que la goroutine se marca como terminada al finalizar
	client := &http.Client{}
	var bulkData []Email // Slice para almacenar los correos electrónicos a indexar en un lote

	// Procesar los archivos recibidos desde el canal
	for path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error leyendo archivo %s: %v\n", path, err)
			continue
		}

		// Parsear el contenido del archivo a un objeto Email
		email := parseEmail(string(content))
		if email.MessageID == "" {
			log.Printf("Saltando archivo sin Message-ID: %s\n", path)
			continue
		}

		// Establecer el campo "folder" que contiene la ruta del archivo
		email.Folder = filepath.Join(filepath.Dir(path), filepath.Base(path))
		// Limpiar el cuerpo del correo (eliminar saltos de línea innecesarios)
		email.Body = cleanBody(email.Body)

		// Agregar el correo al lote
		bulkData = append(bulkData, email)

		// Si el lote alcanza el tamaño configurado, enviarlo
		if len(bulkData) >= batchSize {
			if err := sendBulk(indexName, bulkData, client); err != nil {
				log.Printf("Error indexando lote: %v\n", err)
			}
			bulkData = nil // Reiniciar el slice para el siguiente lote
		}
	}

	// Enviar el último lote si queda algún correo sin procesar
	if len(bulkData) > 0 {
		if err := sendBulk(indexName, bulkData, client); err != nil {
			log.Printf("Error indexando lote final: %v\n", err)
		}
	}
}

// sendBulk envía un lote de correos electrónicos a la API de ZincSearch
// Utiliza la API _bulk para enviar los documentos de forma eficiente.
func sendBulk(indexName string, emails []Email, client *http.Client) error {
	zincURL := fmt.Sprintf("http://localhost:4080/api/%s/_bulk", indexName)
	var buffer bytes.Buffer

	// Construir el cuerpo del request en formato bulk
	for _, email := range emails {
		action := map[string]interface{}{
			"index": map[string]string{}, // Acción de indexación
		}
		actionJSON, _ := json.Marshal(action)
		emailJSON, _ := json.Marshal(email)

		// Escribir la acción y el documento en el buffer
		buffer.Write(actionJSON)
		buffer.WriteByte('\n')
		buffer.Write(emailJSON)
		buffer.WriteByte('\n')
	}

	// Crear la solicitud HTTP POST
	req, err := http.NewRequest("POST", zincURL, &buffer)
	if err != nil {
		return fmt.Errorf("error creando solicitud HTTP bulk: %w", err)
	}

	// Configurar los encabezados para la solicitud
	req.Header.Set("Content-Type", "application/json")

	// Obtener las credenciales desde las variables de entorno para autenticación básica
	username := os.Getenv("ZINC_FIRST_ADMIN_USER")
	password := os.Getenv("ZINC_FIRST_ADMIN_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Las credenciales no están configuradas en las variables de entorno.")
	}
	req.SetBasicAuth(username, password) // Autenticación básica

	// Realizar la solicitud HTTP
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando solicitud HTTP bulk: %w", err)
	}
	defer resp.Body.Close()

	// Verificar la respuesta HTTP
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error al indexar lote: %s", body)
	}

	log.Printf("Lote indexado con %d documentos\n", len(emails))
	return nil
}

// parseEmail analiza el contenido del correo electrónico y extrae los campos relevantes
// como Message-ID, Fecha, De, Para, Asunto y Cuerpo.
func parseEmail(content string) Email {
	lines := strings.Split(content, "\n")
	email := Email{}
	var sb strings.Builder
	readingBody := false

	// Recorrer cada línea del contenido del correo
	for _, line := range lines {
		if readingBody {
			sb.WriteString(line)
			sb.WriteString("\n")
		} else if line == "" {
			readingBody = true // Empezar a leer el cuerpo del correo
		} else if strings.HasPrefix(line, "Message-ID:") {
			email.MessageID = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		} else if strings.HasPrefix(line, "Date:") {
			email.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		} else if strings.HasPrefix(line, "From:") {
			email.From = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		} else if strings.HasPrefix(line, "To:") {
			email.To = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		} else if strings.HasPrefix(line, "Subject:") {
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		}
	}

	// Establecer el cuerpo del correo
	email.Body = sb.String()
	return email
}

// cleanBody limpia el cuerpo del correo, eliminando saltos de línea innecesarios
// y dejando solo un espacio limpio entre las líneas del cuerpo.
func cleanBody(body string) string {
	return strings.TrimSpace(strings.ReplaceAll(body, "\n", " "))
}
