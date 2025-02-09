package main

import (
	"github.com/Hanagasumiiii/docker-track/internal/handlers"
	"github.com/Hanagasumiiii/docker-track/internal/storage"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	// TODO: init config
	// TODO: init logger

	dsn := os.Getenv("DATABASE_URL")
	s, err := storage.Connect(dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/containers/add", handlers.Add(s))
	mux.HandleFunc("/containers/delete", handlers.Delete(s))
	mux.HandleFunc("/containers/get", handlers.Get(s))
	mux.HandleFunc("/containers/update", handlers.Update(s))

	handler := cors.Default().Handler(mux)

	if err = http.ListenAndServe(":8081", handler); err != nil {
		log.Fatal(err)
	}
}
