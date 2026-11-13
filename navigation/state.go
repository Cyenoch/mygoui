package navigation

import "github.com/egoist/mygo/ui"

// accordionScope is the expansion state one Accordion shares with its items,
// and the header its arrows move the focus to.
type accordionScope struct {
	variant       AccordionVariant
	multiple      bool
	expanded      *[]string
	hideSeparator bool
	onChange      func(key string, expanded bool)
	focus         *int
	enabled       func(i int) bool
}

func (s *accordionScope) set(key string, open bool) {
	keys := *s.expanded
	for i, k := range keys {
		if k != key {
			continue
		}
		if open {
			return
		}
		*s.expanded = append(keys[:i], keys[i+1:]...)
		return
	}
	if !open {
		return
	}
	if !s.multiple {
		keys = keys[:0]
	}
	*s.expanded = append(keys, key)
}

// disclosureGroupScope coordinates caller-owned open flags among the
// disclosures of one group, and the trigger its arrows move the focus to.
type disclosureGroupScope struct {
	multiple bool
	items    []disclosureItem
	focus    *int
	enabled  func(i int) bool
	// at is the disclosure being built, of n.
	at, n int
}

type disclosureItem struct {
	open    *bool
	changed bool
}

// rove gives header i of n its arrows: back and forth move the focus to the
// previous and next enabled header, round the ends, Home and End to the
// first and the last. *focus is the header to focus as it is built, -1 for
// none. Unlike a focus group, each header stays a Tab stop and what lies
// between them keeps its own Tab order, as APG's accordion has it.
func rove(el ui.Element, focus *int, i, n int, enabled func(int) bool, back, forth ui.Key) {
	if *focus == i {
		el.Focus()
		*focus = -1
	}
	find := func(from, step int) {
		for d := 1; d <= n; d++ {
			if j := ((from+step*d)%n + n) % n; enabled(j) {
				*focus = j
				return
			}
		}
	}
	el.OnShortcut(0, forth, func() { find(i, 1) })
	el.OnShortcut(0, back, func() { find(i, -1) })
	el.OnShortcut(0, ui.KeyHome, func() { find(-1, 1) })
	el.OnShortcut(0, ui.KeyEnd, func() { find(n, -1) })
}
