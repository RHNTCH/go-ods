package ods

import (
	"encoding/xml"
	"strconv"
)

func getRepeatedColumns(el xml.StartElement) int {
	return getPositiveIntAttribute(el, "number-columns-repeated")
}

func getSpannedColumns(el xml.StartElement) int {
	return getPositiveIntAttribute(el, "number-columns-spanned")
}

func getSpannedRows(el xml.StartElement) int {
	return getPositiveIntAttribute(el, "number-rows-spanned")
}

func getPositiveIntAttribute(el xml.StartElement, name string) int {
	for _, attr := range el.Attr {
		if attr.Name.Local == name {
			n, err := strconv.Atoi(attr.Value)
			if err == nil && n > 0 {
				return n
			}
		}
	}

	return 1
}
