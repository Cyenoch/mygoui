package navigation

import (
	"strings"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type Choice[T comparable] struct {
	Value        T
	Label, Group string
	Disabled     bool
}
type SelectState struct {
	prefix    string
	typed     time.Time
	highlight bool
	scroll    ui.ScrollState
}
type SelectOptions[T comparable] struct {
	Label, Placeholder string
	Disabled, Small    bool
	Width              float32
	OnChange           func(T)
}

// Select is a custom native popup; disabled rows never enter the toolkit's
// selectable-value sequence, so keyboard navigation cannot commit them.
func Select[T comparable](c *ui.Context, s *SelectState, value *T, o SelectOptions[T], choices []Choice[T]) ui.Element {
	p := theme.Colors(c)
	sel := ui.SelectBase(c, value)
	label := o.Placeholder
	for _, choice := range choices {
		if choice.Value == *value {
			label = choice.Label
			break
		}
	}
	tr := sel.Trigger.Label(o.Label).Value(label).Disabled(o.Disabled).Justify(ui.Start).
		Padding(theme.Space6, theme.Space12).Gap(theme.Space6).Radius(theme.Radii().XXL).
		MinHeight(theme.HeightInput).FontSize(theme.Scale(c, theme.Body().Size)).Background(p.Input.Alpha(0.5)).TextColor(p.Foreground)
	if o.Small {
		tr.MinHeight(theme.HeightButton)
	}
	if o.Width > 0 {
		tr.Width(o.Width)
	}
	tr.Children(func() {
		ui.Text(c, label).SingleLine().Grow(1)
		ui.Icon(c, lucide.ChevronDown()).Size(theme.IconMedium, theme.IconMedium).TextColor(p.MutedForeground)
	})
	tr.HandleInput(func(ev ui.InputEvent) bool {
		if o.Disabled {
			return false
		}
		text := ""
		if ev.Kind == ui.InputText {
			text = ev.Text
		} else if ev.Kind == ui.InputKeyDown && ev.Mods == 0 {
			if ev.Key >= ui.KeyA && ev.Key <= ui.KeyZ {
				text = string(rune('a' + ev.Key - ui.KeyA))
			}
			if ev.Key >= ui.Key0 && ev.Key <= ui.Key9 {
				text = string(rune('0' + ev.Key - ui.Key0))
			}
		}
		if text == "" {
			return false
		}
		if ev.Time.Sub(s.typed) > time.Second {
			s.prefix = ""
		}
		s.typed = ev.Time
		s.prefix += strings.ToLower(text)
		s.highlight = true
		return true
	})
	sel.Popup(func(panel ui.Element) {
		_, h := c.Size()
		panel.Column().Disabled(o.Disabled).Padding(theme.Space4).Gap(theme.Space4).Margin(theme.Space4, 0, 0, 0).
			MaxHeight(max(0, h-theme.Space16)).Radius(theme.Radii().XXL).Background(p.Popover).TextColor(p.PopoverForeground).Border(theme.HairlineWidth, p.Border)
		// Both buttons remain in place as the viewport moves, avoiding layout jumps.
		mygoui.Button(c.Key("up"), mygoui.ButtonOptions{Label: locales.ComponentScrollChoicesUp(), Icon: lucide.ChevronUp(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconXS,
			Disabled: s.scroll.Y <= 0, OnClick: func() { s.scroll.Y = max(0, s.scroll.Y-theme.HeightRowCompact*3) }}).KeepFocus().AlignSelf(ui.Center)
		viewport := ui.Scroll(c.Key("choices")).TrackScroll(&s.scroll).MinHeight(0).MaxHeight(min(theme.Scale(c, 240), max(0, h-112)))
		mygoui.ScrollEdges(c, viewport, &s.scroll, mygoui.ScrollEdgesOptions{Surface: p.Popover})
		viewport.Children(func() {
			group := ""
			for i, choice := range choices {
				if choice.Group != group {
					if i > 0 {
						ui.Box(c).Height(theme.HairlineWidth).Background(p.Border).MarginY(theme.Space4)
					}
					group = choice.Group
					if group != "" {
						ui.Text(c, group).Padding(theme.Space4, theme.Space8).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(p.MutedForeground)
					}
				}
				var item ui.Element
				if choice.Disabled {
					item = ui.Row(c.Key(choice.Value)).Disabled(true).TextColor(p.DisabledForeground).Label(choice.Label)
				} else {
					c.Key(choice.Value)
					item = sel.Item(choice.Value).Role(ui.RoleListItem).Label(choice.Label).Selected(choice.Value == *value)
					if item.Highlighted() {
						item.Background(p.Accent).TextColor(p.AccentForeground)
					}
				}
				item.Padding(theme.Space6, theme.Space8).Gap(theme.Space8).MinHeight(theme.Scale(c, 28)).Radius(theme.Radii().XL).Children(func() {
					ui.Text(c, choice.Label).Grow(1).SingleLine()
					if choice.Value == *value {
						ui.Icon(c, lucide.Check()).Size(theme.IconMedium, theme.IconMedium)
					}
				})
			}
		})
		mygoui.Button(c.Key("down"), mygoui.ButtonOptions{Label: locales.ComponentScrollChoicesDown(), Icon: lucide.ChevronDown(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconXS,
			Disabled: s.scroll.Y >= s.scroll.MaxY, OnClick: func() { s.scroll.Y = min(s.scroll.MaxY, s.scroll.Y+theme.HeightRowCompact*3) }}).KeepFocus().AlignSelf(ui.Center)
	})
	if s.highlight {
		for _, choice := range choices {
			if !choice.Disabled && strings.HasPrefix(strings.ToLower(choice.Label), s.prefix) {
				sel.Highlight(choice.Value)
				break
			}
		}
		s.highlight = false
	}
	tr.OnChange(func() {
		if o.OnChange != nil {
			o.OnChange(*value)
		}
	})
	return tr
}
