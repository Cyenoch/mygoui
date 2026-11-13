package navigation

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// menuTypeWindow is how long typed letters keep building one prefix: the
// toolkit's pause (MyGo ui/input.go, typePause), as a tree's and a list's.
const menuTypeWindow = time.Second

// menuAhead moves the focus to the enabled item whose label starts with
// what was typed: a letter pressed again steps through the items of that
// letter, and several letters in quick succession narrow to one prefix. The
// input handler only notes the item, since the focus is given while the rows
// are built.
type menuAhead struct {
	labels  []string // the enabled items' labels, in order, of the last build
	focused int      // the enabled item holding the focus, -1 for none
	prefix  string
	at      time.Time
	want    int // the item to focus in the next build, -1 for none
	n       int // the enabled item being built
}

// begin starts a build of the items.
func (a *menuAhead) begin() {
	a.labels, a.focused, a.n = a.labels[:0], -1, 0
}

// item notes the next enabled item and focuses it when it was asked for.
func (a *menuAhead) item(label string, el ui.Element) {
	a.labels = append(a.labels, label)
	if a.want == a.n {
		el.Focus()
		a.want = -1
	}
	if el.Focused() {
		a.focused = a.n
	}
	el.HandleInput(a.input)
	a.n++
}

func (a *menuAhead) input(ev ui.InputEvent) bool {
	if ev.Kind != ui.InputText || ev.Mods&(ui.Ctrl|ui.Alt|ui.Cmd) != 0 {
		return false
	}
	typed := strings.ToLower(ev.Text)
	if strings.TrimSpace(typed) == "" {
		return false
	}
	if ev.Time.Sub(a.at) > menuTypeWindow {
		a.prefix = ""
	}
	a.at = ev.Time
	// A first letter, or the same letter again, steps along the items that
	// start with it; more letters narrow the prefix from the current item.
	from := a.focused + 1
	if a.prefix == "" || a.prefix == strings.Repeat(typed, len([]rune(a.prefix))) {
		a.prefix = typed
	} else {
		a.prefix += typed
		from = max(a.focused, 0)
	}
	for i := range a.labels {
		at := (from + i) % len(a.labels)
		if strings.HasPrefix(strings.ToLower(a.labels[at]), a.prefix) {
			a.want = at
			return true
		}
	}
	return true
}
