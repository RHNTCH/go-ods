package ods

import (
	"encoding/xml"

	"github.com/RHNTCH/go-ods/model"
)

func parseRow(decoder *xml.Decoder) model.Row {
	row := model.Row{}

	for {
		token, err := decoder.Token()
		if err != nil {
			return row
		}

		switch el := token.(type) {
		case xml.StartElement:
			if el.Name.Local == "table-cell" {
				repeat := getRepeatedColumns(el)
				cell := parseCell(decoder, el)

				for range repeat {
					row.Cells = append(row.Cells, cell)
				}
			}

		case xml.EndElement:
			if el.Name.Local == "table-row" {
				return row
			}
		}
	}
}
