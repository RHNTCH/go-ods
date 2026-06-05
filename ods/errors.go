package ods

import "errors"

var (
	// ErrEmptyHeader is returned when MakeTable cannot find a non-empty header cell.
	ErrEmptyHeader = errors.New("header is empty")
	// ErrSheetNotFound is returned when a requested sheet does not exist.
	ErrSheetNotFound = errors.New("sheet not found")
)
