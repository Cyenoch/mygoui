package calendar

import (
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type Preset struct {
	Label string
	Range Range
}

// DatePickerState belongs to one picker in one window. Calendar is authoritative;
// invalid editor drafts survive rebuilds without replacing its valid selection.
type DatePickerState struct {
	Calendar                    State
	Open                        bool
	Text, EndText               string
	Invalid, EndInvalid         string
	syncedStart, syncedEnd      Date
	initialized                 bool
	startHandle                 ui.Handle
	invalidKind, endInvalidKind dateInputError
}
type DatePickerOptions struct {
	Calendar  Options
	Label     string
	Range     bool
	Presets   []Preset
	Clearable bool
	OnChange  func()
}

type dateInputError uint8

const (
	dateInputValid dateInputError = iota
	dateInputMalformed
	dateInputUnavailable
	dateInputRangeUnavailable
)

func (kind dateInputError) message() string {
	switch kind {
	case dateInputMalformed:
		return locales.ComponentEnterISODate()
	case dateInputUnavailable:
		return locales.ComponentDateUnavailable()
	case dateInputRangeUnavailable:
		return locales.ComponentChooseAvailableRange()
	default:
		return ""
	}
}

func (s *DatePickerState) sync(o DatePickerOptions) {
	start, end := s.Calendar.Date, Date{}
	if o.Range {
		start, end = s.Calendar.Range.Start, s.Calendar.Range.End
	}
	if !s.initialized || start != s.syncedStart {
		s.Text = start.String()
		s.Invalid = ""
		s.invalidKind = dateInputValid
	}
	if !s.initialized || end != s.syncedEnd {
		s.EndText = end.String()
		s.EndInvalid = ""
		s.endInvalidKind = dateInputValid
	}
	s.syncedStart, s.syncedEnd = start, end
	s.initialized = true
}
func (s *DatePickerState) changed(o DatePickerOptions) {
	s.sync(o)
	if o.Calendar.OnChange != nil {
		o.Calendar.OnChange()
	}
	if o.OnChange != nil {
		o.OnChange()
	}
}
func (s *DatePickerState) edit(end bool, o DatePickerOptions) {
	text, invalid, kind := &s.Text, &s.Invalid, &s.invalidKind
	if end {
		text, invalid, kind = &s.EndText, &s.EndInvalid, &s.endInvalidKind
	}
	*invalid = ""
	*kind = dateInputValid
	raw := strings.TrimSpace(*text)
	d := Date{}
	if raw != "" {
		var err error
		d, err = Parse(raw)
		if err != nil {
			*kind = dateInputMalformed
			*invalid = kind.message()
			return
		}
		if !o.Calendar.allowed(d) {
			*kind = dateInputUnavailable
			*invalid = kind.message()
			return
		}
	}
	if o.Range {
		r := s.Calendar.Range
		if end {
			r.End = d
		} else {
			r.Start = d
		}
		if r.Start.Valid() && r.End.Valid() && !o.Calendar.allowedRange(r) {
			*kind = dateInputRangeUnavailable
			*invalid = kind.message()
			return
		}
		if r == s.Calendar.Range {
			return
		}
		s.Calendar.Range = r
	} else {
		if d == s.Calendar.Date {
			return
		}
		s.Calendar.Date = d
	}
	if d.Valid() {
		s.Calendar.Cursor = d
		s.Calendar.Month = Date{d.Year, d.Month, 1}
	}
	s.changed(o)
}

// DatePicker composes maintained native editors and the existing Popover with
// Calendar. Input uses the explicit YYYY-MM-DD contract; locale controls the
// calendar's CLDR presentation, rather than silently changing parser semantics.
func DatePicker(c *ui.Context, s *DatePickerState, o DatePickerOptions) ui.Element {
	if s == nil {
		return ui.Element{}
	}
	s.sync(o)
	if s.invalidKind != dateInputValid {
		s.Invalid = s.invalidKind.message()
	}
	if s.endInvalidKind != dateInputValid {
		s.EndInvalid = s.endInvalidKind.message()
	}
	p := theme.Colors(c)
	root := ui.Column(c).Gap(theme.Space4).MinWidth(0)
	root.Children(func() {
		if o.Label != "" {
			form.Label(c, o.Label, &s.startHandle)
		}
		ui.Row(c).Gap(theme.Space4).Children(func() {
			label := o.Label
			if label == "" {
				label = locales.ComponentDate()
			}
			if o.Range {
				label = locales.ComponentDateRangeStart(label)
			}
			form.Input(c.Key("start"), &s.Text, form.InputOptions{ControlOptions: form.ControlOptions{Label: label, Error: s.Invalid, ReadOnly: o.Calendar.ReadOnly, Handle: &s.startHandle}, Placeholder: "YYYY-MM-DD", OnChange: func(string) { s.edit(false, o) }, OnSubmit: func(string) { s.edit(false, o) }}).Grow(1).Basis(0)
			if o.Range {
				ui.Text(c, "–").TextColor(p.MutedForeground)
				endLabel := o.Label
				if endLabel == "" {
					endLabel = locales.ComponentDate()
				}
				form.Input(c.Key("end"), &s.EndText, form.InputOptions{ControlOptions: form.ControlOptions{Label: locales.ComponentDateRangeEnd(endLabel), Error: s.EndInvalid, ReadOnly: o.Calendar.ReadOnly}, Placeholder: "YYYY-MM-DD", OnChange: func(string) { s.edit(true, o) }, OnSubmit: func(string) { s.edit(true, o) }}).Grow(1).Basis(0)
			}
			anchor := mygoui.Button(c.Key("calendar-trigger"), mygoui.ButtonOptions{Label: locales.ComponentOpenCalendar(), Icon: lucide.Calendar(), Size: mygoui.ButtonIcon, Variant: mygoui.ButtonOutline, Disabled: o.Calendar.ReadOnly, OnClick: func() { s.Open = !s.Open }})
			if o.Clearable {
				mygoui.Button(c.Key("clear"), mygoui.ButtonOptions{Label: locales.ComponentClearDate(), Icon: lucide.X(), Size: mygoui.ButtonIcon, Variant: mygoui.ButtonGhost, Disabled: o.Calendar.ReadOnly || (s.Text == "" && s.EndText == ""), OnClick: func() {
					s.Calendar.Date = Date{}
					s.Calendar.Dates = nil
					s.Calendar.Range = Range{}
					s.Text = ""
					s.EndText = ""
					s.Invalid = ""
					s.EndInvalid = ""
					s.invalidKind, s.endInvalidKind = dateInputValid, dateInputValid
					s.changed(o)
				}})
			}
			overlay.Popover(c, anchor, &s.Open, overlay.PopoverOptions{Width: max(float32(304), theme.Scale(c, 252)+56), Placement: overlay.Placement{Align: overlay.AlignStart}}, func() {
				calendarOptions := o.Calendar
				calendarOptions.Mode = Single
				if o.Range {
					calendarOptions.Mode = RangeMode
				}
				calendarOptions.OnChange = func() {
					s.changed(o)
					if !o.Range || s.Calendar.Range.End.Valid() {
						s.Calendar.Handle.CancelFocus(c)
						s.Open = false
					}
				}
				_, windowHeight := c.Size()
				ui.ScrollBoth(c).FillWidth().Height(min(max(float32(320), theme.Scale(c, 340)), max(float32(120), windowHeight-96))).Children(func() {
					Calendar(c.Key("picker-grid"), &s.Calendar, calendarOptions)
					if len(o.Presets) > 0 {
						ui.Row(c).Wrap().Gap(theme.Space4).Children(func() {
							for i, preset := range o.Presets {
								valid := o.Calendar.allowed(preset.Range.Start)
								if o.Range {
									valid = o.Calendar.allowedRange(preset.Range)
								}
								mygoui.Button(c.Key(i), mygoui.ButtonOptions{Label: preset.Label, Variant: mygoui.ButtonGhost, Size: mygoui.ButtonSM, Disabled: !valid || o.Calendar.ReadOnly, OnClick: func() {
									if o.Range {
										s.Calendar.Range = preset.Range
									} else {
										s.Calendar.Date = preset.Range.Start
									}
									s.Calendar.Cursor = preset.Range.Start
									s.Calendar.Month = Date{preset.Range.Start.Year, preset.Range.Start.Month, 1}
									s.changed(o)
									s.Open = false
								}})
							}
						})
					}
				})
			})
		})
		// Reserved error line prevents input validation from moving following content.
		message := s.Invalid
		if s.EndInvalid != "" {
			if message != "" {
				message += "; "
			}
			message += s.EndInvalid
		}
		ui.Text(c, message).Height(theme.Scale(c, 36)).MaxLines(2).FontSize(theme.Scale(c, 12)).TextColor(p.Destructive).Role(ui.RoleStatus)
	})
	return root
}
