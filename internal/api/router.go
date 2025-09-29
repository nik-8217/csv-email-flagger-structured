package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/upload", UploadHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/status/{id}", StatusHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/download/{id}", DownloadHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/cleanup", CleanupHandler).Methods(http.MethodPost)
	r.HandleFunc("/swagger.json", SwaggerJSON).Methods(http.MethodGet)
	r.HandleFunc("/healthz", Health).Methods(http.MethodGet)
}
