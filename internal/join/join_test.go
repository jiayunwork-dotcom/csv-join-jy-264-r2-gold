package join

import (
	"testing"

	"csv-join/internal/table"
)

func leftTbl() *table.Table {
	return &table.Table{
		Header: []string{"id", "name"},
		Rows:   [][]string{{"1", "Alice"}, {"2", "Bob"}, {"3", "Carol"}},
	}
}

func rightTbl() *table.Table {
	return &table.Table{
		Header: []string{"id", "city"},
		Rows:   [][]string{{"1", "Beijing"}, {"3", "Shanghai"}},
	}
}

func TestInnerJoin(t *testing.T) {
	res, err := Inner(leftTbl(), rightTbl(), "id")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Rows) != 2 {
		t.Fatalf("rows=%d want 2", len(res.Rows))
	}
	if res.Header[2] != "city" {
		t.Fatalf("header=%v", res.Header)
	}
	if res.Rows[0][2] != "Beijing" || res.Rows[1][0] != "3" {
		t.Fatalf("rows=%v", res.Rows)
	}
}

func TestLeftJoinKeepsUnmatched(t *testing.T) {
	res, err := Left(leftTbl(), rightTbl(), "id")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Rows) != 3 {
		t.Fatalf("rows=%d want 3", len(res.Rows))
	}
	if res.Rows[1][0] != "2" || res.Rows[1][2] != "" {
		t.Fatalf("unmatched row=%v", res.Rows[1])
	}
}

func TestJoinMissingKey(t *testing.T) {
	if _, err := Inner(leftTbl(), rightTbl(), "nope"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestInnerDoesNotMutateInput(t *testing.T) {
	l := leftTbl()
	r := rightTbl()
	origLeft := l.Rows[0][1]
	origRight := r.Rows[0][1]
	res, err := Inner(l, r, "id")
	if err != nil {
		t.Fatal(err)
	}
	res.Rows[0][1] = "CHANGED"
	res.Rows[0][2] = "CHANGED"
	if l.Rows[0][1] != origLeft {
		t.Fatalf("join mutated left row: %q", l.Rows[0][1])
	}
	if r.Rows[0][1] != origRight {
		t.Fatalf("join mutated right row: %q", r.Rows[0][1])
	}
}

func TestInnerResultIndependent(t *testing.T) {
	res1, err := Inner(leftTbl(), rightTbl(), "id")
	if err != nil {
		t.Fatal(err)
	}
	want := res1.Rows[0][1]
	res2, err := Inner(leftTbl(), rightTbl(), "id")
	if err != nil {
		t.Fatal(err)
	}
	res2.Rows[0][1] = "CHANGED"
	if res1.Rows[0][1] != want {
		t.Fatalf("second join mutated first result: %q", res1.Rows[0][1])
	}
}

func TestInnerEmptyNonNil(t *testing.T) {
	l := &table.Table{Header: []string{"id"}, Rows: [][]string{{"9"}}}
	r := &table.Table{Header: []string{"id", "x"}, Rows: [][]string{{"1", "a"}}}
	res, err := Inner(l, r, "id")
	if err != nil {
		t.Fatal(err)
	}
	if res.Rows == nil {
		t.Fatal("empty join must return non-nil Rows")
	}
	if len(res.Rows) != 0 {
		t.Fatalf("rows=%v", res.Rows)
	}
}
