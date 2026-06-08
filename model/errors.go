package model

import "errors"

// ErrUndefinedCellType is returned when a cell has an unsupported or invalid type.
var ErrUndefinedCellType = errors.New("undefined cell type")
