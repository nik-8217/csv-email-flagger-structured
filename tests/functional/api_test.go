package functional

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"csv-email-flagger/internal/api"
	"csv-email-flagger/internal/storage"

	"github.com/gorilla/mux"
)

func newTestServer() *httptest.Server {
	r := mux.NewRouter()
	api.RegisterRoutes(r)
	return httptest.NewServer(r)
}

func createMultipartFile(t *testing.T, field, filename, content string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	io.Copy(part, strings.NewReader(content))
	writer.Close()
	return body, writer.FormDataContentType()
}

func TestUploadAndProcess(t *testing.T) {
	_ = storage.EnsureStorage()
	ts := newTestServer()
	defer ts.Close()

	csvContent := `name,email
Alice,alice@example.com
Bob,not-an-email
`
	body, contentType := createMultipartFile(t, "file", "test.csv", csvContent)
	res, err := http.Post(ts.URL+"/api/upload", contentType, body)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("upload returned %d", res.StatusCode)
	}
}

func TestUploadAndProcess_BlankRows(t *testing.T) {
	_ = storage.EnsureStorage()
	ts := newTestServer()
	defer ts.Close()

	csvContent := `name,email
Alice,alice@example.com

Bob,not-an-email
   ,   
Charlie,charlie@test.com
`
	body, contentType := createMultipartFile(t, "file", "test.csv", csvContent)
	res, err := http.Post(ts.URL+"/api/upload", contentType, body)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("upload returned %d", res.StatusCode)
	}
}

func TestUploadAndProcess_ExistingHasEmailHeader(t *testing.T) {
	_ = storage.EnsureStorage()
	ts := newTestServer()
	defer ts.Close()

	csvContent := `name,email,hasEmail
Alice,alice@example.com,false
Bob,not-an-email,true
`
	body, contentType := createMultipartFile(t, "file", "test.csv", csvContent)
	res, err := http.Post(ts.URL+"/api/upload", contentType, body)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("upload returned %d", res.StatusCode)
	}
}

// Negative test cases
func TestUpload_MissingFile(t *testing.T) {
	_ = storage.EnsureStorage()
	ts := newTestServer()
	defer ts.Close()

	res, err := http.Post(ts.URL+"/api/upload", "multipart/form-data", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 for missing file, got %d", res.StatusCode)
	}
}

func TestUpload_EmptyFile(t *testing.T) {
	_ = storage.EnsureStorage()
	ts := newTestServer()
	defer ts.Close()

	body, contentType := createMultipartFile(t, "file", "empty.csv", "")
	res, err := http.Post(ts.URL+"/api/upload", contentType, body)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("upload returned %d", res.StatusCode)
	}

	// Wait a bit for processing
	time.Sleep(100 * time.Millisecond)

	// Check if job failed due to empty file
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	jobID, ok := response["id"].(string)
	if !ok {
		t.Fatal("no job ID in response")
	}

	// Check job status
	statusRes, err := http.Get(ts.URL + "/api/status/" + jobID)
	if err != nil {
		t.Fatalf("status check failed: %v", err)
	}
	defer statusRes.Body.Close()

	var statusResponse map[string]interface{}
	if err := json.NewDecoder(statusRes.Body).Decode(&statusResponse); err != nil {
		t.Fatalf("failed to decode status response: %v", err)
	}

	// Job should be failed due to empty CSV
	if statusResponse["status"] != "FAILED" {
		t.Errorf("expected job to be failed, got %v", statusResponse["status"])
	}
}

func TestStatus_InvalidID(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/status/invalid-id")
	if err != nil {
		t.Fatalf("status check failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 for invalid ID, got %d", res.StatusCode)
	}
}

func TestDownload_InvalidID(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/download/invalid-id")
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	if res.StatusCode != 400 {
		t.Fatalf("expected 400 for invalid ID, got %d", res.StatusCode)
	}
}

func TestDownload_JobInProgress(t *testing.T) {
	_ = storage.EnsureStorage()
	ts := newTestServer()
	defer ts.Close()

	// Create a large CSV to ensure processing takes time
	largeCSV := `name,email
Alice,alice@example.com
Bob,not-an-email
`
	for i := 0; i < 1000; i++ {
		largeCSV += fmt.Sprintf("User%d,user%d@example.com\n", i, i)
	}

	body, contentType := createMultipartFile(t, "file", "large.csv", largeCSV)
	res, err := http.Post(ts.URL+"/api/upload", contentType, body)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("upload returned %d", res.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	jobID, ok := response["id"].(string)
	if !ok {
		t.Fatal("no job ID in response")
	}

	// Try to download immediately (should be in progress)
	downloadRes, err := http.Get(ts.URL + "/api/download/" + jobID)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	if downloadRes.StatusCode != 423 { // 423 Locked
		t.Fatalf("expected 423 for job in progress, got %d", downloadRes.StatusCode)
	}
}

func TestCleanup(t *testing.T) {
	_ = storage.EnsureStorage()
	ts := newTestServer()
	defer ts.Close()

	res, err := http.Post(ts.URL+"/api/cleanup", "application/json", nil)
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("cleanup returned %d", res.StatusCode)
	}
}

func TestHealth(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("health check returned %d", res.StatusCode)
	}
}
