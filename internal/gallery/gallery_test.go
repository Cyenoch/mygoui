package gallery

import (
	"image"
	"math"
	"strconv"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSearchNoMatchAndRecovery(t *testing.T) {
	s := New()
	s.Search = "no component has this name"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1000, 760)
	if !tt.HasText("No matching components") {
		t.Fatal("missing no-match state")
	}
	if err := tt.Click("Clear search"); err != nil {
		t.Fatal(err)
	}
	if s.Search != "" || !tt.HasText("Button") {
		t.Fatal("clearing search did not restore catalog navigation")
	}
	if s.Selected != "code-block" {
		t.Fatal("search changed the independent selected preview")
	}
}

func TestGalleryUsesActualButtonStateAndDisablesActivation(t *testing.T) {
	s := New()
	s.Selected, s.Search = "button", "Button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1100, 900)
	stageBefore, _ := tt.Find("Interactive preview")
	codeBefore, _ := tt.Find("View Code")
	if err := tt.Click("Default"); err != nil {
		t.Fatal(err)
	}
	if s.DemoBasic.Count != 1 || !tt.HasText("Sample action completed") {
		t.Fatal("preview action did not change caller fixture state")
	}
	stageAfter, _ := tt.Find("Interactive preview")
	codeAfter, _ := tt.Find("View Code")
	if stageBefore != stageAfter || codeBefore != codeAfter {
		t.Fatal("action feedback shifted the specimen or source disclosure")
	}
	if s.Basic.Count != 0 {
		t.Fatal("default preview mutated the independent advanced fixture")
	}
	advanced := ui.NewTester(func(c *ui.Context) { PreviewBasic(c, &s.Basic, "button") }, 800, 700)
	if err := advanced.Click("Default"); err != nil {
		t.Fatal(err)
	}
	if s.Basic.Count != 1 {
		t.Fatal("advanced fixture action did not activate")
	}
	if err := advanced.Click("Disabled"); err != nil {
		t.Fatal(err)
	}
	_ = advanced.Click("Default")
	if s.Basic.Count != 1 {
		t.Fatal("disabled library control activated")
	}
	if err := advanced.Click("Reset fixture"); err != nil {
		t.Fatal(err)
	}
	if s.Basic.Count != 0 || s.Basic.Disabled {
		t.Fatal("reset did not restore independent fixture state")
	}
}

func TestNarrowSearchAndLargeTextDoNotReplaceSource(t *testing.T) {
	defer theme.SetPreferences(theme.Preferences{})
	s := New()
	s.Search = "Trace"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 520, 1000)
	if err := tt.Click("Trace"); err != nil {
		t.Fatal(err)
	}
	if s.Selected != "trace" {
		t.Fatal("compact navigation did not select actual preview")
	}
	s.Trace.Source = "Retained local source"
	s.LargeText = true
	theme.SetPreferences(theme.Preferences{FontSize: 32})
	tt.Frame()
	if s.Trace.Source != "Retained local source" {
		t.Fatal("text fixture change replaced source")
	}
	if _, ok := tt.Find("Trace sample text"); !ok {
		t.Fatal("large-text narrow preview lost its editable source")
	}
}

// The toolbar must occupy its own layout band, including after either pane scrolls.
// Comparing pixels catches paint escaping a viewport; clicking the same coordinates
// catches an invisible scroll child covering the toolbar's hit targets.
func TestScrolledPanesPreserveToolbarPaintAndInput(t *testing.T) {
	s := New()
	s.Selected = "typography"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1440, 760)
	header, ok := tt.Find("Gallery toolbar")
	if !ok {
		t.Fatal("missing toolbar")
	}
	preview, ok := tt.Find("Component preview")
	if !ok {
		t.Fatal("missing preview viewport")
	}
	if preview.Y < header.Y+header.H || preview.Y+preview.H > 761 {
		t.Fatalf("preview escapes body below toolbar: header=%+v preview=%+v", header, preview)
	}
	before := image.NewRGBA(tt.Image().Bounds())
	copy(before.Pix, tt.Image().Pix)
	tt.Scroll(preview.X+preview.W/2, preview.Y+preview.H/2, 0, 5000)
	nav, ok := tt.Find("Component navigation")
	if !ok {
		t.Fatal("missing catalog viewport")
	}
	tt.Scroll(nav.X+nav.W/2, nav.Y+nav.H/2, 0, 5000)
	tt.Move(1439, 759)
	if s.previewScroll.Y <= 0 {
		t.Fatal("preview fixture did not scroll")
	}
	after := tt.Image()
	for y := int(header.Y); y < int(header.Y+header.H)-1; y++ {
		for x := int(header.X); x < int(header.X+header.W); x++ {
			if before.RGBAAt(x, y) != after.RGBAAt(x, y) {
				t.Fatalf("pane scroll repainted toolbar at %d,%d", x, y)
			}
		}
	}
	if err := tt.Click("Use light theme"); err != nil {
		t.Fatal(err)
	}
	if s.Dark {
		t.Fatal("scrolled content intercepted toolbar appearance action")
	}
}

func TestGalleryViewportGeometryAcrossWidthsAndTextSizes(t *testing.T) {
	defer theme.SetPreferences(theme.Preferences{})
	for _, font := range []float32{16, 32} {
		for _, width := range []int{390, 768, 1440} {
			t.Run(strconv.Itoa(width)+"/font-"+strconv.Itoa(int(font)), func(t *testing.T) {
				theme.SetPreferences(theme.Preferences{FontSize: font})
				s := New()
				s.Search, s.Selected = "Input", "input"
				tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, width, 1000)
				header, ok := tt.Find("Gallery toolbar")
				if !ok {
					t.Fatal("missing toolbar")
				}
				pane, ok := tt.Find("Component preview")
				if !ok || pane.Y < header.Y+header.H || pane.H <= 100 || pane.Y+pane.H > 1001 || pane.X < 0 || pane.X+pane.W > float32(width)+1 {
					t.Fatalf("unbounded or unusable preview at font %.0f: toolbar=%+v preview=%+v", font, header, pane)
				}
				stage, stageOK := tt.Find("Interactive preview")
				label, ok := tt.Find("Sample input")
				wantLeft := stage.X + max(float32(24), (stage.W-480)/2)
				if !ok || !stageOK || math.Abs(float64(label.X-wantLeft)) > 1 {
					t.Fatalf("form label lost leading edge inside intentional specimen measure: stage=%+v label=%+v", stage, label)
				}
				for _, name := range []string{"Search components", "Use light theme", "Display options"} {
					r, found := tt.Find(name)
					if !found || r.X < 0 || r.X+r.W > float32(width)+1 || r.Y < 0 || r.Y+r.H > header.Y+header.H+1 {
						t.Fatalf("toolbar control %q escapes at font %.0f: %+v", name, font, r)
					}
				}
				if err := tt.Click("Use light theme"); err != nil {
					t.Fatal(err)
				}
				if s.Dark {
					t.Fatal("appearance control is not usable")
				}
			})
		}
	}
}

func TestDocumentationCodeOutlineAndAdjacentNavigation(t *testing.T) {
	s := New()
	s.Selected = "dialog"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1440, 1000)
	if !tt.HasText("Go / MyGo") || !tt.HasText("Open Dialog") {
		t.Fatal("missing native source documentation")
	}
	if err := tt.Click("View Code"); err != nil {
		t.Fatal(err)
	}
	if !s.SourceOpen || !tt.HasText("Hide Code") || !tt.HasText("Actual preview source") {
		t.Fatal("source expansion did not change")
	}
	// The outline is outside the scrolling documentation pane and remains reachable.
	if err := tt.Click("Jump to Behavior and API"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if s.previewScroll.Y <= 0 {
		t.Fatal("outline jump did not move the documentation viewport")
	}
	s.previewScroll.Y = s.previewScroll.MaxY
	tt.Frame()
	if err := tt.Click("Direction →"); err != nil {
		t.Fatal(err)
	}
	if s.Selected != "direction" || !tt.HasText("Direction") {
		t.Fatal("next component did not open the selected documentation")
	}
}

func TestEveryDefaultKeepsDocumentAtTop(t *testing.T) {
	for _, entry := range Catalog {
		t.Run(string(entry.ID), func(t *testing.T) {
			for _, size := range []struct{ width, height int }{{390, 1000}, {768, 1000}, {1120, 820}} {
				s := New()
				s.Selected = entry.ID
				tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, size.width, size.height)
				tt.Frame()
				if s.previewScroll.Y != 0 {
					t.Fatalf("initial document scrolled away from heading at width %d: %g", size.width, s.previewScroll.Y)
				}
				if !tt.HasText("View Code") {
					t.Fatal("component lost actual preview/source disclosure")
				}
				stage, ok := tt.Find("Interactive preview")
				if !ok || stage.X < 0 || stage.X+stage.W > float32(size.width)+1 {
					t.Fatalf("component preview escapes document width %d: %+v", size.width, stage)
				}
			}
		})
	}
}

func TestRemainingFamiliesUseRealCallerOwnedActions(t *testing.T) {
	t.Run("progress", func(t *testing.T) {
		s := New()
		s.Selected = "progress"
		tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1120, 820)
		if err := tt.Click("Increase"); err != nil {
			t.Fatal(err)
		}
		if s.DemoDataLayout.Value != 50 || s.DataLayout.Value != 40 {
			t.Fatal("progress did not update only its default fixture")
		}
	})
	t.Run("tabs", func(t *testing.T) {
		s := New()
		s.Selected = "tabs"
		tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1120, 820)
		if err := tt.Click("Security"); err != nil {
			t.Fatal(err)
		}
		if s.DemoNavigation.Tab != "security" || !tt.HasText("Review the local sample security settings.") {
			t.Fatal("tab did not select its real panel")
		}
	})
	t.Run("accordion", func(t *testing.T) {
		s := New()
		s.Selected = "accordion"
		tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1120, 820)
		if err := tt.Click("Is it accessible?"); err != nil {
			t.Fatal(err)
		}
		if len(s.DemoNavigation.Accordion.Open) != 1 || !tt.HasText("Native disclosure triggers support keyboard focus and expansion.") {
			t.Fatal("accordion did not open actual content")
		}
	})
	t.Run("table data", func(t *testing.T) {
		s := NewDataLayoutState()
		tt := ui.NewTester(func(c *ui.Context) { dataLayoutSample(c, &s, "data-table", true) }, 800, 700)
		before := len(s.Rows)
		if err := tt.Click("Append local row"); err != nil {
			t.Fatal(err)
		}
		if len(s.Rows) != before+1 || s.Revision != 1 {
			t.Fatal("append action did not update actual data and query revision")
		}
	})
}

func TestMessageSpecimenUsesGroupedConversationAndReaction(t *testing.T) {
	s := New()
	s.Selected = "message"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1120, 1000)
	for _, text := range []string{"Deploying to prod real quick.", "It's 4:55 PM. On a Friday.", "It's a one-line change.", "Delivered", "Alright, let me take a look.", "Oliver is typing..."} {
		if !tt.HasText(text) {
			t.Fatalf("missing full conversation content %q", text)
		}
	}
	before, _ := tt.Find("Interactive preview")
	if err := tt.Click("Thumbs up, 1"); err != nil {
		t.Fatal(err)
	}
	if !s.DemoConversation.Selected || s.Conversation.Selected {
		t.Fatal("reaction did not update independent default conversation state")
	}
	after, _ := tt.Find("Interactive preview")
	if before != after {
		t.Fatal("reaction changed conversation specimen geometry")
	}
}

func TestSidebarSpecimenContainsCompleteShell(t *testing.T) {
	s := New()
	s.Selected = "sidebar"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1440, 1000)
	for _, text := range []string{"Acme Studio", "Alex Morgan", "alex@example.com", "Playground", "History", "Building your application", "Project activity"} {
		if !tt.HasText(text) {
			t.Fatalf("missing complete sidebar region %q", text)
		}
	}
	if err := tt.Click("Acme Studio"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Personal workspace"); err != nil {
		t.Fatal(err)
	}
	if s.DemoNavigation.Team != "Personal workspace" {
		t.Fatal("workspace selector did not change actual sample")
	}
	if err := tt.Click("Toggle Sample workspace"); err != nil {
		t.Fatal(err)
	}
	if !s.DemoNavigation.Sidebar.Collapsed {
		t.Fatal("external main-pane toggle did not collapse the sample navigation")
	}
	if !tt.HasText("Building your application") || !tt.HasText("Project activity") {
		t.Fatal("collapsing sidebar removed the adjacent main content")
	}
}

func TestNarrowLargeTextSidebarContentScrollsIndependently(t *testing.T) {
	defer theme.SetPreferences(theme.Preferences{})
	theme.SetPreferences(theme.Preferences{FontSize: 32})
	s := NewNavigationState()
	tt := ui.NewTester(func(c *ui.Context) { theme.ApplyPreset(c, theme.Neutral); navigationSample(c, &s, "sidebar", false) }, 342, 600)
	before, ok := tt.Find("Three documents are ready for review.")
	if !ok {
		t.Fatal("missing real content card")
	}
	header, _ := tt.Find("Building your application")
	tt.Scroll(180, 360, 0, 800)
	after, _ := tt.Find("Three documents are ready for review.")
	if after.Y >= before.Y {
		t.Fatalf("inner content did not scroll: before=%+v after=%+v", before, after)
	}
	headerAfter, _ := tt.Find("Building your application")
	if header != headerAfter {
		t.Fatal("inner content scroll moved fixed main header")
	}
}
