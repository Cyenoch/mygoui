package navigation

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

func TestSidebarBoundedShellFixedHeaderFooterAndIndependentScroll(t *testing.T) {
	var s SidebarState
	items := make([]SidebarItem, 32)
	for i := range items {
		items[i] = SidebarItem{ID: SidebarID(fmt.Sprint(i)), Label: fmt.Sprintf("Destination %02d", i), Icon: lucide.File()}
	}
	var shell ui.Handle
	var shellBounds, navBounds ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Sidebar(c, &s, SidebarOptions{Label: "Workspace", Width: 256,
			Header: func(collapsed bool) {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Choose team", Icon: lucide.Building2()}).FillWidth()
			},
			Footer: func(collapsed bool) {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Account menu", Icon: lucide.User()}).FillWidth()
			},
			ContentHeader: func() { ui.Text(c, "Application header") },
			Content: func() {
				for i := range 40 {
					ui.Text(c, fmt.Sprintf("Content %02d", i)).Padding(theme.Space8)
				}
			},
		}, []SidebarGroup{{ID: "platform", Label: "Platform", Items: items}}).Height(360).Bind(&shell)
		shellBounds = shell.Bounds(c)
		navBounds = s.scroll.Bounds(c)
	}, 900, 500)
	tt.Frame()
	if shellBounds.H != 360 || shellBounds.W < 890 {
		t.Fatalf("shell initial bounds=%+v", shellBounds)
	}
	header, _ := tt.Find("Choose team")
	footer, _ := tt.Find("Account menu")
	toggle, _ := tt.Find("Toggle Workspace")
	if toggle.X < 256 || toggle.Y > 12 {
		t.Fatalf("toggle is not in external fixed main header: %+v", toggle)
	}
	if navBounds.Y < header.Y+header.H || navBounds.Y+navBounds.H > footer.Y+1 {
		t.Fatalf("navigation overlaps fixed header/footer: %+v %+v %+v", header, navBounds, footer)
	}
	if footer.Y+footer.H > 360 || footer.Y < 300 {
		t.Fatalf("footer not anchored in bounded sidebar: %+v", footer)
	}
	if s.scroll.MaxY <= 0 || s.contentScroll.MaxY <= 0 {
		t.Fatal("navigation and content did not get independently bounded scrolling")
	}
	tt.Scroll(navBounds.X+navBounds.W/2, navBounds.Y+navBounds.H/2, 0, 900)
	if s.scroll.Y <= 0 || s.contentScroll.Y != 0 {
		t.Fatal("navigation wheel scrolled wrong pane")
	}
	headerAfter, _ := tt.Find("Choose team")
	footerAfter, _ := tt.Find("Account menu")
	if headerAfter != header || footerAfter != footer {
		t.Fatal("navigation scrolling moved fixed team/account controls")
	}
	tt.Click("Choose team")
	tt.Click("Toggle Workspace")
	tt.Advance(250 * time.Millisecond)
	if !s.Collapsed {
		t.Fatal("external toggle did not collapse sidebar")
	}
	footerAfter, _ = tt.Find("Account menu")
	if footerAfter.W > 56 || footerAfter.H == 0 {
		t.Fatalf("collapsed footer disappeared or escaped rail: %+v", footerAfter)
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 2})
	tt.Click("Toggle Workspace")
	tt.Frame()
	headerAfter, _ = tt.Find("Choose team")
	footerAfter, _ = tt.Find("Account menu")
	if headerAfter.Y+headerAfter.H > footerAfter.Y || footerAfter.Y+footerAfter.H > 361 {
		t.Fatal("200% text clipped fixed header/footer")
	}
}

func TestSidebarNestedDisclosureFocusDisabledAndRail(t *testing.T) {
	s := SidebarState{Selected: "history"}
	open := true
	actions := 0
	groups := []SidebarGroup{{ID: "platform", Label: "Platform", Items: []SidebarItem{
		{ID: "playground", Label: "Playground", Icon: lucide.SquareTerminal(), Open: &open, Children: []SidebarItem{
			{ID: "history", Label: "History", OnSelect: func() { actions++ }},
			{ID: "disabled", Label: "Unavailable", Disabled: true, OnSelect: func() { actions += 100 }},
			{ID: "starred", Label: "Starred", OnSelect: func() { actions += 10 }},
		}},
		{ID: "models", Label: "Models", Icon: lucide.Bot()},
	}}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Sidebar(c, &s, SidebarOptions{Label: "Navigation"}, groups)
	}, 900, 500)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if err := tt.Click("History"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 || !tt.Focused("History") {
		t.Fatal("nested leaf did not route/focus")
	}
	open = false
	tt.Frame()
	if !tt.Focused("Playground") || tt.HasText("History") {
		t.Fatal("closing branch failed to restore focus or remove content")
	}
	tt.Key(0, ui.KeySpace)
	if !open || !tt.HasText("Starred") {
		t.Fatal("native disclosure keyboard failed")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("History") {
		t.Fatal("expanded branch did not enter native focus group")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Starred") {
		t.Fatal("native traversal failed to skip disabled child")
	}
	tt.Key(0, ui.KeyEnter)
	if actions != 11 || s.Selected != "starred" {
		t.Fatal("keyboard leaf navigation failed")
	}
	tt.Click("Toggle Navigation")
	if !s.Collapsed || tt.HasText("History") {
		t.Fatal("rail retained expanded text")
	}
	tt.Click("Playground")
	if s.Collapsed || !open || !tt.HasText("History") {
		t.Fatal("rail branch did not reopen expanded navigation")
	}
	tt.Click("Unavailable")
	if actions != 11 {
		t.Fatal("disabled leaf executed")
	}
	tt.SetDirection(ui.RightToLeft)
	branch, _ := tt.Find("Playground")
	child, _ := tt.Find("History")
	if child.X+child.W >= branch.X+branch.W {
		t.Fatalf("RTL nested branch lacks logical inset: %+v %+v", branch, child)
	}
}

func TestSidebarMobileSheetTrapDismissSelectionAndPersistentFooter(t *testing.T) {
	var s SidebarState
	actions := 0
	groups := []SidebarGroup{{ID: "platform", Label: "Platform", Items: []SidebarItem{
		{ID: "home", Label: "Home", Icon: lucide.House(), OnSelect: func() { actions++ }},
		{ID: "disabled", Label: "Unavailable", Disabled: true},
	}}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Sidebar(c, &s, SidebarOptions{Label: "Navigation", Header: func(collapsed bool) { mygoui.Button(c, mygoui.ButtonOptions{Label: "Team menu"}) },
			Footer: func(collapsed bool) { mygoui.Button(c, mygoui.ButtonOptions{Label: "Account menu"}) },
			Content: func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Main content action", OnClick: func() { actions += 100 }})
			},
		}, groups).Height(420)
	}, 390, 500)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 2})
	if tt.HasText("Platform") {
		t.Fatal("closed mobile sheet left sidebar visible")
	}
	if err := tt.Click("Toggle Navigation"); err != nil {
		t.Fatal(err)
	}
	if !s.Open || !tt.HasText("Platform") || !tt.HasText("Account menu") {
		t.Fatal("mobile sheet did not build actual navigation anatomy")
	}
	for range 8 {
		tt.Key(0, ui.KeyTab)
		if tt.Focused("Main content action") || tt.Focused("Toggle Navigation") {
			t.Fatal("mobile focus escaped sheet")
		}
	}
	header, _ := tt.Find("Team menu")
	footer, _ := tt.Find("Account menu")
	if header.Y+header.H >= footer.Y || math.Abs(float64(footer.Y+footer.H-492)) > 10 {
		t.Fatalf("sheet fixed controls not bounded at 200%%: %+v %+v", header, footer)
	}
	tt.Key(0, ui.KeyEscape)
	if s.Open || !tt.Focused("Toggle Navigation") {
		t.Fatalf("Escape dismissal: open=%v toggle=%v close=%v home=%v main=%v", s.Open, tt.Focused("Toggle Navigation"), tt.Focused("Close Navigation"), tt.Focused("Home"), tt.Focused("Main content action"))
	}
	tt.Click("Toggle Navigation")
	tt.Click("Home")
	if s.Open || s.Selected != "home" || actions != 1 {
		t.Fatal("mobile leaf did not route and dismiss sheet")
	}
	if !tt.Focused("Toggle Navigation") {
		t.Fatal("selection dismissal lost opener")
	}
}
