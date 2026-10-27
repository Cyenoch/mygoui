package mygoui

import (
	"time"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// SpinnerSize selects a native loading indicator size.
type SpinnerSize uint8

const (
	SpinnerMedium SpinnerSize = iota
	SpinnerSmall
	SpinnerLarge
	SpinnerExtraLarge
)

// SpinnerOptions configures a loading indicator.
type SpinnerOptions struct {
	Size  SpinnerSize
	Color StatusColor
	Label string
}

// Spinner uses MyGo's maintained arc indicator. The toolkit owns frame
// scheduling and reduced-motion poses. The native arc approximates COSS's
// rotating loading icon.
func Spinner(c *ui.Context, o SpinnerOptions) ui.Element {
	size := float32(16)
	switch o.Size {
	case SpinnerSmall:
		size = 12
	case SpinnerLarge:
		size = 24
	case SpinnerExtraLarge:
		size = 32
	}
	label := o.Label
	if label == "" {
		label = locale.M().Loading
	}
	return loading.Arc(c, loading.Options{
		Size: theme.Scale(c, size), Color: statusForeground(theme.Colors(c), o.Color),
		// COSS's standard spin completes one turn each second.
		Duration: time.Second,
	}).Label(label)
}

func statusForeground(t theme.Palette, color StatusColor) ui.Color {
	switch color {
	case ColorAccent:
		return t.Primary
	case ColorSuccess:
		return t.Success
	case ColorWarning:
		return t.Warning
	case ColorDanger:
		return t.Danger
	case ColorInfo:
		return t.Info
	}
	return t.Foreground
}
