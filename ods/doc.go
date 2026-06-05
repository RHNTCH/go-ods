// Package ods provides streaming and in-memory APIs for reading ODS spreadsheets.
//
// Use Open to create a Reader for an .ods file. For low-memory, one-pass
// processing, use ForEachRow, ForSheet, ForSheets, or the sheet and row
// iterators. For convenient row and column access, use MakeTable.
//
// The parser reads content.xml sequentially from the ODS zip archive. Readers
// are single-pass: after rows or sheets have been consumed, open a new Reader
// to process the same file again.
package ods
