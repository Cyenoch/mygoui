package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type buttonsState struct {
	Ready bool

	Clicks     int
	LastAction string

	DefaultOff bool
	DefaultOn  bool
	GhostOff   bool
	GhostOn    bool

	BoldOn     bool
	ItalicOn   bool
	BookmarkOn bool

	AlignLeft   bool
	AlignCenter bool
	AlignRight  bool

	OverflowOn bool
	FiltersOn  bool

	GridOn bool
	ListOn bool

	CloseClicks int
	Presses     int
}

func (s *buttonsState) prepare() {
	if s.Ready {
		return
	}
	s.Ready = true
	s.DefaultOn = true
	s.GhostOn = true
	s.BoldOn = true
	s.AlignCenter = true
	s.ListOn = true
	s.LastAction = "None"
}

func (s *buttonsState) hit(label string) {
	s.Clicks++
	s.LastAction = label
}

// buttonVariantsBody renders one button per HeroUI variant and the shared
// activation feedback.
func buttonVariantsBody(c *ui.Context, st *buttonsState) {
	ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			variants := []struct {
				label   string
				variant mygoui.ButtonVariant
			}{
				{"Primary", mygoui.ButtonPrimary},
				{"Secondary", mygoui.ButtonSecondary},
				{"Tertiary", mygoui.ButtonTertiary},
				{"Outline", mygoui.ButtonOutline},
				{"Ghost", mygoui.ButtonGhost},
				{"Danger", mygoui.ButtonDanger},
				{"Danger Soft", mygoui.ButtonDangerSoft},
			}
			for _, item := range variants {
				label, variant := item.label, item.variant
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: label, Variant: variant,
					OnClick: func() { st.hit(label) },
				})
			}
		})
		ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: fmt.Sprintf("Actions fired %d", st.Clicks),
				Color: mygoui.ColorAccent, Variant: mygoui.ChipSoft, Size: mygoui.ChipSmall,
			})
			mygoui.Typography(c, "Last action: "+st.LastAction, mygoui.TypographyOptions{
				Type: theme.TypeBodyXS, Muted: true,
			})
		})
	})
}

// buttonSizesBody renders the three button heights.
func buttonSizesBody(c *ui.Context, st *buttonsState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		for _, item := range []struct {
			label string
			size  theme.ButtonSize
		}{
			{"Small", theme.ButtonSmall},
			{"Medium", theme.ButtonMedium},
			{"Large", theme.ButtonLarge},
		} {
			label, size := item.label, item.size
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: label, Variant: mygoui.ButtonSecondary, Size: size,
				OnClick: func() { st.hit(label) },
			})
		}
	})
}

// buttonIconsBody renders leading, trailing and icon-only layouts.
func buttonIconsBody(c *ui.Context, st *buttonsState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Download", Icon: lucide.Download(),
			OnClick: func() { st.hit("Download") },
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Continue", Variant: mygoui.ButtonSecondary, EndIcon: lucide.ArrowRight(),
			OnClick: func() { st.hit("Continue") },
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Notifications", Icon: lucide.Bell(), IconOnly: true,
			Variant: mygoui.ButtonSecondary, OnClick: func() { st.hit("Notifications") },
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Settings", Icon: lucide.Settings(), IconOnly: true,
			Variant: mygoui.ButtonGhost, OnClick: func() { st.hit("Settings") },
		})
	})
}

// buttonLoadingBody renders the pending state: a leading spinner, the full
// readable face and no activation.
func buttonLoadingBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{Label: "Saving", Loading: true})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Publishing", Variant: mygoui.ButtonSecondary, Loading: true,
		})
	})
}

// buttonDisabledBody renders the dimmed, non-activating state.
func buttonDisabledBody(c *ui.Context) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Unavailable", Variant: mygoui.ButtonSecondary, Disabled: true,
		})
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Delete", Variant: mygoui.ButtonDangerSoft, Icon: lucide.Trash2(), Disabled: true,
		})
	})
}

// buttonGroupBody renders a horizontal joined group.
func buttonGroupBody(c *ui.Context, st *buttonsState) {
	mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Separators: true},
		mygoui.ButtonOptions{
			Label: "Day", Variant: mygoui.ButtonSecondary,
			OnClick: func() { st.hit("Group day") },
		},
		mygoui.ButtonOptions{
			Label: "Week", Variant: mygoui.ButtonSecondary,
			OnClick: func() { st.hit("Group week") },
		},
		mygoui.ButtonOptions{
			Label: "Month", Variant: mygoui.ButtonSecondary,
			OnClick: func() { st.hit("Group month") },
		},
	)
}

// buttonGroupVerticalBody renders a full-width vertical group.
func buttonGroupVerticalBody(c *ui.Context, st *buttonsState) {
	ui.Column(c).FillWidth().MinWidth(theme.Scale(c, 180)).Children(func() {
		mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{
			Orientation: mygoui.GroupVertical, Separators: true, FullWidth: true,
		},
			mygoui.ButtonOptions{
				Label: "Profile", Variant: mygoui.ButtonSecondary,
				OnClick: func() { st.hit("Group profile") },
			},
			mygoui.ButtonOptions{
				Label: "Billing", Variant: mygoui.ButtonSecondary,
				OnClick: func() { st.hit("Group billing") },
			},
			mygoui.ButtonOptions{
				Label: "Sessions", Variant: mygoui.ButtonSecondary,
				OnClick: func() { st.hit("Group sessions") },
			},
		)
	})
}

// buttonGroupOutlineBody renders outline segments whose shared borders
// collapse to a single hairline.
func buttonGroupOutlineBody(c *ui.Context, st *buttonsState) {
	mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{},
		mygoui.ButtonOptions{
			Label: "Copy", Variant: mygoui.ButtonOutline,
			OnClick: func() { st.hit("Group copy") },
		},
		mygoui.ButtonOptions{
			Label: "Move", Variant: mygoui.ButtonOutline,
			OnClick: func() { st.hit("Group move") },
		},
		mygoui.ButtonOptions{
			Label: "Delete", Variant: mygoui.ButtonOutline,
			OnClick: func() { st.hit("Group delete") },
		},
	)
}

// closeButtonBody renders the dismiss action inside the notice it closes.
func closeButtonBody(c *ui.Context, st *buttonsState) {
	ui.Row(c).FillWidth().AlignItems(ui.Center).Wrap().Gap(theme.Space12).Children(func() {
		mygoui.Chip(c, mygoui.ChipOptions{
			Label: "Dismissible", Color: mygoui.ColorWarning,
			Variant: mygoui.ChipSoft, Size: mygoui.ChipSmall,
		})
		mygoui.Typography(c, "Release notes are up to date.", mygoui.TypographyOptions{
			Type: theme.TypeBodySM,
		})
		mygoui.CloseButton(c, mygoui.CloseButtonOptions{
			Label: "Dismiss notice", OnClick: func() { st.CloseClicks++ },
		})
		mygoui.Typography(c, fmt.Sprintf("Dismissed %d times", st.CloseClicks), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// closeButtonDisabledBody renders the dimmed dismiss control.
func closeButtonDisabledBody(c *ui.Context) {
	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space12).Children(func() {
		mygoui.Typography(c, "Notice locked while publishing", mygoui.TypographyOptions{
			Type: theme.TypeBodySM,
		})
		mygoui.CloseButton(c, mygoui.CloseButtonOptions{Label: "Dismiss notice", Disabled: true})
	})
}

// toggleDefaultBody renders the default surface in both states.
func toggleDefaultBody(c *ui.Context, st *buttonsState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Default off", Selected: &st.DefaultOff, OnChange: func(bool) {},
		})
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Default on", Selected: &st.DefaultOn, OnChange: func(bool) {},
		})
	})
}

// toggleGhostBody renders the ghost surface in both states.
func toggleGhostBody(c *ui.Context, st *buttonsState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Ghost off", Variant: mygoui.ToggleGhost,
			Selected: &st.GhostOff, OnChange: func(bool) {},
		})
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Ghost on", Variant: mygoui.ToggleGhost,
			Selected: &st.GhostOn, OnChange: func(bool) {},
		})
	})
}

// toggleIconsBody renders icon and icon-only toggles.
func toggleIconsBody(c *ui.Context, st *buttonsState) {
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Bold", Icon: lucide.Bold(), Selected: &st.BoldOn,
			OnChange: func(bool) {},
		})
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Italic", Icon: lucide.Italic(), Variant: mygoui.ToggleGhost,
			Selected: &st.ItalicOn, OnChange: func(bool) {},
		})
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Bookmark", Icon: lucide.Bookmark(), IconOnly: true,
			Selected: &st.BookmarkOn, OnChange: func(bool) {},
		})
	})
}

// toggleDisabledBody renders disabled toggles in both states.
func toggleDisabledBody(c *ui.Context) {
	off, on := false, true
	ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Unavailable", Selected: &off, Disabled: true,
		})
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: "Locked on", Selected: &on, Disabled: true,
		})
	})
}

// toggleGroupBody renders the joined single-selection group and its state.
func toggleGroupBody(c *ui.Context, st *buttonsState) {
	mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{
		Orientation: mygoui.Horizontal, Separators: true,
	},
		mygoui.ToggleButtonOptions{
			Label: "Left", Icon: lucide.AlignLeft(), Selected: &st.AlignLeft,
			OnChange: func(on bool) {
				if on {
					st.AlignCenter, st.AlignRight = false, false
				}
			},
		},
		mygoui.ToggleButtonOptions{
			Label: "Center", Icon: lucide.AlignCenter(), Selected: &st.AlignCenter,
			OnChange: func(on bool) {
				if on {
					st.AlignLeft, st.AlignRight = false, false
				}
			},
		},
		mygoui.ToggleButtonOptions{
			Label: "Right", Icon: lucide.AlignRight(), Selected: &st.AlignRight,
			OnChange: func(on bool) {
				if on {
					st.AlignLeft, st.AlignCenter = false, false
				}
			},
		},
	)
	alignment := "None"
	switch {
	case st.AlignLeft:
		alignment = "Left"
	case st.AlignCenter:
		alignment = "Center"
	case st.AlignRight:
		alignment = "Right"
	}
	mygoui.Typography(c, "Alignment: "+alignment, mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Muted: true,
	})
}

// toggleGroupVerticalBody renders a full-width vertical group.
func toggleGroupVerticalBody(c *ui.Context, st *buttonsState) {
	mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{
		Orientation: mygoui.Vertical, FullWidth: true,
	},
		mygoui.ToggleButtonOptions{
			Label: "Overflow", Icon: lucide.Layers(), Selected: &st.OverflowOn,
			FullWidth: true,
			OnChange: func(on bool) {
				if on {
					st.FiltersOn = false
				}
			},
		},
		mygoui.ToggleButtonOptions{
			Label: "Filters", Icon: lucide.Filter(), Selected: &st.FiltersOn,
			FullWidth: true,
			OnChange: func(on bool) {
				if on {
					st.OverflowOn = false
				}
			},
		},
	)
}

// toggleGroupDetachedBody renders separated toggles.
func toggleGroupDetachedBody(c *ui.Context, st *buttonsState) {
	mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{Detached: true},
		mygoui.ToggleButtonOptions{
			Label: "Grid", Icon: lucide.Layers(), Variant: mygoui.ToggleGhost,
			Selected: &st.GridOn, OnChange: func(bool) {},
		},
		mygoui.ToggleButtonOptions{
			Label: "List", Icon: lucide.Filter(), Variant: mygoui.ToggleGhost,
			Selected: &st.ListOn, OnChange: func(bool) {},
		},
	)
}

// pressableBody renders the pressable surface and its press count.
func pressableBody(c *ui.Context, st *buttonsState) {
	t := theme.Colors(c)
	mygoui.PressableFeedback(c, mygoui.PressableFeedbackOptions{
		Label:   "Pressable card",
		OnClick: func() { st.Presses++ },
	}, func() {
		ui.Row(c).AlignItems(ui.Center).Gap(theme.Space12).
			Padding(theme.Space12).MinWidth(theme.Scale(c, 320)).Children(func() {
			ui.Icon(c, lucide.Sparkles()).Size(theme.Scale(c, theme.IconLG), theme.Scale(c, theme.IconLG)).
				TextColor(t.Accent).Shrink(0)
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
				mygoui.Typography(c, "Press feedback", mygoui.TypographyOptions{Type: theme.TypeBodySM})
				mygoui.Typography(c, "A 0.985 press scale and surface fade follow the pointer.",
					mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
			})
		})
	})
	mygoui.Typography(c, fmt.Sprintf("Presses: %d", st.Presses), mygoui.TypographyOptions{
		Type: theme.TypeBodyXS, Muted: true,
	})
}

// pressableFullWidthBody renders the full-width feedback row.
func pressableFullWidthBody(c *ui.Context, st *buttonsState) {
	t := theme.Colors(c)
	mygoui.PressableFeedback(c, mygoui.PressableFeedbackOptions{
		Label: "Pressable row", FullWidth: true,
		OnClick: func() { st.Presses++ },
	}, func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space12).
			Padding(theme.Space12).Children(func() {
			ui.Icon(c, lucide.Mail()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
				TextColor(t.Muted).Shrink(0)
			mygoui.Typography(c, "Open release notes", mygoui.TypographyOptions{
				Type: theme.TypeBodySM,
			})
			ui.Spacer(c)
			mygoui.Chip(c, mygoui.ChipOptions{
				Label: "Press", Color: mygoui.ColorAccent,
				Variant: mygoui.ChipSoft, Size: mygoui.ChipSmall,
			})
		})
	})
}

// pressableDisabledBody renders the feedback-free disabled surface.
func pressableDisabledBody(c *ui.Context) {
	t := theme.Colors(c)
	mygoui.PressableFeedback(c, mygoui.PressableFeedbackOptions{
		Label: "Disabled pressable", Disabled: true, FullWidth: true,
	}, func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(theme.Space12).
			Padding(theme.Space12).Children(func() {
			ui.Icon(c, lucide.Lock()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
				TextColor(t.Muted).Shrink(0)
			mygoui.Typography(c, "Unavailable while publishing", mygoui.TypographyOptions{
				Type: theme.TypeBodySM, Muted: true,
			})
		})
	})
}
