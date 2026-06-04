package model

import (
	"testing"
)

func TestTableWidth(t *testing.T) {
	table := Table{
		Headers: map[int]string{
			0: "A",
			1: "B",
			2: "C",
		},
	}
	got := table.Width()
	want := 3

	if got != want {
		t.Fatalf("Width() = %d, want %d", got, want)
	}

}

func TestTableHeight(t *testing.T) {
	table := Table{
		Rows: []Row{
			{Cells: []Cell{{Formatted: "row1"}}},
			{Cells: []Cell{{Formatted: "row2"}}},
		},
	}

	got := table.Height()
	want := 2

	if got != want {
		t.Fatalf("Height() = %d, want %d", got, want)
	}
}

func TestColumnByIndex(t *testing.T) {
	table := Table{
		Columns: map[int][]Cell{
			0: {{Formatted: "A1"}, {Formatted: "A2"}},
		},
	}

	got, ok := table.ColumnByIndex(0)

	if !ok {
		t.Fatal("ColumnByIndex() ok = false, want true")
	}

	if got[0].Formatted != "A1" {
		t.Fatalf("ColumnByIndex() = %s, want A1", got[0].Formatted)
	}

	if len(got) != 2 {
		t.Fatalf("len(ColumnByIndex()) = %d, want 2", len(got))
	}

	_, ok = table.ColumnByIndex(99)
	if ok {
		t.Fatal("ColumnByIndex() ok = true, want false")
	}

}
func TestColumnIndexes(t *testing.T) {
	table := Table{
		HeaderIndexes: map[string][]int{
			"object_name": {0},
			"object_tag":  {1, 3},
		},
	}

	t.Run("existing unique header", func(t *testing.T) {
		indexes, ok := table.ColumnIndexes("object_name")
		if !ok {
			t.Fatal("ok = false, want true")
		}
		if len(indexes) != 1 || indexes[0] != 0 {
			t.Fatalf("indexes = %v, want [0]", indexes)
		}
	})

	t.Run("existing duplicate header", func(t *testing.T) {
		indexes, ok := table.ColumnIndexes("object_tag")
		if !ok {
			t.Fatal("ok = false, want true")
		}
		if len(indexes) != 2 || indexes[0] != 1 || indexes[1] != 3 {
			t.Fatalf("indexes = %v, want [1 3]", indexes)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		_, ok := table.ColumnIndexes("free_col")
		if ok {
			t.Fatal("ok = true, want false")
		}
	})
}

func TestColumnByName(t *testing.T) {
	table := Table{
		Headers: map[int]string{
			0: "object_name",
			1: "object_tag",
			2: "object_name",
		},
		HeaderIndexes: map[string][]int{
			"object_name": {0, 2},
			"object_tag":  {1},
		},
		Columns: map[int][]Cell{
			1: {{Formatted: "TT101"}, {Formatted: "TT102"}},
		},
	}

	_, ok := table.ColumnByName("object_name")

	if ok {
		t.Fatal("ColumnByName() ok = true, want false")
	}

	_, ok = table.ColumnByName("object_id")

	if ok {
		t.Fatal("ColumnByName() ok = true, want false")
	}

	column, ok := table.ColumnByName("object_tag")

	if !ok {
		t.Fatal("ok = false, want true")
	}

	if len(column) != 2 {
		t.Fatalf("len(ColumnByName()) = %d, want 2", len(column))
	}

	if column[0].Formatted != "TT101" {
		t.Fatalf("ColumnByName()[0] = %s, want TT101", column[0].Formatted)
	}
}

func TestCell(t *testing.T) {
	// Table has 2 data rows and 3 columns.
	table := Table{
		Headers: map[int]string{
			0: "object_name",
			1: "object_tag",
			2: "object_name",
		},
		HeaderIndexes: map[string][]int{
			"object_name": {0, 2},
			"object_tag":  {1},
		},
		Columns: map[int][]Cell{
			0: {{Formatted: "Temp1"}, {Formatted: "Temp2"}},
			1: {{Formatted: "TT101"}, {Formatted: "TT102"}},
			2: {{Formatted: "Temp1"}, {Formatted: "Temp2"}},
		},
		Rows: []Row{
			{
				Cells: []Cell{
					{Formatted: "Temp1"},
					{Formatted: "TT101"},
					{Formatted: "Temp1"},
				},
			},
			{
				Cells: []Cell{
					{Formatted: "Temp2"},
					{Formatted: "TT102"},
					{Formatted: "Temp2"},
				},
			},
		},
	}

	t.Run("regular cell", func(t *testing.T) {
		cell, ok := table.Cell(0, 1)
		if !ok {
			t.Fatal("ok = false, want true")
		}

		if cell.Formatted != "TT101" {
			t.Fatalf("Cell() = %s, want TT101", cell.Formatted)
		}
	})

	t.Run("non-existing cell", func(t *testing.T) {
		_, ok := table.Cell(0, 3)

		if ok {
			t.Fatal("ok = true, want false")
		}

		_, ok = table.Cell(2, 0)

		if ok {
			t.Fatal("ok = true, want false")
		}

		_, ok = table.Cell(2, 3)

		if ok {
			t.Fatal("ok = true, want false")
		}

		_, ok = table.Cell(-1, 0)

		if ok {
			t.Fatal("ok = true, want false")
		}

		_, ok = table.Cell(0, -1)

		if ok {
			t.Fatal("ok = true, want false")
		}
	})

}

func TestCellByName(t *testing.T) {
	table := Table{
		Headers: map[int]string{
			0: "object_name",
			1: "object_tag",
			2: "object_name",
		},
		HeaderIndexes: map[string][]int{
			"object_name": {0, 2},
			"object_tag":  {1},
		},
		Columns: map[int][]Cell{
			0: {{Formatted: "Temp1"}, {Formatted: "Temp2"}},
			1: {{Formatted: "TT101"}, {Formatted: "TT102"}},
			2: {{Formatted: "Temp1"}, {Formatted: "Temp2"}},
		},
		Rows: []Row{
			{
				Cells: []Cell{
					{Formatted: "Temp1"},
					{Formatted: "TT101"},
					{Formatted: "Temp1"},
				},
			},
			{
				Cells: []Cell{
					{Formatted: "Temp2"},
					{Formatted: "TT102"},
					{Formatted: "Temp2"},
				},
			},
		},
	}

	t.Run("non-unique column name", func(t *testing.T) {
		_, ok := table.CellByName(0, "object_name")

		if ok {
			t.Fatal("ok = true, want false")
		}
	})

	t.Run("unique column name", func(t *testing.T) {
		cell, ok := table.CellByName(0, "object_tag")

		if !ok {
			t.Fatal("ok = false, want true")
		}

		if cell.Formatted != "TT101" {
			t.Fatalf("Cell() = %s, want TT101", cell.Formatted)
		}
	})

	t.Run("invalid index", func(t *testing.T) {
		_, ok := table.CellByName(-1, "object_tag")

		if ok {
			t.Fatal("ok = true, want false")
		}
	})

	t.Run("invalid column name", func(t *testing.T) {
		_, ok := table.CellByName(0, "object_id")

		if ok {
			t.Fatal("ok = true, want false")
		}
	})
}
