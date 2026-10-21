package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestNativeSelectUsesOSMenuWithGroupsAndDisabledChoices(t *testing.T) {
	value := "apple"
	readonly, disabled := false, false
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		NativeSelect(c, &value, NativeSelectOptions{
			ControlOptions: ControlOptions{Label: "Fruit", ReadOnly: readonly, Disabled: disabled},
			OnChange:       func(string) { changes++ },
		}, []SelectOption{
			{Value: "apple", Label: "Apple", Group: "Available"},
			{Value: "banana", Label: "Banana", Group: "Available"},
			{Value: "pear", Label: "Pear", Disabled: true},
		})
	}, 300, 100)
	tt.Click("Fruit")
	if len(tt.Menu()) == 0 {
		t.Fatal("select did not open platform menu")
	}
	if err := tt.ChooseMenuItem("Available", "Banana"); err != nil {
		t.Fatal(err)
	}
	if value != "banana" || changes != 1 {
		t.Fatalf("native grouped choice: %q/%d", value, changes)
	}
	tt.Click("Fruit")
	_ = tt.ChooseMenuItem("Pear")
	if value != "banana" || changes != 1 {
		t.Fatal("disabled native menu option accepted selection")
	}
	tt.CloseMenu()
	readonly = true
	tt.Frame()
	tt.Click("Fruit")
	_ = tt.ChooseMenuItem("Available", "Apple")
	if value != "banana" || changes != 1 {
		t.Fatal("read-only select accepted choice")
	}
	tt.CloseMenu()
	readonly, disabled = false, true
	tt.Frame()
	tt.Click("Fruit")
	if len(tt.Menu()) != 0 {
		t.Fatal("disabled select opened platform menu")
	}
}
