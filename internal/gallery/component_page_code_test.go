package gallery

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestCodeBlockCopyControl pins the compact header and the copy feedback: the
// header stays 40 DIP or shorter, Copy writes the untouched snippet, the
// button reports Copied, and it reverts after the feedback window.
func TestCodeBlockCopyControl(t *testing.T) {
	const code = "package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"
	p := ComponentPage{ID: "copyblock", Title: "Copy block", Group: "Foundations", Code: code}
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 700, 700)
	tt.Frame()

	header, ok := tt.Find("Go code header")
	if !ok || header.W <= 0 || header.H <= 0 || header.H > 40 {
		t.Fatalf("code header is not compact (rect %+v)", header)
	}
	if _, ok := tt.Find("Copy code"); !ok {
		t.Fatal("code header has no Copy code button")
	}
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != code {
		t.Fatalf("clipboard = %q, want the raw snippet %q", got, code)
	}
	if _, ok := tt.Find("Copied"); !ok {
		t.Fatal("Copy button does not report Copied")
	}

	tt.Advance(codeCopyFeedback)
	if _, ok := tt.Find("Copy code"); !ok {
		t.Fatal("Copy code button did not revert after the feedback window")
	}
	if _, ok := tt.Find("Copied"); ok {
		t.Fatal("Copy button still says Copied after the feedback window")
	}
}

// TestCodeBlockHugsSnippet proves the panel adds only its compact header and
// 12 DIP body padding, with no extra blank line under the last source line.
func TestCodeBlockHugsSnippet(t *testing.T) {
	code := strings.Join([]string{
		"mygoui.Button(c, mygoui.ButtonOptions{",
		"    Label: \"Save\",",
		"    Variant: mygoui.ButtonPrimary,",
		"})",
		"",
	}, "\n")
	p := ComponentPage{ID: "hugblock", Title: "Hug block", Group: "Foundations", Code: code}
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 700, 700)
	tt.Frame()
	block, ok := tt.Find("Go code")
	if !ok || block.H <= 0 {
		t.Fatalf("code block is missing (rect %+v)", block)
	}
	lines := float32(len(strings.Split(strings.TrimRight(code, "\n"), "\n")))
	if block.H >= lines*24+60 {
		t.Fatalf("code block wastes space: height %v for %v lines (rect %+v)", block.H, lines, block)
	}
}
