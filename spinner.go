package mygoui

import (
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
)

// SpinnerOptions sets the indicator's unscaled size and accessible name.
type SpinnerOptions struct {
	Size  float32
	Label string
}

// Spinner uses MyGo's native arc indicator. The toolkit owns paint-only frame
// scheduling, visibility, reduced-motion presentation and progress semantics.
func Spinner(c *ui.Context, o SpinnerOptions) ui.Element {
	size := o.Size
	if size <= 0 {
		size = theme.IconMedium
	}
	label := o.Label
	if label == "" {
		label = locale.ComponentLoading()
	}
	return loading.Arc(c, loading.Options{
		Size: theme.Scale(c, size), Color: theme.Colors(c).MutedForeground,
	}).Label(label)
}
