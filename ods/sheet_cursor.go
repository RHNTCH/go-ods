package ods

import (
	"encoding/xml"

	"github.com/RHNTCH/go-ods/model"
)

type SheetCursor struct {
	Sheet   *model.Sheet
	decoder *xml.Decoder
}
