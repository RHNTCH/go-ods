package ods

import "github.com/RHNTCH/go-ods/model"

// ForEachRow calls fn for every row in every sheet.
func (r *Reader) ForEachRow(
	fn func(sheet SheetInfo, row model.Row) error,
) error {
	return r.ForEachSheet(func(sheet *SheetCursor) error {
		rows := sheet.Rows()

		for rows.Next() {
			if err := fn(
				SheetInfo{Name: sheet.Sheet.Name},
				rows.Row(),
			); err != nil {
				return err
			}
		}

		return rows.Err()
	})
}

// ForEachSheet calls fn for every sheet in the file.
func (r *Reader) ForEachSheet(fn func(sheet *SheetCursor) error) error {
	sheets := r.Sheets()

	for sheets.Next() {
		if err := fn(sheets.Sheet()); err != nil {
			return err
		}
	}

	return sheets.Err()
}

// ForSheet calls fn for the first sheet matching name.
func (r *Reader) ForSheet(name string, fn func(sheet *SheetCursor) error) error {
	sheets := r.Sheets()

	for sheets.Next() {
		sheet := sheets.Sheet()
		if sheet.Sheet.Name != name {
			continue
		}

		if err := fn(sheet); err != nil {
			return err
		}

		return nil
	}

	if err := sheets.Err(); err != nil {
		return err
	}

	return ErrSheetNotFound
}

// ForSheets calls fn for each sheet whose name is listed in names.
func (r *Reader) ForSheets(names []string, fn func(sheet *SheetCursor) error) error {
	targetSheets := make(map[string]bool, len(names))
	for _, name := range names {
		targetSheets[name] = true
	}

	if len(targetSheets) == 0 {
		return nil
	}

	sheets := r.Sheets()

	for sheets.Next() {
		sheet := sheets.Sheet()

		if !targetSheets[sheet.Sheet.Name] {
			continue
		}

		if err := fn(sheet); err != nil {
			return err
		}

		delete(targetSheets, sheet.Sheet.Name)
		if len(targetSheets) == 0 {
			return nil
		}
	}

	if err := sheets.Err(); err != nil {
		return err
	}

	return ErrSheetNotFound
}
