package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestVerticalGroupsMatchTheirWidestControl(t *testing.T) {
	for _, kind := range []string{"buttons", "toggles", "mixed"} {
		t.Run(kind, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Column(c).Fill().Padding(20).Children(func() {
					switch kind {
					case "buttons":
						mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Orientation: mygoui.GroupVertical},
							mygoui.ButtonOptions{Label: "A"}, mygoui.ButtonOptions{Label: "Very long item"})
					case "toggles":
						mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{Orientation: mygoui.Vertical},
							mygoui.ToggleButtonOptions{Label: "A", Variant: mygoui.ToggleOutline},
							mygoui.ToggleButtonOptions{Label: "Very long item", Variant: mygoui.ToggleOutline})
					case "mixed":
						items := make([]mygoui.GroupItem, 0, 2)
						for _, label := range []string{"A", "Very long item"} {
							items = append(items, mygoui.GroupItem{Build: func(corners [4]float32) ui.Element {
								return mygoui.Button(c, mygoui.ButtonOptions{Label: label, Corners: &corners})
							}})
						}
						mygoui.Group(c, mygoui.GroupOptions{Orientation: mygoui.Vertical}, items...)
					}
				})
			}, 300, 160)
			a, b := findRect(t, tt, "A"), findRect(t, tt, "Very long item")
			if a.W != b.W || a.X != b.X {
				t.Fatalf("vertical controls do not join: short=%+v long=%+v", a, b)
			}
			if b.W >= 260 {
				t.Fatalf("intrinsic group unexpectedly fills its parent: %+v", b)
			}
		})
	}
}

func TestVerticalGroupsJoinFixedWidthAndTextControls(t *testing.T) {
	for _, kind := range []string{"buttons", "toggles", "mixed"} {
		t.Run(kind, func(t *testing.T) {
			calls := 0
			selected := false
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Fill().Padding(20).Children(func() {
					switch kind {
					case "buttons":
						mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Orientation: mygoui.GroupVertical},
							mygoui.ButtonOptions{Label: "Short", IconOnly: true, Icon: lucide.Plus(), OnClick: func() { calls++ }},
							mygoui.ButtonOptions{Label: "Very long item"})
					case "toggles":
						mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{Orientation: mygoui.Vertical},
							mygoui.ToggleButtonOptions{Label: "Short", Selected: &selected, Variant: mygoui.ToggleOutline, IconOnly: true, Icon: lucide.Plus(), OnChange: func(bool) { calls++ }},
							mygoui.ToggleButtonOptions{Label: "Very long item", Variant: mygoui.ToggleOutline})
					case "mixed":
						mygoui.Group(c, mygoui.GroupOptions{Orientation: mygoui.Vertical},
							mygoui.GroupItem{Build: func(corners [4]float32) ui.Element {
								return mygoui.Button(c, mygoui.ButtonOptions{Label: "Short", Corners: &corners, OnClick: func() { calls++ }}).Width(80)
							}},
							mygoui.GroupItem{Build: func(corners [4]float32) ui.Element {
								return mygoui.Button(c, mygoui.ButtonOptions{Label: "Very long item", Corners: &corners})
							}})
					}
				})
			}, 300, 160)
			short, long := findRect(t, tt, "Short"), findRect(t, tt, "Very long item")
			if short.X != long.X || short.W != long.W || long.W >= 260 {
				t.Fatalf("fixed-width control breaks intrinsic vertical group: short=%+v long=%+v", short, long)
			}
			tt.ClickAt(short.X+short.W-8, short.Y+short.H/2)
			if calls != 1 {
				t.Fatal("expanded face is not interactive")
			}
		})
	}
}

func TestFramePanelsKeepInsetsAndIndependentActions(t *testing.T) {
	calls := 0
	var palette theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		palette = theme.Colors(c)
		ui.Column(c).Fill().Padding(20).Background(palette.Background).Children(func() {
			mygoui.Frame(c, mygoui.FrameOptions{Label: "Frame"}, func() {
				mygoui.FramePanel(c, func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "First", OnClick: func() { calls++ },
					})
				}).Label("First panel")
				mygoui.FramePanel(c, func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Second", OnClick: func() { calls += 10 },
					})
				}).Label("Second panel")
			})
		})
	}, 360, 260)
	frame := findRect(t, tt, "Frame")
	first := findRect(t, tt, "First panel")
	second := findRect(t, tt, "Second panel")
	if first.X-frame.X != 4 || first.Y-frame.Y != 4 ||
		frame.X+frame.W-first.X-first.W != 4 {
		t.Fatalf("frame panel inset: frame=%+v panel=%+v, want 4 DIP", frame, first)
	}
	if second.Y-first.Y-first.H != 4 || frame.Y+frame.H-second.Y-second.H != 4 {
		t.Fatalf("panel gap or bottom inset: frame=%+v panels=%+v %+v", frame, first, second)
	}
	action := findRect(t, tt, "First")
	// Native borders occupy layout space outside the panel's 20 DIP padding.
	if action.X-first.X != 21 || action.Y-first.Y != 21 ||
		first.Y+first.H-action.Y-action.H != 21 {
		t.Fatalf("panel content inset: panel=%+v action=%+v, want 20 DIP plus border", first, action)
	}
	img := tt.Image()
	muted := palette.MutedBackground.Alpha(0.72)
	wantPixel(t, img, int(frame.X+frame.W/2), int(frame.Y)+1,
		muted.Over(palette.Background), "muted frame inset")
	wantPixel(t, img, int(first.X)+10, int(first.Y+first.H/2), palette.Background, "panel face")
	wantPixel(t, img, int(first.X), int(first.Y+first.H/2),
		palette.Border.Over(palette.Background), "panel border")
	wantPixel(t, img, int(frame.X)+1, int(frame.Y)+1, palette.Background, "rounded outer corner")
	if err := tt.Click("First"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Second"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	if calls != 21 || !tt.Focused("Second") {
		t.Fatalf("panel actions calls=%d focus=%v, want 21 and Second focused", calls, tt.Focused("Second"))
	}
}

func TestFrameHeaderAndFooterKeepNativeLayout(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Frame(c, mygoui.FrameOptions{}, func() {
			mygoui.FrameHeader(c, func() {
				mygoui.FrameTitle(c, "Settings")
				mygoui.FrameDescription(c, "Account preferences")
			}).Label("Header")
			mygoui.FrameFooter(c, func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Apply"})
			}).Label("Footer")
		})
	}, 300, 180)
	header := findRect(t, tt, "Header")
	title := findRect(t, tt, "Settings")
	description := findRect(t, tt, "Account preferences")
	footer := findRect(t, tt, "Footer")
	action := findRect(t, tt, "Apply")
	if title.X-header.X != 20 || title.Y-header.Y != 16 ||
		description.Y != title.Y+title.H || header.Y+header.H-description.Y-description.H != 16 {
		t.Fatalf("header padding or text layout: header=%+v title=%+v description=%+v", header, title, description)
	}
	if title.H != 20 || description.H != 20 {
		t.Fatalf("frame text lines = %v and %v, want 20 DIP", title.H, description.H)
	}
	if action.X-footer.X != 20 || action.Y-footer.Y != 16 ||
		footer.Y+footer.H-action.Y-action.H != 16 {
		t.Fatalf("footer padding: footer=%+v action=%+v", footer, action)
	}
}

func TestGroupJoinsMixedControlsAndPreservesFocusAcrossReorder(t *testing.T) {
	for _, tc := range []struct {
		name        string
		orientation mygoui.Orientation
		direction   ui.Direction
	}{
		{name: "horizontal"},
		{name: "right to left", direction: ui.RightToLeft},
		{name: "vertical", orientation: mygoui.Vertical},
	} {
		t.Run(tc.name, func(t *testing.T) {
			value := ""
			reversed := false
			calls := 0
			var supplied [3][4]float32
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Column(c).Fill().Padding(24).Children(func() {
					items := []mygoui.GroupItem{
						{Key: "start", Build: func(corners [4]float32) ui.Element {
							supplied[0] = corners
							return mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Start", Corners: &corners, OnClick: func() { calls++ },
							})
						}},
						{Key: "query", Build: func(corners [4]float32) ui.Element {
							supplied[1] = corners
							return ui.TextInput(c, &value).Label("Query").Height(32).
								Radius(corners[0], corners[1], corners[2], corners[3])
						}},
						{Key: "end", Build: func(corners [4]float32) ui.Element {
							supplied[2] = corners
							return mygoui.Button(c, mygoui.ButtonOptions{
								Label: "End", Corners: &corners, OnClick: func() { calls += 10 },
							})
						}},
					}
					if reversed {
						items[0], items[2] = items[2], items[0]
					}
					mygoui.Group(c, mygoui.GroupOptions{
						Orientation: tc.orientation, FullWidth: true, Label: "Controls",
					}, items...)
				})
			}, 480, 200)
			tt.SetDirection(tc.direction)
			radius := theme.RadiusScale().LG
			wantCorners := [3][4]float32{{radius, 0, 0, radius}, {}, {0, radius, radius, 0}}
			if tc.orientation == mygoui.Vertical {
				wantCorners = [3][4]float32{{radius, radius, 0, 0}, {}, {0, 0, radius, radius}}
			} else if tc.direction == ui.RightToLeft {
				wantCorners[0], wantCorners[2] = wantCorners[2], wantCorners[0]
			}
			if supplied != wantCorners {
				t.Fatalf("supplied corners = %v, want %v", supplied, wantCorners)
			}
			rects := []ui.Rect{
				findRect(t, tt, "Start"), findRect(t, tt, "Query"), findRect(t, tt, "End"),
			}
			for i, rect := range rects {
				wantWidth := float32(144)
				if tc.orientation == mygoui.Vertical {
					wantWidth = 432
				}
				if rect.W != wantWidth || rect.H != 32 {
					t.Fatalf("control %d = %+v, want %v x 32 DIP", i, rect, wantWidth)
				}
				if i == 0 {
					continue
				}
				previous := rects[i-1]
				switch {
				case tc.orientation == mygoui.Vertical:
					if rect.Y != previous.Y+previous.H {
						t.Fatalf("vertical controls are not adjacent: %+v", rects)
					}
				case tc.direction == ui.RightToLeft:
					if rect.X+rect.W != previous.X {
						t.Fatalf("RTL controls are not adjacent: %+v", rects)
					}
				default:
					if rect.X != previous.X+previous.W {
						t.Fatalf("horizontal controls are not adjacent: %+v", rects)
					}
				}
			}
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("Start") {
				t.Fatal("first control did not receive keyboard focus")
			}
			tt.Key(0, ui.KeyEnter)
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("Query") {
				t.Fatal("Tab did not reach the native editor")
			}
			tt.Type("typed")
			reversed = true
			tt.Frame()
			if !tt.Focused("Query") {
				t.Fatal("keyed group reorder lost editor focus")
			}
			tt.Type("!")
			if value != "typed!" {
				t.Fatalf("editor value after reorder = %q, want typed!", value)
			}
			if err := tt.Click("End"); err != nil {
				t.Fatal(err)
			}
			if calls != 11 {
				t.Fatalf("independent group actions = %d, want 11", calls)
			}
		})
	}
}

func TestGroupUnkeyedControlsKeepDistinctNativeIdentity(t *testing.T) {
	var first, second bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Group(c, mygoui.GroupOptions{},
			mygoui.GroupItem{Build: func(corners [4]float32) ui.Element {
				return mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
					Label: "First", Selected: &first, Corners: &corners,
				})
			}},
			mygoui.GroupItem{Build: func(corners [4]float32) ui.Element {
				return mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
					Label: "Second", Selected: &second, Corners: &corners,
				})
			}},
		)
	}, 240, 100)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !first || second || !tt.Focused("First") {
		t.Fatalf("first activation states=%v %v focus=%v", first, second, tt.Focused("First"))
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if !first || !second || !tt.Focused("Second") {
		t.Fatalf("second activation states=%v %v focus=%v", first, second, tt.Focused("Second"))
	}
}
