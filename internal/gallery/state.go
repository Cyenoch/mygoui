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

	// searchKey is the lowercased title, description and group the sidebar
	// search matches, and navRank is the group's navigation position. Both are
	// prepared at registration so no build pass derives them again.
	searchKey string
	navRank   int
}

var (
	pagesMu   sync.RWMutex
	pages     []Page
	pagesByID = map[string]Page{}
)

// RegisterPage adds a page to the catalog. Call it from a page file's init.
func RegisterPage(p Page) {
	p.searchKey = strings.ToLower(p.Title + " " + p.Description + " " + p.Group)
	p.navRank = navRankOf(p.Group)
	pagesMu.Lock()
	defer pagesMu.Unlock()
	if _, ok := pagesByID[p.ID]; ok {
		panic("gallery: duplicate page " + p.ID)
	}
	pagesByID[p.ID] = p
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

// pageFor returns the page with id.
func pageFor(id string) (Page, bool) {
	pagesMu.RLock()
	defer pagesMu.RUnlock()
	p, ok := pagesByID[id]
	return p, ok
}

// State is the window-owned gallery state. Page state is typed and stored by
// page ID, so sibling pages do not share a struct or initialization order.
type State struct {
	Selected       string
	Dark           bool
	Search         string
	OpenGroup      string
	Primary        int
	Radius         float32
	NavOpen        bool
	SidebarClosed  bool
	AppearanceOpen bool
	SystemFont     bool
	Panel          int
	SearchFocus    ui.Handle
	Sidebar        ui.ScrollState
	NavScroll      ui.ScrollState
	Content        ui.ScrollState
	pages          map[string]any
}

// New returns the default gallery state.
func New() *State {
	s := &State{Dark: false, Selected: "foundations", Radius: 10, pages: map[string]any{}}
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
// returns every page. The result is one fresh slice the caller may sort; the
// per-page haystack was lowercased at registration, so a build pass only
// compares it.
func (s *State) filterPages() []Page {
	query := strings.ToLower(strings.TrimSpace(s.Search))
	pagesMu.RLock()
	defer pagesMu.RUnlock()
	out := make([]Page, 0, len(pages))
	for _, page := range pages {
		if query == "" || strings.Contains(page.searchKey, query) {
			out = append(out, page)
		}
	}
	return out
}
