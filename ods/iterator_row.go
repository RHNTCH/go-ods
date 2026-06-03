package ods

import (
	"encoding/xml"
	"io"

	"github.com/RHNTCH/go-ods/model"
)

type RowIterator struct {
	decoder *xml.Decoder

	current model.Row
	err     error
}

func (s *SheetCursor) Rows() *RowIterator {
	return &RowIterator{
		decoder: s.decoder,
	}
}

func (it *RowIterator) Row() model.Row {
	return it.current
}

func (it *RowIterator) Err() error {
	return it.err
}

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
				it.current = parseRow(it.decoder)
				return true
			}

		case xml.EndElement:
			if el.Name.Local == "table" {
				return false
			}
		}
	}
}
