package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "server/config"
    "server/internal/handlers"
    "server/internal/services"
    customMiddleware "server/internal/middleware"
)

func main() {
    config, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Cannot load config:", err)
    }

    searchService := services.NewSearchService(config)
    searchHandler := handlers.NewSearchHandler(searchService)

    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(customMiddleware.EnableCors)

    r.Route("/api", func(r chi.Router) {
        r.Post("/search", searchHandler.Handle)
    })

    log.Printf("Server listening on port %s", config.ServerPort)
    log.Fatal(http.ListenAndServe(":"+config.ServerPort, r))
}