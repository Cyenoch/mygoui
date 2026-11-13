package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func rateTester(value *float64, o RateOptions) *ui.Tester {
	o.AriaLabel = "Rating"
	return ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { Rate(c, value, o) })
	}, 240, 120)
}

// rateClick clicks the given fraction across the rating's width.
func rateClick(t *testing.T, tt *ui.Tester, fraction float32) {
	t.Helper()
	r, ok := tt.Find("Rating")
	if !ok {
		t.Fatal("the rating is not rendered")
	}
	tt.ClickAt(r.X+r.W*fraction, r.Y+r.H/2)
}

func TestRateArrowsStepWithinBounds(t *testing.T) {
	value, changes := 4.0, 0
	tt := rateTester(&value, RateOptions{OnChange: func() { changes++ }})
	rateClick(t, tt, 0.7) // focuses; the fourth star keeps the value
	tt.Key(0, ui.KeyRight)
	if value != 5 || changes != 1 {
		t.Fatalf("Right gave %v with %d changes, want 5 and 1", value, changes)
	}
	tt.Key(0, ui.KeyRight)
	if value != 5 || changes != 1 {
		t.Fatalf("Right past the last star gave %v with %d changes", value, changes)
	}
	tt.Key(0, ui.KeyLeft)
	if value != 4 {
		t.Fatalf("Left gave %v, want 4", value)
	}
	tt.Key(0, ui.KeyHome)
	if value != 0 {
		t.Fatalf("Home gave %v, want 0", value)
	}
	tt.Key(0, ui.KeyEnd)
	if value != 5 {
		t.Fatalf("End gave %v, want 5", value)
	}
}

func TestRateArrowsFollowDirection(t *testing.T) {
	value := 2.0
	tt := rateTester(&value, RateOptions{AllowHalf: true})
	tt.SetDirection(ui.RightToLeft)
	rateClick(t, tt, 0.5)
	start := value
	tt.Key(0, ui.KeyLeft)
	if value != start+0.5 {
		t.Fatalf("Left in a right-to-left layout moved %v to %v, want a half step up", start, value)
	}
	// The first star is at the right edge.
	rateClick(t, tt, 0.95)
	if value != 0.5 {
		t.Fatalf("a click at the right edge of a right-to-left rating gave %v, want 0.5", value)
	}
}

func TestRateClickChoosesStarAndHalf(t *testing.T) {
	value := 0.0
	tt := rateTester(&value, RateOptions{AllowHalf: true, AllowClear: true})
	rateClick(t, tt, 0.95)
	if value != 5 {
		t.Fatalf("a click on the last star's far half gave %v, want 5", value)
	}
	rateClick(t, tt, 0.05)
	if value != 0.5 {
		t.Fatalf("a click on the first star's near half gave %v, want 0.5", value)
	}
	rateClick(t, tt, 0.05)
	if value != 0 {
		t.Fatalf("a click on the current rating gave %v, want it cleared", value)
	}
}

func TestRateReadOnlyIgnoresInput(t *testing.T) {
	value := 2.0
	tt := rateTester(&value, RateOptions{ReadOnly: true})
	rateClick(t, tt, 0.95)
	tt.Key(0, ui.KeyRight)
	if value != 2 {
		t.Fatalf("a read-only rating changed to %v", value)
	}
}

// Enter and Space on the focused rating must not pick the star the pointer
// happens to rest on.
func TestRateKeyboardActivationIgnoresRestingPointer(t *testing.T) {
	value := 4.0
	tt := rateTester(&value, RateOptions{})
	rateClick(t, tt, 0.7) // focuses; the fourth star keeps the value
	r, _ := tt.Find("Rating")
	tt.Move(r.X+r.W*0.1, r.Y+r.H/2)
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyEnter)
	if value != 4 {
		t.Fatalf("Space and Enter with the pointer on the first star gave %v, want 4", value)
	}
}

// PageUp and PageDown move a rating by a whole star, from a half star too.
func TestRatePageKeysMoveAWholeStar(t *testing.T) {
	stars := 1.0
	tt := formTextTester(t, func(c *ui.Context) {
		Rate(c, &stars, RateOptions{AriaLabel: "Stars", AllowHalf: true})
	})
	formTextClick(t, tt, "Stars")
	stars = 1
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyPageUp)
	if stars != 2.5 {
		t.Fatalf("Right then PageUp rated %v, want 2.5", stars)
	}
	tt.Key(0, ui.KeyPageDown)
	if stars != 1.5 {
		t.Fatalf("PageDown rated %v, want 1.5", stars)
	}
}
