package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DescriptionsItem is one labelled value. Content replaces the Value text
// with caller-built content, as a badge or a link. Span is the number of
// columns the item takes; values below one take one.
type DescriptionsItem struct {
	Label, Value string
	Content      func()
	Span         int
}

// DescriptionsOptions configures a grid of labelled values. Columns is the
// number of items in a row, three when zero; the grid drops to fewer when
// its width leaves a column less than MinColumnWidth, 200 DIPs when zero.
// Bordered draws the values as a ruled table with the labels in aligned,
// tinted columns; otherwise each label sits above its value. A value that
// is empty shows a muted dash. Values wrap; Truncate cuts each at one line
// and shows the whole in a tooltip. Each item is a group named by its label.
type DescriptionsOptions struct {
	Title          string
	Items          []DescriptionsItem
	Columns        int
	MinColumnWidth float32
	Bordered       bool
	Truncate       bool
	// Loading keeps the labels and shows a skeleton bar in place of each
	// value, for a record whose values have not arrived.
	Loading bool
}

// Descriptions lays read-only label and value pairs out in rows of Columns
// items. The item ending a row takes the columns left in it, so a row is
// never ragged. Below 320 DIPs a column a bordered table stacks each label
// over its value. It builds every item: page a record of unbounded length.
func Descriptions(c *ui.Context, o DescriptionsOptions) ui.Element {
	t := theme.Colors(c)
	want := o.Columns
	if want <= 0 {
		want = 3
	}
	minWidth := theme.Scale(c, descriptionsMinColumnWidth)
	if o.MinColumnWidth > 0 {
		minWidth = o.MinColumnWidth
	}
	stackBelow := theme.Scale(c, descriptionsStackWidth)
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space(3)).Role(ui.RoleGroup).Label(o.Title).Busy(o.Loading)
	// The width of the last layout; the first build has none and lays out
	// the columns asked for.
	width := root.Bounds().W
	columns, stacked := descriptionShape(want, width, minWidth, stackBelow)
	services := c.Services()
	root.Draw(func(_ *ui.Painter, r ui.Rect) {
		// A width that changes the shape builds the view again.
		if r.W <= 0 {
			return
		}
		if cols, st := descriptionShape(want, r.W, minWidth, stackBelow); cols != columns || st != stacked {
			services.Invalidate()
		}
	})
	layout := descriptionsLayout{
		columns: columns, loading: o.Loading, truncate: o.Truncate,
		stacked: stacked && o.Bordered,
		// A label never takes more than 160 DIPs or two fifths of its column,
		// so the value keeps room.
		labelMax: theme.Scale(c, 160),
	}
	if width > 0 {
		layout.labelMax = min(layout.labelMax, width/float32(columns)*0.4)
	}
	return root.Children(func() {
		if o.Title != "" {
			style.Text(ui.Text(c, o.Title), c, theme.H6()).TextColor(t.Foreground)
		}
		grid := func() {
			if o.Bordered {
				descriptionsTable(c, t, o.Items, layout)
			} else {
				descriptionsGrid(c, t, o.Items, layout)
			}
		}
		if o.Loading {
			SkeletonGroup(c, SkeletonGroupOptions{}, grid).Gap(0)
		} else {
			grid()
		}
	})
}

const (
	// descriptionsMinColumnWidth is the narrowest a column gets before the
	// grid drops a column.
	descriptionsMinColumnWidth = 200
	// descriptionsStackWidth is the column width under which a table puts a
	// label over its value instead of beside it.
	descriptionsStackWidth = 320
)

type descriptionsLayout struct {
	columns  int
	labelMax float32
	loading  bool
	truncate bool
	stacked  bool
}

// descriptionShape is the columns the width leaves and whether a table
// stacks its labels over the values.
func descriptionShape(want int, width, minWidth, stackBelow float32) (columns int, stacked bool) {
	if width <= 0 {
		return want, false
	}
	columns = min(want, max(1, int(width/minWidth)))
	return columns, width/float32(columns) < stackBelow
}

// descriptionSpans calls place with the columns each item takes: its own
// span, or all that is left of the row when the next item does not fit.
func descriptionSpans(items []DescriptionsItem, columns int, place func(i, span int)) {
	column := 0
	for i := range items {
		span, left := descriptionSpan(items[i], columns), columns-column
		if i == len(items)-1 || descriptionSpan(items[i+1], columns) > left-span {
			span = left
		}
		place(i, span)
		column = (column + span) % columns
	}
}

func descriptionSpan(it DescriptionsItem, columns int) int {
	return min(max(it.Span, 1), columns)
}

func descriptionsGrid(c *ui.Context, t theme.Palette, items []DescriptionsItem, l descriptionsLayout) {
	columns := l.columns
	ui.Grid(c).Columns(columns).FillWidth().MinWidth(0).
		GapX(theme.Space(6)).GapY(theme.Space(4)).Children(func() {
		descriptionSpans(items, columns, func(i, span int) {
			it := items[i]
			ui.Column(c.Key(i)).ColumnSpan(span).MinWidth(0).Gap(theme.Space(0.5)).
				AlignItems(ui.Start).Role(ui.RoleGroup).Label(it.Label).Children(func() {
				style.Text(ui.Text(c, it.Label), c, theme.XS()).TextColor(t.MutedForeground).MinWidth(0)
				descriptionValue(c, t, it, l, i)
			})
		})
	})
}

// descriptionsTable rules the cells with the grid's own background showing
// through one-hairline gaps, so shared edges are never drawn twice. A value
// cell is the group named by its label.
func descriptionsTable(c *ui.Context, t theme.Palette, items []DescriptionsItem, l descriptionsLayout) {
	columns := l.columns
	tracks := make([]ui.Track, 0, 2*columns)
	for range columns {
		tracks = append(tracks, ui.FitContent(), ui.Fr(1))
	}
	ui.Grid(c).ColumnTracks(tracks...).FillWidth().MinWidth(0).Gap(theme.Hairline).
		Background(t.Border).Border(theme.Hairline, t.Border).
		Radius(theme.RadiusScale().LG).Clip().Children(func() {
		descriptionSpans(items, columns, func(i, span int) {
			it := items[i]
			if l.stacked {
				// One cell per item: its label over its value, rule between.
				ui.Column(c.Key(i)).ColumnSpan(2 * span).MinWidth(0).Gap(theme.Hairline).
					Role(ui.RoleGroup).Label(it.Label).Children(func() {
					descriptionsLabelCell(c, t, it, l, 2*i, 1)
					descriptionsValueCell(c, t, it, l, 2*i+1, 1)
				})
				return
			}
			descriptionsLabelCell(c, t, it, l, 2*i, 1)
			descriptionsValueCell(c, t, it, l, 2*i+1, 2*span-1)
		})
	})
}

func descriptionsLabelCell(c *ui.Context, t theme.Palette, it DescriptionsItem, l descriptionsLayout, key, span int) {
	ui.Row(c.Key(key)).ColumnSpan(span).MinWidth(0).Padding(theme.Space(2), theme.Space(3)).
		AlignItems(ui.Start).Background(t.MutedBackground.Over(t.Background)).Children(func() {
		label := style.Text(ui.Text(c, it.Label), c, theme.SM()).TextColor(t.MutedForeground).MinWidth(0)
		if !l.stacked {
			label.MaxWidth(l.labelMax)
		}
	})
}

func descriptionsValueCell(c *ui.Context, t theme.Palette, it DescriptionsItem, l descriptionsLayout, key, span int) {
	cell := ui.Row(c.Key(key)).ColumnSpan(span).MinWidth(0).Padding(theme.Space(2), theme.Space(3)).
		AlignItems(ui.Start).Background(t.Background)
	if l.stacked {
		cell.Grow(1)
	}
	if !l.stacked {
		cell.Role(ui.RoleGroup).Label(it.Label)
	}
	cell.Children(func() { descriptionValue(c, t, it, l, key/2) })
}

// loadingBarWidths vary the placeholder widths by index, so a column of
// them does not read as a barcode.
var loadingBarWidths = [...]float32{96, 128, 72, 112}

// loadingBar is a skeleton bar inside a box of one line of style, so the
// placeholder takes the height of the text it stands for.
func loadingBar(c *ui.Context, line theme.TextStyle, width float32) {
	height := line.LineHeight * theme.TextScale(c)
	ui.Row(c).Height(height).AlignItems(ui.Center).MinWidth(0).Children(func() {
		Skeleton(c, SkeletonOptions{
			Width: theme.Scale(c, width), Height: height * 0.7, Animation: SkeletonNone,
		}).Shrink(1).MinWidth(0)
	})
}

func descriptionValue(c *ui.Context, t theme.Palette, it DescriptionsItem, l descriptionsLayout, i int) {
	if l.loading {
		loadingBar(c, theme.SM(), loadingBarWidths[i%len(loadingBarWidths)])
		return
	}
	if it.Content != nil {
		it.Content()
		return
	}
	if it.Value == "" {
		style.Text(ui.Text(c, "\u2014"), c, theme.SM()).TextColor(t.MutedForeground)
		return
	}
	value := style.Text(ui.Text(c, it.Value), c, theme.SM()).TextColor(t.Foreground).MinWidth(0).Selectable()
	if l.truncate {
		value.SingleLine().Tooltip(it.Value)
	}
}
