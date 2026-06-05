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

// Sheets returns an iterator over sheets in the ODS file.
func (r *Reader) Sheets() *SheetIterator {
	return &SheetIterator{
		decoder: r.decoder,
	}
}
