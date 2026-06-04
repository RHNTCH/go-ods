package model

// Width returns width of the table by num of headers
func (t Table) Width() int {
	return len(t.Headers)
}

// Height returns height of table by num of rows including header
func (t Table) Height() int {
	return len(t.Rows)
}

// ColumnIndexes returns indexes of columns by it's header
func (t Table) ColumnIndexes(name string) ([]int, bool) {
	val, ok := t.HeaderIndexes[name]

	if !ok {
		return nil, false
	}

	return val, true
}

// ColumnByIndex returns columns by it's index
func (t Table) ColumnByIndex(index int) ([]Cell, bool) {
	column, ok := t.Columns[index]
	return column, ok
}

// ColumnByName returns column by it's header.
// Also returns nil, if column's name isn't unique
func (t Table) ColumnByName(name string) ([]Cell, bool) {
	indexes := t.HeaderIndexes[name]

	if len(indexes) != 1 {
		return nil, false
	}

	return t.ColumnByIndex(indexes[0])
}

// Cell returns type Cell by it's coordinates
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

// CellByName returns cell by header and row num
// if header is not unique CellByName returns false
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
