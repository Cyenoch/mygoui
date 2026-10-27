package date

import (
	"strings"
	"time"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TimeFieldState is the caller-owned time of day shown by a TimeField. Only
// the clock fields matter; the date part is preserved so the value can be
// stored back into a full timestamp, and is otherwise ignored.
type TimeFieldState struct {
	Value time.Time
}

// TimeFieldOptions configures a time field. The editor is ISO 8601
// HH:MM by default, or hh:mm AM/PM when Hour12 is set.
type TimeFieldOptions struct {
	Label, Description, Error string
	Required, Invalid         bool
	Disabled, FullWidth       bool

	// Width overrides the 104 DIP default. FullWidth wins over Width.
	Width  float32
	Hour12 bool
	// Placeholder shows while the field is empty. Empty uses "HH:MM", or
	// "hh:mm AM" with Hour12.
	Placeholder string

	Handle   *ui.Handle
	OnChange func(time.Time)
}

// TimeField renders a COSS-style time field. A complete time commits as it is
// typed, keeping the date of the state's value (or today's date when it holds
// none). A nil state uses a component-local copy.
func TimeField(c *ui.Context, o TimeFieldOptions, st *TimeFieldState) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	root := ui.Box(c).MinWidth(0)
	if o.FullWidth {
		root.FillWidth()
	}
	state := st
	if state == nil {
		state = ui.Local(root, timeStateKey{}, func() TimeFieldState { return TimeFieldState{} })
	}
	width := o.Width
	if width == 0 {
		width = 104
	}
	width = theme.Scale(c, width)
	loc := c.Now().Location()
	if !state.Value.IsZero() {
		loc = state.Value.Location()
	}
	placeholder := o.Placeholder
	if placeholder == "" {
		if o.Hour12 {
			placeholder = "hh:mm AM"
		} else {
			placeholder = "HH:MM"
		}
	}
	spec := timeFieldSpec{
		field: fieldSpec{
			label:       o.Label,
			description: o.Description,
			error:       o.Error,
			required:    o.Required,
			invalid:     o.Invalid,
			disabled:    o.Disabled,
			fullWidth:   o.FullWidth,
		},
		state:       state,
		hour12:      o.Hour12,
		placeholder: placeholder,
		width:       width,
		loc:         loc,
		now:         c.Now(),
		onChange:    o.OnChange,
	}
	buf := newTimeBuffer(root, spec)

	var field ui.Element
	root.Children(func() {
		field = timeField(c, spec, buf)
	})
	if o.Handle != nil {
		field.Bind(o.Handle)
	}
	return root
}

type timeStateKey struct{}
type timeBufferKey struct{}

type timeFieldSpec struct {
	field       fieldSpec
	state       *TimeFieldState
	hour12      bool
	placeholder string
	width       float32
	loc         *time.Location
	now         time.Time
	onChange    func(time.Time)
}

func timeField(c *ui.Context, spec timeFieldSpec, buf *timeBuffer) ui.Element {
	buf.refresh(spec.state.Value, spec.hour12)
	invalid := spec.field.invalid || buf.invalid
	errorText := spec.field.error
	if errorText == "" && buf.invalid {
		errorText = locale.M().Invalid
	}
	return formField(c, spec.field, errorText, func() {
		group := fieldGroup(c, invalid, spec.field.disabled)
		if spec.width > 0 && !spec.field.fullWidth {
			group = group.Width(spec.width)
		}
		if spec.field.fullWidth {
			group = group.FillWidth()
		}
		group.Children(func() {
			input := fieldInput(c, &buf.text, spec.placeholder)
			input.OnChange(func() { buf.commit(spec) })
		})
	})
}

type timeBuffer struct {
	text      string
	committed time.Time
	invalid   bool
	hour12    bool
}

func newTimeBuffer(root ui.Element, spec timeFieldSpec) *timeBuffer {
	return ui.Local(root, timeBufferKey{}, func() timeBuffer {
		return timeBuffer{
			text:      formatClock(spec.state.Value, spec.hour12),
			committed: spec.state.Value,
			hour12:    spec.hour12,
		}
	})
}

func (b *timeBuffer) refresh(value time.Time, hour12 bool) {
	if sameClock(value, b.committed) && hour12 == b.hour12 {
		return
	}
	b.text = formatClock(value, hour12)
	b.committed = value
	b.invalid = false
	b.hour12 = hour12
}

func (b *timeBuffer) commit(spec timeFieldSpec) {
	if b.text == "" {
		b.invalid = false
		if !b.committed.IsZero() {
			b.committed = time.Time{}
			spec.state.Value = time.Time{}
			if spec.onChange != nil {
				spec.onChange(time.Time{})
			}
		}
		return
	}
	hour, minute, ok := parseClock(b.text, spec.hour12)
	if !ok {
		b.invalid = true
		return
	}
	b.invalid = false
	y, m, d := spec.state.Value.Date()
	if spec.state.Value.IsZero() {
		y, m, d = spec.now.Date()
	}
	next := time.Date(y, m, d, hour, minute, 0, 0, spec.loc)
	if sameClock(next, b.committed) {
		return
	}
	b.committed = next
	spec.state.Value = next
	if spec.onChange != nil {
		spec.onChange(next)
	}
}

func formatClock(t time.Time, hour12 bool) string {
	if t.IsZero() {
		return ""
	}
	if hour12 {
		return t.Format("03:04 PM")
	}
	return t.Format("15:04")
}

func parseClock(text string, hour12 bool) (int, int, bool) {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0, 0, false
	}
	if hour12 {
		text = strings.ToUpper(text)
		for _, layout := range [...]string{"03:04 PM", "3:04 PM"} {
			if t, err := time.Parse(layout, text); err == nil {
				return t.Hour(), t.Minute(), true
			}
		}
		return 0, 0, false
	}
	t, err := time.Parse("15:04", text)
	if err != nil {
		return 0, 0, false
	}
	return t.Hour(), t.Minute(), true
}

func sameClock(a, b time.Time) bool {
	if a.IsZero() || b.IsZero() {
		return a.IsZero() && b.IsZero()
	}
	ah, am, _ := a.Clock()
	bh, bm, _ := b.Clock()
	return ah == bh && am == bm
}
