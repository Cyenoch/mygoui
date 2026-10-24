package mygoui

import (
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// SpinnerSize is HeroUI's indicator scale.
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
// scheduling and reduced-motion poses. HeroUI's dual-gradient arc is
// approximated by the native arc, whose stroke and cadence stay platform
// native.
func Spinner(c *ui.Context, o SpinnerOptions) ui.Element {
	size := float32(24)
	switch o.Size {
	case SpinnerSmall:
		size = 16
	case SpinnerLarge:
		size = 32
	case SpinnerExtraLarge:
		size = 40
	}
	label := o.Label
	if label == "" {
		label = locale.M().Loading
	}
	return loading.Arc(c, loading.Options{
		Size: theme.Scale(c, size), Color: statusForeground(theme.Colors(c), o.Color),
	}).Label(label)
}

func statusForeground(t theme.Palette, color StatusColor) ui.Color {
	switch color {
	case ColorAccent:
		return t.Accent
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
