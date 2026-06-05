package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/RHNTCH/go-ods/ods"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: go run ./examples <file.ods> <sheet-name>")
	}

	path := os.Args[1]
	sheetName := os.Args[2]

	reader, err := ods.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	count, err := reader.CountRows()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total rows:", count)

	if err := reader.Close(); err != nil {
		log.Fatal(err)
	}

	reader, err = ods.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	table, err := reader.MakeTable(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sheet:", table.Name)
	fmt.Println("Width:", table.Width())
	fmt.Println("Height:", table.Height())
	fmt.Println("Headers:", table.Headers)

	if table.Height() > 0 {
		indexes := make([]int, 0, len(table.Headers))
		for index := range table.Headers {
			indexes = append(indexes, index)
		}
		sort.Ints(indexes)

		for _, index := range indexes {
			header := table.Headers[index]
			if header == "" {
				continue
			}

			cell, ok := table.Cell(0, index)
			if !ok {
				continue
			}

			fmt.Printf("First row %q: %s\n", header, cell.Formatted)
		}
	}
}
