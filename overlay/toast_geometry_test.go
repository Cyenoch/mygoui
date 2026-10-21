package overlay

import (
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

// These tests exercise the native viewport, selectable paragraphs and buttons.
// Geometry reports and optional PNGs describe the same frame that receives input.
func toastGeometryCapture(t *testing.T, tt *myui.Tester, name string) {
	t.Helper()
	dir := os.Getenv("MYGOUI_NOTIFICATION_GEOMETRY_CAPTURE_DIR")
	if dir == "" {
		dir = os.Getenv("MYGOUI_NOTIFICATION_CAPTURE_DIR")
	}
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "geometry-"+name+".png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	encodeErr := png.Encode(file, tt.Image())
	closeErr := file.Close()
	if encodeErr != nil {
		t.Fatal(encodeErr)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	t.Logf("notification geometry capture: %s", path)
}

func toastGeometryPixels(t *testing.T, tt *myui.Tester, r myui.Rect) image.Rectangle {
	t.Helper()
	img := tt.Image()
	region := image.Rect(
		int(math.Ceil(float64(r.X))),
		int(math.Ceil(float64(r.Y))),
		int(math.Floor(float64(r.X+r.W))),
		int(math.Floor(float64(r.Y+r.H))),
	).Intersect(img.Bounds())
	if region.Empty() {
		t.Fatalf("empty painted region: %+v", r)
	}
	background := img.RGBAAt(region.Min.X, region.Min.Y)
	ink := image.Rectangle{}
	for y := region.Min.Y; y < region.Max.Y; y++ {
		for x := region.Min.X; x < region.Max.X; x++ {
			p := img.RGBAAt(x, y)
			delta := math.Abs(float64(p.R)-float64(background.R)) +
				math.Abs(float64(p.G)-float64(background.G)) + math.Abs(float64(p.B)-float64(background.B))
			if delta > 80 {
				ink = ink.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}
	if ink.Empty() {
		t.Fatalf("content has geometry but no painted ink: %+v", r)
	}
	return ink
}

func toastGeometryInside(t *testing.T, outer, inner myui.Rect, label string) {
	t.Helper()
	if inner.W <= 0 || inner.H <= 0 || inner.X < outer.X-1 || inner.Y < outer.Y-1 ||
		inner.X+inner.W > outer.X+outer.W+1 || inner.Y+inner.H > outer.Y+outer.H+1 {
		t.Fatalf("%s outside card: card=%+v content=%+v", label, outer, inner)
	}
}

func TestToastGeometrySingleLineSharesCenterAndBalancedInsets(t *testing.T) {
	for _, tc := range []struct {
		name  string
		dark  bool
		scale float32
	}{
		{name: "light-16", scale: 1},
		{name: "dark-16", dark: true, scale: 1},
		{name: "dark-32", dark: true, scale: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newNotificationFixture(t, 640, 480)
			f.tt.SetDark(tc.dark)
			f.tt.SetPreferences(myui.Preferences{TextScale: tc.scale, ReduceMotion: true})
			closed := 0
			f.add(ToastOptions{ID: "center", Title: "Sample error", Kind: ToastError, Timeout: -1,
				OnClose: func() { closed++ }})
			card := notificationCard(t, f, "center")
			title := notificationRect(t, f.tt, "Sample error")
			icon := notificationRect(t, f.tt, "error")
			close := notificationRect(t, f.tt, "Close toast")
			center := title.Y + title.H/2
			for _, r := range []myui.Rect{icon, close} {
				if !notificationNear(center, r.Y+r.H/2) {
					t.Fatalf("single-line centers differ by >1 DIP: title=%+v icon=%+v close=%+v", title, icon, close)
				}
			}
			if !notificationNear(title.Y-card.Y, card.Y+card.H-title.Y-title.H) {
				t.Fatalf("uneven single-line vertical insets: card=%+v title=%+v", card, title)
			}
			if card.H > 2*theme.Space10+2+max(title.H, close.H)+1 {
				t.Fatalf("single line retains excess vertical padding: card=%+v title=%+v", card, title)
			}
			t.Logf("center report: card=%+v title=%+v icon=%+v close=%+v titleInk=%v iconInk=%v closeInk=%v",
				card, title, icon, close, toastGeometryPixels(t, f.tt, title),
				toastGeometryPixels(t, f.tt, icon), toastGeometryPixels(t, f.tt, close))
			toastGeometryCapture(t, f.tt, tc.name)
			notificationClick(t, f.tt, "Close toast")
			f.tt.Frame()
			if closed != 1 || len(f.state.ids) != 0 || f.tt.FramePending() {
				t.Fatalf("native close callback=%d ids=%q pending=%v", closed, f.state.ids, f.tt.FramePending())
			}
		})
	}
}

func TestToastGeometryWrappedTitleKeepsEveryLineAndIndependentContent(t *testing.T) {
	for _, tc := range []struct {
		name  string
		width int
		scale float32
		dark  bool
		rtl   bool
		title string
	}{
		{name: "words-390", width: 390, scale: 1, title: "Every word in this notification title must remain visible on its own wrapped line"},
		{name: "unicode-390-dark", width: 390, scale: 2, dark: true, title: "项目保存成功 👩🏽‍💻 é 完成通知"},
		{name: "token-200", width: 200, scale: 1, title: "UnbrokenNotificationIdentifier0123456789abcdef"},
		{name: "rtl-200-32", width: 200, scale: 2, rtl: true, dark: true, title: "تم حفظ 👩‍💻"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newNotificationFixture(t, tc.width, 1600)
			f.tt.SetDark(tc.dark)
			f.tt.SetPreferences(myui.Preferences{TextScale: tc.scale, ReduceMotion: true})
			if tc.rtl {
				f.tt.SetDirection(myui.RightToLeft)
			}
			actions, closed := 0, 0
			f.add(ToastOptions{ID: "wrap", Title: tc.title, Description: "Description remains separate.",
				Kind: ToastInfo, Action: "Open", Timeout: -1,
				OnAction: func() { actions++ }, OnClose: func() { closed++ }})
			card := notificationCard(t, f, "wrap")
			title := notificationRect(t, f.tt, tc.title)
			description := notificationRect(t, f.tt, "Description remains separate.")
			icon := notificationRect(t, f.tt, "info")
			close := notificationRect(t, f.tt, "Close toast")
			action := notificationRect(t, f.tt, "Open")
			var titleHeight, line float32
			f.pending = append(f.pending, func(c *myui.Context) {
				_, titleHeight = c.MeasureText(title.W, myui.Span{
					Text: tc.title, Font: c.Theme().Font, Size: theme.Scale(c, 16), Weight: 700,
				})
				_, line = c.MeasureText(0, myui.Span{Text: "M", Font: c.Theme().Font, Size: theme.Scale(c, 16)})
			})
			f.tt.Frame()
			row := max(line, max(icon.H, close.H))
			if titleHeight <= line || !notificationNear(title.H, titleHeight+max(float32(0), row-line)) {
				t.Fatalf("title lost wrapped lines: title=%+v measured=%g line=%g row=%g", title, titleHeight, line, row)
			}
			firstCenter := title.Y + row/2
			if !notificationNear(firstCenter, icon.Y+icon.H/2) || !notificationNear(firstCenter, close.Y+close.H/2) {
				t.Fatalf("controls do not align with first text line: title=%+v icon=%+v close=%+v row=%g", title, icon, close, row)
			}
			if description.Y < title.Y+title.H+theme.Space4-1 || action.Y < description.Y+description.H+theme.Space10-1 {
				t.Fatalf("independent content overlaps: title=%+v description=%+v action=%+v", title, description, action)
			}
			for label, r := range map[string]myui.Rect{"title": title, "description": description, "action": action, "close": close, "icon": icon} {
				toastGeometryInside(t, card, r, label)
			}
			lastLine := myui.Rect{X: title.X, Y: title.Y + title.H - line, W: title.W, H: line}
			t.Logf("wrap report: card=%+v title=%+v measuredHeight=%g line=%g lastLineInk=%v", card, title, titleHeight, line,
				toastGeometryPixels(t, f.tt, lastLine))
			toastGeometryCapture(t, f.tt, tc.name)
			// Native paragraph selection must copy the original Unicode string,
			// including hard-to-break tokens and joined emoji, without truncation.
			f.tt.ClickAt(title.X+title.W/2, title.Y+row/2)
			f.tt.Key(myui.Cmd, myui.KeyA)
			f.tt.Key(myui.Cmd, myui.KeyC)
			if got := f.tt.Clipboard(); got != tc.title {
				t.Fatalf("native title copy=%q want=%q", got, tc.title)
			}
			notificationClick(t, f.tt, "Open")
			f.tt.Key(0, myui.KeyEnter)
			if actions != 1 || closed != 1 || len(f.state.ids) != 0 {
				t.Fatalf("native callback ownership: actions=%d close=%d ids=%q", actions, closed, f.state.ids)
			}
		})
	}
}

func TestToastGeometryEmptyContentAndScaledLoading(t *testing.T) {
	for _, tc := range []struct {
		name, title, description, action string
	}{
		{name: "title", title: "Only title"},
		{name: "description", description: "Only description"},
		{name: "action", action: "Continue"},
		{name: "empty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newNotificationFixture(t, 390, 844)
			f.tt.SetPreferences(myui.Preferences{TextScale: 2, ReduceMotion: true})
			closed := 0
			f.add(ToastOptions{ID: "empty", Title: tc.title, Description: tc.description, Action: tc.action,
				Kind: ToastLoading, Timeout: -1, OnClose: func() { closed++ }})
			card := notificationCard(t, f, "empty")
			loading := notificationRect(t, f.tt, "Loading")
			close := notificationRect(t, f.tt, "Close toast")
			if !notificationNear(loading.W, 36) || !notificationNear(loading.H, 36) ||
				!notificationNear(loading.Y+loading.H/2, close.Y+close.H/2) {
				t.Fatalf("loading ignored actual scaled icon slot: loading=%+v close=%+v", loading, close)
			}
			toastGeometryInside(t, card, close, "close")
			for _, text := range []string{tc.title, tc.description, tc.action} {
				if text != "" {
					toastGeometryInside(t, card, notificationRect(t, f.tt, text), text)
				}
			}
			toastGeometryPixels(t, f.tt, loading)
			toastGeometryCapture(t, f.tt, "empty-"+tc.name)
			notificationClick(t, f.tt, "Close toast")
			if closed != 1 || f.tt.FramePending() {
				t.Fatalf("empty-content close=%d pending=%v", closed, f.tt.FramePending())
			}
		})
	}
}

func TestToastGeometryShortViewportScrollsFullMessageWithReachableClose(t *testing.T) {
	f := newNotificationFixture(t, 200, 180)
	f.tt.SetPreferences(myui.Preferences{TextScale: 2, ReduceMotion: true})
	f.tt.SetDirection(myui.RightToLeft)
	title := strings.Repeat("完整通知👩‍💻", 8)
	closed, actions := 0, 0
	f.add(ToastOptions{ID: "short", Title: title, Description: "End of the complete description", Kind: ToastWarning,
		Action: "Read", Timeout: -1, OnAction: func() { actions++ }, OnClose: func() { closed++ }})
	card := notificationCard(t, f, "short")
	close := notificationRect(t, f.tt, "Close toast")
	toastGeometryInside(t, myui.Rect{W: 200, H: 180}, card, "card")
	toastGeometryInside(t, card, close, "pinned close")
	before := notificationImage(f.tt)
	// The body owns native scrolling, while the close column stays pinned.
	x := card.X + card.W/2
	y := card.Y + card.H/2
	f.tt.Scroll(x, y, 0, 100000)
	if !notificationSameRect(close, notificationRect(t, f.tt, "Close toast")) {
		t.Fatal("scrolling long text moved its close control")
	}
	action := notificationRect(t, f.tt, "Read")
	toastGeometryInside(t, card, action, "scrolled action")
	if notificationPixelDifference(before, notificationImage(f.tt), before.Bounds()) < 40 {
		t.Fatal("native scroll did not reveal the end of the message")
	}
	toastGeometryPixels(t, f.tt, action)
	toastGeometryCapture(t, f.tt, "short-rtl-scrolled")
	notificationClick(t, f.tt, "Read")
	if actions != 1 || closed != 1 {
		t.Fatalf("scrolled native action=%d close=%d", actions, closed)
	}
	f.add(ToastOptions{ID: "short-close", Title: title, Timeout: -1, OnClose: func() { closed++ }})
	notificationClick(t, f.tt, "Close toast")
	if closed != 2 || len(f.state.ids) != 0 {
		t.Fatalf("short-viewport close=%d ids=%q", closed, f.state.ids)
	}
}

func TestToastGeometryTitleUpdatePreservesNativeSpringContinuity(t *testing.T) {
	f := newNotificationFixture(t, 390, 1000)
	f.add(ToastOptions{ID: "resize", Title: "Saved", Kind: ToastInfo, Timeout: -1})
	f.tt.Advance(2 * time.Second)
	before := notificationCard(t, f, "resize")
	title := "A changed title wraps over several lines while the existing native card smoothly adopts its new height"
	f.add(ToastOptions{ID: "resize", Title: title, Kind: ToastInfo, Timeout: -1})
	notificationAssertContinuous(t, before, notificationCard(t, f, "resize"))
	toastGeometryCapture(t, f.tt, "title-growth-first")
	f.tt.Advance(80 * time.Millisecond)
	// Advance can leave Now beyond its last presented refresh. Commit its
	// endpoint before comparing a same-time update with the current pose.
	f.tt.Frame()
	middle := notificationCard(t, f, "resize")
	if notificationSameRect(before, middle) {
		t.Fatal("title growth skipped its intermediate spring geometry")
	}
	toastGeometryCapture(t, f.tt, "title-growth-mid")
	f.add(ToastOptions{ID: "resize", Title: "Saved again", Kind: ToastSuccess, Timeout: -1})
	notificationAssertContinuous(t, middle, notificationCard(t, f, "resize"))
	toastGeometryCapture(t, f.tt, "title-reversal-first")
	f.tt.Advance(2 * time.Second)
	settled := notificationCard(t, f, "resize")
	if !notificationNear(before.H, settled.H) || f.tt.FramePending() {
		t.Fatalf("title reversal did not settle: before=%+v middle=%+v settled=%+v pending=%v",
			before, middle, settled, f.tt.FramePending())
	}
	t.Logf("title continuity: before=%+v middle=%+v settled=%+v", before, middle, settled)
	toastGeometryCapture(t, f.tt, "title-reversal-settled")
}

func TestToastGeometryVeryShortViewportKeepsCloseReachable(t *testing.T) {
	f := newNotificationFixture(t, 200, 60)
	f.tt.SetPreferences(myui.Preferences{TextScale: 2, ReduceMotion: true})
	closed := 0
	f.add(ToastOptions{ID: "tiny", Title: "A complete notification in a very short viewport", Kind: ToastError,
		Timeout: -1, OnClose: func() { closed++ }})
	card := notificationCard(t, f, "tiny")
	close := notificationRect(t, f.tt, "Close toast")
	toastGeometryInside(t, myui.Rect{W: 200, H: 60}, card, "very short card")
	toastGeometryInside(t, card, close, "very short close")
	toastGeometryPixels(t, f.tt, close)
	toastGeometryCapture(t, f.tt, "very-short-close")
	notificationClick(t, f.tt, "Close toast")
	if closed != 1 || len(f.state.ids) != 0 || f.tt.FramePending() {
		t.Fatalf("very short close callback=%d ids=%q pending=%v", closed, f.state.ids, f.tt.FramePending())
	}
}

// Opt-in captures exercise each host palette with native light/dark, text scale
// and direction. The same visible cards also prove containment and action input.
func TestToastGeometryCaptureMatrix(t *testing.T) {
	if os.Getenv("MYGOUI_NOTIFICATION_GEOMETRY_CAPTURE_DIR") == "" {
		t.Skip("set MYGOUI_NOTIFICATION_GEOMETRY_CAPTURE_DIR for the headless capture matrix")
	}
	for _, palette := range []struct {
		name   string
		preset theme.Preset
	}{
		{name: "neutral", preset: theme.Neutral},
		{name: "blue", preset: theme.Blue},
	} {
		for _, appearance := range []struct {
			name string
			dark bool
		}{
			{name: "light"},
			{name: "dark", dark: true},
		} {
			for _, text := range []struct {
				name  string
				scale float32
			}{
				{name: "16", scale: 1},
				{name: "32", scale: 2},
			} {
				for _, direction := range []struct {
					name  string
					value myui.Direction
					title string
				}{
					{name: "ltr", value: myui.LeftToRight, title: "Every word remains visible when this notification title wraps"},
					{name: "rtl", value: myui.RightToLeft, title: "تم حفظ المشروع وجميع التغييرات بنجاح 👩‍💻"},
				} {
					name := palette.name + "-" + appearance.name + "-" + text.name + "-" + direction.name
					t.Run(name, func(t *testing.T) {
						var state ToastState
						show := false
						actions, closed := 0, 0
						tt := myui.NewTester(func(c *myui.Context) {
							theme.ApplyPreset(c, palette.preset)
							if show {
								show = false
								Toast(c, &state, ToastOptions{ID: "matrix", Title: direction.title,
									Description: "Description remains independently readable.", Kind: ToastInfo, Action: "Open", Timeout: -1,
									OnAction: func() { actions++ }, OnClose: func() { closed++ }})
							}
							Toaster(c)
						}, 390, 1000)
						tt.SetDark(appearance.dark)
						tt.SetPreferences(myui.Preferences{TextScale: text.scale, ReduceMotion: true})
						tt.SetDirection(direction.value)
						show = true
						tt.Frame()
						card := notificationRect(t, tt, "Toast: matrix")
						toastGeometryInside(t, myui.Rect{W: 390, H: 1000}, card, "card")
						for _, label := range []string{direction.title, "Description remains independently readable.", "info", "Open", "Close toast"} {
							r := notificationRect(t, tt, label)
							toastGeometryInside(t, card, r, label)
							toastGeometryPixels(t, tt, r)
						}
						if direction.value == myui.RightToLeft && !notificationNear(card.X, theme.Space16) ||
							direction.value == myui.LeftToRight && !notificationNear(card.X+card.W, 390-theme.Space16) {
							t.Fatalf("notification missed its trailing edge: direction=%s card=%+v", direction.name, card)
						}
						toastGeometryCapture(t, tt, "matrix-"+name)
						notificationClick(t, tt, "Open")
						if actions != 1 || closed != 1 || len(state.ids) != 0 || tt.FramePending() {
							t.Fatalf("matrix native action=%d close=%d ids=%q pending=%v", actions, closed, state.ids, tt.FramePending())
						}
					})
				}
			}
		}
	}
}
