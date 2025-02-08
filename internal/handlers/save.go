package handlers

import (
	"docker-track/internal/models"
	"encoding/json"
	"net/http"
)

type Request struct {
	Ip     string `json:"ip"`
	Status string `json:"status"`
}

type ContainerSaver interface {
	SaveContainer(container models.Container) error
}

func New(saver ContainerSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save"

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
