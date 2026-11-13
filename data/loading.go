package data

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const defaultLoadingRows = 5

// skeletonWidths varies the placeholder bars so a column does not read as a
// barcode.
var skeletonWidths = [...]float32{60, 85, 45, 70}

func loadingRows(n int) int {
	if n <= 0 {
		return defaultLoadingRows
	}
	return n
}

func tableCorner(o TableOptions) float32 {
	if o.Variant == TableSecondary {
		return theme.RadiusScale().XL
	}
	return 0
}

// wireSkeleton detaches the native table from the caller's selection and
// sort while it shows placeholders.
func (s *TableState) wireSkeleton() {
	s.List.Selected, s.List.Selection = nil, nil
	s.List.Sort = &s.nativeSort
	s.List.Key, s.List.CellMenu = nil, nil
}

// buildSkeletonCell builds a placeholder cell with the real cell's box, so
// rows and columns keep the table's metrics.
func (tl *tableLayout) buildSkeletonCell(c *ui.Context, row, col int) {
	cell := ui.Box(c).FillWidth().FillHeight().MinWidth(0).
		FixedLineHeight(theme.Scale(c, theme.Space20)).
		PaddingY(theme.Scale(c, theme.Space4)).Justify(ui.Center)
	switch cellColumn(tl, col).Align {
	case ui.Center:
		cell.AlignItems(ui.Center)
	case ui.End:
		cell.AlignItems(ui.End)
	}
	switch {
	case row > 0:
		cell.Draw(tl.paintRow)
	case tl.o.Variant == TableSecondary:
		cell.Draw(tl.paintFirst)
	}
	cell.Children(func() {
		bar := ui.Box(c).WidthPercent(skeletonWidths[(row+2*col)%len(skeletonWidths)]).MinWidth(0)
		bar.Children(func() {
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Height: theme.Scale(c, 12), Animation: mygoui.SkeletonNone})
		})
	})
}
