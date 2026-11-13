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

// TestStepsArrowsMoveAlongPressableSteps protects keyboard movement among
// pressable steps: the arrows follow the steps, in reading order, and Home
// and End reach the ends.
func TestStepsArrowsMoveAlongPressableSteps(t *testing.T) {
	for _, tc := range []struct {
		name        string
		dir         ui.Direction
		orientation StepsOrientation
		forth, back ui.Key
	}{
		{"horizontal", ui.LeftToRight, StepsHorizontal, ui.KeyRight, ui.KeyLeft},
		{"right to left", ui.RightToLeft, StepsHorizontal, ui.KeyLeft, ui.KeyRight},
		{"vertical", ui.LeftToRight, StepsVertical, ui.KeyDown, ui.KeyUp},
	} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			Steps(c, StepsOptions{Items: stepItems(), Orientation: tc.orientation, OnSelect: func(int) {}})
		}, 480, 240)
		tt.SetDirection(tc.dir)
		if err := tt.Click("Account"); err != nil {
			t.Fatal(err)
		}
		for _, step := range []struct {
			key  ui.Key
			want string
		}{{tc.forth, "Profile"}, {tc.back, "Account"}, {tc.back, "Review"}, {ui.KeyHome, "Account"}, {ui.KeyEnd, "Review"}} {
			tt.Key(0, step.key)
			if !tt.Focused(step.want) {
				t.Fatalf("%s: key %v did not focus %s", tc.name, step.key, step.want)
			}
		}
	}
}
