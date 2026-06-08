package model

// Value returns the typed value of a cell.
// It returns ErrUndefinedCellType for unsupported or invalid cell types.
func (c Cell) Value() (any, error) {
	switch c.Type {
	case CellTypeString:
		return c.String, nil

	case CellTypeFloat:
		return c.Number, nil

	case CellTypeBool:
		return c.Bool, nil

	case CellTypeDate:
		return c.Time, nil

	case CellTypeEmpty:
		return nil, nil

	case CellTypeUnknown:
		return nil, ErrUndefinedCellType

	default:
		return nil, ErrUndefinedCellType
	}
}
