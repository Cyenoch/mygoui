package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestResultShowsOutcomeAndActions pins the outcome page: a bare status is
// only its icon, and caller copy and actions appear and work.
func TestResultShowsOutcomeAndActions(t *testing.T) {
	retried, detailed := 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		o := mygoui.ResultOptions{Status: mygoui.ColorDanger}
		if detailed {
			o.Title, o.Description = "Payment failed", "The card was declined."
			o.Actions = func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Retry", OnClick: func() { retried++ }})
			}
			o.Content = func() { ui.Text(c, "Code 402") }
		}
		mygoui.Result(c, o).Label("Outcome")
	}, 480, 400)
	if got := tt.Texts(); len(got) != 1 {
		t.Fatalf("a bare status shows %v", got)
	}
	detailed = true
	tt.Frame()
	for _, text := range []string{"Payment failed", "The card was declined.", "Code 402"} {
		if !tt.HasText(text) {
			t.Fatalf("result shows %v, want %q", tt.Texts(), text)
		}
	}
	title, action, detail := findRect(t, tt, "Payment failed"), findRect(t, tt, "Retry"), findRect(t, tt, "Code 402")
	if title.Y >= action.Y || action.Y >= detail.Y {
		t.Fatalf("order: title=%+v action=%+v detail=%+v", title, action, detail)
	}
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	if retried != 1 {
		t.Fatalf("action fired %d times", retried)
	}
}
