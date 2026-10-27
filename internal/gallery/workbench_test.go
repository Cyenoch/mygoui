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
			s.Panel = 1
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
			if s.Selected != "composer" || s.NavOpen || s.Panel != 0 {
				t.Fatalf("selected=%q navigation open=%v panel=%d", s.Selected, s.NavOpen, s.Panel)
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

func TestWorkbenchViewsRetainPreviewState(t *testing.T) {
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
	if s.Panel != 1 || !tt.HasText("Go code") {
		t.Fatal("right arrow did not select the source view")
	}
	tt.Key(0, ui.KeyLeft)
	if s.Panel != 0 || !tt.HasText("Actions fired 1") {
		t.Fatal("returning to preview lost the live demo state")
	}
}

func TestWorkbenchViewFocusSurvivesResize(t *testing.T) {
	s := New()
	s.Selected = "button"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 1.75, ReduceMotion: true})
	if err := tt.Click("Code"); err != nil {
		t.Fatal(err)
	}
	tt.SetSize(360, 900)
	if !tt.Focused("Code") || s.Panel != 1 {
		t.Fatal("resizing lost the selected view or its native focus")
	}
	tt.Key(0, ui.KeyRight)
	if s.Panel != 2 || !tt.Focused("API") {
		t.Fatal("keyboard selection stopped working after the views wrapped")
	}
	tt.SetSize(1180, 900)
	if !tt.Focused("API") || s.Panel != 2 {
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

func TestDesktopCategoryAndSidebarKeyboard(t *testing.T) {
	s := New()
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 1180, 900)
	tt.Frame()
	if navHas(tt, "Button") || !navHas(tt, "Actions") {
		t.Fatal("navigation did not start with collapsed categories")
	}
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	if s.OpenGroup != "Actions" || !navHas(tt, "Button") {
		t.Fatal("category did not expand")
	}
	tt.Key(0, ui.KeySpace)
	if s.OpenGroup != "" || navHas(tt, "Button") {
		t.Fatal("Space did not collapse the focused category")
	}
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Button") {
		t.Fatal("Tab did not reach the first component")
	}
	tt.Key(0, ui.KeyEnter)
	if s.Selected != "button" {
		t.Fatal("Enter did not select the focused component")
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
			if err := tt.Click("Appearance"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeySpace)
			if s.AppearanceOpen {
				t.Fatal("keyboard could not close appearance controls")
			}
			tt.Key(0, ui.KeySpace)
			if !s.AppearanceOpen {
				t.Fatal("keyboard could not open appearance controls")
			}
			if err := tt.Click("Light"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("Dark") {
				t.Fatal("Dark is not reachable with Tab")
			}
			tt.Key(0, ui.KeySpace)
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("System") {
				t.Fatal("System font is not reachable with Tab")
			}
			tt.Key(0, ui.KeySpace)
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("Radius 4") {
				t.Fatal("radius is not reachable with Tab")
			}
			tt.Key(0, ui.KeySpace)
			prefs := theme.PreferencesSnapshot()
			if !s.Dark || !s.SystemFont || s.Radius != 4 || prefs.FontFamily != "" || prefs.Radius != 4 {
				t.Fatalf("appearance did not apply: state=%+v preferences=%+v", s, prefs)
			}
			if err := tt.Click("Cal Sans"); err != nil {
				t.Fatal(err)
			}
			if theme.PreferencesSnapshot().FontFamily != fonts.Family {
				t.Fatal("Cal Sans did not return")
			}
		})
	}
}

func TestNativeTitleBarInsetsClearControls(t *testing.T) {
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
					if !ok || title.H != max(float32(TitleBarHeight), bar.Height) {
						t.Fatalf("native title-bar height: %+v, insets %+v", title, bar)
					}
					brand, found := tt.Find("coss ui")
					if !found {
						t.Fatal("brand is missing")
					}
					for _, label := range []string{"coss ui", "MyGo"} {
						r, ok := tt.Find(label)
						if !ok && label == "MyGo" {
							continue
						}
						if !ok || r.X < 0 || r.X+r.W > float32(width) {
							t.Fatalf("%s leaves the window: %+v", label, r)
						}
						if r.Y < title.Y+title.H &&
							(r.X < bar.Left || r.X+r.W > float32(width)-bar.Right) {
							t.Fatalf("%s overlaps window controls: %+v, insets %+v", label, r, bar)
						}
					}
					if brand.Y < title.Y+title.H {
						center := brand.Y + brand.H/2
						if diff := center - (title.Y + title.H/2); diff < -1 || diff > 1 {
							t.Fatalf("brand is not centered in native title bar: %+v, %+v", brand, title)
						}
					}
					controls := []string{"Appearance", "Preview", "Code", "Toggle navigation"}
					compact := float32(width) < scale*docsCompactBelow
					if !compact {
						controls = append(controls, "Search components")
					}
					for _, label := range controls {
						r, ok := tt.Find(label)
						if !ok || r.Y < title.H || r.X < 0 || r.X+r.W > float32(width) {
							t.Fatalf("%s overlaps title-bar controls or leaves window: %+v, insets %+v", label, r, bar)
						}
					}
					if err := tt.Click("Appearance"); err != nil {
						t.Fatal(err)
					}
					if compact {
						if err := tt.Click("Toggle navigation"); err != nil {
							t.Fatal(err)
						}
					}
					for _, label := range []string{"Dark", "System", "Radius 16", "Search components"} {
						r, ok := tt.Find(label)
						if !ok || r.Y < title.H || r.X < 0 || r.X+r.W > float32(width) {
							t.Fatalf("%s overlaps title-bar controls or leaves window: %+v, insets %+v", label, r, bar)
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
	tt.SetPreferences(ui.Preferences{TextScale: 1.75, ReduceMotion: true})
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80})
	windowed, _ := tt.Find("coss ui")
	tt.SetTitleBar(ui.TitleBar{})
	fullscreen, _ := tt.Find("coss ui")
	if fullscreen.X >= windowed.X {
		t.Fatalf("fullscreen did not release the native inset: windowed=%+v fullscreen=%+v", windowed, fullscreen)
	}
	tt.SetSize(360, 900)
	tt.SetTitleBar(ui.TitleBar{Height: 52, Right: 138})
	brand, ok := tt.Find("coss ui")
	if !ok || brand.X+brand.W > 360-138 {
		t.Fatalf("resized window did not respect its new right inset: %+v", brand)
	}
	tt.Key(ui.Cmd, ui.KeyK)
	if !tt.Focused("Search components") {
		t.Fatal("search lost native focus after title-bar changes")
	}
}

func TestTitleBarKeepsVeryLargeTextOutOfNativeBand(t *testing.T) {
	s := New()
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 360, 900)
	tt.SetPreferences(ui.Preferences{TextScale: 2.5, ReduceMotion: true})
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80, Right: 138})
	title, _ := tt.Find("Gallery title bar")
	brand, ok := tt.Find("coss ui")
	if title.H != TitleBarHeight || !ok || brand.Y < title.Y+title.H ||
		brand.X < 0 || brand.X+brand.W > 360 {
		t.Fatalf("large brand overlaps the native band or window edge: title=%+v brand=%+v", title, brand)
	}
}

func TestCompactWorkbenchWithEnlargedText(t *testing.T) {
	const width = 360
	labels := []string{"Preview", "Code", "API", "Tokens"}
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
				if s.Panel != i {
					t.Fatalf("%s did not select its view", label)
				}
			}
			if err := tt.Click("Preview"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyRight)
			if s.Panel != 1 {
				t.Fatal("keyboard did not move from Preview to Code")
			}
		})
	}
}
