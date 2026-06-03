package ods

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/RHNTCH/go-ods/model"
)

func parseCell(
	decoder *xml.Decoder,
	start xml.StartElement,
) model.Cell {
	cell := model.Cell{
		Type: model.CellTypeEmpty,
	}

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "value-type":
			switch attr.Value {
			case "string":
				cell.Type = model.CellTypeString
			case "float":
				cell.Type = model.CellTypeFloat
			case "boolean":
				cell.Type = model.CellTypeBool
			case "date":
				cell.Type = model.CellTypeDate
			}

		case "value":
			cell.Raw = attr.Value

		case "boolean-value":
			cell.Raw = attr.Value

		case "date-value":
			cell.Raw = attr.Value

		case "formula":
			cell.Formula = attr.Value
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			fillTypedValue(&cell)
			return cell
		}

		switch el := token.(type) {
		case xml.StartElement:
			if el.Name.Local == "p" {
				var content string

				err := decoder.DecodeElement(&content, &el)
				if err == nil {
					if cell.Formatted != "" {
						cell.Formatted += "\n"
					}

					cell.Formatted += content
				}
			}

		case xml.EndElement:
			if el.Name.Local == "table-cell" {
				fillTypedValue(&cell)
				return cell
			}
		}
	}
}

func fillTypedValue(cell *model.Cell) {
	switch cell.Type {
	case model.CellTypeString:
		cell.String = cell.Formatted

	case model.CellTypeFloat:
		v, err := strconv.ParseFloat(cell.Raw, 64)
		if err == nil {
			cell.Number = v
		}

	case model.CellTypeBool:
		v, err := strconv.ParseBool(cell.Raw)
		if err == nil {
			cell.Bool = v
		}

	case model.CellTypeDate:
		t, err := time.Parse(time.RFC3339, cell.Raw)
		if err == nil {
			cell.Time = t
		}
	}
}
