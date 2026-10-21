package overlay

import (
	"bytes"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
	"slices"
	"testing"
	"time"
)

func TestDialogTrapsFocusRestoresTriggerAndClosesOnlyTopOverlay(t *testing.T) {
	var open, child bool
	backgroundActions := 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Column(c).Padding(16).Gap(12).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Open dialog", OnClick: func() { open = true }})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Background", OnClick: func() { backgroundActions++ }})
			Dialog(c, &open, DialogOptions{Title: "Edit profile", HideClose: true}, func() {
				myui.Button(c, "First").AutoFocus()
				anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Details", OnClick: func() { child = true }})
				Popover(c, anchor, &child, PopoverOptions{Title: "Details panel"}, func() { myui.Button(c, "Inside popup").AutoFocus() })
			})
		})
	}, 640, 480)
	if err := tt.Click("Open dialog"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("First") {
		t.Fatal("dialog did not move focus to its first control")
	}
	tt.Key(myui.Shift, myui.KeyTab)
	if !tt.Focused("Details") {
		t.Fatal("Shift+Tab escaped dialog instead of wrapping")
	}
	tt.Key(0, myui.KeyTab)
	if !tt.Focused("First") {
		t.Fatal("Tab did not wrap inside modal")
	}
	if backgroundActions != 0 {
		t.Fatal("background activated under modal")
	}
	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, myui.KeyEscape)
	if child || !open {
		t.Fatalf("Escape closed wrong layer: dialog=%v popup=%v", open, child)
	}
	tt.Key(0, myui.KeyEscape)
	if open {
		t.Fatal("Escape did not close dialog")
	}
	if !tt.Focused("Open dialog") {
		t.Fatal("closing dialog did not restore its trigger focus")
	}
}

func TestAlertDialogOutsidePressAndEnterNeverImplicitlyAccept(t *testing.T) {
	open := false
	accepted, cancelled := 0, 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{Label: "Delete file", OnClick: func() { open = true }})
		AlertDialog(c, &open, AlertDialogOptions{Title: "Delete report?", ActionLabel: "Delete", Destructive: true,
			OnAction: func() { accepted++ }, OnCancel: func() { cancelled++ }}, nil)
	}, 640, 480)
	if err := tt.Click("Delete file"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Cancel") {
		t.Fatal("Cancel must receive initial focus")
	}
	tt.ClickAt(630, 470)
	if !open || accepted != 0 {
		t.Fatal("outside press dismissed or accepted alert")
	}
	tt.Key(0, myui.KeyEnter)
	if open || accepted != 0 || cancelled != 1 {
		t.Fatalf("Enter on Cancel: open=%v accepted=%d cancelled=%d", open, accepted, cancelled)
	}
	if err := tt.Click("Delete file"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, myui.KeyEscape)
	if open || accepted != 0 || cancelled != 2 {
		t.Fatalf("Escape: open=%v accepted=%d cancelled=%d", open, accepted, cancelled)
	}
}

func TestSheetPhysicalSidesAndFullHeight(t *testing.T) {
	for _, side := range []Side{SideLeft, SideRight} {
		t.Run(map[Side]string{SideLeft: "left", SideRight: "right"}[side], func(t *testing.T) {
			open := true
			var bounds myui.Rect
			tt := myui.NewTester(func(c *myui.Context) {
				theme.Apply(c)
				panel := Sheet(c, &open, SheetOptions{DialogOptions: DialogOptions{Title: "Filters"}, Side: side}, nil)
				bounds = panel.Bounds()
			}, 640, 480)
			tt.Frame()
			if bounds.H != 480 {
				t.Fatalf("side sheet height=%v, want480", bounds.H)
			}
			if side == SideLeft && bounds.X != 0 || side == SideRight && bounds.X+bounds.W != 640 {
				t.Fatalf("sheet missed physical edge: %+v", bounds)
			}
		})
	}
}

func TestDialogKeepsReadableMeasureAndVisibleActions(t *testing.T) {
	for _, viewport := range []struct {
		name          string
		width, height int
		scale         float32
	}{
		{name: "normal", width: 1280, height: 800, scale: 1},
		{name: "narrow", width: 320, height: 600, scale: 1},
		{name: "large text", width: 320, height: 600, scale: 2},
	} {
		t.Run(viewport.name, func(t *testing.T) {
			open := true
			var bounds myui.Rect
			var saved bool
			tt := myui.NewTester(func(c *myui.Context) {
				theme.Apply(c)
				panel := Dialog(c, &open, DialogOptions{
					Title: "Edit sample", HideClose: true,
					Footer: func() {
						mygoui.Button(c, mygoui.ButtonOptions{
							Label: "Save sample", OnClick: func() { saved = true },
						})
					},
				}, func() {
					for range 20 {
						myui.Text(c, "Scrollable sample content")
					}
				})
				bounds = panel.Bounds()
			}, viewport.width, viewport.height)
			tt.SetPreferences(myui.Preferences{TextScale: viewport.scale})
			tt.Frame()
			if bounds.W <= 0 || bounds.X <= 0 || bounds.X+bounds.W >= float32(viewport.width) {
				t.Fatalf("dialog lost its window inset: %+v", bounds)
			}
			if viewport.width > 1000 && bounds.W >= float32(viewport.width)/2 {
				t.Fatalf("dialog expanded beyond a readable measure: %+v", bounds)
			}
			footer, ok := tt.Find("Save sample")
			if !ok || footer.Y < bounds.Y || footer.Y+footer.H > bounds.Y+bounds.H {
				t.Fatalf("scrolling content displaced footer: panel=%+v footer=%+v", bounds, footer)
			}
			if err := tt.Click("Save sample"); err != nil {
				t.Fatal(err)
			}
			if !saved {
				t.Fatal("footer action was clipped or covered")
			}
		})
	}
}

func TestModalSurfacesClipChildrenAtExposedCorners(t *testing.T) {
	for _, surface := range []struct {
		name string
		side Side
	}{
		{name: "dialog", side: SideDefault},
		{name: "right sheet", side: SideRight},
		{name: "bottom drawer", side: SideBottom},
		{name: "top sheet", side: SideTop},
	} {
		t.Run(surface.name, func(t *testing.T) {
			open := true
			var bounds myui.Rect
			tt := myui.NewTester(func(c *myui.Context) {
				theme.Apply(c)
				o := DialogOptions{Title: "Surface", HideClose: true, Header: func() {
					// Consumer-supplied edge-to-edge content must respect the panel's corners.
					myui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).
						Background(myui.RGB(255, 0, 0)).PassThrough()
				}}
				var panel myui.Element
				switch surface.side {
				case SideDefault:
					panel = Dialog(c, &open, o, nil).Height(160)
				case SideBottom:
					panel = Drawer(c, &open, DrawerOptions{DialogOptions: o, Height: 160}, nil)
				default:
					panel = Sheet(c, &open, SheetOptions{DialogOptions: o, Side: surface.side}, nil)
					if surface.side == SideTop {
						panel.Height(160)
					}
				}
				bounds = panel.Bounds()
			}, 640, 480)
			tt.Frame()
			img := tt.Image()
			isContent := func(x, y float32) bool {
				pixel := img.RGBAAt(int(x), int(y))
				return pixel.R == 255 && pixel.G == 0 && pixel.B == 0
			}
			if !isContent(bounds.X+bounds.W/2, bounds.Y+bounds.H/2) {
				t.Fatal("custom header content did not paint inside panel")
			}
			switch surface.side {
			case SideDefault:
				if isContent(bounds.X, bounds.Y) {
					t.Fatal("dialog content escaped rounded corner")
				}
			case SideRight:
				if !isContent(bounds.X+bounds.W-1, bounds.Y) {
					t.Fatal("side sheet rounded its attached corner")
				}
			case SideBottom:
				if bounds.X != 0 || bounds.W != 640 || bounds.Y+bounds.H != 480 {
					t.Fatalf("drawer lost intentional full-width edge placement: %+v", bounds)
				}
				if isContent(bounds.X, bounds.Y) {
					t.Fatal("drawer content escaped exposed corner")
				}
				if !isContent(bounds.X, bounds.Y+bounds.H-1) {
					t.Fatal("drawer rounded its attached corner")
				}
			case SideTop:
				if isContent(bounds.X, bounds.Y+bounds.H-1) {
					t.Fatal("sheet content escaped exposed corner")
				}
				if !isContent(bounds.X, bounds.Y) {
					t.Fatal("sheet rounded its attached corner")
				}
			}
		})
	}
}

func TestModalExitRetainsPaintAfterExternalDataClears(t *testing.T) {
	for _, family := range []string{"dialog", "sheet", "drawer", "alert"} {
		t.Run(family, func(t *testing.T) {
			type content struct {
				text  string
				color myui.Color
			}
			data := &content{text: "First sample", color: myui.RGB(220, 40, 30)}
			open := false
			bodyCalls, actions := 0, 0
			var bounds myui.Rect
			tt := myui.NewTester(func(c *myui.Context) {
				theme.Apply(c)
				myui.Button(c, "Open sample").OnClick(func() { open = true })
				body := func() {
					bodyCalls++
					if data == nil {
						t.Fatal("closing modal rebuilt cleared caller data")
					}
					myui.Column(c).Height(64).Background(data.color).Children(func() {
						myui.Text(c, data.text)
						myui.Button(c, "Inside action").OnClick(func() { actions++ })
					})
				}
				o := DialogOptions{Title: "Sample", HideClose: true}
				var panel myui.Element
				switch family {
				case "dialog":
					panel = Dialog(c, &open, o, body)
				case "sheet":
					panel = Sheet(c, &open, SheetOptions{DialogOptions: o}, body)
				case "drawer":
					panel = Drawer(c, &open, DrawerOptions{DialogOptions: o}, body)
				case "alert":
					panel = AlertDialog(c, &open, AlertDialogOptions{Title: "Sample"}, body)
				}
				if panel.Valid() {
					bounds = panel.Bounds()
				}
			}, 640, 480)
			if err := tt.Click("Open sample"); err != nil {
				t.Fatal(err)
			}
			tt.Advance(time.Second)
			tt.Frame()
			before := slices.Clone(tt.Image().Pix)
			action, ok := tt.Find("Inside action")
			if !ok {
				t.Fatal("body action absent")
			}
			calls := bodyCalls
			open, data = false, nil
			tt.Frame()
			if bodyCalls != calls {
				t.Fatal("closed surface invoked body")
			}
			if !bytes.Equal(before, tt.Image().Pix) {
				t.Fatal("exit started with changed content or geometry")
			}
			if tt.HasText("Inside action") {
				t.Fatal("exit copy remained in interactive tree")
			}
			if !tt.Focused("Open sample") {
				t.Fatal("exit delayed opener focus restoration")
			}
			tt.ClickAt(action.X+action.W/2, action.Y+action.H/2)
			if actions != 0 {
				t.Fatal("exit copy accepted pointer input")
			}
			tt.Advance(80 * time.Millisecond)
			if bytes.Equal(before, tt.Image().Pix) {
				t.Fatal("exit did not animate")
			}
			if bounds.W <= 0 {
				t.Fatal("sample panel was never laid out")
			}
			data = &content{text: "Replacement sample", color: myui.RGB(30, 80, 220)}
			open = true
			tt.Frame()
			if !tt.HasText("Replacement sample") || tt.HasText("First sample") {
				t.Fatal("reopen reused stale presentation content")
			}
			for range 3 {
				open, data = false, nil
				tt.Frame()
				data = &content{text: "Latest sample", color: myui.RGB(30, 180, 80)}
				open = true
				tt.Frame()
			}
			open, data = false, nil
			tt.Frame()
			tt.Advance(time.Second)
			if tt.FramePending() {
				t.Fatal("closed surface retained animation work")
			}
		})
	}
}

func TestModalReducedMotionClosesImmediately(t *testing.T) {
	open := false
	builds, actions := 0, 0
	var markerBounds myui.Rect
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Button(c, "Open sample").OnClick(func() { open = true })
		Dialog(c, &open, DialogOptions{Title: "Sample", HideClose: true}, func() {
			builds++
			// Explicit caller paint remains distinguishable even when the panel
			// and the window share the same semantic surface color.
			marker := myui.Column(c).Height(96).Padding(8).Background(myui.RGB(220, 40, 30))
			marker.Children(func() {
				myui.Text(c, "Retained content")
				myui.Button(c, "Inside action").OnClick(func() { actions++ })
			})
			markerBounds = marker.Bounds()
		})
	}, 640, 480)
	tt.SetPreferences(myui.Preferences{TextScale: 1, ReduceMotion: true})
	empty := *tt.Image()
	empty.Pix = slices.Clone(empty.Pix)
	if err := tt.Click("Open sample"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	x, y := int(markerBounds.X+4), int(markerBounds.Y+markerBounds.H/2)
	before := tt.Image().RGBAAt(x, y)
	if markerBounds.W <= 8 || markerBounds.H <= 0 || before.R != 220 || before.G != 40 || before.B != 30 || before.A != 255 {
		t.Fatalf("sample did not establish visible caller paint: bounds=%+v pixel=%+v", markerBounds, before)
	}
	action, ok := tt.Find("Inside action")
	if !ok {
		t.Fatal("open sample action absent")
	}
	calls := builds
	open = false
	tt.Frame()
	after := tt.Image().RGBAAt(x, y)
	if after != empty.RGBAAt(x, y) {
		t.Fatalf("reduced-motion close left paint: before=%+v after=%+v empty=%+v", before, after, empty.RGBAAt(x, y))
	}
	if builds != calls || tt.FramePending() {
		t.Fatal("reduced-motion close retained work")
	}
	if tt.HasText("Retained content") || tt.HasText("Inside action") {
		t.Fatal("reduced-motion close retained an interactive subtree")
	}
	if !tt.Focused("Open sample") {
		t.Fatal("reduced-motion close lost opener focus")
	}
	tt.ClickAt(action.X+action.W/2, action.Y+action.H/2)
	if actions != 0 {
		t.Fatal("reduced-motion close retained a pointer target")
	}
}

func TestModalHeadingReservesCloseTargetAtNarrowLargeText(t *testing.T) {
	open := true
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		Dialog(c, &open, DialogOptions{Title: "Sample overlay"}, nil).Label("Sample dialog")
	}, 320, 600)
	tt.SetPreferences(myui.Preferences{TextScale: 2})
	title, hasTitle := tt.Find("Sample overlay")
	close, hasClose := tt.Find("Close")
	if !hasTitle || !hasClose {
		t.Fatal("heading or close target missing")
	}
	if title.X < close.X+close.W && title.X+title.W > close.X && title.Y < close.Y+close.H && title.Y+title.H > close.Y {
		t.Fatalf("heading overlaps close target: title=%+v close=%+v", title, close)
	}
}

func TestSheetRetainsNativeBodyScrollWithPinnedClose(t *testing.T) {
	open := true
	var body myui.ScrollState
	actions := 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		Sheet(c, &open, SheetOptions{Side: SideBottom, DialogOptions: DialogOptions{
			Title: "Scrollable sheet", BodyScroll: &body, BodyLabel: "Sheet body",
			Footer: func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Pinned close", OnClick: func() { open = false }})
			},
		}}, func() {
			myui.Column(c).FillWidth().PaddingX(8).Children(func() {
				for range 30 {
					myui.Text(c, "Exact retained source paragraph").Height(60)
				}
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Last action", OnClick: func() { actions++ }}).FillWidth()
			})
		})
	}, 390, 844)
	tt.SetPreferences(myui.Preferences{TextScale: 2, ReduceMotion: true})
	tt.Frame()
	pinned, ok := tt.Find("Pinned close")
	if !ok {
		t.Fatal("pinned Close missing")
	}
	for range 30 {
		viewport, ok := tt.Find("Sheet body")
		if !ok || viewport.H <= 0 {
			t.Fatalf("native body unavailable: %+v", viewport)
		}
		action, _ := tt.Find("Last action")
		if action.W > 0 && action.H > 0 && action.Y >= viewport.Y && action.Y+action.H <= viewport.Y+viewport.H {
			break
		}
		tt.Scroll(viewport.X+2, viewport.Y+viewport.H/2, 0, 160)
	}
	if body.Y <= 0 || body.MaxY <= 0 {
		t.Fatal("shared Sheet native body did not retain scroll observation")
	}
	after, _ := tt.Find("Pinned close")
	if after != pinned {
		t.Fatal("body scrolling moved pinned Close")
	}
	if err := tt.Click("Last action"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 {
		t.Fatal("native scrolled action unreachable")
	}
	if err := tt.Click("Pinned close"); err != nil {
		t.Fatal(err)
	}
	if open {
		t.Fatal("pinned Close did not dismiss scrolled Sheet")
	}
}
