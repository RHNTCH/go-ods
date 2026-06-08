package model

import "errors"

// ErrUndefinedCellType is returned when a cell has an unsupported or invalid type.
var ErrUndefinedCellType = errors.New("undefined cell type")
var ErrColumnNotFound = errors.New("column not found")
var ErrColumnNotUnique = errors.New("column not unique")
var ErrRowOutOfRange = errors.New("row index out of range")
