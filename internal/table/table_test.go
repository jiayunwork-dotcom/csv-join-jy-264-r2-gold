package table

import (
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	in := "id,name\n1,Alice\n2,Bob\n"
	tbl, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(tbl.Header) != 2 || tbl.Header[0] != "id" {
		t.Fatalf("header=%v", tbl.Header)
	}
	if len(tbl.Rows) != 2 || tbl.Rows[0][1] != "Alice" {
		t.Fatalf("rows=%v", tbl.Rows)
	}
}

func TestParseStripsBOM(t *testing.T) {
	in := "\xef\xbb\xbfid,name\n1,Alice\n"
	tbl, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if tbl.Header[0] != "id" {
		t.Fatalf("BOM not stripped, header=%v", tbl.Header)
	}
}

func TestParseSkipsHashComments(t *testing.T) {
	in := "# note\nid,name\n1,Alice\n"
	tbl, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(tbl.Rows) != 1 || tbl.Header[0] != "id" {
		t.Fatalf("comment not skipped: header=%v rows=%v", tbl.Header, tbl.Rows)
	}
}

func TestParseRejectsExtraFields(t *testing.T) {
	in := "id,name\n1,Alice,extra\n"
	if _, err := Parse(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for extra fields")
	}
}

func TestParseEmpty(t *testing.T) {
	tbl, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if tbl == nil || tbl.Rows == nil {
		t.Fatal("empty input must return a non-nil table with empty rows")
	}
	if len(tbl.Rows) != 0 {
		t.Fatalf("rows=%v", tbl.Rows)
	}
}
