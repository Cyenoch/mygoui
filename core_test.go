package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
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

func TestButtonLoadingPreservesNaturalWidthAndSiblingPosition(t *testing.T) {
	for _, tc := range []struct {
		name      string
		icons     bool
		iconOnly  bool
		custom    bool
		textScale float32
	}{
		{name: "label", textScale: 1},
		{name: "icons and label", icons: true, textScale: 1},
		{name: "icon only", icons: true, iconOnly: true, textScale: 1},
		{name: "custom content", custom: true, textScale: 1},
		{name: "enlarged text", icons: true, textScale: 1.5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pending := false
			calls := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Row(c).Padding(20).Gap(12).Children(func() {
					options := mygoui.ButtonOptions{
						Label: "Export selected records", Loading: pending, IconOnly: tc.iconOnly,
						OnClick: func() { calls++ },
					}
					if tc.icons {
						options.Icon, options.EndIcon = lucide.Download(), lucide.ChevronDown()
					}
					if tc.custom {
						options.Content = func() { ui.Text(c, "Selection: 123 records") }
					}
					mygoui.Button(c, options)
					ui.Text(c, "After")
				})
			}, 600, 160)
			tt.SetPreferences(ui.Preferences{TextScale: tc.textScale})
			rest := findRect(t, tt, "Export selected records")
			sibling := findRect(t, tt, "After")
			pending = true
			tt.Frame()
			if got := findRect(t, tt, "Export selected records"); got != rest {
				t.Fatalf("pending geometry = %+v, want %+v", got, rest)
			}
			if got := findRect(t, tt, "After"); got != sibling {
				t.Fatalf("pending shifted sibling from %+v to %+v", sibling, got)
			}
			spinner := findRect(t, tt, locale.M().Loading)
			if spinner.X+spinner.W/2 != rest.X+rest.W/2 ||
				spinner.Y+spinner.H/2 != rest.Y+rest.H/2 {
				t.Fatalf("pending spinner is not centered: spinner=%+v button=%+v", spinner, rest)
			}
			if err := tt.Click("Export selected records"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyEnter)
			if calls != 0 {
				t.Fatalf("pending button activated %d times", calls)
			}
			pending = false
			tt.Frame()
			if got := findRect(t, tt, "Export selected records"); got != rest {
				t.Fatalf("restored geometry = %+v, want %+v", got, rest)
			}
			if err := tt.Click("Export selected records"); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatalf("restored button activations = %d, want 1", calls)
			}
		})
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
			mygoui.Badge(c, mygoui.BadgeOptions{Label: "9", Variant: mygoui.BadgeDestructive})
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
