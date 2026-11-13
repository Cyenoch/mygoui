package mygoui_test

import (
	"bytes"
	"image"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// highlightFrame renders one highlighted paragraph hugging its text, and
// returns the frame with the paragraph's box.
func highlightFrame(t *testing.T, o mygoui.HighlightOptions) (*image.RGBA, ui.Rect) {
	t.Helper()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Row(c).AlignItems(ui.Start).Children(func() { mygoui.Highlight(c, o) })
	}, 320, 60)
	a, ok := tt.Access(o.Text)
	if !ok {
		t.Fatalf("%q does not read as one paragraph: %v", o.Text, tt.Texts())
	}
	return tt.Image(), a.Bounds
}

func sameHighlight(t *testing.T, a, b mygoui.HighlightOptions) bool {
	t.Helper()
	x, _ := highlightFrame(t, a)
	y, _ := highlightFrame(t, b)
	return bytes.Equal(x.Pix, y.Pix)
}

// TestHighlightMarksMatches pins which text is marked: blank queries mark
// nothing, case matters unless IgnoreCase folds it (also where a letter's
// cases differ in length), and overlapping queries share one mark.
func TestHighlightMarksMatches(t *testing.T) {
	plain := func(text string) mygoui.HighlightOptions { return mygoui.HighlightOptions{Text: text} }
	marked := func(text string, fold bool, queries ...string) mygoui.HighlightOptions {
		return mygoui.HighlightOptions{Text: text, Queries: queries, IgnoreCase: fold}
	}
	for _, tc := range []struct {
		name string
		a, b mygoui.HighlightOptions
		same bool
	}{
		{"a match is marked", marked("go build", false, "go"), plain("go build"), false},
		{"blank queries mark nothing", marked("go build", false, "", "  "), plain("go build"), true},
		{"queries are trimmed", marked("go build", false, " go "), marked("go build", false, "go"), true},
		{"case matters by default", marked("Go build", false, "go"), plain("Go build"), true},
		{"IgnoreCase folds case", marked("Go build", true, "gO"), marked("Go build", false, "Go"), true},
		{"folding across byte lengths", marked("ſtraße", true, "STR"), marked("ſtraße", false, "ſtr"), true},
		{"overlapping queries join", marked("abcde", false, "abc", "bcd"), marked("abcde", false, "abcd"), true},
		{"every occurrence is marked", marked("go go", false, "go"), marked("go go", false, "go go"), false},
	} {
		if got := sameHighlight(t, tc.a, tc.b); got != tc.same {
			t.Errorf("%s: frames equal = %v, want %v", tc.name, got, tc.same)
		}
	}
}

// TestHighlightMarksWholeMultiByteCharacters pins the placement in text of
// several bytes a character: the mark covers the matched character alone.
func TestHighlightMarksWholeMultiByteCharacters(t *testing.T) {
	const text = "日本語"
	base, box := highlightFrame(t, mygoui.HighlightOptions{Text: text})
	hit, _ := highlightFrame(t, mygoui.HighlightOptions{Text: text, Queries: []string{"本"}})
	lo, hi := float32(-1), float32(-1)
	for x := 0; x < base.Rect.Dx(); x++ {
		for y := 0; y < base.Rect.Dy(); y++ {
			if base.RGBAAt(x, y) != hit.RGBAAt(x, y) {
				if lo < 0 {
					lo = float32(x)
				}
				hi = float32(x)
				break
			}
		}
	}
	// Three characters of one width: the mark covers the middle one and
	// stops short of the centers of its neighbors.
	third := box.W / 3
	if lo < box.X+third/2 || lo > box.X+third+1 || hi < box.X+2*third-1 || hi > box.X+2.5*third {
		t.Fatalf("mark spans x=%v..%v, want the middle third of %+v", lo, hi, box)
	}
}

// TestHighlightFollowsItsOptions pins the refresh: the marks are those of
// the current options, including a query slice the caller reuses.
func TestHighlightFollowsItsOptions(t *testing.T) {
	queries := []string{"go"}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Highlight(c, mygoui.HighlightOptions{Text: "go build", Queries: queries})
	}, 200, 40)
	before := append([]byte(nil), tt.Image().Pix...)
	queries[0] = "build"
	tt.Frame()
	if bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("a query changed in place kept the old marks")
	}
}
