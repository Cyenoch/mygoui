package navigation

// accordionScope is the expansion state one Accordion shares with its items.
type accordionScope struct {
	variant       AccordionVariant
	multiple      bool
	expanded      *[]string
	hideSeparator bool
	onChange      func(key string, expanded bool)
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
// disclosures of one group.
type disclosureGroupScope struct {
	multiple bool
	items    []disclosureItem
}

type disclosureItem struct {
	open    *bool
	changed bool
}
