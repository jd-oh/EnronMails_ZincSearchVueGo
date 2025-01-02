package main

import (
	"archive/tar"        // Para leer archivos tar
	"bytes"              // Manejo de buffers
	"compress/gzip"      // Descompresión de archivos gzip
	"encoding/json"      // Codificación y decodificación JSON
	"fmt"                // Funciones de formato de texto
	"io"                 // Lectura y escritura de streams
	"net/http"           // Enviar solicitudes HTTP
	"os"                 // Manejo de archivos y sistema operativo
	"runtime/pprof"      // Perfilado de rendimiento de CPU
	"strings"            // Manipulación de cadenas de texto
	"sync"               // Para sincronización de goroutines
	"time"               // Medir tiempo y control de espera
	_ "net/http/pprof"   // Habilitar servidor HTTP para depuración de perfiles
)

// Email estructura para almacenar datos de correo
type Email struct {
	MessageID string `json:"message_id"`
	Date      string `json:"date"`
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

var client = &http.Client{} // Cliente HTTP reutilizable

func main() {
	// Iniciar perfilado de CPU y guardar en archivo
	f, err := os.Create("cpu_profile2.prof")
	if err != nil {
		fmt.Println("Error creating profile file:", err)
		return
	}
	defer f.Close()

	pprof.StartCPUProfile(f)         // Inicia el perfilado de CPU
	defer pprof.StopCPUProfile()     // Asegura que el perfilado se detenga al finalizar

	// Configuración inicial
	tgzFile := "enron_mail_20110402.tgz" // Archivo TGZ que contiene los correos
	indexName := "emails"                // Nombre del índice en ZincSearch

	start := time.Now()                  // Marca el inicio del procesamiento
	err = processTgz(tgzFile, indexName) // Procesa el archivo TGZ
	if err != nil {
		fmt.Printf("Error processing tgz: %v\n", err)
	}
	duration := time.Since(start) // Calcula el tiempo total de procesamiento
	fmt.Printf("Processing completed in %s\n", duration)
}

func processTgz(tgzPath, indexName string) error {
	file, err := os.Open(tgzPath) // Abre el archivo TGZ
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file) // Crea un lector GZIP para descomprimir el archivo
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tarReader := tar.NewReader(gzr) // Crea un lector TAR para leer el contenido

	// Itera sobre los archivos en el archivo TAR
	var emails []Email
	for {
		header, err := tarReader.Next()
		if err == io.EOF { // EOF indica que se llegó al final del archivo
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %w", err)
		}
		if header.Typeflag == tar.TypeDir { // Ignora directorios
			continue
		}

		// Procesar archivos regulares
		if header.Typeflag == tar.TypeReg {
			buffer := make([]byte, 4096) // Usar un buffer de tamaño fijo
			content := []byte{}
			for {
				n, err := tarReader.Read(buffer)
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("Error reading tar file: %w", err)
				}
				content = append(content, buffer[:n]...)
			}

			email := parseEmail(string(content)) // Convierte el contenido en un objeto Email
			if email.MessageID == "" {           // Verifica que el correo tenga un Message-ID
				fmt.Printf("Skipping file due to missing Message-ID: %s\n", header.Name)
				continue
			}

			email.Body = cleanBody(email.Body) // Limpia el cuerpo del correo
			emails = append(emails, email)     // Añade el correo a la lista
		}
	}

	// Procesar todos los correos concurrentemente
	indexDocumentsConcurrently(indexName, emails)
	return nil
}

// Optimización: Usa un mapa para evitar múltiples "HasPrefix"
func parseEmail(content string) Email {
	lines := strings.Split(content, "\n") // Divide el contenido en líneas
	email := Email{}
	readingBody := false

	emailFields := map[string]*string{
		"Message-ID:": &email.MessageID,
		"Date:":       &email.Date,
		"From:":       &email.From,
		"To:":         &email.To,
		"Subject:":    &email.Subject,
	}

	for _, line := range lines {
		if readingBody {
			email.Body += line + "\n"
		} else if line == "" { // Detecta el inicio del cuerpo del correo
			readingBody = true
		} else {
			for prefix, field := range emailFields {
				if strings.HasPrefix(line, prefix) {
					*field = strings.TrimSpace(strings.TrimPrefix(line, prefix))
				}
			}
		}
	}
	return email
}

func cleanBody(body string) string {
	return strings.Join(strings.Fields(body), " ") // Elimina saltos de línea y múltiples espacios
}

func indexDocumentsConcurrently(indexName string, emails []Email) {
	var wg sync.WaitGroup
	for _, email := range emails {
		wg.Add(1)
		go func(email Email) {
			defer wg.Done()
			err := indexDocument(indexName, email)
			if err != nil {
				fmt.Printf("Error indexing document (%s): %v\n", email.MessageID, err)
			}
		}(email)
	}
	wg.Wait() // Espera que todas las goroutines terminen
}

func indexDocument(indexName string, email Email) error {
	zincURL := "http://localhost:4080/api/" + indexName + "/_doc" // URL del índice en ZincSearch

	jsonData, err := json.Marshal(email) // Convierte el correo a JSON
	if err != nil {
		return fmt.Errorf("failed to marshal email: %w", err)
	}

	req, err := http.NewRequest("POST", zincURL, bytes.NewBuffer(jsonData)) // Crea una solicitud HTTP
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin123") // Autenticación básica

	resp, err := client.Do(req) // Envía la solicitud
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated { // Verifica errores
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to index document: %s", body)
	}
	fmt.Printf("Indexed document: %s\n", email.MessageID)
	return nil
}
