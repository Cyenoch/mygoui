package gallery

import (
	"strconv"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/egoist/mygo/ui"
)

type OverlayState struct {
	Open, Disabled, Checked bool
	Text, Result, ToastID   string
	LoadingID               string
	Radio                   overlay.MenuChoice
	Menu                    overlay.MenuState
	Bar                     overlay.MenubarState
	Hover                   overlay.HoverCardState
	Toasts                  overlay.ToastState
	ToastIDs                []string
	ToastSequence           int
}

func NewOverlayState() OverlayState {
	return OverlayState{Text: "Editable overlay sample", Radio: "compact"}
}

func PreviewOverlay(c *ui.Context, s *OverlayState, id ComponentID) {
	ui.Column(c).FillWidth().MinWidth(0).Justify(ui.Start).AlignItems(ui.Start).Gap(12).Children(func() {
		ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
			mygoui.Toggle(c, &s.Disabled, mygoui.ToggleOptions{Label: "Disable trigger"})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset overlay", Variant: mygoui.ButtonOutline,
				OnClick: func() { s.Open = false; s.Disabled = false; s.Text = "Editable overlay sample"; s.Result = "" }})
		})
		body := func() { ui.TextInput(c, &s.Text).Label("Overlay sample text").FillWidth() }
		open := func() { s.Open = true }
		result := func() { s.Result = "Sample action completed"; s.Open = false; c.AddToast(ui.Toast{Title: s.Result}) }
		ui.Column(c).FillWidth().MinWidth(0).MinHeight(220).Justify(ui.Center).AlignItems(ui.Center).Gap(12).Children(func() {
			switch id {
			case "dialog", "sheet", "drawer", "alert-dialog":
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Open sample overlay", Disabled: s.Disabled, OnClick: open})
				o := overlay.DialogOptions{Title: "Sample overlay", Description: "This edits a local fixture only.",
					Footer: func() { mygoui.Button(c, mygoui.ButtonOptions{Label: "Save sample", OnClick: result}) }}
				switch id {
				case "dialog":
					overlay.Dialog(c, &s.Open, o, body)
				case "sheet":
					overlay.Sheet(c, &s.Open, overlay.SheetOptions{DialogOptions: o, Side: overlay.SideRight}, body)
				case "drawer":
					ui.Text(c, "Native edge modal. Swipe and velocity snap behavior are not implemented.")
					overlay.Drawer(c, &s.Open, overlay.DrawerOptions{DialogOptions: o, Side: overlay.SideBottom}, body)
				case "alert-dialog":
					overlay.AlertDialog(c, &s.Open, overlay.AlertDialogOptions{Title: "Apply sample change?",
						Description: "Only the local gallery result changes.", ActionLabel: "Apply sample", CancelLabel: "Cancel",
						OnAction: result, OnCancel: func() { s.Result = "Sample cancelled" }}, nil)
				}
			case "popover":
				anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Open sample popover", Disabled: s.Disabled,
					OnClick: func() { s.Open = !s.Open }})
				overlay.Popover(c, anchor, &s.Open, overlay.PopoverOptions{Title: "Sample popover", Description: "Native editable anchored content"}, body)
			case "tooltip":
				anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Focus or hover for tooltip", Disabled: s.Disabled,
					OnClick: func() { s.Result = "Tooltip anchor activated"; c.AddToast(ui.Toast{Title: s.Result}) }})
				overlay.Tooltip(c, anchor, overlay.TooltipOptions{Text: "Native tooltip also appears on keyboard focus"})
			case "hover-card":
				anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Focus or hover for preview", Disabled: s.Disabled,
					Variant: mygoui.ButtonLink, OnClick: func() { s.Result = "Preview link activated"; c.AddToast(ui.Toast{Title: s.Result}) }})
				overlay.HoverCard(c, anchor, &s.Hover, overlay.HoverCardOptions{PopoverOptions: overlay.PopoverOptions{
					Title: "Sample preview", Description: "Focus remains available in this native content."}}, body)
			case "dropdown-menu":
				overlay.DropdownMenu(c, &s.Menu, overlay.MenuOptions{Label: "Sample menu", Disabled: s.Disabled}, sampleMenu(c, s))
			case "context-menu":
				anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Right click or Shift+F10", Disabled: s.Disabled,
					Variant: mygoui.ButtonOutline, OnClick: func() { s.Result = "Use the context menu to select an action"; c.AddToast(ui.Toast{Title: s.Result}) }})
				overlay.ContextMenu(anchor, sampleMenu(c, s))
			case "menubar":
				overlay.Menubar(c, &s.Bar, []overlay.BarMenu{
					{ID: "sample", Label: "Sample", Disabled: s.Disabled, Items: sampleMenu(c, s)},
					{ID: "view", Label: "View", Disabled: s.Disabled, Items: []overlay.MenuItem{{ID: "view-action", Label: "Show sample result", OnSelect: result}}},
					{ID: "disabled", Label: "Unavailable", Disabled: true},
				})
			case "toast":
				ui.TextInput(c, &s.Text).Label("Toast sample description").FillWidth()
				ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Show sample toast", Disabled: s.Disabled, OnClick: func() {
						s.ToastID = overlay.Toast(c, &s.Toasts, overlay.ToastOptions{Title: "Sample changed", Description: s.Text,
							Kind: overlay.ToastSuccess, Action: "Undo", OnAction: func() { s.Result = "Sample change undone" }})
					}})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Show error toast", Disabled: s.Disabled, Variant: mygoui.ButtonDestructive,
						OnClick: func() {
							s.ToastID = overlay.Toast(c, &s.Toasts, overlay.ToastOptions{Title: "Sample error", Kind: overlay.ToastError})
						}})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Dismiss sample toast", Disabled: s.ToastID == "", Variant: mygoui.ButtonOutline,
						OnClick: func() { overlay.DismissToast(c, s.ToastID); s.ToastID = "" }})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Insert three notifications", Variant: mygoui.ButtonOutline, Disabled: s.Disabled, OnClick: func() {
						for i := range 3 {
							s.ToastSequence++
							description := "Short local update."
							if i == 1 {
								description = "A longer local notification demonstrates independent card height, readable expanded stacking and stable reflow when its middle neighbor disappears."
							}
							id := overlay.Toast(c, &s.Toasts, overlay.ToastOptions{Title: "Notification " + strconv.Itoa(s.ToastSequence), Description: description, Timeout: -1, Action: "Review", OnAction: func() { c.AddToast(ui.Toast{Title: "Notification reviewed"}) }})
							s.ToastIDs = append(s.ToastIDs, id)
						}
					}})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Start loading notification", Variant: mygoui.ButtonOutline, Disabled: s.Disabled, OnClick: func() {
						s.LoadingID = overlay.Toast(c, &s.Toasts, overlay.ToastOptions{ID: "gallery-loading", Title: "Saving local sample", Description: s.Text, Kind: overlay.ToastLoading})
					}})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Finish loading notification", Variant: mygoui.ButtonOutline, Disabled: s.Disabled || s.LoadingID == "", OnClick: func() {
						overlay.Toast(c, &s.Toasts, overlay.ToastOptions{ID: s.LoadingID, Title: "Local sample saved", Description: "The same notification now reports success.", Kind: overlay.ToastSuccess})
					}})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Fail loading notification", Variant: mygoui.ButtonOutline, Disabled: s.Disabled || s.LoadingID == "", OnClick: func() {
						overlay.Toast(c, &s.Toasts, overlay.ToastOptions{ID: s.LoadingID, Title: "Local sample could not be saved", Description: "Retry the local sample operation.", Kind: overlay.ToastError})
					}})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Close middle notification", Variant: mygoui.ButtonOutline, Disabled: len(s.ToastIDs) < 3, OnClick: func() {
						overlay.DismissToast(c, s.ToastIDs[len(s.ToastIDs)-2])
					}})
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Dismiss all sample notifications", Variant: mygoui.ButtonOutline, Disabled: len(s.ToastIDs) == 0 && s.ToastID == "" && s.LoadingID == "", OnClick: func() {
						for _, id := range s.ToastIDs {
							overlay.DismissToast(c, id)
						}
						if s.ToastID != "" {
							overlay.DismissToast(c, s.ToastID)
						}
						if s.LoadingID != "" {
							overlay.DismissToast(c, s.LoadingID)
						}
						s.LoadingID = ""
						s.ToastIDs = nil
						s.ToastID = ""
					}})
				})
			}
		})
	})
}

func sampleMenu(c *ui.Context, s *OverlayState) []overlay.MenuItem {
	return []overlay.MenuItem{
		{ID: "action", Label: "Record sample action", OnSelect: func() { s.Result = "Menu action selected"; c.AddToast(ui.Toast{Title: s.Result}) }},
		{ID: "checked", Label: "Show sample details", Kind: overlay.MenuCheckbox, Checked: &s.Checked,
			OnSelect: func() { s.Result = "Sample details changed" }},
		{ID: "divider", Kind: overlay.MenuSeparator},
		{ID: "density", Label: "Density", Items: []overlay.MenuItem{
			{ID: "compact", Label: "Compact", Kind: overlay.MenuRadio, Radio: &s.Radio, Value: "compact"},
			{ID: "roomy", Label: "Roomy", Kind: overlay.MenuRadio, Radio: &s.Radio, Value: "roomy"},
		}},
		{ID: "unavailable", Label: "Unavailable action", Disabled: true},
	}
}
