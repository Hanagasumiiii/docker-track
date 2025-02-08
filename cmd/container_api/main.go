package main

import (
	"docker-track/internal/handlers"
	"docker-track/internal/storage"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO: init config
	// TODO: init logger

	s, err := storage.Connect("host=localhost user=q password=q dbname=q sslmode=disable")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("DONE")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/containers/add", handlers.Add(s))
	mux.HandleFunc("/containers/delete", handlers.Delete(s))
	mux.HandleFunc("/containers/get", handlers.Get(s))

	if err = http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

	// TODO: init router: http.ServeMux
	// TODO: run server
}
