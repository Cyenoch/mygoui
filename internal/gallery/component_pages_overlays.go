package gallery

import (
	"github.com/egoist/mygo/ui"
)

// overlaysStateOf returns the shared overlay state. Every open flag stays in
// this state, so a component page opened directly starts closed.
func overlaysStateOf(s *State) *overlaysState {
	return StateOf[overlaysState](s, "overlays")
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "modal", Title: "Modal", Group: "Overlays", Order: 10,
		Description: "A focus-trapped modal across every size, with a scrollable body and footer actions.",
		Examples: []ComponentExample{
			{Title: "Sizes", Description: "Every panel size, from small to full screen, with a two-action footer.",
				Demo: func(c *ui.Context, s *State) { modalSizesExample(c, overlaysStateOf(s)) }},
			{Title: "Prevent dismiss", Description: "The backdrop and Escape stay inert while the switch is on.",
				Demo: func(c *ui.Context, s *State) { modalDismissExample(c, overlaysStateOf(s)) }},
			{Title: "Backdrop", Description: "A blurred or transparent backdrop changes what stays visible behind the panel.",
				Demo: func(c *ui.Context, s *State) { modalBackdropExample(c, overlaysStateOf(s)) }},
		},
		Code: `open := false
overlay.Modal(c, &open, overlay.ModalOptions{
    Size:        overlay.ModalMedium,
    Title:       "Edit profile",
    Description: "Changes apply immediately.",
    Body: func() {
        mygoui.Paragraph(c, "Focus stays inside the panel.", true)
    },
    Footer: func() {
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Save changes", Variant: mygoui.ButtonPrimary,
            OnClick: func() { open = false },
        })
    },
})`,
		Params: []ComponentParam{
			{"Size", "overlay.ModalSize", "ModalMedium", "Small through full screen; the body scrolls when short."},
			{"Backdrop", "overlay.ModalBackdrop", "ModalBackdropOpaque", "Opaque dim, blur or transparent backdrop."},
			{"Title", "string", "required", "Heading and accessible name of the dialog."},
			{"Body", "func()", "required", "Scrollable content of the panel."},
			{"Footer", "func()", "nil", "Action row below the body."},
			{"PreventDismiss", "bool", "false", "Ignores backdrop presses and Escape while set."},
			{"HideClose", "bool", "false", "Removes the top-right close button."},
			{"OnDismiss", "func()", "nil", "Runs after a backdrop, Escape or close-button dismissal."},
		},
		Tokens: []string{"Overlay", "OverlayForeground", "Muted", "Backdrop"},
	})

	RegisterComponent(ComponentPage{
		ID: "alertdialog", Title: "AlertDialog", Group: "Overlays", Order: 20,
		Description: "An explicit confirmation that ignores backdrop presses and focuses Cancel first.",
		Examples: []ComponentExample{
			{Title: "Destructive confirm", Description: "The cancel action receives focus first; the action only runs when chosen.",
				Demo: func(c *ui.Context, s *State) { alertDestructiveExample(c, overlaysStateOf(s)) }},
			{Title: "Informational", Description: "A non-destructive confirmation with its own status color.",
				Demo: func(c *ui.Context, s *State) { alertInfoExample(c, overlaysStateOf(s)) }},
		},
		Code: `open := false
overlay.AlertDialog(c, &open, overlay.AlertDialogOptions{
    Status:      mygoui.ColorDanger,
    Title:       "Delete project?",
    Description: "This removes every deployment and cannot be undone.",
    CancelLabel: "Keep project",
    ActionLabel: "Delete permanently",
    Destructive: true,
    OnAction:    deleteProject,
})`,
		Params: []ComponentParam{
			{"Status", "mygoui.StatusColor", "ColorDefault", "Icon fill and accent of the confirmation."},
			{"Title", "string", "required", "Question shown to the user."},
			{"CancelLabel", "string", "locale Close", "Label of the safe choice, which receives focus first."},
			{"ActionLabel", "string", "\"\"", "Label of the confirm action; empty with no OnAction omits it."},
			{"Destructive", "bool", "false", "Renders the confirm action as a danger button."},
			{"OnAction", "func()", "nil", "Runs when the confirm action is chosen."},
			{"HideIcon", "bool", "false", "Removes the status icon above the title."},
		},
		Tokens: []string{"Overlay", "OverlayForeground", "Muted", "Danger", "Backdrop"},
	})

	RegisterComponent(ComponentPage{
		ID: "drawer", Title: "Drawer", Group: "Overlays", Order: 30,
		Description: "A modal edge panel from any side, with an optional visual drag handle.",
		Examples: []ComponentExample{
			{Title: "Edge placement", Description: "Open the panel from the left, right, top or bottom edge.",
				Demo: func(c *ui.Context, s *State) { drawerEdgesExample(c, overlaysStateOf(s)) }},
			{Title: "Handle and footer", Description: "A bottom drawer with the visual drag handle and an action.",
				Demo: func(c *ui.Context, s *State) { drawerDetailExample(c, overlaysStateOf(s)) }},
		},
		Code: `open := false
overlay.Drawer(c, &open, overlay.DrawerOptions{
    Side:     overlay.DrawerRight,
    Title:    "Filters",
    Width:    theme.Scale(c, 360),
    Handle:   false,
    OnDismiss: func() { note = "drawer dismissed" },
    Body: func() {
        mygoui.Paragraph(c, "The panel slides in from the selected edge.", true)
    },
})`,
		Params: []ComponentParam{
			{"Side", "overlay.DrawerSide", "DrawerBottom", "Bottom, left, right or top edge of the window."},
			{"Width", "float32", "384", "Panel width for left and right drawers, capped at 85% of the window."},
			{"Height", "float32", "content", "Panel extent for top and bottom drawers, capped at 85%."},
			{"Handle", "bool", "false", "Shows the visual drag handle on top and bottom drawers."},
			{"PreventDismiss", "bool", "false", "Ignores backdrop presses and Escape while set."},
			{"OnDismiss", "func()", "nil", "Runs after the drawer is dismissed."},
		},
		Tokens: []string{"Overlay", "OverlayForeground", "Muted", "Separator", "Backdrop"},
	})

	RegisterComponent(ComponentPage{
		ID: "bottomsheet", Title: "BottomSheet", Group: "Overlays", Order: 40,
		Description: "HeroUI's bottom panel with a centered handle, overlay surface and footer.",
		Examples: []ComponentExample{
			{Title: "With footer", Description: "The handle, overlay surface and a footer action.",
				Demo: func(c *ui.Context, s *State) { sheetFooterExample(c, overlaysStateOf(s)) }},
			{Title: "Fixed height", Description: "A taller sheet with the visual handle removed.",
				Demo: func(c *ui.Context, s *State) { sheetTallExample(c, overlaysStateOf(s)) }},
		},
		Code: `open := false
overlay.BottomSheet(c, &open, overlay.BottomSheetOptions{
    Title: "Quick actions",
    Body: func() {
        mygoui.Paragraph(c, "The sheet caps at 85% of the window height.", true)
    },
    Footer: func() {
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Close sheet", Variant: mygoui.ButtonPrimary,
            OnClick: func() { open = false },
        })
    },
})`,
		Params: []ComponentParam{
			{"Title", "string", "required", "Heading of the sheet."},
			{"Height", "float32", "content", "Fixed panel height, capped at 85% of the window."},
			{"HideHandle", "bool", "false", "Removes the visual drag handle."},
			{"PreventDismiss", "bool", "false", "Ignores backdrop presses and Escape while set."},
			{"HideClose", "bool", "false", "Removes the top-right close button."},
			{"OnDismiss", "func()", "nil", "Runs after the sheet is dismissed."},
		},
		Tokens: []string{"Overlay", "OverlayForeground", "Muted", "Backdrop"},
	})

	RegisterComponent(ComponentPage{
		ID: "dialog", Title: "Dialog", Group: "Overlays", Order: 50,
		Description: "The shared centered dialog panel with a heading, description and explicit footer.",
		Examples: []ComponentExample{
			{Title: "With footer", Description: "The centered panel with a heading, description and actions.",
				Demo: func(c *ui.Context, s *State) { dialogFooterExample(c, overlaysStateOf(s)) }},
			{Title: "Edge placement", Description: "A concrete side spans the window edge instead of centering.",
				Demo: func(c *ui.Context, s *State) { dialogEdgeExample(c, overlaysStateOf(s)) }},
		},
		Code: `open := false
overlay.Dialog(c, &open, overlay.DialogOptions{
    Width: theme.Scale(c, 420),
    Title: "Rename workspace",
    Description: "The new name appears in navigation and invitations.",
    Body: func() {
        mygoui.Paragraph(c, "Dialogs are the smaller centered surface.", true)
    },
    Footer: func() {
        mygoui.Button(c, mygoui.ButtonOptions{
            Label: "Rename", Variant: mygoui.ButtonPrimary,
            OnClick: func() { open = false },
        })
    },
})`,
		Params: []ComponentParam{
			{"Placement", "mygoui.Placement", "SideAuto", "SideAuto centers the panel; an edge side spans that edge like a sheet."},
			{"Width", "float32", "content", "Fixed panel width for the centered form."},
			{"Title", "string", "required", "Heading of the dialog."},
			{"Body", "func()", "required", "Content between the header and the footer."},
			{"PreventDismiss", "bool", "false", "Ignores backdrop presses and Escape while set."},
			{"OnDismiss", "func()", "nil", "Runs after the dialog is dismissed."},
		},
		Tokens: []string{"Overlay", "OverlayForeground", "Backdrop"},
	})

	RegisterComponent(ComponentPage{
		ID: "popover", Title: "Popover", Group: "Overlays", Order: 60,
		Description: "A non-modal anchored panel with an optional arrow, closed by outside presses or Escape.",
		Examples: []ComponentExample{
			{Title: "With arrow", Description: "The panel is non-modal: a press that dismisses it still reaches the page.",
				Demo: func(c *ui.Context, s *State) { popoverArrowExample(c, overlaysStateOf(s)) }},
			{Title: "Edge placement", Description: "The panel prefers the right edge and aligns with its trigger's end.",
				Demo: func(c *ui.Context, s *State) { popoverEdgeExample(c, overlaysStateOf(s)) }},
		},
		Code: `open := false
anchor := mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Show keyboard shortcuts", Icon: lucide.Info(),
    OnClick: func() { open = !open },
})
overlay.Popover(c, anchor, &open, overlay.PopoverOptions{
    Title:     "Keyboard shortcuts",
    Arrow:     true,
    Width:     theme.Scale(c, 320),
    Placement: mygoui.Placement{Side: mygoui.SideBottom, Align: mygoui.AlignStart},
    Body: func() {
        shortcutRow(c, "Open command menu", "⌘K")
    },
})`,
		Params: []ComponentParam{
			{"Placement", "mygoui.Placement", "SideBottom", "Preferred side and alignment; SideAuto falls back to below the anchor."},
			{"Arrow", "bool", "false", "Draws the 12 DIP arrow at the anchor-facing edge."},
			{"Width", "float32", "content", "Fixed panel width; zero sizes it to its content."},
			{"MinWidth", "float32", "180", "Keeps small popovers usable."},
			{"Title", "string", "\"\"", "Heading and fallback accessible name of the popup."},
			{"Body", "func()", "required", "Content of the anchored panel."},
		},
		Tokens: []string{"Overlay", "OverlayForeground", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "tooltip", Title: "Tooltip", Group: "Overlays", Order: 70,
		Description: "A short label shown on pointer rest or keyboard focus of a real control.",
		Examples: []ComponentExample{
			{Title: "Placement", Description: "Tab to the button to show the tip; Escape dismisses it.",
				Demo: func(c *ui.Context, s *State) { tooltipPlacementExample(c, overlaysStateOf(s)) }},
			{Title: "Arrow and rich content", Description: "An arrowed tip builds a small content block instead of plain text.",
				Demo: func(c *ui.Context, s *State) { tooltipRichExample(c, overlaysStateOf(s)) }},
		},
		Code: `anchor := mygoui.Button(c, mygoui.ButtonOptions{
    Label: "Copy share link", Icon: lucide.Copy(),
    Variant: mygoui.ButtonSecondary,
})
overlay.Tooltip(c, anchor, overlay.TooltipOptions{
    Text:      "Copies a link that expires in 7 days",
    Arrow:     true,
    Placement: mygoui.Placement{Side: mygoui.SideTop, Align: mygoui.AlignCenter},
})`,
		Params: []ComponentParam{
			{"Text", "string", "required", "Plain label of the tip."},
			{"Content", "func()", "nil", "Rich content replacing the plain text."},
			{"Placement", "mygoui.Placement", "SideTop", "Preferred side; SideAuto falls back above the anchor."},
			{"Arrow", "bool", "false", "Draws HeroUI's overlay arrow."},
		},
		Tokens: []string{"Overlay", "OverlayForeground"},
	})

	RegisterComponent(ComponentPage{
		ID: "toast", Title: "Toast", Group: "Overlays", Order: 80,
		Description: "A caller-owned toast manager with variants, actions, in-place updates and dismissal.",
		Examples: []ComponentExample{
			{Title: "Variants", Description: "Default, accent, success, warning and danger surfaces.",
				Demo: func(c *ui.Context, s *State) { toastVariantsExample(c, overlaysStateOf(s)) }},
			{Title: "Actions and updates", Description: "An inline action, an in-place update and explicit dismissal.",
				Demo: func(c *ui.Context, s *State) { toastActionsExample(c, overlaysStateOf(s)) }},
			{Title: "Placement", Description: "Move the whole toast region between the top and bottom edges.",
				Demo: func(c *ui.Context, s *State) { toastPlacementExample(c, overlaysStateOf(s)) }},
		},
		Code: `var toasts overlay.ToastState
toasts.View(c) // install the region once per page

overlay.Toast(c, &toasts, overlay.ToastOptions{
    ID:          "invite-sent",
    Title:       "Invite sent",
    Description: "Ada will receive an email shortly.",
    Variant:     overlay.ToastSuccess,
    Placement:   overlay.ToastBottom,
    Timeout:     -1,
    Action:      "Undo",
    OnAction:    revokeInvite,
})`,
		Params: []ComponentParam{
			{"ID", "string", "generated", "Reusing a visible ID updates that toast in place and restarts its lifetime."},
			{"Variant", "overlay.ToastVariant", "ToastDefault", "Default, accent, success, warning or danger surface."},
			{"Placement", "overlay.ToastPlacement", "ToastBottom", "Moves the whole region to the top or bottom edge."},
			{"Timeout", "time.Duration", "4s", "Lifetime of the toast; -1 keeps it until dismissed."},
			{"Action", "string", "\"\"", "Label of an inline action button."},
			{"OnAction", "func()", "nil", "Runs when the inline action is pressed."},
			{"OnClose", "func()", "nil", "Runs when the toast is dismissed or expires."},
		},
		Tokens: []string{"Overlay", "OverlayForeground", "Surface", "Accent", "Success", "Danger"},
	})
}
