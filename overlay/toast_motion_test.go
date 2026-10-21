package overlay

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

// This harness sends model events through the real native viewport. It never
// retains build values, substitutes a toast engine, sleeps, or opens a window.
type notificationFixture struct {
	tt      *myui.Tester
	state   ToastState
	pending []func(*myui.Context)
	page    int
}

func newNotificationFixture(t *testing.T, width, height int) *notificationFixture {
	t.Helper()
	f := &notificationFixture{}
	f.tt = myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Button(c, "Page action").OnClick(func() { f.page++ })
		pending := f.pending
		f.pending = nil
		for _, event := range pending {
			event(c)
		}
		Toaster(c)
	}, width, height)
	return f
}

func (f *notificationFixture) add(o ToastOptions) {
	f.pending = append(f.pending, func(c *myui.Context) { Toast(c, &f.state, o) })
	f.tt.Frame()
}

func (f *notificationFixture) dismiss(id string) {
	f.pending = append(f.pending, func(c *myui.Context) { DismissToast(c, id) })
	f.tt.Frame()
}

func notificationRect(t *testing.T, tt *myui.Tester, label string) myui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("missing %q; native labels=%q", label, tt.Texts())
	}
	return r
}

func notificationClick(t *testing.T, tt *myui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}

func notificationFocus(t *testing.T, tt *myui.Tester, label string) {
	t.Helper()
	for range 12 {
		if tt.Focused(label) {
			return
		}
		tt.Key(0, myui.KeyTab)
	}
	t.Fatalf("Tab did not reach %q", label)
}

func notificationNear(a, b float32) bool { return math.Abs(float64(a-b)) <= 1 }

func notificationSameRect(a, b myui.Rect) bool {
	return notificationNear(a.X, b.X) && notificationNear(a.Y, b.Y) &&
		notificationNear(a.W, b.W) && notificationNear(a.H, b.H)
}

func notificationImage(tt *myui.Tester) *image.RGBA {
	frame := tt.Image()
	snapshot := image.NewRGBA(frame.Bounds())
	for y := frame.Bounds().Min.Y; y < frame.Bounds().Max.Y; y++ {
		copy(snapshot.Pix[snapshot.PixOffset(0, y):], frame.Pix[frame.PixOffset(0, y):frame.PixOffset(0, y)+frame.Bounds().Dx()*4])
	}
	return snapshot
}

func notificationPixelDifference(a, b *image.RGBA, r image.Rectangle) int {
	r = r.Intersect(a.Bounds()).Intersect(b.Bounds())
	n := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if a.RGBAAt(x, y) != b.RGBAAt(x, y) {
				n++
			}
		}
	}
	return n
}

func notificationCapture(t *testing.T, tt *myui.Tester, name string) {
	t.Helper()
	dir := os.Getenv("MYGOUI_NOTIFICATION_CAPTURE_DIR")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name+".png")
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
	t.Logf("notification capture: %s", path)
}

func TestToastNotificationArrivalKeepsPageFocusAndAnnouncesOnce(t *testing.T) {
	f := newNotificationFixture(t, 640, 600)
	notificationFocus(t, f.tt, "Page action")
	f.add(ToastOptions{ID: "arrival", Title: "Event created", Description: "The local sample has been saved.", Timeout: -1})
	if !f.tt.Focused("Page action") {
		t.Fatal("toast arrival stole page focus")
	}
	if got := f.tt.Announcements(); !slices.Equal(got, []string{"Event created", "The local sample has been saved."}) {
		t.Fatalf("arrival announcements=%q", got)
	}
	f.tt.Advance(time.Second)
	f.tt.Frame()
	if got := f.tt.Announcements(); len(got) != 0 {
		t.Fatalf("animation/rebuild repeated announcements: %q", got)
	}
	f.tt.Key(0, myui.KeyEscape)
	if !f.tt.HasText("Event created") {
		t.Fatal("page Escape dismissed a toast outside the focus scope")
	}
}

func TestToastNotificationSameIDUpdatePreservesFocusedAction(t *testing.T) {
	for _, result := range []struct {
		name string
		kind ToastKind
	}{{"success", ToastSuccess}, {"error", ToastError}} {
		t.Run(result.name, func(t *testing.T) {
			f := newNotificationFixture(t, 640, 600)
			oldClosed, newClosed, actions := 0, 0, 0
			f.add(ToastOptions{ID: "upload", Kind: ToastLoading, Title: "Uploading sample", Action: "Inspect upload", Timeout: -1,
				OnClose: func() { oldClosed++ }})
			f.tt.Advance(time.Second)
			notificationFocus(t, f.tt, "Inspect upload")
			f.tt.Announcements()
			f.add(ToastOptions{ID: "upload", Kind: result.kind, Title: "Upload " + result.name, Action: "Inspect upload", Timeout: time.Second,
				OnAction: func() { actions++ }, OnClose: func() { newClosed++ }})
			if !f.tt.Focused("Inspect upload") || f.tt.HasText("Uploading sample") {
				t.Fatal("same-ID update replaced focus identity or retained old content")
			}
			if got := f.tt.Announcements(); !slices.Equal(got, []string{"Upload " + result.name}) {
				t.Fatalf("update announcements=%q", got)
			}
			f.tt.Advance(3 * time.Second)
			if newClosed != 0 || !f.tt.HasText("Upload "+result.name) {
				t.Fatal("same-ID update lost focus pause")
			}
			f.tt.Key(0, myui.KeyEnter)
			f.tt.Advance(time.Second)
			if actions != 1 || newClosed != 1 || oldClosed != 0 {
				t.Fatalf("action=%d current close=%d superseded close=%d", actions, newClosed, oldClosed)
			}
		})
	}
}

func TestToastNotificationSameIDUpdatePreservesHoverPause(t *testing.T) {
	for _, result := range []struct {
		name string
		kind ToastKind
	}{{"success", ToastSuccess}, {"error", ToastError}} {
		t.Run(result.name, func(t *testing.T) {
			f := newNotificationFixture(t, 640, 600)
			closed := 0
			f.add(ToastOptions{ID: "hover", Kind: ToastLoading, Title: "Saving sample", Timeout: -1})
			f.tt.Advance(time.Second)
			r := notificationRect(t, f.tt, "Saving sample")
			f.tt.Move(r.X+r.W/2, r.Y+r.H/2)
			f.tt.Advance(time.Second)
			f.add(ToastOptions{ID: "hover", Kind: result.kind, Title: "Sample " + result.name, Timeout: time.Second, OnClose: func() { closed++ }})
			f.tt.Advance(3 * time.Second)
			if closed != 0 || !f.tt.HasText("Sample "+result.name) {
				t.Fatal("same-ID update restarted expiry beneath a stationary pointer")
			}
			f.tt.Move(0, 0)
			f.tt.Advance(2 * time.Second)
			if closed != 1 || f.tt.HasText("Sample "+result.name) {
				t.Fatalf("leaving hover close=%d texts=%q", closed, f.tt.Texts())
			}
		})
	}
}

func TestToastNotificationActionCloseOnceAndExternalDataClear(t *testing.T) {
	f := newNotificationFixture(t, 640, 600)
	title, description := "Source can disappear", "Retained exit keeps this complete message visible."
	var events []string
	f.add(ToastOptions{ID: "clear", Title: title, Description: description, Action: "Clear source", Timeout: -1,
		OnAction: func() { events = append(events, "action") },
		OnClose: func() {
			events = append(events, "close")
			title, description = "", ""
		}})
	f.tt.Advance(time.Second)
	button := notificationRect(t, f.tt, "Clear source")
	before := notificationImage(f.tt)
	notificationClick(t, f.tt, "Clear source")
	after := notificationImage(f.tt)
	if !slices.Equal(events, []string{"action", "close"}) || title != "" || description != "" {
		t.Fatalf("close did not clear source once after action: events=%q title=%q description=%q", events, title, description)
	}
	if f.tt.HasText("Source can disappear") || f.tt.HasText("Clear source") {
		t.Fatal("retained exit subtree remains in semantic lookup")
	}
	// Input over a retained exit cannot replay actions, including keyboard input.
	f.tt.ClickAt(button.X+button.W/2, button.Y+button.H/2)
	f.tt.Key(0, myui.KeyEnter)
	f.tt.Advance(time.Second)
	if !slices.Equal(events, []string{"action", "close"}) {
		t.Fatalf("retained exit control reactivated: %q", events)
	}
	empty := notificationImage(f.tt)
	region := image.Rect(200, 350, 640, 600)
	if notificationPixelDifference(after, empty, region) < 40 {
		t.Fatal("last toast disappeared without retained exit paint after source cleared")
	}
	if changed, painted := notificationPixelDifference(before, after, region), notificationPixelDifference(before, empty, region); changed > painted/3+100 {
		t.Fatalf("first exit frame lost painted content: changed=%d painted=%d", changed, painted)
	}
	if f.tt.FramePending() {
		t.Fatal("last exit kept requesting animation frames")
	}
}

func TestToastNotificationSimultaneousExpiryAndArrival(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	closed := 0
	f.add(ToastOptions{ID: "survivor", Title: "Survivor", Timeout: -1})
	f.add(ToastOptions{ID: "expires", Title: "Expires now", Timeout: time.Second, OnClose: func() {
		closed++
		f.pending = append(f.pending, func(c *myui.Context) {
			Toast(c, &f.state, ToastOptions{ID: "arrives", Title: "Arrived at expiry", Timeout: -1})
		})
	}})
	f.tt.Move(0, 0)
	f.tt.Advance(time.Second)
	f.tt.Frame()
	f.tt.Advance(time.Second)
	if closed != 1 || !f.tt.HasText("Survivor") || !f.tt.HasText("Arrived at expiry") || f.tt.HasText("Expires now") {
		t.Fatalf("same-timestamp lifecycle: close=%d texts=%q", closed, f.tt.Texts())
	}
	if len(f.state.ids) != 2 || !slices.Contains(f.state.ids, "survivor") || !slices.Contains(f.state.ids, "arrives") {
		t.Fatalf("retained IDs after simultaneous lifecycle=%q", f.state.ids)
	}
}

func TestToastNotificationCapacityUpdateDoesNotEvict(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	closed := make([]int, ToastCapacity+4)
	f.pending = append(f.pending, func(c *myui.Context) {
		for i := range ToastCapacity + 4 {
			Toast(c, &f.state, ToastOptions{ID: fmt.Sprintf("bounded-%d", i), Title: fmt.Sprintf("Bounded message %d", i), Timeout: -1,
				OnClose: func() { closed[i]++ }})
		}
	})
	f.tt.Frame()
	f.tt.Advance(time.Second)
	for i, count := range closed {
		want := 0
		if i < 4 {
			want = 1
		}
		if count != want {
			t.Fatalf("message %d close=%d want=%d", i, count, want)
		}
	}
	f.add(ToastOptions{ID: "bounded-35", Title: "Updated newest", Timeout: -1})
	f.tt.Advance(time.Second)
	if len(f.state.ids) != ToastCapacity || closed[4] != 0 {
		t.Fatalf("same-ID update evicted unrelated data: count=%d oldest surviving close=%d", len(f.state.ids), closed[4])
	}
	if f.tt.HasText("Bounded message 32") || !f.tt.HasText("Bounded message 33") || !f.tt.HasText("Bounded message 34") || !f.tt.HasText("Updated newest") {
		t.Fatalf("viewport did not retain exactly the newest three messages: %q", f.tt.Texts())
	}
}

func notificationCard(t *testing.T, f *notificationFixture, id string) myui.Rect {
	t.Helper()
	return notificationRect(t, f.tt, "Toast: "+id)
}

func notificationExpand(t *testing.T, f *notificationFixture, id string) {
	t.Helper()
	r := notificationCard(t, f, id)
	f.tt.Move(r.X+r.W/2, r.Y+r.H/2)
	f.tt.Advance(2 * time.Second)
}

func notificationAssertContinuous(t *testing.T, before, first myui.Rect) {
	t.Helper()
	if !notificationSameRect(before, first) {
		t.Fatalf("layout jumped on its first frame: before=%+v first=%+v", before, first)
	}
}

func TestToastMotionEntryAndUnequalHeightInsertion(t *testing.T) {
	f := newNotificationFixture(t, 640, 850)
	f.add(ToastOptions{ID: "short", Title: "Short card", Action: "Inspect short", Timeout: -1})
	f.tt.Advance(2 * time.Second)
	notificationFocus(t, f.tt, "Inspect short")
	f.tt.Advance(2 * time.Second)
	before := notificationCard(t, f, "short")
	f.add(ToastOptions{ID: "tall", Title: "Unequal height card", Description: "A much longer message wraps onto several lines so a stack cannot compute its layout from a fixed nominal card height. The complete description remains readable.", Timeout: -1})
	first := notificationCard(t, f, "short")
	entryFirst := notificationImage(f.tt)
	notificationCapture(t, f.tt, "unequal-insertion-first")
	notificationAssertContinuous(t, before, first)
	f.tt.Advance(80 * time.Millisecond)
	mid := notificationCard(t, f, "short")
	entryMid := notificationImage(f.tt)
	notificationCapture(t, f.tt, "unequal-insertion-mid")
	f.tt.Advance(2 * time.Second)
	end, tall := notificationCard(t, f, "short"), notificationCard(t, f, "tall")
	notificationCapture(t, f.tt, "unequal-insertion-settled")
	if tall.H <= end.H+10 {
		t.Fatalf("fixture did not exercise unequal heights: short=%+v tall=%+v", end, tall)
	}
	if !(end.Y < mid.Y && mid.Y < first.Y) {
		t.Fatalf("existing card did not move continuously for insertion: first=%+v mid=%+v end=%+v", first, mid, end)
	}
	if gap := tall.Y - (end.Y + end.H); math.Abs(float64(gap-10)) > 1 {
		t.Fatalf("expanded unequal-height gap=%v, want 10 DIP", gap)
	}
	if notificationPixelDifference(entryFirst, entryMid, image.Rect(180, 250, 640, 850)) < 40 {
		t.Fatal("entry first and middle frames paint identically")
	}
	if !f.tt.Focused("Inspect short") || f.tt.FramePending() {
		t.Fatal("insertion replaced focused identity or failed to settle")
	}
}

func TestToastMotionMiddleAndLastDismissalRetainInertExit(t *testing.T) {
	f := newNotificationFixture(t, 640, 850)
	closed, actions := 0, 0
	for _, o := range []ToastOptions{
		{ID: "old", Title: "Old card", Action: "Keep focused", Timeout: -1},
		{ID: "middle", Title: "Middle card", Description: "This card has a wrapped explanation with a different height.", Action: "Middle action", Timeout: -1, OnAction: func() { actions++ }, OnClose: func() { closed++ }},
		{ID: "new", Title: "Newest card", Timeout: -1},
	} {
		f.add(o)
	}
	f.tt.Advance(2 * time.Second)
	notificationFocus(t, f.tt, "Keep focused")
	f.tt.Advance(2 * time.Second)
	old := notificationCard(t, f, "old")
	button := notificationRect(t, f.tt, "Middle action")
	f.dismiss("middle")
	notificationCapture(t, f.tt, "middle-exit-first")
	notificationAssertContinuous(t, old, notificationCard(t, f, "old"))
	if f.tt.HasText("Middle card") || closed != 1 {
		t.Fatal("middle close retained a live semantic subtree or lost callback")
	}
	f.tt.ClickAt(button.X+button.W/2, button.Y+button.H/2)
	f.tt.Advance(80 * time.Millisecond)
	notificationCapture(t, f.tt, "middle-exit-mid")
	f.tt.Advance(2 * time.Second)
	if actions != 0 || closed != 1 {
		t.Fatalf("exit snapshot accepted pointer: actions=%d close=%d", actions, closed)
	}
	f.dismiss("old")
	f.tt.Advance(2 * time.Second)
	f.tt.Move(0, 0)
	f.tt.Advance(2 * time.Second)
	before := notificationImage(f.tt)
	f.dismiss("new")
	first := notificationImage(f.tt)
	notificationCapture(t, f.tt, "last-exit-first")
	f.tt.Advance(80 * time.Millisecond)
	notificationCapture(t, f.tt, "last-exit-mid")
	f.tt.Advance(2 * time.Second)
	empty := notificationImage(f.tt)
	region := image.Rect(180, 500, 640, 850)
	painted := notificationPixelDifference(before, empty, region)
	if notificationPixelDifference(before, first, region) > painted/3+100 || notificationPixelDifference(first, empty, region) < 40 {
		t.Fatal("last exit lost visual continuity or failed to retain painted content")
	}
	if f.tt.FramePending() {
		t.Fatal("empty notification deck keeps requesting animation frames")
	}
}

func TestToastMotionHoverGapKeepsExpansionAndPauseStable(t *testing.T) {
	f := newNotificationFixture(t, 640, 850)
	closed := 0
	for i := range 3 {
		f.add(ToastOptions{ID: fmt.Sprintf("gap-%d", i), Title: fmt.Sprintf("Gap card %d", i), Timeout: 10 * time.Second, OnClose: func() { closed++ }})
	}
	f.tt.Advance(2 * time.Second)
	notificationExpand(t, f, "gap-2")
	upper, lower := notificationCard(t, f, "gap-1"), notificationCard(t, f, "gap-2")
	if gap := lower.Y - upper.Y - upper.H; gap < 8 || gap > 12 {
		t.Fatalf("hover did not create independent cards and real gap: upper=%+v lower=%+v", upper, lower)
	}
	x, y := lower.X+lower.W/2, (upper.Y+upper.H+lower.Y)/2
	f.tt.Move(x, y)
	f.tt.Advance(20 * time.Second)
	for range 5 {
		f.tt.Frame()
		if !notificationSameRect(lower, notificationCard(t, f, "gap-2")) {
			t.Fatal("stationary inter-card pointer oscillated expansion geometry")
		}
	}
	if closed != 0 || f.tt.FramePending() {
		t.Fatalf("inter-card hover lost expiry pause or kept animating: closed=%d pending=%v", closed, f.tt.FramePending())
	}
	f.tt.Move(0, 0)
	f.tt.Advance(12 * time.Second)
	if closed != 3 || f.tt.HasText("Gap card") {
		t.Fatalf("leaving gap did not resume expiry: closed=%d texts=%q", closed, f.tt.Texts())
	}
}

func TestToastMotionDragCancelCommitAndSameIDUpdate(t *testing.T) {
	for _, result := range []struct {
		name string
		kind ToastKind
	}{{"success", ToastSuccess}, {"error", ToastError}} {
		t.Run(result.name, func(t *testing.T) {
			f := newNotificationFixture(t, 640, 700)
			oldClosed, newClosed := 0, 0
			f.add(ToastOptions{ID: "drag", Title: "Drag loading", Kind: ToastLoading, OnClose: func() { oldClosed++ }})
			f.tt.Advance(2 * time.Second)
			notificationExpand(t, f, "drag")
			rest := notificationCard(t, f, "drag")
			// The border strip is card background, outside selectable text/buttons.
			x, y := rest.X+rest.W/2, rest.Y+3
			f.tt.Press(x, y)
			f.tt.Advance(100 * time.Millisecond)
			f.tt.Move(x+24, y)
			f.tt.Advance(200 * time.Millisecond)
			dragged := notificationCard(t, f, "drag")
			if dragged.X <= rest.X+5 {
				t.Fatalf("background drag did not follow pointer: rest=%+v dragged=%+v", rest, dragged)
			}
			f.add(ToastOptions{ID: "drag", Title: "Drag " + result.name, Kind: result.kind, Timeout: -1, OnClose: func() { newClosed++ }})
			updated := notificationCard(t, f, "drag")
			if !notificationNear(updated.X, dragged.X) {
				t.Fatalf("same-ID update reset drag pose: before=%+v after=%+v", dragged, updated)
			}
			notificationCapture(t, f.tt, "drag-update-"+result.name)
			f.tt.Release(x+24, y)
			f.tt.Advance(2 * time.Second)
			if got := notificationCard(t, f, "drag"); !notificationNear(got.X, rest.X) || oldClosed != 0 || newClosed != 0 {
				t.Fatalf("cancel did not return to rest with new payload: card=%+v old=%d new=%d", got, oldClosed, newClosed)
			}
			rest = notificationCard(t, f, "drag")
			x, y = rest.X+rest.W/2, rest.Y+3
			f.tt.Press(x, y)
			f.tt.Move(x+rest.W*0.6, y)
			f.tt.Release(x+rest.W*0.6, y)
			f.tt.Advance(2 * time.Second)
			if f.tt.HasText("Drag "+result.name) || newClosed != 1 || oldClosed != 0 || f.tt.FramePending() {
				t.Fatalf("committed drag lifecycle: old=%d new=%d pending=%v texts=%q", oldClosed, newClosed, f.tt.FramePending(), f.tt.Texts())
			}
		})
	}
}

func TestToastMotionButtonGestureDoesNotSwipeOrActivate(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	actions, closed := 0, 0
	f.add(ToastOptions{ID: "button", Title: "Button owns gesture", Action: "Keep action", Timeout: -1,
		OnAction: func() { actions++ }, OnClose: func() { closed++ }})
	f.tt.Advance(2 * time.Second)
	notificationExpand(t, f, "button")
	card := notificationCard(t, f, "button")
	b := notificationRect(t, f.tt, "Keep action")
	x, y := b.X+b.W/2, b.Y+b.H/2
	f.tt.Press(x, y)
	f.tt.Move(x-180, y)
	f.tt.Release(x-180, y)
	f.tt.Advance(2 * time.Second)
	if actions != 0 || closed != 0 || !notificationSameRect(card, notificationCard(t, f, "button")) {
		t.Fatalf("button drag activated or moved card: action=%d close=%d", actions, closed)
	}
	notificationClick(t, f.tt, "Keep action")
	f.tt.Advance(2 * time.Second)
	if actions != 1 || closed != 1 {
		t.Fatalf("ordinary action after cancelled button drag: action=%d close=%d", actions, closed)
	}
}

func TestToastMotionRapidExpansionReversalStartsAtCurrentPose(t *testing.T) {
	f := newNotificationFixture(t, 640, 850)
	for i := range 3 {
		f.add(ToastOptions{ID: fmt.Sprintf("reverse-%d", i), Title: fmt.Sprintf("Reversal card %d", i), Timeout: -1})
	}
	f.tt.Advance(2 * time.Second)
	notificationCapture(t, f.tt, "stack-three-collapsed")
	collapsed := notificationCard(t, f, "reverse-0")
	newest := notificationCard(t, f, "reverse-2")
	f.tt.Move(newest.X+newest.W/2, newest.Y+newest.H/2)
	f.tt.Advance(80 * time.Millisecond)
	// Advance leaves Now at its endpoint, which need not be a refresh time.
	// Commit that time before comparing input at the same instant.
	f.tt.Frame()
	notificationCapture(t, f.tt, "stack-three-mid")
	part := notificationCard(t, f, "reverse-0")
	if notificationSameRect(collapsed, part) {
		t.Fatal("hover expansion never entered an intermediate pose")
	}
	f.tt.Move(0, 0)
	notificationAssertContinuous(t, part, notificationCard(t, f, "reverse-0"))
	f.tt.Advance(40 * time.Millisecond)
	f.tt.Frame()
	part = notificationCard(t, f, "reverse-0")
	newest = notificationCard(t, f, "reverse-2")
	f.tt.Move(newest.X+newest.W/2, newest.Y+newest.H/2)
	notificationAssertContinuous(t, part, notificationCard(t, f, "reverse-0"))
	notificationCapture(t, f.tt, "reversal")
	f.tt.Advance(2 * time.Second)
	notificationCapture(t, f.tt, "stack-three-expanded")
	f.tt.Move(0, 0)
	f.tt.Advance(2 * time.Second)
	if !notificationSameRect(collapsed, notificationCard(t, f, "reverse-0")) || f.tt.FramePending() {
		t.Fatal("rapid hover reversal did not settle back to the original collapsed pose")
	}
}

func TestToastNotificationNarrowLargeTextThemesAndRTL(t *testing.T) {
	for _, tc := range []struct {
		name string
		dark bool
		rtl  bool
	}{{name: "light"}, {name: "dark", dark: true}, {name: "rtl", dark: true, rtl: true}} {
		t.Run(tc.name, func(t *testing.T) {
			f := newNotificationFixture(t, 390, 844)
			f.tt.SetDark(tc.dark)
			f.tt.SetPreferences(myui.Preferences{TextScale: 2, ReduceMotion: true})
			if tc.rtl {
				f.tt.SetDirection(myui.RightToLeft)
			}
			actions := 0
			f.add(ToastOptions{ID: "narrow", Title: "Event created", Description: "The local sample has been saved. A wrapped description remains readable at two hundred percent text size.", Action: "Undo", Timeout: -1, OnAction: func() { actions++ }})
			card := notificationCard(t, f, "narrow")
			for _, label := range []string{"Toast: narrow", "Event created", "The local sample has been saved. A wrapped description remains readable at two hundred percent text size.", "Undo", "Close toast"} {
				r := notificationRect(t, f.tt, label)
				if r.X < -0.5 || r.Y < -0.5 || r.X+r.W > 390.5 || r.Y+r.H > 844.5 || r.W <= 0 || r.H <= 0 {
					t.Fatalf("%q falls outside 390-DIP viewport at 200%%: %+v", label, r)
				}
			}
			closeButton := notificationRect(t, f.tt, "Close toast")
			action := notificationRect(t, f.tt, "Undo")
			if closeButton.X < action.X+action.W && closeButton.X+closeButton.W > action.X && closeButton.Y < action.Y+action.H && closeButton.Y+closeButton.H > action.Y {
				t.Fatalf("action and close overlap: action=%+v close=%+v", action, closeButton)
			}
			if tc.rtl && card.X > 17 || !tc.rtl && card.X+card.W < 373 {
				t.Fatalf("toast did not use trailing window edge: RTL=%v card=%+v", tc.rtl, card)
			}
			notificationCapture(t, f.tt, "narrow-"+tc.name+"-200")
			notificationClick(t, f.tt, "Undo")
			if actions != 1 || f.tt.HasText("Event created") || f.tt.FramePending() {
				t.Fatal("large-text action failed or reduced-motion exit retained animation")
			}
		})
	}
}

func TestToastNotificationReducedMotionLoadingAndDirectNativeToasts(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	f.tt.SetPreferences(myui.Preferences{ReduceMotion: true})
	f.add(ToastOptions{ID: "sticky", Title: "Sticky loading", Kind: ToastLoading})
	f.tt.Move(0, 0)
	f.tt.Advance(10 * time.Second)
	if !f.tt.HasText("Sticky loading") || f.tt.FramePending() {
		t.Fatal("loading default expired or reduced-motion indicator keeps looping")
	}
	f.pending = append(f.pending, func(c *myui.Context) {
		c.AddToast(myui.Toast{ID: "direct", Title: "Native toast", Timeout: -1})
	})
	f.tt.Frame()
	if !f.tt.HasText("Native toast") {
		t.Fatal("Toaster failed to render direct native AddToast")
	}
	f.add(ToastOptions{ID: "sticky", Title: "Loading complete", Kind: ToastSuccess})
	f.tt.Advance(5 * time.Second)
	if f.tt.HasText("Loading complete") || !f.tt.HasText("Native toast") {
		t.Fatal("same-ID loading completion did not restore native default expiration")
	}
	f.dismiss("direct")
	notificationCapture(t, f.tt, "reduced-motion")
	if f.tt.HasText("Native toast") || f.tt.FramePending() {
		t.Fatal("reduced-motion final exit did not finish immediately")
	}
}

func TestToastNotificationFocusedUpdateEscapeClosesOnlyCurrentToast(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	oldClosed, currentClosed := 0, 0
	f.add(ToastOptions{ID: "other", Title: "Other toast stays", Timeout: -1})
	f.add(ToastOptions{ID: "escape", Title: "Loading keyboard task", Kind: ToastLoading, Action: "Inspect keyboard task", OnClose: func() { oldClosed++ }})
	f.tt.Advance(2 * time.Second)
	notificationFocus(t, f.tt, "Inspect keyboard task")
	f.add(ToastOptions{ID: "escape", Title: "Keyboard task failed", Kind: ToastError, Action: "Inspect keyboard task", Timeout: time.Second, OnClose: func() { currentClosed++ }})
	f.tt.Advance(3 * time.Second)
	if !f.tt.Focused("Inspect keyboard task") || currentClosed != 0 {
		t.Fatal("focused same-ID update lost native pause or focus")
	}
	f.tt.Key(0, myui.KeyEscape)
	f.tt.Advance(2 * time.Second)
	if currentClosed != 1 || oldClosed != 0 || f.tt.HasText("Keyboard task failed") || !f.tt.HasText("Other toast stays") {
		t.Fatalf("Escape scope/update callback: old=%d current=%d texts=%q", oldClosed, currentClosed, f.tt.Texts())
	}
}

func TestToastNotificationCloseCallbackCanClearEntireSource(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	closed := 0
	model := []string{"First source", "Second source"}
	for i, title := range model {
		f.add(ToastOptions{ID: fmt.Sprintf("clear-all-%d", i), Title: title, Timeout: -1, OnClose: func() {
			closed++
			model = nil
			f.pending = append(f.pending, func(c *myui.Context) { DismissToast(c, "") })
		}})
	}
	f.tt.Advance(2 * time.Second)
	f.dismiss("clear-all-1")
	f.tt.Frame()
	f.tt.Advance(2 * time.Second)
	if model != nil || closed != 2 || len(f.state.ids) != 0 || f.tt.HasText("First source") || f.tt.HasText("Second source") {
		t.Fatalf("callback clear-all lifecycle: model=%q close=%d ids=%q texts=%q", model, closed, f.state.ids, f.tt.Texts())
	}
	if f.tt.FramePending() {
		t.Fatal("callback-cleared deck did not release retained animations")
	}
}

func TestToastMotionEntryFirstMiddleAndSettledPaint(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	// The empty native viewport already exists; the native enter policy applies.
	blank := notificationImage(f.tt)
	f.add(ToastOptions{ID: "entry", Title: "Event created", Description: "The local sample has been saved.", Timeout: -1})
	first := notificationImage(f.tt)
	notificationCapture(t, f.tt, "entry-first")
	f.tt.Advance(80 * time.Millisecond)
	mid := notificationImage(f.tt)
	notificationCapture(t, f.tt, "entry-mid")
	f.tt.Advance(2 * time.Second)
	settled := notificationImage(f.tt)
	notificationCapture(t, f.tt, "entry-settled")
	region := image.Rect(180, 400, 640, 700)
	if notificationPixelDifference(blank, settled, region) < 100 {
		t.Fatal("settled toast did not paint its card and content")
	}
	if notificationPixelDifference(first, mid, region) < 40 || notificationPixelDifference(mid, settled, region) < 40 {
		t.Fatal("entry omitted either its first or middle animated paint")
	}
	if f.tt.FramePending() {
		t.Fatal("settled entry keeps requesting animation frames")
	}
}

func TestToastNotificationSelectableBodyDoesNotSwipe(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	closed := 0
	f.add(ToastOptions{ID: "selection", Title: "Selectable description", Description: "Select this readable body text", Timeout: -1, OnClose: func() { closed++ }})
	f.tt.Advance(2 * time.Second)
	notificationExpand(t, f, "selection")
	card := notificationCard(t, f, "selection")
	body := notificationRect(t, f.tt, "Select this readable body text")
	y := body.Y + body.H/2
	f.tt.Press(body.X+2, y)
	f.tt.Move(body.X+min(body.W-2, 150), y)
	f.tt.Release(body.X+min(body.W-2, 150), y)
	f.tt.Key(myui.Cmd, myui.KeyC)
	f.tt.Advance(2 * time.Second)
	if closed != 0 || !notificationSameRect(card, notificationCard(t, f, "selection")) {
		t.Fatal("text selection was treated as a card swipe")
	}
	if f.tt.Clipboard() == "" {
		t.Fatal("description did not expose native selectable text and Copy")
	}
}

func TestToastNotificationCardTracksSemanticTheme(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	f.tt.SetPreferences(myui.Preferences{ReduceMotion: true})
	f.tt.SetDark(false)
	f.add(ToastOptions{ID: "theme", Title: "Semantic surface", Timeout: -1})
	card := notificationCard(t, f, "theme")
	x, y := int(card.X+card.W/2), int(card.Y+3)
	light := f.tt.Image().RGBAAt(x, y)
	f.tt.SetDark(true)
	f.tt.Advance(2 * time.Second)
	dark := f.tt.Image().RGBAAt(x, y)
	if int(light.R)+int(light.G)+int(light.B) <= int(dark.R)+int(dark.G)+int(dark.B)+100 {
		t.Fatalf("notification surface did not track semantic appearance: light=%v dark=%v", light, dark)
	}
	if !notificationSameRect(card, notificationCard(t, f, "theme")) {
		t.Fatal("theme update changed card identity or geometry")
	}
}

func TestToastMotionLoadingExitStopsNativeLoop(t *testing.T) {
	f := newNotificationFixture(t, 640, 700)
	closed := 0
	f.add(ToastOptions{ID: "loop", Title: "Loading exit", Kind: ToastLoading, OnClose: func() { closed++ }})
	f.tt.Advance(2 * time.Second)
	spinner := notificationRect(t, f.tt, "Loading")
	region := image.Rect(int(spinner.X), int(spinner.Y), int(math.Ceil(float64(spinner.X+spinner.W))), int(math.Ceil(float64(spinner.Y+spinner.H))))
	before := notificationImage(f.tt)
	// The native spinner advances one spoke every 75ms on a paint-only
	// deadline. FramePending reports view requests, not those paint deadlines.
	// 113ms always crosses a spoke boundary without sampling a full turn.
	f.tt.Advance(113 * time.Millisecond)
	if notificationPixelDifference(before, notificationImage(f.tt), region) == 0 {
		t.Fatal("normal-motion loading indicator did not repaint its next spoke")
	}
	f.dismiss("loop")
	f.tt.Advance(2 * time.Second)
	if closed != 1 || f.tt.HasText("Loading exit") || f.tt.FramePending() {
		t.Fatalf("loading exit retained loop or callback: close=%d pending=%v texts=%q", closed, f.tt.FramePending(), f.tt.Texts())
	}
	settled := notificationImage(f.tt)
	f.tt.Advance(113 * time.Millisecond)
	if notificationPixelDifference(settled, notificationImage(f.tt), settled.Bounds()) != 0 {
		t.Fatal("loading exit continued painting after its native transition settled")
	}
}

func TestToastNotificationRetainedExitPassesPointerThrough(t *testing.T) {
	var state ToastState
	var pending func(*myui.Context)
	underlay, closed := 0, 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Button(c, "Underlying action").Absolute().Right(16).Bottom(16).Size(384, 160).OnClick(func() { underlay++ })
		if pending != nil {
			event := pending
			pending = nil
			event(c)
		}
		Toaster(c)
	}, 640, 700)
	pending = func(c *myui.Context) {
		Toast(c, &state, ToastOptions{ID: "pass-through", Title: "Exit is inert", Timeout: -1, OnClose: func() { closed++ }})
	}
	tt.Frame()
	tt.Advance(2 * time.Second)
	card := notificationRect(t, tt, "Toast: pass-through")
	pending = func(c *myui.Context) { DismissToast(c, "pass-through") }
	tt.Frame()
	// No clock advancement: the native exit snapshot is still at its first pose.
	tt.ClickAt(card.X+card.W/2, card.Y+card.H/2)
	if underlay != 1 || closed != 1 {
		t.Fatalf("retained exit intercepted pointer: underlay=%d close=%d", underlay, closed)
	}
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Fatal("inert exit did not release animation")
	}
}

func TestToastNotificationModalBlocksPointerAndKeyboardActivation(t *testing.T) {
	var state ToastState
	var pending func(*myui.Context)
	open := false
	actions, closed, approvals := 0, 0, 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		if pending != nil {
			event := pending
			pending = nil
			event(c)
		}
		// Consumer host order: the modal scrim is built after notifications.
		Toaster(c)
		AlertDialog(c, &open, AlertDialogOptions{
			Title: "Approve scoped operation", Description: "The decision must retain exclusive input.",
			CancelLabel: "Deny operation", ActionLabel: "Approve operation",
			OnAction: func() { approvals++ },
		}, nil)
	}, 640, 850)
	pending = func(c *myui.Context) {
		Toast(c, &state, ToastOptions{ID: "behind-modal", Title: "Background notification", Action: "Background toast action", Timeout: -1,
			OnAction: func() { actions++ }, OnClose: func() { closed++ }})
	}
	tt.Frame()
	tt.Advance(2 * time.Second)
	action := notificationRect(t, tt, "Background toast action")
	closeButton := notificationRect(t, tt, "Close toast")
	open = true
	tt.Frame()
	tt.Advance(2 * time.Second)
	tt.ClickAt(action.X+action.W/2, action.Y+action.H/2)
	tt.ClickAt(closeButton.X+closeButton.W/2, closeButton.Y+closeButton.H/2)
	if !open || actions != 0 || closed != 0 || approvals != 0 {
		t.Fatalf("modal did not block toast pointer input: open=%v actions=%d close=%d approvals=%d", open, actions, closed, approvals)
	}
	for range 8 {
		tt.Key(0, myui.KeyTab)
		if tt.Focused("Background toast action") || tt.Focused("Close toast") {
			t.Fatal("modal focus traversal escaped into notifications")
		}
		if !tt.Focused("Deny operation") && !tt.Focused("Approve operation") {
			t.Fatal("approval dialog lost exclusive keyboard focus")
		}
	}
	notificationFocus(t, tt, "Approve operation")
	tt.Key(0, myui.KeyEnter)
	tt.Advance(2 * time.Second)
	if open || approvals != 1 || actions != 0 || closed != 0 {
		t.Fatalf("modal keyboard activation reached wrong target: open=%v approvals=%d toast action=%d close=%d", open, approvals, actions, closed)
	}
	notificationClick(t, tt, "Background toast action")
	tt.Advance(2 * time.Second)
	if actions != 1 || closed != 1 || approvals != 1 {
		t.Fatalf("toast did not regain input after modal closed: actions=%d close=%d approvals=%d", actions, closed, approvals)
	}
}
