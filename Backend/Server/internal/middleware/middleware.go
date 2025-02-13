package middleware

//Implementa funcionalidades de middleware, como EnableCors, que añade los encabezados 
//necesarios para permitir solicitudes CORS. Esto es útil para que clientes alojados en 
//otros dominios puedan acceder a la API sin problemas.
import "net/http"

func EnableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}