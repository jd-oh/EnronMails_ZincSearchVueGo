package services

import (
    "log"
    "os"
    "path/filepath"
)

type FolderService struct{}

func NewFolderService() *FolderService {
    return &FolderService{}
}

// GetFolders recorre la carpeta "enron_mail_20110402" y devuelve un mapa donde cada llave es el nombre de una persona 
// y el valor es un slice con los nombres de las subcarpetas (por ejemplo, "contacts", "all_documents", etc.).
func (s *FolderService) GetFolders() map[string][]string {
    folders := make(map[string][]string)

    // Obtener el directorio de trabajo actual.
    cwd, err := os.Getwd()
    if err != nil {
        log.Println("Error obteniendo directorio de trabajo:", err)
        return folders
    }

    // Ajusta la ruta según tu estructura.
    // Si se ejecuta desde Backend/Server/cmd, sube dos niveles y entra a Indexer/enron_mail_20110402
    baseDir := filepath.Join(cwd, "..", "..", "Indexer", "enron_mail_20110402")

    // Si existe un directorio "maildir" dentro de enron_mail_20110402, asumimos que es donde están los nombres.
    // En ese caso, actualizamos la ruta base.
    maildirPath := filepath.Join(baseDir, "maildir")
    if info, err := os.Stat(maildirPath); err == nil && info.IsDir() {
        baseDir = maildirPath
    }

    // Leer los directorios principales (cada persona)
    personDirs, err := os.ReadDir(baseDir)
    if err != nil {
        log.Println("Error leyendo baseDir:", err)
        return folders
    }

    for _, personEntry := range personDirs {
        if personEntry.IsDir() {
            personName := personEntry.Name()
            personFolderPath := filepath.Join(baseDir, personName)

            // Leer las subcarpetas dentro de la carpeta de la persona.
            subEntries, err := os.ReadDir(personFolderPath)
            if err != nil {
                log.Printf("No se pudo leer '%s': %v", personFolderPath, err)
                continue
            }

            var subFolders []string
            for _, subEntry := range subEntries {
                if subEntry.IsDir() {
                    subFolders = append(subFolders, subEntry.Name())
                }
            }
            folders[personName] = subFolders
        }
    }
    return folders
}