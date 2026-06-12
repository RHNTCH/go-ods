package ods

import (
	"encoding/xml"
	"io"

	"github.com/RHNTCH/go-ods/model"
)

type RowIterator struct {
	decoder *xml.Decoder

	current     model.Row
	err         error
	mergedCells map[int]mergedCell
}

// Rows returns an iterator over rows in the current sheet.
func (s *SheetCursor) Rows() *RowIterator {
	return &RowIterator{
		decoder:     s.decoder,
		mergedCells: make(map[int]mergedCell),
	}
}

// Row returns the current row.
func (it *RowIterator) Row() model.Row {
	return it.current
}

// Err returns the first error encountered by the iterator.
func (it *RowIterator) Err() error {
	return it.err
}

// Next advances the iterator to the next row.
func (it *RowIterator) Next() bool {
	for {
		token, err := it.decoder.Token()

		if err == io.EOF {
			return false
		}

		if err != nil {
			it.err = err
			return false
		}

		switch el := token.(type) {
		case xml.StartElement:
			if el.Name.Local == "table-row" {
				it.current = parseRow(it.decoder, it.mergedCells)
				return true
			}

		case xml.EndElement:
			if el.Name.Local == "table" {
				return false
			}
		}
	}
}
