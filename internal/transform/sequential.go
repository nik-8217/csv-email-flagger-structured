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
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if rowIdx == 0 {
			if err := cw.Write(rec); err != nil {
				return err
			}
			rowIdx++
			continue
		}
		nonEmpty := false
		for _, f := range rec {
			if strings.TrimSpace(f) != "" {
				nonEmpty = true
				break
			}
		}
		if !nonEmpty {
			continue
		}
		hasEmail := emailRe.MatchString(strings.Join(rec, " "))
		rec = append(rec, fmt.Sprintf("%t", hasEmail))
		if err := cw.Write(rec); err != nil {
			return err
		}
		rowIdx++
	}
	return nil
}
