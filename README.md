# go-ods

`go-ods` is a small Go library for reading `.ods` spreadsheets.

The project was created for automation tasks where signal lists and other
engineering data are stored in LibreOffice/OpenOffice Calc files and then need
to be converted into configuration files for other tools.

The library reads `content.xml` from the ODS archive and parses it with Go's
streaming XML decoder. This keeps memory usage low and makes it possible to
process rows one by one.

## Status

This project is experimental and focused on reading ODS files.

Currently supported:

- opening `.ods` files;
- iterating over sheets;
- iterating over rows;
- reading cell text, raw values, formulas, and basic typed values;
- preserving unsupported ODS value types for later handling;
- handling repeated columns;
- building an in-memory table from a sheet with `MakeTable`;
- simple callback helpers for common processing tasks.

Not implemented yet:

- writing ODS files;
- styles and formatting;
- merged cells;
- charts, images, and other embedded objects;
- full OpenDocument specification coverage.

## Installation

```bash
go get github.com/RHNTCH/go-ods
```

## Basic Usage

`MakeTable` is the easiest way to read a sheet when you want convenient
in-memory access by rows, columns, and header names.

```go
package main

import (
	"fmt"
	"log"

	"github.com/RHNTCH/go-ods/ods"
)

func main() {
	reader, err := ods.Open("signals.ods")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable("AP")
	if err != nil {
		log.Fatal(err)
	}

	tag, ok := table.CellByName(0, "object_tag")
	if !ok {
		log.Fatal("object_tag column not found")
	}

	fmt.Println(tag.Formatted)
}
```

## Iterator API

The lower-level API uses iterators. This style gives more control over the
parsing process.

```go
reader, err := ods.Open("signals.ods")
if err != nil {
	log.Fatal(err)
}
defer reader.Close()

sheets := reader.Sheets()
for sheets.Next() {
	sheet := sheets.Sheet()
	fmt.Println("Sheet:", sheet.Sheet.Name)

	rows := sheet.Rows()
	for rows.Next() {
		row := rows.Row()
		fmt.Println("Cells:", len(row.Cells))
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

if err := sheets.Err(); err != nil {
	log.Fatal(err)
}
```

## Callback Helpers

For one-pass processing, the library provides callback helpers. This API keeps
memory usage low and is a better fit for large files or generator pipelines.

```go
err := reader.ForEachSheet(func(sheet *ods.SheetCursor) error {
	fmt.Println("Sheet:", sheet.Sheet.Name)
	return nil
})
```

`ForEachRow` calls your function once for every parsed row:

```go
err := reader.ForEachRow(func(sheet ods.SheetInfo, row model.Row) error {
	// Process one row here.
	return nil
})
```

If the callback returns an error, iteration stops and the error is returned to
the caller.

To process a specific sheet, use `ForSheet`:

```go
err := reader.ForSheet("AP", func(sheet *ods.SheetCursor) error {
	rows := sheet.Rows()
	for rows.Next() {
		row := rows.Row()
		fmt.Println(len(row.Cells))
	}

	return rows.Err()
})
```

To process a set of sheets, use `ForSheets`:

```go
err := reader.ForSheets([]string{"AP", "DI"}, func(sheet *ods.SheetCursor) error {
	fmt.Println("Sheet:", sheet.Sheet.Name)
	return nil
})
```

## In-Memory Table API

`MakeTable` reads one sheet into memory and returns a `model.Table`.

```go
table, err := reader.MakeTable("AP")
if err != nil {
	log.Fatal(err)
}

fmt.Println(table.Headers)
fmt.Println(table.HeaderIndexes)
```

`model.Table` provides helper methods for common access patterns:

```go
fmt.Println(table.Width())
fmt.Println(table.Height())

column, ok := table.ColumnByName("object_id")
cell, ok := table.CellByName(0, "object_tag")
```

## Typed Cell Values

Cells preserve their displayed text, raw XML value, formula, and parsed type.
Use `Value` to retrieve a value using its Go type:

```go
value, err := cell.Value()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("%T: %v\n", value, value)
```

Currently supported typed values:

| ODS value type | Go value returned by `Cell.Value()` |
| --- | --- |
| `string` | `string` |
| `float` | `float64` |
| `boolean` | `bool` |
| `date` | `time.Time` |
| empty cell | `nil` |

Unsupported ODS value types are preserved as `CellTypeUnknown`. Their original
type remains available in `Cell.ValueType`, while `Cell.Value()` returns
`model.ErrUndefinedCellType`. The raw and formatted cell values remain
available.

This API assumes that the first row of the sheet is a header row. The logical
table width is determined by the last non-empty header cell. Empty header cells
inside that width are preserved in `Headers`, but they are not added to
`HeaderIndexes`.

This behavior helps ignore empty formatted cells that can appear in ODS files
when a sheet has styles applied far beyond the actual data range.

`MakeTable` stores both row-oriented and column-oriented views of the sheet.
This uses more memory than streaming iteration, but makes it convenient to read
data by row, column index, or column name. For large files or one-pass
processing, prefer the streaming and callback APIs.

## Streaming Notes

`go-ods` reads the XML stream from start to finish. This means a `Reader` is
single-pass: after rows or sheets have been read, the stream cannot be rewound.

If you need to process the same file again, open a new `Reader`.

```go
reader, err := ods.Open("signals.ods")
if err != nil {
	log.Fatal(err)
}
defer reader.Close()
```

## Example

See [examples/example1.go](examples/example1.go).

```bash
go run ./examples signals.ods AP
```

## Testing

```bash
go test ./...
go test ./model -cover
go test ./ods -cover
```

## Benchmarks

Run synthetic benchmarks:

```bash
go test ./ods -bench=. -benchmem
```

To benchmark a real file, pass the path and sheet name through environment
variables:

```bash
GO_ODS_BENCH_FILE="../testdata/signals.ods" \
GO_ODS_BENCH_SHEET="AP" \
go test ./ods -bench=RealFile -benchmem
```

## Roadmap

- tests for more OpenDocument edge cases;
- support for additional ODS value types;
- support for repeated rows;
- CLI tools and configuration generators built on top of the library.
