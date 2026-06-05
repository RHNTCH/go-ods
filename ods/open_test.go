package ods

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func writeTestZip(t *testing.T, files map[string]string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.ods")

	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	for name, content := range files {
		zipFile, err := zipWriter.Create(name)
		if err != nil {
			t.Fatal(err)
		}

		_, err = zipFile.Write([]byte(content))
		if err != nil {
			t.Fatal(err)
		}
	}

	return path
}

func TestOpenReturnsErrorForMissingFile(t *testing.T) {
	_, err := Open(filepath.Join(t.TempDir(), "missing.ods"))
	if err == nil {
		t.Fatal("Open() err = nil, want error")
	}
}

func TestOpenReturnsErrorWithoutContentXML(t *testing.T) {
	path := writeTestZip(t, map[string]string{
		"meta.xml": "<meta></meta>",
	})

	_, err := Open(path)
	if err == nil {
		t.Fatal("Open() err = nil, want error")
	}
}
