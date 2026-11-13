package gallery

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// languageToken is one non-color theme token: its role and the value the
// published theme holds now.
type languageToken struct {
	name  string
	role  string
	value func(theme.Theme) string
}

// languageTokens lists the motion, spacing and effect tokens in display order.
var languageTokens = []languageToken{
	{"Spacing.Unit", "One spacing step; theme.Space(n) is n steps.", func(t theme.Theme) string { return dip(t.Spacing.Unit) }},
	{"Motion.State", "A control changing state in place: fills, thumbs, marks.", func(t theme.Theme) string { return ms(t.Motion.State) }},
	{"Motion.Reveal", "Content opening or a selection sliding to a new place.", func(t theme.Theme) string { return ms(t.Motion.Reveal) }},
	{"Motion.Roll", "A value rolling to the next one.", func(t theme.Theme) string { return ms(t.Motion.Roll) }},
	{"Motion.Shift", "Neighbours making room while an item is dragged.", func(t theme.Theme) string { return ms(t.Motion.Shift) }},
	{"Motion.Page", "One page replacing another.", func(t theme.Theme) string { return ms(t.Motion.Page) }},
	{"Motion.Popup", "Anchored surfaces: menus, popovers, tooltips.", func(t theme.Theme) string { return ms(t.Motion.Popup) }},
	{"Motion.Dialog", "Centered and edge dialogs, with their backdrop.", func(t theme.Theme) string { return ms(t.Motion.Dialog) }},
	{"Motion.Drawer", "Drawers and sheets travelling from an edge.", func(t theme.Theme) string { return ms(t.Motion.Drawer) }},
	{"Motion.Ease", "The default curve; ease-out in the built-in theme.", easing},
	{"Motion.EaseTravel", "Curve for marks that travel; ease-in-out in the built-in theme.", easing},
	{"Motion.EaseDrawer", "Curve for drawers and sheets; a long settle in the built-in theme.", easing},
	{"Effects.HoverOpacity", "How much of a fill is left under the pointer.", func(t theme.Theme) string { return percent(t.Effects.HoverOpacity) }},
	{"Effects.PopupScale", "The scale a surface arrives from.", func(t theme.Theme) string { return percent(t.Effects.PopupScale) }},
	{"Effects.BackdropBlur", "Blur behind a modal surface.", func(t theme.Theme) string { return dip(t.Effects.BackdropBlur) }},
	{"Effects.FadeEdge", "Depth of a scroll edge fade.", func(t theme.Theme) string { return dip(t.Effects.FadeEdge) }},
	{"Effects.SpinnerPeriod", "One turn of a spinner.", func(t theme.Theme) string { return ms(t.Effects.SpinnerPeriod) }},
	{"Effects.SweepPeriod", "One pass of an indeterminate bar.", func(t theme.Theme) string { return ms(t.Effects.SweepPeriod) }},
	{"Effects.SweepBand", "Width of the indeterminate band, as a share of the track.", func(t theme.Theme) string { return percent(t.Effects.SweepBand) }},
	{"Effects.LoadingDelay", "How long loading stays hidden before it shows.", func(t theme.Theme) string { return ms(t.Effects.LoadingDelay) }},
	{"Effects.RefreshOpacity", "Opacity of content while it refreshes.", func(t theme.Theme) string { return percent(t.Effects.RefreshOpacity) }},
	{"Effects.ShimmerPeriod", "One loop of a loading placeholder's shimmer or pulse.", func(t theme.Theme) string { return ms(t.Effects.ShimmerPeriod) }},
	{"Effects.ShimmerPeak", "How far the shimmer band lightens a placeholder in light appearance.", func(t theme.Theme) string { return percent(t.Effects.ShimmerPeak) }},
	{"Effects.ShimmerPeakDark", "How far the shimmer band lightens a placeholder in dark appearance.", func(t theme.Theme) string { return percent(t.Effects.ShimmerPeakDark) }},
}

// pageLanguageTokens names the language tokens a component page lists after
// its colors.
var pageLanguageTokens = map[string][]string{
	"checkbox":         {"Motion.State"},
	"switch":           {"Motion.State"},
	"slider":           {"Motion.State", "Motion.Reveal"},
	"togglebutton":     {"Motion.State"},
	"tabs":             {"Motion.Reveal", "Motion.EaseTravel"},
	"segmentedcontrol": {"Motion.Reveal", "Motion.EaseTravel"},
	"numberflip":       {"Motion.Roll"},
	"sortablelist":     {"Motion.Shift"},
	"documenttabs":     {"Motion.Shift", "Effects.FadeEdge"},
	"carousel":         {"Motion.Page", "Motion.State", "Motion.Reveal"},
	"navigationstack":  {"Motion.Page"},
	"popover":          {"Motion.Popup", "Effects.PopupScale"},
	"tooltip":          {"Motion.Popup", "Effects.PopupScale"},
	"menu":             {"Motion.Popup", "Effects.PopupScale"},
	"dialog":           {"Motion.Dialog", "Motion.EaseTravel", "Effects.PopupScale", "Effects.BackdropBlur"},
	"alertdialog":      {"Motion.Dialog", "Motion.EaseTravel", "Effects.PopupScale", "Effects.BackdropBlur"},
	"drawer":           {"Motion.Drawer", "Motion.EaseDrawer", "Effects.BackdropBlur"},
	"sheet":            {"Motion.Drawer", "Motion.EaseDrawer", "Effects.BackdropBlur"},
	"spinner":          {"Effects.SpinnerPeriod"},
	"progressbar":      {"Effects.SweepPeriod", "Effects.SweepBand"},
	"scrollfade":       {"Effects.FadeEdge"},
	"mask":             {"Effects.FadeEdge"},
	"skeleton":         {"Effects.ShimmerPeriod", "Effects.ShimmerPeak", "Effects.ShimmerPeakDark"},
	"shimmer":          {"Effects.ShimmerPeriod", "Effects.ShimmerPeak", "Effects.ShimmerPeakDark"},
	"reveal":           {"Motion.Reveal", "Motion.Ease", "Effects.FadeEdge"},
	"table":            {"Effects.LoadingDelay", "Effects.RefreshOpacity"},
	"loadingoverlay":   {"Effects.LoadingDelay", "Motion.State"},
}

func init() {
	for _, token := range languageTokens {
		tokenRoles[token.name] = token.role
	}
}

// languageValue reports the published value of a language token.
func languageValue(name string) (string, bool) {
	for _, token := range languageTokens {
		if token.name == name {
			return token.value(theme.Current()), true
		}
	}
	return "", false
}

// languageTokenCards is the theming page's specimen of every language token.
func languageTokenCards(c *ui.Context, s *State) {
	names := make([]string, len(languageTokens))
	for i, token := range languageTokens {
		names[i] = token.name
	}
	tokenCards(c, s, names)
}

func ms(d time.Duration) string { return strconv.FormatInt(d.Milliseconds(), 10) + " ms" }

func dip(v float32) string { return strconv.FormatFloat(float64(v), 'g', -1, 32) + " DIP" }

func percent(v float32) string { return fmt.Sprintf("%.0f%%", v*100) }

func easing(theme.Theme) string { return "ui.Easing" }
