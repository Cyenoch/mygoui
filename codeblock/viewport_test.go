package codeblock_test

import (
	"fmt"
	"image/png"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/codeblock"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/egoist/mygo/ui/markdown"
)

// MYGOUI_SOURCE_VIEWPORT_CAPTURE=1 writes the original screenshot scenario's
// headless proof frames to /tmp/source-viewport-*.png when Main runs this test.
func TestSourceViewportHighlightedBothAxesBoundaries(t *testing.T) {
	var source strings.Builder
	source.WriteString("package main\n\n")
	for i := range 26 {
		fmt.Fprintf(&source, "var sourceLine%02d = []string{\"native highlighted Go source\", \"a long value preserving both axis overflow and scrollbar dragging\"}\n", i)
	}
	source.WriteString("var finalSourceLine = 42\n")
	text := source.String()
	state := codeblock.State{}
	o := codeblock.Options{Title: "Actual preview source", Footer: "Original source", Language: "go", LineNumbers: true, MaxHeight: 187, Embedded: true}
	for line := range strings.SplitSeq(strings.TrimSuffix(text, "\n"), "\n") {
		o.Lines = append(o.Lines, codeblock.Line{Text: line, Variant: codeblock.Highlighted})
	}
	var surface ui.Color
	var syntax *editor.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(ui.DarkTheme())
		theme.Apply(c)
		colors := theme.Colors(c)
		surface, syntax = colors.Muted.Over(colors.Background), editor.DefaultTheme(c.Theme())
		ui.Column(c).Fill().Padding(24).Background(theme.Colors(c).Background).Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Radius(theme.Radii().LG).Border(1, theme.Colors(c).Border).Clip().Children(func() {
				ui.Box(c).FillWidth().Height(54).Background(theme.Colors(c).Card)
				codeblock.View(c, text, &state, o)
			})
		})
	}, 620, 400)
	capture := func(name string) {
		t.Helper()
		if os.Getenv("MYGOUI_SOURCE_VIEWPORT_CAPTURE") != "1" {
			return
		}
		path := "/tmp/source-viewport-" + name + ".png"
		f, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		err = png.Encode(f, tt.Image())
		closeErr := f.Close()
		if err != nil {
			t.Fatal(err)
		}
		if closeErr != nil {
			t.Fatal(closeErr)
		}
		t.Log(path)
	}
	viewport := func() ui.Rect {
		t.Helper()
		r, ok := tt.Find("Code content")
		if !ok {
			t.Fatal("missing native source viewport")
		}
		return r
	}
	assertLanes := func() {
		t.Helper()
		r := viewport()
		img := tt.Image()
		lane := float32(12)
		if state.Wrap {
			lane = 0
		}
		// Sample the entire physical bottom lane away from corners.
		for y := int(math.Ceil(float64(r.Y + r.H - lane))); y < int(r.Y+r.H); y++ {
			for x := int(r.X) + 24; x < int(r.X+r.W)-24; x++ {
				if got := img.RGBAAt(x, y); !nearInk(got, surface) {
					t.Fatalf("code/fade reappeared beneath clip in horizontal bar lane at %d,%d: %v want %v", x, y, got, surface)
				}
			}
		}
		for y := int(r.Y) + 24; y < int(r.Y+r.H)-24; y++ {
			if got := img.RGBAAt(int(r.X+r.W)-2, y); !nearInk(got, surface) {
				t.Fatalf("code/fade entered native vertical bar lane at y=%d: %v", y, got)
			}
		}
	}
	if state.Scroll.MaxX <= 0 || state.Scroll.MaxY <= 0 {
		t.Fatalf("screenshot scenario lacks both-axis overflow: %+v", state.Scroll)
	}
	requireTokenInk(t, tt, "package main", syntax.TokenColor(editor.TokenKeyword))
	assertLanes()
	capture("start")
	initialViewport := viewport()
	tt.Move(initialViewport.X+initialViewport.W/2, initialViewport.Y+initialViewport.H/2)
	capture("start-scrollbars")
	tt.Move(619, 399)
	// Native source selection must still exclude the decorative gutter.
	first, _ := tt.Find("package main")
	tt.Press(first.X, first.Y+first.H/2)
	tt.Move(first.X+first.W, first.Y+first.H/2)
	tt.Release(first.X+first.W, first.Y+first.H/2)
	tt.Command("copy")
	if got := tt.Clipboard(); got != "package main" {
		t.Fatalf("source selection = %q", got)
	}
	tt.Move(619, 399)
	state.Scroll.X, state.Scroll.Y = state.Scroll.MaxX/2, state.Scroll.MaxY/2
	tt.Frame()
	assertLanes()
	capture("middle")
	state.Scroll.X, state.Scroll.Y = 0, math.MaxFloat32
	tt.Frame()
	assertLanes()
	last, ok := tt.Find("var finalSourceLine = 42")
	r := viewport()
	if !ok || last.Y < r.Y || last.Y+last.H > r.Y+r.H-12 {
		t.Fatalf("last source line is not fully revealable above horizontal scrollbar: last=%+v viewport=%+v", last, r)
	}
	requireTokenInk(t, tt, "var finalSourceLine = 42", syntax.TokenColor(editor.TokenKeyword))
	capture("end")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	capture("end-scrollbars")
	tt.Move(619, 399)
	// Native horizontal and vertical thumbs remain independently draggable.
	state.Scroll.X, state.Scroll.Y = 0, 0
	tt.Frame()
	r = viewport()
	tt.Press(r.X+20, r.Y+r.H-6)
	tt.Move(r.X+r.W/2, r.Y+r.H-6)
	tt.Release(r.X+r.W/2, r.Y+r.H-6)
	if state.Scroll.X <= 0 {
		t.Fatal("source horizontal scrollbar drag was intercepted")
	}
	tt.Press(r.X+r.W-6, r.Y+10)
	tt.Move(r.X+r.W-6, r.Y+r.H/2)
	tt.Release(r.X+r.W-6, r.Y+r.H/2)
	if state.Scroll.Y <= 0 {
		t.Fatal("source vertical scrollbar drag was intercepted")
	}
	tt.Move(619, 399)
	tt.SetSize(420, 400)
	tt.Frame()
	assertLanes()
	capture("resized")
	if err := tt.Click("Code content"); err != nil {
		t.Fatal(err)
	}
	state.Scroll.X, state.Scroll.Y = 0, 0
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyPageDown)
	if state.Scroll.X <= 0 || state.Scroll.Y <= 0 {
		t.Fatal("bounded source lost native keyboard scrolling")
	}
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	if got := tt.Clipboard(); got != text {
		t.Fatal("scrolling or selection changed original source copy")
	}
	if err := tt.Click("Wrap lines"); err != nil {
		t.Fatal(err)
	}
	if state.Scroll.MaxX != 0 {
		t.Fatalf("wrapping retained horizontal overflow: %+v", state.Scroll)
	}
	capture("wrapped")
	state.Wrap = false
	text, o.Lines = "package main\n", nil
	tt.Frame()
	if state.Scroll.MaxX != 0 || state.Scroll.MaxY != 0 || state.Scroll.X != 0 || state.Scroll.Y != 0 {
		t.Fatalf("shortened source retained stale scroll extents: %+v", state.Scroll)
	}
	assertLanes()
	capture("shrink")
	text = source.String()
	o.MaxHeight = 18
	tt.Frame()
	assertLanes()
	capture("tiny")
}

func TestMarkdownSourceUsesMaintainedHorizontalViewport(t *testing.T) {
	text := "```go\nvar nativeSource = \"" + strings.Repeat("wide highlighted source ", 20) + "\"\n```"
	var options markdown.Options
	options.CodeBlock = func(c *ui.Context, block *markdown.CodeBlockView) {
		codeblock.Markdown(c, block, &options)
		// Markdown's public adapter retains its State locally; the observable
		// paragraph position below proves the maintained scroller actually moves.
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() { markdown.View(c, text, &options) })
	}, 360, 240)
	paragraph := "var nativeSource = \"" + strings.Repeat("wide highlighted source ", 20) + "\""
	before, ok := tt.Find(paragraph)
	if !ok {
		t.Fatal("missing maintained Markdown source paragraph")
	}
	if err := tt.Click("Code content"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	after, ok := tt.Find(paragraph)
	if !ok || after.X >= before.X {
		t.Fatal("maintained Markdown source viewport did not scroll horizontally")
	}
	if err := tt.Click("Copy code"); err != nil {
		t.Fatal(err)
	}
	want := strings.TrimSuffix(strings.TrimPrefix(text, "```go\n"), "```")
	if got := tt.Clipboard(); got != want {
		t.Fatalf("maintained Markdown body copy = %q, want %q", got, want)
	}
	if err := tt.Click("Wrap lines"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.Find("Disable line wrap"); !ok {
		t.Fatal("Markdown source wrap action lost retained state")
	}
}
