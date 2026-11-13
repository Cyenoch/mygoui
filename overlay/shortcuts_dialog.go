package overlay

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Shortcut is one row of a shortcut sheet: what the keys do and the chord.
type Shortcut struct {
	Label string
	// Hotkey is the chord as the platform names it: glyphs on macOS, words
	// elsewhere, so a sheet built from the keys the application binds
	// cannot drift from them.
	Hotkey form.Hotkey
	// Keys, when set, replaces Hotkey with explicit key caps in order, for
	// a chord that is not a single key press, such as "G then I".
	Keys []string
}

func (s Shortcut) caps() []string {
	if s.Keys != nil {
		return s.Keys
	}
	return s.Hotkey.Keys()
}

// ShortcutGroup is a titled section of shortcuts.
type ShortcutGroup struct {
	Title string
	Items []Shortcut
}

// ShortcutsDialogState is the caller-owned state of a shortcut sheet: whether
// it shows and the search text, which each opening clears. Keep one per
// window.
type ShortcutsDialogState struct {
	Open  bool
	Query string

	shown bool
}

// ShortcutsDialogOptions configures a shortcut sheet.
type ShortcutsDialogOptions struct {
	Title     string
	Groups    []ShortcutGroup
	OnDismiss func()
}

// shortcutsTwoColumns is the window width from which the groups sit in two
// columns.
const shortcutsTwoColumns = 800

// ShortcutsDialog lists keyboard shortcuts by group in a dialog while
// state.Open is true, in two columns on a wide window. Its search field
// takes the focus and narrows the rows to those whose label or keys contain
// every word typed, whatever the case ("ctrl k", "cmd"), marking the match in
// the labels; the title and the field stay in place while the rows scroll.
// The rows are built whole, which suits the few hundred a shortcut sheet
// holds.
func ShortcutsDialog(c *ui.Context, state *ShortcutsDialogState, o ShortcutsDialogOptions) ui.Element {
	if state == nil {
		panic("overlay: ShortcutsDialog requires a caller-owned state")
	}
	if state.Open && !state.shown {
		state.Query = ""
	}
	state.shown = state.Open
	w, _ := c.Size()
	columns := 1
	var width float32
	if w >= shortcutsTwoColumns {
		columns, width = 2, 720
	}
	return Dialog(c, &state.Open, DialogOptions{
		Title:     o.Title,
		Width:     width,
		OnDismiss: o.OnDismiss,
		Header:    func() { shortcutsHeader(c, state, o.Title) },
		Body:      func() { shortcutsBody(c, o.Groups, strings.Fields(strings.ToLower(state.Query)), columns) },
	})
}

func shortcutsHeader(c *ui.Context, state *ShortcutsDialogState, title string) {
	if title != "" {
		style.Text(ui.Text(c, title), c, theme.H4()).FontWeight(theme.WeightBold).
			Role(ui.RoleHeading).Level(2)
	}
	// The header keeps clear of the close button beside the title; the
	// field sits below the button and takes that room back.
	overlap := theme.Space(2) + theme.Scale(c, closeButtonSize()) - theme.Space(6)
	ui.Box(c).FillWidth().MinWidth(0).Margin(theme.Space(2), 0, 0, 0).MarginEnd(-max(overlap, 0)).Children(func() {
		form.SearchField(c, &state.Query, form.SearchFieldOptions{
			Placeholder: locale.M().Search, FullWidth: true,
		})
	})
}

// shortcutsBodyHeight keeps the dialog from collapsing around a short result
// as the search narrows it.
const shortcutsBodyHeight = 240

// shortcutHits is a group with the indices of its rows that match.
type shortcutHits struct {
	group int
	items []int
}

func (h shortcutHits) weight() int { return len(h.items) + 1 }

func shortcutsBody(c *ui.Context, groups []ShortcutGroup, words []string, columns int) {
	var hits []shortcutHits
	for i, group := range groups {
		var items []int
		for j := range group.Items {
			if group.Items[j].matches(words) {
				items = append(items, j)
			}
		}
		if len(items) > 0 {
			hits = append(hits, shortcutHits{group: i, items: items})
		}
	}
	list := ui.Column(c).FillWidth().MinWidth(0).MinHeight(theme.Scale(c, shortcutsBodyHeight)).
		Padding(theme.Space(2), 0, 0, 0)
	list.Children(func() {
		if len(hits) == 0 {
			mygoui.EmptyState(c, mygoui.EmptyStateOptions{Description: locale.M().NoResults})
			return
		}
		if columns < 2 {
			shortcutColumn(c, groups, hits, words)
			return
		}
		// Each group goes to the shorter column, so the two stay level.
		var split [2][]shortcutHits
		var weight [2]int
		for _, h := range hits {
			side := 0
			if weight[1] < weight[0] {
				side = 1
			}
			split[side] = append(split[side], h)
			weight[side] += h.weight()
		}
		ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space(8)).Children(func() {
			for i := range split {
				ui.Column(c.Key(i)).Grow(1).Basis(0).MinWidth(0).Children(func() {
					shortcutColumn(c, groups, split[i], words)
				})
			}
		})
	})
}

func shortcutColumn(c *ui.Context, groups []ShortcutGroup, hits []shortcutHits, words []string) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space(4)).Children(func() {
		for _, h := range hits {
			group := groups[h.group]
			section := ui.Column(c.Key(h.group)).FillWidth().MinWidth(0).Role(ui.RoleList).Label(group.Title)
			section.Children(func() {
				if group.Title != "" {
					ts := theme.XS()
					ts.Weight = theme.WeightMedium
					style.Text(ui.Text(c, group.Title), c, ts).TextColor(t.MutedForeground).
						Margin(0, 0, theme.Space(1), 0)
				}
				for _, j := range h.items {
					shortcutRow(c.Key(j), group.Items[j], words)
				}
			})
		}
	})
}

func shortcutRow(c *ui.Context, item Shortcut, words []string) {
	row := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Space(4)).
		MinHeight(theme.Scale(c, theme.Space(8))).Role(ui.RoleListItem)
	row.Children(func() {
		mygoui.Highlight(c, mygoui.HighlightOptions{Text: item.Label, Queries: words, IgnoreCase: true}).
			Grow(1).Basis(0).MinWidth(0)
		if caps := item.caps(); len(caps) > 0 {
			mygoui.KbdGroup(c, mygoui.KbdOptions{Keys: caps}).AlignSelf(ui.Center)
		}
	})
}

// capAliases are the words a user types for a macOS modifier glyph.
var capAliases = map[string]string{"⌘": "cmd command", "⌃": "ctrl control", "⌥": "alt option", "⇧": "shift"}

// matches reports whether every word, already lowercase, is in the label or
// the keys: in a cap, a word for a modifier glyph, or the caps run together
// as typed ("⌘k").
func (s Shortcut) matches(words []string) bool {
	if len(words) == 0 {
		return true
	}
	caps := s.caps()
	joined := ""
	for _, w := range words {
		if containsFold(s.Label, w) || capsContain(caps, w) {
			continue
		}
		if joined == "" {
			joined = strings.Join(caps, "")
		}
		if !containsFold(joined, w) {
			return false
		}
	}
	return true
}

func capsContain(caps []string, word string) bool {
	for _, k := range caps {
		if containsFold(k, word) || containsFold(capAliases[k], word) {
			return true
		}
	}
	return false
}

// containsFold reports whether s contains lower, which is already lowercase,
// without allocating: the rows are matched in every frame the dialog shows.
func containsFold(s, lower string) bool {
	for {
		if hasPrefixFold(s, lower) {
			return true
		}
		if s == "" {
			return false
		}
		_, size := utf8.DecodeRuneInString(s)
		s = s[size:]
	}
}

func hasPrefixFold(s, lower string) bool {
	for _, want := range lower {
		got, size := utf8.DecodeRuneInString(s)
		if size == 0 || unicode.ToLower(got) != want {
			return false
		}
		s = s[size:]
	}
	return true
}
