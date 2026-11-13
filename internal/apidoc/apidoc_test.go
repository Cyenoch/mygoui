package apidoc

import (
	"bytes"
	"os"
	"testing"
)

// A stale example table would show code that no longer runs, or no Code tab
// at all, on every component page.
func TestGeneratedGalleryTablesAreCurrent(t *testing.T) {
	tables := []struct {
		file     string
		generate func(string) ([]byte, error)
	}{
		{"api_gen.go", Generate},
		{"example_source_gen.go", GenerateExamples},
	}
	for _, table := range tables {
		// A stale API table would show option fields that no longer exist,
		// or hide new ones, on every component page.
		want, err := table.generate("../..")
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile("../gallery/" + table.file)
		if err != nil {
			t.Fatal(err)
		}
		// Checkouts on Windows may convert line endings.
		got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
		if !bytes.Equal(got, want) {
			t.Fatalf("internal/gallery/%s is stale; run `make generate`", table.file)
		}
	}
}
