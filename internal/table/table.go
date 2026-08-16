package table

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// Table is a CSV with a header and data rows.
type Table struct {
	Header []string
	Rows   [][]string
}

func colIndex(header []string, name string) (int, bool) {
	for i, h := range header {
		if h == name {
			return i, true
		}
	}
	return -1, false
}

// Col returns the index of name, or -1 if missing.
func (t *Table) Col(name string) int {
	i, ok := colIndex(t.Header, name)
	if !ok {
		return -1
	}
	return i
}

// Parse reads a CSV. UTF-8 BOM is stripped. Lines starting with '#' are skipped.
// Each data row must have exactly as many fields as the header.
func Parse(r io.Reader) (*Table, error) {
	br := bufio.NewReader(r)
	head, err := br.Peek(3)
	if err == nil && bytes.Equal(head, []byte{0xEF, 0xBB, 0xBF}) {
		if _, err := br.Discard(3); err != nil {
			return nil, err
		}
	}
	cr := csv.NewReader(br)
	cr.FieldsPerRecord = -1
	cr.ReuseRecord = false

	var recs [][]string
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("table: %w", err)
		}
		if len(rec) == 0 {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(rec[0]), "#") {
			continue
		}
		recs = append(recs, rec)
	}
	if len(recs) == 0 {
		return &Table{Header: []string{}, Rows: [][]string{}}, nil
	}
	header := recs[0]
	rows := make([][]string, 0, len(recs)-1)
	for i, rec := range recs[1:] {
		if len(rec) != len(header) {
			return nil, fmt.Errorf("table: line %d: expected %d fields, got %d", i+2, len(header), len(rec))
		}
		row := append([]string(nil), rec...)
		rows = append(rows, row)
	}
	return &Table{Header: append([]string(nil), header...), Rows: rows}, nil
}
