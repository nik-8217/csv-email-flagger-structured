package jobs

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"csv-email-flagger/internal/storage"
	"csv-email-flagger/internal/transform"
	"csv-email-flagger/pkg/logger"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateAndQueue handles file upload, saves it, creates job, and spawns processing goroutine
func CreateAndQueue(r *http.Request) (string, string, error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return "", "", err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", "", errors.New("missing file")
	}
	defer file.Close()

	id := uuid.NewString()
	inPath, err := storage.SaveUpload(file, id)
	if err != nil {
		return "", "", err
	}

	mode := os.Getenv("PROCESS_MODE")
	if mode == "" {
		mode = "sequential"
	}

	j := &Job{
		ID:        id,
		Status:    StatusQueued,
		InputPath: inPath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Mode:      mode,
	}
	Jobs.Create(j)

	// process in background
	go processJob(j)

	return id, mode, nil
}

func processJob(j *Job) {
	log := logger.Log.WithFields(logrus.Fields{"job_id": j.ID, "mode": j.Mode})
	Jobs.SetStatus(j.ID, StatusInProgress, nil)

	// Open input file
	in, err := os.Open(j.InputPath)
	if err != nil {
		Jobs.SetStatus(j.ID, StatusFailed, err)
		log.WithError(err).Error("failed to open input file")
		return
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil {
			log.WithError(closeErr).Warn("failed to close input file")
		}
	}()

	// Create output file
	outPath := storage.GetProcessedFilePath(j.ID)
	out, err := os.Create(outPath)
	if err != nil {
		Jobs.SetStatus(j.ID, StatusFailed, err)
		log.WithError(err).Error("failed to create output file")
		return
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			log.WithError(closeErr).Warn("failed to close output file")
		}
	}()

	// Process depending on mode
	if j.Mode == "parallel" {
		err = transform.TransformParallel(in, out, 4)
	} else {
		err = transform.TransformSequential(in, out)
	}

	if err != nil {
		// Clean up output file on error
		if removeErr := os.Remove(outPath); removeErr != nil {
			log.WithError(removeErr).Warn("failed to remove output file after error")
		}
		Jobs.SetStatus(j.ID, StatusFailed, err)
		log.WithError(err).Error("processing failed")
		return
	}

	// Update job with output path and mark as done
	j.Output = outPath
	Jobs.SetStatus(j.ID, StatusDone, nil)
	log.Info("job completed successfully")
}

// Get retrieves job by id
func Get(id string) (*Job, bool) {
	return Jobs.Get(id)
}

// ServeDownload streams the processed file if available
func ServeDownload(w http.ResponseWriter, r *http.Request, id string) {
	j, ok := Jobs.Get(id)
	if !ok {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	switch j.Status {
	case StatusDone:
		f, err := os.Open(j.Output)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/csv")
		http.ServeContent(w, r, filepath.Base(j.Output), time.Now(), f)
	case StatusFailed:
		http.Error(w, "invalid id", http.StatusBadRequest)
	default:
		http.Error(w, "job in progress", http.StatusLocked)
	}
}
