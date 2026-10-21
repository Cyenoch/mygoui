package navigation

import (
	"strings"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCommandProjectDetailSearchAndSourceAvailability(t *testing.T) {
	var state CommandState
	selected := ""
	items := []CommandItem{
		{ID: "alpha", Label: "Alpha", Detail: "Local · /projects/alpha", Icon: lucide.Folder(), Group: "Projects", OnSelect: func() { selected = "alpha" }},
		{ID: "beta", Label: "Beta", Detail: "Local · /projects/target", Icon: lucide.Folder(), Group: "Projects", OnSelect: func() { selected = "beta" }},
		{ID: "azure", Label: "Azure DevOps repository", Detail: "Clone project/repository", Badge: "Setup Required", DisabledReason: "Configure source control first", Disabled: true, OnSelect: func() { selected = "unavailable" }},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, &state, CommandOptions{Label: "Search projects", AutoFocus: true}, items).FillWidth()
	}, 390, 420)
	if !tt.Focused("Search projects") {
		t.Fatal("entry search did not receive requested focus")
	}
	if !tt.HasText("Setup Required") {
		t.Fatal("unavailable source has no visible setup explanation")
	}
	if err := tt.Click("Azure DevOps repository"); err != nil {
		t.Fatal(err)
	}
	if selected != "" {
		t.Fatal("unavailable source executed")
	}
	if err := tt.Click("Search projects"); err != nil {
		t.Fatal(err)
	}
	tt.Type("/projects/target")
	tt.Key(0, ui.KeyEnter)
	if selected != "beta" {
		t.Fatalf("project detail did not select exact target: %q", selected)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("Azure")
	tt.Key(0, ui.KeyEnter)
	if selected != "beta" {
		t.Fatal("filtered unavailable source executed")
	}
}

func TestCommandBackPreservesNativeEditingAndComposition(t *testing.T) {
	var state CommandState
	back, selected := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, &state, CommandOptions{Label: "Search sources", AutoFocus: true, OnBack: func() { back++ }}, []CommandItem{
			{ID: "source", Label: "New project", Keywords: "你", OnSelect: func() { selected++ }},
		})
	}, 420, 360)
	tt.Type("x")
	tt.Key(0, ui.KeyBackspace)
	if state.Query != "" || back != 0 {
		t.Fatalf("Backspace bypassed native editing: query=%q back=%d", state.Query, back)
	}
	tt.Compose("ni", 2)
	tt.Key(0, ui.KeyEnter)
	if selected != 0 || back != 0 {
		t.Fatal("IME candidate confirmation executed navigation")
	}
	tt.Type("你")
	tt.Command("selectAll")
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace)
	if back != 1 || selected != 0 {
		t.Fatalf("empty-query Backspace=%d action=%d", back, selected)
	}
	if err := tt.Click("Back"); err != nil {
		t.Fatal(err)
	}
	if back != 2 || !tt.Focused("Search sources") {
		t.Fatal("Back button did not restore search focus")
	}
}

func TestCommandProjectDetailsWrapWithinNarrowRow(t *testing.T) {
	defer theme.SetPreferences(theme.Preferences{})
	theme.SetPreferences(theme.Preferences{FontSize: 20})
	var state CommandState
	detail := strings.Repeat("project-directory/", 7) + "repository"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Command(c, &state, CommandOptions{Label: "Search", Height: 300}, []CommandItem{
			{ID: "target", Label: "Target repository", Detail: detail, Icon: lucide.Folder(), OnSelect: func() {}},
		}).FillWidth()
	}, 300, 420)
	row, rowOK := tt.Find("Target repository")
	text, textOK := tt.Find(detail)
	if !rowOK || !textOK || text.H <= 20 {
		t.Fatalf("project path did not wrap: row=%+v text=%+v", row, text)
	}
	if text.X < row.X || text.X+text.W > row.X+row.W+1 || row.X+row.W > 301 {
		t.Fatalf("wrapped project detail crossed row boundary: row=%+v text=%+v", row, text)
	}
}
