package gallery

import (
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// exampleFramesFixture is a self-contained page for the docs chrome tests,
// so catalog pages can move to ComponentExample frames without breaking them.
func exampleFramesFixture() ComponentPage {
	return ComponentPage{
		ID: "examples-fixture", Title: "Examples fixture", Group: "Foundations",
		Examples: []ComponentExample{
			{
				Title: "Button variants", Description: "Two named frames share one snippet.",
				Demo: func(c *ui.Context, s *State) {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Primary"})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Secondary"})
				},
			},
			{
				Title: "Second variant", Description: "Its own title and panel.",
				Demo: func(c *ui.Context, s *State) {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Ghost", Variant: mygoui.ButtonGhost})
				},
			},
		},
		Code: "package main\n\nimport \"github.com/Cyenoch/mygoui\"\n\n" +
			"func main() {\n\tprintln(\"hi\")\n}\n",
		Tokens: []string{"Accent"},
	}
}

// singleExampleFixture is a one-example page: its code area belongs to the
// example's own panel rather than a separate code panel.
func singleExampleFixture() ComponentPage {
	return ComponentPage{
		ID: "single-fixture", Title: "Single example", Group: "Foundations",
		Examples: []ComponentExample{
			{
				Title: "Solo demo", Description: "One example frame.",
				Demo: func(c *ui.Context, s *State) {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Solo action"})
				},
			},
		},
		Code: "package main\n\nfunc main() {}\n",
	}
}

// TestDocsHeaderActionsCopyRealContent pins the two functional header pills:
// Copy Markdown writes a generated document and Copy import writes the
// module import, each reporting Copied for the feedback window.
func TestDocsHeaderActionsCopyRealContent(t *testing.T) {
	state := New()
	state.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	for _, label := range []string{"Copy Markdown", "Copy import"} {
		if _, ok := tt.Find(label); !ok {
			t.Fatalf("header has no %q action", label)
		}
	}

	if err := tt.Click("Copy Markdown"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	markdown := tt.Clipboard()
	for _, want := range []string{
		"# Button", "## Usage", "```go", "## API Reference",
		"| Name | Type | Default | Description |", "## Theme tokens", "`Accent`",
	} {
		if !strings.Contains(markdown, want) {
			t.Fatalf("copied Markdown is missing %q:\n%s", want, markdown)
		}
	}
	if _, ok := tt.Find("Copied"); !ok {
		t.Fatal("Copy Markdown does not report Copied")
	}
	tt.Advance(codeCopyFeedback)
	if _, ok := tt.Find("Copy Markdown"); !ok {
		t.Fatal("Copy Markdown did not revert after the feedback window")
	}

	if err := tt.Click("Copy import"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != mygouiImport {
		t.Fatalf("Copy import copied %q, want %q", got, mygouiImport)
	}
}

// TestDocsSectionOrder pins the docs order: Usage, Examples, API Reference
// and Theme tokens.
func TestDocsSectionOrder(t *testing.T) {
	state := New()
	state.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	last := float32(-1)
	for _, heading := range []string{"Usage", "Examples", "API Reference", "Theme tokens"} {
		r, ok := tt.Find(heading)
		if !ok {
			t.Fatalf("section %q is missing", heading)
		}
		if r.Y <= last {
			t.Fatalf("section %q is out of order (y %v after %v)", heading, r.Y, last)
		}
		last = r.Y
	}
}

// TestDocsTOCScrollsTheContent clicks every rail link and proves the content
// scroll moves the target into view.
func TestDocsTOCScrollsTheContent(t *testing.T) {
	state := New()
	state.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	if _, ok := tt.Find("On this page"); !ok {
		t.Fatal("wide docs page has no On this page rail")
	}
	for _, item := range []struct {
		label, heading string
		nearTop        bool
	}{
		{"Go to Usage", "Usage", true},
		{"Go to Examples", "Examples", true},
		{"Go to API Reference", "API Reference", true},
		{"Go to Theme tokens", "Theme tokens", false},
	} {
		state.Content.Y = 0
		tt.Frame()
		if err := tt.Click(item.label); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if state.Content.Y <= 0 {
			t.Fatalf("clicking %q did not scroll the content", item.label)
		}
		r, ok := tt.Find(item.heading)
		if !ok {
			t.Fatalf("%q is missing after the rail click", item.heading)
		}
		if r.Y < 0 || r.Y > 900 {
			t.Fatalf("clicking %q left %q outside the viewport (y %v)", item.label, item.heading, r.Y)
		}
		if item.nearTop && r.Y > 200 {
			t.Fatalf("clicking %q put %q at y %v, want it near the top", item.label, item.heading, r.Y)
		}
	}
}

// TestDemoTitlesSitOutsidePreviewPanels pins the HeroUI example shape: the
// title and description are outside the bordered panel, never repeated
// inside it, and the panel holds only the live body.
func TestDemoTitlesSitOutsidePreviewPanels(t *testing.T) {
	p := exampleFramesFixture()
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()
	title, ok := tt.Find("Button variants")
	if !ok {
		t.Fatal("demo title is missing")
	}
	panel, ok := tt.Find("Example panel: Button variants")
	if !ok {
		t.Fatal("demo panel is missing")
	}
	if title.Y+title.H > panel.Y || title.X < panel.X || title.X >= panel.X+panel.W {
		t.Fatalf("demo title is not outside the panel (title %+v, panel %+v)", title, panel)
	}
	if title.Y >= panel.Y && title.Y+title.H <= panel.Y+panel.H &&
		title.X >= panel.X && title.X+title.W <= panel.X+panel.W {
		t.Fatalf("demo title is inside the panel (title %+v, panel %+v)", title, panel)
	}
	count := 0
	for _, text := range tt.Texts() {
		if text == "Button variants" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("demo title appears %d times, want exactly once outside the panel", count)
	}
	body, ok := tt.Find("Primary")
	if !ok || body.X < panel.X || body.X+body.W > panel.X+panel.W ||
		body.Y < panel.Y || body.Y+body.H > panel.Y+panel.H {
		t.Fatalf("demo body is outside its panel (body %+v, panel %+v)", body, panel)
	}
}

// TestSingleExampleOwnsItsCode checks a single-frame page puts the code area
// in the same panel under a divider, while a multi-frame page shares one code
// panel below a "Code" label.
func TestSingleExampleOwnsItsCode(t *testing.T) {
	single := New()
	tt := ui.NewTester(func(c *ui.Context) {
		buildComponentPage(c, single, singleExampleFixture())
	}, 1180, 1200)
	tt.Frame()
	code, ok := tt.Find("Go code")
	if !ok {
		t.Fatal("single example code is missing")
	}
	panel, ok := tt.Find("Example panel: Solo demo")
	if !ok || code.Y <= panel.Y || code.Y+code.H > panel.Y+panel.H {
		t.Fatalf("single example code is outside its example panel (code %+v, panel %+v)", code, panel)
	}
	if _, ok := tt.Find("Example code panel"); ok {
		t.Fatal("single example page renders a separate code panel")
	}

	multi := New()
	fixture := exampleFramesFixture()
	tt = ui.NewTester(func(c *ui.Context) { buildComponentPage(c, multi, fixture) }, 1180, 1400)
	tt.Frame()
	label, ok := tt.Find("Code")
	if !ok {
		t.Fatal("multi-frame page has no Code label")
	}
	codePanel, ok := tt.Find("Example code panel")
	if !ok {
		t.Fatal("multi-frame page has no code panel")
	}
	if label.Y >= codePanel.Y {
		t.Fatalf("Code label is not above the code panel (label %+v, panel %+v)", label, codePanel)
	}
	code, ok = tt.Find("Go code")
	if !ok || code.Y < codePanel.Y || code.Y+code.H > codePanel.Y+codePanel.H {
		t.Fatalf("shared code is outside the code panel (code %+v, panel %+v)", code, codePanel)
	}
}

// TestDocsRailHidesOnNarrowWindows checks the rail breakpoint and that the
// main column keeps the page title at every width.
func TestDocsRailHidesOnNarrowWindows(t *testing.T) {
	for _, tc := range []struct {
		width int
		rail  bool
	}{
		{1100, true},
		{999, false},
		{760, false},
		{640, false},
	} {
		state := New()
		state.Selected = "button"
		tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, tc.width, 900)
		tt.Frame()
		if _, ok := tt.Find("On this page"); ok != tc.rail {
			t.Fatalf("width %d: rail present = %v, want %v", tc.width, ok, tc.rail)
		}
		if !tt.HasText("Button") {
			t.Fatalf("width %d: page title is missing", tc.width)
		}
	}
}

// TestUsageBarIsACompactImportLine pins the Usage section to one short bar:
// the import line, a single copy icon and no line numbers.
func TestUsageBarIsACompactImportLine(t *testing.T) {
	state := New()
	p := exampleFramesFixture()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()
	bar, ok := tt.Find("Usage import")
	if !ok || bar.H < 40 || bar.H > 48 {
		t.Fatalf("usage bar rect = %+v, want a 40-44 DIP row", bar)
	}
	if !tt.HasText(mygouiImport) {
		t.Fatal("usage bar lost the import line")
	}
	if line, ok := tt.Find("line 1"); ok && line.Y >= bar.Y && line.Y < bar.Y+bar.H {
		t.Fatal("usage bar renders code line numbers")
	}
	copy, ok := tt.Find("Copy import line")
	if !ok || copy.W > 40 {
		t.Fatalf("usage copy control = %+v, want an icon-only button", copy)
	}
	if err := tt.Click("Copy import line"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != mygouiImport {
		t.Fatalf("usage copy wrote %q, want %q", got, mygouiImport)
	}
	if _, ok := tt.Find("Copied"); !ok {
		t.Fatal("usage copy control does not report Copied")
	}
}

// TestComponentPageImportStatement covers the optional Import field: empty
// derives the block from the usage snippet and a full statement passes
// through.
func TestComponentPageImportStatement(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want string
	}{
		{"default", "", mygouiImport},
		{"path", "github.com/Cyenoch/mygoui/form", `import "github.com/Cyenoch/mygoui/form"`},
		{"statement", `import "github.com/Cyenoch/mygoui/list"`, `import "github.com/Cyenoch/mygoui/list"`},
	} {
		if got := (ComponentPage{Import: tc.in}).importStatement(); got != tc.want {
			t.Fatalf("%s: importStatement() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// TestComponentPageImportStatementDerivesFromCode checks that the usage bar
// names the packages a page's snippet actually uses, while always keeping the
// module root, and that repeated qualifiers do not duplicate a path.
func TestComponentPageImportStatementDerivesFromCode(t *testing.T) {
	for _, tc := range []struct {
		id   string
		want []string
	}{
		{"input", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoui/form"`}},
		{"modal", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoui/overlay"`}},
		{"typography", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoui/theme"`}},
		{"link", []string{`"github.com/Cyenoch/mygoui"`, `"github.com/Cyenoch/mygoicons/lucide"`}},
	} {
		p, ok := componentPage(tc.id)
		if !ok {
			t.Fatalf("%s: page is not registered", tc.id)
		}
		got := p.importStatement()
		if !strings.HasPrefix(got, "import (\n") || !strings.HasSuffix(got, "\n)") {
			t.Fatalf("%s: several imports are not a parenthesized block:\n%s", tc.id, got)
		}
		for _, want := range tc.want {
			if count := strings.Count(got, want); count != 1 {
				t.Fatalf("%s: import block names %s %d times, want once:\n%s", tc.id, want, count, got)
			}
		}
	}
	// A snippet that names no subpackage keeps one line.
	if got := (ComponentPage{Code: "mygoui.Button(c, mygoui.ButtonOptions{})"}).importStatement(); got != mygouiImport {
		t.Fatalf("a single package import = %q, want %q", got, mygouiImport)
	}
}

// TestUsageAndHeaderCopyDerivedImport proves both copy controls write the
// derived block, not a stale module import.
func TestUsageAndHeaderCopyDerivedImport(t *testing.T) {
	p, ok := componentPage("input")
	if !ok {
		t.Fatal("input page is not registered")
	}
	want := p.importStatement()

	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()
	if err := tt.Click("Copy import line"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != want {
		t.Fatalf("usage copy wrote %q, want %q", got, want)
	}

	state = New()
	state.Selected = "input"
	tt = ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	if err := tt.Click("Copy import"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if got := tt.Clipboard(); got != want {
		t.Fatalf("header copy wrote %q, want %q", got, want)
	}
}

// TestUsageBarGrowsForMultipleImports checks that a page using several
// packages grows the bar to the import block's lines while keeping the copy
// control centered and the text intact.
func TestUsageBarGrowsForMultipleImports(t *testing.T) {
	p, ok := componentPage("input")
	if !ok {
		t.Fatal("input page is not registered")
	}
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()
	bar, ok := tt.Find("Usage import")
	if !ok {
		t.Fatal("usage bar is missing")
	}
	if bar.H <= 48 {
		t.Fatalf("usage bar height = %v, want the multi-line import block", bar.H)
	}
	if !tt.HasText(p.importStatement()) {
		t.Fatalf("usage bar does not show the import block:\n%s", p.importStatement())
	}
	copy, ok := tt.Find("Copy import line")
	if !ok {
		t.Fatal("usage copy is missing")
	}
	if d := (copy.Y + copy.H/2) - (bar.Y + bar.H/2); d < -1 || d > 1 {
		t.Fatalf("copy center is %.1f DIP off the grown bar's center", d)
	}
}

// TestExampleCodeControls pins the icon-only copy control in the code area
// and that the area belongs to a bordered panel.
func TestExampleCodeControls(t *testing.T) {
	state := New()
	p := exampleFramesFixture()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()
	code, ok := tt.Find("Go code")
	if !ok {
		t.Fatal("code area is missing")
	}
	panel, ok := tt.Find("Example code panel")
	if !ok || code.Y <= panel.Y || code.Y+code.H > panel.Y+panel.H {
		t.Fatalf("code area is outside its panel (code %+v, panel %+v)", code, panel)
	}
	copy, ok := tt.Find("Copy code")
	if !ok || copy.W > 40 {
		t.Fatalf("code copy control = %+v, want an icon-only button", copy)
	}
}

// TestComponentExampleFrames proves the ComponentExample contract: two named
// examples render two panels with their titles outside, and the one page
// snippet appears once in the shared code panel.
func TestComponentExampleFrames(t *testing.T) {
	p := ComponentPage{
		ID: "exampleframes", Title: "Example frames", Group: "Foundations",
		Examples: []ComponentExample{
			{
				Title: "First variant", Description: "The first named demo.",
				Demo: func(c *ui.Context, s *State) {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "First action"})
				},
			},
			{
				Title: "Second variant", Description: "The second named demo.",
				Demo: func(c *ui.Context, s *State) {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Second action"})
				},
			},
		},
		Code: "package main\n\nfunc main() {}\n",
	}
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 1180, 1400)
	tt.Frame()

	panels := 0
	codeHeaders := 0
	for _, text := range tt.Texts() {
		if strings.HasPrefix(text, "Example panel: ") {
			panels++
		}
		if text == "Go code header" {
			codeHeaders++
		}
	}
	if panels != 2 {
		t.Fatalf("rendered %d example panels, want 2", panels)
	}
	for _, example := range p.Examples {
		title, ok := tt.Find(example.Title)
		if !ok {
			t.Fatalf("example title %q is missing", example.Title)
		}
		panel, ok := tt.Find("Example panel: " + example.Title)
		if !ok {
			t.Fatalf("example panel %q is missing", example.Title)
		}
		if title.Y+title.H > panel.Y {
			t.Fatalf("example title %q is not above its panel (title %+v, panel %+v)", example.Title, title, panel)
		}
	}
	if _, ok := tt.Find("Example code panel"); !ok {
		t.Fatal("shared code panel is missing")
	}
	if codeHeaders != 1 {
		t.Fatalf("rendered %d code headers, want the snippet once", codeHeaders)
	}
}

// TestCopyControlsAlignAndDoNotStretchRows pins the copy icon geometry: a
// square 28-32 DIP button, centered on its header/bar within 1 DIP and inset
// 8-12 DIP from the rounded right edge. Collapsed/expanded code keeps it.
func TestCopyControlsAlignAndDoNotStretchRows(t *testing.T) {
	check := func(t *testing.T, row, copy ui.Rect, label string) {
		t.Helper()
		if copy.W < 28 || copy.W > 32 || abs32(copy.W-copy.H) > 0.5 {
			t.Fatalf("%s: copy control is %vx%v, want a 28-32 square", label, copy.W, copy.H)
		}
		if d := (copy.Y + copy.H/2) - (row.Y + row.H/2); d < -1 || d > 1 {
			t.Fatalf("%s: copy center is %.1f DIP off the row center", label, d)
		}
		if inset := row.X + row.W - (copy.X + copy.W); inset < 8 || inset > 12 {
			t.Fatalf("%s: copy right inset = %.1f, want 8-12", label, inset)
		}
	}
	for _, width := range []int{1180, 640} {
		state := New()
		fixture := exampleFramesFixture()
		tt := ui.NewTester(func(c *ui.Context) {
			buildComponentPage(c, state, fixture)
		}, width, 1400)
		tt.Frame()
		usage, ok := tt.Find("Usage import")
		if !ok {
			t.Fatalf("%d: usage bar is missing", width)
		}
		usageCopy, ok := tt.Find("Copy import line")
		if !ok {
			t.Fatalf("%d: usage copy is missing", width)
		}
		check(t, usage, usageCopy, "usage bar")

		header, ok := tt.Find("Go code header")
		if !ok {
			t.Fatalf("%d: code header is missing", width)
		}
		copy, ok := tt.Find("Copy code")
		if !ok {
			t.Fatalf("%d: code copy is missing", width)
		}
		check(t, header, copy, "code header")
		if _, ok := tt.Find("Expand code"); ok {
			if err := tt.Click("Expand code"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			header, _ = tt.Find("Go code header")
			copy, _ = tt.Find("Copy code")
			check(t, header, copy, "expanded code header")
		}
	}
}

// TestRoundedSurfacesClipCorners samples the four outside-corner pixels of
// the rounded gallery surfaces: they take the page background, while the
// pixels inside the corner take the surface color, in both modes.
func TestRoundedSurfacesClipCorners(t *testing.T) {
	for _, dark := range []bool{false, true} {
		state := New()
		pal := theme.ColorsForMode(dark)
		fixture := exampleFramesFixture()
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, dark)
			ui.Column(c).Fill().Padding(theme.Scale(c, theme.Space40)).
				Background(theme.Colors(c).Background).Children(func() {
				buildComponentPage(c, state, fixture)
			})
		}, 1180, 1600)
		tt.Frame()

		check := func(name string, r ui.Rect, radius float32, inside ui.Color) {
			t.Helper()
			img := tt.Image()
			bg := opaque(pal.Background)
			corners := [][2]int{
				{int(r.X) + 2, int(r.Y) + 2},
				{int(r.X+r.W) - 3, int(r.Y) + 2},
				{int(r.X) + 2, int(r.Y+r.H) - 3},
				{int(r.X+r.W) - 3, int(r.Y+r.H) - 3},
			}
			for _, p := range corners {
				if got := img.RGBAAt(p[0], p[1]); dist8(got, bg) > 3 {
					t.Fatalf("dark=%v %s: outside corner (%d,%d) = %v, want page background %v",
						dark, name, p[0], p[1], got, bg)
				}
			}
			x, y := int(r.X+radius/2), int(r.Y+radius/2)
			if got := img.RGBAAt(x, y); dist8(got, opaque(inside)) > 3 {
				t.Fatalf("dark=%v %s: inside corner = %v, want surface %v", dark, name, got, inside)
			}
		}

		usage, ok := tt.Find("Usage import")
		if !ok {
			t.Fatal("usage bar is missing")
		}
		check("usage bar", usage, theme.Space12, pal.SurfaceSecondary)
		previewPanel, ok := tt.Find("Example panel: Button variants")
		if !ok {
			t.Fatal("example preview panel is missing")
		}
		check("example preview panel", previewPanel, 20, pal.Surface)
		codePanel, ok := tt.Find("Example code panel")
		if !ok {
			t.Fatal("example code panel is missing")
		}
		check("example code panel", codePanel, 20, pal.Surface)
		card, ok := tt.Find("token card Accent")
		if !ok {
			t.Fatal("token card is missing")
		}
		check("token card", card, theme.Space12, pal.SurfaceSecondary)
	}
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
