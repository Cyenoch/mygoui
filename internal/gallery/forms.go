package gallery

import (
	"strconv"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/egoist/mygo/ui"
)

// FormsState contains editable gallery fixtures, never production input or credentials.
type FormsState struct {
	Text, Notes, Choice, Result                                   string
	Checked, Mixed, Enabled, Disabled, ReadOnly, Invalid, Loading bool
	Value, Low, High                                              float64
	OTP                                                           form.OTPState
	InputHandle                                                   ui.Handle
}

func NewFormsState() FormsState {
	return FormsState{Text: "Native editable sample", Notes: "Edit this multiline fixture.\nSelection and undo are native.",
		Choice: "balanced", Value: 40, Low: 20, High: 70}
}

func PreviewForm(c *ui.Context, s *FormsState, id ComponentID) {
	ui.Column(c).FillWidth().MinWidth(0).Justify(ui.Start).AlignItems(ui.Stretch).Gap(16).Children(func() {
		ui.Row(c).FillWidth().Wrap().Justify(ui.Start).AlignContent(ui.Start).Gap(8).Children(func() {
			mygoui.Toggle(c, &s.Disabled, mygoui.ToggleOptions{Label: "Disabled"})
			mygoui.Toggle(c, &s.ReadOnly, mygoui.ToggleOptions{Label: "Read only"})
			mygoui.Toggle(c, &s.Invalid, mygoui.ToggleOptions{Label: "Invalid"})
			mygoui.Toggle(c, &s.Loading, mygoui.ToggleOptions{Label: "Loading"})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Reset form", Variant: mygoui.ButtonOutline,
				OnClick: func() { *s = NewFormsState() }})
		})
		errorText := ""
		if s.Invalid {
			errorText = "Sample validation error. Change the Invalid fixture to clear it."
		}
		o := form.ControlOptions{Label: "Sample value", Description: "Editable local fixture", Error: errorText,
			Disabled: s.Disabled, ReadOnly: s.ReadOnly, Invalid: s.Invalid, Loading: s.Loading}
		switch id {
		case "input":
			form.Field(c, form.FieldOptions{Label: "Sample input", Description: "Native selection, paste, undo and IME", Error: errorText}, func() {
				form.Input(c, &s.Text, form.InputOptions{ControlOptions: o, Placeholder: "Enter sample text",
					OnSubmit: func(value string) { s.Result = "Submitted: " + value; c.AddToast(ui.Toast{Title: s.Result}) }})
			})
		case "textarea":
			form.Field(c, form.FieldOptions{Label: "Sample notes", Error: errorText}, func() {
				form.Textarea(c, &s.Notes, form.TextareaOptions{InputOptions: form.InputOptions{ControlOptions: o}, MinLines: 3, MaxLines: 8})
			})
		case "label":
			o.Handle = &s.InputHandle
			form.Input(c, &s.Text, form.InputOptions{ControlOptions: o})
			form.Label(c, "Focus sample input", &s.InputHandle)
		case "field":
			form.FieldSet(c, form.FieldSetOptions{Legend: "Sample preferences", Description: "Native labels connect to their fields", Disabled: s.Disabled}, func() {
				form.Field(c, form.FieldOptions{Label: "Display name", Description: "Choose any sample name", Error: errorText}, func() {
					form.Input(c, &s.Text, form.InputOptions{ControlOptions: o})
				})
				form.Field(c, form.FieldOptions{Label: "Notes", Orientation: form.FieldResponsive}, func() {
					form.Textarea(c, &s.Notes, form.TextareaOptions{InputOptions: form.InputOptions{ControlOptions: o}, MinLines: 2, MaxLines: 5})
				})
			})
		case "input-group":
			form.InputGroup(c, &s.Text, form.InputGroupOptions{InputOptions: form.InputOptions{ControlOptions: o},
				InlineStart: func() { form.InputGroupText(c, "https://") },
				BlockEnd: func() {
					ui.Row(c).FillWidth().Wrap().Gap(8).Children(func() {
						form.InputGroupText(c, "Local sample address")
						form.InputGroupButton(c, mygoui.ButtonOptions{Label: "Use sample", Variant: mygoui.ButtonGhost,
							Size: mygoui.ButtonXS, Disabled: s.Disabled, OnClick: func() { s.Result = "Used: " + s.Text; c.AddToast(ui.Toast{Title: s.Result}) }})
					})
				},
			})
		case "input-otp":
			ui.ScrollHorizontal(c).FillWidth().MinWidth(0).Height(112).Children(func() {
				form.InputOTP(c, &s.OTP, form.OTPOptions{ControlOptions: o, Length: 6, SeparatorAfter: 3,
					OnComplete: func(value string) {
						s.Result = "Completed sample code: " + value
						c.AddToast(ui.Toast{Title: s.Result})
					}})
			})
			ui.Text(c, "Fixture value: "+s.OTP.Value)
		case "checkbox":
			mygoui.Toggle(c, &s.Mixed, mygoui.ToggleOptions{Label: "Mixed state"})
			o.Label = "Include sample details"
			form.Checkbox(c, &s.Checked, form.CheckboxOptions{ControlOptions: o, Mixed: &s.Mixed})
			ui.Text(c, "Checked: "+strconv.FormatBool(s.Checked)+"; mixed: "+strconv.FormatBool(s.Mixed))
		case "switch":
			o.Label = "Enable sample option"
			form.Switch(c, &s.Enabled, form.SwitchOptions{ControlOptions: o})
			ui.Text(c, "Enabled: "+strconv.FormatBool(s.Enabled))
		case "radio-group":
			form.RadioGroup(c, &s.Choice, form.RadioGroupOptions{ControlOptions: o}, []form.RadioItem{
				{Value: "fast", Label: "Fast", Description: "Less fixture detail"},
				{Value: "balanced", Label: "Balanced", Description: "Regular fixture detail"},
				{Value: "complete", Label: "Complete", Description: "More fixture detail"},
				{Value: "unavailable", Label: "Unavailable", Disabled: true},
			})
			ui.Text(c, "Selected: "+s.Choice)
		case "slider":
			o.Label = "Sample level"
			form.Slider(c, &s.Value, form.SliderOptions{ControlOptions: o, Min: 0, Max: 100, Step: 5, Unit: "%"})
			form.RangeSlider(c, &s.Low, &s.High, form.RangeSliderOptions{ControlOptions: o, Min: 0, Max: 100, Step: 5})
			ui.Text(c, "Value: "+strconv.FormatFloat(s.Value, 'f', 0, 64)+"; range: "+
				strconv.FormatFloat(s.Low, 'f', 0, 64)+"–"+strconv.FormatFloat(s.High, 'f', 0, 64))
		case "native-select":
			form.NativeSelect(c, &s.Choice, form.NativeSelectOptions{ControlOptions: o, Placeholder: "Choose sample mode"}, []form.SelectOption{
				{Value: "fast", Label: "Fast", Group: "Available"}, {Value: "balanced", Label: "Balanced", Group: "Available"},
				{Value: "complete", Label: "Complete", Group: "Available"}, {Value: "unavailable", Label: "Unavailable", Disabled: true},
			})
			ui.Text(c, "Selected: "+s.Choice)
		}
	})
}
