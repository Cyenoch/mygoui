package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDialogAndSheetShowCloseByDefault(t *testing.T) {
	for _, tc := range []struct {
		name  string
		build func(*ui.Context, *bool)
	}{
		{"dialog", func(c *ui.Context, open *bool) { Dialog(c, open, DialogOptions{Title: "Panel"}) }},
		{"sheet", func(c *ui.Context, open *bool) { Sheet(c, open, SheetOptions{Title: "Panel"}) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				tc.build(c, &open)
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			if err := tt.Click("Close"); err != nil {
				t.Fatal(err)
			}
			if open || !tt.Focused("Open") {
				t.Fatalf("close: open=%v opener focused=%v", open, tt.Focused("Open"))
			}
		})
	}
}

func TestSheetAndDrawerReducedMotionAndDismissal(t *testing.T) {
	for _, drawer := range []bool{false, true} {
		name := "sheet"
		if drawer {
			name = "drawer"
		}
		t.Run(name, func(t *testing.T) {
			var open, prevent bool
			dismissals := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				if drawer {
					Drawer(c, &open, DrawerOptions{Side: DrawerRight, Title: "Panel", ShowClose: true,
						PreventDismiss: prevent, OnDismiss: func() { dismissals++ }})
				} else {
					Sheet(c, &open, SheetOptions{Title: "Panel",
						PreventDismiss: prevent, OnDismiss: func() { dismissals++ }})
				}
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			prevent = true
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			first, ok := tt.Find("Panel")
			if !ok || tt.FramePending() {
				t.Fatal("reduced-motion surface did not open at rest")
			}
			tt.Advance(drawerMotionDuration)
			tt.Frame()
			settled, _ := tt.Find("Panel")
			if first != settled {
				t.Fatalf("reduced-motion surface moved: first=%+v settled=%+v", first, settled)
			}
			tt.Key(0, ui.KeyEscape)
			tt.ClickAt(4, 4)
			if !open || dismissals != 0 {
				t.Fatalf("prevented dismissal: open=%v dismissals=%d", open, dismissals)
			}
			if err := tt.Click("Close"); err != nil {
				t.Fatal(err)
			}
			if open || dismissals != 1 || !tt.Focused("Open") || tt.FramePending() {
				t.Fatalf("explicit close: open=%v dismissals=%d focused=%v pending=%v",
					open, dismissals, tt.Focused("Open"), tt.FramePending())
			}
		})
	}
}

func TestSheetGeometry(t *testing.T) {
	for _, tc := range []struct {
		name   string
		width  int
		side   mygoui.Side
		inset  bool
		custom float32
		want   ui.Rect
	}{
		{"default", 640, mygoui.SideAuto, false, 0, ui.Rect{X: 192, Y: 0, W: 448, H: 480}},
		{"narrow", 360, mygoui.SideRight, false, 0, ui.Rect{X: 48, Y: 0, W: 312, H: 480}},
		{"custom", 640, mygoui.SideLeft, false, 240, ui.Rect{X: 0, Y: 0, W: 240, H: 480}},
		{"inset", 640, mygoui.SideRight, true, 0, ui.Rect{X: 176, Y: 16, W: 448, H: 448}},
		{"inset on narrow window", 360, mygoui.SideLeft, true, 0, ui.Rect{X: 0, Y: 0, W: 312, H: 480}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				if p := Sheet(c, &open, SheetOptions{Side: tc.side, Width: tc.custom, Inset: tc.inset}); p.Valid() {
					panel = p.Bounds()
				}
			}, tc.width, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			if panel != tc.want {
				t.Fatalf("sheet geometry=%+v, want %+v", panel, tc.want)
			}
		})
	}
}

func TestSheetReopenDuringExit(t *testing.T) {
	var open bool
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		Sheet(c, &open, SheetOptions{
			Title: "Panel",
			Body:  func() { ui.Button(c, "Inside").OnClick(func() { actions++ }) },
		})
	}, 640, 480)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	motionSettle(tt, 2*dialogMotionDuration)
	settled, _ := tt.Find("Inside")
	tt.Key(0, ui.KeyEscape)
	tt.Advance(dialogMotionDuration / 3)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	motionSettle(tt, 2*dialogMotionDuration)
	reopened, ok := tt.Find("Inside")
	if !ok || reopened != settled {
		t.Fatalf("reopened sheet did not settle: got=%+v want=%+v", reopened, settled)
	}
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 || !open {
		t.Fatalf("reopened sheet action=%d open=%v", actions, open)
	}
	tt.Key(0, ui.KeyEscape)
	if open || !tt.Focused("Open") {
		t.Fatal("reopened sheet did not restore focus")
	}
}

func TestSheetDoesNotStrandScrimWhenDismissedDuringEntry(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Box(c).Fill().Background(ui.RGB(240, 240, 240)).Children(func() {
			ui.Button(c, "Open").OnClick(func() { open = true })
			Sheet(c, &open, SheetOptions{Title: "Panel"})
		})
	}, 640, 480)
	initial := tt.Image().RGBAAt(4, 400)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	motionSettle(tt, 2*dialogMotionDuration)
	if open || tt.HasText("Panel") || !tt.Focused("Open") {
		t.Fatal("sheet did not dismiss during entry")
	}
	if after := tt.Image().RGBAAt(4, 400); after != initial {
		t.Fatalf("scrim remained after dismissal: before=%v after=%v", initial, after)
	}
}

func TestEdgeSurfacesSlideAndRetainInertExits(t *testing.T) {
	for _, tc := range []struct {
		name       string
		side       mygoui.Side
		drawerSide DrawerSide
		dx         float32
		dy         float32
	}{
		{"left", mygoui.SideLeft, DrawerLeft, -1, 0},
		{"right", mygoui.SideRight, DrawerRight, 1, 0},
		{"top", mygoui.SideTop, DrawerTop, 0, -1},
		{"bottom", mygoui.SideBottom, DrawerBottom, 0, 1},
	} {
		for _, drawer := range []bool{false, true} {
			name := "sheet"
			if drawer {
				name = "drawer"
			}
			t.Run(name+"/"+tc.name, func(t *testing.T) {
				var open bool
				actions, dismissals := 0, 0
				duration := dialogMotionDuration
				if drawer {
					duration = drawerMotionDuration
				}
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Column(c).Fill().Background(theme.Colors(c).Background).Children(func() {
						ui.Button(c, "Open").OnClick(func() { open = true })
						marker := ui.RGB(230, 30, 30)
						if drawer {
							marker = ui.RGB(255, 0, 0)
						}
						body := func() {
							ui.Box(c).Size(32, 24).Background(marker)
							ui.Button(c, "Inside").OnClick(func() { actions++ })
						}
						if drawer {
							Drawer(c, &open, DrawerOptions{Side: tc.drawerSide, Title: "Panel",
								OnDismiss: func() { dismissals++ }, Body: body})
						} else {
							Sheet(c, &open, SheetOptions{Side: tc.side, Title: "Panel",
								OnDismiss: func() { dismissals++ }, Body: body})
						}
					})
				}, 640, 480)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				first, ok := tt.Find("Inside")
				if !ok {
					t.Fatal("sheet body is missing")
				}
				tt.Advance(duration / 2)
				tt.Frame()
				middle, _ := tt.Find("Inside")
				tt.Advance(duration)
				tt.Frame()
				settled, _ := tt.Find("Inside")
				distance := func(r ui.Rect) float32 {
					return (r.X-settled.X)*tc.dx + (r.Y-settled.Y)*tc.dy
				}
				if firstDistance, middleDistance := distance(first), distance(middle); firstDistance < 25 ||
					middleDistance <= 0 || middleDistance >= firstDistance {
					t.Fatalf("sheet did not slide into place: first=%+v middle=%+v settled=%+v", first, middle, settled)
				}
				redBounds := func() ui.Rect {
					img := tt.Image()
					minX, minY := img.Bounds().Max.X, img.Bounds().Max.Y
					maxX, maxY := -1, -1
					for y := 0; y < img.Bounds().Max.Y; y++ {
						for x := 0; x < img.Bounds().Max.X; x++ {
							pixel := img.RGBAAt(x, y)
							if int(pixel.R)-int(pixel.G) < 20 || int(pixel.R)-int(pixel.B) < 20 {
								continue
							}
							minX, minY = min(minX, x), min(minY, y)
							maxX, maxY = max(maxX, x), max(maxY, y)
						}
					}
					if maxX < 0 {
						return ui.Rect{}
					}
					return ui.Rect{X: float32(minX), Y: float32(minY), W: float32(maxX - minX + 1), H: float32(maxY - minY + 1)}
				}
				painted := redBounds()
				if painted.W == 0 {
					t.Fatal("sheet marker is not painted")
				}
				tt.Key(0, ui.KeyEscape)
				if open || dismissals != 1 || !tt.Focused("Open") {
					t.Fatalf("dismissal: open=%v dismissals=%d opener focused=%v", open, dismissals, tt.Focused("Open"))
				}
				tt.ClickAt(settled.X+settled.W/2, settled.Y+settled.H/2)
				if actions != 0 {
					t.Fatal("the departing snapshot accepted a click")
				}
				exitStep := duration / 2
				if drawer {
					exitStep = duration / 40
					if tc.drawerSide == DrawerLeft {
						exitStep = duration / 80
					}
				}
				tt.Advance(exitStep)
				tt.Frame()
				departing := redBounds()
				travel := (departing.X-painted.X)*tc.dx + (departing.Y-painted.Y)*tc.dy
				if departing.W == 0 || travel < 5 {
					t.Fatalf("exit snapshot did not move toward its edge: settled=%+v departing=%+v", painted, departing)
				}
				tt.Advance(duration)
				tt.Frame()
				if remaining := redBounds(); remaining.W != 0 {
					t.Fatalf("exit snapshot remained painted: %+v", remaining)
				}
			})
		}
	}
}
