package main

import (
	"fmt"
	"log"

	"github.com/RHNTCH/go-ods/ods"
)

func main() {

	reader, err := ods.Open("testdata/СА MNS_PNS_RP_04_03_object_id.ods")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// ------------------------------------------------
	// Count rows
	// ------------------------------------------------

	count, err := reader.CountRows()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total rows:", count)

	reader.Close()

	reader, err = ods.Open("testdata/СА MNS_PNS_RP_04_03_object_id.ods")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// ------------------------------------------------
	// Custom callback processing
	// ------------------------------------------------
	AP_table, err := reader.MakeTable("Брдыщ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Names", AP_table.Name)
	fmt.Println()
	fmt.Println("Headers", AP_table.Headers)
	fmt.Println()
	fmt.Println("Header indexes", AP_table.HeaderIndexes)
	fmt.Println()
	fmt.Println("Val from A2", AP_table.Columns[0][0].Formatted)

}
