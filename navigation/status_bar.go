package navigation

import (
	"slices"
	"strconv"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatusBarTone is the face of a status bar: the neutral strip, or a warning
// or error face on the whole bar.
type StatusBarTone uint8

const (
	StatusBarNeutral StatusBarTone = iota
	StatusBarWarning
	StatusBarError
)

// StatusBarOptions configures a window's status strip. Message is the
// transient text at the leading edge: the caller sets it, and clears it when
// it no longer applies. Items are the permanent trailing items.
type StatusBarOptions struct {
	Message string
	Label   string
	Tone    StatusBarTone
	Items   []StatusBarItemOptions
}

// StatusBarProgress is a small progress track in an item: Value runs from 0
// to 1.
type StatusBarProgress struct {
	Value         float64
	Indeterminate bool
}

// StatusBarItemOptions configures one trailing item: text, an icon, a
// progress track, or any of them. OnClick makes it a button, whose Label or
// Text names it for assistive technology and whose tooltip is Tooltip, or
// Label. Label also names an icon standing alone, before the text when there
// is one. Separator draws a short rule between groups and ignores the rest.
//
// When the bar is too narrow, the message keeps a minimum width and items
// leave in order of Priority, the lowest first and the later of equal ones
// before the earlier; a rule stays only between items that remain.
type StatusBarItemOptions struct {
	Text      string
	Icon      *ui.SVG
	Label     string
	Tooltip   string
	Progress  *StatusBarProgress
	Priority  int
	Separator bool
	OnClick   func()
}

// statusBarMessageMin is the least width of the message, before scaling.
const statusBarMessageMin = 96

// StatusBar builds a compact strip for the bottom of a window. The message
// is a status region, so assistive technology reads each change of it; the
// trailing items are one focus stop whose Left and Right arrows move among
// the pressable ones. The row measures itself from the layout of the frame
// before, so a resize settles in the next frame; meanwhile it clips what
// does not fit.
func StatusBar(c *ui.Context, o StatusBarOptions) ui.Element {
	t := theme.Colors(c)
	face, ink := t.MutedBackground, legibleInk(t.MutedForeground, t.Foreground, t.MutedBackground.Over(t.Background))
	switch o.Tone {
	case StatusBarWarning:
		face = t.Warning
		ink = readableOn(face)
	case StatusBarError:
		face = t.Destructive
		ink = readableOn(face)
	}
	pad, gap := scale(c, theme.Space8), scale(c, theme.Space8)
	root := ui.Row(c).FillWidth().MinWidth(0).Shrink(0).AlignItems(ui.Center).Clip().
		Height(scale(c, theme.Space24)).PaddingX(pad).Gap(gap).
		Background(face).TextColor(ink).
		BorderWidth(hairline, 0, 0, 0).BorderColor(t.Border).Role(ui.RoleGroup)
	if o.Label != "" {
		root.Label(o.Label)
	}
	avail := root.LayoutBounds().W - 2*pad
	root.DrawOver(func(_ *ui.Painter, r ui.Rect) { invalidateOnResize(c, avail+2*pad, r.W) })
	root.Children(func() {
		// The region is built while empty too: a live region reads the
		// changes of an element that was already there.
		style.Text(ui.Text(c, o.Message), c, theme.XS()).SingleLine().
			Grow(1).Basis(0).MinWidth(scale(c, statusBarMessageMin)).Role(ui.RoleStatus).
			ContentTransition(ui.ContentTransition{})
		if len(o.Items) == 0 {
			return
		}
		row := ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(scale(c, theme.Space2)).FocusGroup(ui.Horizontal)
		row.Children(func() {
			items := make([]ui.Element, len(o.Items))
			widths := make([]float32, len(o.Items))
			keys := statusBarKeys(o.Items)
			for i := range o.Items {
				items[i] = buildStatusBarItem(c.Key(keys[i]), o.Items[i], t, ink, o.Tone)
				widths[i] = items[i].LayoutBounds().W
				w := widths[i]
				items[i].DrawOver(func(_ *ui.Painter, r ui.Rect) { invalidateOnResize(c, w, r.W) })
			}
			room := avail - gap - scale(c, statusBarMessageMin)
			shown := statusBarFit(o.Items, widths, room, scale(c, theme.Space2))
			for i, show := range shown {
				if !show {
					items[i].Absolute().Invisible()
				}
			}
		})
	})
	return root
}

// statusBarKeys is the identity of each item: its Label, else its Text, else
// its place, which a repeated identity falls back to as well.
func statusBarKeys(items []StatusBarItemOptions) []string {
	keys := make([]string, len(items))
	seen := make(map[string]bool, len(items))
	for i, it := range items {
		key := it.Label
		if key == "" {
			key = it.Text
		}
		if key == "" || seen[key] {
			key = "#" + strconv.Itoa(i)
		}
		seen[key], keys[i] = true, key
	}
	return keys
}

// statusBarFit returns which items show in room: the lowest priorities leave
// first until the rest fit, and a separator shows only between items that
// show. Every item shows while a width is still unknown.
func statusBarFit(items []StatusBarItemOptions, widths []float32, room, gap float32) []bool {
	shown := make([]bool, len(items))
	for i, w := range widths {
		shown[i] = true
		if w <= 0 || room <= 0 {
			return shown
		}
	}
	separator := make([]bool, len(items))
	for i := range items {
		separator[i] = items[i].Separator
	}
	// order lists the items by when they leave.
	order := make([]int, 0, len(items))
	for i := range items {
		if !items[i].Separator {
			order = append(order, i)
		}
	}
	slices.SortStableFunc(order, func(a, b int) int {
		if items[a].Priority != items[b].Priority {
			return items[a].Priority - items[b].Priority
		}
		return b - a
	})
	for leaving := 0; ; leaving++ {
		settleSeparators(separator, shown)
		total, count := float32(0), 0
		for i, show := range shown {
			if show {
				total += widths[i]
				count++
			}
		}
		if total+gap*float32(count-1) <= room+0.5 || leaving == len(order) {
			return shown
		}
		shown[order[leaving]] = false
	}
}

// settleSeparators shows each separator that has a shown item before it and
// after it, with no other separator between, and hides the others. Items
// that are separators are marked in separator; the rest of shown is theirs.
func settleSeparators(separator, shown []bool) {
	lastSeparator, sawItem := -1, false
	for i := range shown {
		switch {
		case !separator[i]:
			if shown[i] {
				sawItem = true
				if lastSeparator >= 0 {
					shown[lastSeparator] = true
				}
				lastSeparator = -1
			}
		case sawItem && lastSeparator < 0:
			lastSeparator = i
			shown[i] = false
		default:
			shown[i] = false
		}
	}
}

// buildStatusBarItem builds one trailing item.
func buildStatusBarItem(c *ui.Context, o StatusBarItemOptions, t theme.Palette, ink ui.Color, tone StatusBarTone) ui.Element {
	if o.Separator {
		return ui.Box(c).Width(hairline).Height(scale(c, theme.Space12)).Shrink(0).
			Margin(0, scale(c, theme.Space4)).Background(ink.Alpha(0.32)).Role(ui.RoleNone)
	}
	var item ui.Element
	if o.OnClick != nil {
		item = ui.ButtonBase(c).FocusRing(true)
		hover := t.Accent
		if tone != StatusBarNeutral {
			hover = ink.Alpha(0.18)
		}
		item = style.Background(item, theme.RadiusScale().SM, ui.Transparent, hover, hover)
		if item.Hovered() && tone == StatusBarNeutral {
			item.TextColor(t.Foreground)
		}
		item.OnClick(o.OnClick)
	} else {
		item = ui.Box(c)
	}
	item.Row().AlignItems(ui.Center).Shrink(0).Gap(scale(c, theme.Space4)).
		Height(scale(c, theme.Space20)).PaddingX(scale(c, theme.Space6))
	// A text item is read as its text; a name is for the button, and for
	// an icon standing alone.
	switch {
	case o.Label != "" && o.Text != "":
		item.Label(o.Label + ", " + o.Text)
	case o.Label != "":
		item.Label(o.Label)
	case o.OnClick != nil:
		item.Label(o.Text)
	}
	switch {
	case o.Tooltip != "":
		item.Tooltip(o.Tooltip)
	case o.OnClick != nil && o.Label != "":
		item.Tooltip(o.Label)
	}
	item.Children(func() {
		if o.Icon != nil {
			glyph := scale(c, theme.IconXS)
			ui.Icon(c, o.Icon).Size(glyph, glyph).Shrink(0)
		}
		if o.Text != "" {
			style.Text(ui.Text(c, o.Text), c, theme.XS()).SingleLine().FontFeatures("tnum")
		}
		if p := o.Progress; p != nil {
			// A toned face keeps one ink in both appearances, so the track
			// follows the ink, not the window's appearance.
			onFace := *c.Theme()
			onFace.Dark = ink == ui.RGB(255, 255, 255) || (tone == StatusBarNeutral && onFace.Dark)
			ui.Themed(c, &onFace).Background(ui.Transparent).Width(scale(c, 64)).Shrink(0).Children(func() {
				mygoui.ProgressBar(c, mygoui.ProgressOptions{
					Value: p.Value, Max: 1, Indeterminate: p.Indeterminate, Size: mygoui.ProgressSmall,
					Color: mygoui.IndicatorDefault,
				}).Label(statusBarProgressName(o))
			})
		}
	})
	return item
}

// statusBarProgressName names a progress track by its item.
func statusBarProgressName(o StatusBarItemOptions) string {
	if o.Label != "" {
		return o.Label
	}
	return o.Text
}
