package gallery

// formsPageState keeps the shared form specimen state and fills its defaults
// once, so every forms page edits the same caller-owned values.
func formsPageState(s *State) *formsState {
	st := StateOf[formsState](s, "forms")
	st.initState()
	return st
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "form", Title: "Form", Group: "Forms", Order: 10,
		Description: "The container that groups fields and owns their top-level error.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "A labelled field group with supporting copy.", Demo: formsExampleDemo(formsFormExampleDefault)},
			{Title: "Invalid", Description: "The form reports one message above its fields.", Demo: formsExampleDemo(formsFormExampleInvalid)},
			{Title: "With actions", Description: "Submit and secondary actions sit under the fields.", Demo: formsExampleDemo(formsFormExampleActions)},
		},
		Code: `form.Form(c, form.FormOptions{Label: "Account", Invalid: invalid}, func() {
    form.FieldGroup(c, func() {
        form.TextField(c, form.TextFieldOptions{
            Value: &email, Label: "Email", Required: true, FullWidth: true,
        })
    })
    form.FieldActions(c, func() {
        mygoui.Button(c, mygoui.ButtonOptions{Label: "Submit form", OnClick: submit})
    })
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Accessible name for the form container."},
			{"Invalid", "bool", "false", "Shows the localized invalid message when Error is empty."},
			{"Error", "string", "", "Overrides the localized message at the top of the form."},
		},
		Tokens: []string{"Danger", "Foreground", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "fieldset", Title: "Fieldset", Group: "Forms", Order: 20,
		Description: "A legend and supporting copy around a group of related controls.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "A legend and supporting copy group related check boxes.", Demo: formsExampleDemo(formsFieldsetExampleDefault)},
			{Title: "Disabled", Description: "Disabling the fieldset disables every descendant.", Demo: formsExampleDemo(formsFieldsetExampleDisabled)},
		},
		Code: `form.Fieldset(c, form.FieldsetOptions{
    Legend: "Notifications", Description: "Choose what we send.", FullWidth: true,
}, func() {
    form.Checkbox(c, form.CheckboxOptions{
        Checked: &updates, Label: "Product updates",
    })
})`,
		Params: []ComponentParam{
			{"Legend", "string", "", "Group title above the controls."},
			{"Description", "string", "", "Supporting copy under the legend."},
			{"Disabled", "bool", "false", "Disables every descendant control."},
			{"FullWidth", "bool", "false", "Stretches the fieldset across its row."},
		},
		Tokens: []string{"Foreground", "Muted", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "label", Title: "Label", Group: "Forms", Order: 30,
		Description: "The field label, with required, invalid and disabled states.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "The plain field label.", Demo: formsExampleDemo(formsLabelExampleDefault)},
			{Title: "Required", Description: "The danger-colored required marker follows the text.", Demo: formsExampleDemo(formsLabelExampleRequired)},
			{Title: "Invalid", Description: "The label takes the danger color.", Demo: formsExampleDemo(formsLabelExampleInvalid)},
			{Title: "Disabled", Description: "The label dims and drops its focus target.", Demo: formsExampleDemo(formsLabelExampleDisabled)},
		},
		Code: `form.Label(c, "Email address", form.LabelOptions{Required: true})
form.Field(c, form.FieldOptions{
    Label: "Workspace", Description: "Lowercase letters and dashes.",
}, func() {
    form.Input(c, &name, form.InputOptions{FullWidth: true})
})`,
		Params: []ComponentParam{
			{"Required", "bool", "false", "Appends the danger-colored required marker."},
			{"Invalid", "bool", "false", "Paints the label with the danger color."},
			{"Disabled", "bool", "false", "Dims the label and removes target focus."},
			{"Target", "*ui.Handle", "nil", "Control that receives focus when the label is clicked."},
		},
		Tokens: []string{"Foreground", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "description", Title: "Description", Group: "Forms", Order: 40,
		Description: "Muted help text below a label, alone or through a field.",
		Examples: []ComponentExample{
			{Title: "Standalone", Description: "Muted help copy follows its label.", Demo: formsExampleDemo(formsDescriptionExampleStandalone)},
			{Title: "With field", Description: "The field supplies the label and supporting copy.", Demo: formsExampleDemo(formsDescriptionExampleField)},
			{Title: "Hidden while invalid", Description: "The error replaces the description while the field is invalid.", Demo: formsExampleDemo(formsDescriptionExampleHidden)},
		},
		Code: `form.Label(c, "API key", form.LabelOptions{})
form.Description(c, "Rotate keys every ninety days.")

form.Field(c, form.FieldOptions{
    Label: "Slug", Description: "Appears in the URL.",
}, func() { form.Input(c, &slug, form.InputOptions{FullWidth: true}) })`,
		Params: []ComponentParam{
			{"text", "string", "required", "Help copy rendered in the muted color."},
			{"Field.Description", "string", "", "The same copy through a field option, hidden while invalid."},
			{"Field.Invalid", "bool", "false", "Hides the field's description while showing its error."},
		},
		Tokens: []string{"Muted", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "fielderror", Title: "FieldError", Group: "Forms", Order: 50,
		Description: "Validation copy in the danger color, announced as a status.",
		Examples: []ComponentExample{
			{Title: "Standalone", Description: "Validation copy sits under its label with status semantics.", Demo: formsExampleDemo(formsFieldErrorExampleStandalone)},
			{Title: "Through a field", Description: "The field root supplies the invalid state and message.", Demo: formsExampleDemo(formsFieldErrorExampleField)},
		},
		Code: `form.FieldError(c, "Enter an address containing @.")

form.Field(c, form.FieldOptions{
    Label: "Username", Invalid: true, Error: "This username is already taken.",
}, func() { form.Input(c, &username, form.InputOptions{FullWidth: true}) })`,
		Params: []ComponentParam{
			{"text", "string", "required", "Validation copy; empty text renders nothing."},
			{"Field.Error", "string", "", "Error text set through a field root."},
			{"Field.Invalid", "bool", "false", "Shows the localized invalid message when Error is empty."},
		},
		Tokens: []string{"Danger", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "errormessage", Title: "ErrorMessage", Group: "Forms", Order: 60,
		Description: "Low-level error text for custom layouts and form-level failures.",
		Examples: []ComponentExample{
			{Title: "Custom layout", Description: "The low-level message follows a caller-built field.", Demo: formsExampleDemo(formsErrorMessageExampleCustom)},
			{Title: "Form-level", Description: "One failure message reports the whole form.", Demo: formsExampleDemo(formsErrorMessageExampleForm)},
		},
		Code: `form.ErrorMessage(c, "Use at least twelve characters.")

form.Form(c, form.FormOptions{Label: "Profile", Invalid: invalid}, func() {
    form.Input(c, &name, form.InputOptions{FullWidth: true})
})`,
		Params: []ComponentParam{
			{"text", "string", "required", "Low-level error copy; empty text renders nothing."},
			{"Form.Error", "string", "", "Form-level error shown above its fields."},
			{"Form.Invalid", "bool", "false", "Uses the localized message when Error is empty."},
		},
		Tokens: []string{"Danger", "Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "input", Title: "Input", Group: "Forms", Order: 70,
		Description: "A single-line editor, bare or wrapped in a field.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "A bare editor with an accessible name.", Demo: formsExampleDemo(formsInputExampleDefault)},
			{Title: "With description", Description: "The field supplies the label, required marker and help copy.", Demo: formsExampleDemo(formsInputExampleDescription)},
			{Title: "Invalid", Description: "The danger border pairs with the error message.", Demo: formsExampleDemo(formsInputExampleInvalid)},
			{Title: "Password", Description: "The editor masks the typed text.", Demo: formsExampleDemo(formsInputExamplePassword)},
			{Title: "Disabled", Description: "The surface dims and stops editing.", Demo: formsExampleDemo(formsInputExampleDisabled)},
		},
		Code: `form.Input(c, &email, form.InputOptions{
    Label: "Email", Description: "We send account notices here.",
    Required: true, FullWidth: true,
})

form.Input(c, &query, form.InputOptions{
    AriaLabel: "Search projects", Placeholder: "Search projects",
})`,
		Params: []ComponentParam{
			{"Variant", "form.ControlVariant", "ControlPrimary", "Primary field surface or secondary fill."},
			{"Placeholder", "string", "", "Hint shown while the field is empty."},
			{"Required", "bool", "false", "Marks the field label as required."},
			{"Disabled", "bool", "false", "Stops editing and dims the surface."},
			{"Invalid", "bool", "false", "Paints the danger border and error."},
			{"Password", "bool", "false", "Masks the text as a password."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldBorderHover", "FieldPlaceholder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "textarea", Title: "TextArea", Group: "Forms", Order: 80,
		Description: "A multiline editor that grows to its line limit, then scrolls.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "The editor grows to its line limit, then scrolls.", Demo: formsExampleDemo(formsTextAreaExampleDefault)},
			{Title: "Locked", Description: "Text stays selectable; typing cannot change it.", Demo: formsExampleDemo(formsTextAreaExampleReadOnly)},
			{Title: "Invalid", Description: "The danger border pairs with the error message.", Demo: formsExampleDemo(formsTextAreaExampleInvalid)},
		},
		Code: `form.TextArea(c, &bio, form.TextAreaOptions{
    Label: "Bio", Description: "Grows to six lines, then scrolls.",
    MinLines: 3, MaxLines: 6, FullWidth: true,
})`,
		Params: []ComponentParam{
			{"MinLines", "int", "1", "Minimum visible line count."},
			{"MaxLines", "int", "0", "Grows to this many lines, then scrolls."},
			{"Required", "bool", "false", "Marks the field label as required."},
			{"Disabled", "bool", "false", "Stops editing and dims the surface."},
			{"Invalid", "bool", "false", "Paints the danger border and error."},
			{"FullWidth", "bool", "false", "Stretches the editor across its row."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldHover", "FieldPlaceholder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "textfield", Title: "TextField", Group: "Forms", Order: 90,
		Description: "A labelled, described and validated single-line field in one call.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "Label and editor in one call.", Demo: formsExampleDemo(formsTextFieldExampleDefault)},
			{Title: "With description", Description: "Supporting copy sits between the label and the editor.", Demo: formsExampleDemo(formsTextFieldExampleDescription)},
			{Title: "Invalid", Description: "The field supplies the danger border and error.", Demo: formsExampleDemo(formsTextFieldExampleInvalid)},
			{Title: "Password", Description: "The editor masks the value and marks it required.", Demo: formsExampleDemo(formsTextFieldExamplePassword)},
			{Title: "Disabled", Description: "The caller-owned value cannot change.", Demo: formsExampleDemo(formsTextFieldExampleDisabled)},
		},
		Code: `form.TextField(c, form.TextFieldOptions{
    Value: &email, Label: "Email", Description: "We only use this for account notices.",
    Required: true, Invalid: invalid, Error: "Enter an address containing @.",
    FullWidth: true,
})`,
		Params: []ComponentParam{
			{"Value", "*string", "required", "Caller-owned text the field edits in place."},
			{"Label", "string", "", "Visible label and accessible name."},
			{"Description", "string", "", "Help text hidden while the field is invalid."},
			{"Required", "bool", "false", "Marks the label as required."},
			{"Invalid", "bool", "false", "Paints the danger border and error."},
			{"Password", "bool", "false", "Masks the text as a password."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldBorderHover", "FieldPlaceholder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "inputgroup", Title: "InputGroup", Group: "Forms", Order: 100,
		Description: "One editor framed by inline or full-width addon regions.",
		Examples: []ComponentExample{
			{Title: "Prefix and suffix", Description: "Inline addon regions frame the editor.", Demo: formsExampleDemo(formsInputGroupExamplePrefixSuffix)},
			{Title: "Block addons", Description: "Full-width rows sit above and below the editor.", Demo: formsExampleDemo(formsInputGroupExampleBlock)},
			{Title: "Invalid", Description: "The danger shell pairs with the error message.", Demo: formsExampleDemo(formsInputGroupExampleInvalid)},
		},
		Code: `form.InputGroup(c, &price, form.InputGroupOptions{
    Label: "Price", FullWidth: true,
    Prefix: func() { muted("$") },
    Suffix: func() { maxPriceButton() },
})`,
		Params: []ComponentParam{
			{"Prefix", "func()", "nil", "Inline content before the editor."},
			{"Suffix", "func()", "nil", "Inline content after the editor."},
			{"BlockPrefix", "func()", "nil", "Full-width row above the input."},
			{"BlockSuffix", "func()", "nil", "Full-width row below the input."},
			{"Variant", "form.ControlVariant", "ControlPrimary", "Primary field surface or secondary fill."},
			{"Invalid", "bool", "false", "Paints the danger border and error."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldPlaceholder", "DefaultHover", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "inputotp", Title: "InputOTP", Group: "Forms", Order: 110,
		Description: "One-time code slots painted over a single native editor.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "One native editor drives the six painted slots.", Demo: formsExampleDemo(formsOTPExampleDefault)},
			{Title: "Invalid", Description: "Every slot takes the danger border and the error shows.", Demo: formsExampleDemo(formsOTPExampleInvalid)},
			{Title: "Disabled", Description: "The slots dim and stop accepting input.", Demo: formsExampleDemo(formsOTPExampleDisabled)},
		},
		Code: `form.InputOTP(c, &otp, form.OTPOptions{
    Label: "Verification code", Length: 6,
    OnChange:   func(value string) { note(len(value)) },
    OnComplete: func(value string) { verify(value) },
})`,
		Params: []ComponentParam{
			{"Length", "int", "6", "Number of digit slots, capped at twelve."},
			{"Filter", "func(rune) bool", "printable ASCII", "Accepts or rejects each typed character."},
			{"Invalid", "bool", "false", "Paints every slot with the danger border."},
			{"Disabled", "bool", "false", "Stops editing and dims the slots."},
			{"OnChange", "func(string)", "nil", "Reports the accepted value on every edit."},
			{"OnComplete", "func(string)", "nil", "Runs once each time the value fills Length."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldFocus", "AccentSoft", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "searchfield", Title: "SearchField", Group: "Forms", Order: 120,
		Description: "A native search editor with an in-field clear action.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "The in-field clear action empties the query.", Demo: formsExampleDemo(formsSearchFieldExampleDefault)},
			{Title: "Invalid", Description: "The danger border pairs with the error message.", Demo: formsExampleDemo(formsSearchFieldExampleInvalid)},
			{Title: "Disabled", Description: "Editing stops and the clear action is disabled.", Demo: formsExampleDemo(formsSearchFieldExampleDisabled)},
		},
		Code: `form.SearchField(c, &query, form.SearchFieldOptions{
    Label: "Search components", Placeholder: "Search by name",
    FullWidth: true, OnClear: func() { resetResults() },
})`,
		Params: []ComponentParam{
			{"Placeholder", "string", "", "Hint shown while the query is empty."},
			{"OnClear", "func()", "nil", "Runs after the clear button or Escape empties the value."},
			{"OnSubmit", "func()", "nil", "Runs when Enter is pressed."},
			{"Invalid", "bool", "false", "Paints the danger border and error."},
			{"Disabled", "bool", "false", "Stops editing; the clear action is disabled."},
			{"FullWidth", "bool", "false", "Stretches the field across its row."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldPlaceholder", "DefaultHover", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "numberfield", Title: "NumberField", Group: "Forms", Order: 130,
		Description: "A numeric editor between stepper buttons, with optional bounds.",
		Examples: []ComponentExample{
			{Title: "Bounded", Description: "Typed values and the steppers share one clamped number.", Demo: formsExampleDemo(formsNumberFieldExampleBounded)},
			{Title: "Disabled", Description: "Stepping and typed changes stop.", Demo: formsExampleDemo(formsNumberFieldExampleDisabled)},
		},
		Code: `form.NumberField(c, &seats, form.NumberFieldOptions{
    Label: "Seats", Description: "Between one and twelve.",
    Min: new(float64(1)), Max: new(float64(12)), Step: 1,
})`,
		Params: []ComponentParam{
			{"Min", "*float64", "nil", "Lower bound; typed values clamp to it."},
			{"Max", "*float64", "nil", "Upper bound; typed values clamp to it."},
			{"Step", "float64", "1", "Increment of the stepper buttons and arrow keys."},
			{"Format", "func(float64) string", "step decimals", "Replaces the displayed number."},
			{"Disabled", "bool", "false", "Stops stepping and editable text."},
			{"Invalid", "bool", "false", "Paints the danger border and error."},
		},
		Tokens: []string{"FieldBackground", "FieldBorder", "FieldForeground", "FieldPlaceholder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "checkbox", Title: "Checkbox", Group: "Forms", Order: 140,
		Description: "A single caller-owned check box with mixed and invalid states.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "A single caller-owned check box.", Demo: formsExampleDemo(formsCheckboxExampleDefault)},
			{Title: "Secondary", Description: "The flat secondary control face.", Demo: formsExampleDemo(formsCheckboxExampleSecondary)},
			{Title: "With description", Description: "Supporting copy indents under the label.", Demo: formsExampleDemo(formsCheckboxExampleDescription)},
			{Title: "Mixed", Description: "The indeterminate mark shows a partial selection.", Demo: formsExampleDemo(formsCheckboxExampleMixed)},
			{Title: "Invalid", Description: "The danger control pairs with the error message.", Demo: formsExampleDemo(formsCheckboxExampleInvalid)},
			{Title: "Disabled", Description: "Activation is removed and the row dims.", Demo: formsExampleDemo(formsCheckboxExampleDisabled)},
		},
		Code: `form.Checkbox(c, form.CheckboxOptions{
    Checked: &releaseNotes, Label: "Release notes",
    Description: "Email me when a new version ships.",
})`,
		Params: []ComponentParam{
			{"Checked", "*bool", "required", "Caller-owned value flipped by the row."},
			{"Indeterminate", "*bool", "nil", "Paints the mixed mark until activation clears it."},
			{"Variant", "form.ControlVariant", "ControlPrimary", "Primary field face or the flat secondary face."},
			{"Description", "string", "", "Supporting line indented under the label."},
			{"Required", "bool", "false", "Marks the label as required."},
			{"Invalid", "bool", "false", "Paints the danger control and error."},
			{"Disabled", "bool", "false", "Removes activation and dims the row."},
		},
		Tokens: []string{"Accent", "AccentForeground", "AccentHover", "FieldBorder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "checkboxgroup", Title: "CheckboxGroup", Group: "Forms", Order: 150,
		Description: "Related check boxes with an optional mixed select-all row.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "Related check boxes share one label.", Demo: formsExampleDemo(formsCheckboxGroupExampleDefault)},
			{Title: "Select all", Description: "The select-all row derives its mixed state.", Demo: formsExampleDemo(formsCheckboxGroupExampleSelectAll)},
			{Title: "Horizontal", Description: "The choices wrap into a row.", Demo: formsExampleDemo(formsCheckboxGroupExampleHorizontal)},
		},
		Code: `form.CheckboxGroup(c, form.CheckboxGroupOptions{
    Label: "Notifications", State: &checks, SelectAllLabel: "All notifications",
}, func() {
    form.Checkbox(c, form.CheckboxOptions{Checked: &updates, Label: "Product updates"})
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Group label above the check boxes."},
			{"Description", "string", "", "Supporting line under the label."},
			{"Orientation", "form.GroupOrientation", "GroupVertical", "Vertical column or wrapping horizontal row."},
			{"SelectAllLabel", "string", "", "Adds a select-all row with the derived mixed state."},
			{"State", "*form.CheckboxGroupState", "nil", "Items backing the select-all row."},
			{"Invalid", "bool", "false", "Marks the group's error."},
		},
		Tokens: []string{"Foreground", "Muted", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "radio", Title: "Radio", Group: "Forms", Order: 160,
		Description: "One native radio button sharing a caller-owned group value.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "One choice inside its shared group.", Demo: formsExampleDemo(formsRadioExampleDefault)},
			{Title: "With description", Description: "Supporting copy indents under the choice.", Demo: formsExampleDemo(formsRadioExampleDescription)},
			{Title: "Invalid", Description: "The danger control pairs with the error message.", Demo: formsExampleDemo(formsRadioExampleInvalid)},
			{Title: "Disabled", Description: "Selection is removed and the row dims.", Demo: formsExampleDemo(formsRadioExampleDisabled)},
		},
		Code: `form.Radio(c, form.RadioOptions[string]{
    State: &plan, Value: "pro", Label: "Pro",
    Description: "Unlimited projects and priority support.",
})`,
		Params: []ComponentParam{
			{"State", "*form.RadioGroupState[T]", "required", "Caller-owned single choice shared by the group."},
			{"Value", "T", "required", "Value this radio selects."},
			{"Label", "string", "", "Visible label of the choice."},
			{"Description", "string", "", "Supporting line indented under the label."},
			{"Invalid", "bool", "false", "Paints the danger control and error."},
			{"Disabled", "bool", "false", "Removes selection and dims the row."},
		},
		Tokens: []string{"Accent", "AccentForeground", "FieldBorder", "FieldHover", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "radiogroup", Title: "RadioGroup", Group: "Forms", Order: 170,
		Description: "One tab stop and one roving choice across a set of radios.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "One tab stop and one roving choice.", Demo: formsExampleDemo(formsRadioGroupExampleDefault)},
			{Title: "Horizontal", Description: "The choices wrap into a row.", Demo: formsExampleDemo(formsRadioGroupExampleHorizontal)},
			{Title: "Invalid", Description: "The group and each radio take the danger outline.", Demo: formsExampleDemo(formsRadioGroupExampleInvalid)},
			{Title: "Disabled", Description: "The group and its choices dim together.", Demo: formsExampleDemo(formsRadioGroupExampleDisabled)},
		},
		Code: `form.RadioGroup(c, form.RadioGroupOptions{Label: "Plan", FullWidth: true}, func() {
    form.Radio(c, form.RadioOptions[string]{State: &plan, Value: "free", Label: "Free"})
    form.Radio(c, form.RadioOptions[string]{State: &plan, Value: "pro", Label: "Pro"})
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Group label and native group name."},
			{"Description", "string", "", "Supporting line under the label."},
			{"Orientation", "form.GroupOrientation", "GroupVertical", "Vertical column or wrapping horizontal row."},
			{"Invalid", "bool", "false", "Marks the group error; set each radio Invalid too."},
			{"Disabled", "bool", "false", "Disables the group and dims it."},
			{"Required", "bool", "false", "Marks the group label as required."},
		},
		Tokens: []string{"Foreground", "Muted", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "switch", Title: "Switch", Group: "Forms", Order: 180,
		Description: "An on/off track in three sizes with either label position.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "An on/off row with supporting copy.", Demo: formsExampleDemo(formsSwitchExampleDefault)},
			{Title: "Sizes", Description: "Small, medium and large tracks.", Demo: formsExampleDemo(formsSwitchExampleSizes)},
			{Title: "Label before", Description: "The label sits before the track.", Demo: formsExampleDemo(formsSwitchExampleLabelBefore)},
			{Title: "Invalid", Description: "The row takes the danger mark and error.", Demo: formsExampleDemo(formsSwitchExampleInvalid)},
			{Title: "Disabled", Description: "Activation is removed and the row dims.", Demo: formsExampleDemo(formsSwitchExampleDisabled)},
		},
		Code: `form.Switch(c, form.SwitchOptions{
    On: &wifi, Label: "Wi-Fi",
    Description: "Connect automatically on known networks.",
})`,
		Params: []ComponentParam{
			{"On", "*bool", "required", "Caller-owned value flipped by the row."},
			{"Size", "form.SwitchSize", "SwitchMedium", "Small, medium or large track."},
			{"LabelPosition", "form.SwitchLabelPosition", "SwitchLabelAfter", "Puts the label before or after the track."},
			{"Description", "string", "", "Supporting line indented under the label."},
			{"Invalid", "bool", "false", "Marks the row invalid and shows its error."},
			{"Disabled", "bool", "false", "Removes activation and dims the row."},
		},
		Tokens: []string{"Accent", "AccentHover", "AccentForeground", "DefaultHover"},
	})

	RegisterComponent(ComponentPage{
		ID: "switchgroup", Title: "SwitchGroup", Group: "Forms", Order: 190,
		Description: "Related switches under one label and description.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "Related switches share one label and description.", Demo: formsExampleDemo(formsSwitchGroupExampleDefault)},
			{Title: "Horizontal", Description: "The switches wrap into a row.", Demo: formsExampleDemo(formsSwitchGroupExampleHorizontal)},
			{Title: "Disabled", Description: "Every switch dims and stops responding.", Demo: formsExampleDemo(formsSwitchGroupExampleDisabled)},
		},
		Code: `form.SwitchGroup(c, form.SwitchGroupOptions{
    Label: "Privacy", Description: "Applies to every workspace.", FullWidth: true,
}, func() {
    form.Switch(c, form.SwitchOptions{On: &analytics, Label: "Usage analytics"})
})`,
		Params: []ComponentParam{
			{"Label", "string", "", "Group label above the switches."},
			{"Description", "string", "", "Supporting line under the label."},
			{"Orientation", "form.GroupOrientation", "GroupVertical", "Vertical column or wrapping horizontal row."},
			{"Invalid", "bool", "false", "Marks the group error."},
			{"Disabled", "bool", "false", "Disables every switch and dims the group."},
			{"FullWidth", "bool", "false", "Stretches the group across its row."},
		},
		Tokens: []string{"Foreground", "Muted", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "controlfield", Title: "ControlField", Group: "Forms", Order: 200,
		Description: "A whole label-and-description row that toggles one trailing indicator.",
		Examples: []ComponentExample{
			{Title: "Switch", Description: "The switch trails the whole activation row.", Demo: formsExampleDemo(formsControlFieldExampleSwitch)},
			{Title: "Checkbox", Description: "The check box trails the whole activation row.", Demo: formsExampleDemo(formsControlFieldExampleCheckbox)},
			{Title: "Invalid", Description: "The danger indicator pairs with the error.", Demo: formsExampleDemo(formsControlFieldExampleInvalid)},
		},
		Code: `form.ControlField(c, form.ControlFieldOptions{
    Checked: &autoSave, Label: "Auto-save",
    Description: "Save drafts a moment after you stop typing.",
    Kind: form.ControlFieldSwitch,
})`,
		Params: []ComponentParam{
			{"Kind", "form.ControlFieldKind", "ControlFieldSwitch", "Switch or check box indicator."},
			{"Checked", "*bool", "required", "Caller-owned value flipped by the whole row."},
			{"Label", "string", "", "Row title and accessible name."},
			{"Description", "string", "", "Supporting line under the title."},
			{"Invalid", "bool", "false", "Paints the danger indicator and error."},
			{"Disabled", "bool", "false", "Removes activation and dims the row."},
		},
		Tokens: []string{"Accent", "AccentForeground", "Foreground", "Muted", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "slider", Title: "Slider", Group: "Forms", Order: 210,
		Description: "Single, range and vertical tracks with native drag and keys.",
		Examples: []ComponentExample{
			{Title: "Default", Description: "A single thumb with labelled marks.", Demo: formsExampleDemo(formsSliderExampleDefault)},
			{Title: "Range", Description: "Two thumbs stay ordered; drag either end.", Demo: formsExampleDemo(formsSliderExampleRange)},
			{Title: "Vertical", Description: "A vertical single-thumb track.", Demo: formsExampleDemo(formsSliderExampleVertical)},
			{Title: "Disabled", Description: "Dragging and keyboard changes stop.", Demo: formsExampleDemo(formsSliderExampleDisabled)},
		},
		Code: `form.Slider(c, form.SliderOptions{
    State: form.SliderState{Value: &volume},
    Label: "Volume", Min: 0, Max: 100, Step: 5,
    Format: func(v float64) string { return fmt.Sprintf("%.0f%%", v) },
})`,
		Params: []ComponentParam{
			{"Min", "float64", "0", "Lower end of the track."},
			{"Max", "float64", "100", "Upper end of the track."},
			{"Step", "float64", "0", "Snap increment; zero keeps the native step."},
			{"Marks", "[]form.SliderMark", "nil", "Labelled ticks under a horizontal track."},
			{"Vertical", "bool", "false", "Draws a vertical single-thumb track."},
			{"Disabled", "bool", "false", "Stops dragging and keyboard changes."},
		},
		Tokens: []string{"Accent", "AccentForeground", "Default", "Foreground", "Muted"},
	})
}
