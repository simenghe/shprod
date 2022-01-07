package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Shopify production engineer intern challenge
func main() {
	const port = ":8080"
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	fmt.Println("Starting server at PORT :", port)
	log.Fatalln(http.ListenAndServe(port, r))
}
