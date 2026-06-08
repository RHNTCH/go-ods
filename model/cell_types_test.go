package model

import (
	"errors"
	"testing"
	"time"
)

func TestCellValue(t *testing.T) {
	date := time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		cell  Cell
		want  any
		isNil bool
	}{
		{
			name: "string",
			cell: Cell{Type: CellTypeString, String: "TT101"},
			want: "TT101",
		},
		{
			name: "float",
			cell: Cell{Type: CellTypeFloat, Number: 3.5},
			want: float64(3.5),
		},
		{
			name: "boolean",
			cell: Cell{Type: CellTypeBool, Bool: true},
			want: true,
		},
		{
			name: "date",
			cell: Cell{Type: CellTypeDate, Time: date},
			want: date,
		},
		{
			name:  "empty",
			cell:  Cell{Type: CellTypeEmpty},
			isNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cell.Value()
			if err != nil {
				t.Fatal(err)
			}

			if tt.isNil {
				if got != nil {
					t.Fatalf("Value() = %v, want nil", got)
				}
				return
			}

			if got != tt.want {
				t.Fatalf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellValueReturnsErrUndefinedCellType(t *testing.T) {
	tests := []struct {
		name string
		cell Cell
	}{
		{
			name: "unknown ODS type",
			cell: Cell{
				Type:      CellTypeUnknown,
				ValueType: "currency",
			},
		},
		{
			name: "invalid cell type",
			cell: Cell{
				Type: CellType(999),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.cell.Value()

			if value != nil {
				t.Fatalf("Value() = %v, want nil", value)
			}

			if !errors.Is(err, ErrUndefinedCellType) {
				t.Fatalf("err = %v, want ErrUndefinedCellType", err)
			}
		})
	}
}
