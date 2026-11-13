package form

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

var imageSelectItems = []ImageSelectItem{
	{Value: "light", Label: "Light"},
	{Value: "dark", Label: "Dark"},
	{Value: "system", Label: "System", Disabled: true},
}

func TestImageSelectSingleChoosesWithClickAndArrows(t *testing.T) {
	value, changes := "light", 0
	tt := ui.NewTester(func(c *ui.Context) {
		ImageSelect(c, ImageSelectOptions{
			AriaLabel: "Theme", Items: imageSelectItems, Value: &value,
			OnChange: func() { changes++ },
		})
	}, 520, 200)
	if err := tt.Click("Dark"); err != nil {
		t.Fatal(err)
	}
	if value != "dark" || changes != 1 {
		t.Fatalf("a click chose %q with %d changes, want dark and 1", value, changes)
	}
	tt.Key(0, ui.KeyLeft)
	if value != "light" {
		t.Fatalf("Left chose %q, want light", value)
	}
	if err := tt.Click("System"); err != nil {
		t.Fatal(err)
	}
	if value != "light" {
		t.Fatalf("a disabled tile chose %q", value)
	}
}

func TestImageSelectMultipleTogglesTiles(t *testing.T) {
	values := []string{"light"}
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ImageSelect(c, ImageSelectOptions{
			AriaLabel: "Themes", Items: imageSelectItems, Values: &values,
			OnChange: func() { changes++ },
		})
	}, 520, 200)
	if err := tt.Click("Dark"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, []string{"light", "dark"}) {
		t.Fatalf("a click left %v, want light and dark", values)
	}
	// Space toggles the focused tile off again.
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(values, []string{"light"}) || changes != 2 {
		t.Fatalf("Space left %v with %d changes, want [light] and 2", values, changes)
	}
}
