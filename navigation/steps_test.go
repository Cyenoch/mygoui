package navigation

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func stepItems() []StepItem {
	return []StepItem{{Title: "Account"}, {Title: "Profile"}, {Title: "Review"}}
}

func TestStepsShowCompletedCurrentAndUpcoming(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Steps(c, StepsOptions{Label: "Signup", Items: stepItems(), Current: 1})
	}, 480, 120)
	// A completed step replaces its number with a check.
	if tt.HasText("1") || !tt.HasText("2") || !tt.HasText("3") {
		t.Fatalf("step numbers shown: %q", tt.Texts())
	}
	for title, want := range map[string]ui.Current{
		"Account": ui.CurrentNone, "Profile": ui.CurrentStep, "Review": ui.CurrentNone,
	} {
		info, ok := tt.Access(title)
		if !ok || info.Current != want {
			t.Errorf("%s: found %v, current %v, want %v", title, ok, info.Current, want)
		}
	}
}

func TestStepsSelectReportsThePressedStep(t *testing.T) {
	var selected []int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Steps(c, StepsOptions{
			Items: stepItems(), Current: 0, Orientation: StepsVertical,
			OnSelect: func(i int) { selected = append(selected, i) },
		})
	}, 320, 240)
	if err := tt.Click("Review"); err != nil {
		t.Fatal(err)
	}
	if info, ok := tt.AccessByRole(ui.RoleButton, "Review"); !ok || !tt.Focused("Review") {
		t.Fatalf("pressable step: button %v (%+v), focused %v", ok, info.Role, tt.Focused("Review"))
	}
	tt.Key(0, ui.KeySpace)
	if len(selected) != 2 || selected[0] != 2 || selected[1] != 2 {
		t.Fatalf("selected %v, want [2 2]", selected)
	}
}
