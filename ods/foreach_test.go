package ods

import (
	"errors"
	"testing"

	"github.com/RHNTCH/go-ods/model"
)

const twoSheetsODS = `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>object_id</p></table-cell>
    </table-row>
    <table-row>
      <table-cell value-type="string"><p>1</p></table-cell>
    </table-row>
  </table>
  <table name="DI">
    <table-row>
      <table-cell value-type="string"><p>object_id</p></table-cell>
    </table-row>
    <table-row>
      <table-cell value-type="string"><p>2</p></table-cell>
    </table-row>
  </table>
</document>
`

func TestForEachSheet(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	var names []string

	err = reader.ForEachSheet(func(sheet *SheetCursor) error {
		names = append(names, sheet.Sheet.Name)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(names) != 2 {
		t.Fatalf("len(names) = %d, want 2", len(names))
	}

	if names[0] != "AP" || names[1] != "DI" {
		t.Fatalf("names = %v, want [AP DI]", names)
	}
}

func TestForEachSheetReturnsCallbackError(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	wantErr := errors.New("callback error")

	err = reader.ForEachSheet(func(sheet *SheetCursor) error {
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
}

func TestForSheets(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	var names []string

	err = reader.ForSheets([]string{"DI"}, func(sheet *SheetCursor) error {
		names = append(names, sheet.Sheet.Name)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(names) != 1 {
		t.Fatalf("len(names) = %d, want 1", len(names))
	}

	if names[0] != "DI" {
		t.Fatalf("names = %v, want [DI]", names)
	}

}

func TestForSheetsReturnsCallbackError(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	wantErr := errors.New("callback error")

	err = reader.ForSheets([]string{"DI"}, func(sheet *SheetCursor) error {
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
}

func TestForSheetsReturnsErrSheetNotFound(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	err = reader.ForSheets([]string{"APN"}, func(sheet *SheetCursor) error {
		t.Fatal("callback should not be called")
		return nil
	})

	if !errors.Is(err, ErrSheetNotFound) {
		t.Fatalf("err = %v, want ErrSheetNotFound", err)
	}
}

func TestForEachRow(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	var rows []string

	err = reader.ForEachRow(func(sheet SheetInfo, row model.Row) error {
		rows = append(rows, sheet.Name+":"+row.Cells[0].Formatted)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"AP:object_id",
		"AP:1",
		"DI:object_id",
		"DI:2",
	}

	if len(rows) != len(want) {
		t.Fatalf("len(rows) = %d, want %d", len(rows), len(want))
	}

	for i := range want {
		if rows[i] != want[i] {
			t.Fatalf("rows[%d] = %s, want %s", i, rows[i], want[i])
		}
	}
}

func TestForEachRowReturnsCallbackError(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	wantErr := errors.New("callback error")

	err = reader.ForEachRow(func(sheet SheetInfo, row model.Row) error {
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
}

func TestCountRows(t *testing.T) {
	path := writeTestODS(t, twoSheetsODS)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	count, err := reader.CountRows()
	if err != nil {
		t.Fatal(err)
	}

	if count != 4 {
		t.Fatalf("CountRows() = %d, want 4", count)
	}
}
