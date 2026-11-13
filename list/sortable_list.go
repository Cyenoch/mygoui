package list

import (
	"fmt"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SortableItem is one row of a SortableList. ID is its stable identity: the
// row keeps its focus and state as the order changes. Content replaces the
// Label text, which still names the row's handle.
type SortableItem struct {
	ID      string
	Label   string
	Content func()
}

// SortableListOptions configures a reorderable list. Items are caller-owned
// and never modified: OnMove reports that the item at from goes to index to,
// counted after its removal, and the caller applies it.
type SortableListOptions struct {
	Items    []SortableItem
	Label    string
	Disabled bool
	OnMove   func(from, to int)
}

type (
	sortStateKey struct{}
	sortShiftKey struct{}
)

// sortState is the drag in progress: the row held, how far the pointer took
// it, and the rows' boxes from the frame before.
type sortState struct {
	active int
	// activeID is the held item, so a changed list cancels the drag.
	activeID string
	dy       float32
	bounds   []ui.Rect
	// settle places the rows at once in the build after a drop, where the
	// caller's new order already holds them.
	settle bool
	// cancelled is the item whose drag Escape cancelled while the pointer
	// still holds it, so that it neither drags again nor drops.
	cancelled string
}

// SortableList shows rows the user reorders by dragging a row's grip handle,
// the other rows making room as it passes, or with Alt+Up and Alt+Down while
// the handle has the focus. Escape cancels a drag. Each row tells assistive
// technology its place, the grip its keys, and a move is read out.
func SortableList(c *ui.Context, o SortableListOptions) ui.Element {
	gap := theme.Scale(c, theme.Space(2))
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(gap).Role(ui.RoleList)
	if o.Label != "" {
		root.Label(o.Label)
	}
	st := ui.Local(root, sortStateKey{}, func() sortState { return sortState{active: -1} })
	n := len(o.Items)
	if st.active >= n || len(st.bounds) != n || o.Disabled ||
		(st.active >= 0 && o.Items[st.active].ID != st.activeID) {
		st.active = -1
	}
	over := st.active
	if st.active >= 0 {
		st.dy = sortClamp(st.bounds, st.active, st.dy)
		over = sortOver(st.bounds, st.active, st.dy)
	}
	settle := st.settle
	st.settle = false
	if len(st.bounds) != n {
		st.bounds = make([]ui.Rect, n)
	}
	from, to := -1, -1
	root.Children(func() {
		for i, it := range o.Items {
			row := sortRow(c.Key(it.ID), i == st.active).PositionInSet(i+1, n)
			if it.Label != "" {
				row.Label(it.Label)
			}
			held, cancel, step := false, false, 0
			row.Children(func() {
				var dy float32
				dy, held, cancel, step = sortHandle(c, it, o.Disabled || o.OnMove == nil)
				if i == st.active {
					st.dy += dy
				}
				body := ui.Row(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center)
				body.Children(func() {
					if it.Content != nil {
						it.Content()
						return
					}
					style.Text(ui.Text(c, it.Label), c, theme.SM())
				})
			})
			shift := sortShiftOf(st.bounds, st.active, over, i, gap)
			switch {
			case i == st.active:
				row.Translate(0, st.dy).ZIndex(1)
			case settle:
				row.Translate(0, row.Jump(sortShiftKey{}, 0))
			default:
				m := theme.Motion()
				row.Translate(0, row.AnimateTo(sortShiftKey{}, shift, ui.Animation{Duration: m.Shift, Ease: m.Ease}))
			}
			switch {
			case st.cancelled == it.ID:
				if !held {
					st.cancelled = ""
				}
			case cancel && st.active == i:
				st.active, st.dy, st.cancelled = -1, 0, it.ID
				c.Invalidate()
			case held && st.active < 0:
				st.active, st.activeID, st.dy = i, it.ID, 0
			case !held && st.active == i:
				from, to = i, over
				st.active, st.dy, st.settle = -1, 0, true
				c.Invalidate()
			case step != 0 && st.active < 0 && i+step >= 0 && i+step < n:
				from, to = i, i+step
			}
			st.bounds[i] = row.LayoutBounds()
		}
	})
	// The caller reorders its items here, after the rows were built from them.
	if from != to && o.OnMove != nil {
		c.Announce(fmt.Sprintf(locale.M().MovedTo, o.Items[from].Label, to+1, n))
		o.OnMove(from, to)
	}
	return root
}

func sortRow(c *ui.Context, lifted bool) ui.Element {
	t := theme.Colors(c)
	row := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(1))).
		Padding(theme.Scale(c, theme.Space(1.5)), theme.Scale(c, theme.Space(3)), theme.Scale(c, theme.Space(1.5)), theme.Scale(c, theme.Space(1.5))).
		Radius(theme.RadiusScale().LG).Background(t.Card).TextColor(t.CardForeground).
		Border(theme.Hairline, t.Border).Role(ui.RoleListItem)
	if lifted {
		row = style.Shadow(row, c, theme.OverlayShadows(c.Theme().Dark))
	}
	return row
}

// sortHandle builds a row's grip and reports how far the pointer holding it
// moved, whether Escape cancelled that drag, and the step a key asked for.
func sortHandle(c *ui.Context, it SortableItem, disabled bool) (dy float32, held, cancel bool, step int) {
	t := theme.Colors(c)
	side := theme.Scale(c, theme.Space(7))
	radius := theme.RadiusScale().MD
	label := locale.M().Reorder
	if it.Label != "" {
		label += " " + it.Label
	}
	h := ui.ButtonBase(c).Size(side, side).Center().Shrink(0).FocusRing(true).
		Label(label).TextColor(t.MutedForeground).Cursor(ui.CursorGrab).
		KeyShortcuts(ui.Alt, ui.KeyUp).KeyShortcuts(ui.Alt, ui.KeyDown)
	h = style.Disabled(h, disabled)
	h = style.Background(h, radius, ui.Transparent, t.Accent, t.Accent)
	h.Children(func() {
		glyph := theme.Scale(c, theme.IconMD)
		ui.Icon(c, lucide.GripVertical()).Size(glyph, glyph).Shrink(0)
	})
	if disabled {
		return 0, false, false, 0
	}
	if _, dy, held = h.Dragged(); held {
		h.Cursor(ui.CursorGrabbing)
		return dy, true, h.Shortcut(0, ui.KeyEscape), 0
	}
	switch {
	case h.Shortcut(ui.Alt, ui.KeyUp):
		step = -1
	case h.Shortcut(ui.Alt, ui.KeyDown):
		step = 1
	}
	return 0, false, false, step
}

// sortClamp keeps the held row between the first and the last row.
func sortClamp(bounds []ui.Rect, active int, dy float32) float32 {
	a, first, last := bounds[active], bounds[0], bounds[len(bounds)-1]
	return min(max(dy, first.Y-a.Y), last.Y+last.H-a.Y-a.H)
}

// sortOver returns the place whose row center is nearest the held row's.
func sortOver(bounds []ui.Rect, active int, dy float32) int {
	center := bounds[active].Y + bounds[active].H/2 + dy
	best, far := active, float32(-1)
	for i, b := range bounds {
		d := b.Y + b.H/2 - center
		if d < 0 {
			d = -d
		}
		if far < 0 || d < far {
			best, far = i, d
		}
	}
	return best
}

// sortShiftOf is how far row i moves aside for the held row going to over.
func sortShiftOf(bounds []ui.Rect, active, over, i int, gap float32) float32 {
	if active < 0 {
		return 0
	}
	room := bounds[active].H + gap
	switch {
	case active < i && i <= over:
		return -room
	case over <= i && i < active:
		return room
	}
	return 0
}
