package model

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
