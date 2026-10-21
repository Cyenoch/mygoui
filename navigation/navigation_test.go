package navigation

import (
	"slices"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

func TestAccordionSelectionReorderAndDisabledKeyboard(t *testing.T) {
	var s AccordionState[string]
	mode := AccordionSingle
	items := []AccordionItem[string]{
		{Value: "a", Label: "Alpha", Content: func() {}},
		{Value: "x", Label: "Unavailable", Disabled: true},
		{Value: "b", Label: "Beta", Content: func() {}},
	}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); Accordion(c, &s, AccordionOptions{Mode: mode}, items) }, 420, 340)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Beta") {
		t.Fatal("arrow did not skip disabled header")
	}
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(s.Open, []string{"b"}) {
		t.Fatalf("single open=%v", s.Open)
	}
	slices.Reverse(items)
	tt.Frame()
	if !tt.Focused("Beta") || !slices.Equal(s.Open, []string{"b"}) {
		t.Fatal("reorder lost stable identity")
	}
	mode = AccordionMultiple
	tt.Frame()
	tt.Click("Alpha")
	if len(s.Open) != 2 {
		t.Fatalf("multiple open=%v", s.Open)
	}
	tt.Click("Unavailable")
	if len(s.Open) != 2 {
		t.Fatal("disabled header changed open values")
	}
}

func TestCollapsibleExternalCloseReturnsFocusAndDisablesExit(t *testing.T) {
	s := CollapsibleState{Open: true}
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Collapsible(c, &s, CollapsibleOptions{Label: "Details"}, func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Inside", OnClick: func() { actions++ }}) })
	}, 420, 340)
	tt.Click("Inside")
	if !tt.Focused("Inside") {
		t.Fatal("content did not receive focus")
	}
	s.Open = false
	tt.Frame()
	if !tt.Focused("Details") {
		t.Fatal("external close failed to restore trigger focus")
	}
	tt.Click("Inside")
	tt.Advance(250 * time.Millisecond)
	if actions != 1 || tt.HasText("Inside") {
		t.Fatal("outgoing content was actionable or remained after transition")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Click("Details")
	if !tt.HasText("Inside") {
		t.Fatal("reduced motion disclosure failed to open")
	}
}

func TestBreadcrumbOverflowAndPassiveCurrent(t *testing.T) {
	var s BreadcrumbState
	selected := ""
	items := []BreadcrumbItem{
		{ID: "root", Label: "Root", OnSelect: func() { selected = "root" }},
		{ID: "one", Label: "First", OnSelect: func() { selected = "first" }},
		{ID: "two", Label: "Second", OnSelect: func() { selected = "second" }},
		{ID: "three", Label: "Third", OnSelect: func() { selected = "third" }},
		{ID: "current", Label: "Current", OnSelect: func() { selected = "bad" }},
	}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); Breadcrumb(c, &s, BreadcrumbOptions{MaxVisible: 3}, items) }, 600, 300)
	tt.Click("Current")
	if selected != "" {
		t.Fatal("current page activated")
	}
	tt.Click("More ancestors")
	tt.Click("Second")
	if selected != "second" {
		t.Fatalf("overflow action=%q", selected)
	}
}

func TestSelectCustomPopupTypeaheadDisabledAndReorder(t *testing.T) {
	var s SelectState
	value := "a"
	changes := 0
	choices := []Choice[string]{{Value: "a", Label: "Apple", Group: "Fruit"}, {Value: "x", Label: "Blocked", Disabled: true}, {Value: "b", Label: "Banana", Group: "Fruit"}, {Value: "c", Label: "Cherry", Group: "Fruit"}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Select(c, &s, &value, SelectOptions[string]{Label: "Fruit", Width: 240, OnChange: func(string) { changes++ }}, choices)
	}, 420, 400)
	tt.Click("Fruit")
	if len(tt.Menu()) != 0 || !tt.HasText("Banana") {
		t.Fatal("Select did not use the custom popup")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if value != "b" || changes != 1 {
		t.Fatalf("keyboard disabled skip=%s/%d", value, changes)
	}
	slices.Reverse(choices)
	tt.Frame()
	tt.Click("Fruit")
	tt.Type("ch")
	tt.Key(0, ui.KeyEnter)
	if value != "c" {
		t.Fatalf("typeahead after reorder=%q", value)
	}
	tt.Click("Fruit")
	tt.Key(0, ui.KeyEscape)
	if !tt.Focused("Fruit") {
		t.Fatal("Escape lost select trigger focus")
	}
}

func TestComboboxFilterChipsRemovalAndEmpty(t *testing.T) {
	var s ComboboxState
	var values []string
	choices := []Choice[string]{{Value: "a", Label: "Apple"}, {Value: "x", Label: "Apricot", Disabled: true}, {Value: "b", Label: "Banana"}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Combobox(c, &s, &values, ComboboxOptions[string]{Label: "Fruit", Multiple: true, Width: 260}, choices)
	}, 500, 400)
	tt.Click("Fruit")
	tt.Type("ban")
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(values, []string{"b"}) {
		t.Fatalf("filtered choice=%v", values)
	}
	tt.Click("Fruit")
	tt.Type("app")
	tt.Key(0, ui.KeyEnter)
	if !slices.Equal(values, []string{"b", "a"}) {
		t.Fatalf("multiple choice=%v", values)
	}
	tt.Click("Remove Banana")
	if !slices.Equal(values, []string{"a"}) || !tt.Focused("Fruit") {
		t.Fatal("chip removal lost focus or value")
	}
	tt.Key(0, ui.KeyBackspace)
	if len(values) != 0 {
		t.Fatal("empty-query Backspace failed to remove chip")
	}
	tt.Type("zzzz")
	if !tt.HasText("No matches") {
		t.Fatal("empty feedback missing")
	}
	tt.Key(0, ui.KeyEnter)
	if len(values) != 0 {
		t.Fatal("empty state committed phantom choice")
	}
}

func TestCommandRanksActualActionsShortcutsAndEmpty(t *testing.T) {
	var s CommandState
	action := ""
	items := []CommandItem{
		{ID: "contains", Label: "Recent files", Group: "Files", OnSelect: func() { action = "recent" }},
		{ID: "prefix", Label: "File search", Group: "Files", Shortcut: "⌘F", Mods: ui.Cmd, Key: ui.KeyF, OnSelect: func() { action = "search" }},
		{ID: "disabled", Label: "File delete", Disabled: true, OnSelect: func() { action = "bad" }},
		{ID: "settings", Label: "Settings", Group: "Application", Keywords: "preferences", OnSelect: func() { action = "settings" }},
	}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); Command(c, &s, CommandOptions{Label: "Commands"}, items) }, 500, 420)
	tt.Click("Commands")
	tt.Type("file")
	tt.Key(0, ui.KeyEnter)
	if action != "search" {
		t.Fatalf("prefix ranking=%q", action)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("preferences")
	tt.Key(0, ui.KeyEnter)
	if action != "settings" {
		t.Fatalf("keyword action=%q", action)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("zzzz")
	if !tt.HasText("No commands found") {
		t.Fatal("command empty feedback missing")
	}
	tt.Key(0, ui.KeyEnter)
	if action != "settings" {
		t.Fatal("empty command ran action")
	}
	s.Query = ""
	tt.Frame()
	tt.Key(ui.Cmd, ui.KeyF)
	if action != "search" {
		t.Fatal("displayed shortcut did not execute")
	}
}
func TestCommandDialogEscapeRestoresOpener(t *testing.T) {
	var s CommandState
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{Label: "Open commands", OnClick: func() { open = true }})
		CommandDialog(c, &open, &s, overlay.DialogOptions{Title: "Palette"}, CommandOptions{Label: "Search commands"}, nil)
	}, 500, 400)
	tt.Click("Open commands")
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if open || !tt.Focused("Open commands") {
		t.Fatal("command dialog failed to dismiss/restore focus")
	}
}

func TestNavigationMenuKeyboardAndRealLeafActions(t *testing.T) {
	var s NavigationMenuState
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		NavigationMenu(c, &s, []NavigationMenuItem{{ID: "docs", Label: "Docs", Items: []overlay.MenuItem{{ID: "guide", Label: "Guide", OnSelect: func() { actions++ }}}}, {ID: "off", Label: "Disabled", Disabled: true}, {ID: "home", Label: "Home", OnSelect: func() { actions += 10 }}})
	}, 600, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if actions != 1 {
		t.Fatal("popup keyboard action did not execute")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Home") {
		t.Fatal("top navigation did not skip disabled entry")
	}
	tt.Key(0, ui.KeyEnter)
	if actions != 11 {
		t.Fatal("leaf navigation action failed")
	}
}

func TestPaginationBoundariesAndActualPages(t *testing.T) {
	page := 1
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Pagination(c, PaginationOptions{Page: page, Pages: 12, OnChange: func(n int) { page = n; changes++ }})
	}, 390, 220)
	tt.Click("Previous")
	if page != 1 || changes != 0 {
		t.Fatal("previous boundary activated")
	}
	tt.Click("12")
	if page != 12 {
		t.Fatalf("last page=%d", page)
	}
	tt.Click("Next")
	if changes != 1 {
		t.Fatal("next boundary activated")
	}
	tt.Click("12")
	if changes != 1 {
		t.Fatal("current page redundantly activated")
	}
	tt.Click("Previous")
	if page != 11 {
		t.Fatal("previous failed to change actual page")
	}
}

func TestSidebarCollapsedExpandedNarrowActions(t *testing.T) {
	var s SidebarState
	actions := 0
	groupOpen := true
	groups := []SidebarGroup{{ID: "main", Label: "Workspace", Open: &groupOpen, Items: []SidebarItem{{ID: "files", Label: "Files", OnSelect: func() { actions++ }}, {ID: "disabled", Label: "Unavailable", Disabled: true}, {ID: "settings", Label: "Settings", Action: &mygoui.ButtonOptions{Label: "Configure", OnClick: func() { actions += 10 }}}}}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Sidebar(c, &s, SidebarOptions{Label: "Navigation", NarrowAt: 600}, groups)
	}, 800, 500)
	tt.Click("Files")
	tt.Click("Configure")
	if s.Selected != "files" || actions != 11 {
		t.Fatal("sidebar actions did not execute independently")
	}
	tt.Click("Toggle Navigation")
	if !s.Collapsed {
		t.Fatal("sidebar did not collapse")
	}
	tt.Click("Settings")
	if s.Selected != "settings" {
		t.Fatal("collapsed item did not select")
	}
	tt.SetSize(390, 500)
	tt.Click("Toggle Navigation")
	if !s.Open || !tt.HasText("Workspace") {
		t.Fatal("narrow sidebar did not open real sheet")
	}
	tt.Click("Files")
	if s.Open || s.Selected != "files" || actions != 12 {
		t.Fatal("narrow item failed to select and dismiss")
	}
}

func TestTabsRoveSkipDisabledAndRetainValueAcrossReorder(t *testing.T) {
	value := "a"
	vertical := false
	tabs := []Tab[string]{{Value: "a", Label: "Alpha", Content: func() {}}, {Value: "x", Label: "Unavailable", Disabled: true, Content: func() { t.Error("disabled panel constructed") }}, {Value: "b", Label: "Beta", Content: func() {}}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		orientation := ui.Horizontal
		if vertical {
			orientation = ui.Vertical
		}
		Tabs(c, &value, TabsOptions[string]{Label: "Views", Variant: TabsLine, Orientation: orientation}, tabs)
		mygoui.Button(c, mygoui.ButtonOptions{Label: "After", OnClick: func() {}})
	}, 600, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if value != "b" || !tt.Focused("Beta") {
		t.Fatalf("roving selection=%q", value)
	}
	slices.Reverse(tabs)
	tt.Frame()
	if value != "b" || !tt.Focused("Beta") {
		t.Fatal("reorder lost selected tab identity")
	}
	vertical = true
	tt.Frame()
	tt.Key(0, ui.KeyEnd)
	if value != "a" {
		t.Fatal("vertical End did not skip disabled")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("tabs were more than one Tab stop")
	}
}
