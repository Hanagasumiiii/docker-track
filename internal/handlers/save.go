package handlers

import (
	"encoding/json"
	"github.com/Hanagasumiiii/docker-track/internal/models"
	"net/http"
)

type ContainerSaver interface {
	SaveContainer(container models.Container) error
}

func Add(saver ContainerSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//const op = "handlers.save"

		var container models.Container
		if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
			http.Error(w, "JSON decoding failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		err := saver.SaveContainer(container)
		if err != nil {
			http.Error(w, "JSON saving failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
