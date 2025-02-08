package main

import (
	"docker-track/internal/storage"
	"fmt"
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

	err = s.SaveContainer("192.168.11.11", "in work")
	if err != nil {
		panic(err)
	}

	// TODO: init router: http.ServeMux
	// TODO: run server
}
