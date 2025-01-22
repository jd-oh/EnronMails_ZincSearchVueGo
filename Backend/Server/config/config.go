package config

import "os"

type Config struct {
    ZincSearchURL      string
    ZincSearchUser     string
    ZincSearchPassword string
    ServerPort         string
}

func LoadConfig() (*Config, error) {
    return &Config{
        ZincSearchURL:      getEnvOrDefault("ZINC_SEARCH_URL", "http://localhost:4080"),
        ZincSearchUser:     os.Getenv("ZINC_FIRST_ADMIN_USER"),
        ZincSearchPassword: os.Getenv("ZINC_FIRST_ADMIN_PASSWORD"),
        ServerPort:         getEnvOrDefault("SERVER_PORT", "8080"),
    }, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}