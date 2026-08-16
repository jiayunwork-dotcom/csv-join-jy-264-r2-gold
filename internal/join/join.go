package join

import (
	"fmt"

	"csv-join/internal/table"
)

// Result is an inner or left join of two tables on a shared key column.
type Result struct {
	Header []string
	Rows   [][]string
}

// Inner keeps only left rows that have a matching key on the right.
func Inner(left, right *table.Table, key string) (*Result, error) {
	return join(left, right, key, false)
}

// Left keeps every left row; missing right columns are empty strings.
func Left(left, right *table.Table, key string) (*Result, error) {
	return join(left, right, key, true)
}

func join(left, right *table.Table, key string, keepUnmatched bool) (*Result, error) {
	li := left.Col(key)
	if li < 0 {
		return nil, fmt.Errorf("join: left has no column %q", key)
	}
	ri := right.Col(key)
	if ri < 0 {
		return nil, fmt.Errorf("join: right has no column %q", key)
	}

	index := map[string][][]string{}
	for _, row := range right.Rows {
		k := row[ri]
		cp := append([]string(nil), row...)
		index[k] = append(index[k], cp)
	}

	header := make([]string, 0, len(left.Header)+len(right.Header)-1)
	header = append(header, left.Header...)
	for i, h := range right.Header {
		if i == ri {
			continue
		}
		header = append(header, h)
	}

	var out [][]string
	for _, lrow := range left.Rows {
		matches := index[lrow[li]]
		if len(matches) == 0 {
			if keepUnmatched {
				row := append([]string(nil), lrow...)
				for i := range right.Header {
					if i == ri {
						continue
					}
					row = append(row, "")
				}
				out = append(out, row)
			}
			continue
		}
		for _, rrow := range matches {
			row := append([]string(nil), lrow...)
			for i := range right.Header {
				if i == ri {
					continue
				}
				row = append(row, rrow[i])
			}
			out = append(out, row)
		}
	}
	return &Result{Header: header, Rows: out}, nil
}
