package model

import (
	"errors"
	"reflect"
	"testing"
	"time"
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

func TestRequireColumns(t *testing.T) {
	table := Table{
		Name: "AP",
		HeaderIndexes: map[string][]int{
			"object_id":  {0},
			"object_tag": {1},
			"duplicate":  {2, 3},
		},
	}

	t.Run("all columns exist and are unique", func(t *testing.T) {
		err := table.RequireColumns("object_id", "object_tag")
		if err != nil {
			t.Fatalf("RequireColumns() err = %v, want nil", err)
		}
	})

	t.Run("empty requirements", func(t *testing.T) {
		err := table.RequireColumns()
		if err != nil {
			t.Fatalf("RequireColumns() err = %v, want nil", err)
		}
	})

	t.Run("missing column", func(t *testing.T) {
		err := table.RequireColumns("missing")
		if !errors.Is(err, ErrColumnNotFound) {
			t.Fatalf("err = %v, want ErrColumnNotFound", err)
		}
	})

	t.Run("duplicate column", func(t *testing.T) {
		err := table.RequireColumns("duplicate")
		if !errors.Is(err, ErrColumnNotUnique) {
			t.Fatalf("err = %v, want ErrColumnNotUnique", err)
		}
	})
}

func TestRowMap(t *testing.T) {
	table := Table{
		Name: "AP",
		Headers: map[int]string{
			0: "object_id",
			1: "",
			2: "object_tag",
		},
		HeaderIndexes: map[string][]int{
			"object_id":  {0},
			"object_tag": {2},
		},
		Rows: []Row{
			{
				Cells: []Cell{
					{Type: CellTypeFloat, Number: 1},
					{Type: CellTypeEmpty},
					{Type: CellTypeString, String: "TT101"},
				},
			},
		},
	}

	row, err := table.RowMap(0)
	if err != nil {
		t.Fatal(err)
	}

	if len(row) != 2 {
		t.Fatalf("len(RowMap()) = %d, want 2", len(row))
	}

	if row["object_id"].Number != 1 {
		t.Fatalf(`RowMap()["object_id"].Number = %v, want 1`, row["object_id"].Number)
	}

	if row["object_tag"].String != "TT101" {
		t.Fatalf(`RowMap()["object_tag"].String = %s, want TT101`, row["object_tag"].String)
	}

	if _, ok := row[""]; ok {
		t.Fatal(`RowMap()[""] exists, want missing`)
	}
}

func TestRowMapReturnsErrRowOutOfRange(t *testing.T) {
	table := Table{
		Rows: []Row{{}},
	}

	for _, rowIndex := range []int{-1, 1} {
		_, err := table.RowMap(rowIndex)
		if !errors.Is(err, ErrRowOutOfRange) {
			t.Fatalf("RowMap(%d) err = %v, want ErrRowOutOfRange", rowIndex, err)
		}
	}
}

func TestRowMapReturnsColumnErrors(t *testing.T) {
	t.Run("missing header index", func(t *testing.T) {
		table := Table{
			Name: "AP",
			Headers: map[int]string{
				0: "object_id",
			},
			HeaderIndexes: map[string][]int{},
			Rows: []Row{
				{Cells: []Cell{{Type: CellTypeFloat, Number: 1}}},
			},
		}

		_, err := table.RowMap(0)
		if !errors.Is(err, ErrColumnNotFound) {
			t.Fatalf("err = %v, want ErrColumnNotFound", err)
		}
	})

	t.Run("duplicate header", func(t *testing.T) {
		table := Table{
			Name: "AP",
			Headers: map[int]string{
				0: "object_id",
				1: "object_id",
			},
			HeaderIndexes: map[string][]int{
				"object_id": {0, 1},
			},
			Rows: []Row{
				{Cells: []Cell{{}, {}}},
			},
		}

		_, err := table.RowMap(0)
		if !errors.Is(err, ErrColumnNotUnique) {
			t.Fatalf("err = %v, want ErrColumnNotUnique", err)
		}
	})
}

func TestValueMap(t *testing.T) {
	date := time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC)
	table := Table{
		Headers: map[int]string{
			0: "object_tag",
			1: "value",
			2: "enabled",
			3: "updated_at",
			4: "description",
			5: "",
		},
		HeaderIndexes: map[string][]int{
			"object_tag":  {0},
			"value":       {1},
			"enabled":     {2},
			"updated_at":  {3},
			"description": {4},
		},
		Rows: []Row{
			{Cells: []Cell{
				{Type: CellTypeString, String: "TT101"},
				{Type: CellTypeFloat, Number: 3.5},
				{Type: CellTypeBool, Bool: true},
				{Type: CellTypeDate, Time: date},
				{Type: CellTypeEmpty},
				{Type: CellTypeString, String: "ignored"},
			}},
		},
	}

	got, err := table.ValueMap(0)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]any{
		"object_tag":  "TT101",
		"value":       3.5,
		"enabled":     true,
		"updated_at":  date,
		"description": nil,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ValueMap() = %#v, want %#v", got, want)
	}
}

func TestValueMapReturnsRowMapErrors(t *testing.T) {
	tests := []struct {
		name    string
		table   Table
		row     int
		wantErr error
	}{
		{
			name:    "row out of range",
			table:   Table{Rows: []Row{{}}},
			row:     1,
			wantErr: ErrRowOutOfRange,
		},
		{
			name: "duplicate header",
			table: Table{
				Headers:       map[int]string{0: "object_tag"},
				HeaderIndexes: map[string][]int{"object_tag": {0, 1}},
				Rows:          []Row{{Cells: []Cell{{Type: CellTypeString}}}},
			},
			wantErr: ErrColumnNotUnique,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.table.ValueMap(tt.row)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValueMapReturnsErrUndefinedCellType(t *testing.T) {
	table := Table{
		Headers:       map[int]string{0: "price"},
		HeaderIndexes: map[string][]int{"price": {0}},
		Rows: []Row{
			{Cells: []Cell{{Type: CellTypeUnknown, ValueType: "currency"}}},
		},
	}

	values, err := table.ValueMap(0)
	if values != nil {
		t.Fatalf("ValueMap() = %#v, want nil", values)
	}
	if !errors.Is(err, ErrUndefinedCellType) {
		t.Fatalf("err = %v, want ErrUndefinedCellType", err)
	}
}
