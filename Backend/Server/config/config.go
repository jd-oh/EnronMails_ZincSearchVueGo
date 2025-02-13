package config


//Contiene la definición de la estructura de configuración y la función de carga que obtiene 
//parámetros desde variables de entorno o valores por defecto. Esto incluye la URL de ZincSearch, 
//credenciales de autenticación y el puerto del servidor.

import "os"

type Config struct {
    ZincSearchURL      string
    ZincSearchUser     string
    ZincSearchPassword string
    ServerPort         string
    EndpointIndex      string
}

func LoadConfig() (*Config, error) {
    return &Config{
        ZincSearchURL:      getEnvOrDefault("ZINC_SEARCH_URL", "http://localhost:4080"),
        ZincSearchUser:     os.Getenv("ZINC_FIRST_ADMIN_USER"),
        ZincSearchPassword: os.Getenv("ZINC_FIRST_ADMIN_PASSWORD"),
        ServerPort:         getEnvOrDefault("SERVER_PORT", "8080"),
        EndpointIndex:      getEnvOrDefault("ENDPOINT_INDEX", "/api/emails"),
    }, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}