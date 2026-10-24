package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestButtonActivationAndPendingState(t *testing.T) {
	calls := 0
	loading, disabled := false, false
	var handle ui.Handle
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Save", Handle: &handle, Loading: loading, Disabled: disabled,
			OnClick: func() { calls++ },
		})
	}, 320, 160)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || !tt.Focused("Save") {
		t.Fatalf("pointer activation calls=%d focus=%v", calls, tt.Focused("Save"))
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if calls != 3 {
		t.Fatalf("keyboard activations=%d, want 3", calls)
	}
	loading = true
	tt.Frame()
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	if calls != 3 {
		t.Fatalf("pending control activated: %d", calls)
	}
	loading, disabled = false, true
	tt.Frame()
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("disabled control activated: %d", calls)
	}
}

func TestCardCompositionKeepsActionsIndependent(t *testing.T) {
	calls := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.Card(c, mygoui.CardOptions{}, func() {
			mygoui.CardHeader(c, func() {
				mygoui.CardTitle(c, "Scheduled reports")
				mygoui.CardDescription(c, "Choose a schedule and destination.")
			})
			mygoui.CardContent(c, func() {
				mygoui.Typography(c, "Every Monday", mygoui.TypographyOptions{Type: theme.TypeBodySM})
			})
			mygoui.CardFooter(c, func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Edit", Variant: mygoui.ButtonSecondary,
					OnClick: func() { calls++ }})
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: func() { calls += 10 }})
			})
		})
	}, 360, 420)
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	if err := tt.Click("Edit"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if calls != 11 {
		t.Fatalf("card action dispatch=%d, want 11", calls)
	}
	if !tt.HasText("Scheduled reports") || !tt.HasText("Every Monday") {
		t.Fatal("card content missing")
	}
}

func TestChipAndBadgeKeepAccessibleText(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Row(c).Children(func() {
			mygoui.Chip(c, mygoui.ChipOptions{Label: "Beta", Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft})
			mygoui.Badge(c, mygoui.BadgeOptions{Label: "9", Color: mygoui.ColorDanger, Variant: mygoui.ChipPrimary})
		})
	}, 320, 120)
	if !tt.HasText("Beta") || !tt.HasText("9") {
		t.Fatal("chip or badge lost its label")
	}
}

func TestToggleButtonUpdatesCallerStateWithoutCallback(t *testing.T) {
	selected := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{Label: "Bold", Selected: &selected})
	}, 320, 160)
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	if !selected {
		t.Fatal("toggle did not update the caller-owned bool")
	}
}

func TestFullWidthButtonGroupSharesWidth(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Padding(20).Children(func() {
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{FullWidth: true},
				mygoui.ButtonOptions{Label: "Left", OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Right", OnClick: func() {}})
		})
	}, 400, 120)
	left, okLeft := tt.Find("Left")
	right, okRight := tt.Find("Right")
	if !okLeft || !okRight {
		t.Fatal("group labels missing")
	}
	if right.X <= left.X || right.Y != left.Y {
		t.Fatalf("group buttons are not horizontal: left=%v right=%v", left, right)
	}
	if left.W > 180 || right.W > 180 {
		t.Fatalf("buttons did not share width: left=%v right=%v", left, right)
	}
}

func TestProgressDefaultsToZeroToHundred(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ProgressBar(c, mygoui.ProgressOptions{
			Label: "Upload", Value: 60, ShowOutput: true,
		})
	}, 360, 120)
	if !tt.HasText("60%") {
		t.Fatalf("default progress output missing; texts: %v", tt.Texts())
	}
}
