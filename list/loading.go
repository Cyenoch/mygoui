package list

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/busy"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const defaultLoadingRows = 5

// skeletonWidths varies the placeholder bars so rows do not read as a
// barcode.
var skeletonWidths = [...]float32{60, 85, 45, 70}

// treeSkeletonDepths indents placeholder rows at three depths, so they read
// as a tree.
var treeSkeletonDepths = [...]float32{0, 1, 2, 2, 1, 0}

func treeSkeletonIndent(i int) float32 {
	return treeSkeletonDepths[i%len(treeSkeletonDepths)] * theme.Space16
}

// buildSkeletonRows builds n placeholder rows of the given height inside the
// list being built, one shimmer sweep for all of them. indent gives a row's
// extra start padding in DIPs before scaling.
func buildSkeletonRows(c *ui.Context, n int, height, padX float32, indent func(i int) float32) {
	if n <= 0 {
		n = defaultLoadingRows
	}
	mygoui.SkeletonGroup(c, mygoui.SkeletonGroupOptions{}, func() {
		for i := range n {
			row := ui.Row(c).FillWidth().Shrink(0).AlignItems(ui.Center).
				Height(height).PaddingStart(padX + theme.Scale(c, indent(i))).PaddingEnd(padX)
			row.Children(func() {
				bar := ui.Box(c).WidthPercent(skeletonWidths[i%len(skeletonWidths)]).MinWidth(0)
				bar.Children(func() {
					mygoui.Skeleton(c, mygoui.SkeletonOptions{Height: theme.Scale(c, 12), Animation: mygoui.SkeletonNone})
				})
			})
		}
	}).Gap(0)
}

// refreshing takes a list's input while it loads new content and draws the
// progress line over its top edge. It returns the opacity its rows take, so
// the line stays at full strength above them.
func refreshing(c *ui.Context, list ui.Element, loading bool) float32 {
	shown := busy.Show(list, loading)
	list.IgnoreInput(loading)
	list.Children(func() { busy.Line(c, shown, 0) })
	return busy.Dimmed(shown)
}

// rowDim carries the rows' opacity from one build to the next: the list
// builds its rows before the list element that owns the fade exists.
type rowDim struct{ faded float32 }

func (d *rowDim) value() float32 { return 1 - d.faded }

// set asks for one more frame when the value moved, so the rows catch up
// with the frame that settled it.
func (d *rowDim) set(c *ui.Context, opacity float32) {
	if faded := 1 - opacity; faded != d.faded {
		d.faded = faded
		c.AnimationFrame()
	}
}
