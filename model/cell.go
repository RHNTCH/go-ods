package model

import "time"

const (
	CellTypeEmpty CellType = iota
	CellTypeString
	CellTypeFloat
	CellTypeBool
	CellTypeDate
	CellTypeUnknown
)

// Cell contains raw, formatted, typed, and formula data parsed from an ODS cell.
type Cell struct {
	Type CellType

	// typed values
	String string
	Number float64
	Bool   bool
	Time   time.Time

	// raw xml values
	ValueType string
	Raw       string

	// displayed text
	Formatted string

	// formulas
	Formula string
}
