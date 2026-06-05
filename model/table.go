package model

// Table is an in-memory sheet representation optimized for row and column access.
type Table struct {
	Name string

	Headers       map[int]string
	HeaderIndexes map[string][]int
	Rows          []Row
	Columns       map[int][]Cell
}
