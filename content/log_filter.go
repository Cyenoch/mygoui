package content

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

// LogLevels is a set of levels for LogViewerState.Levels: the listing
// narrows to the levels in it. The zero set shows every level.
type LogLevels uint8

// LevelSet returns the set holding levels.
func LevelSet(levels ...LogLevel) LogLevels {
	var s LogLevels
	for _, l := range levels {
		s |= 1 << l
	}
	return s
}

// Shows reports whether a line of level l is listed: any level in the zero
// set. A plain line follows the Info level.
func (s LogLevels) Shows(l LogLevel) bool {
	if l == LogPlain {
		l = LogInfo
	}
	return s == 0 || s.has(l)
}

func (s LogLevels) has(l LogLevel) bool { return s&(1<<l) != 0 }

// toggled returns the set with level l added or removed.
func (s LogLevels) toggled(l LogLevel, on bool) LogLevels {
	if on {
		return s | 1<<l
	}
	return s &^ (1 << l)
}

// logFilter is the index of the lines a query and a level set list, kept
// with what it was built from, so a frame at rest scans nothing and a log
// that grew is scanned from where the last scan stopped.
type logFilter struct {
	ready    bool
	query    string // lower case
	levels   LogLevels
	revision int
	scanned  int
	rows     []int // original index of each listed line
}

// filtering reports whether the state lists fewer lines than the log holds.
func (s *LogViewerState) filtering() bool { return s.Query != "" || s.Levels != 0 }

// count is the number of lines listed of the n the log holds.
func (s *LogViewerState) count(n int) int {
	if !s.filtering() {
		return n
	}
	return len(s.filter.rows)
}

// original is the index in the log of the listed line row.
func (s *LogViewerState) original(row int) int {
	if !s.filtering() {
		return row
	}
	return s.filter.rows[row]
}

// refresh brings the listing up to date with the state and the log's n
// lines, takes the list to where the new listing begins when the query or
// the levels changed it, and returns the number of lines listed.
func (s *LogViewerState) refresh(o *LogViewerOptions, n int) int {
	if s.sync(o, n) {
		if o.Follow {
			s.List.ScrollToEnd()
		} else {
			s.List.ScrollTo(0, ui.Start)
		}
	}
	return s.count(n)
}

// sync brings the index up to date with the query, the levels, the
// revision and the n lines of the log, and reports whether the query or the
// levels changed, which takes the list back to the start. A log that was
// replaced or shrank is listed again in place.
func (s *LogViewerState) sync(o *LogViewerOptions, n int) (changed bool) {
	f := &s.filter
	if !s.filtering() {
		changed = f.ready
		*f = logFilter{}
		return changed
	}
	query := strings.ToLower(s.Query)
	changed = !f.ready || f.levels != s.Levels || query != f.query
	switch {
	case !f.ready || f.revision != s.Revision || n < f.scanned || f.levels != s.Levels ||
		!strings.Contains(query, f.query):
		f.rows, f.scanned = f.rows[:0], 0
	case query != f.query:
		// A longer query lists a subset of the lines listed so far.
		kept := f.rows[:0]
		for _, i := range f.rows {
			if containsFold(o.line(i), query) {
				kept = append(kept, i)
			}
		}
		f.rows = kept
	}
	f.ready, f.query, f.levels, f.revision = true, query, s.Levels, s.Revision
	for i := f.scanned; i < n; i++ {
		if s.Levels != 0 {
			level := LogPlain
			if o.Level != nil {
				level = o.Level(i)
			}
			if !s.Levels.Shows(level) {
				continue
			}
		}
		if containsFold(o.line(i), query) {
			f.rows = append(f.rows, i)
		}
	}
	f.scanned = n
	return changed
}

// containsFold reports whether s holds lower, a lower-case string, whatever
// the case of s, without allocating.
func containsFold(s, lower string) bool {
	if lower == "" {
		return true
	}
	for i := 0; i < len(s); {
		if hasPrefixFold(s[i:], lower) {
			return true
		}
		_, w := utf8.DecodeRuneInString(s[i:])
		i += w
	}
	return false
}

func hasPrefixFold(s, lower string) bool {
	for lower != "" {
		if s == "" {
			return false
		}
		a, w := utf8.DecodeRuneInString(s)
		b, v := utf8.DecodeRuneInString(lower)
		if unicode.ToLower(a) != b {
			return false
		}
		s, lower = s[w:], lower[v:]
	}
	return true
}
