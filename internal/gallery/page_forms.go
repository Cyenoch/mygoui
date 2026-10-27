package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// formsState is the specimen's caller-owned data. Every control written to a
// pointer lives here so the page keeps its values across frames.
type formsState struct {
	ready bool

	name, email, username, password, bare, bio string
	accountID, disabledSearch                  string
	price, website                             string
	lastEdit                                   string

	otp     form.OTPState
	otpNote string
	search  string

	quantity         form.NumberFieldState
	quantityDisabled form.NumberFieldState

	terms          bool
	partial        bool
	partialChecked bool
	secondaryCheck bool
	features       [3]bool
	checks         form.CheckboxGroupState
	channels       [2]bool

	plan      string
	planRadio form.RadioGroupState[string]
	shipping  string
	shipRadio form.RadioGroupState[string]
	seat      string
	seatRadio form.RadioGroupState[string]
	billing   string
	billRadio form.RadioGroupState[string]

	wifi, bluetooth, airplane, locked bool
	analytics, sync, syncInvalid      bool
	controlSwitch, controlCheckbox    bool
	volume, rangeLow, rangeHigh, zoom float64
	lockedVolume                      float64
}

// initState fills the specimen's defaults once. It is lazy because StateOf
// hands back a zero value on the page's first frame.
func (st *formsState) initState() {
	if st.ready {
		return
	}
	st.ready = true

	st.name = "Ada Lovelace"
	st.email = "ada@example.com"
	st.username = "ada"
	st.password = "correct-horse-battery"
	st.accountID = "acct_1042"
	st.disabledSearch = "archived query"
	st.bio = "Building native interfaces with MyGo and Coss."
	st.price = "499.00"
	st.website = "coss.com"
	st.lastEdit = "Nothing changed yet"

	st.otp.Value = "1234"
	st.otpNote = "Four of six digits entered."

	st.quantity.Value = 2
	st.quantityDisabled.Value = 8

	st.terms = true
	st.partial = true
	st.partialChecked = false
	st.features = [3]bool{true, false, true}
	st.checks.Items = []*bool{&st.features[0], &st.features[1], &st.features[2]}
	st.channels = [2]bool{true, false}

	st.plan = "pro"
	st.planRadio = form.RadioGroupState[string]{Selected: &st.plan}
	st.shipping = ""
	st.shipRadio = form.RadioGroupState[string]{Selected: &st.shipping}
	st.seat = "standard"
	st.seatRadio = form.RadioGroupState[string]{Selected: &st.seat}
	st.billing = "monthly"
	st.billRadio = form.RadioGroupState[string]{Selected: &st.billing}

	st.wifi = true
	st.analytics = true
	st.volume = 64
	st.rangeLow = 25
	st.rangeHigh = 75
	st.zoom = 40
	st.lockedVolume = 30
}

// formsFieldStack keeps fields close to their natural reading width while
// still filling narrow cards.
func formsFieldStack(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, 520)).
		Gap(theme.Space16).Children(body)
}

// formsMutedText is a small helper for inline addon text.
func formsMutedText(c *ui.Context, text string) {
	mygoui.Typography(c, text, mygoui.TypographyOptions{
		Type: theme.TypeBodySM, Muted: true,
	})
}

// formsExampleDemo adapts one single-variant forms specimen to a named
// ComponentExample demo. The shared field stack keeps the control at a
// readable width inside the example panel.
func formsExampleDemo(body func(c *ui.Context, st *formsState)) func(c *ui.Context, s *State) {
	return func(c *ui.Context, s *State) {
		st := formsPageState(s)
		formsFieldStack(c, func() { body(c, st) })
	}
}

// formsFormExampleDefault shows a labelled field group with supporting copy.
func formsFormExampleDefault(c *ui.Context, st *formsState) {
	form.Form(c, form.FormOptions{Label: "Account form"}, func() {
		form.FieldGroup(c, func() {
			form.TextField(c, form.TextFieldOptions{
				Value:       &st.email,
				Label:       "Email",
				Description: "We send account notices here.",
				Required:    true,
				FullWidth:   true,
			})
		})
	})
}

// formsFormExampleInvalid shows the form's top-level failure message.
func formsFormExampleInvalid(c *ui.Context, st *formsState) {
	form.Form(c, form.FormOptions{Label: "Account form", Invalid: true}, func() {
		form.FieldGroup(c, func() {
			form.TextField(c, form.TextFieldOptions{
				Value:     &st.email,
				Label:     "Email",
				Required:  true,
				FullWidth: true,
			})
		})
	})
}

// formsFormExampleActions shows the field actions row under the fields.
func formsFormExampleActions(c *ui.Context, st *formsState) {
	form.Form(c, form.FormOptions{Label: "Account form"}, func() {
		form.FieldGroup(c, func() {
			form.TextField(c, form.TextFieldOptions{
				Value:     &st.email,
				Label:     "Email",
				Required:  true,
				FullWidth: true,
			})
		})
		form.FieldActions(c, func() {
			mygoui.Button(c.Key("form-submit-example"), mygoui.ButtonOptions{
				Label: "Submit form", Size: theme.ButtonSmall,
				OnClick: func() { st.lastEdit = "Form submitted" },
			})
			mygoui.Button(c.Key("form-reset-example"), mygoui.ButtonOptions{
				Label:   "Reset",
				Variant: mygoui.ButtonGhost,
				Size:    theme.ButtonSmall,
				OnClick: func() { st.email = "ada@example.com" },
			})
		})
	})
}

// formsFieldsetExampleDefault groups check boxes under a legend and copy.
func formsFieldsetExampleDefault(c *ui.Context, st *formsState) {
	form.Fieldset(c, form.FieldsetOptions{
		Legend:      "Notifications",
		Description: "Choose what we send to your inbox.",
		FullWidth:   true,
	}, func() {
		form.Checkbox(c, form.CheckboxOptions{
			Checked:     &st.terms,
			Label:       "Product updates",
			Description: "New features and fixes.",
		})
		form.Checkbox(c, form.CheckboxOptions{
			Checked:     &st.partialChecked,
			Label:       "Security alerts",
			Description: "Sign-ins and policy changes.",
		})
	})
}

// formsFieldsetExampleDisabled shows a fieldset that disables its controls.
func formsFieldsetExampleDisabled(c *ui.Context, st *formsState) {
	form.Fieldset(c, form.FieldsetOptions{
		Legend:      "Managed by your administrator",
		Description: "Administrators control these settings.",
		Disabled:    true,
		FullWidth:   true,
	}, func() {
		form.Checkbox(c, form.CheckboxOptions{
			Checked:  &st.features[0],
			Label:    "Weekly digest",
			Disabled: true,
		})
	})
}

// formsLabelExampleDefault shows the default field label.
func formsLabelExampleDefault(c *ui.Context, st *formsState) {
	form.Label(c, "Display name", form.LabelOptions{})
}

// formsLabelExampleRequired shows the required marker.
func formsLabelExampleRequired(c *ui.Context, st *formsState) {
	form.Label(c, "Email address", form.LabelOptions{Required: true})
}

// formsLabelExampleInvalid paints the label with the danger color.
func formsLabelExampleInvalid(c *ui.Context, st *formsState) {
	form.Label(c, "Account ID", form.LabelOptions{Invalid: true})
}

// formsLabelExampleDisabled dims the label and drops its focus target.
func formsLabelExampleDisabled(c *ui.Context, st *formsState) {
	form.Label(c, "Managed by your organization", form.LabelOptions{Disabled: true})
}

// formsDescriptionExampleStandalone shows help copy under a label.
func formsDescriptionExampleStandalone(c *ui.Context, st *formsState) {
	ui.Column(c).Gap(theme.Space4).Children(func() {
		form.Label(c, "API key", form.LabelOptions{})
		form.Description(c, "Rotate keys every ninety days.")
	})
}

// formsDescriptionExampleField shows the description a field supplies.
func formsDescriptionExampleField(c *ui.Context, st *formsState) {
	form.Field(c, form.FieldOptions{
		Label:       "Slug",
		Description: "Appears in the workspace URL.",
	}, func() {
		form.Input(c, &st.username, form.InputOptions{FullWidth: true})
	})
}

// formsDescriptionExampleHidden shows the error replacing the description.
func formsDescriptionExampleHidden(c *ui.Context, st *formsState) {
	form.Field(c, form.FieldOptions{
		Label:       "Slug",
		Description: "Appears in the workspace URL.",
		Invalid:     true,
		Error:       "This slug is already taken.",
	}, func() {
		form.Input(c, &st.username, form.InputOptions{Invalid: true, FullWidth: true})
	})
}

// formsFieldErrorExampleStandalone shows validation copy under a label.
func formsFieldErrorExampleStandalone(c *ui.Context, st *formsState) {
	ui.Column(c).Gap(theme.Space4).Children(func() {
		form.Label(c, "Email", form.LabelOptions{})
		form.FieldError(c, "Enter an address containing @.")
	})
}

// formsFieldErrorExampleField shows the same copy through a field root.
func formsFieldErrorExampleField(c *ui.Context, st *formsState) {
	form.Field(c, form.FieldOptions{
		Label:   "Username",
		Invalid: true,
		Error:   "This username is already taken.",
	}, func() {
		form.Input(c, &st.username, form.InputOptions{Invalid: true, FullWidth: true})
	})
}

// formsErrorMessageExampleCustom shows low-level error copy in a layout.
func formsErrorMessageExampleCustom(c *ui.Context, st *formsState) {
	ui.Column(c).Gap(theme.Space8).Children(func() {
		form.Label(c, "Password", form.LabelOptions{Required: true})
		form.Input(c, &st.password, form.InputOptions{
			Password: true, Invalid: true, FullWidth: true,
		})
		form.ErrorMessage(c, "Use at least twelve characters.")
	})
}

// formsErrorMessageExampleForm shows a form-level failure above its fields.
func formsErrorMessageExampleForm(c *ui.Context, st *formsState) {
	form.Form(c, form.FormOptions{Label: "Profile", Invalid: true}, func() {
		form.Input(c, &st.name, form.InputOptions{FullWidth: true})
	})
}

// formsInputExampleDefault shows a bare editor with an accessible name.
func formsInputExampleDefault(c *ui.Context, st *formsState) {
	form.Input(c, &st.bare, form.InputOptions{
		AriaLabel:   "Project name",
		Placeholder: "Bare input without a label",
		FullWidth:   true,
	})
}

// formsInputExampleDescription shows the field composition with help copy.
func formsInputExampleDescription(c *ui.Context, st *formsState) {
	form.Input(c, &st.username, form.InputOptions{
		Label:       "Username",
		Description: "Letters and digits only.",
		Required:    true,
		FullWidth:   true,
	})
}

// formsInputExampleInvalid shows the danger border and error message.
func formsInputExampleInvalid(c *ui.Context, st *formsState) {
	form.Input(c, &st.email, form.InputOptions{
		Label:     "Email",
		Invalid:   true,
		Error:     "Enter an address containing @.",
		Variant:   form.ControlSecondary,
		FullWidth: true,
	})
}

// formsInputExamplePassword masks the typed text.
func formsInputExamplePassword(c *ui.Context, st *formsState) {
	form.Input(c, &st.password, form.InputOptions{
		Label:     "Password",
		Password:  true,
		FullWidth: true,
	})
}

// formsInputExampleDisabled stops editing and dims the surface.
func formsInputExampleDisabled(c *ui.Context, st *formsState) {
	form.Input(c, &st.accountID, form.InputOptions{
		Label:     "Account ID",
		Disabled:  true,
		FullWidth: true,
	})
}

// formsTextAreaExampleDefault shows the growing multiline editor.
func formsTextAreaExampleDefault(c *ui.Context, st *formsState) {
	form.TextArea(c, &st.bio, form.TextAreaOptions{
		Label:       "Bio",
		Description: "Grows to six lines, then scrolls.",
		Placeholder: "Tell your team what you are working on",
		MinLines:    3,
		MaxLines:    6,
		FullWidth:   true,
	})
}

// formsTextAreaExampleReadOnly shows selectable but uneditable text.
func formsTextAreaExampleReadOnly(c *ui.Context, st *formsState) {
	form.TextArea(c, &st.bio, form.TextAreaOptions{
		Label:     "Meeting notes",
		ReadOnly:  true,
		MinLines:  2,
		MaxLines:  4,
		FullWidth: true,
	})
}

// formsTextAreaExampleInvalid shows the danger border and error message.
func formsTextAreaExampleInvalid(c *ui.Context, st *formsState) {
	form.TextArea(c, &st.bio, form.TextAreaOptions{
		Label:     "Release notes",
		Invalid:   true,
		Error:     "Keep the summary under 400 characters.",
		MinLines:  2,
		FullWidth: true,
	})
}

// formsTextFieldExampleDefault shows the complete field in one call.
func formsTextFieldExampleDefault(c *ui.Context, st *formsState) {
	form.TextField(c, form.TextFieldOptions{
		Value:       &st.name,
		Label:       "Full name",
		Placeholder: "Ada Lovelace",
		FullWidth:   true,
	})
}

// formsTextFieldExampleDescription adds supporting copy under the label.
func formsTextFieldExampleDescription(c *ui.Context, st *formsState) {
	form.TextField(c, form.TextFieldOptions{
		Value:       &st.email,
		Label:       "Email",
		Description: "We only use this for account notices.",
		Placeholder: "you@example.com",
		Required:    true,
		FullWidth:   true,
	})
}

// formsTextFieldExampleInvalid shows validation copy from the field.
func formsTextFieldExampleInvalid(c *ui.Context, st *formsState) {
	form.TextField(c, form.TextFieldOptions{
		Value:       &st.username,
		Label:       "Username",
		Error:       "This username is already taken.",
		Invalid:     true,
		Placeholder: "username",
		FullWidth:   true,
	})
}

// formsTextFieldExamplePassword masks the value with a required label.
func formsTextFieldExamplePassword(c *ui.Context, st *formsState) {
	form.TextField(c, form.TextFieldOptions{
		Value:       &st.password,
		Label:       "Password",
		Description: "At least twelve characters.",
		Password:    true,
		Required:    true,
		FullWidth:   true,
	})
}

// formsTextFieldExampleDisabled shows a caller-owned value that cannot change.
func formsTextFieldExampleDisabled(c *ui.Context, st *formsState) {
	form.TextField(c, form.TextFieldOptions{
		Value:       &st.accountID,
		Label:       "Account ID",
		Description: "Managed by your organization.",
		Disabled:    true,
		FullWidth:   true,
	})
}

// formsInputGroupExamplePrefixSuffix shows inline addon regions.
func formsInputGroupExamplePrefixSuffix(c *ui.Context, st *formsState) {
	form.InputGroup(c, &st.price, form.InputGroupOptions{
		Label:       "Price",
		Description: "List price in US dollars.",
		Placeholder: "0.00",
		FullWidth:   true,
		Prefix: func() {
			formsMutedText(c, "$")
		},
		Suffix: func() {
			mygoui.Button(c.Key("max-price-example"), mygoui.ButtonOptions{
				Label:   "Max",
				Tooltip: "Set the maximum price",
				Variant: mygoui.ButtonGhost,
				Size:    theme.ButtonSmall,
				OnClick: func() { st.price = "999.00" },
			})
		},
	})
}

// formsInputGroupExampleBlock shows full-width rows above and below the editor.
func formsInputGroupExampleBlock(c *ui.Context, st *formsState) {
	form.InputGroup(c, &st.website, form.InputGroupOptions{
		Label:       "Public URL",
		Description: "Changes apply as soon as you save.",
		Placeholder: "your-workspace",
		FullWidth:   true,
		Prefix: func() {
			formsMutedText(c, "https://")
		},
		Suffix: func() {
			formsMutedText(c, ".mygo.dev")
		},
		BlockPrefix: func() {
			mygoui.Typography(c, "Address preview", mygoui.TypographyOptions{
				Type: theme.TypeBodyXS, Muted: true,
			})
		},
		BlockSuffix: func() {
			formsMutedText(c, fmt.Sprintf("%d characters", len(st.website)))
		},
	})
}

// formsInputGroupExampleInvalid paints the danger shell and error message.
func formsInputGroupExampleInvalid(c *ui.Context, st *formsState) {
	form.InputGroup(c, &st.website, form.InputGroupOptions{
		Label:     "Public URL",
		Invalid:   true,
		Error:     "Enter a workspace name without spaces.",
		FullWidth: true,
		Prefix: func() {
			formsMutedText(c, "https://")
		},
		Suffix: func() {
			formsMutedText(c, ".mygo.dev")
		},
	})
}

// formsOTPExampleDefault shows the six painted code slots.
func formsOTPExampleDefault(c *ui.Context, st *formsState) {
	form.InputOTP(c, &st.otp, form.OTPOptions{
		Label:       "Verification code",
		Description: "We sent a six-digit code to ada@example.com.",
		Length:      6,
	})
}

// formsOTPExampleInvalid paints every slot with the danger border.
func formsOTPExampleInvalid(c *ui.Context, st *formsState) {
	form.InputOTP(c, &st.otp, form.OTPOptions{
		Label:   "Verification code",
		Invalid: true,
		Error:   "That code is incorrect or expired.",
		Length:  6,
	})
}

// formsOTPExampleDisabled stops editing and dims the slots.
func formsOTPExampleDisabled(c *ui.Context, st *formsState) {
	form.InputOTP(c, &st.otp, form.OTPOptions{
		Label:       "Verification code",
		Description: "Confirmed and locked for this session.",
		Length:      6,
		Disabled:    true,
	})
}

// formsSearchFieldExampleDefault shows the search editor and clear action.
func formsSearchFieldExampleDefault(c *ui.Context, st *formsState) {
	form.SearchField(c, &st.search, form.SearchFieldOptions{
		Label:       "Search components",
		Placeholder: "Search by name",
		FullWidth:   true,
	})
}

// formsSearchFieldExampleInvalid shows the danger border and error message.
func formsSearchFieldExampleInvalid(c *ui.Context, st *formsState) {
	form.SearchField(c, &st.search, form.SearchFieldOptions{
		Label:     "Filter",
		Invalid:   true,
		Error:     "Enter at least two characters.",
		FullWidth: true,
	})
}

// formsSearchFieldExampleDisabled stops editing and disables the clear action.
func formsSearchFieldExampleDisabled(c *ui.Context, st *formsState) {
	form.SearchField(c, &st.disabledSearch, form.SearchFieldOptions{
		Label:       "Archived queries",
		Placeholder: "Search past sessions",
		FullWidth:   true,
		Disabled:    true,
	})
}

// formsNumberFieldExampleBounded shows stepper buttons between bounds. The
// field keeps the native library's intrinsic 200 DIP group instead of stretching.
func formsNumberFieldExampleBounded(c *ui.Context, st *formsState) {
	ui.Column(c).Width(theme.Scale(c, 200)).Children(func() {
		form.NumberField(c, &st.quantity, form.NumberFieldOptions{
			Label:       "Seats",
			Description: "Between one and twelve, in whole seats.",
			Min:         new(float64(1)),
			Max:         new(float64(12)),
			Step:        1,
		})
	})
}

// formsNumberFieldExampleDisabled stops stepping and typed changes.
func formsNumberFieldExampleDisabled(c *ui.Context, st *formsState) {
	ui.Column(c).Width(theme.Scale(c, 200)).Children(func() {
		form.NumberField(c, &st.quantityDisabled, form.NumberFieldOptions{
			Label:       "Seats",
			Description: "Billing controls this value.",
			Min:         new(float64(1)),
			Max:         new(float64(12)),
			Step:        1,
			Disabled:    true,
		})
	})
}

// formsCheckboxExampleDefault shows a single caller-owned check box.
func formsCheckboxExampleDefault(c *ui.Context, st *formsState) {
	form.Checkbox(c, form.CheckboxOptions{
		Checked: &st.terms,
		Label:   "Release notes",
	})
}

// formsCheckboxExampleSecondary shows the flat secondary control face.
func formsCheckboxExampleSecondary(c *ui.Context, st *formsState) {
	form.Checkbox(c, form.CheckboxOptions{
		Checked:     &st.secondaryCheck,
		Label:       "Only on Wi-Fi",
		Description: "Uses the flat secondary control face.",
		Variant:     form.ControlSecondary,
	})
}

// formsCheckboxExampleDescription adds supporting copy under the label.
func formsCheckboxExampleDescription(c *ui.Context, st *formsState) {
	form.Checkbox(c, form.CheckboxOptions{
		Checked:     &st.analytics,
		Label:       "Product updates",
		Description: "Email me when a new version ships.",
	})
}

// formsCheckboxExampleMixed shows the indeterminate mark.
func formsCheckboxExampleMixed(c *ui.Context, st *formsState) {
	form.Checkbox(c, form.CheckboxOptions{
		Checked:       &st.partialChecked,
		Indeterminate: &st.partial,
		Label:         "Parent folder",
		Description:   "Mixed while only some children are selected.",
	})
}

// formsCheckboxExampleInvalid paints the danger control and error.
func formsCheckboxExampleInvalid(c *ui.Context, st *formsState) {
	form.Checkbox(c, form.CheckboxOptions{
		Checked:  &st.syncInvalid,
		Label:    "Background sync",
		Invalid:  true,
		Error:    "Background sync requires a network connection.",
		Required: true,
	})
}

// formsCheckboxExampleDisabled removes activation and dims the row.
func formsCheckboxExampleDisabled(c *ui.Context, st *formsState) {
	form.Checkbox(c, form.CheckboxOptions{
		Checked:  &st.locked,
		Label:    "Locked by policy",
		Disabled: true,
	})
}

// formsCheckboxGroupExampleDefault shows related check boxes under one label.
func formsCheckboxGroupExampleDefault(c *ui.Context, st *formsState) {
	form.CheckboxGroup(c, form.CheckboxGroupOptions{
		Label:       "Delivery channels",
		Description: "Pick every channel you want to use.",
	}, func() {
		form.Checkbox(c.Key("channel-email-example"), form.CheckboxOptions{
			Checked: &st.channels[0],
			Label:   "Email",
		})
		form.Checkbox(c.Key("channel-sms-example"), form.CheckboxOptions{
			Checked: &st.channels[1],
			Label:   "SMS",
		})
	})
}

// formsCheckboxGroupExampleSelectAll adds the mixed select-all row.
func formsCheckboxGroupExampleSelectAll(c *ui.Context, st *formsState) {
	form.CheckboxGroup(c, form.CheckboxGroupOptions{
		Label:          "Notifications",
		State:          &st.checks,
		SelectAllLabel: "All notifications",
	}, func() {
		labels := [3]string{"Product updates", "Weekly digest", "Security alerts"}
		for i := range st.features {
			form.Checkbox(c.Key(fmt.Sprintf("select-all-%d", i)), form.CheckboxOptions{
				Checked: &st.features[i],
				Label:   labels[i],
			})
		}
	})
}

// formsCheckboxGroupExampleHorizontal lays the check boxes in a wrapping row.
func formsCheckboxGroupExampleHorizontal(c *ui.Context, st *formsState) {
	form.CheckboxGroup(c, form.CheckboxGroupOptions{
		Label:       "Delivery channels",
		Orientation: form.GroupHorizontal,
	}, func() {
		form.Checkbox(c.Key("horizontal-email"), form.CheckboxOptions{
			Checked: &st.channels[0],
			Label:   "Email",
		})
		form.Checkbox(c.Key("horizontal-sms"), form.CheckboxOptions{
			Checked: &st.channels[1],
			Label:   "SMS",
		})
	})
}

// formsRadioExampleDefault shows one choice inside a live group.
func formsRadioExampleDefault(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{Label: "Plan"}, func() {
		form.Radio(c, form.RadioOptions[string]{
			State: &st.planRadio,
			Value: "pro",
			Label: "Pro",
		})
	})
}

// formsRadioExampleDescription adds supporting copy under the choice.
func formsRadioExampleDescription(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{Label: "Plan"}, func() {
		form.Radio(c, form.RadioOptions[string]{
			State:       &st.planRadio,
			Value:       "pro",
			Label:       "Pro",
			Description: "Unlimited projects and priority support.",
		})
	})
}

// formsRadioExampleInvalid paints the danger control and error.
func formsRadioExampleInvalid(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{Label: "Plan"}, func() {
		form.Radio(c, form.RadioOptions[string]{
			State:   &st.planRadio,
			Value:   "free",
			Label:   "Free",
			Invalid: true,
			Error:   "This plan is unavailable in your region.",
		})
	})
}

// formsRadioExampleDisabled removes selection and dims the row.
func formsRadioExampleDisabled(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{Label: "Plan"}, func() {
		form.Radio(c, form.RadioOptions[string]{
			State:    &st.planRadio,
			Value:    "team",
			Label:    "Team",
			Disabled: true,
		})
	})
}

// formsRadioGroupExampleDefault shows the vertical single-choice group.
func formsRadioGroupExampleDefault(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{
		Label:       "Plan",
		Description: "Change at any time.",
		FullWidth:   true,
	}, func() {
		plans := []struct{ value, label string }{
			{"free", "Free"},
			{"pro", "Pro"},
			{"team", "Team"},
		}
		for _, plan := range plans {
			form.Radio(c.Key("group-plan-"+plan.value), form.RadioOptions[string]{
				State: &st.planRadio,
				Value: plan.value,
				Label: plan.label,
			})
		}
	})
}

// formsRadioGroupExampleHorizontal lays the choices in a wrapping row.
func formsRadioGroupExampleHorizontal(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{
		Label:       "Billing",
		Orientation: form.GroupHorizontal,
		FullWidth:   true,
	}, func() {
		for _, choice := range []struct{ value, label string }{
			{"monthly", "Monthly"},
			{"yearly", "Yearly"},
		} {
			form.Radio(c.Key("group-billing-"+choice.value), form.RadioOptions[string]{
				State: &st.billRadio,
				Value: choice.value,
				Label: choice.label,
			})
		}
	})
}

// formsRadioGroupExampleInvalid marks the group and each radio invalid.
func formsRadioGroupExampleInvalid(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{
		Label:       "Shipping",
		Description: "Choose a delivery speed.",
		Orientation: form.GroupHorizontal,
		Invalid:     true,
		Error:       "Choose a shipping speed to continue.",
		FullWidth:   true,
	}, func() {
		for _, choice := range []struct{ value, label string }{
			{"standard", "Standard"},
			{"express", "Express"},
		} {
			form.Radio(c.Key("group-shipping-"+choice.value), form.RadioOptions[string]{
				State:   &st.shipRadio,
				Value:   choice.value,
				Label:   choice.label,
				Invalid: true,
			})
		}
	})
}

// formsRadioGroupExampleDisabled disables the group and dims it.
func formsRadioGroupExampleDisabled(c *ui.Context, st *formsState) {
	form.RadioGroup(c, form.RadioGroupOptions{
		Label:       "Seat type",
		Description: "Fixed by your contract.",
		Orientation: form.GroupHorizontal,
		Disabled:    true,
		FullWidth:   true,
	}, func() {
		for _, choice := range []struct{ value, label string }{
			{"standard", "Standard"},
			{"accessible", "Accessible"},
		} {
			form.Radio(c.Key("group-seat-"+choice.value), form.RadioOptions[string]{
				State:    &st.seatRadio,
				Value:    choice.value,
				Label:    choice.label,
				Disabled: true,
			})
		}
	})
}

// formsSwitchExampleDefault shows an on/off row with supporting copy.
func formsSwitchExampleDefault(c *ui.Context, st *formsState) {
	form.Switch(c, form.SwitchOptions{
		On:          &st.wifi,
		Label:       "Wi-Fi",
		Description: "Connect automatically on known networks.",
	})
}

// formsSwitchExampleSizes shows the three track sizes.
func formsSwitchExampleSizes(c *ui.Context, st *formsState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		form.Switch(c, form.SwitchOptions{
			On: &st.bluetooth, Label: "Small", Size: form.SwitchSmall,
		})
		form.Switch(c, form.SwitchOptions{
			On: &st.wifi, Label: "Medium",
		})
		form.Switch(c, form.SwitchOptions{
			On: &st.airplane, Label: "Large", Size: form.SwitchLarge,
		})
	})
}

// formsSwitchExampleLabelBefore puts the label before the track.
func formsSwitchExampleLabelBefore(c *ui.Context, st *formsState) {
	form.Switch(c, form.SwitchOptions{
		On:            &st.bluetooth,
		Label:         "Bluetooth",
		LabelPosition: form.SwitchLabelBefore,
	})
}

// formsSwitchExampleInvalid marks the row invalid and shows its error.
func formsSwitchExampleInvalid(c *ui.Context, st *formsState) {
	form.Switch(c, form.SwitchOptions{
		On:      &st.syncInvalid,
		Label:   "Background sync",
		Invalid: true,
		Error:   "Background sync requires a network connection.",
	})
}

// formsSwitchExampleDisabled removes activation and dims the row.
func formsSwitchExampleDisabled(c *ui.Context, st *formsState) {
	form.Switch(c, form.SwitchOptions{
		On:       &st.locked,
		Label:    "Locked by policy",
		Disabled: true,
	})
}

// formsSwitchGroupExampleDefault shows a vertical group under one label.
func formsSwitchGroupExampleDefault(c *ui.Context, st *formsState) {
	form.SwitchGroup(c, form.SwitchGroupOptions{
		Label:       "Privacy",
		Description: "Applies to every workspace you belong to.",
		FullWidth:   true,
	}, func() {
		form.Switch(c.Key("privacy-analytics"), form.SwitchOptions{
			On:          &st.analytics,
			Label:       "Usage analytics",
			Description: "Share anonymous performance data.",
		})
		form.Switch(c.Key("privacy-sync"), form.SwitchOptions{
			On:          &st.sync,
			Label:       "Sync settings",
			Description: "Keep preferences across devices.",
		})
	})
}

// formsSwitchGroupExampleHorizontal lays the group in a wrapping row.
func formsSwitchGroupExampleHorizontal(c *ui.Context, st *formsState) {
	form.SwitchGroup(c, form.SwitchGroupOptions{
		Label:       "Privacy",
		Orientation: form.GroupHorizontal,
		FullWidth:   true,
	}, func() {
		form.Switch(c.Key("horizontal-analytics"), form.SwitchOptions{
			On: &st.analytics, Label: "Usage analytics",
		})
		form.Switch(c.Key("horizontal-sync"), form.SwitchOptions{
			On: &st.sync, Label: "Sync settings",
		})
	})
}

// formsSwitchGroupExampleDisabled disables every switch and dims the group.
func formsSwitchGroupExampleDisabled(c *ui.Context, st *formsState) {
	form.SwitchGroup(c, form.SwitchGroupOptions{
		Label:       "Managed by policy",
		Description: "Your administrator controls both switches.",
		Orientation: form.GroupHorizontal,
		Disabled:    true,
		FullWidth:   true,
	}, func() {
		form.Switch(c.Key("policy-auto-update"), form.SwitchOptions{
			On: &st.locked, Label: "Auto-update", Disabled: true,
		})
		form.Switch(c.Key("policy-diagnostics"), form.SwitchOptions{
			On: &st.syncInvalid, Label: "Diagnostics", Disabled: true,
		})
	})
}

// formsControlFieldExampleSwitch shows a switch trailing the whole row.
func formsControlFieldExampleSwitch(c *ui.Context, st *formsState) {
	form.ControlField(c, form.ControlFieldOptions{
		Checked:     &st.controlSwitch,
		Label:       "Auto-save",
		Description: "Save drafts a moment after you stop typing.",
		Kind:        form.ControlFieldSwitch,
	})
}

// formsControlFieldExampleCheckbox shows a check box trailing the whole row.
func formsControlFieldExampleCheckbox(c *ui.Context, st *formsState) {
	form.ControlField(c, form.ControlFieldOptions{
		Checked:     &st.controlCheckbox,
		Label:       "Offline mode",
		Description: "Keep the last synced copy on this device.",
		Kind:        form.ControlFieldCheckbox,
	})
}

// formsControlFieldExampleInvalid paints the danger indicator and error.
func formsControlFieldExampleInvalid(c *ui.Context, st *formsState) {
	form.ControlField(c, form.ControlFieldOptions{
		Checked:     &st.controlCheckbox,
		Label:       "Offline mode",
		Description: "Keep the last synced copy on this device.",
		Kind:        form.ControlFieldCheckbox,
		Invalid:     true,
		Error:       "Offline mode needs a downloaded workspace.",
	})
}

// formsSliderExampleDefault shows a single thumb with labelled marks.
func formsSliderExampleDefault(c *ui.Context, st *formsState) {
	form.Slider(c, form.SliderOptions{
		State:       form.SliderState{Value: &st.volume},
		Label:       "Volume",
		Description: "Marks show the soft, balanced and loud stops.",
		Min:         0,
		Max:         100,
		Step:        5,
		Marks: []form.SliderMark{
			{Value: 0, Label: "0"},
			{Value: 50, Label: "50"},
			{Value: 100, Label: "100"},
		},
		Format: func(v float64) string { return fmt.Sprintf("%.0f%%", v) },
	})
}

// formsSliderExampleRange shows the ordered two-thumb track.
func formsSliderExampleRange(c *ui.Context, st *formsState) {
	form.Slider(c, form.SliderOptions{
		State:       form.SliderState{Low: &st.rangeLow, High: &st.rangeHigh},
		Label:       "Price range",
		Description: "Thumbs stay ordered; drag either end.",
		Min:         0,
		Max:         100,
		Step:        5,
		Format:      func(v float64) string { return fmt.Sprintf("$%.0f", v) },
	})
}

// formsSliderExampleVertical shows the vertical single-thumb track.
func formsSliderExampleVertical(c *ui.Context, st *formsState) {
	ui.Column(c).Width(theme.Scale(c, 160)).Children(func() {
		form.Slider(c, form.SliderOptions{
			State:    form.SliderState{Value: &st.zoom},
			Label:    "Zoom",
			Min:      0,
			Max:      100,
			Step:     10,
			Vertical: true,
			Height:   theme.Scale(c, 160),
		})
	})
}

// formsSliderExampleDisabled stops dragging and keyboard changes.
func formsSliderExampleDisabled(c *ui.Context, st *formsState) {
	form.Slider(c, form.SliderOptions{
		State:       form.SliderState{Value: &st.lockedVolume},
		Label:       "Locked level",
		Description: "Set by your administrator.",
		Min:         0,
		Max:         100,
		Step:        10,
		Disabled:    true,
	})
}
