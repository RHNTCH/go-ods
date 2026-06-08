package model

import "fmt"

// Width returns the number of logical table columns.
func (t Table) Width() int {
	return len(t.Headers)
}

// Height returns the number of data rows in the table.
func (t Table) Height() int {
	return len(t.Rows)
}

// ColumnIndexes returns all column indexes matching a header name.
func (t Table) ColumnIndexes(name string) ([]int, bool) {
	val, ok := t.HeaderIndexes[name]

	if !ok {
		return nil, false
	}

	return val, true
}

// ColumnByIndex returns a column by its zero-based index.
func (t Table) ColumnByIndex(index int) ([]Cell, bool) {
	column, ok := t.Columns[index]
	return column, ok
}

// ColumnByName returns a column by header name.
// It returns false if the header is missing or not unique.
func (t Table) ColumnByName(name string) ([]Cell, bool) {
	indexes := t.HeaderIndexes[name]

	if len(indexes) != 1 {
		return nil, false
	}

	return t.ColumnByIndex(indexes[0])
}

// Cell returns a cell by zero-based row and column indexes.
func (t Table) Cell(rowIndex, columnIndex int) (Cell, bool) {

	if rowIndex < 0 || rowIndex >= len(t.Rows) {
		return Cell{}, false
	}

	row := t.Rows[rowIndex]

	if columnIndex < 0 || columnIndex >= len(row.Cells) {
		return Cell{}, false
	}

	return row.Cells[columnIndex], true
}

// CellByName returns a cell by row index and header name.
// It returns false if the header is missing, not unique, or the row is out of range.
func (t Table) CellByName(rowIndex int, columnName string) (Cell, bool) {
	columnIndexes, ok := t.ColumnIndexes(columnName)

	if !ok {
		return Cell{}, false
	}

	if len(columnIndexes) != 1 {
		return Cell{}, false
	}

	return t.Cell(rowIndex, columnIndexes[0])
}

// RequireColumns verifies that each named column exists and is unique.
func (t Table) RequireColumns(names ...string) error {
	for _, name := range names {
		indexes, ok := t.HeaderIndexes[name]

		if !ok {
			return fmt.Errorf(
				"sheet %s: %w: %s",
				t.Name,
				ErrColumnNotFound,
				name,
			)
		}

		if len(indexes) != 1 {
			return fmt.Errorf(
				"sheet %s: %w: %s",
				t.Name,
				ErrColumnNotUnique,
				name,
			)
		}
	}

	return nil
}

// RowMap returns a row mapped by its non-empty, unique header names.
func (t Table) RowMap(rowIndex int) (map[string]Cell, error) {
	if rowIndex < 0 || rowIndex >= len(t.Rows) {
		return nil, ErrRowOutOfRange
	}

	rowMap := make(map[string]Cell, len(t.Headers))

	for i, cell := range t.Rows[rowIndex].Cells {
		header := t.Headers[i]
		if header == "" {
			continue
		}

		if err := t.RequireColumns(header); err != nil {
			return nil, err
		}

		rowMap[header] = cell
	}

	return rowMap, nil
}

// ValueMap returns a row mapped by header names to typed cell values.
func (t Table) ValueMap(rowIndex int) (map[string]any, error) {
	rowMap, err := t.RowMap(rowIndex)
	if err != nil {
		return nil, err
	}

	valueMap := make(map[string]any, len(rowMap))

	for column, value := range rowMap {
		valueMap[column], err = value.Value()
		if err != nil {
			return nil, err
		}
	}

	return valueMap, nil
}
