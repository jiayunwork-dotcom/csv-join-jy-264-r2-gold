package write

import (
	"bufio"
	"encoding/csv"
	"io"

	"csv-join/internal/join"
)

// CSV writes header then rows. The writer is flushed before return.
func CSV(w io.Writer, res *join.Result) (err error) {
	bw := bufio.NewWriter(w)
	defer func() {
		if ferr := bw.Flush(); ferr != nil && err == nil {
			err = ferr
		}
	}()
	cw := csv.NewWriter(bw)
	if err := cw.Write(res.Header); err != nil {
		return err
	}
	for _, row := range res.Rows {
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}
	return nil
}
