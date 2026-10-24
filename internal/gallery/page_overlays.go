package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// overlaysState owns every open flag on the page. The toast manager is held
// here as a mutable value because ToastState retains identities and payloads;
// the page keeps a stable pointer to this state through StateOf.
type overlaysState struct {
	modalOpen         bool
	modalSize         overlay.ModalSize
	preventDismiss    bool
	modalNote         string
	modalPreventOpen  bool
	modalBackdropOpen bool
	modalBackdrop     overlay.ModalBackdrop

	alertOpen       bool
	alertResult     string
	alertInfoOpen   bool
	alertInfoResult string

	drawerOpen       bool
	drawerSide       overlay.DrawerSide
	drawerNote       string
	drawerDetailOpen bool

	sheetOpen      bool
	sheetTallOpen  bool
	dialogOpen     bool
	dialogEdgeOpen bool

	popoverOpen     bool
	popoverEdgeOpen bool

	toasts    overlay.ToastState
	toastTop  bool
	toastNote string
	updatePct int
}

func modalSizeName(size overlay.ModalSize) string {
	switch size {
	case overlay.ModalSmall:
		return "Small"
	case overlay.ModalLarge:
		return "Large"
	case overlay.ModalExtraLarge:
		return "Extra large"
	case overlay.ModalFull:
		return "Full screen"
	}
	return "Medium"
}

func drawerSideName(side overlay.DrawerSide) string {
	switch side {
	case overlay.DrawerTop:
		return "Top"
	case overlay.DrawerLeft:
		return "Left"
	case overlay.DrawerRight:
		return "Right"
	}
	return "Bottom"
}

func shortcutRow(c *ui.Context, label, keys string) {
	ui.Row(c).FillWidth().AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(theme.Space8).Children(func() {
		mygoui.Typography(c, label, mygoui.TypographyOptions{Type: theme.TypeBodySM})
		mygoui.Kbd(c, mygoui.KbdOptions{Keys: []string{keys}})
	})
}

func toastPlacement(top bool) overlay.ToastPlacement {
	if top {
		return overlay.ToastTop
	}
	return overlay.ToastBottom
}

// modalSizesExample opens the shared modal at every size.
func modalSizesExample(c *ui.Context, st *overlaysState) {
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			for _, size := range []overlay.ModalSize{
				overlay.ModalSmall, overlay.ModalMedium, overlay.ModalLarge,
				overlay.ModalExtraLarge, overlay.ModalFull,
			} {
				size := size
				label := modalSizeName(size)
				mygoui.Button(c.Key(label), mygoui.ButtonOptions{
					Label:   "Open " + label + " modal",
					Variant: mygoui.ButtonSecondary,
					Size:    theme.ButtonSmall,
					OnClick: func() {
						st.modalSize = size
						st.modalOpen = true
					},
				})
			}
		})
		note := st.modalNote
		if note == "" {
			note = "open a size to inspect the panel and footer"
		}
		mygoui.Chip(c, mygoui.ChipOptions{Label: note, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft})
	})

	overlay.Modal(c, &st.modalOpen, overlay.ModalOptions{
		Size:        st.modalSize,
		Title:       modalSizeName(st.modalSize) + " modal",
		Description: "A scrollable panel with a two-action footer.",
		OnDismiss:   func() { st.modalNote = "modal dismissed" },
		Body: func() {
			mygoui.Paragraph(c, "The body scrolls when the window is short and focus stays inside the panel.", true)
			ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
				mygoui.Chip(c, mygoui.ChipOptions{Label: "Scrollable body", Variant: mygoui.ChipSoft})
				mygoui.Chip(c, mygoui.ChipOptions{Label: modalSizeName(st.modalSize), Variant: mygoui.ChipSoft})
			})
		},
		Footer: func() {
			ui.Row(c).Justify(ui.End).Gap(theme.Space8).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Cancel", Variant: mygoui.ButtonTertiary,
					OnClick: func() {
						st.modalOpen = false
						st.modalNote = "modal cancelled"
					},
				})
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Save changes", Variant: mygoui.ButtonPrimary,
					OnClick: func() {
						st.modalOpen = false
						st.modalNote = "changes saved from the " + modalSizeName(st.modalSize) + " modal"
					},
				})
			})
		},
	})
}

// modalDismissExample keeps the backdrop and Escape inert.
func modalDismissExample(c *ui.Context, st *overlaysState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Open protected modal", Variant: mygoui.ButtonSecondary,
			OnClick: func() {
				st.modalPreventOpen = true
				st.modalNote = "backdrop and Escape are now ignored"
			},
		})
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label:    "Prevent dismiss",
			Variant:  mygoui.ToggleGhost,
			Selected: &st.preventDismiss,
		})
	})
	overlay.Modal(c, &st.modalPreventOpen, overlay.ModalOptions{
		Title:          "Unsaved changes",
		Description:    "Backdrop presses and Escape are ignored while the switch is on.",
		PreventDismiss: st.preventDismiss,
		Body: func() {
			mygoui.Paragraph(c, "Close the panel with its own button, or turn the switch off first.", true)
		},
		Footer: func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Keep editing", Variant: mygoui.ButtonPrimary,
				OnClick: func() { st.modalPreventOpen = false },
			})
		},
	})
}

// modalBackdropExample switches between the blur and transparent backdrops.
func modalBackdropExample(c *ui.Context, st *overlaysState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, variant := range []overlay.ModalBackdrop{
			overlay.ModalBackdropBlur, overlay.ModalBackdropTransparent,
		} {
			variant := variant
			label := "Blur"
			if variant == overlay.ModalBackdropTransparent {
				label = "Transparent"
			}
			mygoui.Button(c.Key(label), mygoui.ButtonOptions{
				Label:   "Open " + label + " backdrop",
				Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					st.modalBackdrop = variant
					st.modalBackdropOpen = true
				},
			})
		}
	})
	overlay.Modal(c, &st.modalBackdropOpen, overlay.ModalOptions{
		Backdrop:    st.modalBackdrop,
		Title:       "Backdrop variant",
		Description: "The backdrop changes what is visible behind the panel.",
		Body: func() {
			mygoui.Paragraph(c, "A blur dims and blurs the page; a transparent backdrop leaves it visible.", true)
		},
	})
}

// alertDestructiveExample asks for an explicit destructive confirmation.
func alertDestructiveExample(c *ui.Context, st *overlaysState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label:   "Delete project",
				Icon:    lucide.Trash(),
				Variant: mygoui.ButtonDangerSoft,
				OnClick: func() { st.alertOpen = true },
			})
			if st.alertResult != "" {
				mygoui.Chip(c, mygoui.ChipOptions{
					Label: st.alertResult, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft,
				})
			}
		})
	})
	overlay.AlertDialog(c, &st.alertOpen, overlay.AlertDialogOptions{
		Status:      mygoui.ColorDanger,
		Title:       "Delete project?",
		Description: "This removes every deployment and cannot be undone.",
		CancelLabel: "Keep project",
		ActionLabel: "Delete permanently",
		Destructive: true,
		OnCancel:    func() { st.alertResult = "Project kept" },
		OnAction:    func() { st.alertResult = "Project deleted" },
		Body: func() {
			mygoui.Paragraph(c, "Members lose access immediately. Backups are retained for 30 days.", true)
		},
	})
}

// alertInfoExample confirms a non-destructive action.
func alertInfoExample(c *ui.Context, st *overlaysState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label:   "Publish release",
				Variant: mygoui.ButtonSecondary,
				OnClick: func() { st.alertInfoOpen = true },
			})
			if st.alertInfoResult != "" {
				mygoui.Chip(c, mygoui.ChipOptions{
					Label: st.alertInfoResult, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft,
				})
			}
		})
	})
	overlay.AlertDialog(c, &st.alertInfoOpen, overlay.AlertDialogOptions{
		Status:      mygoui.ColorAccent,
		Title:       "Publish release?",
		Description: "The build becomes available to every workspace.",
		CancelLabel: "Not yet",
		ActionLabel: "Publish now",
		OnCancel:    func() { st.alertInfoResult = "Release left unpublished" },
		OnAction:    func() { st.alertInfoResult = "Release published" },
	})
}

// drawerEdgesExample opens the same panel from each window edge.
func drawerEdgesExample(c *ui.Context, st *overlaysState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			for _, side := range []overlay.DrawerSide{
				overlay.DrawerLeft, overlay.DrawerRight, overlay.DrawerTop, overlay.DrawerBottom,
			} {
				side := side
				mygoui.Button(c.Key(drawerSideName(side)), mygoui.ButtonOptions{
					Label:   "Open " + drawerSideName(side) + " drawer",
					Variant: mygoui.ButtonSecondary,
					Size:    theme.ButtonSmall,
					OnClick: func() {
						st.drawerSide = side
						st.drawerOpen = true
					},
				})
			}
		})
		note := st.drawerNote
		if note == "" {
			note = "choose an edge to open the drawer"
		}
		mygoui.Chip(c, mygoui.ChipOptions{Label: note, Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft})
	})
	side := st.drawerSide
	overlay.Drawer(c, &st.drawerOpen, overlay.DrawerOptions{
		Side:        side,
		Title:       drawerSideName(side) + " drawer",
		Description: "Edge panel with a focus trap and a caller-owned open flag.",
		Width:       theme.Scale(c, 360),
		Height:      theme.Scale(c, 320),
		Handle:      side == overlay.DrawerTop || side == overlay.DrawerBottom,
		OnDismiss:   func() { st.drawerNote = drawerSideName(side) + " drawer dismissed" },
		Body: func() {
			mygoui.Paragraph(c, "The drawer slides from "+drawerSideName(side)+". Backdrop presses and Escape dismiss it.", true)
		},
	})
}

// drawerDetailExample adds a handle and a footer action.
func drawerDetailExample(c *ui.Context, st *overlaysState) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: "Open filters drawer", Variant: mygoui.ButtonSecondary,
		OnClick: func() { st.drawerDetailOpen = true },
	})
	overlay.Drawer(c, &st.drawerDetailOpen, overlay.DrawerOptions{
		Side:        overlay.DrawerBottom,
		Title:       "Filters",
		Description: "A bottom drawer with the visual drag handle.",
		Height:      theme.Scale(c, 280),
		Handle:      true,
		Body: func() {
			mygoui.Paragraph(c, "The handle is visual only; the backdrop, Escape and the footer button dismiss the drawer.", true)
		},
		Footer: func() {
			ui.Row(c).Justify(ui.End).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Apply filters", Variant: mygoui.ButtonPrimary,
					OnClick: func() { st.drawerDetailOpen = false },
				})
			})
		},
	})
}

// sheetFooterExample pairs the sheet with a footer action.
func sheetFooterExample(c *ui.Context, st *overlaysState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label:   "Open bottom sheet",
			Variant: mygoui.ButtonSecondary,
			OnClick: func() { st.sheetOpen = true },
		})
		mygoui.Typography(c, "Dismiss by pressing the backdrop or Escape.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
	overlay.BottomSheet(c, &st.sheetOpen, overlay.BottomSheetOptions{
		Title:       "Quick actions",
		Description: "Native HeroUI's bottom sheet surface.",
		Body: func() {
			mygoui.Paragraph(c, "The sheet caps at 85% of the window height and keeps the drag handle purely visual.", true)
			ui.Row(c).Wrap().Gap(theme.Space8).Children(func() {
				mygoui.Chip(c, mygoui.ChipOptions{Label: "Duplicate", Variant: mygoui.ChipSoft})
				mygoui.Chip(c, mygoui.ChipOptions{
					Label: "Delete", Color: mygoui.ColorDanger, Variant: mygoui.ChipSoft,
				})
			})
		},
		Footer: func() {
			ui.Row(c).Justify(ui.End).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Close sheet", Variant: mygoui.ButtonPrimary,
					OnClick: func() { st.sheetOpen = false },
				})
			})
		},
	})
}

// sheetTallExample fixes the panel height and hides the handle.
func sheetTallExample(c *ui.Context, st *overlaysState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Open tall sheet", Variant: mygoui.ButtonSecondary,
			OnClick: func() { st.sheetTallOpen = true },
		})
		mygoui.Typography(c, "A fixed height with no drag handle.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
	overlay.BottomSheet(c, &st.sheetTallOpen, overlay.BottomSheetOptions{
		Title:      "Fixed-height sheet",
		Height:     theme.Scale(c, 320),
		HideHandle: true,
		Body: func() {
			mygoui.Paragraph(c, "The panel keeps its height regardless of its content.", true)
		},
	})
}

// dialogFooterExample is the smaller centered surface with actions.
func dialogFooterExample(c *ui.Context, st *overlaysState) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label:   "Open dialog",
		Variant: mygoui.ButtonSecondary,
		OnClick: func() { st.dialogOpen = true },
	})
	overlay.Dialog(c, &st.dialogOpen, overlay.DialogOptions{
		Title:       "Rename workspace",
		Description: "The new name appears in navigation and invitations.",
		Width:       theme.Scale(c, 420),
		Body: func() {
			mygoui.Paragraph(c, "Dialogs are the smaller centered surface; use a drawer or the bottom sheet when the content needs an edge.", true)
			mygoui.Chip(c, mygoui.ChipOptions{Label: "20 DIP padding", Variant: mygoui.ChipSoft})
		},
		Footer: func() {
			ui.Row(c).Justify(ui.End).Gap(theme.Space8).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Close dialog", Variant: mygoui.ButtonTertiary,
					OnClick: func() { st.dialogOpen = false },
				})
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Rename", Variant: mygoui.ButtonPrimary,
					OnClick: func() { st.dialogOpen = false },
				})
			})
		},
	})
}

// dialogEdgeExample spans the bottom edge instead of centering.
func dialogEdgeExample(c *ui.Context, st *overlaysState) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label:   "Open edge dialog",
		Variant: mygoui.ButtonSecondary,
		OnClick: func() { st.dialogEdgeOpen = true },
	})
	overlay.Dialog(c, &st.dialogEdgeOpen, overlay.DialogOptions{
		Placement:   mygoui.Placement{Side: mygoui.SideBottom},
		Title:       "Edge dialog",
		Description: "An edge side spans the window edge like a sheet.",
		Body: func() {
			mygoui.Paragraph(c, "SideAuto centers the panel; a concrete side changes its placement.", true)
		},
	})
}

// popoverArrowExample anchors an arrowed panel below its trigger.
func popoverArrowExample(c *ui.Context, st *overlaysState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		anchor := mygoui.Button(c, mygoui.ButtonOptions{
			Label:   "Show keyboard shortcuts",
			Icon:    lucide.Info(),
			Variant: mygoui.ButtonSecondary,
			OnClick: func() { st.popoverOpen = !st.popoverOpen },
		})
		overlay.Popover(c, anchor, &st.popoverOpen, overlay.PopoverOptions{
			Title:       "Keyboard shortcuts",
			Description: "A compact anchored panel.",
			Arrow:       true,
			Width:       theme.Scale(c, 320),
			Placement:   mygoui.Placement{Side: mygoui.SideBottom, Align: mygoui.AlignStart},
			Body: func() {
				ui.Column(c).Gap(theme.Space6).Children(func() {
					shortcutRow(c, "Open command menu", "\u2318K")
					shortcutRow(c, "Toggle sidebar", "\u2318B")
					shortcutRow(c, "Close popover", "Esc")
				})
			},
			Footer: func() {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Got it", Variant: mygoui.ButtonPrimary, Size: theme.ButtonSmall,
					OnClick: func() { st.popoverOpen = false },
				})
			},
		})
	})
}

// popoverEdgeExample places its panel at the right edge without an arrow.
func popoverEdgeExample(c *ui.Context, st *overlaysState) {
	anchor := mygoui.Button(c, mygoui.ButtonOptions{
		Label:   "Show edge popover",
		Variant: mygoui.ButtonSecondary,
		OnClick: func() { st.popoverEdgeOpen = !st.popoverEdgeOpen },
	})
	overlay.Popover(c, anchor, &st.popoverEdgeOpen, overlay.PopoverOptions{
		Title:     "Edge placement",
		Arrow:     false,
		Width:     theme.Scale(c, 260),
		Placement: mygoui.Placement{Side: mygoui.SideRight, Align: mygoui.AlignEnd},
		Body: func() {
			mygoui.Paragraph(c, "The panel prefers the right edge and aligns its end with the trigger.", true)
		},
	})
}

// tooltipPlacementExample anchors a tip at each side.
func tooltipPlacementExample(c *ui.Context, _ *overlaysState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		for _, placement := range []struct {
			label string
			place mygoui.Placement
		}{
			{"Top tip", mygoui.Placement{Side: mygoui.SideTop, Align: mygoui.AlignCenter}},
			{"Right tip", mygoui.Placement{Side: mygoui.SideRight, Align: mygoui.AlignCenter}},
			{"Bottom tip", mygoui.Placement{Side: mygoui.SideBottom, Align: mygoui.AlignCenter}},
			{"Left tip", mygoui.Placement{Side: mygoui.SideLeft, Align: mygoui.AlignCenter}},
		} {
			placement := placement
			anchor := mygoui.Button(c, mygoui.ButtonOptions{
				Label: placement.label, Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
			})
			overlay.Tooltip(c, anchor, overlay.TooltipOptions{
				Text: placement.label, Placement: placement.place,
			})
		}
	})
}

// tooltipRichExample shows arrowed rich content.
func tooltipRichExample(c *ui.Context, _ *overlaysState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		anchor := mygoui.Button(c, mygoui.ButtonOptions{
			Label:   "Copy share link",
			Icon:    lucide.Copy(),
			Variant: mygoui.ButtonSecondary,
		})
		overlay.Tooltip(c, anchor, overlay.TooltipOptions{
			Arrow:     true,
			Placement: mygoui.Placement{Side: mygoui.SideTop, Align: mygoui.AlignCenter},
			Content: func() {
				ui.Column(c).Gap(theme.Space4).Children(func() {
					mygoui.Typography(c, "Copy share link", mygoui.TypographyOptions{Type: theme.TypeBodySM})
					shortcutRow(c, "Expires in 7 days", "\u2318C")
				})
			},
		})
	})
}

// toastVariantsExample fires each visual variant.
func toastVariantsExample(c *ui.Context, st *overlaysState) {
	st.toasts.View(c)
	show := func(o overlay.ToastOptions) {
		o.Timeout = -1
		o.Placement = toastPlacement(st.toastTop)
		overlay.Toast(c, &st.toasts, o)
	}
	ui.Column(c).Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Default toast", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					show(overlay.ToastOptions{
						ID: "demo-default", Title: "Draft saved",
						Description: "Your changes are available offline.",
					})
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Accent toast", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					show(overlay.ToastOptions{
						ID: "demo-accent", Title: "New version available",
						Description: "Reload to pick up mygoui 0.5.", Variant: overlay.ToastAccent,
					})
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Success toast", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					show(overlay.ToastOptions{
						ID: "demo-success", Title: "Invite sent",
						Description: "Ada will receive an email shortly.", Variant: overlay.ToastSuccess,
					})
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Warning toast", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					show(overlay.ToastOptions{
						ID: "demo-warning", Title: "Storage almost full",
						Description: "92% of the workspace quota is used.", Variant: overlay.ToastWarning,
					})
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Danger toast", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					show(overlay.ToastOptions{
						ID: "demo-danger", Title: "Deploy failed",
						Description: "The build exited with status 1.", Variant: overlay.ToastDanger,
					})
				},
			})
		})
		mygoui.Chip(c, mygoui.ChipOptions{
			Label: fmt.Sprintf("%d toast(s) retained", len(st.toasts.Items())),
			Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft,
		})
	})
}

// toastActionsExample adds an inline action and an in-place update.
func toastActionsExample(c *ui.Context, st *overlaysState) {
	show := func(o overlay.ToastOptions) {
		o.Timeout = -1
		o.Placement = toastPlacement(st.toastTop)
		overlay.Toast(c, &st.toasts, o)
	}
	ui.Column(c).Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Action toast", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					show(overlay.ToastOptions{
						ID: "demo-action", Title: "Invite sent",
						Description: "Undo the invitation within 30 seconds.",
						Action:      "Undo",
						OnAction: func() {
							st.toastNote = "undo pressed; invitation revoked"
							show(overlay.ToastOptions{
								ID: "demo-action-done", Title: "Invitation revoked",
								Description: "Ada no longer has pending access.",
								Variant:     overlay.ToastSuccess,
							})
						},
					})
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Update toast", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
				OnClick: func() {
					st.updatePct += 25
					if st.updatePct > 100 {
						st.updatePct = 25
					}
					show(overlay.ToastOptions{
						ID: "demo-update", Title: fmt.Sprintf("Uploading\u2026 %d%%", st.updatePct),
						Description: "The same ID updates this toast in place.",
						Variant:     overlay.ToastAccent,
					})
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Dismiss update", Variant: mygoui.ButtonTertiary, Size: theme.ButtonSmall,
				OnClick: func() {
					st.toasts.Dismiss(c, "demo-update")
					st.toastNote = "upload toast dismissed"
				},
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Dismiss all toasts", Variant: mygoui.ButtonTertiary, Size: theme.ButtonSmall,
				OnClick: func() {
					for _, item := range st.toasts.Items() {
						st.toasts.Dismiss(c, item.ID)
					}
					st.toastNote = "all toasts dismissed"
				},
			})
		})
		if st.toastNote != "" {
			mygoui.Typography(c, st.toastNote, mygoui.TypographyOptions{
				Type: theme.TypeBodyXS, Muted: true,
			})
		}
	})
}

// toastPlacementExample moves the whole region to the top edge.
func toastPlacementExample(c *ui.Context, st *overlaysState) {
	show := func(o overlay.ToastOptions) {
		o.Timeout = -1
		o.Placement = toastPlacement(st.toastTop)
		overlay.Toast(c, &st.toasts, o)
	}
	mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
		Label:    "Anchor toasts at top",
		Variant:  mygoui.ToggleGhost,
		Selected: &st.toastTop,
		OnChange: func(top bool) {
			show(overlay.ToastOptions{
				ID: "demo-placement", Title: "Toast region moved",
				Description: "New toasts anchor at the " + map[bool]string{true: "top", false: "bottom"}[top] + " edge.",
				Variant:     overlay.ToastDefault,
			})
		},
	})
}
