package ods

import "github.com/RHNTCH/go-ods/model"

// CountRows returns the total number of rows across all sheets, including header rows.
func (r *Reader) CountRows() (int, error) {
	count := 0

	err := r.ForEachRow(func(sheet SheetInfo, row model.Row) error {
		count++
		return nil
	})

	return count, err
}
