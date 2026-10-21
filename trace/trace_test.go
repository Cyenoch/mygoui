package trace_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/Cyenoch/mygoui/trace"
)

func TestStaticLinesSkipKeyboardFocusAndActionsActivate(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		trace.Trace(c).Children(func() {
			trace.Line(c, trace.Pending).Children(func() {
				trace.Title(c, "Pending: inspect source", trace.Pending)
			})
			trace.Action(c, trace.Options{Title: "Open result"}, func() { clicked++ })
			trace.Link(c, trace.Options{Title: "Read documentation"}, "https://ui.shadcn.com/docs/components/base/trace")
		})
	}, 420, 240)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Open result") || tt.Focused("Pending: inspect source") {
		t.Fatal("static trace text became a tab stop")
	}
	tt.Key(0, ui.KeyEnter)
	if clicked != 1 {
		t.Fatalf("Enter action calls=%d", clicked)
	}
	tt.Key(0, ui.KeySpace)
	if clicked != 2 {
		t.Fatalf("Space action calls=%d", clicked)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Read documentation") {
		t.Fatal("link missing keyboard focus")
	}
	tt.Key(0, ui.KeyEnter)
	if urls := tt.OpenedURLs(); len(urls) != 1 || urls[0] != "https://ui.shadcn.com/docs/components/base/trace" {
		t.Fatalf("link did not open actual URL: %v", urls)
	}
}

func TestCollapsibleKeyboardControlledFocusAndNestedState(t *testing.T) {
	var parent, nested trace.State
	open := true
	show := true
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			if show {
				trace.Collapsible(c.Key("source-parent"), &parent, trace.Options{
					Title: "Ran source inspection", Open: &open,
					OnOpenChange: func(bool) { changes++ },
				}, func() {
					trace.Collapsible(c.Key("source-nested"), &nested, trace.Options{
						Title: "Read source output", DefaultOpen: true,
					}, func() { ui.Button(c, "Copy source output") })
				})
			}
		})
	}, 420, 300)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Ran source inspection") {
		t.Fatal("disclosure missing keyboard focus")
	}
	tt.Key(0, ui.KeySpace)
	if open || changes != 1 || tt.HasText("Copy source output") || !nested.Open {
		t.Fatal("Space collapse did not preserve nested state")
	}
	if !tt.Focused("Ran source inspection") {
		t.Fatal("collapse lost trigger focus")
	}
	tt.Key(0, ui.KeyEnter)
	if !open || !tt.HasText("Copy source output") {
		t.Fatal("Enter reopen lost nested expanded content")
	}
	if err := tt.Click("Copy source output"); err != nil {
		t.Fatal(err)
	}
	open = false
	tt.Frame()
	if !tt.Focused("Ran source inspection") {
		t.Fatal("controlled close did not return focus from hidden body")
	}
	show = false
	tt.Frame()
	show, open = true, true
	tt.Frame()
	if !nested.Open || !tt.HasText("Copy source output") {
		t.Fatal("remount lost caller-owned nested expansion")
	}
}

func TestDefaultOpenInitializesOnceAndDisabledStaysClosed(t *testing.T) {
	var state, disabled trace.State
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		trace.Trace(c).Children(func() {
			trace.Collapsible(c.Key("default"), &state, trace.Options{
				Title: "Read source", DefaultOpen: true,
			}, func() { ui.Text(c, "Exact source detail") })
			trace.Collapsible(c.Key("disabled"), &disabled, trace.Options{
				Title: "Unavailable result", Disabled: true,
			}, func() { ui.Text(c, "Disabled detail") })
		})
	}, 420, 240)
	if !state.Open || !tt.HasText("Exact source detail") {
		t.Fatal("DefaultOpen did not initialize expansion")
	}
	if err := tt.Click("Read source"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state.Open || tt.HasText("Exact source detail") {
		t.Fatal("DefaultOpen overrode the user's collapse on rebuild")
	}
	if err := tt.Click("Unavailable result"); err != nil {
		t.Fatal(err)
	}
	if disabled.Open || tt.HasText("Disabled detail") {
		t.Fatal("disabled disclosure activated")
	}
}

func TestTitleTruncatesWhileDescriptionKeepsWidth(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		trace.Line(c, trace.Default).Children(func() {
			trace.Title(c, strings.Repeat("Read a long source name ", 15), trace.Default).Label("Long trace title")
			trace.Description(c, "128 results")
		})
	}, 240, 100)
	title, ok := tt.Find("Long trace title")
	if !ok {
		t.Fatal("full accessible title missing")
	}
	description, ok := tt.Find("128 results")
	if !ok || description.W <= 0 || description.X+description.W > 240 || title.X+title.W > description.X {
		t.Fatalf("title displaced untruncated description: title=%+v description=%+v", title, description)
	}
}

func TestRunningShimmerStopsForReducedMotion(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		trace.Line(c, trace.Running).Children(func() {
			trace.Title(c.Key("running-title"), "Running source inspection", trace.Running)
		})
	}, 420, 100)
	before := bytes.Clone(tt.Image().Pix)
	tt.Advance(600 * time.Millisecond)
	if bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("running title did not shimmer")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	before = bytes.Clone(tt.Image().Pix)
	tt.Advance(600 * time.Millisecond)
	if !bytes.Equal(before, tt.Image().Pix) || tt.FramePending() || !tt.HasText("Running source inspection") {
		t.Fatal("reduced motion kept shimmer or hid status words")
	}
}
