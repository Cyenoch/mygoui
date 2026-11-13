// Package gallery is the interactive native workbench for mygoui.
package gallery

import (
	"github.com/Cyenoch/mygoui"
	"sort"
	"strings"
	"sync"

	"github.com/egoist/mygo"
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
	Selected           string
	Dark               bool
	OnAppearanceChange func(bool)
	Window             *mygo.Window
	// Flyouts opens anchored popups as native flyouts of Window.
	Flyouts       mygoui.FlyoutHost
	zoomStep      int
	Search        string
	Primary       int
	Radius        float32
	NavOpen       bool
	SidebarClosed bool
	// areaLaid is the width the content area had in the last frame's layout.
	areaLaid       float32
	AppearanceOpen bool
	SystemFont     bool
	// SystemFontOnly records that Cal Sans could not be registered, as in
	// browsers: the gallery uses the system font and offers no font choice.
	SystemFontOnly bool
	SearchFocus    ui.Handle
	Sidebar        ui.ScrollState
	Content        ui.ScrollState
	pages          map[string]any
	sidebarGroups  groupCache
	nav            navMetrics
	overviewGroups groupCache
}

func (s *State) setDark(dark bool) {
	if s.Dark == dark {
		return
	}
	s.Dark = dark
	if s.OnAppearanceChange != nil {
		s.OnAppearanceChange(dark)
	}
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
// groupCache keeps the page groups of one search. The registered pages are
// fixed after init, so only the search and the page count (tests register
// pages) invalidate it.
type groupCache struct {
	query  string
	count  int
	groups []pageGroup
	valid  bool
}

// of returns the cached groups for query, building them when the search or
// the registered pages changed.
func (g *groupCache) of(query string, build func() []pageGroup) []pageGroup {
	pagesMu.RLock()
	count := len(pages)
	pagesMu.RUnlock()
	if !g.valid || g.query != query || g.count != count {
		g.query, g.count, g.groups, g.valid = query, count, build(), true
	}
	return g.groups
}

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
