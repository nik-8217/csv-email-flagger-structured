package transform

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"sync"
)

type Row struct {
	Index int
	Data  []string
}
type Result struct {
	Index int
	Data  []string
	Err   error
	Skip  bool
}

func TransformParallel(in io.Reader, out io.Writer, workerCount int) error {
	cr := csv.NewReader(in)
	cr.FieldsPerRecord = -1
	cw := csv.NewWriter(out)
	defer cw.Flush()

	rowChan := make(chan Row, 1000)
	resChan := make(chan Result, 1000)
	var wg sync.WaitGroup
	var headerProcessed bool
	var headerMutex sync.Mutex

	// start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowChan {
				if row.Index == 0 {
					// Process header row
					headerMutex.Lock()
					if !headerProcessed {
						// Add hasEmail column to header if not already present
						hasEmailHeader := "hasEmail"
						headerExists := false
						for _, field := range row.Data {
							if strings.TrimSpace(strings.ToLower(field)) == "hasemail" {
								headerExists = true
								break
							}
						}
						if !headerExists {
							row.Data = append(row.Data, hasEmailHeader)
						}
						headerProcessed = true
					}
					headerMutex.Unlock()
					resChan <- Result{Index: row.Index, Data: row.Data}
					continue
				}

				// Skip completely empty rows (all fields are empty or whitespace)
				isEmpty := true
				for _, field := range row.Data {
					if strings.TrimSpace(field) != "" {
						isEmpty = false
						break
					}
				}
				if isEmpty {
					// Send a special result to indicate this row should be skipped
					resChan <- Result{Index: row.Index, Data: nil, Skip: true}
					continue
				}

				// Check for email in the row data
				hasEmail := IsValidEmail(strings.Join(row.Data, " "))
				row.Data = append(row.Data, fmt.Sprintf("%t", hasEmail))
				resChan <- Result{Index: row.Index, Data: row.Data}
			}
		}()
	}

	// close result channel when workers finish
	go func() {
		wg.Wait()
		close(resChan)
	}()

	// feed rows
	go func() {
		idx := 0
		for {
			rec, err := cr.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				resChan <- Result{Index: idx, Err: fmt.Errorf("error reading CSV row %d: %w", idx+1, err)}
				break
			}
			rowChan <- Row{Index: idx, Data: rec}
			idx++
		}
		close(rowChan)
	}()

	// maintain order
	pending := make(map[int]Result)
	next := 0
	rowsWritten := 0

	for res := range resChan {
		if res.Err != nil {
			return res.Err
		}
		pending[res.Index] = res
		for {
			if r, ok := pending[next]; ok {
				// Skip empty rows
				if r.Skip {
					delete(pending, next)
					next++
					continue
				}
				if err := cw.Write(r.Data); err != nil {
					return fmt.Errorf("error writing data row %d: %w", next+1, err)
				}
				rowsWritten++
				delete(pending, next)
				next++
			} else {
				break
			}
		}
	}

	// Ensure we processed at least a header
	if rowsWritten == 0 {
		return fmt.Errorf("CSV file appears to be empty or invalid")
	}

	return nil
}
