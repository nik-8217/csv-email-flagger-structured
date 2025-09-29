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
}

func TransformParallel(in io.Reader, out io.Writer, workerCount int) error {
	cr := csv.NewReader(in)
	cr.FieldsPerRecord = -1
	cw := csv.NewWriter(out)
	defer cw.Flush()

	rowChan := make(chan Row, 1000)
	resChan := make(chan Result, 1000)
	var wg sync.WaitGroup

	// start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowChan {
				if row.Index == 0 {
					resChan <- Result{Index: row.Index, Data: row.Data}
					continue
				}
				nonEmpty := false
				for _, f := range row.Data {
					if strings.TrimSpace(f) != "" {
						nonEmpty = true
						break
					}
				}
				if !nonEmpty {
					continue
				}
				hasEmail := emailRe.MatchString(strings.Join(row.Data, " "))
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
				resChan <- Result{Index: idx, Err: err}
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
	for res := range resChan {
		if res.Err != nil {
			return res.Err
		}
		pending[res.Index] = res
		for {
			if r, ok := pending[next]; ok {
				if err := cw.Write(r.Data); err != nil {
					return err
				}
				delete(pending, next)
				next++
			} else {
				break
			}
		}
	}
	return nil
}
