package codeblock_test

import (
	"image/color"
	"strings"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/codeblock"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func tester(text *string, state *codeblock.State, options *codeblock.Options) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() { codeblock.View(c, *text, state, *options) })
	}, 340, 440)
}

func TestCopyPreservesSourceAndAcknowledgesForTwoSeconds(t *testing.T) {
	for _, tc := range []struct{ name, source string }{
		{name: "blank lines and final newline", source: "first\n\nthird\n"},
		{name: "no final newline", source: "first\n\nthird"},
		{name: "only blank lines", source: "\n\n"},
		{name: "Unicode and tabs", source: "\t世界🙂\r\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state := codeblock.State{}
			o := codeblock.Options{Title: "Output", Footer: "Complete"}
			source := tc.source
			tt := tester(&source, &state, &o)
			if err := tt.Click("Copy code"); err != nil {
				t.Fatal(err)
			}
			if got := tt.Clipboard(); got != source {
				t.Fatalf("clipboard = %q, want %q", got, source)
			}
			if got := tt.Announcements(); len(got) != 1 || got[0] != "Copied" {
				t.Fatalf("announcements = %q", got)
			}
			tt.Advance(1999 * time.Millisecond)
			if _, ok := tt.Find("Copied"); !ok {
				t.Fatal("copied acknowledgment ended before two seconds")
			}
			tt.Advance(time.Millisecond)
			if _, ok := tt.Find("Copy code"); !ok {
				t.Fatal("copy action did not reset after two seconds")
			}
		})
	}
}

func TestStreamingCopyWaitsForCompletionAndUsesLatestText(t *testing.T) {
	text := "received\n"
	state := codeblock.State{}
	o := codeblock.Options{Title: "Output", Open: true}
	tt := tester(&text, &state, &o)
	_ = tt.Click("Copy code")
	if tt.Clipboard() != "" {
		t.Fatal("incomplete stream copied")
	}
	text += "\ncomplete\n"
	o.Open = false
	tt.Frame()
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != text {
		t.Fatalf("stale copy: %q", tt.Clipboard())
	}
}

func TestWrapAndKeyboardScrollPreserveSource(t *testing.T) {
	text := strings.Repeat("a long source line ", 24) + "\n" + strings.Repeat("more\n", 40)
	state := codeblock.State{}
	o := codeblock.Options{Title: "Output", MaxHeight: 160}
	tt := tester(&text, &state, &o)
	if state.Scroll.MaxX <= 0 || state.Scroll.MaxY <= 0 {
		t.Fatalf("scroll extents = %+v", state.Scroll)
	}
	if err := tt.Click("Code content"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyPageDown)
	if state.Scroll.X <= 0 || state.Scroll.Y <= 0 {
		t.Fatalf("keyboard did not scroll: %+v", state.Scroll)
	}
	if err := tt.Click("Wrap lines"); err != nil {
		t.Fatal(err)
	}
	if !state.Wrap {
		t.Fatal("wrap did not toggle")
	}
	tt.Frame()
	if state.Scroll.MaxX != 0 {
		t.Fatalf("wrapped content scrolls horizontally: %g", state.Scroll.MaxX)
	}
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != text {
		t.Fatal("wrap changed copied text")
	}
}

func TestFloatingActionsRemainKeyboardReachable(t *testing.T) {
	text := "output"
	state := codeblock.State{}
	o := codeblock.Options{}
	tt := tester(&text, &state, &o)
	tt.Move(339, 439)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Code content") {
		t.Fatal("Tab did not focus content")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Wrap lines") {
		t.Fatal("Tab did not reach floating wrap")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Copy code") {
		t.Fatal("Tab did not reach floating copy")
	}
	tt.Key(0, ui.KeyEnter)
	if tt.Clipboard() != text {
		t.Fatal("keyboard copy did not copy source")
	}
}

func TestNumberedDiffSelectionExcludesGuttersAndPreservesBlanks(t *testing.T) {
	text := "old\n\n世界🙂\nlast\n"
	state := codeblock.State{}
	o := codeblock.Options{Title: "Change", LineNumbers: true, StartLine: 42, Lines: []codeblock.Line{
		{Text: "old", Variant: codeblock.Removed},
		{Text: "", Variant: codeblock.Highlighted},
		{Text: "世界🙂", Variant: codeblock.Added},
		{Text: "last"},
	}}
	tt := tester(&text, &state, &o)
	for _, label := range []string{"42", "45", "Added line", "Removed line"} {
		if _, ok := tt.Find(label); !ok {
			t.Fatalf("missing %q", label)
		}
	}
	first, ok := tt.Find("old")
	if !ok {
		t.Fatal("missing first line")
	}
	last, ok := tt.Find("last")
	if !ok {
		t.Fatal("missing last line")
	}
	tt.Press(first.X, first.Y+first.H/2)
	tt.Move(last.X+last.W, last.Y+last.H/2)
	tt.Release(last.X+last.W, last.Y+last.H/2)
	tt.Command("copy")
	if got := tt.Clipboard(); got != text {
		t.Fatalf("selection = %q, want %q", got, text)
	}
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != text {
		t.Fatal("copy included line numbers or markers")
	}
}

func TestSelectionSurvivesPlainOutputAppend(t *testing.T) {
	text := "alpha beta"
	state := codeblock.State{}
	o := codeblock.Options{Title: "Output"}
	tt := tester(&text, &state, &o)
	r, _ := tt.Find(text)
	tt.Press(r.X, r.Y+r.H/2)
	tt.Move(r.X+r.W, r.Y+r.H/2)
	tt.Release(r.X+r.W, r.Y+r.H/2)
	text += "\ngamma\n"
	tt.Frame()
	tt.Command("copy")
	if got := tt.Clipboard(); got != "alpha beta" {
		t.Fatalf("append changed selection: %q", got)
	}
}

func TestVariantsSizesAndDownloadAreReal(t *testing.T) {
	text := "failed\n"
	state := codeblock.State{}
	o := codeblock.Options{Title: "Output"}
	tt := tester(&text, &state, &o)
	if _, ok := tt.Find("Download"); ok {
		t.Fatal("download without handler")
	}
	before, _ := tt.Find("Code block")
	backgroundBefore := tt.Image().RGBAAt(int(before.X+before.W/2), int(before.Y+1))
	o.Size, o.Variant = codeblock.SizeSmall, codeblock.Destructive
	var downloaded string
	o.OnDownload = func(value string) { downloaded = value }
	tt.Frame()
	after, _ := tt.Find("Code block")
	backgroundAfter := tt.Image().RGBAAt(int(after.X+after.W/2), int(after.Y+1))
	if backgroundAfter == backgroundBefore {
		t.Fatal("destructive variant did not change visible surface")
	}
	if after.H >= before.H {
		t.Fatalf("small size did not reduce height: %g -> %g", before.H, after.H)
	}
	if err := tt.Click("Download"); err != nil {
		t.Fatal(err)
	}
	if downloaded != text {
		t.Fatalf("download = %q", downloaded)
	}
}

func requireTokenInk(t *testing.T, tt *ui.Tester, source string, want ui.Color) {
	t.Helper()
	r, ok := tt.Find(source)
	if !ok {
		t.Fatalf("missing source paragraph %q", source)
	}
	img := tt.Image()
	pixels := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			got := img.RGBAAt(x, y)
			if nearInk(got, want) {
				pixels++
			}
		}
	}
	if pixels < 3 {
		t.Fatalf("source %q has no painted token ink matching %v", source, want)
	}
}

func nearInk(got color.RGBA, want ui.Color) bool {
	return abs(int(got.R)-int(want.R)) < 12 &&
		abs(int(got.G)-int(want.G)) < 12 && abs(int(got.B)-int(want.B)) < 12
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func TestNativeSyntaxPaintSurvivesLineTreatmentsAndThemeChanges(t *testing.T) {
	text := "package main\nfunc greeting() string {\nreturn \"hello\"\n/* comment\ncontinued */\n}\n"
	state := codeblock.State{}
	o := codeblock.Options{Title: "main.go", Language: "go"}
	dark := true
	var syntax *editor.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		if dark {
			c.SetTheme(ui.DarkTheme())
		} else {
			c.SetTheme(ui.LightTheme())
		}
		theme.Apply(c)
		syntax = editor.DefaultTheme(c.Theme())
		ui.Column(c).Fill().Children(func() { codeblock.View(c, text, &state, o) })
	}, 440, 440)
	for _, tc := range []struct {
		name                    string
		numbered, wrapped, dark bool
	}{
		{name: "source", dark: true},
		{name: "numbered diff", numbered: true, dark: true},
		{name: "wrapped diff", numbered: true, wrapped: true, dark: true},
		{name: "light diff", numbered: true, wrapped: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dark, o.LineNumbers, state.Wrap = tc.dark, tc.numbered, tc.wrapped
			if tc.numbered {
				o.Lines = []codeblock.Line{
					{Text: "package main", Variant: codeblock.Highlighted},
					{Text: "func greeting() string {", Variant: codeblock.Added},
					{Text: "return \"hello\"", Variant: codeblock.Removed},
					{Text: "/* comment"},
					{Text: "continued */"},
					{Text: "}"},
				}
			}
			tt.Frame()
			requireTokenInk(t, tt, "package main", syntax.TokenColor(editor.TokenKeyword))
			requireTokenInk(t, tt, "func greeting() string {", syntax.TokenColor(editor.TokenFunction))
			requireTokenInk(t, tt, "return \"hello\"", syntax.TokenColor(editor.TokenString))
			requireTokenInk(t, tt, "continued */", syntax.TokenColor(editor.TokenComment))
			if err := tt.Click("Copy code"); err != nil {
				t.Fatal(err)
			}
			if tt.Clipboard() != text {
				t.Fatalf("copy = %q, want exact source %q", tt.Clipboard(), text)
			}
			tt.Advance(2 * time.Second)
		})
	}
	// Removing a multiline opener must invalidate tokens after the edited line.
	text = strings.Replace(text, "/* comment", "// comment", 1)
	tt.Frame()
	requireTokenInk(t, tt, "continued */", syntax.Text)
	text = "var message = `hello\nraw continuation`\n"
	tt.Frame()
	requireTokenInk(t, tt, "raw continuation`", syntax.TokenColor(editor.TokenString))
	// A changed, unknown language clears cached syntax without guessing from Title.
	o.Language = "unrecognized-language"
	tt.Frame()
	requireTokenInk(t, tt, "raw continuation`", syntax.Text)
}

func TestSyntaxSelectionSurvivesAppendWrapNumbersAndTheme(t *testing.T) {
	text := "package main\nfunc main() {}\n"
	state := codeblock.State{}
	o := codeblock.Options{Title: "main.go", Language: "go", Open: true}
	dark := true
	tt := ui.NewTester(func(c *ui.Context) {
		if dark {
			c.SetTheme(ui.DarkTheme())
		} else {
			c.SetTheme(ui.LightTheme())
		}
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() { codeblock.View(c, text, &state, o) })
	}, 340, 440)
	r, ok := tt.Find("package main")
	if !ok {
		t.Fatal("missing package source")
	}
	tt.Press(r.X, r.Y+r.H/2)
	tt.Move(r.X+r.W, r.Y+r.H/2)
	tt.Release(r.X+r.W, r.Y+r.H/2)
	text += "// streamed\n"
	tt.Frame()
	state.Wrap, o.LineNumbers, dark = true, true, false
	tt.Frame()
	tt.Command("copy")
	if got := tt.Clipboard(); got != "package main" {
		t.Fatalf("transition changed selection: %q", got)
	}
	o.Open = false
	tt.Frame()
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != text {
		t.Fatalf("completion copied stale source: %q", tt.Clipboard())
	}
}
