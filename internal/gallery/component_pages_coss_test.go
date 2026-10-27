package gallery

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestJoinedExamplesKeepIntrinsicActionWidths(t *testing.T) {
	for _, tc := range []struct {
		page, example, action string
	}{
		{"group", "Mixed controls", "Find"},
		{"buttongroup", "Vertical", "Profile"},
	} {
		t.Run(tc.page, func(t *testing.T) {
			page, _ := componentPage(tc.page)
			var example ComponentExample
			for _, candidate := range page.Examples {
				if candidate.Title == tc.example {
					example = candidate
				}
			}
			if example.Demo == nil {
				t.Fatal("example missing")
			}
			s := New()
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Column(c).Fill().AlignItems(ui.Center).Children(func() { example.Demo(c, s) })
			}, 800, 240)
			r, ok := tt.Find(tc.action)
			if !ok || r.W > 160 {
				t.Fatalf("action stretched instead of using content width: %+v", r)
			}
		})
	}
}

func TestPressableIconUsesForegroundInk(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(fmt.Sprint(dark), func(t *testing.T) {
			var ink ui.Color
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, dark)
				ink = theme.Colors(c).MutedForeground
				ui.Column(c).Fill().Background(theme.Colors(c).Background).Children(func() {
					pressableBody(c, &buttonsState{})
				})
			}, 400, 150)
			r, ok := tt.Find("Press feedback icon")
			if !ok {
				t.Fatal("press feedback icon missing")
			}
			img := tt.Image()
			for y := int(r.Y); y < int(r.Y+r.H); y++ {
				for x := int(r.X); x < int(r.X+r.W); x++ {
					p := img.RGBAAt(x, y)
					if math.Abs(float64(p.R)-float64(ink.R)) < 12 {
						return
					}
				}
			}
			t.Fatal("icon has no legible foreground-colored pixels")
		})
	}
}

func TestDefaultGroupIsCenteredAndInteractive(t *testing.T) {
	page, ok := componentPage("group")
	if !ok || len(page.Examples) < 3 || page.Examples[0].Title != "Default" {
		t.Fatal("Group does not start with the default example")
	}
	for _, width := range []int{360, 1180} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			s := New()
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Fill().Background(theme.Colors(c).Background).Children(func() {
					page.Examples[0].Demo(c, s)
				})
			}, width, 300)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			group, ok := tt.Find("File actions")
			if !ok || math.Abs(float64(group.X+group.W/2-float32(width)/2)) > 1 {
				t.Fatalf("file actions are not centered: %+v", group)
			}
			files, filesOK := tt.Find("Files")
			media, mediaOK := tt.Find("Media")
			more, moreOK := tt.Find("More")
			if !filesOK || !mediaOK || !moreOK {
				t.Fatalf("missing actions: Files=%v Media=%v More=%v", filesOK, mediaOK, moreOK)
			}
			if media.X-files.X-files.W > 2 || more.X-media.X-media.W > 2 ||
				files.Y != media.Y || media.Y != more.Y {
				t.Fatalf("actions are not joined: Files=%+v Media=%+v More=%+v", files, media, more)
			}
			img := tt.Image()
			for _, seam := range []float32{media.X, more.X} {
				y := int(media.Y + media.H/2)
				line := img.RGBAAt(int(math.Round(float64(seam)))-1, y)
				face := img.RGBAAt(int(seam)+2, y)
				if line == face {
					t.Fatalf("joined outline is missing the separator at x=%v", seam)
				}
			}
			for _, label := range []string{"Files", "Media", "More"} {
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
				if got := cossStateOf(s, "group-default").Last; got != label {
					t.Fatalf("action %q recorded %q", label, got)
				}
			}
			if cossStateOf(s, "group").Last != "" {
				t.Fatal("default group changed the other group examples")
			}
			if err := tt.Click("Files"); err != nil {
				t.Fatal(err)
			}
			cossStateOf(s, "group-default").Last = ""
			tt.Key(0, ui.KeyEnter)
			if cossStateOf(s, "group-default").Last != "Files" {
				t.Fatal("focused group action did not activate with Enter")
			}
		})
	}
}

func TestDefaultSegmentedControlIsIndependent(t *testing.T) {
	page, ok := componentPage("segmentedcontrol")
	if !ok || len(page.Examples) < 3 || page.Examples[0].Title != "Default" {
		t.Fatal("SegmentedControl does not start with the default example")
	}
	s := New()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() { page.Examples[0].Demo(c, s) })
	}, 360, 300)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	billing, ok := tt.Find("Billing period")
	if !ok || math.Abs(float64(billing.X+billing.W/2-180)) > 1 {
		t.Fatalf("billing period is not centered: %+v", billing)
	}
	if err := tt.Click("Yearly"); err != nil {
		t.Fatal(err)
	}
	if cossStateOf(s, "segment-billing").Choice != 1 || cossStateOf(s, "segmentedcontrol").Choice != 0 {
		t.Fatal("billing selection did not change independently")
	}
	tt.Key(0, ui.KeyLeft)
	if cossStateOf(s, "segment-billing").Choice != 0 {
		t.Fatal("keyboard selection did not return to Monthly")
	}
}

func TestEveryPageRoutesThroughNativeNavigation(t *testing.T) {
	for _, width := range []int{640, 1180} {
		for _, page := range Pages() {
			t.Run(fmt.Sprintf("%d/%s", width, page.ID), func(t *testing.T) {
				s := New()
				s.Search = page.Title
				tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, width, 900)
				tt.Frame()
				if width < docsCompactBelow {
					if err := tt.Click("Toggle navigation"); err != nil {
						t.Fatal(err)
					}
					r, ok := tt.Find(page.Title)
					if !ok {
						t.Fatalf("missing navigation link %q", page.Title)
					}
					if r.Y+r.H > 440 {
						s.NavScroll.Y = r.Y - 200
						tt.Frame()
					}
				} else {
					revealSidebarEntry(t, tt, s, page.Title, float32(width), 900)
				}
				if err := tt.Click(page.Title); err != nil {
					t.Fatal(err)
				}
				tt.Frame()
				if s.Selected != page.ID || s.NavOpen {
					t.Fatalf("route = %q open=%v, want %q", s.Selected, s.NavOpen, page.ID)
				}
				if cp, ok := componentPage(page.ID); ok {
					for _, example := range cp.Examples {
						if _, ok := tt.Find("Example panel: " + example.Title); !ok {
							t.Fatalf("routed page lost example %q", example.Title)
						}
					}
					if sourceURL(page.ID) != "" {
						s.Panel = 1
						tt.Frame()
						if !tt.HasText("Coss source") {
							t.Fatal("source component lost its upstream reference")
						}
					}
				}
			})
		}
	}
}

func TestCossCatalogHasEveryUpstreamComponent(t *testing.T) {
	for id := range cossSources {
		cp, ok := componentPage(id)
		if !ok || len(cp.Examples) < 2 {
			t.Errorf("%s: upstream component needs a page with multiple specimens", id)
		}
	}
}

func TestCossCanonicalTokenResolver(t *testing.T) {
	cp, _ := componentPage("theming")
	for _, dark := range []bool{false, true} {
		palette := theme.ColorsForMode(dark)
		fields := reflect.ValueOf(palette)
		for _, token := range cp.Tokens {
			got, ok := paletteColor(&palette, token)
			field := fields.FieldByName(token)
			if !ok || !field.IsValid() || got != field.Interface().(ui.Color) {
				t.Errorf("dark=%v: %s does not resolve to its actual palette role", dark, token)
			}
			if tokenRoles[token] == "" {
				t.Errorf("%s has no role description", token)
			}
		}
	}
}

func TestCossExamplesUseNativeInputAndCallerState(t *testing.T) {
	for _, width := range []int{640, 1180} {
		newPage := func(t *testing.T, id string) (*State, *ui.Tester, func(string)) {
			s := New()
			s.Selected = id
			tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, width, 900)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			click := func(label string) {
				t.Helper()
				revealContent(t, tt, s, label, float32(width), 900)
				if err := tt.Click(label); err != nil {
					t.Fatal(err)
				}
			}
			return s, tt, click
		}
		t.Run(fmt.Sprintf("%d/frame", width), func(t *testing.T) {
			s, _, click := newPage(t, "frame")
			click("Save project")
			if cossStateOf(s, "frame").Count != 1 {
				t.Fatal("frame footer did not save")
			}
		})
		t.Run(fmt.Sprintf("%d/group", width), func(t *testing.T) {
			s, tt, click := newPage(t, "group")
			click("Repository query")
			tt.Type("mygoui")
			click("Find")
			if st := cossStateOf(s, "group"); st.Last != "mygoui" {
				t.Fatalf("joined editor/action = %q", st.Last)
			}
			click("Security")
			if cossStateOf(s, "group").Last != "Security" {
				t.Fatal("vertical group did not activate")
			}
		})
		t.Run(fmt.Sprintf("%d/command", width), func(t *testing.T) {
			s, tt, click := newPage(t, "command")
			click("Search actions")
			tt.Type("new")
			tt.Key(0, ui.KeyEnter)
			if st := cossStateOf(s, "command"); st.Last != "create" || st.Count != 1 {
				t.Fatalf("command did not execute filtered native highlight: %+v", st)
			}
			tt.Key(ui.Cmd, ui.KeyA)
			tt.Type("unmatched")
			if !tt.HasText("No matching actions") {
				t.Fatal("command empty state missing")
			}
		})
		t.Run(fmt.Sprintf("%d/contextmenu", width), func(t *testing.T) {
			s, tt, _ := newPage(t, "contextmenu")
			revealContent(t, tt, s, "Document actions", float32(width), 900)
			if err := tt.RightClick("Document actions"); err != nil {
				t.Fatal(err)
			}
			if err := tt.ChooseMenuItem("Duplicate"); err != nil {
				t.Fatal(err)
			}
			if cossStateOf(s, "contextmenu").Last != "duplicate" {
				t.Fatal("native context menu did not activate")
			}
		})
		t.Run(fmt.Sprintf("%d/segmentedcontrol", width), func(t *testing.T) {
			s, tt, click := newPage(t, "segmentedcontrol")
			click("Grid view")
			if cossStateOf(s, "segmentedcontrol").Choice != 1 || !tt.HasText("Current view: Grid") {
				t.Fatal("segment did not update composed content")
			}
			click("Map")
			if cossStateOf(s, "segmentedcontrol").Choice != 1 {
				t.Fatal("disabled segment changed selection")
			}
		})
		t.Run(fmt.Sprintf("%d/sidebar", width), func(t *testing.T) {
			s, _, click := newPage(t, "sidebar")
			click("Projects")
			click("Collapse workspace")
			st := cossStateOf(s, "sidebar")
			if st.Name != "projects" || !st.Checked {
				t.Fatalf("sidebar selected=%s collapsed=%v", st.Name, st.Checked)
			}
			click("Settings")
			if st.Name != "settings" {
				t.Fatal("collapsed sidebar lost native activation")
			}
		})
		t.Run(fmt.Sprintf("%d/composer", width), func(t *testing.T) {
			s, tt, click := newPage(t, "composer")
			click("Message draft")
			tt.Type("Hello")
			tt.Key(ui.Shift, ui.KeyEnter)
			tt.Type("Coss")
			tt.Key(0, ui.KeyEnter)
			st := cossStateOf(s, "composer")
			if st.Last != "Hello\nCoss" || st.Draft != "" || st.Count != 1 {
				t.Fatalf("composer sent=%q draft=%q count=%d", st.Last, st.Draft, st.Count)
			}
		})
		t.Run(fmt.Sprintf("%d/sheet", width), func(t *testing.T) {
			s, tt, click := newPage(t, "sheet")
			click("Open right sheet")
			if !cossStateOf(s, "sheet").Open[0] || !tt.HasText("Workspace settings") {
				t.Fatal("sheet did not open")
			}
			if err := tt.Click("Save settings"); err != nil {
				t.Fatal(err)
			}
			if st := cossStateOf(s, "sheet"); st.Open[0] || st.Count != 1 {
				t.Fatal("sheet did not close after save")
			}
			click("Open left sheet")
			tt.Key(0, ui.KeyEscape)
			if cossStateOf(s, "sheet").Open[1] {
				t.Fatal("left sheet did not dismiss")
			}
		})
		t.Run(fmt.Sprintf("%d/previewcard", width), func(t *testing.T) {
			s, tt, _ := newPage(t, "previewcard")
			r := revealContent(t, tt, s, "Ada Lovelace", float32(width), 900)
			tt.Move(r.X+r.W/2, r.Y+r.H/2)
			tt.Advance(300 * time.Millisecond)
			if err := tt.Click("Follow updates"); err != nil {
				t.Fatal(err)
			}
			if cossStateOf(s, "previewcard").PreviewHits != 1 {
				t.Fatal("preview content did not retain native hit testing")
			}
			tt.Key(0, ui.KeyEscape)
			if cossStateOf(s, "previewcard").Preview[0].Open {
				t.Fatal("preview did not dismiss")
			}
		})
		t.Run(fmt.Sprintf("%d/mask", width), func(t *testing.T) {
			s, tt, click := newPage(t, "scrollfade")
			click("Inspect build 02")
			if cossStateOf(s, "scrollfade").Count != 2 || !tt.Focused("Inspect build 02") {
				t.Fatal("mask lost native focus or hit testing")
			}
			r := revealContent(t, tt, s, "Masked activity", float32(width), 900)
			tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 120)
			if cossStateOf(s, "scrollfade").Scroll[3].Y <= 0 {
				t.Fatal("masked viewport did not scroll")
			}
		})
	}
}

func TestThemeOverridesReachEverySpecimen(t *testing.T) {
	s := New()
	if s.Radius != 10 || s.Primary != 0 {
		t.Fatal("gallery must start with Coss defaults")
	}
	s.Selected = "theming"
	tt := ui.NewTester(func(c *ui.Context) { View(c, s) }, 640, 900)
	tt.Frame()
	for _, label := range []string{"Radius 16", "Primary 2", "Dark"} {
		revealContent(t, tt, s, label, 640, 900)
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	prefs := theme.PreferencesSnapshot()
	if !s.Dark || s.Radius != 16 || s.Primary != 1 || prefs.Radius != 16 || prefs.Accent != primarySwatches[1] {
		t.Fatalf("theme overrides state=%+v preferences=%+v", s, prefs)
	}
	s.Selected = "button"
	s.Content.Y = 0
	tt.Frame()
	if got := theme.ColorsForMode(true).Primary; got != primarySwatches[1] {
		t.Fatalf("new page lost primary override: %v", got)
	}
}
