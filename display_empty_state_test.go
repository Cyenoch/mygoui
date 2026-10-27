package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestEmptyStateDefaultsToTheLocalePlaceholder pins the default copy: an
// empty state with nothing to say still names what is missing, while caller
// content replaces it.
func TestEmptyStateDefaultsToTheLocalePlaceholder(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.EmptyState(c, mygoui.EmptyStateOptions{})
	}, 240, 80)
	placeholder := locale.M().Empty
	if placeholder == "" {
		t.Fatal("the locale publishes no empty placeholder")
	}
	if !tt.HasText(placeholder) {
		t.Fatalf("empty state shows %v, want %q", tt.Texts(), placeholder)
	}

	tt = ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.EmptyState(c, mygoui.EmptyStateOptions{
			Title: "No saved searches", Description: "Save one to see it here.",
		})
	}, 240, 80)
	if tt.HasText(placeholder) {
		t.Fatalf("caller copy kept the locale placeholder: %v", tt.Texts())
	}
	if !tt.HasText("No saved searches") {
		t.Fatalf("empty state shows %v, want the caller title", tt.Texts())
	}
}
