package ods

import (
	"encoding/xml"
	"strconv"
)

func getRepeatedColumns(el xml.StartElement) int {
	repeat := 1

	for _, attr := range el.Attr {
		if attr.Name.Local == "number-columns-repeated" {
			n, err := strconv.Atoi(attr.Value)
			if err == nil && n > 0 {
				repeat = n
			}
		}
	}

	return repeat
}

func getRepeatedRows(el xml.StartElement) int {
	repeat := 1

	for _, attr := range el.Attr {
		if attr.Name.Local == "number-rows-repeated" {
			n, err := strconv.Atoi(attr.Value)
			if err == nil && n > 0 {
				repeat = n
			}
		}
	}

	return repeat
}
