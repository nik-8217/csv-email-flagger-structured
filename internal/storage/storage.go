package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

const (
	StorageDir      = "storage"
	UploadSuffix    = ".upload"
	ProcessedSuffix = ".csv"
)

func EnsureStorage() error {
	return os.MkdirAll(StorageDir, 0o755)
}

func SaveUpload(file multipart.File, id string) (string, error) {
	path := filepath.Join(StorageDir, id+UploadSuffix)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return path, err
}

// CleanupOldFiles removes files older than the specified duration
func CleanupOldFiles(maxAge time.Duration) error {
	files, err := os.ReadDir(StorageDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-maxAge)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			filePath := filepath.Join(StorageDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				// Log error but continue with other files
				continue
			}
		}
	}

	return nil
}

// GetProcessedFilePath returns the expected path for a processed file
func GetProcessedFilePath(id string) string {
	return filepath.Join(StorageDir, id+ProcessedSuffix)
}

// CleanupJobFiles removes both upload and processed files for a job
func CleanupJobFiles(id string) error {
	uploadPath := filepath.Join(StorageDir, id+UploadSuffix)
	processedPath := filepath.Join(StorageDir, id+ProcessedSuffix)

	var errors []error

	if err := os.Remove(uploadPath); err != nil && !os.IsNotExist(err) {
		errors = append(errors, err)
	}

	if err := os.Remove(processedPath); err != nil && !os.IsNotExist(err) {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors[0] // Return first error
	}

	return nil
}
