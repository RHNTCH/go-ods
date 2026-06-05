package ods

import (
	"os"
	"testing"

	"github.com/RHNTCH/go-ods/model"
)

const benchmarkODSContent = `
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
`

func BenchmarkMakeTable(b *testing.B) {
	benchmarkMakeTable(b, writeTestODS(b, benchmarkODSContent), "AP")
}

func BenchmarkForSheetRows(b *testing.B) {
	benchmarkForSheetRows(b, writeTestODS(b, benchmarkODSContent), "AP")
}

func BenchmarkForEachRow(b *testing.B) {
	benchmarkForEachRow(b, writeTestODS(b, benchmarkODSContent))
}

func BenchmarkMakeTableRealFile(b *testing.B) {
	path := os.Getenv("GO_ODS_BENCH_FILE")
	if path == "" {
		b.Skip("set GO_ODS_BENCH_FILE to benchmark a real ODS file")
	}

	sheetName := os.Getenv("GO_ODS_BENCH_SHEET")
	if sheetName == "" {
		sheetName = "AP"
	}

	benchmarkMakeTable(b, path, sheetName)
}

func BenchmarkForSheetRowsRealFile(b *testing.B) {
	path := os.Getenv("GO_ODS_BENCH_FILE")
	if path == "" {
		b.Skip("set GO_ODS_BENCH_FILE to benchmark a real ODS file")
	}

	sheetName := os.Getenv("GO_ODS_BENCH_SHEET")
	if sheetName == "" {
		sheetName = "AP"
	}

	benchmarkForSheetRows(b, path, sheetName)
}

func BenchmarkForEachRowRealFile(b *testing.B) {
	path := os.Getenv("GO_ODS_BENCH_FILE")
	if path == "" {
		b.Skip("set GO_ODS_BENCH_FILE to benchmark a real ODS file")
	}

	benchmarkForEachRow(b, path)
}

func benchmarkMakeTable(b *testing.B, path, sheetName string) {
	b.Helper()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader, err := Open(path)
		if err != nil {
			b.Fatal(err)
		}

		_, err = reader.MakeTable(sheetName)
		if err != nil {
			b.Fatal(err)
		}

		if err := reader.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkForSheetRows(b *testing.B, path, sheetName string) {
	b.Helper()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader, err := Open(path)
		if err != nil {
			b.Fatal(err)
		}

		err = reader.ForSheet(sheetName, func(sheet *SheetCursor) error {
			rows := sheet.Rows()

			for rows.Next() {
				_ = rows.Row()
			}

			return rows.Err()
		})
		if err != nil {
			b.Fatal(err)
		}

		if err := reader.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkForEachRow(b *testing.B, path string) {
	b.Helper()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader, err := Open(path)
		if err != nil {
			b.Fatal(err)
		}

		err = reader.ForEachRow(func(sheet SheetInfo, row model.Row) error {
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}

		if err := reader.Close(); err != nil {
			b.Fatal(err)
		}
	}
}
