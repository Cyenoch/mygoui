package form

import (
	"math"
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// NumberFieldState is the caller-owned value of a number field. Value is the
// committed number; Text is the editor's locale-independent text ("." decimal
// separator, no grouping). Change Value between builds to update the field.
type NumberFieldState struct {
	Value float64
	Text  string

	handle      ui.Handle
	focused     bool
	initialized bool
	lastValue   float64
}

// NumberFieldOptions configures a number field. Min and Max are nil when
// unbounded; Step defaults to 1. Format replaces the default display, which is
// rounded to the step's decimals.
type NumberFieldOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth                 bool
	Min, Max                  *float64
	Step                      float64
	Format                    func(float64) string
	OnChange                  func()
}

// numberStep normalizes a step option.
func numberStep(step float64) float64 {
	step = math.Abs(step)
	if step == 0 || math.IsNaN(step) || math.IsInf(step, 0) {
		return 1
	}
	return step
}

// numberDecimals returns the decimals a step needs to be shown exactly.
func numberDecimals(step float64) int {
	s := strconv.FormatFloat(step, 'f', -1, 64)
	if i := strings.IndexByte(s, '.'); i >= 0 {
		return len(s) - i - 1
	}
	return 0
}

// numberRound rounds v to the step's decimals, avoiding binary residue such as
// 0.30000000000000004.
func numberRound(v, step float64) float64 {
	decimals := numberDecimals(step)
	if decimals <= 0 {
		return math.Round(v)
	}
	p := math.Pow(10, float64(decimals))
	return math.Round(v*p) / p
}

// numberFormatValue renders v without locale formatting.
func numberFormatValue(v float64, step float64, format func(float64) string) string {
	v = numberRound(v, step)
	if format != nil {
		return format(v)
	}
	return strconv.FormatFloat(v, 'f', numberDecimals(step), 64)
}

// numberBounds resolves the optional bounds into a usable range.
func numberBounds(o NumberFieldOptions) (lo, hi float64) {
	lo, hi = -math.MaxFloat64, math.MaxFloat64
	if o.Min != nil {
		lo = *o.Min
	}
	if o.Max != nil {
		hi = *o.Max
	}
	if hi < lo {
		lo, hi = hi, lo
	}
	return lo, hi
}

// numberClamp keeps v inside the bounds.
func numberClamp(v, lo, hi float64) float64 {
	switch {
	case v < lo:
		return lo
	case v > hi:
		return hi
	}
	return v
}

// NumberField renders a native numeric editor between increment and decrement
// buttons. Typed text commits only as a finite in-range number; out-of-range
// numbers clamp and invalid text returns to the last value.
func NumberField(c *ui.Context, state *NumberFieldState, o NumberFieldOptions) ui.Element {
	if state == nil {
		state = &NumberFieldState{}
	}
	step := numberStep(o.Step)
	lo, hi := numberBounds(o)
	format := func(v float64) string { return numberFormatValue(v, step, o.Format) }

	// Settle the text before the editor reads it, so blur and programmatic
	// values format in the frame that builds them.
	focused := state.handle.Focused(c)
	if !state.initialized {
		state.Value = numberRound(numberClamp(state.Value, lo, hi), step)
		state.Text = format(state.Value)
		state.lastValue = state.Value
		state.initialized = true
	} else if focused != state.focused {
		if !focused {
			state.Text = format(state.Value)
			state.lastValue = state.Value
		}
	} else if !focused && state.Value != state.lastValue {
		state.Value = numberRound(numberClamp(state.Value, lo, hi), step)
		state.Text = format(state.Value)
		state.lastValue = state.Value
	}
	state.focused = focused

	apply := func(v float64) {
		v = numberRound(numberClamp(v, lo, hi), step)
		changed := v != state.Value
		state.Value = v
		state.Text = format(v)
		state.lastValue = v
		if changed && o.OnChange != nil {
			o.OnChange()
		}
	}

	t := theme.Colors(c)
	height := theme.Scale(c, 32)
	buttonWidth := theme.Scale(c, 36)

	build := func() ui.Element {
		row := ui.Grid(c).ColumnTracks(ui.Fixed(buttonWidth), ui.Fr(1), ui.Fixed(buttonWidth)).
			AlignItems(ui.Center).MinHeight(height).Shrink(0)
		if o.FullWidth {
			row.FillWidth()
		} else {
			row.MinWidth(theme.Scale(c, inputDefaultWidth))
		}
		row = fieldSurface(row, c, o.Variant, o.Invalid, o.Disabled)

		var editor ui.Element
		row.Children(func() {
			button := func(label string, icon *ui.SVG, delta float64, atBound bool) {
				b := ui.ButtonBase(c).Width(buttonWidth).Height(height).Center().Shrink(0).
					Padding(0).Radius(0).FocusRing(false).Label(label).Tooltip(label)
				b.TextColor(t.FieldForeground)
				blocked := o.Disabled || o.ReadOnly || atBound
				b = style.Disabled(b, blocked)
				b = style.Background(b, 0, ui.Transparent, t.Accent, t.Accent)

				b.Children(func() {
					ui.Icon(c, icon).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
				})
				if !blocked {
					b.OnClick(func() { apply(state.Value + delta) })
				}
			}
			// Grid children must be built in visual column order.
			button(locale.M().Decrease, lucide.Minus(), -step, state.Value <= lo)
			editor = ui.TextInputBase(c, &state.Text).MinWidth(0).TextAlign(ui.Center).
				PaddingX(theme.Space12).PaddingY(theme.Scale(c, 5)).FontFeatures("tnum")
			editor = fieldText(editor, c)
			if o.Placeholder != "" {
				editor.Placeholder(o.Placeholder)
			}
			if name := inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder); name != "" {
				editor.Label(name)
			}
			if o.ReadOnly {
				editor.ReadOnly(true)
			}
			if o.Min != nil || o.Max != nil {
				editor.Range(lo, hi, state.Value)
			}
			editor.OnChange(func() {
				if o.Disabled || o.ReadOnly {
					return
				}
				raw := strings.TrimSpace(state.Text)
				if raw == "" {
					return
				}
				v, err := strconv.ParseFloat(raw, 64)
				if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
					state.Text = format(state.Value)
					return
				}
				next := numberRound(numberClamp(v, lo, hi), step)
				if next != v {
					state.Text = format(next)
				}
				if next != state.Value {
					state.Value = next
					state.lastValue = next
					if o.OnChange != nil {
						o.OnChange()
					}
				}
			})
			editor.OnTextKey(func(k ui.TextKey) bool {
				if o.Disabled || o.ReadOnly {
					return false
				}
				switch k.Key {
				case ui.KeyUp:
					apply(state.Value + step)
				case ui.KeyDown:
					apply(state.Value - step)
				default:
					return false
				}
				k.SetText(state.Text)
				return true
			})

			button(locale.M().Increase, lucide.Plus(), step, state.Value >= hi)
		})

		if editor.Valid() {
			editor.Bind(&state.handle)
		}
		return row
	}

	fieldOptions := InputOptions{
		Label:       o.Label,
		Description: o.Description,
		Error:       o.Error,
		Required:    o.Required,
		Disabled:    o.Disabled,
		Invalid:     o.Invalid,
		FullWidth:   o.FullWidth,
	}
	if !inputNeedsField(fieldOptions) {
		return build()
	}
	return inputField(c, fieldOptions, true, func() { build() })
}
