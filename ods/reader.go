package ods

import (
	"archive/zip"
	"encoding/xml"
	"io"
)

type Reader struct {
	decoder   *xml.Decoder
	zipReader *zip.ReadCloser
	content   io.ReadCloser
}

func (r *Reader) Sheets() *SheetIterator {
	return &SheetIterator{
		decoder: r.decoder,
	}
}
