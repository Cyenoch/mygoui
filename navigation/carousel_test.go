package navigation

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// carouselFixture owns a carousel's index and records the slides built in
// the last frame.
type carouselFixture struct {
	index   int
	count   int
	perView int
	loop    bool
	built   []int
	changes []int
}

func (f *carouselFixture) view(c *ui.Context) {
	theme.Apply(c)
	f.built = f.built[:0]
	Carousel(c, CarouselOptions{
		Label: "Photos", Index: &f.index, Count: f.count, PerView: f.perView, Loop: f.loop,
		Slide: func(i int) {
			f.built = append(f.built, i)
			ui.Text(c, fmt.Sprintf("Slide %d", i+1))
		},
		OnChange: func(i int) { f.changes = append(f.changes, i) },
	})
}

func TestCarouselStepsWithButtonsAndArrowKeys(t *testing.T) {
	f := &carouselFixture{count: 3}
	tt := ui.NewTester(f.view, 480, 200)
	if previous, ok := tt.AccessByRole(ui.RoleButton, "Previous"); !ok || !previous.Disabled {
		t.Fatalf("Previous on the first slide: found %v, disabled %v", ok, previous.Disabled)
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if next, _ := tt.AccessByRole(ui.RoleButton, "Next"); f.index != 2 || !next.Disabled {
		t.Fatalf("index %d, Next disabled %v", f.index, next.Disabled)
	}
	tt.Key(0, ui.KeyRight) // no loop: stays on the last slide
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyLeft)
	if want := []int{1, 2, 1, 0}; !slices.Equal(f.changes, want) || f.index != 0 {
		t.Fatalf("changes %v, want %v; index %d", f.changes, want, f.index)
	}

	f.loop = true
	if err := tt.Click("Previous"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if want := []int{1, 2, 1, 0, 2, 0}; !slices.Equal(f.changes, want) {
		t.Fatalf("looping changes %v, want %v", f.changes, want)
	}
}

func TestCarouselBuildsOnlyTheSlidesInView(t *testing.T) {
	f := &carouselFixture{count: 100, perView: 3, index: 40}
	tt := ui.NewTester(f.view, 640, 200)
	if !slices.Equal(f.built, []int{40, 41, 42}) {
		t.Fatalf("built slides %v", f.built)
	}
	// 98 positions, of which the pips show a window around the current one.
	pips := 0
	for _, a := range tt.AccessAll() {
		if a.Current == ui.CurrentPage {
			pips++
			if a.Name != "Page 41" || a.PosInSet != 41 || a.SetSize != 98 {
				t.Errorf("current pip %q is %d of %d", a.Name, a.PosInSet, a.SetSize)
			}
		} else if a.SetSize == 98 {
			pips++
		}
	}
	if pips != 7 {
		t.Fatalf("%d pips built, want 7", pips)
	}
	if err := tt.Click("Page 43"); err != nil {
		t.Fatal(err)
	}
	if f.index != 42 || !slices.Equal(f.changes, []int{42}) {
		t.Fatalf("a pip chose index %d, changes %v", f.index, f.changes)
	}
	// In the pips the arrows move the focus; they do not step the slides.
	tt.Key(0, ui.KeyRight)
	if f.index != 42 || !tt.Focused("Page 44") {
		t.Fatalf("Right on a pip: index %d, Page 44 focused %v", f.index, tt.Focused("Page 44"))
	}
}

func TestPipsPagerOwnsItsIndexWithoutACaller(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		PipsPager(c, PipsPagerOptions{Label: "Onboarding", Count: 4})
	}, 240, 80)
	if _, ok := tt.AccessByRole(ui.RoleGroup, "Onboarding"); !ok {
		t.Fatal("no group named Onboarding")
	}
	if err := tt.Click("Page 3"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if page, _ := tt.Access("Page 4"); page.Current != ui.CurrentPage || !tt.Focused("Page 4") {
		t.Fatalf("Page 4: current %v, focused %v", page.Current, tt.Focused("Page 4"))
	}
}

func TestCarouselDragTurnsOnePageOnlyPastTheThreshold(t *testing.T) {
	f := &carouselFixture{count: 3}
	tt := ui.NewTester(f.view, 480, 200)
	view, ok := tt.Find("Photos")
	if !ok {
		t.Fatal("no slides")
	}
	x, y := center(view)
	drag := func(by float32) {
		tt.Press(x, y)
		for i := 1; i <= 4; i++ {
			tt.Move(x+by*float32(i)/4, y)
		}
		tt.Release(x+by, y)
		tt.Advance(time.Second)
	}
	drag(-view.W / 8)
	if f.index != 0 || len(f.changes) != 0 {
		t.Fatalf("a short drag turned the page: index %d, changes %v", f.index, f.changes)
	}
	drag(-view.W / 2)
	if f.index != 1 || !slices.Equal(f.changes, []int{1}) {
		t.Fatalf("index %d, changes %v after dragging toward the next slide", f.index, f.changes)
	}
	drag(view.W / 2)
	if f.index != 0 || !slices.Equal(f.changes, []int{1, 0}) {
		t.Fatalf("index %d, changes %v after dragging back", f.index, f.changes)
	}
	drag(view.W / 2) // pulled past the first slide: nowhere to go
	if f.index != 0 || len(f.changes) != 2 || tt.FramePending() {
		t.Fatalf("index %d, changes %v at the start, pending %v", f.index, f.changes, tt.FramePending())
	}
}

func TestCarouselSidewaysSwipeTurnsOnePage(t *testing.T) {
	f := &carouselFixture{count: 4}
	tt := ui.NewTester(f.view, 480, 200)
	view, _ := tt.Find("Photos")
	x, y := center(view)
	for range 6 { // one swipe, with its momentum
		tt.Scroll(x, y, 30, 2)
		tt.Advance(16 * time.Millisecond)
	}
	if f.index != 1 || len(f.changes) != 1 {
		t.Fatalf("one swipe went to %d with changes %v", f.index, f.changes)
	}
	tt.Advance(time.Second)
	tt.Scroll(x, y, -60, 0)
	tt.Advance(time.Second)
	tt.Scroll(x, y, 0, 80) // vertical: the page's
	if f.index != 0 || len(f.changes) != 2 {
		t.Fatalf("index %d, changes %v", f.index, f.changes)
	}
}

func TestCarouselArrowsAreCenteredOnTheSlides(t *testing.T) {
	f := &carouselFixture{count: 3}
	tt := ui.NewTester(f.view, 480, 200)
	view, _ := tt.Find("Photos")
	_, mid := center(view)
	for _, name := range []string{"Previous", "Next"} {
		b, _ := tt.Find(name)
		if _, y := center(b); y < mid-1 || y > mid+1 {
			t.Errorf("%s is centered at y=%v, the slides at %v", name, y, mid)
		}
	}
}

func TestCarouselJumpsWithHomeAndEndAndSaysItsPlaceAmongManyPages(t *testing.T) {
	f := &carouselFixture{count: 12}
	tt := ui.NewTester(f.view, 480, 200)
	if !tt.HasText("1 / 12") {
		t.Fatal("twelve pages show no position beside the pips")
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnd)
	if f.index != 11 || !tt.HasText("12 / 12") {
		t.Fatalf("End: index %d", f.index)
	}
	tt.Key(0, ui.KeyHome)
	if f.index != 0 {
		t.Fatalf("Home: index %d", f.index)
	}
	few := &carouselFixture{count: 3}
	if few2 := ui.NewTester(few.view, 480, 200); few2.HasText("1 / 3") {
		t.Fatal("pips that hold every page also show a count")
	}
}

// TestCarouselKeepsTheFocusAndAnnouncesThePageAStepReaches protects the
// carousel for keyboard and screen reader users: Next disabling at the last
// slide hands the focus to the slides, whose arrows step on, and each step
// announces the page it reached.
func TestCarouselKeepsTheFocusAndAnnouncesThePageAStepReaches(t *testing.T) {
	f := &carouselFixture{count: 3}
	tt := ui.NewTester(f.view, 480, 200)
	tt.UseScreenReader()
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if a, _ := tt.AccessByRole(ui.RoleGroup, "Photos"); f.index != 2 || !a.Focused {
		t.Fatalf("index %d, slides focused %v", f.index, a.Focused)
	}
	tt.Key(0, ui.KeyLeft)
	if want := []string{"Page 2", "Page 3", "Page 2"}; !slices.Equal(tt.Announcements(), want) || f.index != 1 {
		t.Fatalf("announced %q, want %q; index %d", tt.Announcements(), want, f.index)
	}
}

// TestCarouselTellsItsRolesInTheLocale pins the role descriptions: the
// viewport reads as a carousel and each slide in view as a slide, in the
// locale's words.
func TestCarouselTellsItsRolesInTheLocale(t *testing.T) {
	f := &carouselFixture{count: 3, perView: 2}
	tt := ui.NewTester(f.view, 480, 200)
	if a, _ := tt.AccessByRole(ui.RoleGroup, "Photos"); a.RoleDescription != "carousel" {
		t.Fatalf("the viewport's role description is %q, want carousel", a.RoleDescription)
	}
	slides := 0
	for _, a := range tt.AccessAll() {
		if a.Role == ui.RoleGroup && a.SetSize == 3 {
			slides++
			if a.RoleDescription != "slide" {
				t.Fatalf("slide %d's role description is %q, want slide", a.PosInSet, a.RoleDescription)
			}
		}
	}
	if slides != 2 {
		t.Fatalf("found %d slides in view, want 2", slides)
	}
}
