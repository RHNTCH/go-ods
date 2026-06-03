package ods

import (
	"archive/zip"
	"encoding/xml"
	"errors"
)

func Open(path string) (*Reader, error) {

	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	var contentFile *zip.File

	for _, f := range zr.File {

		if f.Name == "content.xml" {
			contentFile = f
			break
		}
	}

	if contentFile == nil {
		zr.Close()
		return nil, errors.New("content.xml not found")
	}

	rc, err := contentFile.Open()
	if err != nil {
		zr.Close()
		return nil, err
	}

	decoder := xml.NewDecoder(rc)

	return &Reader{
		decoder:   decoder,
		zipReader: zr,
		content:   rc,
	}, nil
}
