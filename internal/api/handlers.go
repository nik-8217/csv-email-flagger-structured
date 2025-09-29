package api

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "csv-email-flagger/internal/jobs"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    id, mode, err := jobs.CreateAndQueue(r)
    if err != nil {
        writeErr(w, http.StatusBadRequest, err)
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"id": id, "mode": mode})
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if _, err := uuid.Parse(id); err != nil {
        writeErr(w, http.StatusBadRequest, errors.New("invalid id"))
        return
    }
    j, ok := jobs.Get(id)
    if !ok {
        writeErr(w, http.StatusBadRequest, errors.New("invalid id"))
        return
    }
    writeJSON(w, http.StatusOK, j)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    jobs.ServeDownload(w, r, id)
}

func SwaggerJSON(w http.ResponseWriter, r *http.Request) {
    spec := `{"openapi":"3.0.3","info":{"title":"CSV Email Flagger API","version":"1.0.0"}}`
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(spec))
}

func Health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}

func writeErr(w http.ResponseWriter, code int, err error) {
    writeJSON(w, code, map[string]string{"error": err.Error()})
}
func writeJSON(w http.ResponseWriter, code int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(v)
}
