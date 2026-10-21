package navigation

import (
	"strconv"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

func TestSelectEndRevealsLastChoiceAndScrollControlsWork(t *testing.T) {
	var s SelectState
	value := 0
	choices := make([]Choice[int], 40)
	for i := range choices {
		choices[i] = Choice[int]{Value: i, Label: "Choice " + strconv.Itoa(i)}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Select(c, &s, &value, SelectOptions[int]{Label: "Choices", Width: 240}, choices)
	}, 420, 400)
	tt.Click("Choices")
	tt.Key(0, ui.KeyEnd)
	if s.scroll.Y <= 0 {
		t.Fatal("End did not scroll highlight into view")
	}
	tt.Key(0, ui.KeyEnter)
	if value != 39 {
		t.Fatalf("End committed %d instead of last value", value)
	}
	tt.Click("Choices")
	before := s.scroll.Y
	tt.Click("Scroll choices up")
	if s.scroll.Y >= before {
		t.Fatal("scroll-up control did not move actual viewport")
	}
}

func TestAccordionArrowsIgnoreFocusablePanelChildren(t *testing.T) {
	s := AccordionState[string]{Open: []string{"a"}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Accordion(c, &s, AccordionOptions{}, []AccordionItem[string]{
			{Value: "a", Label: "Alpha", Content: func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Inside", OnClick: func() {}}) }},
			{Value: "b", Label: "Beta"},
		})
	}, 500, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Beta") {
		t.Fatal("header arrow entered panel focus targets")
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Inside") {
		t.Fatal("Tab failed to enter expanded panel content")
	}
}

func TestTabsLargeTextWrapWithinAvailableWidth(t *testing.T) {
	value := "account"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Width(308).Children(func() {
			Tabs(c, &value, TabsOptions[string]{Label: "Settings"}, []Tab[string]{{Value: "account", Label: "Account"}, {Value: "security", Label: "Security"}, {Value: "unavailable", Label: "Unavailable", Disabled: true}})
		})
	}, 390, 500)
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	for _, label := range []string{"Account", "Security", "Unavailable"} {
		r, ok := tt.Find(label)
		if !ok || r.X < 0 || r.X+r.W > 309 {
			t.Fatalf("large-text tab %q outside header: %+v", label, r)
		}
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if value != "security" {
		t.Fatal("wrapped header lost native roving selection")
	}
}
