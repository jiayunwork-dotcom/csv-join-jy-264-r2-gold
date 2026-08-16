// Command csv-join joins two CSV files on a shared key column.
package main

import (
	"flag"
	"fmt"
	"os"

	"csv-join/internal/join"
	"csv-join/internal/table"
	"csv-join/internal/write"
)

func main() {
	leftPath := flag.String("left", "", "left CSV path (required)")
	rightPath := flag.String("right", "", "right CSV path (required)")
	key := flag.String("key", "id", "join column name present in both files")
	how := flag.String("how", "inner", "join type: inner or left")
	out := flag.String("out", "-", "output path, or - for stdout")
	flag.Parse()

	if *leftPath == "" || *rightPath == "" {
		fatal("missing required -left and -right")
	}

	left, err := load(*leftPath)
	if err != nil {
		fatal("left: %v", err)
	}
	right, err := load(*rightPath)
	if err != nil {
		fatal("right: %v", err)
	}

	var res *join.Result
	switch *how {
	case "inner":
		res, err = join.Inner(left, right, *key)
	case "left":
		res, err = join.Left(left, right, *key)
	default:
		fatal("unknown -how %q (want inner or left)", *how)
	}
	if err != nil {
		fatal("%v", err)
	}

	w := os.Stdout
	if *out != "-" {
		f, err := os.Create(*out)
		if err != nil {
			fatal("create %s: %v", *out, err)
		}
		defer f.Close()
		w = f
	}
	if err := write.CSV(w, res); err != nil {
		fatal("write: %v", err)
	}
}

func load(path string) (*table.Table, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return table.Parse(f)
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "csv-join: "+format+"\n", a...)
	os.Exit(1)
}
