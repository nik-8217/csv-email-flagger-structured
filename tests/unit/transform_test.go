package unit

import (
	"bytes"
	"strings"
	"testing"

	"csv-email-flagger/internal/transform"
)

var sampleCSV = `name,email
Alice,alice@example.com
Bob,not-an-email
`
var expectedCSV = `name,email
Alice,alice@example.com,true
Bob,not-an-email,false
`

func TestTransformSequential(t *testing.T) {
	in := strings.NewReader(sampleCSV)
	var out bytes.Buffer
	if err := transform.TransformSequential(in, &out); err != nil {
		t.Fatalf("sequential transform failed: %v", err)
	}
	if out.String() != expectedCSV {
		t.Errorf("unexpected output\nGot:\n%s\nWant:\n%s", out.String(), expectedCSV)
	}
}

func TestTransformParallel(t *testing.T) {
	in := strings.NewReader(sampleCSV)
	var out bytes.Buffer
	if err := transform.TransformParallel(in, &out, 2); err != nil {
		t.Fatalf("parallel transform failed: %v", err)
	}
	if out.String() != expectedCSV {
		t.Errorf("unexpected output\nGot:\n%s\nWant:\n%s", out.String(), expectedCSV)
	}
}
