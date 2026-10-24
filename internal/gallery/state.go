// Package gallery is the interactive native workbench for mygoui.
package gallery

import (
	"sort"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
)

// Page is one registered specimen page. Pages register themselves so family
// files stay independent and the shell never imports their internals.
type Page struct {
	ID          string
	Title       string
	Description string
	Group       string
	Order       int
	Build       func(c *ui.Context, s *State) ui.Element
}

var (
	pagesMu sync.RWMutex
	pages   []Page
)

// RegisterPage adds a page to the catalog. Call it from a page file's init.
func RegisterPage(p Page) {
	pagesMu.Lock()
	defer pagesMu.Unlock()
	for _, existing := range pages {
		if existing.ID == p.ID {
			panic("gallery: duplicate page " + p.ID)
		}
	}
	pages = append(pages, p)
	sort.SliceStable(pages, func(i, j int) bool {
		if pages[i].Order != pages[j].Order {
			return pages[i].Order < pages[j].Order
		}
		return pages[i].Title < pages[j].Title
	})
}

// Pages returns the registered pages.
func Pages() []Page {
	pagesMu.RLock()
	defer pagesMu.RUnlock()
	out := make([]Page, len(pages))
	copy(out, pages)
	return out
}

// State is the window-owned gallery state. Page state is typed and stored by
// page ID, so sibling pages do not share a struct or initialization order.
type State struct {
	Selected string
	// ActiveComponent is the catalog entry the navigation highlights; its
	// PageID is what the content pane shows.
	ActiveComponent string
	Dark            bool
	Search          string
	Accent          int
	Radius          float32
	NavOpen         bool
	Sidebar         ui.ScrollState
	NavScroll       ui.ScrollState
	Content         ui.ScrollState
	pages           map[string]any
}

// New returns the default gallery state.
func New() *State {
	s := &State{Dark: false, Selected: "foundations", Radius: 8, pages: map[string]any{}}
	if list := Pages(); len(list) > 0 {
		s.Selected = list[0].ID
	}
	return s
}

// StateOf returns the typed state of a page, creating it on first use.
func StateOf[T any](s *State, page string) *T {
	if s.pages == nil {
		s.pages = map[string]any{}
	}
	if value, ok := s.pages[page].(*T); ok {
		return value
	}
	value := new(T)
	s.pages[page] = value
	return value
}

// filterPages returns the pages matching the current search. An empty query
// returns every page.
func (s *State) filterPages() []Page {
	query := strings.ToLower(strings.TrimSpace(s.Search))
	if query == "" {
		return Pages()
	}
	var out []Page
	for _, page := range Pages() {
		haystack := strings.ToLower(page.Title + " " + page.Description + " " + page.Group)
		if strings.Contains(haystack, query) {
			out = append(out, page)
		}
	}
	return out
}
