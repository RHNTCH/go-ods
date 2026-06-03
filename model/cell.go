package model

import "time"

const (
	CellTypeEmpty CellType = iota
	CellTypeString
	CellTypeFloat
	CellTypeBool
	CellTypeDate
)

type Cell struct {
	Type CellType

	// typed values
	String string
	Number float64
	Bool   bool
	Time   time.Time

	// raw xml values
	Raw string

	// displayed text
	Formatted string

	// formulas
	Formula string
}
