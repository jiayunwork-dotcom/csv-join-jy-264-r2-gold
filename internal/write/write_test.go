package write

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"csv-join/internal/join"
)

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write fail")
}

func TestCSVOK(t *testing.T) {
	res := &join.Result{
		Header: []string{"id", "name"},
		Rows:   [][]string{{"1", "Alice"}},
	}
	var buf bytes.Buffer
	if err := CSV(&buf, res); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "id,name") || !strings.Contains(got, "1,Alice") {
		t.Fatalf("missing content: %q", got)
	}
}

func TestCSVFlushError(t *testing.T) {
	res := &join.Result{Header: []string{"id"}, Rows: [][]string{{"1"}}}
	if err := CSV(failWriter{}, res); err == nil {
		t.Fatal("expected flush/write error to propagate")
	}
}
