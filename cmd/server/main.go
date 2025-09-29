package main

import (
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"

    "csv-email-flagger/internal/api"
    "csv-email-flagger/internal/storage"
    "csv-email-flagger/pkg/logger"
)

func main() {
    logger.Init()
    log := logger.Log

    if err := storage.EnsureStorage(); err != nil {
        log.WithError(err).Fatal("failed to ensure storage")
    }

    r := mux.NewRouter()
    api.RegisterRoutes(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.WithFields(logrus.Fields{
        "port": port,
        "mode": os.Getenv("PROCESS_MODE"),
    }).Info("server starting")

    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.WithError(err).Fatal("server stopped unexpectedly")
    }
}
