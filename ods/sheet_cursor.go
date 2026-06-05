package ods

import (
	"encoding/xml"

	"github.com/RHNTCH/go-ods/model"
)

// SheetCursor represents the current sheet and provides row iteration.
type SheetCursor struct {
	Sheet   *model.Sheet
	decoder *xml.Decoder
}
