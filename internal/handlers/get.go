package handlers

import (
	"encoding/json"
	"github.com/Hanagasumiiii/docker-track/internal/models"
	"net/http"
)

type ContainerGetter interface {
	GetContainers() ([]models.Container, error)
}

func Get(getter ContainerGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		containers, err := getter.GetContainers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(containers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
