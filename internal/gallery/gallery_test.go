package gallery

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/fonts"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestMain(m *testing.M) {
	if err := fonts.Register(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestCatalogIsUniqueAndComplete(t *testing.T) {
	pages := Pages()

	// The catalog the docs promise, filed by the group the sidebar shows.
	// Overview is the only page registered without examples.
	want := map[string][]string{
		"Overview": {"overview"},
		"Foundations": {
			"theming", "typography", "surface", "separator", "kbd", "card", "frame", "link",
			"linkbutton", "listgroup", "emptystate", "scrollshadow", "scrollarea", "scrollfade", "header",
			"resizablepanel",
		},
		"Actions": {
			"button", "buttongroup", "group", "closebutton", "togglebutton",
			"togglebuttongroup", "pressablefeedback", "splitbutton",
		},
		"Display": {
			"chip", "badge", "avatar", "avatargroup", "tag", "taggroup", "alert",
			"spinner", "skeleton", "skeletongroup", "progressbar",
			"progresscircle", "meter", "image",
			"copybutton", "snippet", "qrcode",
			"descriptions", "statistic", "numberflip", "timeline", "statusdot", "highlight", "spoiler",
			"loadingoverlay", "banner",
		},
		"Forms": {
			"form", "field", "fieldset", "label", "description", "fielderror", "errormessage",
			"input", "textarea", "textfield", "inputgroup", "inputotp", "searchfield",
			"numberfield", "checkbox", "checkboxgroup", "radio", "radiogroup",
			"switch", "switchgroup", "controlfield", "slider", "select", "combobox",
			"autocomplete", "editabletext", "hotkeyinput", "rate", "imageselect",
			"passwordfield", "tagsinput", "rangeslider", "checkboxcard", "radiocard",
			"pathfield", "maskinput",
		},
		"Collections": {"listbox", "sortablelist", "tree"},
		"Navigation": {
			"accordion", "contextmenu", "segmentedcontrol", "disclosure", "disclosuregroup", "tabs", "breadcrumbs",
			"pagination", "toolbar", "menu", "dropdown", "submenu", "steps",
			"overflowlist", "carousel",
		},
		"Overlays": {
			"dialog", "alertdialog", "drawer", "popover",
			"tooltip", "toast", "previewcard", "sheet",
			"floatingpanel", "floatingsheet",
			"popconfirm", "teachingtip",
		},
		"Dates": {
			"calendar", "rangecalendar", "datefield", "timefield", "datepicker",
			"daterangepicker",
		},
		"Color": {
			"colorarea", "colorslider", "colorfield", "colorswatch",
			"colorswatchpicker", "colorpicker",
		},
		"Data":    {"table"},
		"Viewers": {"highlighter", "markdown", "imageviewer", "diffviewer", "jsonviewer", "logviewer"},
		"Patterns": {
			"command", "composer", "sidebar", "result", "guidecard", "upload", "keyvalueeditor",
			"menubar", "documenttabs", "navigationstack", "statusbar",
			"tour", "shortcutsdialog", "actionbar", "aboutdialog", "settingsgroup",
		},
		"Charts":          {"areachart", "barchart", "candlestickchart", "choroplethchart", "composedchart", "funnelchart", "gaugechart", "heatmapchart", "linechart", "profitlossline", "livelinechart", "piechart", "radarchart", "radialchart", "ringchart", "scatterchart", "sankeychart", "sparkline", "sunburstchart"},
		"Chart Utilities": {"chartlegend", "chartgrid", "chartbackground", "referencearea", "projectionline", "chartband", "chartweekends", "chartdeadline", "trendline", "chartmarkers", "charttooltip", "chartbrush", "chartxaxis", "chartyaxis", "chartcontext", "customindicator"},
	}
	expected := map[string]string{}
	for group, ids := range want {
		for _, id := range ids {
			if _, dup := expected[id]; dup {
				t.Fatalf("page %q is listed twice in the catalog table", id)
			}
			expected[id] = group
		}
	}
	if len(pages) != len(expected) {
		t.Fatalf("registered pages = %d, want %d", len(pages), len(expected))
	}

	seen := map[string]bool{}
	lastOrder := -1
	for _, page := range pages {
		if page.ID == "" || page.Title == "" || page.Group == "" || page.Build == nil {
			t.Fatalf("incomplete page registration: %+v", page)
		}
		if seen[page.ID] {
			t.Fatalf("duplicate page ID %q", page.ID)
		}
		seen[page.ID] = true
		if group, ok := expected[page.ID]; !ok || group != page.Group {
			t.Fatalf("page %q is registered under %q, want %q", page.ID, page.Group, expected[page.ID])
		}
		// Every component page carries the docs anatomy: a live example, a
		// usage snippet, API params and theme tokens.
		if p, ok := componentPage(page.ID); ok {
			if len(p.Examples) == 0 || p.Code == "" || len(p.Params) == 0 || len(p.Tokens) == 0 {
				t.Fatalf("page %q has an empty docs section: %d examples, %d params, %d tokens",
					page.ID, len(p.Examples), len(p.Params), len(p.Tokens))
			}
		}
		if page.Order < lastOrder {
			t.Fatalf("page %q out of order", page.ID)
		}
		lastOrder = page.Order
	}
	first := pages[0]
	if first.ID != "overview" || first.Title != "Overview" || first.Group != "Overview" {
		t.Fatalf("first page = %+v, want the Overview page", first)
	}
}

func TestEveryPageRendersAtDesktopAndCompactWidths(t *testing.T) {
	skipInShort(t)
	for _, page := range Pages() {
		for _, mode := range []struct {
			name   string
			width  int
			height int
			dark   bool
		}{
			{"compact-light", 640, 900, false},
			{"compact-dark", 640, 900, true},
			{"wide-light", 1100, 900, false},
			{"wide-dark", 1100, 900, true},
		} {
			state := New()
			state.Selected = page.ID
			state.Dark = mode.dark
			tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, mode.width, mode.height)
			if !tt.HasText(pageHeading(page)) {
				t.Fatalf("%s/%s: page title missing", mode.name, page.ID)
			}
			// Every named example keeps its own rendered frame, so a page
			// cannot silently lose a demo in any layout.
			if p, ok := componentPage(page.ID); ok {
				for _, example := range p.Examples {
					if _, ok := tt.Find("Example panel: " + example.Title); !ok {
						t.Fatalf("%s/%s: example %q did not render", mode.name, page.ID, example.Title)
					}
				}
			}
		}
	}
}

// pageHeading is the heading a page shows first: its title, or the home
// view's hero.
func pageHeading(page Page) string {
	if page.ID == "overview" {
		return "coss ui, native in Go"
	}
	return page.Title
}

// A component page is one document holding the component's own preview,
// source, reference and tokens.
func TestComponentPageAnatomy(t *testing.T) {
	state := New()
	state.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	for _, want := range []string{
		"Example panel: Variants", "Install command", "Go code", "mygoui.ButtonOptions",
		"Variant", "token card PrimaryForeground",
	} {
		if !tt.HasText(want) {
			t.Fatalf("the page does not show %q", want)
		}
	}
}

func TestPagesRenderWithEnlargedText(t *testing.T) {
	skipInShort(t)
	for _, page := range Pages() {
		state := New()
		state.Selected = page.ID
		tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1000, 900)
		tt.SetPreferences(ui.Preferences{TextScale: 1.75})
		if !tt.HasText(pageHeading(page)) {
			t.Fatalf("%s: title missing at enlarged text", page.ID)
		}
	}
}

func TestCompactNavigationOpensAndSelects(t *testing.T) {
	state := New()
	state.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 600, 900)
	tt.Frame()
	if err := tt.Click("Toggle navigation"); err != nil {
		t.Fatal(err)
	}
	if !state.NavOpen {
		t.Fatal("compact navigation did not open")
	}
	if err := tt.Click("Typography"); err != nil {
		t.Fatal(err)
	}
	if state.Selected != "typography" || state.NavOpen {
		t.Fatalf("compact selection = %q open=%v", state.Selected, state.NavOpen)
	}
}

func TestCompactNavigationScrollsToEveryPage(t *testing.T) {
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 600, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	if err := tt.Click("Toggle navigation"); err != nil {
		t.Fatal(err)
	}
	list, ok := tt.Find("Gallery navigation")
	if !ok || list.W < 250 || list.H < 600 {
		t.Fatalf("compact navigation list = %+v, want the sheet's scrolling body", list)
	}
	all := filterPagesInOrder(Pages())
	last := all[len(all)-1]
	for range 40 {
		if r, ok := tt.Find(last.Title); ok && r.Y+r.H <= list.Y+list.H {
			break
		}
		tt.Wheel(list.X+list.W/2, list.Y+list.H/2, 0, 400)
	}
	if err := tt.Click(last.Title); err != nil {
		t.Fatal(err)
	}
	if state.Selected != last.ID || state.NavOpen {
		t.Fatalf("selected %q open=%v, want %q", state.Selected, state.NavOpen, last.ID)
	}
}

func TestShellSearchFiltersAndSelects(t *testing.T) {
	state := New()
	// The sidebar belongs to component pages; the landing page has none.
	state.Selected = "theming"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 860)
	tt.Frame()
	if err := tt.Click("Button"); err != nil {
		t.Fatal(err)
	}
	if state.Selected != "button" {
		t.Fatalf("sidebar selection = %q", state.Selected)
	}
	state.Search = "card"
	tt.Frame()
	if !navHas(tt, "Card") {
		t.Fatal("searching card lost the Card page")
	}
	if navHas(tt, "Button") {
		t.Fatal("search kept an unrelated page")
	}
	state.Search = ""
	tt.Frame()
	if !navHas(tt, "Button") || !navHas(tt, "Typography") {
		t.Fatal("clearing the search did not restore the full navigation")
	}
}

// navHas reports whether the navigation shows name: its label sits in the
// sidebar column, left of the content pane.
func navHas(tt *ui.Tester, name string) bool {
	r, ok := tt.Find(name)
	return ok && r.X < 270
}

// revealSidebarEntry scrolls the sidebar until name is fully visible and
// returns the on-screen rect. The sidebar builds only the rows near its
// viewport, so a row is found by scrolling to it, half a viewport at a time.
func revealSidebarEntry(t *testing.T, tt *ui.Tester, state *State, name string, width, height float32) ui.Rect {
	t.Helper()
	for y := float32(0); ; y += height / 2 {
		state.Sidebar.Y = y
		tt.Frame()
		if r, ok := tt.Find(name); ok && r.X < 270 && contentVisible(r, width, height) {
			return r
		}
		if y >= state.Sidebar.MaxY {
			t.Fatalf("%q is unreachable in the sidebar (scroll %v/%v)", name, state.Sidebar.Y, state.Sidebar.MaxY)
		}
	}
}

func TestSidebarShowsEveryComponentPage(t *testing.T) {
	skipInShort(t)
	const width, height = 1180, 860
	state := New()
	// The sidebar belongs to component pages; the landing page has none.
	state.Selected = "theming"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, width, height)
	tt.Frame()
	for _, group := range pageGroups(filterPagesInOrder(Pages())) {
		if group.name == "Overview" {
			revealSidebarEntry(t, tt, state, group.pages[0].Title, width, height)
			continue
		}
		revealSidebarEntry(t, tt, state, group.name, width, height)
		for _, page := range group.pages {
			revealSidebarEntry(t, tt, state, page.Title, width, height)
		}
	}
}

func TestSidebarUsesScrollArea(t *testing.T) {
	state := New()
	// The sidebar belongs to component pages; the landing page has none.
	state.Selected = "theming"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 860)
	tt.Frame()
	r, ok := tt.Find("Gallery navigation")
	if !ok {
		t.Fatal("sidebar has no navigation scroll area")
	}
	if r.X >= 270 || r.W < 200 || r.H < 400 {
		t.Fatalf("sidebar scroll area = %+v, want a labeled ScrollArea viewport", r)
	}
}

// TestSidebarRowGeometryAndSelectedPaint pins the Coss nav look: compact rows with a 2 DIP gap, group headings, and the sidebar
// accent roles for the selected page.
func TestSidebarRowGeometryAndSelectedPaint(t *testing.T) {
	state := New()
	state.Selected = "card"
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) { pal = theme.Colors(c); View(c, state) }, 1180, 860)
	tt.Frame()

	first := sidebarRect(t, tt, "Typography")
	if first.H < 26 || first.H > 30 {
		t.Fatalf("first row height = %v, want ~28", first.H)
	}
	typography := sidebarRect(t, tt, "Typography")
	surface := sidebarRect(t, tt, "Surface")
	if gap := surface.Y - (typography.Y + typography.H); gap < 0 || gap > 4 {
		t.Fatalf("row gap = %v, want <= 4", gap)
	}
	if heading := sidebarRect(t, tt, "Actions"); heading.H < 20 || heading.H > 28 {
		t.Fatalf("group heading height = %v, want 20-28", heading.H)
	}

	card := sidebarRect(t, tt, "Card")
	img := tt.Image()
	selected := sampleRow(img, card)
	if want := over(pal.Sidebar, pal.SidebarAccent); dist8(selected, want) > 2 {
		t.Fatalf("selected row fill = %v, want SidebarAccent over Sidebar %v", selected, want)
	}
	if got := sampleRow(img, typography); dist8(got, opaque(pal.Sidebar)) > 2 {
		t.Fatalf("unselected row fill = %v, want Sidebar %v", got, pal.Sidebar)
	}
	if got, d := closestInk(img, card, pal.SidebarAccentForeground); d > 2 {
		t.Fatalf("selected row ink = %v, want SidebarAccentForeground %v (distance %d)", got, pal.SidebarAccentForeground, d)
	}
}

// sidebarRect returns the rect of a sidebar row, failing when the label is
// only in the content pane.
func sidebarRect(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok || r.X >= 270 {
		t.Fatalf("%q has no sidebar row (rect %+v)", name, r)
	}
	if r.W <= 0 || r.H <= 0 {
		t.Fatalf("%q is clipped (rect %+v)", name, r)
	}
	return r
}

// sampleRow reads the row fill four DIPs in, clear of the text.
func sampleRow(img *image.RGBA, r ui.Rect) color.RGBA {
	return img.RGBAAt(int(r.X)+4, int(r.Y+r.H/2))
}

// closestInk returns the pixel inside a row closest to the wanted text color.
func closestInk(img *image.RGBA, r ui.Rect, want ui.Color) (color.RGBA, int) {
	best, bestDist := color.RGBA{}, 1<<30
	for y := int(r.Y) + 6; y < int(r.Y+r.H)-6; y++ {
		for x := int(r.X) + 14; x < int(r.X+r.W)-8; x++ {
			if d := dist8(img.RGBAAt(x, y), opaque(want)); d < bestDist {
				best, bestDist = img.RGBAAt(x, y), d
			}
		}
	}
	return best, bestDist
}

// over composites src over an opaque dst in straight sRGB bytes.
func over(dst, src ui.Color) color.RGBA {
	a := float32(src.A) / 255
	blend := func(s, d uint8) uint8 { return uint8(float32(s)*a + float32(d)*(1-a) + 0.5) }
	return color.RGBA{blend(src.R, dst.R), blend(src.G, dst.G), blend(src.B, dst.B), 255}
}

func opaque(c ui.Color) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, 255}
}

func dist8(p, want color.RGBA) int {
	diff := func(a, b uint8) int {
		if a > b {
			return int(a - b)
		}
		return int(b - a)
	}
	d := diff(p.R, want.R)
	if v := diff(p.G, want.G); v > d {
		d = v
	}
	if v := diff(p.B, want.B); v > d {
		d = v
	}
	return d
}

// contentVisible reports whether r is fully inside the window and painted.
func contentVisible(r ui.Rect, width, height float32) bool {
	return r.W > 0 && r.H > 0 && r.X >= 0 && r.Y >= 0 && r.X+r.W <= width && r.Y+r.H <= height
}

// revealContent scrolls the content pane until name is fully visible.
func revealContent(t *testing.T, tt *ui.Tester, state *State, name string, width, height float32) ui.Rect {
	t.Helper()
	state.Content.Y = 0
	tt.Frame()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q is missing from the page", name)
	}
	if contentVisible(r, width, height) && r.Y > 2 && r.Y+r.H < height-2 {
		return r
	}
	target := r.Y - 200
	if target < 0 {
		target = 0
	}
	if target > state.Content.MaxY {
		target = state.Content.MaxY
	}
	state.Content.Y = target
	tt.Frame()
	r, ok = tt.Find(name)
	if !ok || !contentVisible(r, width, height) || r.Y+r.H >= height-2 {
		t.Fatalf("%q is unreachable on the page (rect %+v)", name, r)
	}
	return r
}

func TestCodeBlockLineNumbersAndHighlighting(t *testing.T) {
	for _, size := range []struct{ width, height int }{{1100, 900}, {640, 900}} {
		state := New()
		state.Selected = "button"
		width, height := float32(size.width), float32(size.height)
		tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, size.width, size.height)
		tt.Frame()
		if err := tt.Click("Code"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		// Bring the listing under the title bar so every line is on screen.
		state.Content.Y = 180
		tt.Frame()
		block, _ := tt.Find("Go code")
		if block.W <= 0 || block.X < 0 || block.X+block.W > width {
			t.Fatalf("%dx%d: code block lies outside the page (rect %+v)", size.width, size.height, block)
		}
		if !tt.HasText("mygoui") {
			t.Fatalf("%dx%d: code block lost its text", size.width, size.height)
		}
		// A ten-line snippet folds at eight lines.
		checkLines := func(from, to int) {
			t.Helper()
			var right float32
			previous := float32(-1)
			for n := from; n <= to; n++ {
				r, ok := tt.Find("line " + strconv.Itoa(n))
				if !ok {
					t.Fatalf("%dx%d: missing line number %d", size.width, size.height, n)
				}
				if r.W <= 0 || r.H <= 0 {
					t.Fatalf("%dx%d: line number %d is clipped (rect %+v)", size.width, size.height, n, r)
				}
				if n == from {
					right = r.X + r.W
				} else if d := r.X + r.W - right; d < -1 || d > 1 {
					t.Fatalf("%dx%d: line number %d is misaligned (right edge %v, want %v)", size.width, size.height, n, r.X+r.W, right)
				}
				if r.Y <= previous {
					t.Fatalf("%dx%d: line numbers are not stacked at %d (rect %+v)", size.width, size.height, n, r)
				}
				previous = r.Y
			}
		}
		cp, _ := componentPage("button")
		lines := len(strings.Split(exampleSource("button", cp.Examples[0].Title), "\n"))
		checkLines(1, lines)
		if _, ok := tt.Find("line " + strconv.Itoa(lines+1)); ok {
			t.Fatalf("%dx%d: more line numbers than code lines", size.width, size.height)
		}
		revealContent(t, tt, state, "Collapse code", width, height)
		if err := tt.Click("Collapse code"); err != nil {
			t.Fatal(err)
		}
		state.Content.Y = 180
		tt.Frame()
		checkLines(1, codeCollapseLines)
		if _, ok := tt.Find("line " + strconv.Itoa(codeCollapseLines+1)); ok {
			t.Fatalf("%dx%d: Collapse code did not fold the snippet", size.width, size.height)
		}
		if err := tt.Click("Expand code"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		checkLines(1, lines)
	}
}

func TestComponentPageHeaderAndAPIPills(t *testing.T) {
	state := New()
	state.Selected = "button"
	state.Search = "no such page"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	title, ok := tt.Find("Button")
	if !ok || title.X < 270 {
		t.Fatalf("page header is missing the title (rect %+v)", title)
	}
	for _, want := range []string{"Prop", "Type", "Default", "Description"} {
		if !tt.HasText(want) {
			t.Fatalf("API table is missing the %q header", want)
		}
	}
	for _, want := range []string{"mygoui.ButtonVariant", "ButtonPrimary", "theme.ButtonSize"} {
		if !tt.HasText(want) {
			t.Fatalf("API table is missing the %q pill", want)
		}
	}
}

// TestCodeBlockHighlightMatchesThemeTokens proves the syntax colors come from
// the palette: keyword ink is Primary and string ink is Success.
func TestCodeBlockHighlightMatchesThemeTokens(t *testing.T) {
	// The Usage bar repeats the module import, so give this page a distinct
	// import line and let the code area own the strings the test samples.
	p := ComponentPage{
		ID: "highlight", Title: "Highlight", Group: "Foundations",
		Import: "example.com/other",
		Code:   "package main\n\nimport \"github.com/Cyenoch/mygoui\"\n",
	}
	state := New()
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		pal = theme.Colors(c)
		buildComponentPage(c, state, p)
	}, 1180, 900)
	tt.Frame()
	img := tt.Image()
	for _, want := range []struct {
		text  string
		color ui.Color
	}{
		{"package main", pal.Destructive},
		{`import "github.com/Cyenoch/mygoui"`, pal.Success},
	} {
		r, ok := tt.Find(want.text)
		if !ok {
			t.Fatalf("code line %q is missing", want.text)
		}
		if got, d := closestInk(img, r, want.color); d > 2 {
			t.Fatalf("code line %q ink = %v, want %v (distance %d)", want.text, got, want.color, d)
		}
	}
}

func TestAPIShowsFullTypeValues(t *testing.T) {
	for _, size := range []struct{ width, height int }{{1100, 900}, {640, 900}} {
		state := New()
		state.Selected = "avatargroup"
		tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, size.width, size.height)
		tt.Frame()
		if !tt.HasText("[]mygoui.AvatarOptions") {
			t.Fatalf("%dx%d: API lost the full []mygoui.AvatarOptions type", size.width, size.height)
		}
		for _, text := range tt.Texts() {
			if strings.Contains(text, "…") {
				t.Fatalf("%dx%d: truncated API text %q", size.width, size.height, text)
			}
		}
		r := revealContent(t, tt, state, "[]mygoui.AvatarOptions", float32(size.width), float32(size.height))
		if r.W < 120 || r.X < 0 || r.X+r.W > float32(size.width) {
			t.Fatalf("%dx%d: type pill is not fully readable (rect %+v)", size.width, size.height, r)
		}
	}
}

func TestCodeBlockScrollsLongLines(t *testing.T) {
	long := `var Options = mygoui.ButtonOptions{Label: "` + strings.Repeat("x", 120) + `"}`
	p := ComponentPage{
		ID: "longcode", Title: "Long code", Group: "Foundations",
		Code: "package main\n\n" + long + "\n",
	}
	state := New()
	tt := ui.NewTester(func(c *ui.Context) { buildComponentPage(c, state, p) }, 640, 900)
	tt.Frame()
	scroll := StateOf[codeScrollState](state, "code:longcode")
	if scroll.Scroll.MaxX <= 0 {
		t.Fatalf("code block does not scroll a long line (MaxX = %v)", scroll.Scroll.MaxX)
	}
	if !tt.HasText(long) {
		t.Fatal("code block does not hold the whole long line")
	}
	for _, text := range tt.Texts() {
		if strings.Contains(text, "…") {
			t.Fatalf("code block truncated a line: %q", text)
		}
	}
}

// skipInShort leaves a rendering-heavy test out of the race run, which
// passes -short: it drives a single goroutine, and the race detector slows
// the software rasterizer about twentyfold.
func skipInShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("rendering test; skipped with -short")
	}
}

// The sidebar builds only the rows near its viewport, but a row holding the
// keyboard focus stays built when it scrolls away.
func TestSidebarKeepsFocusedRowScrolledOutOfView(t *testing.T) {
	const width, height = 1180, 860
	state := New()
	state.Selected = "theming"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, width, height)
	tt.Frame()
	title := pageGroups(filterPagesInOrder(Pages()))[1].pages[0].Title
	r := revealSidebarEntry(t, tt, state, title, width, height)
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !tt.Focused(title) {
		t.Fatalf("clicking %q did not focus it", title)
	}
	state.Sidebar.Y = state.Sidebar.MaxY
	tt.Frame()
	tt.Frame()
	if !tt.Focused(title) {
		t.Fatalf("%q lost the focus when it scrolled out of view", title)
	}
}
