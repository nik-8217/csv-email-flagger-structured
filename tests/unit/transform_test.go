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
var expectedCSV = `name,email,hasEmail
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

// Test cases for blank/incomplete rows
func TestTransformSequential_BlankRows(t *testing.T) {
	csvWithBlanks := `name,email
Alice,alice@example.com

Bob,not-an-email
   ,   
Charlie,charlie@test.com
`
	expected := `name,email,hasEmail
Alice,alice@example.com,true
Bob,not-an-email,false
Charlie,charlie@test.com,true
`

	in := strings.NewReader(csvWithBlanks)
	var out bytes.Buffer
	if err := transform.TransformSequential(in, &out); err != nil {
		t.Fatalf("sequential transform failed: %v", err)
	}
	if out.String() != expected {
		t.Errorf("unexpected output\nGot:\n%s\nWant:\n%s", out.String(), expected)
	}
}

func TestTransformParallel_BlankRows(t *testing.T) {
	csvWithBlanks := `name,email
Alice,alice@example.com

Bob,not-an-email
   ,   
Charlie,charlie@test.com
`
	expected := `name,email,hasEmail
Alice,alice@example.com,true
Bob,not-an-email,false
Charlie,charlie@test.com,true
`

	in := strings.NewReader(csvWithBlanks)
	var out bytes.Buffer
	if err := transform.TransformParallel(in, &out, 2); err != nil {
		t.Fatalf("parallel transform failed: %v", err)
	}
	if out.String() != expected {
		t.Errorf("unexpected output\nGot:\n%s\nWant:\n%s", out.String(), expected)
	}
}

// Test cases for existing hasEmail header
func TestTransformSequential_ExistingHasEmailHeader(t *testing.T) {
	csvWithExistingHeader := `name,email,hasEmail
Alice,alice@example.com,false
Bob,not-an-email,true
`
	expected := `name,email,hasEmail
Alice,alice@example.com,false,true
Bob,not-an-email,true,false
`

	in := strings.NewReader(csvWithExistingHeader)
	var out bytes.Buffer
	if err := transform.TransformSequential(in, &out); err != nil {
		t.Fatalf("sequential transform failed: %v", err)
	}
	if out.String() != expected {
		t.Errorf("unexpected output\nGot:\n%s\nWant:\n%s", out.String(), expected)
	}
}

// Test cases for empty CSV
func TestTransformSequential_EmptyCSV(t *testing.T) {
	in := strings.NewReader("")
	var out bytes.Buffer
	err := transform.TransformSequential(in, &out)
	if err == nil {
		t.Fatal("expected error for empty CSV, got nil")
	}
	if !strings.Contains(err.Error(), "empty or invalid") {
		t.Errorf("expected error about empty CSV, got: %v", err)
	}
}

func TestTransformParallel_EmptyCSV(t *testing.T) {
	in := strings.NewReader("")
	var out bytes.Buffer
	err := transform.TransformParallel(in, &out, 2)
	if err == nil {
		t.Fatal("expected error for empty CSV, got nil")
	}
	if !strings.Contains(err.Error(), "empty or invalid") {
		t.Errorf("expected error about empty CSV, got: %v", err)
	}
}

// Test cases for malformed CSV
func TestTransformSequential_MalformedCSV(t *testing.T) {
	malformedCSV := `name,email
Alice,alice@example.com
"Bob,not-an-email
Charlie,charlie@test.com
`
	in := strings.NewReader(malformedCSV)
	var out bytes.Buffer
	err := transform.TransformSequential(in, &out)
	if err == nil {
		t.Fatal("expected error for malformed CSV, got nil")
	}
}

// Test cases for various email formats
func TestTransformSequential_EmailFormats(t *testing.T) {
	emailFormats := `name,email
Alice,alice@example.com
Bob,test.email+tag@domain.co.uk
Charlie,user123@subdomain.example.org
David,not-an-email
Eve,user@domain
Frank,user@domain.
Grace,user@.domain
Henry,user@domain.com.
`
	expected := `name,email,hasEmail
Alice,alice@example.com,true
Bob,test.email+tag@domain.co.uk,true
Charlie,user123@subdomain.example.org,true
David,not-an-email,false
Eve,user@domain,false
Frank,user@domain.,false
Grace,user@.domain,false
Henry,user@domain.com.,false
`

	in := strings.NewReader(emailFormats)
	var out bytes.Buffer
	if err := transform.TransformSequential(in, &out); err != nil {
		t.Fatalf("sequential transform failed: %v", err)
	}
	if out.String() != expected {
		t.Errorf("unexpected output\nGot:\n%s\nWant:\n%s", out.String(), expected)
	}
}

// Test cases for case sensitivity in hasEmail header
func TestTransformSequential_CaseInsensitiveHeader(t *testing.T) {
	csvWithCaseVariations := `name,email,HASEMAIL
Alice,alice@example.com
Bob,not-an-email
`
	expected := `name,email,HASEMAIL
Alice,alice@example.com,true
Bob,not-an-email,false
`

	in := strings.NewReader(csvWithCaseVariations)
	var out bytes.Buffer
	if err := transform.TransformSequential(in, &out); err != nil {
		t.Fatalf("sequential transform failed: %v", err)
	}
	if out.String() != expected {
		t.Errorf("unexpected output\nGot:\n%s\nWant:\n%s", out.String(), expected)
	}
}
