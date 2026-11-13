package gallery

import (
	"runtime"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// formInputsState is the caller-owned state behind the form input pages.
type formInputsState struct {
	ready bool

	password, token string
	tokenRevealed   bool

	tags, reviewers []string

	priceLow, priceHigh float64
	hoursLow, hoursHigh float64

	backups, alerts bool
	plan            string
	planState       form.RadioGroupState[string]

	output, input string
	browseError   string

	env     []form.KeyValue
	headers []form.KeyValue
	secrets []form.KeyValue
	nextID  int

	date, plate, phone string
	strength           string
}

func formInputsStateOf(s *State) *formInputsState {
	st := StateOf[formInputsState](s, "form-inputs")
	if !st.ready {
		st.ready = true
		st.token = "ghp_2f9c41d07ab3"
		st.tags = []string{"native", "desktop"}
		st.reviewers = []string{"ada", "grace", "linus"}
		st.priceLow, st.priceHigh = 120, 480
		st.hoursLow, st.hoursHigh = 9, 17
		st.backups = true
		st.plan = "team"
		st.planState = form.RadioGroupState[string]{Selected: &st.plan}
		st.output = "/Users/ada/Projects/atlas/build/release/darwin-arm64/atlas-2026.11.dmg"
		st.input = "/etc/hosts"
		st.env = []form.KeyValue{{Key: "GOOS", Value: "darwin"}, {Key: "GOARCH", Value: "arm64"}}
		st.headers = []form.KeyValue{
			{Key: "Accept", Value: "application/json"},
			{Key: "Authorization", Value: "Bearer …"},
			{Key: "Accept", Value: "text/plain"},
		}
		st.plate = "AB1234"
		st.secrets = []form.KeyValue{
			{ID: "row-1", Key: "API_URL", Value: "https://api.example.com"},
			{ID: "row-2", Key: "API_TOKEN", Value: "ghp_2f9c41d07ab3", Secret: true},
		}
		st.nextID = 3
	}
	return st
}

func passwordFieldBasicDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.PasswordField(c, &st.password, form.PasswordFieldOptions{
			Label: "Password", Placeholder: "At least 12 characters", FullWidth: true,
		})
	})
}

func passwordFieldRevealedDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.PasswordField(c, &st.token, form.PasswordFieldOptions{
			Label: "Access token", Revealed: &st.tokenRevealed, ReadOnly: true, FullWidth: true,
		})
		formsMutedText(c, "Revealed: "+strconv.FormatBool(st.tokenRevealed))
	})
}

func passwordFieldInvalidDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.PasswordField(c, &st.password, form.PasswordFieldOptions{
			Label: "New password", Error: "Use at least 12 characters.", FullWidth: true,
		})
	})
}

// passwordStrengthDemo rates by length and variety; the policy is the caller's.
func passwordStrengthDemo(value string) (int, string) {
	if value == "" {
		return 0, ""
	}
	classes := 0
	for _, has := range []func(rune) bool{unicode.IsLower, unicode.IsUpper, unicode.IsDigit, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }} {
		if strings.IndexFunc(value, has) >= 0 {
			classes++
		}
	}
	level := min(len(value)/4, 2) + min(classes/2, 2)
	level = min(max(level, 1), form.PasswordStrengthLevels)
	return level, []string{"", "Weak", "Fair", "Good", "Strong"}[level]
}

func passwordFieldStrengthDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.PasswordField(c, &st.strength, form.PasswordFieldOptions{
			Label: "New password", FullWidth: true, Strength: passwordStrengthDemo,
		})
	})
}

func tagsInputBasicDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.TagsInput(c, &st.tags, form.TagsInputOptions{
			Label: "Topics", Placeholder: "Add a topic", FullWidth: true,
		})
		formsMutedText(c, strings.Join(st.tags, ", "))
	})
}

func tagsInputMaxDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.TagsInput(c, &st.reviewers, form.TagsInputOptions{
			Label: "Reviewers", Description: strconv.Itoa(len(st.reviewers)) + " of 4",
			Max: 4, FullWidth: true,
		})
	})
}

func tagsInputPasteDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.TagsInput(c, &st.tags, form.TagsInputOptions{
			Label: "Topics", Placeholder: "Paste a list", FullWidth: true, KeepDraftOnBlur: true,
		})
	})
}

func rangeSliderBasicDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.RangeSlider(c, &st.priceLow, &st.priceHigh, form.RangeSliderOptions{
			Label: "Price", Max: 1000, Step: 10,
			Format: func(v float64) string { return "$" + strconv.FormatFloat(v, 'f', 0, 64) },
		})
	})
}

func rangeSliderGapDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.RangeSlider(c, &st.hoursLow, &st.hoursHigh, form.RangeSliderOptions{
			Label: "Working hours", Max: 24, Step: 1, MinGap: 4,
			Marks: []form.SliderMark{
				{Value: 0, Label: "0"}, {Value: 6, Label: "6"}, {Value: 12, Label: "12"},
				{Value: 18, Label: "18"}, {Value: 24, Label: "24"},
			},
		})
	})
}

func checkboxCardDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.CheckboxGroup(c, form.CheckboxGroupOptions{Label: "Add-ons"}, func() {
			form.CheckboxCard(c, form.CheckboxCardOptions{
				Checked: &st.backups, Icon: lucide.DatabaseBackup(), FullWidth: true,
				Title: "Nightly backups", Description: "Snapshots kept for 30 days.",
			})
			form.CheckboxCard(c, form.CheckboxCardOptions{
				Checked: &st.alerts, Icon: lucide.BellRing(), FullWidth: true,
				Title: "Usage alerts", Description: "Email when a quota passes 80%.",
			})
			form.CheckboxCard(c, form.CheckboxCardOptions{
				Icon: lucide.ShieldCheck(), FullWidth: true, Disabled: true,
				Title: "Audit log", Description: "Included with the Scale plan.",
			})
		})
	})
}

func choiceCardContentDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	plans := []struct{ value, title, price string }{
		{"free", "Free", "$0"},
		{"team", "Team", "$12"},
		{"scale", "Scale", "$48"},
	}
	formsFieldStack(c, func() {
		form.RadioGroup(c, form.RadioGroupOptions{Label: "Billing", FullWidth: true}, func() {
			for _, p := range plans {
				form.RadioCard(c, form.RadioCardOptions[string]{
					State: &st.planState, Value: p.value, FullWidth: true, Title: p.title,
					HideIndicator: true,
					Content: func() {
						formsMutedText(c, p.price+" per seat, billed monthly")
					},
				})
			}
		})
	})
}

func radioCardDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	plans := []struct{ value, title, description string }{
		{"free", "Free", "One project, community support."},
		{"team", "Team", "Unlimited projects and shared roles."},
		{"scale", "Scale", "Audit log, SSO and a support line."},
	}
	formsFieldStack(c, func() {
		form.RadioGroup(c, form.RadioGroupOptions{Label: "Plan", FullWidth: true}, func() {
			for _, p := range plans {
				form.RadioCard(c, form.RadioCardOptions[string]{
					State: &st.planState, Value: p.value, FullWidth: true,
					Title: p.title, Description: p.description,
				})
			}
		})
	})
}

func radioCardHorizontalDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	form.RadioGroup(c, form.RadioGroupOptions{
		Label: "Plan", Orientation: form.GroupHorizontal,
	}, func() {
		for _, title := range []string{"free", "team", "scale"} {
			form.RadioCard(c, form.RadioCardOptions[string]{
				State: &st.planState, Value: title, Title: strings.ToUpper(title[:1]) + title[1:],
			}).Width(theme.Scale(c, 140))
		}
	})
}

func pathFieldBrowseDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	o := form.PathFieldOptions{Label: "Input file", Placeholder: "Choose a file", FullWidth: true}
	// Browsers have no native dialog, and a view built without a window has
	// nothing to publish the result through.
	if window := s.Window; window != nil && runtime.GOOS != "js" {
		o.OnBrowse = func() {
			go func() {
				paths, err := form.PickFiles(form.PickFilesOptions{Parent: window})
				window.Update(func() {
					st.browseError = ""
					if err != nil {
						st.browseError = err.Error()
					}
					if len(paths) > 0 {
						st.input = paths[0]
					}
				})
			}()
		}
	}
	formsFieldStack(c, func() {
		form.PathField(c, &st.input, o)
		form.FieldError(c, st.browseError)
	})
}

func pathFieldLongDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.PathField(c, &st.output, form.PathFieldOptions{
			Label: "Output", ReadOnly: true, BrowseLabel: "Change", FullWidth: true,
			OnBrowse: func() { st.output = "/Volumes/Releases/atlas-2026.11.dmg" },
		})
	})
}

func keyValueEditorEnvDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.KeyValueEditor(c, form.KeyValueEditorOptions{
			AriaLabel: "Environment variables", Pairs: st.env,
			KeyPlaceholder: "Name", AddLabel: "Add variable",
			OnChange: func(i int, p form.KeyValue) { st.env[i] = p },
			OnAdd:    func() { st.env = append(st.env, form.KeyValue{}) },
			OnRemove: func(i int) { st.env = slices.Delete(st.env, i, i+1) },
		})
	})
}

func keyValueEditorDuplicatesDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.KeyValueEditor(c, form.KeyValueEditorOptions{
			AriaLabel: "Request headers", Pairs: st.headers,
			KeyPlaceholder: "Header", MarkDuplicateKeys: true,
			OnChange: func(i int, p form.KeyValue) { st.headers[i] = p },
			OnAdd:    func() { st.headers = append(st.headers, form.KeyValue{}) },
			OnRemove: func(i int) { st.headers = slices.Delete(st.headers, i, i+1) },
		})
	})
}

func keyValueEditorSecretsDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	newID := func() string {
		st.nextID++
		return "row-" + strconv.Itoa(st.nextID-1)
	}
	formsFieldStack(c, func() {
		form.KeyValueEditor(c, form.KeyValueEditorOptions{
			AriaLabel: "Settings", Pairs: st.secrets, ShowHeaders: true,
			KeyPlaceholder: "Name", AddLabel: "Add setting",
			OnChange: func(i int, p form.KeyValue) { st.secrets[i] = p },
			OnAdd:    func() { st.secrets = append(st.secrets, form.KeyValue{ID: newID()}) },
			OnRemove: func(i int) { st.secrets = slices.Delete(st.secrets, i, i+1) },
			OnPaste: func(pairs []form.KeyValue) {
				for _, p := range pairs {
					p.ID = newID()
					st.secrets = append(st.secrets, p)
				}
			},
		})
	})
}

func maskInputPhoneDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.MaskInput(c, &st.phone, form.MaskInputOptions{
			Mask: "(999) 999-9999", Label: "Phone", FullWidth: true,
		})
		formsMutedText(c, "Value: "+strconv.Quote(st.phone))
	})
}

func maskInputDateDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.MaskInput(c, &st.date, form.MaskInputOptions{
			Mask: "99/99/9999", Label: "Date of birth", Placeholder: "DD/MM/YYYY", FullWidth: true,
		})
		formsMutedText(c, "Value: "+strconv.Quote(st.date))
	})
}

func maskInputPlateDemo(c *ui.Context, s *State) {
	st := formInputsStateOf(s)
	formsFieldStack(c, func() {
		form.MaskInput(c, &st.plate, form.MaskInputOptions{
			Mask: "AA-9999", Label: "Plate", Placeholder: "AB-1234", FullWidth: true,
		})
		formsMutedText(c, "Value: "+strconv.Quote(st.plate))
	})
}

func init() {
	fieldParams := []ComponentParam{
		{"Label", "string", "\"\"", "Visible label; also the accessible name."},
		{"Description", "string", "\"\"", "Help text under the field, hidden while invalid."},
		{"Error", "string", "\"\"", "Validation message; implies Invalid."},
		{"Placeholder", "string", "\"\"", "Shown while the value is empty."},
		{"AriaLabel", "string", "\"\"", "Accessible name without a visible label."},
		{"Disabled", "bool", "false", "Dims the field and stops input."},
		{"FullWidth", "bool", "false", "Fills the parent instead of the default width."},
	}

	RegisterComponent(ComponentPage{
		ID: "passwordfield", Title: "PasswordField", Group: "Forms", Order: 280,
		Description: "A password input with a button that shows or conceals the characters.",
		Examples: []ComponentExample{
			{Title: "Reveal toggle", Demo: passwordFieldBasicDemo},
			{Title: "Caller-owned reveal", Demo: passwordFieldRevealedDemo},
			{Title: "Strength", Description: "The caller rates the value; the field sets no policy.",
				Demo: passwordFieldStrengthDemo},
			{Title: "Invalid", Demo: passwordFieldInvalidDemo},
		},
		Code: `var password string // kept in application state
form.PasswordField(c, &password, form.PasswordFieldOptions{
    Label:    "Password",
    OnSubmit: submit,
})`,
		Params: append(slices.Clone(fieldParams),
			ComponentParam{"Revealed", "*bool", "nil", "Caller-owned reveal flag; nil keeps it per window."},
			ComponentParam{"KeepRevealed", "bool", "false", "Keeps the characters showing after the focus leaves."},
			ComponentParam{"Strength", "func(string) (int, string)", "nil", "Rates the value 0 to 4 with a hint; draws a bar under the field."},
			ComponentParam{"ReadOnly", "bool", "false", "Shows the value without letting it change."},
			ComponentParam{"OnChange", "func()", "nil", "Runs after typing changes the value."},
			ComponentParam{"OnSubmit", "func()", "nil", "Runs on Enter."},
		),
		Tokens: []string{"Input", "Ring", "Accent", "MutedForeground", "Destructive"},
	})

	RegisterComponent(ComponentPage{
		ID: "tagsinput", Title: "TagsInput", Group: "Forms", Order: 290,
		Description: "Free text committed as removable tags on Enter or a comma.",
		Examples: []ComponentExample{
			{Title: "Tags", Description: "A duplicate stays in the editor, flagged. Left from the empty editor walks the tags.",
				Demo: tagsInputBasicDemo},
			{Title: "Paste", Description: "A list splits at commas, tabs and line breaks; the draft stays on blur.",
				Demo: tagsInputPasteDemo},
			{Title: "Maximum", Demo: tagsInputMaxDemo},
		},
		Code: `var tags []string // kept in application state
form.TagsInput(c, &tags, form.TagsInputOptions{
    Label:       "Topics",
    Placeholder: "Add a topic",
    Max:         8,
    OnChange:    save,
})`,
		Params: append(slices.Clone(fieldParams),
			ComponentParam{"Max", "int", "0", "Largest number of tags; 0 is unlimited."},
			ComponentParam{"AllowDuplicates", "bool", "false", "Accepts a tag that is already present."},
			ComponentParam{"KeepDraftOnBlur", "bool", "false", "Leaves the draft in the editor instead of committing it when the focus leaves."},
			ComponentParam{"OnChange", "func()", "nil", "Runs after a tag is added or removed."},
		),
		Tokens: []string{"Input", "Ring", "Secondary", "SecondaryForeground", "Accent"},
	})

	RegisterComponent(ComponentPage{
		ID: "rangeslider", Title: "RangeSlider", Group: "Forms", Order: 300,
		Description: "A two-thumb slider over a low and a high value.",
		Examples: []ComponentExample{
			{Title: "Range", Demo: rangeSliderBasicDemo},
			{Title: "Minimum gap", Description: "The thumbs stay four hours apart.",
				Demo: rangeSliderGapDemo},
		},
		Code: `var low, high float64 // kept in application state
form.RangeSlider(c, &low, &high, form.RangeSliderOptions{
    Label:    "Price",
    Max:      1000,
    Step:     10,
    MinGap:   50,
    OnChange: save,
})`,
		Params: []ComponentParam{
			{"Label", "string", "\"\"", "Visible label; names both thumbs."},
			{"Min", "float64", "0", "Low bound."},
			{"Max", "float64", "100", "High bound; 100 when Min and Max are both zero."},
			{"Step", "float64", "0", "Snaps drags and keys; 0 leaves the value free."},
			{"MinGap", "float64", "0", "Least distance between the thumbs."},
			{"Marks", "[]form.SliderMark", "nil", "Labelled ticks under the track."},
			{"HideOutput", "bool", "false", "Hides the formatted range."},
			{"Format", "func(float64) string", "nil", "Formats the output values."},
			{"Disabled", "bool", "false", "Dims the slider and stops input."},
			{"OnChange", "func()", "nil", "Runs after either value changes."},
		},
		Tokens: []string{"Primary", "Input", "Ring", "Foreground", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "checkboxcard", Title: "CheckboxCard", Group: "Forms", Order: 310,
		Description: "A whole card that is one check box.",
		Examples: []ComponentExample{
			{Title: "In a group", Demo: checkboxCardDemo},
		},
		Code: `form.CheckboxCard(c, form.CheckboxCardOptions{
    Checked:     &updates,
    Title:       "Product updates",
    Description: "A monthly summary of what shipped.",
})`,
		Params: []ComponentParam{
			{"Checked", "*bool", "required", "Caller-owned value; nil disables the card."},
			{"Title", "string", "required", "Card heading and accessible name."},
			{"Description", "string", "\"\"", "Supporting text, read after the title."},
			{"Icon", "*ui.SVG", "nil", "Leading icon."},
			{"Content", "func()", "nil", "Builds more under the texts; display only."},
			{"HideIndicator", "bool", "false", "Leaves out the check box; the card's border shows the choice."},
			{"Invalid", "bool", "false", "Draws the danger border."},
			{"Disabled", "bool", "false", "Dims the card and stops input."},
			{"FullWidth", "bool", "false", "Fills the parent."},
			{"OnChange", "func(bool)", "nil", "Reports the new value."},
		},
		Tokens: []string{"Input", "Primary", "PrimaryForeground", "Accent", "Ring", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "radiocard", Title: "RadioCard", Group: "Forms", Order: 320,
		Description: "A whole card that is one radio button of a group.",
		Examples: []ComponentExample{
			{Title: "In a group", Description: "One Tab stop; the arrows move the choice.",
				Demo: radioCardDemo},
			{Title: "Horizontal", Demo: radioCardHorizontalDemo},
			{Title: "Content", Description: "Content builds under the texts; HideIndicator leaves the border to show the choice.",
				Demo: choiceCardContentDemo},
		},
		Code: `var selected string // kept in application state
choice := form.RadioGroupState[string]{Selected: &selected}
form.RadioGroup(c, form.RadioGroupOptions{Label: "Plan"}, func() {
    form.RadioCard(c, form.RadioCardOptions[string]{
        State: &choice, Value: "free", Title: "Free",
    })
    form.RadioCard(c, form.RadioCardOptions[string]{
        State: &choice, Value: "team", Title: "Team",
        Description: "Unlimited projects and shared roles.",
    })
})`,
		Params: []ComponentParam{
			{"State", "*form.RadioGroupState[T]", "required", "The group's caller-owned selection."},
			{"Value", "T", "required", "The value this card selects."},
			{"Title", "string", "required", "Card heading and accessible name."},
			{"Description", "string", "\"\"", "Supporting text, read after the title."},
			{"Icon", "*ui.SVG", "nil", "Leading icon."},
			{"Content", "func()", "nil", "Builds more under the texts; display only."},
			{"HideIndicator", "bool", "false", "Leaves out the radio; the card's border shows the choice."},
			{"Invalid", "bool", "false", "Draws the danger border."},
			{"Disabled", "bool", "false", "Dims the card and stops input."},
			{"FullWidth", "bool", "false", "Fills the parent."},
			{"OnChange", "func(T)", "nil", "Reports the value chosen."},
		},
		Tokens: []string{"Input", "Primary", "PrimaryForeground", "Accent", "Ring", "Muted"},
	})

	RegisterComponent(ComponentPage{
		ID: "pathfield", Title: "PathField", Group: "Forms", Order: 330,
		Description: "A file-system path editor with a browse button the caller wires to a dialog.",
		Examples: []ComponentExample{
			{Title: "Browse", Demo: pathFieldBrowseDemo},
			{Title: "Long path", Description: "The middle is cut, and the tooltip shows it all, until the editor is focused. Drop a file on a field to set its path.",
				Demo: pathFieldLongDemo},
		},
		Code: `form.PathField(c, &destination, form.PathFieldOptions{
    Label:    "Destination",
    OnBrowse: showAll, // open a dialog, then write the chosen path
    OnChange: save,
})`,
		Params: append(slices.Clone(fieldParams),
			ComponentParam{"OnBrowse", "func()", "nil", "Runs for the browse button; nil disables it."},
			ComponentParam{"BrowseLabel", "string", "Browse", "Replaces the localized button text."},
			ComponentParam{"ReadOnly", "bool", "false", "Shows the path without letting it be typed over."},
			ComponentParam{"OnChange", "func()", "nil", "Runs after typing changes the path."},
			ComponentParam{"OnSubmit", "func()", "nil", "Runs on Enter."},
		),
		Tokens: []string{"Input", "Ring", "Accent", "Foreground", "FieldPlaceholder"},
	})

	RegisterComponent(ComponentPage{
		ID: "keyvalueeditor", Title: "KeyValueEditor", Group: "Patterns", Order: 340,
		Description: "Editable rows of key and value text, as environment variables or headers.",
		Examples: []ComponentExample{
			{Title: "Rows", Demo: keyValueEditorEnvDemo},
			{Title: "Duplicate keys", Demo: keyValueEditorDuplicatesDemo},
			{Title: "Secrets and paste", Description: "Paste KEY=VALUE lines to add rows; Enter in the last value adds one.",
				Demo: keyValueEditorSecretsDemo},
		},
		Code: `var env []form.KeyValue // kept in application state
form.KeyValueEditor(c, form.KeyValueEditorOptions{
    AriaLabel: "Environment variables",
    Pairs:     env,
    OnChange:  func(i int, p form.KeyValue) { env[i] = p },
    OnAdd:     func() { env = append(env, form.KeyValue{}) },
    OnRemove:  func(i int) { env = append(env[:i:i], env[i+1:]...) },
})`,
		Params: []ComponentParam{
			{"Pairs", "[]form.KeyValue", "nil", "Caller-owned rows; never modified. ID keeps a row's state, Secret conceals its value."},
			{"AriaLabel", "string", "\"\"", "Accessible name of the group."},
			{"KeyPlaceholder", "string", "Key", "Placeholder and name of the key editors."},
			{"ValuePlaceholder", "string", "Value", "Placeholder and name of the value editors."},
			{"AddLabel", "string", "Add", "Text of the add button."},
			{"MarkDuplicateKeys", "bool", "false", "Marks repeated non-empty keys invalid, with a message."},
			{"ShowHeaders", "bool", "false", "Titles the key and value columns."},
			{"Disabled", "bool", "false", "Dims the editor and stops input."},
			{"OnChange", "func(int, form.KeyValue)", "nil", "Reports an edited row; nil makes rows read-only."},
			{"OnAdd", "func()", "nil", "Asks for a new last row; nil omits the button."},
			{"OnRemove", "func(int)", "nil", "Asks for a row to go; nil omits the buttons."},
			{"OnPaste", "func([]form.KeyValue)", "nil", "Receives the rows of pasted KEY=VALUE lines."},
		},
		Tokens: []string{"Input", "Ring", "Destructive", "Accent", "FieldPlaceholder"},
	})

	RegisterComponent(ComponentPage{
		ID: "maskinput", Title: "MaskInput", Group: "Forms", Order: 350,
		Description: "A text input that fits what is typed to a pattern and writes its literals.",
		Examples: []ComponentExample{
			{Title: "Date", Description: "9 takes a digit; the slashes are written for you.",
				Demo: maskInputDateDemo},
			{Title: "Letters and digits", Demo: maskInputPlateDemo},
			{Title: "Empty slots", Description: "Underscores show what is left; a pasted +1 (415) 555-0100 keeps the last ten digits. Undo steps through the raw value.",
				Demo: maskInputPhoneDemo},
		},
		Code: `var birthday string // kept in application state; holds "25121999"
form.MaskInput(c, &birthday, form.MaskInputOptions{
    Mask:        "99/99/9999",
    Label:       "Date of birth",
    Placeholder: "DD/MM/YYYY",
    OnChange:    save,
})`,
		Params: append(slices.Clone(fieldParams),
			ComponentParam{"Mask", "string", "required", "9 digit, A letter, * either; other runes are literals."},
			ComponentParam{"ReadOnly", "bool", "false", "Shows the value without letting it change."},
			ComponentParam{"OnChange", "func()", "nil", "Runs after the raw value changes."},
			ComponentParam{"OnSubmit", "func()", "nil", "Runs on Enter."},
		),
		Tokens: []string{"Input", "Ring", "Foreground", "FieldPlaceholder", "Destructive"},
	})
}
