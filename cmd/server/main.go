package main

import (
    "net/http"
	"github.com/go-chi/chi/v5" 
)
//чтобы go.sum не пустой->docker не падал
func main() {
    r := chi.NewRouter()

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    http.ListenAndServe(":8080", r)
}