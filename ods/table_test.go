package ods

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/RHNTCH/go-ods/model"
)

func writeTestODS(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.ods")

	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	contentFile, err := zipWriter.Create("content.xml")
	if err != nil {
		t.Fatal(err)
	}

	_, err = contentFile.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestMakeTable(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>object_id</p></table-cell>
      <table-cell value-type="string"><p>object_tag</p></table-cell>
    </table-row>
    <table-row>
      <table-cell value-type="string"><p>1</p></table-cell>
      <table-cell value-type="string"><p>TT101</p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		t.Fatal(err)
	}

	cell, ok := table.CellByName(0, "object_tag")
	if !ok {
		t.Fatal("CellByName() ok = false, want true")
	}

	if cell.Formatted != "TT101" {
		t.Fatalf("CellByName() = %s, want TT101", cell.Formatted)
	}

	if table.Name != "AP" {
		t.Fatalf("Name = %s, want AP", table.Name)
	}

	if table.Width() != 2 {
		t.Fatalf("Width() = %d, want 2", table.Width())
	}

	if table.Height() != 1 {
		t.Fatalf("Height() = %d, want 1", table.Height())
	}
}

func TestMakeTableIgnoresEmptyColumnsAfterLastHeader(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>object_id</p></table-cell>
      <table-cell value-type="string"><p>object_tag</p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
    </table-row>
    <table-row>
      <table-cell value-type="string"><p>1</p></table-cell>
      <table-cell value-type="string"><p>TT101</p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		t.Fatal(err)
	}

	if table.Width() != 2 {
		t.Fatalf("Width() = %d, want 2", table.Width())
	}

	if len(table.Headers) != 2 {
		t.Fatalf("len(Headers) = %d, want 2", len(table.Headers))
	}

	if _, ok := table.Headers[2]; ok {
		t.Fatal("Headers[2] exists, want missing")
	}
}

func TestMakeTableKeepsEmptyColumnsInsideHeader(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>object_id</p></table-cell>
	  <table-cell value-type="string"><p></p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
      <table-cell value-type="string"><p>object_tag</p></table-cell>
      </table-row>
    <table-row>
      <table-cell value-type="string"><p>1</p></table-cell>
	  <table-cell value-type="string"><p></p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
      <table-cell value-type="string"><p>TT101</p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		t.Fatal(err)
	}

	if table.Width() != 4 {
		t.Fatalf("Width() = %d, want 4", table.Width())
	}

	if len(table.Headers) != 4 {
		t.Fatalf("len(Headers) = %d, want 4", len(table.Headers))
	}

	if table.Headers[1] != "" {
		t.Fatalf(`Headers[1] = %s, want ""`, table.Headers[1])
	}

	if table.Headers[2] != "" {
		t.Fatalf(`Headers[2] = %s, want ""`, table.Headers[2])
	}

	if len(table.HeaderIndexes) != 2 {
		t.Fatalf("len(HeaderIndexes) = %d, want 2", len(table.HeaderIndexes))
	}

	if _, ok := table.HeaderIndexes[""]; ok {
		t.Fatal(`HeaderIndexes[""] exists, want missing`)
	}
}

func TestMakeTablePadsShortRows(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>object_id</p></table-cell>
	  <table-cell value-type="string"><p>object_name</p></table-cell>
      <table-cell value-type="string"><p>object_tag</p></table-cell>
      </table-row>
    <table-row>
      <table-cell value-type="string"><p>1</p></table-cell>
	  <table-cell value-type="string"><p>temp</p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		t.Fatal(err)
	}

	if len(table.Rows[0].Cells) != len(table.Headers) {
		t.Fatalf("len(table.Rows[0].Cells) = %d, want %d", len(table.Rows[0].Cells), len(table.Headers))
	}

	if table.Rows[0].Cells[2].Formatted != "" {
		t.Fatalf(`table.Rows[0].Cells[2].Formatted = %s, want ""`, table.Rows[0].Cells[2].Formatted)
	}

	if len(table.Columns[2]) != 1 {
		t.Fatalf("len(Columns[2]) = %d, want 1", len(table.Columns[2]))
	}

	if table.Columns[2][0].Formatted != "" {
		t.Fatalf(`Columns[2][0].Formatted = %s, want ""`, table.Columns[2][0].Formatted)
	}

}

func TestMakeTableReturnsErrSheetNotFound(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	_, err = reader.MakeTable("APN")
	if !errors.Is(err, ErrSheetNotFound) {
		t.Fatalf("err = %v, want ErrSheetNotFound", err)
	}
}

func TestMakeTableReturnsErrEmptyHeader(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p></p></table-cell>
	  <table-cell value-type="string"><p></p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
      </table-row>
    <table-row>
      <table-cell value-type="string"><p>1</p></table-cell>
	  <table-cell value-type="string"><p>temp</p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	_, err = reader.MakeTable("AP")
	if !errors.Is(err, ErrEmptyHeader) {
		t.Fatalf("err = %v, want ErrEmptyHeader", err)
	}
}

func TestMakeTableKeepsDuplicateHeaders(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>object_id</p></table-cell>
	  <table-cell value-type="string"><p>object_id</p></table-cell>
      <table-cell value-type="string"><p></p></table-cell>
	  <table-cell value-type="string"><p></p></table-cell>
      </table-row>
    <table-row>
      <table-cell value-type="string"><p>1</p></table-cell>
	  <table-cell value-type="string"><p>2</p></table-cell>
	  <table-cell value-type="string"><p>3</p></table-cell>
	  <table-cell value-type="string"><p>4</p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		t.Fatal(err)
	}

	if len(table.HeaderIndexes) != 1 {
		t.Fatalf("len(table.HeaderIndexes) = %d, want 1", len(table.HeaderIndexes))
	}

	if len(table.Headers) != 2 {
		t.Fatalf("len(table.Headers) = %d, want 2", len(table.Headers))
	}

	if table.Headers[0] != "object_id" {
		t.Fatalf(`table.Headers[0] = %v, want "object_id"`, table.Headers[0])
	}

	if table.Headers[1] != "object_id" {
		t.Fatalf(`table.Headers[1] = %v, want "object_id"`, table.Headers[1])
	}

	if _, ok := table.Headers[2]; ok {
		t.Fatal("Headers[2] exists, want missing")
	}

	if _, ok := table.Headers[3]; ok {
		t.Fatal("Headers[3] exists, want missing")
	}

	indexes := table.HeaderIndexes["object_id"]
	if len(indexes) != 2 || indexes[0] != 0 || indexes[1] != 1 {
		t.Fatalf(`HeaderIndexes["object_id"] = %v, want [0 1]`, indexes)
	}

}

func TestMakeTableParsesTypedCells(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>name</p></table-cell>
      <table-cell value-type="string"><p>value</p></table-cell>
      <table-cell value-type="string"><p>enabled</p></table-cell>
      <table-cell value-type="string"><p>created_at</p></table-cell>
      <table-cell value-type="string"><p>formula</p></table-cell>
    </table-row>
    <table-row>
      <table-cell value-type="string"><p>TT101</p></table-cell>
      <table-cell value-type="float" value="12.5"><p>12.5</p></table-cell>
      <table-cell value-type="boolean" boolean-value="true"><p>TRUE</p></table-cell>
      <table-cell value-type="date" date-value="2026-06-05T00:00:00Z"><p>2026-06-05</p></table-cell>
      <table-cell value-type="float" value="25" formula="of:=1+24"><p>25</p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		t.Fatal(err)
	}

	name, ok := table.CellByName(0, "name")
	if !ok {
		t.Fatal("CellByName(name) ok = false, want true")
	}
	if name.Type != model.CellTypeString || name.String != "TT101" {
		t.Fatalf("name cell = %#v, want string TT101", name)
	}

	value, ok := table.CellByName(0, "value")
	if !ok {
		t.Fatal("CellByName(value) ok = false, want true")
	}
	if value.Type != model.CellTypeFloat || value.Number != 12.5 {
		t.Fatalf("value cell = %#v, want float 12.5", value)
	}

	enabled, ok := table.CellByName(0, "enabled")
	if !ok {
		t.Fatal("CellByName(enabled) ok = false, want true")
	}
	if enabled.Type != model.CellTypeBool || !enabled.Bool {
		t.Fatalf("enabled cell = %#v, want bool true", enabled)
	}

	createdAt, ok := table.CellByName(0, "created_at")
	if !ok {
		t.Fatal("CellByName(created_at) ok = false, want true")
	}
	wantTime := time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)
	if createdAt.Type != model.CellTypeDate || !createdAt.Time.Equal(wantTime) {
		t.Fatalf("created_at cell = %#v, want %v", createdAt, wantTime)
	}

	formula, ok := table.CellByName(0, "formula")
	if !ok {
		t.Fatal("CellByName(formula) ok = false, want true")
	}
	if formula.Formula != "of:=1+24" {
		t.Fatalf("Formula = %s, want of:=1+24", formula.Formula)
	}
}

func TestMakeTableExpandsRepeatedColumns(t *testing.T) {
	path := writeTestODS(t, `
<document>
  <table name="AP">
    <table-row>
      <table-cell value-type="string"><p>A</p></table-cell>
      <table-cell value-type="string"><p>B</p></table-cell>
      <table-cell value-type="string"><p>C</p></table-cell>
    </table-row>
    <table-row>
      <table-cell value-type="string" number-columns-repeated="3"><p>x</p></table-cell>
    </table-row>
  </table>
</document>
`)

	reader, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		t.Fatal(err)
	}

	if len(table.Rows[0].Cells) != 3 {
		t.Fatalf("len(table.Rows[0].Cells) = %d, want 3", len(table.Rows[0].Cells))
	}

	for i, cell := range table.Rows[0].Cells {
		if cell.Formatted != "x" {
			t.Fatalf("Rows[0].Cells[%d].Formatted = %s, want x", i, cell.Formatted)
		}
	}
}
