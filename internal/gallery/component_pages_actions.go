package gallery

import (
	"github.com/egoist/mygo/ui"
)

// buttonsStateOf keeps the shared specimen state for the action demos.
func buttonsStateOf(s *State) *buttonsState {
	st := StateOf[buttonsState](s, "buttons")
	st.prepare()
	return st
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "button", Title: "Button", Group: "Actions", Order: 10,
		Description: "Native actions with Coss variants, sizes and states.",
		Examples: []ComponentExample{
			{
				Title:       "Variants",
				Description: "Default, secondary, outline, ghost, link and destructive treatments.",
				Demo: func(c *ui.Context, s *State) {
					buttonVariantsBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Sizes",
				Description: "Small, medium and large heights with matching padding.",
				Demo: func(c *ui.Context, s *State) {
					buttonSizesBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "With icons",
				Description: "Leading and trailing icons plus icon-only layouts.",
				Demo: func(c *ui.Context, s *State) {
					buttonIconsBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Loading",
				Description: "A leading spinner keeps the label and blocks activation while pending.",
				Demo: func(c *ui.Context, s *State) {
					buttonLoadingBody(c)
				},
			},
			{
				Title:       "Disabled",
				Description: "Dimmed surfaces that no longer accept activation.",
				Demo: func(c *ui.Context, s *State) {
					buttonDisabledBody(c)
				},
			},
		},
		Code: `package main

import "github.com/Cyenoch/mygoui"
import "github.com/egoist/mygo/ui"
func save(c *ui.Context) {
    mygoui.Button(c, mygoui.ButtonOptions{
        Label: "Save", Variant: mygoui.ButtonPrimary,
        OnClick: func() {},
    })
}`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Visible and accessible text."},
			{"Variant", "mygoui.ButtonVariant", "ButtonPrimary", "Color and surface treatment."},
			{"Size", "theme.ButtonSize", "theme.ButtonMedium", "Small, medium or large height and padding."},
			{"Icon", "*ui.SVG", "nil", "Leading icon."},
			{"EndIcon", "*ui.SVG", "nil", "Trailing icon."},
			{"IconOnly", "bool", "false", "A square control named by Label."},
			{"FullWidth", "bool", "false", "Stretches the button across its row."},
			{"Loading", "bool", "false", "Leads a spinner and blocks activation."},
			{"Disabled", "bool", "false", "Paints at 50% opacity and removes activation."},
			{"Align", "mygoui.ButtonAlign", "ButtonAlignCenter", "Centers the label or starts it for nav rows."},
			{"Tooltip", "string", "", "Native hover description; empty uses Label when icon-only."},
		},
		Tokens: []string{"Primary", "PrimaryForeground", "Secondary", "Accent", "Destructive"},
	})

	RegisterComponent(ComponentPage{
		ID: "buttongroup", Title: "ButtonGroup", Group: "Actions", Order: 20,
		Description: "Joined buttons that round only their outer corners.",
		Examples: []ComponentExample{
			{
				Title:       "Horizontal",
				Description: "A joined row with the native library's half-height shared-edge separators.",
				Demo: func(c *ui.Context, s *State) {
					buttonGroupBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Vertical",
				Description: "A full-width vertical stack that rounds only its ends.",
				Demo: func(c *ui.Context, s *State) {
					buttonGroupVerticalBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Outline",
				Description: "Outline segments collapse their shared borders to one hairline.",
				Demo: func(c *ui.Context, s *State) {
					buttonGroupOutlineBody(c, buttonsStateOf(s))
				},
			},
		},
		Code: `mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Separators: true},
    mygoui.ButtonOptions{Label: "Day", Variant: mygoui.ButtonSecondary},
    mygoui.ButtonOptions{Label: "Week", Variant: mygoui.ButtonSecondary},
)`,
		Params: []ComponentParam{
			{"Orientation", "mygoui.ButtonGroupOrientation", "GroupHorizontal", "Horizontal or vertical stack."},
			{"FullWidth", "bool", "false", "Stretches the group across its row."},
			{"Separators", "bool", "false", "Paints a half-height divider at each shared edge."},
		},
		Tokens: []string{"Default", "Border", "Focus"},
	})

	RegisterComponent(ComponentPage{
		ID: "closebutton", Title: "CloseButton", Group: "Actions", Order: 30,
		Description: "A compact dismiss action for dialogs, chips and panels.",
		Examples: []ComponentExample{
			{
				Title:       "In context",
				Description: "A dismiss action beside the notice it closes, reporting each press.",
				Demo: func(c *ui.Context, s *State) {
					closeButtonBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Disabled",
				Description: "The control dims and stops accepting activation.",
				Demo: func(c *ui.Context, s *State) {
					closeButtonDisabledBody(c)
				},
			},
		},
		Code: `mygoui.CloseButton(c, mygoui.CloseButtonOptions{
    Label: "Dismiss", OnClick: dismiss,
})`,
		Params: []ComponentParam{
			{"Label", "string", "Close", "Accessible name announced by assistive technology."},
			{"Disabled", "bool", "false", "Dims the button and removes activation."},
			{"OnClick", "func()", "nil", "Runs on click, Space or Return."},
		},
		Tokens: []string{"Default", "DefaultHover", "Focus"},
	})

	RegisterComponent(ComponentPage{
		ID: "togglebutton", Title: "ToggleButton", Group: "Actions", Order: 40,
		Description: "A button that stays pressed until its caller-owned state flips.",
		Examples: []ComponentExample{
			{
				Title:       "Default",
				Description: "The default surface in its unselected and selected states.",
				Demo: func(c *ui.Context, s *State) {
					toggleDefaultBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Ghost",
				Description: "A transparent surface that gains the soft accent when pressed.",
				Demo: func(c *ui.Context, s *State) {
					toggleGhostBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "With icons",
				Description: "Icon and icon-only toggles keep the same pressed paint.",
				Demo: func(c *ui.Context, s *State) {
					toggleIconsBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Disabled",
				Description: "Disabled toggles keep their caller state but cannot flip it.",
				Demo: func(c *ui.Context, s *State) {
					toggleDisabledBody(c)
				},
			},
		},
		Code: `selected := true
mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
    Label: "Bold", Variant: mygoui.ToggleGhost, Selected: &selected,
    OnChange: func(on bool) { bold = on },
})`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Visible and accessible text."},
			{"Variant", "mygoui.ToggleButtonVariant", "ToggleDefault", "Default surface or ghost until selected."},
			{"Selected", "*bool", "nil", "Caller-owned pressed state."},
			{"Size", "theme.ButtonSize", "theme.ButtonMedium", "Small, medium or large."},
			{"Icon", "*ui.SVG", "nil", "Leading icon."},
			{"IconOnly", "bool", "false", "A square control named by Label."},
			{"FullWidth", "bool", "false", "Stretches the toggle across its row."},
			{"Disabled", "bool", "false", "Paints at 50% opacity and keeps the caller state."},
			{"OnChange", "func(bool)", "nil", "Reports the new pressed state."},
		},
		Tokens: []string{"AccentSoft", "AccentSoftForeground", "Default", "DefaultHover"},
	})

	RegisterComponent(ComponentPage{
		ID: "togglebuttongroup", Title: "ToggleButtonGroup", Group: "Actions", Order: 50,
		Description: "Joined toggles for segmented choices, each with its own state.",
		Examples: []ComponentExample{
			{
				Title:       "Horizontal",
				Description: "One joined control with single-selection behavior driven by caller state.",
				Demo: func(c *ui.Context, s *State) {
					toggleGroupBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Vertical",
				Description: "A full-width vertical group for narrow rails.",
				Demo: func(c *ui.Context, s *State) {
					toggleGroupVerticalBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Detached",
				Description: "Separated toggles that keep their own corners.",
				Demo: func(c *ui.Context, s *State) {
					toggleGroupDetachedBody(c, buttonsStateOf(s))
				},
			},
		},
		Code: `mygoui.ToggleButtonGroup(c, mygoui.ToggleButtonGroupOptions{Separators: true},
    mygoui.ToggleButtonOptions{Label: "Day", Selected: &day},
    mygoui.ToggleButtonOptions{Label: "Week", Selected: &week},
)`,
		Params: []ComponentParam{
			{"Orientation", "mygoui.Orientation", "Horizontal", "Horizontal or vertical group."},
			{"FullWidth", "bool", "false", "Stretches the group across its row."},
			{"Detached", "bool", "false", "Separates the toggles instead of joining them."},
			{"Separators", "bool", "false", "Paints a half-height divider at each shared edge."},
		},
		Tokens: []string{"AccentSoft", "AccentSoftForeground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "pressablefeedback", Title: "PressableFeedback", Group: "Actions", Order: 60,
		Description: "A native press highlight for content that is not a button.",
		Examples: []ComponentExample{
			{
				Title:       "Press feedback",
				Description: "A pressable surface with the native control's 0.985 press scale.",
				Demo: func(c *ui.Context, s *State) {
					pressableBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Full width",
				Description: "The highlight stretches across its row while the content stays inline.",
				Demo: func(c *ui.Context, s *State) {
					pressableFullWidthBody(c, buttonsStateOf(s))
				},
			},
			{
				Title:       "Disabled",
				Description: "Disabled content keeps its layout and drops the feedback.",
				Demo: func(c *ui.Context, s *State) {
					pressableDisabledBody(c)
				},
			},
		},
		Code: `mygoui.PressableFeedback(c, mygoui.PressableFeedbackOptions{
    Label: "Copy link", OnClick: copyLink, FullWidth: true,
}, func() {
    mygoui.Typography(c, "Copy link", mygoui.TypographyOptions{
        Type: theme.TypeBodySM,
    })
})`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Accessible name for the pressable area."},
			{"FullWidth", "bool", "false", "Stretches the highlight across its row."},
			{"Disabled", "bool", "false", "Removes hover and press feedback."},
			{"OnClick", "func()", "nil", "Runs when the content is pressed."},
		},
		Tokens: []string{"Default", "DefaultHover"},
	})
}
