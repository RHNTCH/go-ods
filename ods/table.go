package ods

import "github.com/RHNTCH/go-ods/model"

// MakeTable reads a sheet into memory using the first row as headers.
func (r *Reader) MakeTable(sheetName string) (model.Table, error) {
	table := model.Table{
		Name: sheetName,

		Headers:       make(map[int]string),
		HeaderIndexes: make(map[string][]int),
		Columns:       make(map[int][]model.Cell),
	}

	err := r.ForSheet(sheetName, func(sheet *SheetCursor) error {
		table.Name = sheet.Sheet.Name

		rows := sheet.Rows()
		rowIndex := 0
		width := 0

		for rows.Next() {
			row := rows.Row()
			lastHeaderIndex := -1

			if rowIndex == 0 {
				for i, cell := range row.Cells {
					if cell.Formatted != "" {
						lastHeaderIndex = i
					}
				}

				if lastHeaderIndex == -1 {
					return ErrEmptyHeader
				}

				width = lastHeaderIndex + 1

				for i := 0; i < width; i++ {
					header := row.Cells[i].Formatted
					table.Headers[i] = header
					if header != "" {
						table.HeaderIndexes[header] = append(table.HeaderIndexes[header], i)
					}
				}
			} else {
				trimmed := model.Row{
					Cells: make([]model.Cell, 0, width),
				}

				for i := 0; i < width; i++ {
					cell := model.Cell{}
					if i < len(row.Cells) {
						cell = row.Cells[i]
					}

					trimmed.Cells = append(trimmed.Cells, cell)
					table.Columns[i] = append(table.Columns[i], cell)
				}

				table.Rows = append(table.Rows, trimmed)
			}
			rowIndex++
		}

		return rows.Err()
	})

	if err != nil {
		return table, err
	}

	return table, nil
}
