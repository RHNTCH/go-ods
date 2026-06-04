package ods

import "errors"

var (
	ErrEmptyHeader   = errors.New("header is empty")
	ErrSheetNotFound = errors.New("sheet not found")
)
