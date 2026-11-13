package gallery

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// layoutState is the caller-owned state behind the layout demos: each
// view's pane sizes and flags.
type layoutState struct {
	view int

	files, details                   float32
	filesCollapsed, detailsCollapsed bool

	fixedCollapsed bool
	sheetOpen      bool

	shellSidebar, shellInspector                   float32
	shellSidebarCollapsed, shellInspectorCollapsed bool
	shellSidebarOpen, shellInspectorOpen           bool
}

func layoutDemo(body func(c *ui.Context, st *layoutState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, StateOf[layoutState](s, "layout")) }
}

// layoutFrame is the window a layout demo stands in.
func layoutFrame(c *ui.Context, height float32, body func()) {
	t := theme.Colors(c)
	frame := ui.Column(c).FillWidth().MinWidth(0).Clip().Radius(theme.RadiusScale().LG).
		BorderWidth(theme.Hairline).BorderColor(t.Border).Background(t.Background)
	if height > 0 {
		frame.Height(theme.Scale(c, height))
	}
	frame.Children(body)
}

func layoutPane(c *ui.Context, title string, lines ...string) {
	ui.Column(c).Fill().MinWidth(0).Padding(theme.Space12).Gap(theme.Space8).Children(func() {
		mygoui.Typography(c, title, mygoui.TypographyOptions{Type: theme.TypeH6})
		for _, line := range lines {
			mygoui.Paragraph(c.Key(line), line, true)
		}
	})
}

func layoutIconButton(c *ui.Context, label string, icon *ui.SVG) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Icon: icon, IconOnly: true,
		Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
	})
}

func titleBarWindowBody(c *ui.Context, _ *layoutState) {
	layoutFrame(c, 0, func() {
		layout.TitleBar(c, layout.TitleBarOptions{
			Title: "Quarterly report", Borderless: true,
			Leading: func() { layoutIconButton(c, "Back", lucide.ChevronLeft()) },
			Trailing: func() {
				layoutIconButton(c, "Share", lucide.Share())
				layoutIconButton(c, "More", lucide.Ellipsis())
			},
		})
	})
}

func titleBarCenterBody(c *ui.Context, st *layoutState) {
	layoutFrame(c, 0, func() {
		layout.TitleBar(c, layout.TitleBarOptions{
			Label: "Document views", Borderless: true,
			Leading: func() { layoutIconButton(c, "Back", lucide.ChevronLeft()) },
			Center: func() {
				ui.Row(c).FillWidth().Justify(ui.Center).Children(func() {
					navigation.SegmentedControl(c, navigation.SegmentedControlOptions{
						Label: "View", Selected: &st.view,
						Items: []navigation.TabOptions{{Label: "Outline"}, {Label: "Pages"}, {Label: "Notes"}},
					})
				})
			},
			Trailing: func() { layoutIconButton(c, "Share", lucide.Share()) },
		})
	})
}

func splitViewResizableBody(c *ui.Context, st *layoutState) {
	layoutFrame(c, 320, func() {
		layout.TitleBar(c, layout.TitleBarOptions{
			Title: "Project", Span: layout.SpanEnd,
			Leading: func() {
				layout.PaneToggle(c, layout.PaneToggleOptions{
					Collapsed: &st.filesCollapsed, Compact: layout.CompactNever,
				})
			},
			Trailing: func() {
				layout.PaneToggle(c, layout.PaneToggleOptions{
					Collapsed: &st.detailsCollapsed, Compact: layout.CompactNever, End: true,
				})
			},
		})
		layout.SplitView(c, layout.SplitViewOptions{
			Compact: layout.CompactNever, ContentMin: 160,
			Start: &layout.PaneOptions{
				Label: "Files", Size: &st.files, Width: 180, Min: 120, Collapsed: &st.filesCollapsed,
				Body: func() { layoutPane(c, "Files", "main.go", "go.mod", "README.md") },
			},
			Content: func() { layoutPane(c, "main.go", "Drag a divider, or focus it and press the arrow keys.") },
			End: &layout.PaneOptions{
				Label: "Details", Size: &st.details, Width: 180, Min: 120, Collapsed: &st.detailsCollapsed,
				Body: func() { layoutPane(c, "Details", "Go source", "1.2 KB") },
			},
		})
	})
}

func splitViewFixedBody(c *ui.Context, st *layoutState) {
	layoutFrame(c, 260, func() {
		layout.TitleBar(c, layout.TitleBarOptions{
			Title: "Settings", Span: layout.SpanEnd,
			Leading: func() {
				layout.PaneToggle(c, layout.PaneToggleOptions{
					Collapsed: &st.fixedCollapsed, Compact: layout.CompactNever,
				})
			},
		})
		layout.SplitView(c, layout.SplitViewOptions{
			Compact: layout.CompactNever, ContentMin: 160,
			Start: &layout.PaneOptions{
				Label: "Sections", Width: 180, Fixed: true, Collapsed: &st.fixedCollapsed,
				Body: func() { layoutPane(c, "Sections", "General", "Appearance", "Shortcuts") },
			},
			Content: func() { layoutPane(c, "General", "The pane keeps its width behind a hairline.") },
		})
	})
}

func splitViewCompactBody(c *ui.Context, st *layoutState) {
	layoutFrame(c, 220, func() {
		layout.TitleBar(c, layout.TitleBarOptions{
			Title: "Inbox", Span: layout.SpanEnd,
			Leading: func() {
				layout.PaneToggle(c, layout.PaneToggleOptions{
					Open: &st.sheetOpen, Compact: layout.CompactAlways, Label: "Show mailboxes",
				})
			},
		})
		layout.SplitView(c, layout.SplitViewOptions{
			Compact: layout.CompactAlways,
			Start: &layout.PaneOptions{
				Label: "Mailboxes", Width: 240, Open: &st.sheetOpen,
				Body: func() { layoutPane(c, "Mailboxes", "Inbox", "Drafts", "Sent") },
			},
			Content: func() { layoutPane(c, "Inbox", "The pane opens over the window from its edge.") },
		})
	})
}

func scaffoldBody(c *ui.Context, st *layoutState) {
	layoutFrame(c, 400, func() {
		layout.Scaffold(c, layout.ScaffoldOptions{
			Compact: layout.CompactNever, ContentMin: 160,
			StartBar: layout.TitleBarOptions{Title: "Mail"},
			Start: &layout.PaneOptions{
				Label: "Mailboxes", Size: &st.shellSidebar, Width: 170, Min: 120,
				Collapsed: &st.shellSidebarCollapsed, Open: &st.shellSidebarOpen,
				Body: func() { layoutPane(c, "Mailboxes", "Inbox", "Drafts", "Sent") },
			},
			TitleBar: layout.TitleBarOptions{
				Title: "Inbox",
				Leading: func() {
					layout.PaneToggle(c, layout.PaneToggleOptions{
						Collapsed: &st.shellSidebarCollapsed, Open: &st.shellSidebarOpen,
						Compact: layout.CompactNever,
					})
				},
				Trailing: func() {
					layout.PaneToggle(c, layout.PaneToggleOptions{
						Collapsed: &st.shellInspectorCollapsed, Open: &st.shellInspectorOpen,
						Compact: layout.CompactNever, End: true,
					})
				},
			},
			Content: func() { layoutPane(c, "Release notes", "The sidebar shares the title band with the bar.") },
			End: &layout.PaneOptions{
				Label: "Details", Size: &st.shellInspector, Width: 170, Min: 120,
				Collapsed: &st.shellInspectorCollapsed, Open: &st.shellInspectorOpen,
				Body: func() { layoutPane(c, "Details", "3 messages") },
			},
			StatusBar: func() {
				navigation.StatusBar(c, navigation.StatusBarOptions{Message: "Up to date"})
			},
		})
	})
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "titlebar", Title: "TitleBar", Group: "Patterns", Order: 210,
		Description: "A draggable window title bar that keeps clear of the system window controls.",
		Examples: []ComponentExample{
			{Title: "Window bar", Description: "The bar leaves the room this window's controls take.",
				Demo: layoutDemo(titleBarWindowBody)},
			{Title: "Custom center", Demo: layoutDemo(titleBarCenterBody)},
		},
		Code: `// The window draws its own title bar:
// mygo.WindowOptions{TitleBarStyle: mygo.TitleBarHiddenInset, TitleBarHeight: 44}
layout.TitleBar(c, layout.TitleBarOptions{
    Title:     "Quarterly report",
    MinHeight: 44,
    Leading: func() {
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Back", Icon: lucide.ChevronLeft(), IconOnly: true,
            Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
        })
    },
    Trailing: func() { /* window-level actions */ },
})`,
		Params: []ComponentParam{
			{"Leading / Trailing", "func()", "nil", "Controls at the start and end; they keep their width and take their own presses."},
			{"Center", "func()", "nil", "Fills the room between them and aligns its own content."},
			{"Title", "string", "\"\"", "Shown when Center is nil; truncates."},
			{"Label", "string", "Title", "Accessible name of the bar."},
			{"Span", "layout.TitleBarSpan", "SpanWindow", "Which window edges the bar reaches: SpanWindow, SpanStart or SpanEnd."},
			{"MinHeight", "float32", "0", "The native band the window was created with, in window DIPs."},
			{"Borderless", "bool", "false", "Drops the hairline under the bar."},
		},
		Tokens: []string{"Background", "Foreground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "splitview", Title: "SplitView", Group: "Patterns", Order: 220,
		Description: "Side panes beside the content, with dividers, collapse and a compact sheet.",
		Examples: []ComponentExample{
			{Title: "Resizable panes", Description: "Enter or a double click on a divider collapses its pane.",
				Demo: layoutDemo(splitViewResizableBody)},
			{Title: "Fixed pane", Demo: layoutDemo(splitViewFixedBody)},
			{Title: "Compact", Demo: layoutDemo(splitViewCompactBody)},
		},
		Code: `var width float32
var collapsed, open bool

layout.SplitView(c, layout.SplitViewOptions{
    Start: &layout.PaneOptions{
        Label: "Files", Size: &width, Width: 240, Max: 400,
        Collapsed: &collapsed, Open: &open,
        Body: func() { /* the pane */ },
    },
    Content: func() { /* the content */ },
})

// In the title bar: flips collapsed, or open in a compact window.
layout.PaneToggle(c, layout.PaneToggleOptions{Collapsed: &collapsed, Open: &open})`,
		Params: []ComponentParam{
			{"Start / End", "*layout.PaneOptions", "nil", "Side panes in layout direction."},
			{"Content", "func()", "nil", "What the panes sit beside."},
			{"ContentMin", "float32", "240", "Width the content keeps while a divider moves."},
			{"Compact", "layout.CompactMode", "CompactAuto", "CompactAuto is compact below layout.CompactWidth; CompactNever and CompactAlways force a mode."},
			{"Pane: Label / Body", "string / func()", "\"\" / nil", "Accessible name of the pane and its divider, and its content."},
			{"Pane: Size", "*float32", "nil", "Caller-owned width; zero shows Width, nil is kept by the view."},
			{"Pane: Width / Min / Max", "float32", "256 / 160 / none", "Default width and the divider's bounds."},
			{"Pane: Fixed", "bool", "false", "Keeps Width behind a hairline, without a divider."},
			{"Pane: Collapsed / Open", "*bool", "nil", "Hides the docked pane; shows it as a sheet while compact."},
		},
		Tokens: []string{"Background", "Border", "Ring", "Popover"},
	})

	RegisterComponent(ComponentPage{
		ID: "scaffold", Title: "Scaffold", Group: "Patterns", Order: 230,
		Description: "A window shell: title bar, full-height sidebar, inspector and status strip.",
		Examples: []ComponentExample{
			{Title: "Window shell", Demo: layoutDemo(scaffoldBody)},
		},
		Code: `var sidebar float32
var collapsed, open bool

layout.Scaffold(c, layout.ScaffoldOptions{
    StartBar: layout.TitleBarOptions{Title: "Mail"},
    Start: &layout.PaneOptions{
        Label: "Mailboxes", Size: &sidebar,
        Collapsed: &collapsed, Open: &open,
        Body: func() { /* navigation */ },
    },
    TitleBar: layout.TitleBarOptions{
        Title: "Inbox", MinHeight: 44,
        Leading: func() {
            layout.PaneToggle(c, layout.PaneToggleOptions{Collapsed: &collapsed, Open: &open})
        },
    },
    Content: func() { /* the page */ },
    StatusBar: func() {
        navigation.StatusBar(c, navigation.StatusBarOptions{Message: "Ready"})
    },
})`,
		Params: []ComponentParam{
			{"TitleBar", "layout.TitleBarOptions", "", "The bar over the content."},
			{"Start / StartBar", "*layout.PaneOptions / layout.TitleBarOptions", "nil", "Full-height start pane and its share of the title band."},
			{"End", "*layout.PaneOptions", "nil", "Pane beside the content, under the title bar."},
			{"Content", "func()", "nil", "The window's content."},
			{"StatusBar", "func()", "nil", "Strip across the bottom."},
			{"ContentMin / Compact", "float32 / layout.CompactMode", "240 / CompactAuto", "As in SplitView."},
		},
		Tokens: []string{"Background", "Foreground", "Sidebar", "SidebarForeground", "Border"},
	})
}
