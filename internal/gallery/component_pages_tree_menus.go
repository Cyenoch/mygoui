package gallery

import (
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// treeMenusState is the caller-owned state behind the tree, menu bar, split
// button and overflow list pages.
type treeMenusState struct {
	ready bool

	files           []list.TreeNode
	fileTree        list.TreeState
	checkTree       list.TreeState
	checked         int
	catalog         []list.TreeNode
	catalogTree     list.TreeState
	guideTree       list.TreeState
	lazy            []list.TreeNode
	lazyTree        list.TreeState
	lazyLoading     bool
	lastFileAction  string
	menubarOpen     string
	lastCommand     string
	lastSplitAction string
	lastOverflow    string
	overflowBold    bool
	splitChoice     string
	overflowItalic  bool
	splitLoading    bool
}

func treeMenusStateOf(s *State) *treeMenusState {
	st := StateOf[treeMenusState](s, "tree-menus")
	if !st.ready {
		st.ready = true
		st.files = treeDemoFiles()
		st.fileTree.SetExpanded("src", true)
		st.fileTree.Selected.Add("src/main.go")
		st.guideTree.SetExpanded("src", true)
		st.guideTree.SetExpanded("src/ui", true)
		st.lazy = []list.TreeNode{{Key: "remote", Label: "remote", Icon: lucide.Folder(), Lazy: true}}
		st.checkTree.SetExpanded("src", true)
		st.checkTree.SetChecked("src/main.go", true)
		st.checked = 1
		st.catalog = treeDemoCatalog(40, 250)
		st.catalogTree.SetExpanded("Group 1", true)
	}
	return st
}

func treeMenusDemo(body func(c *ui.Context, st *treeMenusState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, treeMenusStateOf(s)) }
}

func treeDemoFiles() []list.TreeNode {
	file := func(dir, name string) list.TreeNode {
		return list.TreeNode{Key: dir + name, Label: name, Icon: lucide.FileText()}
	}
	return []list.TreeNode{
		{Key: "src", Label: "src", Icon: lucide.Folder(), Children: []list.TreeNode{
			{Key: "src/ui", Label: "ui", Icon: lucide.Folder(), Children: []list.TreeNode{
				file("src/ui/", "button.go"), file("src/ui/", "menu.go"),
			}},
			file("src/", "main.go"),
			file("src/", "routes.go"),
		}},
		{Key: "docs", Label: "docs", Icon: lucide.Folder(), Children: []list.TreeNode{
			file("docs/", "design.md"),
		}},
		{Key: "vendor", Label: "vendor", Icon: lucide.Folder(), Disabled: true, Children: []list.TreeNode{
			file("vendor/", "modules.txt"),
		}},
		file("", "go.mod"),
	}
}

func treeDemoCatalog(groups, perGroup int) []list.TreeNode {
	nodes := make([]list.TreeNode, groups)
	for g := range nodes {
		group := "Group " + strconv.Itoa(g+1)
		kids := make([]list.TreeNode, perGroup)
		for i := range kids {
			label := "Item " + strconv.Itoa(g*perGroup+i+1)
			kids[i] = list.TreeNode{Key: group + "/" + label, Label: label}
		}
		nodes[g] = list.TreeNode{Label: group, Children: kids}
	}
	return nodes
}

func treeDemoFrame(c *ui.Context, body func()) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 320)).
		Border(theme.Hairline, t.Border).Radius(theme.RadiusScale().XL).Clip().Children(body)
}

func treeFilesBody(c *ui.Context, st *treeMenusState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		treeDemoFrame(c, func() {
			list.Tree(c, list.TreeOptions{
				Label: "Project files", Nodes: st.files, State: &st.fileTree,
				Selection: list.SelectionMultiple, Height: theme.Scale(c, 248),
				OnContextMenu: func(key string, m *ui.Menu) {
					for _, action := range [...]string{"Rename", "Duplicate", "Delete"} {
						if m.Item(action).Chosen() {
							st.lastFileAction = action + " " + key
						}
					}
				},
			})
		})
		if st.lastFileAction != "" {
			formsMutedText(c, st.lastFileAction)
		}
	})
}

func treeGuidesBody(c *ui.Context, st *treeMenusState) {
	treeDemoFrame(c, func() {
		list.Tree(c, list.TreeOptions{
			Label: "Project files with guides", Nodes: st.files, State: &st.guideTree,
			Selection: list.SelectionSingle, Guides: true, Height: theme.Scale(c, 248),
		})
	})
}

func treeLazyBody(c *ui.Context, st *treeMenusState) {
	st.lazy[0].Loading = st.lazyLoading
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		treeDemoFrame(c, func() {
			list.Tree(c, list.TreeOptions{
				Label: "Remote files", Nodes: st.lazy, State: &st.lazyTree,
				Selection: list.SelectionSingle, MaxHeight: theme.Scale(c, 248),
				OnExpandedChange: func(key string, open bool) {
					if open && key == "remote" && len(st.lazy[0].Children) == 0 {
						st.lazyLoading = true
					}
				},
			})
		})
		if st.lazyLoading {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Finish loading", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					st.lazy[0].Children = []list.TreeNode{
						{Key: "remote/a.txt", Label: "a.txt", Icon: lucide.FileText()},
						{Key: "remote/b.txt", Label: "b.txt", Icon: lucide.FileText()},
					}
					st.lazyLoading = false
					st.lazyTree.Revision++
				},
			})
		}
	})
}

func treeCheckableBody(c *ui.Context, st *treeMenusState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		treeDemoFrame(c, func() {
			list.Tree(c, list.TreeOptions{
				Label: "Files to export", Nodes: st.files, State: &st.checkTree,
				Checkable: true, Height: theme.Scale(c, 248),
				OnCheckedChange: func(keys []string) { st.checked = len(keys) },
			})
		})
		formsMutedText(c, "Checked: "+strconv.Itoa(st.checked))
	})
}

func treeLargeBody(c *ui.Context, st *treeMenusState) {
	treeDemoFrame(c, func() {
		list.Tree(c, list.TreeOptions{
			Label: "Catalog", Nodes: st.catalog, State: &st.catalogTree,
			Selection: list.SelectionMultiple, Height: theme.Scale(c, 248),
		})
	})
}

func menubarDemoMenus(st *treeMenusState) []navigation.MenubarMenuOptions {
	menu := func(sections ...navigation.MenuSectionOptions) navigation.MenuOptions {
		return navigation.MenuOptions{Sections: sections, OnSelect: func(key string) { st.lastCommand = key }}
	}
	return []navigation.MenubarMenuOptions{
		{Label: "File", Mnemonic: 'F', Menu: menu(
			navigation.MenuSectionOptions{Key: "new", Items: []navigation.MenuItemOptions{
				{Key: "New Tab", Label: "New Tab", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyT}},
				{Key: "New Window", Label: "New Window", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyN}},
				{Label: "Open Recent", Submenu: []navigation.MenuItemOptions{
					{Label: "roadmap.md"}, {Label: "budget.csv"},
				}},
			}},
			navigation.MenuSectionOptions{Key: "close", Items: []navigation.MenuItemOptions{
				{Separator: true},
				{Key: "Close Tab", Label: "Close Tab", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyW}},
			}},
		)},
		{Label: "Edit", Mnemonic: 'E', Menu: menu(navigation.MenuSectionOptions{Items: []navigation.MenuItemOptions{
			{Key: "Undo", Label: "Undo", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyZ}},
			{Key: "Redo", Label: "Redo", Accelerator: form.Hotkey{Mods: ui.Cmd | ui.Shift, Key: ui.KeyZ}},
			{Separator: true},
			{Key: "Cut", Label: "Cut", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyX}},
			{Key: "Copy", Label: "Copy", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyC}},
			{Key: "Paste", Label: "Paste", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyV}, Disabled: true},
		}})},
		{Label: "View", Mnemonic: 'V', Menu: menu(navigation.MenuSectionOptions{Items: []navigation.MenuItemOptions{
			{Label: "Zoom In", Shortcut: "⌘+"},
			{Label: "Zoom Out", Shortcut: "⌘−"},
			{Label: "Full Screen", Shortcut: "⌃⌘F"},
		}})},
		{Label: "Tools", Disabled: true},
	}
}

func menubarBody(c *ui.Context, st *treeMenusState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		navigation.Menubar(c, navigation.MenubarOptions{
			Label: "Application", Open: &st.menubarOpen, Menus: menubarDemoMenus(st),
		})
		if st.lastCommand != "" {
			formsMutedText(c, "Last command: "+st.lastCommand)
		}
	})
}

func splitButtonDemoMenu(st *treeMenusState) navigation.MenuOptions {
	return navigation.MenuOptions{
		Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
			{Label: "Save As…", Shortcut: "⇧⌘S"},
			{Label: "Save a Copy…"},
			{Label: "Save All"},
		}}},
		OnSelect: func(key string) { st.lastSplitAction = key },
	}
}

func splitButtonVariantsBody(c *ui.Context, st *treeMenusState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		ui.Row(c).Gap(theme.Space12).Wrap().AlignItems(ui.Center).Children(func() {
			for _, variant := range []mygoui.ButtonVariant{mygoui.ButtonPrimary, mygoui.ButtonSecondary, mygoui.ButtonOutline} {
				c.Key(int(variant))
				navigation.SplitButton(c, navigation.SplitButtonOptions{
					Label: "Save", MenuLabel: "Save options", Variant: variant,
					OnClick: func() { st.lastSplitAction = "Save" },
					Menu:    splitButtonDemoMenu(st),
				})
			}
		})
		if st.lastSplitAction != "" {
			formsMutedText(c, "Last action: "+st.lastSplitAction)
		}
	})
}

func splitButtonStatesBody(c *ui.Context, st *treeMenusState) {
	ui.Row(c).Gap(theme.Space12).Wrap().AlignItems(ui.Center).Children(func() {
		navigation.SplitButton(c, navigation.SplitButtonOptions{
			Label: "Save", MenuLabel: "Save options", Icon: lucide.Save(), Size: theme.ButtonSmall,
			Variant: mygoui.ButtonOutline, ActionDisabled: true, Menu: splitButtonDemoMenu(st),
		})
		navigation.SplitButton(c, navigation.SplitButtonOptions{
			Label: "Publish", MenuLabel: "Publish options", Disabled: true, Menu: splitButtonDemoMenu(st),
		})
	})
}

func splitButtonRememberBody(c *ui.Context, st *treeMenusState) {
	if st.splitChoice == "" {
		st.splitChoice = "Save"
	}
	choices := []string{"Save", "Save As…", "Save a Copy…"}
	items := make([]navigation.MenuItemOptions, len(choices))
	for i, label := range choices {
		items[i] = navigation.MenuItemOptions{Key: label, Label: label}
	}
	navigation.SplitButton(c, navigation.SplitButtonOptions{
		Label: st.splitChoice, MenuLabel: "Save options", Variant: mygoui.ButtonOutline,
		Selected: &st.splitChoice, OnClick: func() { st.lastSplitAction = st.splitChoice },
		Menu: navigation.MenuOptions{Sections: []navigation.MenuSectionOptions{{Items: items}}},
	})
}

func splitButtonLoadingBody(c *ui.Context, st *treeMenusState) {
	ui.Row(c).Gap(theme.Space12).Wrap().AlignItems(ui.Center).Children(func() {
		navigation.SplitButton(c, navigation.SplitButtonOptions{
			Label: "Save", MenuLabel: "Save options", Icon: lucide.Save(),
			Loading: st.splitLoading, Menu: splitButtonDemoMenu(st),
		})
		form.Switch(c, form.SwitchOptions{On: &st.splitLoading, Label: "Loading"})
	})
}

func overflowDemoItems() []navigation.OverflowItem {
	return []navigation.OverflowItem{
		{Label: "Overview", Icon: lucide.House()},
		{Label: "Activity"},
		{Label: "Deployments"},
		{Label: "Analytics"},
		{Label: "Members", Disabled: true},
		{Label: "Billing"},
		{Label: "Settings", Icon: lucide.Settings()},
	}
}

func overflowDemo(c *ui.Context, st *treeMenusState, collapse navigation.OverflowCollapse, label string) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		for _, width := range []float32{560, 300} {
			c.Key(int(width))
			ui.Box(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, width)).
				Padding(theme.Space4).Border(theme.Hairline, t.Border).Radius(theme.RadiusScale().XL).
				Children(func() {
					navigation.OverflowList(c, navigation.OverflowListOptions{
						Label: label, Items: overflowDemoItems(), Collapse: collapse,
						OnSelect: func(key string) { st.lastOverflow = key },
					})
				})
		}
		if st.lastOverflow != "" {
			formsMutedText(c, "Opened: "+st.lastOverflow)
		}
	})
}

func overflowToolbarBody(c *ui.Context, st *treeMenusState) {
	t := theme.Colors(c)
	items := []navigation.OverflowItem{
		{Key: "bold", Label: "Bold", Icon: lucide.Bold(), Checked: &st.overflowBold, Priority: 3},
		{Key: "italic", Label: "Italic", Icon: lucide.Italic(), Checked: &st.overflowItalic, Priority: 3},
		{Separator: true},
		{Key: "link", Label: "Link", Icon: lucide.Link(), Priority: 2},
		{Key: "quote", Label: "Quote", Icon: lucide.Quote(), Priority: 1},
		{Key: "code", Label: "Code", Icon: lucide.Code(), Priority: 1},
		{Separator: true},
		{Key: "help", Label: "Help", Icon: lucide.CircleHelp(), Pinned: true},
	}
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
		for _, width := range []float32{560, 260} {
			c.Key(int(width))
			ui.Box(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, width)).
				Padding(theme.Space4).Border(theme.Hairline, t.Border).Radius(theme.RadiusScale().XL).
				Children(func() {
					navigation.OverflowList(c, navigation.OverflowListOptions{Label: "Formatting", Items: items})
				})
		}
	})
}

func overflowEndBody(c *ui.Context, st *treeMenusState) {
	overflowDemo(c, st, navigation.OverflowCollapseEnd, "Project sections")
}

func overflowStartBody(c *ui.Context, st *treeMenusState) {
	overflowDemo(c, st, navigation.OverflowCollapseStart, "Project trail")
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "tree", Title: "Tree", Group: "Collections", Order: 30,
		Description: "An expandable hierarchy with selection, tri-state checks and virtualized rows.",
		Examples: []ComponentExample{
			{Title: "File tree", Description: "Click, Cmd/Ctrl-click and Shift-click choose; right-click opens a menu; double-click opens a folder.",
				Demo: treeMenusDemo(treeFilesBody)},
			{Title: "Guides", Demo: treeMenusDemo(treeGuidesBody)},
			{Title: "Lazy children", Description: "Opening the folder asks for its files; the arrow waits as a spinner.",
				Demo: treeMenusDemo(treeLazyBody)},
			{Title: "Checkable", Description: "A branch follows its leaves and shows a mixed check.",
				Demo: treeMenusDemo(treeCheckableBody)},
			{Title: "Ten thousand nodes", Demo: treeMenusDemo(treeLargeBody)},
			{Title: "Loading", Description: "Placeholder rows at three depths; with nodes, a dimmed refresh.", Demo: treeLoadingDemo},
		},
		Code: `package main

import "github.com/Cyenoch/mygoui/list"
import "github.com/egoist/mygo/ui"

var files list.TreeState

func tree(c *ui.Context, nodes []list.TreeNode) {
    list.Tree(c, list.TreeOptions{
        Label:     "Project files",
        Nodes:     nodes,
        State:     &files,
        Selection: list.SelectionSingle,
        Height:    240,
    })
}`,
		Params: []ComponentParam{
			{"Nodes", "[]list.TreeNode", "nil", "Top-level nodes: Key, Label, Icon, Disabled, Children; Lazy marks a branch whose children load on open, Loading shows a spinner for it."},
			{"State", "*list.TreeState", "required", "Expansion, selection, checks and scroll; bump Revision after editing nodes in place."},
			{"Selection", "list.SelectionMode", "SelectionNone", "None, single or multiple; multiple follows desktop rules: Cmd/Ctrl toggles, Shift ranges, Cmd/Ctrl+A all."},
			{"Checkable", "bool", "false", "Adds check boxes; branches follow their leaves."},
			{"Guides", "bool", "false", "Hairlines at each ancestor level."},
			{"Label", "string", "\"\"", "Accessible name of the tree."},
			{"Height", "float32", "0", "Fixed viewport; zero sizes to the rows."},
			{"MaxHeight", "float32", "0", "Caps a content-sized viewport."},
			{"Loading", "bool", "false", "Placeholder rows without nodes; a dimmed refresh line with nodes."},
			{"LoadingRows", "int", "5", "Number of placeholder rows."},
			{"OnAction", "func(key string)", "nil", "A row activated by a click or Enter."},
			{"OnSelectionChange", "func(keys []string)", "nil", "Chosen keys after a user change."},
			{"OnCheckedChange", "func(keys []string)", "nil", "Checked leaf keys after a user change."},
			{"OnExpandedChange", "func(key string, expanded bool)", "nil", "A branch the user opened or closed; where a Lazy node loads."},
			{"OnContextMenu", "func(key string, m *ui.Menu)", "nil", "Fills the menu of a row; an unchosen row is chosen first."},
		},
		Tokens: []string{"Foreground", "MutedForeground", "Accent", "Primary", "PrimaryForeground", "Input", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "menubar", Title: "Menubar", Group: "Patterns", Order: 150,
		Description: "A horizontal bar of application menus that open, switch and close from the keyboard or pointer.",
		Examples: []ComponentExample{
			{Title: "Application menus", Description: "Left and Right switch menus; hovering another title follows an open menu. F10 or Alt focuses the bar and Alt+letter opens a menu; typing in an open menu jumps to an item; the accelerators work while the bar is shown.",
				Demo: treeMenusDemo(menubarBody)},
		},
		Code: `navigation.Menubar(c, navigation.MenubarOptions{
    Label: "Application",
    Menus: []navigation.MenubarMenuOptions{
        {Label: "File", Menu: navigation.MenuOptions{
            Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
                {Label: "Save", Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyS}, OnSelect: save},
            }}},
        }},
        {Label: "Edit", Mnemonic: 'E', Menu: navigation.MenuOptions{OnSelect: remove}},
    },
})`,
		Params: []ComponentParam{
			{"Menus", "[]navigation.MenubarMenuOptions", "nil", "Titles in order: Key, Label, Disabled, Mnemonic (underlined while Alt is held, not on macOS) and the Menu each opens."},
			{"Accelerator", "form.Hotkey", "{}", "MenuItemOptions: a key chord that runs the item while the bar is built; shown in the platform's glyphs."},
			{"Open", "*string", "nil", "Key of the menu showing; nil keeps it internally."},
			{"Label", "string", "\"\"", "Accessible name of the bar."},
		},
		Tokens: []string{"Background", "Border", "Foreground", "Accent", "Popover", "PopoverForeground", "MutedForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "splitbutton", Title: "SplitButton", Group: "Actions", Order: 70,
		Description: "A primary action joined to an arrow that opens a menu of alternatives.",
		Examples: []ComponentExample{
			{Title: "Variants", Demo: treeMenusDemo(splitButtonVariantsBody)},
			{Title: "Remembered choice", Description: "The alternative chosen last becomes the action.",
				Demo: treeMenusDemo(splitButtonRememberBody)},
			{Title: "Loading", Description: "The action shows a spinner and takes no press; the menu stays usable.",
				Demo: treeMenusDemo(splitButtonLoadingBody)},
			{Title: "Disabled", Description: "The action alone, or the whole control.",
				Demo: treeMenusDemo(splitButtonStatesBody)},
		},
		Code: `navigation.SplitButton(c, navigation.SplitButtonOptions{
    Label:     "Save",
    MenuLabel: "Save options",
    OnClick:   save,
    Menu: navigation.MenuOptions{
        Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
            {Key: "save-as", Label: "Save As…"},
            {Key: "save-copy", Label: "Save a Copy…"},
        }}},
        OnSelect: remove,
    },
})`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Text and accessible name of the action."},
			{"OnClick", "func()", "nil", "Runs the action."},
			{"Menu", "navigation.MenuOptions", "{}", "Alternatives the arrow opens."},
			{"MenuLabel", "string", "\"More\"", "Accessible name of the arrow."},
			{"Variant", "mygoui.ButtonVariant", "ButtonPrimary", "Face of both segments."},
			{"Size", "theme.ButtonSize", "theme.ButtonMedium", "Height and padding of both segments."},
			{"Icon", "*ui.SVG", "nil", "Leading icon of the action."},
			{"Loading", "bool", "false", "Spinner on the action, which takes no press; the menu stays usable."},
			{"Disabled", "bool", "false", "Stops both segments."},
			{"ActionDisabled", "bool", "false", "Stops the action and keeps the menu."},
			{"Open", "*bool", "nil", "Caller-owned menu state; nil keeps it internally."},
			{"Selected", "*string", "nil", "Caller-owned key of the alternative chosen last; the menu marks it and a choice writes it."},
			{"Placement", "mygoui.Placement", "bottom end", "Where the menu opens along the arrow."},
		},
		Tokens: []string{"Primary", "PrimaryForeground", "Secondary", "SecondaryForeground", "Input", "Popover", "Accent"},
	})

	RegisterComponent(ComponentPage{
		ID: "overflowlist", Title: "OverflowList", Group: "Navigation", Order: 160,
		Description: "A row of actions that shows what fits and collapses the rest into a menu.",
		Examples: []ComponentExample{
			{Title: "Collapse from the end", Demo: treeMenusDemo(overflowEndBody)},
			{Title: "Collapse from the start", Demo: treeMenusDemo(overflowStartBody)},
			{Title: "Checked, pinned and prioritised", Description: "Labels drop to icons first; Priority picks who goes into the menu first, Pinned never does.",
				Demo: treeMenusDemo(overflowToolbarBody)},
		},
		Code: `navigation.OverflowList(c, navigation.OverflowListOptions{
    Label: "Project sections",
    Items: []navigation.OverflowItem{
        {Key: "overview", Label: "Overview"},
        {Key: "activity", Label: "Activity"},
        {Key: "settings", Label: "Settings"},
    },
    OnSelect: remove,
})`,
		Params: []ComponentParam{
			{"Items", "[]navigation.OverflowItem", "nil", "Actions in order: Key, Label, Icon, Disabled and OnSelect."},
			{"Checked", "*bool", "nil", "OverflowItem: a toggle, pressed in the row and checked in the menu."},
			{"Pinned", "bool", "false", "OverflowItem: stays in the row."},
			{"Priority", "int", "0", "OverflowItem: the lowest go into the menu first."},
			{"Separator", "bool", "false", "OverflowItem: a rule between groups; it stays only between shown items."},
			{"Collapse", "navigation.OverflowCollapse", "OverflowCollapseEnd", "Which end moves into the menu first."},
			{"Label", "string", "\"\"", "Accessible name of the row."},
			{"OnSelect", "func(key string)", "nil", "An item activated in the row or in the menu."},
		},
		Tokens: []string{"Foreground", "Accent", "Popover", "PopoverForeground", "Border"},
	})
}
