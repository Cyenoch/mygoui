package navigation

import (
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// CommandItem is an action found by its label or keywords. Key defaults to
// Label; keys must be unique across the command list and its groups.
type CommandItem struct {
	Key, Label, Keywords, Shortcut string
	Disabled                       bool
	OnSelect                       func()
}

// CommandGroup is a named group of actions in a command list.
type CommandGroup struct {
	Key, Label string
	Items      []CommandItem
}

// CommandState holds the query, list highlight and focus for one command,
// and a revision for item changes a frame cannot see.
//
// Revision must change when an item or a group is edited in place. Items or
// groups appearing or disappearing through a new slice or a different length
// are noticed without it.
type CommandState struct {
	Query string
	Input ui.Handle
	List  list.ListBoxState
	// Revision tells the command that items changed in place. Increment it
	// when an item's key, label, keywords, shortcut or disabled flag changes
	// without the slice or its length changing.
	Revision uint64

	results commandResults

	keyQuery     string
	keyHighlight string
	keyPending   bool
}

// CommandOptions configures a searchable command list.
type CommandOptions struct {
	Label, Placeholder, Empty string
	Items                     []CommandItem
	Groups                    []CommandGroup
	State                     *CommandState
	MaxHeight                 float32
	OnSelect                  func(key string)
}

// Command filters caller-supplied actions and activates a highlighted action
// by Enter or a pointer click. The query and highlight survive rebuilds.
func Command(c *ui.Context, o CommandOptions) ui.Element {
	if o.State == nil {
		panic("navigation: Command requires a caller-owned State")
	}
	s := o.State
	t := theme.Colors(c)
	activate := func(key string) {
		run := func(items []CommandItem) bool {
			for _, item := range items {
				itemKey := item.Key
				if itemKey == "" {
					itemKey = item.Label
				}
				if itemKey == key && !item.Disabled {
					if item.OnSelect != nil {
						item.OnSelect()
					}
					if o.OnSelect != nil {
						o.OnSelect(key)
					}
					return true
				}
			}
			return false
		}
		if run(o.Items) {
			return
		}
		for _, group := range o.Groups {
			if run(group.Items) {
				return
			}
		}
	}
	empty := o.Empty
	if empty == "" {
		empty = locale.M().NoResults
	}
	height := o.MaxHeight
	if height <= 0 {
		height = 320
	}
	label := o.Label
	if label == "" {
		label = locale.M().Search
	}
	root := ui.Column(c).FillWidth().MinWidth(0).Background(t.Popover).
		TextColor(t.PopoverForeground).Radius(theme.RadiusScale().LG).
		Border(theme.Hairline, t.Border).Clip()
	root.Children(func() {
		var input ui.Element
		ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(theme.Space6, theme.Space12).
			Gap(theme.Space8).Children(func() {
			ui.Icon(c, lucide.Search()).Size(16, 16).TextColor(t.MutedForeground)
			input = ui.TextInputBase(c, &s.Query).FillWidth().Grow(1).MinWidth(0).
				Label(label).Placeholder(o.Placeholder).Bind(&s.Input)
			style.Text(input, c, theme.SM()).TextColor(t.Foreground)
			input.OnTextKey(func(k ui.TextKey) bool {
				if k.Mods != 0 {
					return false
				}
				switch k.Key {
				case ui.KeyDown, ui.KeyUp, ui.KeyEnter:
				default:
					return false
				}
				text := k.Text()
				key := commandKeySelection(o, s, text, k.Key)
				if key != "" {
					if k.Key == ui.KeyEnter {
						s.Query = text
						activate(key)
						if s.Query != text {
							k.SetText(s.Query)
						}
					}
				}
				return true
			})
		})
		r := s.results.update(o, s.Revision, s.Query)
		*r.shortcut = func(shortcut string) {
			style.Text(ui.Text(c, shortcut), c, theme.XS()).TextColor(t.MutedForeground)
		}
		list.ListBox(c, list.ListBoxOptions{
			Label: label, Items: r.rows, Sections: r.sections, State: &s.List, Padding: theme.Space8,
			MaxHeight: height, Empty: empty, NoFocusOnClick: true,
			KeyboardHighlight: s.Input.Focused(c), OnAction: activate,
		})
		if input.Changed() {
			if s.keyPending && s.keyQuery == s.Query && s.keyHighlight != "" {
				s.List.SetHighlight(s.keyHighlight)
			} else {
				s.List.HighlightEdge(false)
			}
		}
		s.keyPending = false
	})
	return root
}

// commandKeySelection uses the editor's current text, which can be ahead of
// Query when several native input events arrive before the next frame.
func commandKeySelection(o CommandOptions, s *CommandState, text string, key ui.Key) string {
	keys := s.results.update(o, s.Revision, text).keys
	highlight := ""
	if s.keyPending && s.keyQuery == text {
		highlight = s.keyHighlight
	} else if text == s.Query {
		highlight = s.List.HighlightKey()
	}
	index := -1
	for i, candidate := range keys {
		if candidate == highlight {
			index = i
			break
		}
	}
	if index < 0 && len(keys) > 0 {
		index = 0
	}
	if key == ui.KeyDown && index >= 0 {
		index = min(index+1, len(keys)-1)
	} else if key == ui.KeyUp && index >= 0 {
		index = max(index-1, 0)
	}
	highlight = ""
	if index >= 0 {
		highlight = keys[index]
	}
	s.keyQuery, s.keyHighlight, s.keyPending = text, highlight, true
	if highlight != "" {
		s.List.SetHighlight(highlight)
	}
	return highlight
}

// commandResults is a command's filtered rows and the enabled keys among
// them in order, rebuilt only when the query, the items or groups, or the
// revision change. The rows' End closures stay valid across frames; they
// draw through shortcut, which Command sets on every build and a long list
// calls later, as its rows come into view during layout.
// It is shared rather than reached through the state, which the caller may
// move.
type commandResults struct {
	items      []CommandItem
	groups     []CommandGroup
	groupItems [][]CommandItem
	revision   uint64
	query      string
	valid      bool

	rows     []list.ListBoxItem
	sections []list.ListBoxSection
	keys     []string
	shortcut *func(string)
}

// update refilters for text when it or the source changed, and returns r.
func (r *commandResults) update(o CommandOptions, revision uint64, text string) *commandResults {
	query := strings.ToLower(strings.TrimSpace(text))
	if r.valid && query == r.query && r.same(o, revision) {
		return r
	}
	r.items, r.groups, r.revision, r.query, r.valid = o.Items, o.Groups, revision, query, true
	r.groupItems = r.groupItems[:0]
	r.keys = nil
	if r.shortcut == nil {
		r.shortcut = new(func(string))
	}
	// Fresh slices: a ListBox built from the previous result may still
	// hold them.
	r.rows = r.filter(o.Items, query)
	r.sections = nil
	for _, group := range o.Groups {
		r.groupItems = append(r.groupItems, group.Items)
		if matches := r.filter(group.Items, query); len(matches) != 0 {
			r.sections = append(r.sections, list.ListBoxSection{Key: group.Key, Label: group.Label, Items: matches})
		}
	}
	return r
}

// same reports whether o's items and groups are the slices, of the same
// lengths, that the result was filtered from, at the same revision. It keeps
// the slices themselves, so the collector cannot give a new slice the
// address of a remembered one.
func (r *commandResults) same(o CommandOptions, revision uint64) bool {
	if revision != r.revision || !sameSlice(o.Items, r.items) || !sameSlice(o.Groups, r.groups) {
		return false
	}
	for i := range o.Groups {
		if !sameSlice(o.Groups[i].Items, r.groupItems[i]) {
			return false
		}
	}
	return true
}

func (r *commandResults) filter(source []CommandItem, query string) []list.ListBoxItem {
	var items []list.ListBoxItem
	for _, item := range source {
		if query != "" && !strings.Contains(strings.ToLower(item.Label+" "+item.Keywords), query) {
			continue
		}
		key := item.Key
		if key == "" {
			key = item.Label
		}
		if !item.Disabled {
			r.keys = append(r.keys, key)
		}
		var end func()
		if shortcut := item.Shortcut; shortcut != "" {
			draw := r.shortcut
			end = func() { (*draw)(shortcut) }
		}
		items = append(items, list.ListBoxItem{Key: key, Label: item.Label, Disabled: item.Disabled, End: end})
	}
	return items
}

// sameSlice reports whether a and b are the same elements of the same array.
func sameSlice[T any](a, b []T) bool {
	return len(a) == len(b) && (len(a) == 0 || &a[0] == &b[0])
}
