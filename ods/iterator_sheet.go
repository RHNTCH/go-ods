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

func (it *SheetIterator) Sheet() *SheetCursor {
	return it.current
}

func (it *SheetIterator) Err() error {
	return it.err
}

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
