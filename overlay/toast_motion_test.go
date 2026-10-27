package overlay

import (
	"fmt"
	"image"
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

// These tests drive the native viewport with keyboard and pointer events.
// Cards are found by their accessible notification title, not their internal
// identity; MyGo owns callbacks, lifetime and exit snapshots.

type toastMotionFixture struct {
	tt      *ui.Tester
	state   ToastState
	pending []func(*ui.Context)
}

func newToastMotionFixture(t *testing.T, width, height int) *toastMotionFixture {
	t.Helper()
	f := &toastMotionFixture{}
	f.tt = ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Page action")
		pending := f.pending
		f.pending = nil
		for _, event := range pending {
			event(c)
		}
		f.state.View(c)
	}, width, height)
	f.tt.SetPreferences(ui.Preferences{TextScale: 1})
	f.tt.Frame()
	return f
}

func (f *toastMotionFixture) add(o ToastOptions) {
	f.pending = append(f.pending, func(c *ui.Context) { Toast(c, &f.state, o) })
	f.tt.Frame()
}

func (f *toastMotionFixture) dismiss(id string) {
	f.pending = append(f.pending, func(c *ui.Context) { f.state.Dismiss(c, id) })
	f.tt.Frame()
}

func toastMotionRect(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("missing toast %q; labels=%q", label, tt.Texts())
	}
	return r
}

func toastMotionNear(a, b float32) bool {
	return math.Abs(float64(a-b)) <= 1
}

func toastMotionSameRect(a, b ui.Rect) bool {
	return toastMotionNear(a.X, b.X) && toastMotionNear(a.Y, b.Y) &&
		toastMotionNear(a.W, b.W) && toastMotionNear(a.H, b.H)
}

func toastMotionFocus(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	for range 16 {
		if tt.Focused(label) {
			return
		}
		tt.Key(0, ui.KeyTab)
	}
	t.Fatalf("Tab did not reach %q; labels=%q", label, tt.Texts())
}

func toastMotionImage(tt *ui.Tester) *image.RGBA {
	frame := tt.Image()
	snapshot := image.NewRGBA(frame.Bounds())
	for y := frame.Bounds().Min.Y; y < frame.Bounds().Max.Y; y++ {
		copy(snapshot.Pix[snapshot.PixOffset(0, y):], frame.Pix[frame.PixOffset(0, y):frame.PixOffset(0, y)+frame.Bounds().Dx()*4])
	}
	return snapshot
}

func toastMotionPixelDifference(a, b *image.RGBA, r image.Rectangle) int {
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

func toastMotionRegion(card ui.Rect) image.Rectangle {
	return image.Rect(int(card.X)-2, int(card.Y)-2, int(math.Ceil(float64(card.X+card.W)))+2, int(math.Ceil(float64(card.Y+card.H)))+2)
}

func TestToastEntryTravelsAndFadesFromWindowEdge(t *testing.T) {
	for _, tc := range []struct {
		name      string
		place     ToastPlacement
		direction float32
	}{{"bottom", ToastBottom, 1}, {"top", ToastTop, -1}} {
		t.Run(tc.name, func(t *testing.T) {
			f := newToastMotionFixture(t, 640, 700)
			blank := toastMotionImage(f.tt)
			f.add(ToastOptions{ID: "arrival", Title: "New notification", Placement: tc.place, Timeout: -1})
			first := toastMotionRect(t, f.tt, "New notification")
			firstImage := toastMotionImage(f.tt)
			// A card now enters from a whole card height away, off the window,
			// so it is compared once it is painted inside it.
			f.tt.Advance(180 * time.Millisecond)
			f.tt.Frame()
			mid := toastMotionRect(t, f.tt, "New notification")
			midImage := toastMotionImage(f.tt)
			f.tt.Advance(2 * time.Second)
			end := toastMotionRect(t, f.tt, "New notification")
			if (first.Y-mid.Y)*tc.direction <= 0 || (mid.Y-end.Y)*tc.direction <= 0 {
				t.Fatalf("entry does not travel from the window edge: first=%+v mid=%+v end=%+v", first, mid, end)
			}
			region := toastMotionRegion(end)
			if toastMotionPixelDifference(blank, firstImage, region) != 0 ||
				toastMotionPixelDifference(firstImage, midImage, region) < 40 {
				t.Fatal("entry does not fade from transparent through an intermediate painted frame")
			}
			if f.tt.FramePending() {
				t.Fatal("arrival continues requesting frames after settling")
			}
		})
	}
}

// TestToastDeckPlacementAndIdleStack pins the deck's resting geometry: the
// newest card hugs the chosen window edge and the older two sit behind it at
// 12 DIP steps and 0.04 narrower each. The bottom cards are compared by their
// bottom edges and the top cards by their top edges, so the idle shrink cannot
// make the assertion depend on where a card is anchored.
func TestToastDeckPlacementAndIdleStack(t *testing.T) {
	const winW, winH = 640, 800
	for _, tc := range []struct {
		name  string
		place ToastPlacement
	}{{"bottom", ToastBottom}, {"top", ToastTop}} {
		t.Run(tc.name, func(t *testing.T) {
			f := newToastMotionFixture(t, winW, winH)
			for _, o := range []ToastOptions{
				{ID: "stack-old", Title: "Oldest", Description: "A longer older message that wraps onto another line.", Placement: tc.place, Timeout: -1},
				{ID: "stack-mid", Title: "Middle", Placement: tc.place, Timeout: -1},
				{ID: "stack-new", Title: "Newest", Placement: tc.place, Timeout: -1},
			} {
				f.add(o)
			}
			f.tt.Advance(2 * time.Second)
			f.tt.Frame()

			newest := toastMotionRect(t, f.tt, "Newest")
			middle := toastMotionRect(t, f.tt, "Middle")
			oldest := toastMotionRect(t, f.tt, "Oldest")
			if !toastMotionNear(newest.X+newest.W/2, winW/2) ||
				!toastMotionNear(middle.X+middle.W/2, winW/2) ||
				!toastMotionNear(oldest.X+oldest.W/2, winW/2) {
				t.Fatalf("deck is not centered: newest=%+v middle=%+v oldest=%+v", newest, middle, oldest)
			}
			if !toastMotionNear(newest.W, 360) {
				t.Fatalf("newest width=%v, want 360", newest.W)
			}
			for i, card := range []ui.Rect{middle, oldest} {
				want := float32(1 - float64(i+1)*0.04)
				if got := card.W / newest.W; math.Abs(float64(got-want)) > 0.01 {
					t.Fatalf("idle depth %d width ratio=%v, want %v", i+1, got, want)
				}
			}

			if tc.place == ToastTop {
				if !toastMotionNear(newest.Y, 16) {
					t.Fatalf("top newest Y=%v, want 16", newest.Y)
				}
				if !toastMotionNear(middle.Y-newest.Y, 12) ||
					!toastMotionNear(oldest.Y-newest.Y, 24) {
					t.Fatalf("top idle stack: newest=%+v middle=%+v oldest=%+v", newest, middle, oldest)
				}
			} else {
				if !toastMotionNear(newest.Y+newest.H, winH-16) {
					t.Fatalf("bottom newest bottom=%v, want %v", newest.Y+newest.H, winH-16)
				}
				if !toastMotionNear(newest.Y+newest.H-(middle.Y+middle.H), 12) ||
					!toastMotionNear(newest.Y+newest.H-(oldest.Y+oldest.H), 24) {
					t.Fatalf("bottom idle stack: newest=%+v middle=%+v oldest=%+v", newest, middle, oldest)
				}
			}
			if f.tt.FramePending() {
				t.Fatal("idle deck keeps requesting animation frames")
			}
		})
	}
}

// TestToastDeckInsertionAndRemovalKeepTheOldCardMoving checks the first frame
// after a new card joins or leaves: the older card keeps its preceding pose,
// then travels through intermediate positions to the new resting slot. A
// height change must animate the declared geometry, not snap it after layout.
func TestToastDeckInsertionAndRemovalKeepTheOldCardMoving(t *testing.T) {
	const winW, winH = 640, 820
	f := newToastMotionFixture(t, winW, winH)
	f.add(ToastOptions{ID: "insert-old", Title: "Existing card", Timeout: -1})
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	before := toastMotionRect(t, f.tt, "Existing card")

	f.add(ToastOptions{ID: "insert-new", Title: "Unequal height card",
		Description: "A much longer message wraps onto several lines so the deck cannot place its cards from a fixed nominal height. The complete description remains readable while the older card slides to its new slot.",
		Timeout:     -1})
	first := toastMotionRect(t, f.tt, "Existing card")
	if !toastMotionSameRect(before, first) {
		t.Fatalf("insertion jumped the existing card on its first frame: before=%+v first=%+v", before, first)
	}
	f.tt.Advance(80 * time.Millisecond)
	f.tt.Frame()
	mid := toastMotionRect(t, f.tt, "Existing card")
	if toastMotionSameRect(mid, first) || !(mid.Y < first.Y) {
		t.Fatalf("existing card did not move continuously up on insertion: first=%+v mid=%+v", first, mid)
	}
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	end := toastMotionRect(t, f.tt, "Existing card")
	tall := toastMotionRect(t, f.tt, "Unequal height card")
	if !(end.Y < mid.Y) || tall.H <= before.H+10 {
		t.Fatalf("insertion did not settle with unequal heights: before=%+v end=%+v mid=%+v tall=%+v", before, end, mid, tall)
	}

	f.dismiss("insert-new")
	removeFirst := toastMotionRect(t, f.tt, "Existing card")
	if !toastMotionSameRect(end, removeFirst) {
		t.Fatalf("removal jumped the existing card on its first frame: before=%+v first=%+v", end, removeFirst)
	}
	f.tt.Advance(80 * time.Millisecond)
	f.tt.Frame()
	removeMid := toastMotionRect(t, f.tt, "Existing card")
	if toastMotionSameRect(removeMid, removeFirst) || !(removeMid.Y > removeFirst.Y) {
		t.Fatalf("existing card did not move continuously back on removal: first=%+v mid=%+v", removeFirst, removeMid)
	}
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if got := toastMotionRect(t, f.tt, "Existing card"); !toastMotionSameRect(before, got) {
		t.Fatalf("existing card did not return to its resting slot: before=%+v after=%+v", before, got)
	}
	if f.tt.HasText("Unequal height card") || f.tt.FramePending() {
		t.Fatalf("removal retained a live subtree or animation: pending=%v labels=%q", f.tt.FramePending(), f.tt.Texts())
	}
}

// TestToastDeckHoverExpandsUnequalCardsAndHoldsGapHover covers the expanded
// deck: each card takes its measured height with a real 12 DIP gap, and a
// stationary pointer in a gap keeps the expansion and the native expiry pause
// alive instead of collapsing and reopening under the pointer.
func TestToastDeckHoverExpandsUnequalCardsAndHoldsGapHover(t *testing.T) {
	const winW, winH = 640, 820
	f := newToastMotionFixture(t, winW, winH)
	closed := 0
	for _, o := range []ToastOptions{
		{ID: "hover-old", Title: "Hover oldest",
			Description: "A much longer older message wraps onto several lines so the expanded deck must use its measured height rather than the newest card's height. It stays readable while hovered.",
			Timeout:     10 * time.Second, OnClose: func() { closed++ }},
		{ID: "hover-mid", Title: "Hover middle", Timeout: 10 * time.Second, OnClose: func() { closed++ }},
		{ID: "hover-new", Title: "Hover newest", Timeout: 10 * time.Second, OnClose: func() { closed++ }},
	} {
		f.add(o)
	}
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()

	newest := toastMotionRect(t, f.tt, "Hover newest")
	f.tt.Move(newest.X+newest.W/2, newest.Y+newest.H/2)
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	newest = toastMotionRect(t, f.tt, "Hover newest")
	middle := toastMotionRect(t, f.tt, "Hover middle")
	oldest := toastMotionRect(t, f.tt, "Hover oldest")
	for _, card := range []ui.Rect{newest, middle, oldest} {
		if !toastMotionNear(card.W, 360) {
			t.Fatalf("expanded card width=%v, want 360: newest=%+v middle=%+v oldest=%+v", card.W, newest, middle, oldest)
		}
	}
	if oldest.H <= middle.H+10 {
		t.Fatalf("expanded cards do not use unequal measured heights: oldest=%+v middle=%+v", oldest, middle)
	}
	gaps := []float32{
		newest.Y - (middle.Y + middle.H),
		middle.Y - (oldest.Y + oldest.H),
	}
	for _, gap := range gaps {
		if !toastMotionNear(gap, 12) {
			t.Fatalf("expanded gaps=%v, want 12 DIP: newest=%+v middle=%+v oldest=%+v", gaps, newest, middle, oldest)
		}
	}

	// Rest the pointer in the gap between the two older cards. The pause and
	// the expansion belong to the whole deck, not only to a card's face.
	gapX := middle.X + middle.W/2
	gapY := middle.Y - 6
	f.tt.Move(gapX, gapY)
	f.tt.Advance(20 * time.Second)
	for range 5 {
		f.tt.Frame()
		if !toastMotionSameRect(middle, toastMotionRect(t, f.tt, "Hover middle")) ||
			!toastMotionSameRect(oldest, toastMotionRect(t, f.tt, "Hover oldest")) {
			t.Fatal("stationary pointer in the inter-card gap oscillated the expanded deck")
		}
	}
	if closed != 0 {
		t.Fatalf("hovering the gap lost the native expiry pause: closed=%d", closed)
	}

	f.tt.Move(0, 0)
	f.tt.Advance(12 * time.Second)
	if closed != 3 || f.tt.HasText("Hover oldest") || f.tt.HasText("Hover middle") || f.tt.HasText("Hover newest") {
		t.Fatalf("leaving the gap did not resume expiry: closed=%d labels=%q", closed, f.tt.Texts())
	}
	if f.tt.FramePending() {
		t.Fatal("empty deck still animates after the last exit")
	}
}

// TestToastDeckRapidHoverReversalStartsAtCurrentPose reverses the deck while
// its expansion is still in flight and then enters it again. Each change of
// direction must begin at the pose already painted, without a first-frame jump
// back to either endpoint, and the final collapse must land at the original
// resting stack.
func TestToastDeckRapidHoverReversalStartsAtCurrentPose(t *testing.T) {
	const winW, winH = 640, 820
	f := newToastMotionFixture(t, winW, winH)
	for i := range 3 {
		f.add(ToastOptions{ID: fmt.Sprintf("reverse-%d", i), Title: fmt.Sprintf("Reversal card %d", i), Timeout: -1})
	}
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	collapsed := toastMotionRect(t, f.tt, "Reversal card 0")
	newest := toastMotionRect(t, f.tt, "Reversal card 2")

	f.tt.Move(newest.X+newest.W/2, newest.Y+newest.H/2)
	f.tt.Advance(80 * time.Millisecond)
	f.tt.Frame()
	part := toastMotionRect(t, f.tt, "Reversal card 0")
	if toastMotionSameRect(collapsed, part) {
		t.Fatalf("hover expansion never entered an intermediate pose: collapsed=%+v part=%+v", collapsed, part)
	}
	f.tt.Move(0, 0)
	if got := toastMotionRect(t, f.tt, "Reversal card 0"); !toastMotionSameRect(part, got) {
		t.Fatalf("leaving hover jumped instead of starting at the current pose: part=%+v got=%+v", part, got)
	}
	f.tt.Advance(40 * time.Millisecond)
	f.tt.Frame()
	part = toastMotionRect(t, f.tt, "Reversal card 0")
	newest = toastMotionRect(t, f.tt, "Reversal card 2")
	f.tt.Move(newest.X+newest.W/2, newest.Y+newest.H/2)
	if got := toastMotionRect(t, f.tt, "Reversal card 0"); !toastMotionSameRect(part, got) {
		t.Fatalf("re-entering hover jumped instead of starting at the current pose: part=%+v got=%+v", part, got)
	}
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	f.tt.Move(0, 0)
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if got := toastMotionRect(t, f.tt, "Reversal card 0"); !toastMotionSameRect(collapsed, got) {
		t.Fatalf("rapid reversal did not settle back to the collapsed deck: collapsed=%+v got=%+v", collapsed, got)
	}
	if f.tt.FramePending() {
		t.Fatal("rapid reversal keeps requesting animation frames after settling")
	}
}

// TestToastDeckFocusExpansionAndSameIDUpdate expands the deck from the keyboard,
// then updates the focused toast: the same native button and card identity keep
// focus, no entrance is replayed, and the focus keeps the updated timeout paused.
func TestToastDeckFocusExpansionAndSameIDUpdate(t *testing.T) {
	const winW, winH = 640, 820
	f := newToastMotionFixture(t, winW, winH)
	actions, closed := 0, 0
	for _, o := range []ToastOptions{
		{ID: "focus-old", Title: "Focused oldest",
			Description: "A longer focused message that wraps so the expanded deck must use its measured height while the keyboard stays in the action.",
			Action:      "Inspect oldest", OnAction: func() {}, Timeout: -1},
		{ID: "focus-mid", Title: "Focused middle", Action: "Inspect middle", OnAction: func() {}, Timeout: -1},
		{ID: "focus-new", Title: "Focused newest", Action: "Inspect newest", OnAction: func() {}, Timeout: -1},
	} {
		f.add(o)
	}
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()

	toastMotionFocus(t, f.tt, "Inspect oldest")
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	oldest := toastMotionRect(t, f.tt, "Focused oldest")
	middle := toastMotionRect(t, f.tt, "Focused middle")
	if !toastMotionNear(oldest.W, 360) || oldest.H <= middle.H+10 {
		t.Fatalf("focus did not expand unequal measured cards: oldest=%+v middle=%+v", oldest, middle)
	}

	f.add(ToastOptions{ID: "focus-old", Title: "Updated oldest",
		Description: "The focused card keeps its identity and pause while the new message replaces the old one.",
		Action:      "Inspect oldest", Timeout: time.Second,
		OnAction: func() { actions++ }, OnClose: func() { closed++ }})
	if got := toastMotionRect(t, f.tt, "Updated oldest"); !toastMotionSameRect(oldest, got) {
		t.Fatalf("same-ID update replayed the entrance: before=%+v after=%+v", oldest, got)
	}
	if !f.tt.Focused("Inspect oldest") || f.tt.HasText("Focused oldest") || !f.tt.HasText("Updated oldest") {
		t.Fatalf("same-ID update lost focus or retained old content: focused=%v labels=%q", f.tt.Focused("Inspect oldest"), f.tt.Texts())
	}
	f.tt.Advance(3 * time.Second)
	if closed != 0 || !f.tt.HasText("Updated oldest") {
		t.Fatalf("focused same-ID update did not keep expiry paused: closed=%d labels=%q", closed, f.tt.Texts())
	}
	f.tt.Key(0, ui.KeyEnter)
	f.tt.Advance(2 * time.Second)
	if actions != 1 || closed != 1 || f.tt.HasText("Updated oldest") {
		t.Fatalf("focused action after update: actions=%d closed=%d labels=%q", actions, closed, f.tt.Texts())
	}
}

// TestToastSameIDUpdateKeepsActiveDragPose checks that an update changes the
// payload without replacing the card or its gesture: the card stays where the
// pointer is holding it and does not paint an entrance pose under the drag.
func TestToastSameIDUpdateKeepsActiveDragPose(t *testing.T) {
	const winW, winH = 640, 700
	f := newToastMotionFixture(t, winW, winH)
	closed := 0
	f.add(ToastOptions{ID: "drag-update", Title: "Dragging card", Action: "Inspect drag", OnAction: func() {}, Timeout: -1,
		OnClose: func() { closed++ }})
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	rest := toastMotionRect(t, f.tt, "Dragging card")

	x, y := rest.X+rest.W/2, rest.Y+3
	f.tt.Press(x, y)
	f.tt.Advance(100 * time.Millisecond)
	f.tt.Move(x+24, y)
	f.tt.Advance(200 * time.Millisecond)
	dragged := toastMotionRect(t, f.tt, "Dragging card")
	if dragged.X <= rest.X+5 {
		t.Fatalf("background drag did not follow the pointer: rest=%+v dragged=%+v", rest, dragged)
	}

	f.add(ToastOptions{ID: "drag-update", Title: "Dragging updated", Action: "Inspect drag", OnAction: func() {}, Timeout: -1,
		OnClose: func() { closed++ }})
	if got := toastMotionRect(t, f.tt, "Dragging updated"); !toastMotionSameRect(dragged, got) {
		t.Fatalf("same-ID update reset the active drag pose: before=%+v after=%+v", dragged, got)
	}
	f.tt.Release(x+24, y)
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if got := toastMotionRect(t, f.tt, "Dragging updated"); !toastMotionSameRect(rest, got) || closed != 0 {
		t.Fatalf("cancelled drag did not return with the new payload: rest=%+v got=%+v closed=%d", rest, got, closed)
	}
}

// TestToastDragCancelsCommitsAndControlsOwnTheGesture checks the swipe
// boundary on the card background and the nested controls: a short drag
// cancels, a long one commits, and dragging the action or close button neither
// moves the deck nor runs the button's activation.
func TestToastDragCancelsCommitsAndControlsOwnTheGesture(t *testing.T) {
	const winW, winH = 640, 700
	f := newToastMotionFixture(t, winW, winH)
	closed := 0
	f.add(ToastOptions{ID: "swipe", Title: "Swipe card", Action: "Inspect swipe", OnAction: func() {}, Timeout: -1,
		OnClose: func() { closed++ }})
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	rest := toastMotionRect(t, f.tt, "Swipe card")

	x, y := rest.X+rest.W/2, rest.Y+3
	f.tt.Press(x, y)
	f.tt.Advance(100 * time.Millisecond)
	f.tt.Move(x+24, y)
	f.tt.Advance(200 * time.Millisecond)
	if dragged := toastMotionRect(t, f.tt, "Swipe card"); dragged.X <= rest.X+5 {
		t.Fatalf("short drag did not follow the pointer: rest=%+v dragged=%+v", rest, dragged)
	}
	f.tt.Release(x+24, y)
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if got := toastMotionRect(t, f.tt, "Swipe card"); !toastMotionSameRect(rest, got) || closed != 0 {
		t.Fatalf("short drag did not cancel: rest=%+v got=%+v closed=%d", rest, got, closed)
	}

	f.tt.Press(x, y)
	f.tt.Move(x+rest.W*0.6, y)
	f.tt.Release(x+rest.W*0.6, y)
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if f.tt.HasText("Swipe card") || closed != 1 {
		t.Fatalf("long drag did not commit dismissal: closed=%d labels=%q", closed, f.tt.Texts())
	}

	// Fresh fixture: action and close own their own horizontal gestures.
	f = newToastMotionFixture(t, winW, winH)
	actions, closed := 0, 0
	f.add(ToastOptions{ID: "controls", Title: "Controls own gestures", Action: "Undo sample", Timeout: -1,
		OnAction: func() { actions++ }, OnClose: func() { closed++ }})
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	card := toastMotionRect(t, f.tt, "Controls own gestures")
	action := toastMotionRect(t, f.tt, "Undo sample")
	f.tt.Press(action.X+action.W/2, action.Y+action.H/2)
	f.tt.Move(action.X+action.W/2-180, action.Y+action.H/2)
	f.tt.Release(action.X+action.W/2-180, action.Y+action.H/2)
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if actions != 0 || closed != 0 || !toastMotionSameRect(card, toastMotionRect(t, f.tt, "Controls own gestures")) {
		t.Fatalf("action drag activated or moved the card: actions=%d closed=%d labels=%q", actions, closed, f.tt.Texts())
	}
	closeButton := toastMotionRect(t, f.tt, "Close")
	f.tt.Press(closeButton.X+closeButton.W/2, closeButton.Y+closeButton.H/2)
	f.tt.Move(closeButton.X+closeButton.W/2-180, closeButton.Y+closeButton.H/2)
	f.tt.Release(closeButton.X+closeButton.W/2-180, closeButton.Y+closeButton.H/2)
	f.tt.Advance(2 * time.Second)
	f.tt.Frame()
	if actions != 0 || closed != 0 || !toastMotionSameRect(card, toastMotionRect(t, f.tt, "Controls own gestures")) {
		t.Fatalf("close drag activated or moved the card: actions=%d closed=%d labels=%q", actions, closed, f.tt.Texts())
	}
	if err := f.tt.Click("Undo sample"); err != nil {
		t.Fatal(err)
	}
	f.tt.Advance(2 * time.Second)
	if actions != 1 || closed != 1 {
		t.Fatalf("ordinary action after control drags: actions=%d closed=%d", actions, closed)
	}
}

// TestToastReducedMotionAndRetainedExit checks the two lifetime extremes:
// reduced motion puts a card straight at its resting edge with no frame left
// pending, while normal motion keeps a painted exit snapshot that no longer
// takes input or exposes the closed message to assistive technology.
func TestToastReducedMotionAndRetainedExit(t *testing.T) {
	t.Run("reduced motion is immediate", func(t *testing.T) {
		const winW, winH = 640, 700
		f := newToastMotionFixture(t, winW, winH)
		f.tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		closed := 0
		f.add(ToastOptions{ID: "reduced", Title: "Reduced card", Timeout: -1, OnClose: func() { closed++ }})
		card := toastMotionRect(t, f.tt, "Reduced card")
		if !toastMotionNear(card.Y+card.H, winH-16) || !toastMotionNear(card.W, 360) {
			t.Fatalf("reduced-motion card not at its resting edge: %+v", card)
		}
		if f.tt.FramePending() {
			t.Fatal("reduced motion still animates the toast entry")
		}
		f.dismiss("reduced")
		if f.tt.HasText("Reduced card") || f.tt.FramePending() || closed != 1 {
			t.Fatalf("reduced-motion close left motion or content: closed=%d pending=%v labels=%q", closed, f.tt.FramePending(), f.tt.Texts())
		}
	})

	t.Run("retained exit is painted but inert", func(t *testing.T) {
		const winW, winH = 640, 820
		var state ToastState
		var pending func(*ui.Context)
		closed, actions, underlay := 0, 0, 0
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Button(c, "Underlay").Absolute().Right(0).Bottom(0).Size(460, 220).OnClick(func() { underlay++ })
			if pending != nil {
				event := pending
				pending = nil
				event(c)
			}
			state.View(c)
		}, winW, winH)
		tt.SetPreferences(ui.Preferences{TextScale: 1})
		tt.Frame()
		empty := toastMotionImage(tt)

		pending = func(c *ui.Context) {
			Toast(c, &state, ToastOptions{ID: "exit", Title: "Exit card", Description: "The retained copy paints this message after the source is gone.",
				Action: "Exit action", Timeout: -1, OnAction: func() { actions++ }, OnClose: func() { closed++ }})
		}
		tt.Frame()
		tt.Advance(2 * time.Second)
		tt.Frame()
		card := toastMotionRect(t, tt, "Exit card")
		action := toastMotionRect(t, tt, "Exit action")
		region := toastMotionRegion(card)
		shown := toastMotionImage(tt)
		if toastMotionPixelDifference(empty, shown, region) < 500 {
			t.Fatal("settled toast did not paint its card over the page")
		}

		pending = func(c *ui.Context) { state.Dismiss(c, "exit") }
		tt.Frame()
		if closed != 1 || tt.HasText("Exit card") {
			t.Fatalf("close did not release the source: closed=%d labels=%q", closed, tt.Texts())
		}
		first := toastMotionImage(tt)
		if toastMotionPixelDifference(empty, first, region) < 500 {
			t.Fatal("first exit frame did not retain the painted card")
		}
		tt.ClickAt(action.X+action.W/2, action.Y+action.H/2)
		if actions != 0 || underlay != 1 {
			t.Fatalf("retained exit intercepted the pointer: actions=%d underlay=%d", actions, underlay)
		}
		tt.Move(0, 0)
		tt.Advance(80 * time.Millisecond)
		tt.Frame()
		mid := toastMotionImage(tt)
		if toastMotionPixelDifference(first, mid, region) < 40 {
			t.Fatal("exit snapshot painted the same first and middle frames, so it did not animate")
		}
		tt.Advance(2 * time.Second)
		tt.Frame()
		settled := toastMotionImage(tt)
		// Compare a strip of the card over the page, clear of the underlay
		// and any focus ring a pass-through click may have given it.
		sample := image.Rect(int(card.X)+30, int(card.Y)+15, int(card.X)+80, int(math.Ceil(float64(card.Y+card.H)))-15)
		if toastMotionPixelDifference(empty, settled, sample) > 20 {
			t.Fatalf("exit snapshot still paints after the native transition settled: %d pixels", toastMotionPixelDifference(empty, settled, sample))
		}
		if tt.FramePending() {
			t.Fatal("inert exit keeps requesting animation frames")
		}
	})
}
