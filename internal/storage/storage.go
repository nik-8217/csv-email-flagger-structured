package storage

import (
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
)

func EnsureStorage() error { return os.MkdirAll("storage", 0o755) }

func SaveUpload(file multipart.File, id string) (string, error) {
    path := filepath.Join("storage", id+".upload")
    out, err := os.Create(path)
    if err != nil {
        return "", err
    }
    defer out.Close()
    _, err = io.Copy(out, file)
    return path, err
}
