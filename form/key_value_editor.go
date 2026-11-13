package form

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// KeyValue is one row of a KeyValueEditor. ID keeps the row's identity, and
// so its focus, caret and reveal state, when rows before it come and go; rows
// without an ID are told apart by position. A Secret row shows its value
// concealed behind a reveal toggle.
type KeyValue struct {
	ID         string
	Key, Value string
	Secret     bool
}

// KeyValueEditorOptions configures an editor of key and value rows, such as
// environment variables or request headers. Pairs are caller-owned and never
// modified: OnChange reports the row at index as edited, OnAdd asks for a new
// row at the end and OnRemove for the row at index to go, and the caller
// applies each. A nil OnChange makes the rows read-only; a nil OnAdd or
// OnRemove omits that action.
//
// Text pasted into an editor as KEY=VALUE lines (a .env file, with comments,
// blank lines and "export" skipped) is reported to OnPaste as the rows it
// holds instead of being typed in; a single line is taken only in a key
// editor. Without OnPaste the text is pasted as typed.
//
// MarkDuplicateKeys marks every non-empty key that occurs more than once as
// invalid, with a message. ShowHeaders titles the columns. Enter in a key
// moves to its value and Enter in the last row's value asks for a new row.
// The rows are not virtualized: the editor suits tens of rows.
type KeyValueEditorOptions struct {
	Pairs             []KeyValue
	AriaLabel         string
	KeyPlaceholder    string
	ValuePlaceholder  string
	AddLabel          string
	MarkDuplicateKeys bool
	ShowHeaders       bool
	Disabled          bool
	OnChange          func(index int, pair KeyValue)
	OnAdd             func()
	OnRemove          func(index int)
	OnPaste           func(pairs []KeyValue)
}

type keyValueState struct {
	// focusRow is the row count at which the row just asked for exists.
	focusRow int
	// valueRow is the row, counted from one, whose value takes the focus.
	valueRow int
	pasted   []KeyValue
	// rows and generation tell when rows without an ID shifted: a removal
	// renews their identity, so a shifted row never inherits the reveal of
	// the row that stood in its place.
	rows, generation int
}

type (
	keyValueStateKey  struct{}
	keyValueRevealKey struct{}
	// keyValueRowKey keeps a row's identity by its ID, else its position
	// within a generation of rows.
	keyValueRowKey struct {
		id                string
		index, generation int
	}
)

// KeyValueEditor renders one row of two native editors per pair, a remove
// button after each row and an add button below. The row added through the
// add button takes the focus in its key.
func KeyValueEditor(c *ui.Context, o KeyValueEditorOptions) ui.Element {
	m := locale.M()
	if o.KeyPlaceholder == "" {
		o.KeyPlaceholder = m.Key
	}
	if o.ValuePlaceholder == "" {
		o.ValuePlaceholder = m.Value
	}
	if o.AddLabel == "" {
		o.AddLabel = m.Add
	}
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).
		Role(ui.RoleGroup).Label(o.AriaLabel)
	st := ui.Local(root, keyValueStateKey{}, func() keyValueState { return keyValueState{} })
	focusLast := st.focusRow == len(o.Pairs)
	if len(o.Pairs) < st.rows {
		st.generation++
	}
	st.rows = len(o.Pairs)
	st.focusRow = 0

	var duplicates map[string]int
	if o.MarkDuplicateKeys {
		duplicates = make(map[string]int, len(o.Pairs))
		for _, p := range o.Pairs {
			duplicates[p.Key]++
		}
	}
	edited, pair, add := -1, KeyValue{}, false
	root.Children(func() {
		if o.ShowHeaders {
			keyValueHeaders(c, o)
		}
		for i, p := range o.Pairs {
			row := keyValueRow{
				index: i, pair: p,
				duplicate: p.Key != "" && duplicates[p.Key] > 1,
				focus:     focusLast && i == len(o.Pairs)-1,
				last:      i == len(o.Pairs)-1,
			}
			next, submitted := row.build(c, o, st)
			if next != p {
				edited, pair = i, next
			}
			add = add || submitted
		}
		if o.OnAdd != nil {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: o.AddLabel, Icon: lucide.Plus(), Variant: mygoui.ButtonOutline,
				Size: theme.ButtonSmall, Disabled: o.Disabled,
				OnClick: func() {
					st.focusRow = len(o.Pairs) + 1
					o.OnAdd()
				},
			}).AlignSelf(ui.Start)
		}
	})
	// The caller edits its pairs here, after the rows were built from them.
	switch {
	case len(st.pasted) > 0:
		pasted := st.pasted
		st.pasted = nil
		o.OnPaste(pasted)
		c.Invalidate()
	case edited >= 0 && o.OnChange != nil:
		o.OnChange(edited, pair)
		c.Invalidate()
	case add && o.OnAdd != nil:
		st.focusRow = len(o.Pairs) + 1
		o.OnAdd()
		c.Invalidate()
	}
	return style.Disabled(root, o.Disabled)
}

// keyValueColumns lays a row out as key, value and the remove button.
func keyValueColumns(c *ui.Context, key any) ui.Element {
	return ui.Grid(c.Key(key)).ColumnTracks(ui.Fr(1), ui.Fr(1), ui.FitContent()).
		FillWidth().MinWidth(0).Gap(theme.Space8).AlignItems(ui.Start)
}

// keyValueActionSize is the width and height of the remove button.
func keyValueActionSize(c *ui.Context) float32 {
	return theme.Scale(c, theme.ButtonHeight(theme.ButtonSmall))
}

// keyValueHeaders titles the key and value columns.
func keyValueHeaders(c *ui.Context, o KeyValueEditorOptions) {
	t := theme.Colors(c)
	keyValueColumns(c, keyValueStateKey{}).Children(func() {
		for _, title := range []string{o.KeyPlaceholder, o.ValuePlaceholder} {
			style.Text(ui.Text(c, title), c, theme.XS()).FontWeight(theme.WeightMedium).
				TextColor(t.MutedForeground).SingleLine().PaddingX(theme.Space12)
		}
		if o.OnRemove != nil {
			ui.Box(c).Width(keyValueActionSize(c))
		} else {
			ui.Box(c)
		}
	})
}

// keyValueRow is one pair being built.
type keyValueRow struct {
	index     int
	pair      KeyValue
	duplicate bool
	focus     bool
	last      bool
}

// build renders the row and returns the pair as the editors now hold it, and
// whether Enter in the value of the last row asked for another row.
func (r keyValueRow) build(c *ui.Context, o KeyValueEditorOptions, st *keyValueState) (next KeyValue, add bool) {
	next = r.pair
	number := strconv.Itoa(r.index + 1)
	rowKey := keyValueRowKey{id: r.pair.ID}
	if rowKey.id == "" {
		rowKey.index, rowKey.generation = r.index, st.generation
	}
	grid := keyValueColumns(c, rowKey)
	revealed := ui.Local(grid, keyValueRevealKey{}, func() bool { return false })
	grid.Children(func() {
		editor := func(value *string, placeholder string, o2 inputCoreOptions) ui.Element {
			o2.Placeholder, o2.Name = placeholder, placeholder+" "+number
			o2.Disabled, o2.ReadOnly, o2.FullWidth = o.Disabled, o.OnChange == nil, true
			e := inputEditor(c, value, o2).MinWidth(0)
			if o.OnPaste != nil {
				e.OnPaste(func(text string) bool {
					pairs, ok := keyValueParse(text, value == &next.Key)
					if ok {
						st.pasted = append(st.pasted, pairs...)
					}
					return ok
				})
			}
			return e
		}
		ui.Column(c).MinWidth(0).Gap(theme.Space4).Children(func() {
			key := editor(&next.Key, o.KeyPlaceholder, inputCoreOptions{Invalid: r.duplicate}).Invalid(r.duplicate)
			if r.focus {
				key.AutoFocus()
			}
			if key.Submitted() && o.OnChange != nil {
				st.valueRow = r.index + 1
			}
			if r.duplicate {
				msg := locale.M().DuplicateKey
				key.Description(msg)
				style.Text(ui.Text(c, msg), c, theme.HelpText()).
					TextColor(theme.Colors(c).Destructive).PaddingX(theme.Space2)
			}
		})
		ui.Box(c).MinWidth(0).Children(func() {
			conceal := r.pair.Secret && !*revealed
			value := editor(&next.Value, o.ValuePlaceholder, inputCoreOptions{Password: conceal})
			if st.valueRow == r.index+1 {
				value.Focus()
				st.valueRow = 0
			}
			if value.Submitted() && r.last && o.OnChange != nil && next != (KeyValue{ID: next.ID, Secret: next.Secret}) {
				add = true
			}
			if r.pair.Secret {
				size := keyValueActionSize(c)
				value.PaddingEnd(size + theme.Scale(c, theme.Space2))
				passwordToggle(c, revealed, size, o.Disabled, r.pair.Key).Absolute().End(theme.Scale(c, theme.Space2)).Top(0)
			}
		})
		if o.OnRemove == nil {
			ui.Box(c)
			return
		}
		name := r.pair.Key
		if name == "" {
			name = number
		}
		ui.Box(c).Height(theme.Scale(c, 32)).Center().Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: fmt.Sprintf(locale.M().RemoveNamed, name), Icon: lucide.X(), IconOnly: true,
				Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall, Disabled: o.Disabled,
				OnClick: func() { o.OnRemove(r.index) },
			})
		})
	})
	return next, add
}

// keyValueParse reads KEY=VALUE lines. Blank lines, # comments and a leading
// "export" are skipped, and a value in matching quotes loses them. It
// reports false unless every other line is a pair, and for a single line
// outside a key editor, which is a value being pasted.
func keyValueParse(text string, inKey bool) ([]KeyValue, bool) {
	var pairs []KeyValue
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		key, value, ok := strings.Cut(line, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" || strings.ContainsAny(key, " \t") {
			return nil, false
		}
		value = strings.TrimSpace(value)
		if n := len(value); n >= 2 && (value[0] == '"' || value[0] == '\'') && value[n-1] == value[0] {
			value = value[1 : n-1]
		}
		pairs = append(pairs, KeyValue{Key: key, Value: value})
	}
	return pairs, len(pairs) > 1 || (len(pairs) == 1 && inKey)
}
