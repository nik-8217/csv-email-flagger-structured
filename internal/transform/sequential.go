package transform

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

var emailRe = EmailRegex()

func TransformSequential(in io.Reader, out io.Writer) error {
	cr := csv.NewReader(in)
	cr.FieldsPerRecord = -1
	cw := csv.NewWriter(out)
	defer cw.Flush()

	rowIdx := 0
	headerAdded := false

	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV row %d: %w", rowIdx+1, err)
		}

		// Handle header row
		if rowIdx == 0 {
			// Add hasEmail column to header if not already present
			hasEmailHeader := "hasEmail"
			headerExists := false
			for _, field := range rec {
				if strings.TrimSpace(strings.ToLower(field)) == "hasemail" {
					headerExists = true
					break
				}
			}
			if !headerExists {
				rec = append(rec, hasEmailHeader)
			}
			headerAdded = true

			if err := cw.Write(rec); err != nil {
				return fmt.Errorf("error writing header: %w", err)
			}
			rowIdx++
			continue
		}

		// Skip completely empty rows (all fields are empty or whitespace)
		isEmpty := true
		for _, field := range rec {
			if strings.TrimSpace(field) != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		// Check for email in the row data
		hasEmail := IsValidEmail(strings.Join(rec, " "))
		rec = append(rec, fmt.Sprintf("%t", hasEmail))

		if err := cw.Write(rec); err != nil {
			return fmt.Errorf("error writing data row %d: %w", rowIdx+1, err)
		}
		rowIdx++
	}

	// Ensure we processed at least a header
	if !headerAdded {
		return fmt.Errorf("CSV file appears to be empty or invalid")
	}

	return nil
}
