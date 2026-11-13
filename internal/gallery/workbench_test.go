package gallery

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/fonts"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestGlobalSearchKeyboardOpensLivePreview(t *testing.T) {
	for _, width := range []int{1180, 360} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			s := New()
			s.Selected = "button"
			tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, width, 900)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			tt.Key(ui.Cmd, ui.KeyK)
			if !tt.Focused("Search components") {
				t.Fatal("search shortcut did not focus the native editor")
			}
			tt.Type("no such component")
			if !tt.HasText("No components found") {
				t.Fatal("search did not show its empty state")
			}
			tt.Key(0, ui.KeyEnter)
			if s.Selected != "button" {
				t.Fatal("empty search changed the selected page")
			}
			tt.Key(0, ui.KeyEscape)
			if s.Search != "" || !tt.HasText("Foundations") {
				t.Fatal("Escape did not clear the search and restore navigation")
			}
			tt.Type("composer")
			if !tt.HasText("Composer") || tt.HasText("Typography") {
				t.Fatal("search did not filter navigation")
			}
			tt.Key(0, ui.KeyEnter)
			if s.Selected != "composer" || s.NavOpen {
				t.Fatalf("selected=%q navigation open=%v", s.Selected, s.NavOpen)
			}
			revealContent(t, tt, s, "Message draft", float32(width), 900)
			if err := tt.Click("Message draft"); err != nil {
				t.Fatal(err)
			}
			tt.Type("Hello from search")
			tt.Key(0, ui.KeyEnter)
			if got := cossStateOf(s, "composer").Last; got != "Hello from search" {
				t.Fatalf("live composer submitted %q", got)
			}
		})
	}
}

// leadView is the Preview or Code choice of the Button page's lead example.
func leadView(s *State) *exampleViewState {
	page, _ := componentPage("button")
	return StateOf[exampleViewState](s, "example:button:"+page.Examples[0].Title)
}

func TestExampleViewsRetainPreviewState(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 900, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	if err := tt.Click("Default"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Preview"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if leadView(s).Tab != 1 || tt.HasText("Actions fired 1") {
		t.Fatal("right arrow did not replace the preview with its source")
	}
	tt.Key(0, ui.KeyLeft)
	if leadView(s).Tab != 0 || !tt.HasText("Actions fired 1") {
		t.Fatal("returning to preview lost the live demo state")
	}
}

func TestExampleViewFocusSurvivesResize(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1.75, ReduceMotion: true})
	if err := tt.Click("Code"); err != nil {
		t.Fatal(err)
	}
	tt.SetSize(360, 900)
	if !tt.Focused("Code") || leadView(s).Tab != 1 {
		t.Fatal("resizing lost the selected view or its native focus")
	}
	tt.Key(0, ui.KeyLeft)
	if leadView(s).Tab != 0 || !tt.Focused("Preview") {
		t.Fatal("keyboard selection stopped working at the compact width")
	}
	tt.SetSize(1180, 900)
	if !tt.Focused("Preview") || leadView(s).Tab != 0 {
		t.Fatal("expanding lost the selected view or its native focus")
	}
}

func TestSearchSurvivesWindowResizing(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Type("comp")
	tt.SetSize(360, 900)
	tt.Key(ui.Cmd, ui.KeyK)
	if !tt.Focused("Search components") || s.Search != "comp" {
		t.Fatalf("compact search lost focus or text: %q", s.Search)
	}
	tt.Type("oser")
	tt.Key(0, ui.KeyEnter)
	if s.Selected != "composer" || s.NavOpen {
		t.Fatalf("compact result selected %q, navigation open=%v", s.Selected, s.NavOpen)
	}
	tt.SetSize(1180, 900)
	tt.Key(ui.Cmd, ui.KeyK)
	if !tt.Focused("Search components") || s.Search != "composer" {
		t.Fatalf("desktop search lost focus or text: %q", s.Search)
	}
	tt.Frame()
	tt.Frame()
}

func TestDesktopSidebarKeyboard(t *testing.T) {
	s := New()
	// The sidebar belongs to component pages; the landing page has none.
	s.Selected = "theming"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.Frame()
	if !navHas(tt, "Button") || !navHas(tt, "Typography") {
		t.Fatal("navigation does not list every group's pages")
	}
	if err := tt.Click("Button"); err != nil {
		t.Fatal(err)
	}
	if s.Selected != "button" {
		t.Fatal("clicking a page did not select it")
	}
	var next string
	for _, group := range pageGroups(filterPagesInOrder(Pages())) {
		for i, page := range group.pages {
			if page.ID == "button" && i+1 < len(group.pages) {
				next = group.pages[i+1].ID
			}
		}
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if s.Selected != next {
		t.Fatalf("Tab and Enter selected %q, want the next page %q", s.Selected, next)
	}
	if err := tt.Click("Toggle navigation"); err != nil {
		t.Fatal(err)
	}
	if !s.SidebarClosed || tt.HasText("Search components") {
		t.Fatal("desktop sidebar did not collapse")
	}
	tt.Key(0, ui.KeySpace)
	if s.SidebarClosed || !tt.HasText("Search components") {
		t.Fatal("Space did not reopen the sidebar")
	}
	tt.Key(0, ui.KeySpace)
	tt.Key(ui.Cmd, ui.KeyK)
	if s.SidebarClosed || !tt.Focused("Search components") {
		t.Fatal("search shortcut did not reopen and focus the sidebar")
	}
}

func TestAppearanceControlsAreKeyboardAccessible(t *testing.T) {
	for _, width := range []int{1180, 360} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			s := New()
			s.Selected = "button"
			tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, width, 900)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			view, _ := tt.Find("Copy Markdown")
			if err := tt.Click("Appearance"); err != nil {
				t.Fatal(err)
			}
			if !s.AppearanceOpen {
				t.Fatal("appearance controls did not open")
			}
			if moved, _ := tt.Find("Copy Markdown"); moved != view || view.W == 0 {
				t.Fatalf("opening appearance moved the page: %+v -> %+v", view, moved)
			}
			if err := tt.Click("Light"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyRight)
			if !s.Dark {
				t.Fatal("arrow keys did not choose Dark")
			}
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			if !s.SystemFont {
				t.Fatal("Tab and arrow keys did not reach the font choice")
			}
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyLeft)
			prefs := theme.PreferencesSnapshot()
			if s.Radius != 4 || prefs.FontFamily != "" || prefs.Radius != 4 {
				t.Fatalf("appearance did not apply: state=%+v preferences=%+v", s, prefs)
			}
			if err := tt.Click("Cal Sans"); err != nil {
				t.Fatal(err)
			}
			if theme.PreferencesSnapshot().FontFamily != fonts.Family {
				t.Fatal("Cal Sans did not return")
			}
			tt.Key(0, ui.KeyEscape)
			if s.AppearanceOpen {
				t.Fatal("Escape did not close appearance controls")
			}
		})
	}
}

func TestAppearanceHidesUnavailableFont(t *testing.T) {
	s := New()
	s.SystemFontOnly, s.AppearanceOpen = true, true
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	if tt.HasText("Cal Sans") || theme.PreferencesSnapshot().FontFamily != "" {
		t.Fatal("an unregistered font is offered or applied")
	}
}

// within reports whether r lies in the window and clear of the native
// control insets wherever it reaches into the title band.
func within(r ui.Rect, width float32, bar ui.TitleBar) bool {
	if r.X < 0 || r.X+r.W > width {
		return false
	}
	return r.Y >= bar.Height || (r.X >= bar.Left && r.X+r.W <= width-bar.Right)
}

func TestNativeTitleBarInsetsClearControls(t *testing.T) {
	skipInShort(t)
	for _, width := range []int{1180, 360} {
		for _, scale := range []float32{1, 1.75} {
			for _, bar := range []ui.TitleBar{
				{Height: 52, Left: 80},
				{Height: 52, Right: 138},
				{Height: 52, Left: 80, Right: 138},
				{Height: 70, Left: 80, Right: 138},
				{}, // Fullscreen has no native button insets.
			} {
				t.Run(fmt.Sprintf("%d-scale-%g-%v", width, scale, bar), func(t *testing.T) {
					s := New()
					s.Selected = "button"
					tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, width, 900)
					tt.SetPreferences(ui.Preferences{TextScale: scale, ReduceMotion: true})
					tt.SetTitleBar(bar)
					title, ok := tt.Find("Gallery title bar")
					if !ok || title.H < max(float32(TitleBarHeight), bar.Height) {
						t.Fatalf("title bar does not cover the native band: %+v, insets %+v", title, bar)
					}
					if side, ok := tt.Find("Sidebar title bar"); ok && side.H != title.H {
						t.Fatalf("sidebar and page title bars differ: %+v, %+v", side, title)
					}
					for _, label := range []string{"coss ui", "Toggle navigation", "GitHub", "Appearance"} {
						r, ok := tt.Find(label)
						if !ok && label == "coss ui" {
							continue
						}
						if !ok || !within(r, float32(width), bar) {
							t.Fatalf("%s overlaps window controls or leaves the window: %+v, insets %+v", label, r, bar)
						}
						// Controls that cannot fit between the window
						// controls are centered below the native band.
						band := max(float32(TitleBarHeight), bar.Height)
						want := title.H / 2
						if r.Y >= band {
							want = (band + title.H) / 2
						}
						if center := r.Y + r.H/2; label != "coss ui" && (center < want-1 || center > want+1) {
							t.Fatalf("%s is not centered in the title bar: %+v, %+v", label, r, title)
						}
					}
					controls := []string{"Preview", "Code"}
					compact := float32(width) < scale*docsCompactBelow
					if !compact {
						controls = append(controls, "Search components")
					}
					for _, label := range controls {
						r, ok := tt.Find(label)
						if !ok || r.Y < title.H || !within(r, float32(width), bar) {
							t.Fatalf("%s overlaps the title bar or leaves the window: %+v, insets %+v", label, r, bar)
						}
					}
					if err := tt.Click("Appearance"); err != nil {
						t.Fatal(err)
					}
					for _, label := range []string{"Dark", "System", "Large"} {
						r, ok := tt.Find(label)
						if !ok || r.Y < title.H || !within(r, float32(width), bar) {
							t.Fatalf("%s overlaps the title bar or leaves the window: %+v, insets %+v", label, r, bar)
						}
					}
					if compact {
						tt.Key(0, ui.KeyEscape)
						if err := tt.Click("Toggle navigation"); err != nil {
							t.Fatal(err)
						}
						for _, label := range []string{"coss ui", "Search components"} {
							r, ok := tt.Find(label)
							if !ok || !within(r, float32(width), bar) {
								t.Fatalf("%s in the navigation sheet overlaps window controls: %+v, insets %+v", label, r, bar)
							}
						}
					}
				})
			}
		}
	}
}

func TestTitleBarUpdatesOnFullscreenAndResize(t *testing.T) {
	s := New()
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80})
	windowed, _ := tt.Find("coss ui")
	tt.SetTitleBar(ui.TitleBar{})
	fullscreen, _ := tt.Find("coss ui")
	if windowed.W == 0 || fullscreen.X >= windowed.X {
		t.Fatalf("fullscreen did not release the native inset: windowed=%+v fullscreen=%+v", windowed, fullscreen)
	}
	tt.SetSize(360, 900)
	tt.SetTitleBar(ui.TitleBar{Height: 52, Right: 138})
	appearance, ok := tt.Find("Appearance")
	if !ok || appearance.X+appearance.W > 360-138 {
		t.Fatalf("resized window did not respect its new right inset: %+v", appearance)
	}
	tt.Key(ui.Cmd, ui.KeyK)
	if !tt.Focused("Search components") {
		t.Fatal("search lost native focus after title-bar changes")
	}
}

func TestTitleBarGrowsForVeryLargeText(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 360, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 2.5, ReduceMotion: true})
	bar := ui.TitleBar{Height: 52, Left: 80, Right: 138}
	tt.SetTitleBar(bar)
	title, _ := tt.Find("Gallery title bar")
	for _, label := range []string{"Toggle navigation", "Appearance"} {
		r, ok := tt.Find(label)
		if !ok || r.Y < title.Y || r.Y+r.H > title.Y+title.H || !within(r, 360, bar) {
			t.Fatalf("%s leaves the title bar or overlaps the native band: title=%+v control=%+v", label, title, r)
		}
	}
}

func TestCompactWorkbenchWithEnlargedText(t *testing.T) {
	const width = 360
	labels := []string{"Preview", "Code"}
	for _, scale := range []float32{1.75, 2.5} {
		t.Run(fmt.Sprint(scale), func(t *testing.T) {
			s := New()
			s.Selected = "button"
			textWidths := make(map[string]float32)
			tt := ui.NewTester(func(c *ui.Context) {
				View(c, s)
				for _, label := range labels {
					textWidths[label], _ = c.MeasureText(0, ui.Span{
						Text: label, Font: theme.FontFamily(c),
						Size: theme.Scale(c, theme.SM().Size), Weight: theme.WeightMedium,
					})
				}
			}, width, 900)
			tt.SetPreferences(ui.Preferences{TextScale: scale, ReduceMotion: true})
			tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80, Right: 138})
			for _, label := range append([]string{"Toggle navigation", "Appearance"}, labels...) {
				if _, isView := textWidths[label]; isView {
					revealContent(t, tt, s, label, width, 900)
				}
				r, ok := tt.Find(label)
				if !ok || r.X < 0 || r.X+r.W > width || r.Y < 0 || r.Y+r.H > 900 {
					t.Fatalf("%s is not visible at enlarged text: %+v", label, r)
				}
				if textWidth, ok := textWidths[label]; ok && r.W+0.5 < textWidth+20*scale {
					t.Fatalf("%s label is clipped: %+v, text width %v", label, r, textWidth)
				}
			}
			for i, label := range labels {
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
				if leadView(s).Tab != i {
					t.Fatalf("%s did not select its view", label)
				}
			}
			if err := tt.Click("Preview"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyRight)
			if leadView(s).Tab != 1 {
				t.Fatal("keyboard did not move from Preview to Code")
			}
		})
	}
}
