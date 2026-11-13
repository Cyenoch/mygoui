package gallery

import (
	"strconv"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// guidanceState is the caller-owned state behind the confirmation, coach
// mark, shortcut sheet, action bar and about dialog demos.
type guidanceState struct {
	deleteOpen, publishOpen bool
	deleted, published      int

	saveOpen  bool
	saveUntil time.Time
	saved     int

	tipOpen   bool
	tour      overlay.TourState
	anchors   [4]ui.Handle
	tourQuery string

	shortcuts overlay.ShortcutsDialogState

	picked  [4]bool
	message ui.Handle

	aboutOpen bool
}

func guidanceDemo(body func(c *ui.Context, st *guidanceState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) { body(c, StateOf[guidanceState](s, "guidance")) }
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "popconfirm", Title: "Popconfirm", Group: "Overlays", Order: 120,
		Description: "A confirmation bubble anchored to the action it guards.",
		Examples: []ComponentExample{
			{Title: "Destructive", Description: "Cancel takes the focus; Escape or an outside press cancels.",
				Demo: guidanceDemo(popconfirmDestructiveBody)},
			{Title: "Placement", Demo: guidanceDemo(popconfirmPlacementBody)},
			{Title: "Pending", Description: "Loading holds the bubble open while the work runs.",
				Demo: guidanceDemo(popconfirmLoadingBody)},
		},
		Code: `var open, deleting bool
anchor := mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Delete", Variant: mygoui.ButtonDestructiveOutline,
    OnClick: func() { open = true },
})
overlay.Popconfirm(c, anchor, &open, overlay.PopconfirmOptions{
    Title:        "Delete this project?",
    Description:  "This cannot be undone.",
    ConfirmLabel: "Delete",
    Destructive:  true,
    Loading:      deleting,
    // Start the work and keep the bubble open until it ends.
    OnConfirm: func() { deleting, open = true, true; deleteProject() },
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "The question, and the bubble's accessible name."},
			{"Description", "string", "\"\"", "Consequence of confirming."},
			{"ConfirmLabel", "string", "\"Confirm\"", "Confirm button text."},
			{"CancelLabel", "string", "\"Cancel\"", "Cancel button text."},
			{"Destructive", "bool", "false", "Warning glyph and danger confirm button; Cancel takes the initial focus, otherwise Confirm does."},
			{"Loading", "bool", "false", "Confirm shows pending, Cancel is off, and Escape and an outside press do not close the bubble."},
			{"OnConfirm", "func()", "nil", "Runs after the confirm button closed the bubble; set *open back to true to keep it open while Loading."},
			{"OnCancel", "func()", "nil", "Runs after Cancel, Escape or an outside press closed it."},
			{"Placement", "mygoui.Placement", "below", "Side, alignment and gap from the anchor."},
			{"Width", "float32", "280", "Bubble width in DIPs."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border", "MutedForeground", "Destructive"},
	})

	RegisterComponent(ComponentPage{
		ID: "teachingtip", Title: "TeachingTip", Group: "Overlays", Order: 130,
		Description: "A non-modal coach mark pointing at the control it explains.",
		Examples: []ComponentExample{
			{Title: "Coach mark", Description: "The page stays usable until the tip is closed.",
				Demo: guidanceDemo(teachingTipBody)},
		},
		Code: `var open bool
anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Filters"})
overlay.TeachingTip(c, anchor, &open, overlay.TeachingTipOptions{
    Title: "Saved filters",
    Body:  "Reuse a filter you applied before.",
    Footer: func() {
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Got it", Size: theme.ButtonSmall,
            OnClick: func() { open = false },
        })
    },
})`,
		Params: []ComponentParam{
			{"Title", "string", "\"\"", "Heading and accessible name."},
			{"Body", "string", "\"\"", "Explanation under the title."},
			{"Footer", "func()", "nil", "Action row under the body."},
			{"Placement", "mygoui.Placement", "below", "Side, alignment and gap from the anchor."},
			{"Width", "float32", "320", "Tip width in DIPs."},
			{"HideClose", "bool", "false", "Leaves out the close button."},
			{"OnDismiss", "func()", "nil", "Runs after the close button, or Escape inside the tip, closed it."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border", "MutedForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "tour", Title: "Tour", Group: "Patterns", Order: 140,
		Description: "A guided walk through marked controls, one spotlight at a time.",
		Examples: []ComponentExample{
			{Title: "Mail app", Demo: guidanceDemo(tourBody)},
		},
		Code: `var tour overlay.TourState
var inbox, search ui.Handle
mygoui.Button(c, mygoui.ButtonOptions{Label: "Inbox", Handle: &inbox})
form.SearchField(c, &query, form.SearchFieldOptions{}).Bind(&search)
mygoui.Button(c, mygoui.ButtonOptions{
    Label:   "Take the tour",
    OnClick: func() { tour = overlay.TourState{Open: true} },
})
// Build the tour after its targets.
overlay.Tour(c, &tour, []overlay.TourStep{
    {Anchor: &inbox, Title: "Inbox", Body: "Everything new lands here.",
        Placement: mygoui.Placement{Side: mygoui.SideRight}},
    {Anchor: &search, Title: "Search", Body: "Find any message by sender or subject."},
}, overlay.TourOptions{OnFinish: save})`,
		Params: []ComponentParam{
			{"TourState.Open", "bool", "false", "Caller-owned: whether the tour shows."},
			{"TourState.Step", "int", "0", "Caller-owned index of the current step."},
			{"TourStep.Anchor", "*ui.Handle", "nil", "The step's target, bound with Bind; without a laid-out target the panel is centred and nothing is spotlit."},
			{"TourStep.Title", "string", "\"\"", "Heading and the panel's accessible name."},
			{"TourStep.Body", "string", "\"\"", "Explanation under the title."},
			{"TourStep.Placement", "mygoui.Placement", "below", "Preferred side, alignment and gap; the panel takes another side where it does not fit."},
			{"Width", "float32", "320", "Panel width in DIPs, the same for every step."},
			{"OnFinish", "func()", "nil", "Runs after Done on the last step."},
			{"OnDismiss", "func()", "nil", "Runs after the close button or Escape."},
		},
		Tokens: []string{"Backdrop", "Popover", "PopoverForeground", "Border", "MutedForeground", "Primary"},
	})

	RegisterComponent(ComponentPage{
		ID: "shortcutsdialog", Title: "ShortcutsDialog", Group: "Patterns", Order: 150,
		Description: "A searchable sheet of keyboard shortcuts by group.",
		Examples: []ComponentExample{
			{Title: "Cheat sheet", Demo: guidanceDemo(shortcutsDialogBody)},
		},
		Code: `var shortcuts overlay.ShortcutsDialogState
overlay.ShortcutsDialog(c, &shortcuts, overlay.ShortcutsDialogOptions{
    Title: "Keyboard shortcuts",
    Groups: []overlay.ShortcutGroup{
        {Title: "General", Items: []overlay.Shortcut{
            {Label: "Open command menu", Hotkey: form.Hotkey{Mods: ui.Ctrl, Key: ui.KeyK}},
            {Label: "Go to inbox", Keys: []string{"G", "I"}},
        }},
    },
})`,
		Params: []ComponentParam{
			{"ShortcutsDialogState.Open", "bool", "false", "Caller-owned: whether the sheet shows."},
			{"ShortcutsDialogState.Query", "string", "\"\"", "Caller-owned search text; each opening clears it. Every word must be in a row's label or keys: \"ctrl k\", \"cmd\"."},
			{"Title", "string", "\"\"", "Heading and accessible name."},
			{"Groups", "[]overlay.ShortcutGroup", "nil", "Titled sections of rows, in two columns on windows 800 DIPs or wider."},
			{"Shortcut.Label", "string", "\"\"", "What the shortcut does; matches of the search are marked."},
			{"Shortcut.Hotkey", "form.Hotkey", "zero", "The chord, shown as the platform names its keys, so the sheet follows the application's bindings."},
			{"Shortcut.Keys", "[]string", "nil", "Explicit key caps in order; replaces Hotkey, for a sequence such as G then I."},
			{"OnDismiss", "func()", "nil", "Runs after Escape, the backdrop or the close button."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border", "MutedBackground", "MutedForeground", "Input", "Ring"},
	})

	RegisterComponent(ComponentPage{
		ID: "actionbar", Title: "ActionBar", Group: "Patterns", Order: 160,
		Description: "A toolbar floating at the bottom of its container while a selection lasts.",
		Examples: []ComponentExample{
			{Title: "Selection", Demo: guidanceDemo(actionBarBody)},
		},
		Code: `var selected []string
ui.Column(c).FillWidth().Height(320).Children(func() {
    // Built last, the bar floats over what the container holds.
    overlay.ActionBar(c, len(selected) > 0, overlay.ActionBarOptions{
        Count: len(selected), Label: "selected",
        Actions: []navigation.OverflowItem{
            {Label: "Archive", Icon: lucide.Archive(), OnSelect: save},
            {Label: "Delete", Icon: lucide.Trash2(), OnSelect: deleteProject},
        },
        OnClose: func() { selected = nil },
    })
})`,
		Params: []ComponentParam{
			{"open", "bool", "", "Second argument: shows the bar while true."},
			{"Label", "string", "\"\"", "Bar text after the Count, and with it the accessible name."},
			{"Count", "int", "0", "Rolls in front of Label when positive; a change is announced."},
			{"Actions", "[]navigation.OverflowItem", "nil", "The bar's controls; those that do not fit fold from the end into a menu."},
			{"OnClose", "func()", "nil", "Adds a close button that calls it; Escape in the bar does too."},
			{"ReturnFocus", "*ui.Handle", "nil", "Takes the focus when the bar closes while holding it."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border", "Separator"},
	})

	RegisterComponent(ComponentPage{
		ID: "aboutdialog", Title: "AboutDialog", Group: "Patterns", Order: 170,
		Description: "An application's identity, links, credits and licence in a dialog.",
		Examples: []ComponentExample{
			{Title: "Application", Demo: guidanceDemo(aboutDialogBody)},
		},
		Code: `var open bool
overlay.AboutDialog(c, &open, overlay.AboutDialogOptions{
    Icon:        func() { ui.Icon(c, lucide.NotebookPen()).Size(48, 48) },
    Name:        "Notes",
    Version:     "Version 2.4.1",
    CopyInfo:    "Notes 2.4.1 (318)\nmacOS 15.2",
    Description: "Plain text, everywhere.",
    Links:       []overlay.AboutLink{{Label: "Website", OnClick: openDocs}},
    Sections: []overlay.AboutSection{
        {Title: "Credits", Body: "Ada Lovelace\nGrace Hopper"},
        {Title: "Licence", Body: "Permission is hereby granted…", Collapsible: true},
    },
})`,
		Params: []ComponentParam{
			{"Icon", "func()", "nil", "Builds the icon or image above the name."},
			{"Name", "string", "\"\"", "Application name and the dialog's accessible name."},
			{"Version", "string", "\"\"", "Version line under the name."},
			{"CopyInfo", "string", "\"\"", "Text of the copy button beside the version, such as version, build and platform; no button without it."},
			{"Description", "string", "\"\"", "One or two sentences about the application."},
			{"Links", "[]overlay.AboutLink", "nil", "Label and OnClick of each link; opening a URL is the caller's work."},
			{"Sections", "[]overlay.AboutSection", "nil", "Titled text blocks such as credits or a licence."},
			{"AboutSection.Collapsible", "bool", "false", "Shows the title as a closed disclosure over a scrolling monospace body."},
			{"OnDismiss", "func()", "nil", "Runs after Escape, the backdrop or the close button."},
		},
		Tokens: []string{"Popover", "PopoverForeground", "Border", "MutedForeground", "Link"},
	})
}

func popconfirmDestructiveBody(c *ui.Context, st *guidanceState) {
	ui.Row(c).Gap(theme.Space16).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		anchor := mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Delete report", Icon: lucide.Trash2(), Variant: mygoui.ButtonDestructiveOutline,
			OnClick: func() { st.deleteOpen = !st.deleteOpen },
		})
		overlay.Popconfirm(c, anchor, &st.deleteOpen, overlay.PopconfirmOptions{
			Title:        "Delete this report?",
			Description:  "Shared links to it stop working.",
			ConfirmLabel: "Delete",
			Destructive:  true,
			OnConfirm:    func() { st.deleted++ },
		})
		formsMutedText(c, "Deleted "+strconv.Itoa(st.deleted)+" times")
	})
}

func popconfirmPlacementBody(c *ui.Context, st *guidanceState) {
	ui.Row(c).Gap(theme.Space16).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		formsMutedText(c, "Published "+strconv.Itoa(st.published)+" times")
		anchor := mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Publish", Variant: mygoui.ButtonOutline,
			OnClick: func() { st.publishOpen = !st.publishOpen },
		})
		overlay.Popconfirm(c, anchor, &st.publishOpen, overlay.PopconfirmOptions{
			Title:        "Publish to everyone?",
			ConfirmLabel: "Publish",
			Placement:    mygoui.Placement{Side: mygoui.SideRight},
			OnConfirm:    func() { st.published++ },
		})
	})
}

func popconfirmLoadingBody(c *ui.Context, st *guidanceState) {
	loading := st.saveOpen && c.Now().Before(st.saveUntil)
	if loading {
		c.After(st.saveUntil.Sub(c.Now()))
	} else if st.saveOpen && !st.saveUntil.IsZero() {
		st.saveOpen, st.saveUntil = false, time.Time{}
		st.saved++
	}
	ui.Row(c).Gap(theme.Space16).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		anchor := mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Sync now", Icon: lucide.RefreshCw(), Variant: mygoui.ButtonOutline,
			OnClick: func() { st.saveOpen = !st.saveOpen },
		})
		overlay.Popconfirm(c, anchor, &st.saveOpen, overlay.PopconfirmOptions{
			Title:        "Replace the remote copy?",
			ConfirmLabel: "Sync",
			Loading:      loading,
			OnConfirm: func() {
				st.saveOpen = true
				st.saveUntil = c.Now().Add(2 * time.Second)
			},
		})
		formsMutedText(c, "Synced "+strconv.Itoa(st.saved)+" times")
	})
}

func teachingTipBody(c *ui.Context, st *guidanceState) {
	ui.Row(c).Gap(theme.Space12).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
		anchor := mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Filters", Icon: lucide.ListFilter(), Variant: mygoui.ButtonOutline,
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Show tip", Variant: mygoui.ButtonGhost, Disabled: st.tipOpen,
			OnClick: func() { st.tipOpen = true },
		})
		overlay.TeachingTip(c, anchor, &st.tipOpen, overlay.TeachingTipOptions{
			Title: "Saved filters",
			Body:  "Reuse a filter you applied before instead of building it again.",
			Footer: func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Got it", Size: theme.ButtonSmall,
					OnClick: func() { st.tipOpen = false },
				})
			},
		})
	})
}

func tourBody(c *ui.Context, st *guidanceState) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space16).AlignItems(ui.Center).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Take the tour", Icon: lucide.Compass(), Variant: mygoui.ButtonOutline,
			OnClick: func() { st.tour = overlay.TourState{Open: true} },
		})
		ui.Row(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 560)).Height(theme.Scale(c, 260)).
			Clip().Radius(theme.RadiusScale().LG).Border(theme.Hairline, t.Border).
			Background(t.Background).Label("Mail").Children(func() {
			ui.Column(c).Width(theme.Scale(c, 148)).Shrink(0).Padding(theme.Space8).Gap(theme.Space2).
				BorderWidth(0, theme.Hairline, 0, 0).BorderColor(t.Border).Children(func() {
				for i, item := range [...]struct {
					label string
					icon  *ui.SVG
				}{{"Inbox", lucide.Inbox()}, {"Drafts", lucide.File()}, {"Archive", lucide.Archive()}} {
					o := mygoui.ButtonOptions{
						Label: item.label, Icon: item.icon, Variant: mygoui.ButtonGhost,
						Align: mygoui.ButtonAlignStart, FullWidth: true,
					}
					if i == 0 {
						o.Handle = &st.anchors[0]
					}
					mygoui.Button(c.Key(i), o)
				}
			})
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
				ui.Row(c).FillWidth().MinWidth(0).Padding(theme.Space8).Gap(theme.Space8).AlignItems(ui.Center).
					BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).Children(func() {
					form.SearchField(c, &st.tourQuery, form.SearchFieldOptions{AriaLabel: "Search mail", Placeholder: "Search", FullWidth: true}).
						Grow(1).Basis(0).MinWidth(0).Bind(&st.anchors[1])
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "New message", Icon: lucide.SquarePen(), IconOnly: true,
						Variant: mygoui.ButtonGhost, Tooltip: "New message", Handle: &st.anchors[2],
					})
				})
				ui.Column(c).FillWidth().Grow(1).Padding(theme.Space12).Gap(theme.Space8).Children(func() {
					for _, subject := range [...]string{"Quarterly report", "Design review", "Release notes"} {
						formsMutedText(c, subject)
					}
				})
				ui.Row(c).FillWidth().Padding(theme.Space8).Justify(ui.End).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Synced", Icon: lucide.CloudCheck(), Variant: mygoui.ButtonGhost,
						Size: theme.ButtonSmall, Handle: &st.anchors[3],
					})
				})
			})
		})
	})
	overlay.Tour(c, &st.tour, []overlay.TourStep{
		{Anchor: &st.anchors[0], Title: "Inbox", Body: "Everything new lands here.",
			Placement: mygoui.Placement{Side: mygoui.SideRight}},
		{Anchor: &st.anchors[1], Title: "Search", Body: "Find any message by sender or subject."},
		{Anchor: &st.anchors[2], Title: "New message", Body: "Start a draft from anywhere.",
			Placement: mygoui.Placement{Align: mygoui.AlignEnd}},
		{Anchor: &st.anchors[3], Title: "Sync", Body: "Shows when your mail is up to date.",
			Placement: mygoui.Placement{Side: mygoui.SideTop, Align: mygoui.AlignEnd}},
	}, overlay.TourOptions{})
}

func shortcutsDialogBody(c *ui.Context, st *guidanceState) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: "Keyboard shortcuts", Icon: lucide.Keyboard(), Variant: mygoui.ButtonOutline,
		OnClick: func() { st.shortcuts.Open = true },
	}).AlignSelf(ui.Center)
	chord := func(mods ui.Modifiers, key ui.Key) form.Hotkey { return form.Hotkey{Mods: mods, Key: key} }
	overlay.ShortcutsDialog(c, &st.shortcuts, overlay.ShortcutsDialogOptions{
		Title: "Keyboard shortcuts",
		Groups: []overlay.ShortcutGroup{
			{Title: "General", Items: []overlay.Shortcut{
				{Label: "Open command menu", Hotkey: chord(ui.Ctrl, ui.KeyK)},
				{Label: "Toggle sidebar", Hotkey: chord(ui.Ctrl, ui.KeyB)},
				{Label: "Show keyboard shortcuts", Hotkey: chord(ui.Ctrl, ui.KeySlash)},
				{Label: "Close dialog", Keys: []string{"Esc"}},
				{Label: "Go to inbox", Keys: []string{"G", "I"}},
			}},
			{Title: "Editing", Items: []overlay.Shortcut{
				{Label: "Undo", Hotkey: chord(ui.Ctrl, ui.KeyZ)},
				{Label: "Redo", Hotkey: chord(ui.Ctrl|ui.Shift, ui.KeyZ)},
				{Label: "Find", Hotkey: chord(ui.Ctrl, ui.KeyF)},
				{Label: "Find and replace", Hotkey: chord(ui.Ctrl|ui.Alt, ui.KeyF)},
				{Label: "Duplicate line", Hotkey: chord(ui.Ctrl|ui.Shift, ui.KeyD)},
			}},
			{Title: "Navigation", Items: []overlay.Shortcut{
				{Label: "Next tab", Hotkey: chord(ui.Ctrl, ui.KeyTab)},
				{Label: "Previous tab", Hotkey: chord(ui.Ctrl|ui.Shift, ui.KeyTab)},
				{Label: "Go to line", Hotkey: chord(ui.Ctrl, ui.KeyG)},
				{Label: "Back", Hotkey: chord(ui.Alt, ui.KeyLeft)},
				{Label: "Forward", Hotkey: chord(ui.Alt, ui.KeyRight)},
			}},
		},
	})
}

func actionBarBody(c *ui.Context, st *guidanceState) {
	t := theme.Colors(c)
	count := 0
	for _, picked := range st.picked {
		if picked {
			count++
		}
	}
	clear := func() { st.picked = [4]bool{} }
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 520)).Height(theme.Scale(c, 280)).
		Clip().Radius(theme.RadiusScale().LG).Border(theme.Hairline, t.Border).
		Background(t.Background).Label("Messages").Children(func() {
		ui.Column(c).Fill().Padding(theme.Space24).Gap(theme.Space12).Children(func() {
			for i, subject := range [...]string{"Quarterly report", "Design review", "Release notes", "Offsite plan"} {
				o := form.CheckboxOptions{Checked: &st.picked[i], Label: subject}
				if i == 0 {
					o.Handle = &st.message
				}
				form.Checkbox(c.Key(i), o)
			}
		})
		overlay.ActionBar(c, count > 0, overlay.ActionBarOptions{
			Label: "selected", Count: count,
			Actions: []navigation.OverflowItem{
				{Label: "Archive", Icon: lucide.Archive(), OnSelect: clear},
				{Label: "Delete", Icon: lucide.Trash2(), OnSelect: clear},
				{Label: "Move", Icon: lucide.FolderInput(), OnSelect: clear},
				{Label: "Share", Icon: lucide.Share2(), OnSelect: clear},
			},
			OnClose:     clear,
			ReturnFocus: &st.message,
		})
	})
}

func aboutDialogBody(c *ui.Context, st *guidanceState) {
	t := theme.Colors(c)
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: "About Notes", Variant: mygoui.ButtonOutline,
		OnClick: func() { st.aboutOpen = true },
	}).AlignSelf(ui.Center)
	overlay.AboutDialog(c, &st.aboutOpen, overlay.AboutDialogOptions{
		Icon: func() {
			ui.Box(c).Size(theme.Space64, theme.Space64).Center().Radius(theme.RadiusScale().XXL).
				Background(t.Primary).TextColor(t.PrimaryForeground).Label("Notes icon").Children(func() {
				ui.Icon(c, lucide.NotebookPen()).Size(theme.Space32, theme.Space32)
			})
		},
		Name:        "Notes",
		Version:     "Version 2.4.1 (build 318)",
		CopyInfo:    "Notes 2.4.1 (318)\nmacOS 15.2",
		Description: "Plain-text notes that stay on your devices.",
		Links: []overlay.AboutLink{
			{Label: "Website", OnClick: func() {}},
			{Label: "Report an issue", OnClick: func() {}},
		},
		Sections: []overlay.AboutSection{
			{Title: "Credits", Body: "Ada Lovelace\nGrace Hopper\nMargaret Hamilton"},
			{Title: "Licence", Collapsible: true, Body: "Permission is hereby granted, free of charge, to any person obtaining a copy " +
				"of this software and associated documentation files, to deal in the Software without restriction, " +
				"including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, " +
				"and/or sell copies of the Software.\n\nThe above copyright notice and this permission notice shall be " +
				"included in all copies or substantial portions of the Software.\n\nThe software is provided \"as is\", " +
				"without warranty of any kind, express or implied."},
		},
	})
}
