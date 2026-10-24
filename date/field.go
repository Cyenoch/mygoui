package date

import (
	"strings"
	"time"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DateFieldOptions configures a HeroUI date field. The value is caller-owned:
// a nil pointer keeps the committed date in the field itself and reports it
// through OnChange only.
type DateFieldOptions struct {
	Label, Description, Error string
	Required, Invalid         bool
	Disabled, FullWidth       bool

	Value *time.Time
	// Width overrides the 160 DIP default. FullWidth wins over Width.
	Width         float32
	Min, Max      time.Time
	DisabledDates func(time.Time) bool
	// Placeholder shows while the field is empty. Empty uses "YYYY-MM-DD".
	Placeholder string

	Handle   *ui.Handle
	OnChange func(time.Time)
}

// DateField renders an ISO 8601 date editor (YYYY-MM-DD) inside HeroUI's date
// input group. A complete, in-bounds date commits as it is typed; partial text
// stays quiet, and an impossible or out-of-bounds date shows the localized
// invalid message.
func DateField(c *ui.Context, o DateFieldOptions) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	root := ui.Box(c).MinWidth(0)
	if o.FullWidth {
		root.FillWidth()
	}
	value := o.Value
	if value == nil {
		value = ui.Local(root, dateValueKey{}, func() time.Time { return time.Time{} })
	}
	spec := dateFieldSpec{
		field: fieldSpec{
			label:       o.Label,
			description: o.Description,
			error:       o.Error,
			required:    o.Required,
			invalid:     o.Invalid,
			disabled:    o.Disabled,
			fullWidth:   o.FullWidth,
		},
		value:         value,
		min:           o.Min,
		max:           o.Max,
		disabledDates: o.DisabledDates,
		placeholder:   o.Placeholder,
		loc:           locationOf(*value, c),
		width:         dateFieldWidth(c, o.Width, o.FullWidth),
		onChange:      o.OnChange,
	}
	buf := newDateBuffer(root, spec)
	var field ui.Element
	root.Children(func() {
		field, _ = dateField(c, spec, buf)
	})
	if o.Handle != nil {
		field.Bind(o.Handle)
	}
	return root
}

type dateValueKey struct{}
type dateBufferKey struct{}

// fieldSpec is HeroUI's date-field/time-field chrome: a 4 DIP gapped column
// with a label, the control, and either its description or its error.
type fieldSpec struct {
	label, description, error string
	required, invalid         bool
	disabled, fullWidth       bool
}

// formField builds the labelled field around one native control, reusing the
// form package's native association and error/description policy.
func formField(c *ui.Context, spec fieldSpec, errorText string, body func()) ui.Element {
	description := spec.description
	if errorText != "" {
		description = ""
	}
	return form.Field(c, form.FieldOptions{
		Label: spec.label, Description: description, Error: errorText,
		Required: spec.required, Disabled: spec.disabled,
		Invalid: spec.invalid || errorText != "", FullWidth: spec.fullWidth,
	}, body)
}

// fieldGroup is HeroUI's date-input-group surface: the field background, its
// hover and focus-within states, the danger ring when invalid, and 12 DIP
// horizontal padding around a 36 DIP row.
func fieldGroup(c *ui.Context, invalid, disabled bool) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().Field
	e := ui.Row(c).MinHeight(theme.Scale(c, 36)).AlignItems(ui.Center).
		PaddingX(theme.Space12).Radius(radius).Clip()
	// Snapshot the states in the build: a retained copy of a leaving
	// element paints in later passes, where e has expired.
	focused, hovered := e.FocusWithin(), e.Hovered()
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		fill := t.FieldBackground
		switch {
		case focused:
			fill = t.FieldFocus
		case hovered:
			fill = t.FieldHover
		}
		p.Fill(r, fill, radius)
		switch {
		case invalid:
			width := float32(1)
			if focused {
				width = 2
			}
			p.Stroke(r, t.Danger, radius, width)
		case focused:
			p.FocusRing(r, [4]float32{radius, radius, radius, radius})
		}
	})
	e = style.Shadow(e, c, theme.FieldShadows(c.Theme().Dark))
	return style.Disabled(e, disabled)
}

// fieldInput styles a native text editor as one date-input-group input.
func fieldInput(c *ui.Context, text *string, placeholder string) ui.Element {
	e := ui.TextInputBase(c, text).Grow(1).MinWidth(0).Padding(0).FocusRing(false)
	e = formstyle.Text(e, c)
	return e.Placeholder(placeholder)
}

type dateFieldSpec struct {
	field         fieldSpec
	value         *time.Time
	min, max      time.Time
	disabledDates func(time.Time) bool
	placeholder   string
	loc           *time.Location
	width         float32
	suffix        func()
	onChange      func(time.Time)
}

// dateField builds the field and returns it with the input group, which
// DatePicker anchors its popover to.
func dateField(c *ui.Context, spec dateFieldSpec, buf *dateBuffer) (ui.Element, ui.Element) {
	buf.refresh(*spec.value)
	errorText := spec.field.error
	invalid := spec.field.invalid || buf.invalid
	if errorText == "" && buf.invalid {
		errorText = locale.M().Invalid
	}
	var group ui.Element
	field := formField(c, spec.field, errorText, func() {
		group = dateGroup(c, spec, buf, invalid)
	})
	return field, group
}

func dateGroup(c *ui.Context, spec dateFieldSpec, buf *dateBuffer, invalid bool) ui.Element {
	placeholder := spec.placeholder
	if placeholder == "" {
		placeholder = "YYYY-MM-DD"
	}
	group := fieldGroup(c, invalid, spec.field.disabled)
	if spec.width > 0 && !spec.field.fullWidth {
		group = group.Width(spec.width)
	}
	if spec.field.fullWidth {
		group = group.FillWidth()
	}
	group.Children(func() {
		input := fieldInput(c, &buf.text, placeholder)
		input.OnChange(func() { buf.commit(spec) })
		if spec.suffix != nil {
			spec.suffix()
		}
	})
	return group
}

// dateBuffer is the editable text and the value it was last synchronised with.
// It lives in the component's build state, not in the caller's options.
type dateBuffer struct {
	text      string
	committed time.Time
	invalid   bool
}

func newDateBuffer(root ui.Element, spec dateFieldSpec) *dateBuffer {
	return ui.Local(root, dateBufferKey{}, func() dateBuffer {
		return dateBuffer{
			text:      formatISO(*spec.value),
			committed: *spec.value,
		}
	})
}

// refresh follows a value changed outside the editor.
func (b *dateBuffer) refresh(value time.Time) {
	if sameDay(value, b.committed) {
		return
	}
	b.text = formatISO(value)
	b.committed = value
	b.invalid = false
}

// commit parses the text after an edit and reports a complete valid date.
func (b *dateBuffer) commit(spec dateFieldSpec) {
	if b.text == "" {
		b.invalid = false
		if !b.committed.IsZero() {
			b.committed = time.Time{}
			*spec.value = time.Time{}
			if spec.onChange != nil {
				spec.onChange(time.Time{})
			}
		}
		return
	}
	if len(b.text) < len("2006-01-02") {
		// Still typing: partial text is not an error.
		b.invalid = false
		return
	}
	parsed, ok := parseISO(b.text, spec.loc)
	if !ok || outOfBounds(parsed, spec.min, spec.max) ||
		(spec.disabledDates != nil && spec.disabledDates(parsed)) {
		b.invalid = true
		return
	}
	b.invalid = false
	if sameDay(parsed, b.committed) {
		return
	}
	b.committed = parsed
	*spec.value = parsed
	if spec.onChange != nil {
		spec.onChange(parsed)
	}
}

func dateFieldWidth(c *ui.Context, width float32, full bool) float32 {
	if full {
		return 0
	}
	if width > 0 {
		return theme.Scale(c, width)
	}
	return theme.Scale(c, 160)
}

func locationOf(value time.Time, c *ui.Context) *time.Location {
	if !value.IsZero() {
		return value.Location()
	}
	return c.Now().Location()
}

func formatISO(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func parseISO(text string, loc *time.Location) (time.Time, bool) {
	d, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(text), loc)
	if err != nil {
		return time.Time{}, false
	}
	return dayStart(d), true
}

func outOfBounds(d, min, max time.Time) bool {
	if !min.IsZero() && compareDay(d, min) < 0 {
		return true
	}
	return !max.IsZero() && compareDay(d, max) > 0
}
