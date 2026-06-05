package ods

import (
	"encoding/xml"
	"io"

	"github.com/RHNTCH/go-ods/model"
)

type SheetIterator struct {
	decoder *xml.Decoder

	current *SheetCursor
	err     error
}

// Sheet returns the current sheet cursor.
func (it *SheetIterator) Sheet() *SheetCursor {
	return it.current
}

// Err returns the first error encountered by the iterator.
func (it *SheetIterator) Err() error {
	return it.err
}

// Next advances the iterator to the next sheet.
func (it *SheetIterator) Next() bool {
	for {
		token, err := it.decoder.Token()

		if err == io.EOF {
			return false
		}

		if err != nil {
			it.err = err
			return false
		}

		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		if start.Name.Local != "table" {
			continue
		}

		meta := &model.Sheet{}

		for _, attr := range start.Attr {
			if attr.Name.Local == "name" {
				meta.Name = attr.Value
				break
			}
		}

		it.current = &SheetCursor{
			Sheet:   meta,
			decoder: it.decoder,
		}

		return true
	}
}
