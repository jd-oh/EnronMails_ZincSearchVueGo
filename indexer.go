package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	_ "net/http/pprof" // Import correcto fuera de main
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

func main() {
	// Iniciar profiling de CPU y guardar en archivo
	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		fmt.Println("Error creating profile file:", err)
		return
	}
	defer f.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Procesamiento del archivo tgz
	tgzFile := "enron_mail_20110402.tgz"
	indexName := "emails"

	start := time.Now()
	err = processTgz(tgzFile, indexName)
	if err != nil {
		fmt.Printf("Error processing tgz: %v\n", err)
	}
	duration := time.Since(start)
	fmt.Printf("Processing completed in %s\n", duration)

	// Esperar 10 segundos para capturar profiling (opcional)
	time.Sleep(10 * time.Second)
}

func processTgz(tgzPath, indexName string) error {
	file, err := os.Open(tgzPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()
	tarReader := tar.NewReader(gzr)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %w", err)
		}
		if header.Typeflag == tar.TypeDir {
			continue
		}

		if header.Typeflag == tar.TypeReg {
			content, err := io.ReadAll(tarReader)
			if err != nil {
				fmt.Printf("Error reading file content: %v\n", err)
				continue
			}

			email := parseEmail(string(content))
			if email.MessageID == "" {
				fmt.Printf("Skipping file due to missing Message-ID: %s\n", header.Name)
				continue
			}

			// Limpieza del cuerpo del email antes de indexar
			email.Body = cleanBody(email.Body)

			err = indexDocument(indexName, email)
			if err != nil {
				fmt.Printf("Error indexing document (%s): %v\n", header.Name, err)
			}
		}
	}
	return nil
}

func parseEmail(content string) Email {
	lines := strings.Split(content, "\n")
	email := Email{}
	readingBody := false
	for _, line := range lines {
		if readingBody {
			email.Body += line + "\n"
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
	return email
}

func cleanBody(body string) string {
	// Remueve caracteres no deseados y espacios adicionales
	return strings.TrimSpace(strings.ReplaceAll(body, "\n", " "))
}

func indexDocument(indexName string, email Email) error {
	zincURL := "http://localhost:4080/api/" + indexName + "/_doc"
	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email: %w", err)
	}

	req, err := http.NewRequest("POST", zincURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to index document: %s", body)
	}
	fmt.Printf("Indexed document: %s\n", email.MessageID)
	return nil
}