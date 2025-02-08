package handlers

import (
	"docker-track/internal/models"
	"encoding/json"
	"net/http"
)

type ContainerDeleter interface {
	DeleteContainer(ip string) error
}

func Delete(deleter ContainerDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var container models.Container
		if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
			http.Error(w, "JSON decoding failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		err := deleter.DeleteContainer(container.Ip)
		if err != nil {
			http.Error(w, "Delete container failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
