package apidoc

import (
	"bytes"
	"os"
	"testing"
)

// A stale API table would show option fields that no longer exist, or hide
// new ones, on every component page.
func TestGalleryAPITableIsCurrent(t *testing.T) {
	want, err := Generate("../..")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("../gallery/api_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	// Checkouts on Windows may convert line endings.
	got = bytes.ReplaceAll(got, []byte("\r\n"), []byte("\n"))
	if !bytes.Equal(got, want) {
		t.Fatal("internal/gallery/api_gen.go is stale; run `make generate`")
	}
}
