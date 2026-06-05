package ods

import (
	"fmt"

	"github.com/RHNTCH/go-ods/model"
)

// PrintRows prints all formatted row values to stdout.
func (r *Reader) PrintRows() error {
	return r.ForEachRow(func(sheet SheetInfo, row model.Row) error {
		fmt.Println("Sheet:", sheet.Name)

		for _, cell := range row.Cells {
			fmt.Print(cell.Formatted, " | ")
		}

		fmt.Println()
		return nil
	})
}
