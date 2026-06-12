package ods

import (
	"encoding/xml"

	"github.com/RHNTCH/go-ods/model"
)

type mergedCell struct {
	cell          model.Cell
	remainingRows int
}

func parseRow(decoder *xml.Decoder, mergedCells map[int]mergedCell) model.Row {
	row := model.Row{}
	columnIndex := 0
	coveredByCurrentRow := 0

	for {
		token, err := decoder.Token()
		if err != nil {
			return row
		}

		switch el := token.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "table-cell":
				repeat := getRepeatedColumns(el)
				columnSpan := getSpannedColumns(el)
				rowSpan := getSpannedRows(el)
				cell := parseCell(decoder, el)

				for range repeat {
					for range columnSpan {
						row.Cells = append(row.Cells, cell)

						if rowSpan > 1 {
							mergedCells[columnIndex] = mergedCell{
								cell:          cell,
								remainingRows: rowSpan - 1,
							}
						}

						columnIndex++
					}
				}

				coveredByCurrentRow += repeat * (columnSpan - 1)

			case "covered-table-cell":
				repeat := getRepeatedColumns(el)

				for range repeat {
					if coveredByCurrentRow > 0 {
						coveredByCurrentRow--
						continue
					}

					cell := model.Cell{Type: model.CellTypeEmpty}
					if merged, ok := mergedCells[columnIndex]; ok {
						cell = merged.cell
						merged.remainingRows--

						if merged.remainingRows == 0 {
							delete(mergedCells, columnIndex)
						} else {
							mergedCells[columnIndex] = merged
						}
					}

					row.Cells = append(row.Cells, cell)
					columnIndex++
				}
			}

		case xml.EndElement:
			if el.Name.Local == "table-row" {
				return row
			}
		}
	}
}
