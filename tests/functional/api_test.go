package functional

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
