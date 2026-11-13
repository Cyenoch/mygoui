package gallery

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestTokenCardsShowThePublishedThemeValue keeps the Tokens view honest for
// motion, spacing and effect tokens: a card shows the value components read,
// so it follows a theme the application installs.
func TestTokenCardsShowThePublishedThemeValue(t *testing.T) {
	t.Cleanup(func() {
		if err := theme.Use(theme.Default()); err != nil {
			t.Error(err)
		}
	})
	custom := theme.Default()
	custom.Motion.Dialog = 320 * time.Millisecond
	if err := theme.Use(custom); err != nil {
		t.Fatal(err)
	}

	state := New()
	state.Selected = "dialog"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1100, 900)
	tt.Frame()
	r := revealContent(t, tt, state, "token card Motion.Dialog", 1100, 900)
	if r.W <= 0 || r.H <= 0 {
		t.Fatalf("the dialog page has no Motion.Dialog card (rect %+v)", r)
	}
	if !tt.HasText("320 ms") {
		t.Fatalf("the Motion.Dialog card does not show the published duration; texts=%q", tt.Texts())
	}
}
