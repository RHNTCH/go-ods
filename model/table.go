package model

type Table struct {
	Name string

	Headers       map[int]string
	HeaderIndexes map[string][]int
	Rows          []Row
	Columns       map[int][]Cell
}
