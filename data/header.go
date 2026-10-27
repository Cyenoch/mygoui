package data

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type headerGeometryKey struct{}

// MyGo 0.7.1 shares header and row colors. Restore header affordances using
// native cells' committed column geometry, without changing native input.
func (tl *tableLayout) decorateHeader(c *ui.Context, table ui.Element) {
	ready := ui.Local(table, headerGeometryKey{}, func() bool { return false })
	if !*ready {
		*ready = true
		c.AnimationFrame()
	}
	cols := tl.nativeColumns()
	x, y, over := table.PointerPosition()
	bounds := table.LayoutBounds()
	px := bounds.X + x
	hover, grip, border := c.Theme().SurfaceHover, c.Theme().Accent, c.Theme().Border
	rtl := c.Direction() == ui.RightToLeft
	focus := table.FocusVisible()
	chrome := tableChrome(tl.o, tl.colors, true)
	table.Draw(func(p *ui.Painter, r ui.Rect) {
		tl.body = tableBody(r, true)
		chrome(p, r)
		if !over || y < 0 || y >= theme.Space(8) {
			return
		}
		p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: theme.Space(8)}, 0, func() {
			for i, box := range tl.headers {
				if cols[i].Sortable && px >= box.X && px < box.X+box.W {
					p.Fill(ui.Rect{X: box.X, Y: r.Y, W: box.W, H: theme.Space(8)}, hover, 0)
				}
			}
		})
	})
	table.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focus {
			p.FocusRing(r, [4]float32{})
		}
		p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: theme.Space(8)}, 0, func() {
			for i, box := range tl.headers {
				col := cols[i]
				if box.W == 0 || col.Fixed {
					continue
				}
				end := box.X + box.W
				lo, hi := end-8, end
				if rtl {
					end = box.X + 1
					lo, hi = box.X, box.X+8
				}
				// Native tables omit the outer edge of a final flex column.
				if col.Width == 0 && tl.st.List.Columns.Widths[col.ID] == 0 &&
					(!rtl && end >= r.X+r.W || rtl && box.X <= r.X) {
					continue
				}
				color := border
				if over && y >= 0 && y < theme.Space(8) && px >= lo && px < hi {
					color = grip
				}
				p.Fill(ui.Rect{X: end - 1, Y: r.Y + theme.Space(8)/4, W: 1, H: theme.Space(8) / 2}, color, 0)
			}
		})
	})
}
