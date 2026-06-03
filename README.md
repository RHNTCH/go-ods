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

```go
package main

import (
	"fmt"
	"log"

	"github.com/RHNTCH/go-ods/model"
	"github.com/RHNTCH/go-ods/ods"
)

func main() {
	reader, err := ods.Open("testdata/testfile_MNS.ods")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	err = reader.ForEachRow(func(sheet ods.SheetInfo, row model.Row) error {
		fmt.Println("Sheet:", sheet.Name)

		for _, cell := range row.Cells {
			fmt.Print(cell.Formatted, " | ")
		}

		fmt.Println()
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
```

## Iterator API

The lower-level API uses iterators. This style gives more control over the
parsing process.

```go
reader, err := ods.Open("testdata/testfile_MNS.ods")
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

For simpler tasks, the library provides callback helpers.

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

This API assumes that the first row of the sheet is a header row. The logical
table width is determined by the last non-empty header cell. Empty header cells
inside that width are preserved in `Headers`, but they are not added to
`HeaderIndexes`.

This behavior helps ignore empty formatted cells that can appear in ODS files
when a sheet has styles applied far beyond the actual data range.

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
