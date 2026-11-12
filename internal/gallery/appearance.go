package gallery

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// appearanceRadii are the radius bases the gallery offers.
var appearanceRadii = []float32{4, 10, 16}

// appearanceButton opens the theme controls in a popover, so changing them
// never moves the page.
func appearanceButton(c *ui.Context, s *State) {
	icon := lucide.Sun()
	if s.Dark {
		icon = lucide.Moon()
	}
	anchor := mygoui.Button(c, mygoui.ButtonOptions{
		Label: "Appearance", Icon: icon, IconOnly: true,
		Variant: buttonFor(s.AppearanceOpen), Size: theme.ButtonSmall,
		OnClick: func() { s.AppearanceOpen = !s.AppearanceOpen },
	}).Role(ui.RoleToggleButton).Checked(s.AppearanceOpen).Shrink(0).AlignSelf(ui.Center)
	overlay.Popover(c, anchor, &s.AppearanceOpen, overlay.PopoverOptions{
		Label:     "Appearance settings",
		Placement: mygoui.Placement{Side: mygoui.SideBottom, Align: mygoui.AlignEnd},
		Width:     theme.Scale(c, 264),
		Flyout:    &s.Flyouts,
		Body:      func(c *ui.Context) { appearanceControls(c, s) },
	})
}

func appearanceControls(c *ui.Context, s *State) {
	ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
		mode := 0
		if s.Dark {
			mode = 1
		}
		appearanceRow(c, "Mode", func() {
			navigation.SegmentedControl(c.Key("mode"), navigation.SegmentedControlOptions{
				Selected: &mode, Label: "Mode",
				Items:    []navigation.TabOptions{{Label: "Light"}, {Label: "Dark"}},
				OnChange: func(i int) { s.setDark(i == 1) },
			})
		})
		if !s.SystemFontOnly {
			font := 0
			if s.SystemFont {
				font = 1
			}
			appearanceRow(c, "Font", func() {
				navigation.SegmentedControl(c.Key("font"), navigation.SegmentedControlOptions{
					Selected: &font, Label: "Font",
					Items:    []navigation.TabOptions{{Label: "Cal Sans"}, {Label: "System"}},
					OnChange: func(i int) { s.SystemFont = i == 1 },
				})
			})
		}
		radius := 1
		items := make([]navigation.TabOptions, len(appearanceRadii))
		for i, r := range appearanceRadii {
			items[i] = navigation.TabOptions{Label: radiusLabel(r)}
			if s.Radius == r {
				radius = i
			}
		}
		appearanceRow(c, "Radius", func() {
			navigation.SegmentedControl(c.Key("radius"), navigation.SegmentedControlOptions{
				Selected: &radius, Label: "Radius", Items: items,
				OnChange: func(i int) { s.Radius = appearanceRadii[i] },
			})
		})
		appearanceRow(c, "Primary", func() { swatches(c, s) })
	})
}

// appearanceRow is one labelled control line in the appearance popover.
func appearanceRow(c *ui.Context, label string, control func()) {
	ui.Column(c.Key(label)).FillWidth().Gap(theme.Space8).Children(func() {
		mygoui.Typography(c, label, mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
		control()
	})
}

func radiusLabel(r float32) string {
	switch r {
	case 4:
		return "Small"
	case 16:
		return "Large"
	}
	return "Medium"
}
