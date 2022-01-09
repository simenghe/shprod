package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

// Shopify production engineer intern challenge
var (
	DB *RealDB
)

func main() {
	const port = ":8080"
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGO_SRV")
	if uri == "" {
		log.Fatal("You must set your 'MONGO_SRV' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	// Start up DB connection
	DB = NewDB(uri)
	defer DB.Disconnect()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", HandleDefault)

	// Group of inventory routes.
	r.Get("/inventory", HandleGetInventory)
	r.Post("/inventory", HandlePostInventory)

	// Replace
	r.Put("/inventory", HandlePutInventory)
	r.Delete("/inventory", HandleDeleteInventory)

	fmt.Println("Starting server at PORT :", port)
	log.Fatalln(http.ListenAndServe(port, r))
}
