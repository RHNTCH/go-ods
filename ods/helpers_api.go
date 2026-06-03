package ods

import (
	"errors"
	"fmt"

	"github.com/RHNTCH/go-ods/model"
)

type SheetInfo struct {
	Name string
}

func (r *Reader) ForEachRow(
	fn func(sheet SheetInfo, row model.Row) error,
) error {
	return r.ForEachSheet(func(sheet *SheetCursor) error {
		rows := sheet.Rows()

		for rows.Next() {
			if err := fn(
				SheetInfo{Name: sheet.Sheet.Name},
				rows.Row(),
			); err != nil {
				return err
			}
		}

		return rows.Err()
	})
}

func (r *Reader) ForEachSheet(fn func(sheet *SheetCursor) error) error {
	sheets := r.Sheets()

	for sheets.Next() {
		if err := fn(sheets.Sheet()); err != nil {
			return err
		}
	}

	return sheets.Err()
}

func (r *Reader) ForSomeSheet(sheetNames []string, fn func(sheet *SheetCursor) error) error {
	targetSheets := make(map[string]bool)
	for _, name := range sheetNames {
		targetSheets[name] = true
	}

	sheets := r.Sheets()

	for sheets.Next() {
		sheet := sheets.Sheet()

		if targetSheets[sheet.Sheet.Name] {
			if err := fn(sheets.Sheet()); err != nil {
				return err
			}
		}

	}
	return sheets.Err()
}

func (r *Reader) CountRows() (int, error) {
	count := 0

	err := r.ForEachRow(func(sheet SheetInfo, row model.Row) error {
		count++
		return nil
	})

	return count, err
}

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

func (r *Reader) MakeTable(sheetName string) (model.Table, error) {

	table := model.Table{
		Name: sheetName,

		Headers:       make(map[int]string),
		HeaderIndexes: make(map[string][]int),
		Columns:       make(map[int][]model.Cell),
	}

	found := false

	err := r.ForSomeSheet([]string{sheetName}, func(sheet *SheetCursor) error {
		found = true
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
					return errors.New("header is empty")
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
				trimmed := model.Row{}

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

	if !found {
		return table, errors.New("sheet not found")
	}

	return table, nil
}

// // Вернуть целую строку как слайс ячеек. Вспомогательная функция, может быть стоит сделать неэкспеортируемой
// func (r *Reader) GetRows(sheetName string) ([]model.Cell, error)

// // Вернуть значения из определенных колонок по имени как слайс ячеек
// func (r *Reader) GetColumnsByName(column ...string) ([]model.Cell, error)

// // Вернуть значения из определенных колонок по номеру как слайс ячеек
// func (r *Reader) GetColumnsByIndex(colNum ...string) ([]model.Cell, error)

// // Вернуть заголовки текущего листа. Вспомогательная функция для GetColumnsByName
// func (r *Reader) getHeader() ([]model.Cell, error)
