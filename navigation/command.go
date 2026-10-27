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

// CommandState holds the query, list highlight and focus for one command.
type CommandState struct {
	Query string
	Input ui.Handle
	List  list.ListBoxState

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
	filter := func(source []CommandItem, query string) []list.ListBoxItem {
		items := make([]list.ListBoxItem, 0, len(source))
		for _, item := range source {
			if query != "" && !strings.Contains(strings.ToLower(item.Label+" "+item.Keywords), query) {
				continue
			}
			key := item.Key
			if key == "" {
				key = item.Label
			}
			shortcut := item.Shortcut
			var end func()
			if shortcut != "" {
				end = func() { style.Text(ui.Text(c, shortcut), c, theme.XS()).TextColor(t.MutedForeground) }
			}
			items = append(items, list.ListBoxItem{Key: key, Label: item.Label, Disabled: item.Disabled, End: end})
		}
		return items
	}
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
		query := strings.ToLower(strings.TrimSpace(s.Query))
		items := filter(o.Items, query)
		groups := make([]list.ListBoxSection, 0, len(o.Groups))
		for _, group := range o.Groups {
			matches := filter(group.Items, query)
			if len(matches) != 0 {
				groups = append(groups, list.ListBoxSection{Key: group.Key, Label: group.Label, Items: matches})
			}
		}
		list.ListBox(c, list.ListBoxOptions{
			Label: label, Items: items, Sections: groups, State: &s.List, Padding: theme.Space8,
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
	keys := commandKeys(o, text)
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

func commandKeys(o CommandOptions, text string) []string {
	query := strings.ToLower(strings.TrimSpace(text))
	var keys []string
	add := func(items []CommandItem) {
		for _, item := range items {
			if item.Disabled || (query != "" &&
				!strings.Contains(strings.ToLower(item.Label+" "+item.Keywords), query)) {
				continue
			}
			key := item.Key
			if key == "" {
				key = item.Label
			}
			keys = append(keys, key)
		}
	}
	add(o.Items)
	for _, group := range o.Groups {
		add(group.Items)
	}
	return keys
}
