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
	// Configurar el archivo de log
	logFile, err := os.OpenFile("process5.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creando archivo de log: %v\n", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Iniciando procesamiento")

	// Iniciar perfilado de CPU
	f, err := os.Create("cpu_profile5.prof")
	if err != nil {
		log.Fatalf("Error creando archivo de perfil: %v", err)
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Configuración
	folderPath := "enron_mail_20110402"
	indexName := "emails"

	start := time.Now()
	log.Printf("Procesando carpeta: %s\n", folderPath)

	// Aumentar el número de workers para maximizar el uso de CPU
	numWorkers := 16 // Puedes ajustar este valor para probar diferentes configuraciones

	err = processFolderConcurrent(folderPath, indexName, numWorkers)
	if err != nil {
		log.Printf("Error procesando carpeta: %v\n", err)
	}

	duration := time.Since(start)
	log.Printf("Procesamiento completado en %s\n", duration)
}

func processFolderConcurrent(folderPath, indexName string, numWorkers int) error {
	files := make(chan string, 10000) // Buffer grande para evitar bloqueos
	wg := &sync.WaitGroup{}

	// Iniciar workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(files, indexName, wg)
	}

	// Agregar archivos al canal
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accediendo a %s: %v", path, err)
			return nil // Continuar procesando
		}
		if !info.IsDir() {
			files <- path
		}
		return nil
	})

	close(files)
	wg.Wait()
	return err
}

func worker(files <-chan string, indexName string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{} // Reutilizar el cliente HTTP para eficiencia

	for path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error leyendo archivo %s: %v\n", path, err)
			continue
		}

		email := parseEmail(string(content))
		if email.MessageID == "" {
			log.Printf("Saltando archivo sin Message-ID: %s\n", path)
			continue
		}

		email.Folder = filepath.Join(filepath.Dir(path), filepath.Base(path))
		email.Body = cleanBody(email.Body)

		err = indexDocument(indexName, email, client)
		if err != nil {
			log.Printf("Error indexando documento (%s): %v\n", path, err)
		}
	}
}

func parseEmail(content string) Email {
	lines := strings.Split(content, "\n")
	email := Email{}
	var sb strings.Builder
	readingBody := false

	for _, line := range lines {
		if readingBody {
			sb.WriteString(line)
			sb.WriteString("\n")
		} else if line == "" {
			readingBody = true
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

	email.Body = sb.String()
	return email
}

func cleanBody(body string) string {
	return strings.TrimSpace(strings.ReplaceAll(body, "\n", " "))
}

func indexDocument(indexName string, email Email, client *http.Client) error {
	zincURL := fmt.Sprintf("http://localhost:4080/api/%s/_doc", indexName)
	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("error al serializar email: %w", err)
	}

	req, err := http.NewRequest("POST", zincURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando solicitud HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin123")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error al indexar documento: %s", body)
	}

	log.Printf("Documento indexado con ruta %s\n", email.Folder)
	return nil
}
