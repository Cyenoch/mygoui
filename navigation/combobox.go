package navigation

import (
	"slices"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type ComboboxState struct {
	Query  string
	input  ui.Handle
	scroll ui.ScrollState
}
type ComboboxOptions[T comparable] struct {
	Label, Placeholder, Empty string
	Disabled, Multiple        bool
	Width                     float32
	OnChange                  func([]T)
}

// Combobox keeps native IME, editing, highlight, option scrolling and keyboard
// selection in ComboboxBase. Labels identify options and must be unique.
func Combobox[T comparable](c *ui.Context, s *ComboboxState, values *[]T, o ComboboxOptions[T], choices []Choice[T]) ui.Element {
	p := theme.Colors(c)
	notify := func() {
		if o.OnChange != nil {
			o.OnChange(*values)
		}
	}
	remove := func(value T) {
		if i := slices.Index(*values, value); i >= 0 {
			*values = slices.Delete(*values, i, i+1)
			notify()
		}
	}
	root := ui.Column(c).Gap(theme.Space4)
	if o.Width > 0 {
		root.Width(o.Width)
	}
	root.Children(func() {
		if o.Multiple && len(*values) > 0 {
			ui.Row(c.Key("chips")).Wrap().Gap(theme.Space4).Children(func() {
				for _, value := range *values {
					for _, choice := range choices {
						if choice.Value == value {
							ui.Row(c.Key(value)).Padding(theme.Space4, theme.Space6).Gap(theme.Space4).Radius(theme.Radii().XXL).Background(p.Input).Children(func() {
								ui.Text(c, choice.Label).FontSize(theme.Scale(c, theme.Small().Size))
								mygoui.Button(c, mygoui.ButtonOptions{Label: locales.ComponentRemoveNamed(choice.Label), Icon: lucide.X(), Size: mygoui.ButtonIconXS, Variant: mygoui.ButtonGhost, Disabled: o.Disabled || choice.Disabled,
									OnClick: func() { remove(value); s.input.Focus(c) }})
							})
							break
						}
					}
				}
			})
		}
		if !o.Multiple && !s.input.Focused(c) && len(*values) > 0 {
			for _, choice := range choices {
				if choice.Value == (*values)[0] {
					s.Query = choice.Label
					break
				}
			}
		}
		cb := ui.ComboboxBase(c.Key("input"), &s.Query)
		if o.Disabled {
			cb.SetOpen(false)
		}
		in := cb.Input.Bind(&s.input).Label(o.Label).Placeholder(o.Placeholder).Disabled(o.Disabled).
			Padding(theme.Space6, theme.Space12).Radius(theme.Radii().XXL).MinHeight(theme.HeightInput).Background(p.Input.Alpha(0.5)).TextColor(p.Foreground)
		in.OnTextKey(func(k ui.TextKey) bool {
			if k.Key == ui.KeyBackspace && k.Mods == 0 && k.Text() == "" && o.Multiple && len(*values) > 0 {
				value := (*values)[len(*values)-1]
				for _, choice := range choices {
					if choice.Value == value && choice.Disabled {
						return false
					}
				}
				remove(value)
				return true
			}
			return false
		})
		query := ""
		if cb.Filtering() {
			query = strings.ToLower(strings.TrimSpace(s.Query))
		}
		count := 0
		cb.Popup(func(panel ui.Element) {
			_, h := c.Size()
			panel.Column().Padding(theme.Space4).Margin(theme.Space4, 0, 0, 0).MaxHeight(max(0, h-theme.Space16)).Radius(theme.Radii().XXL).
				Background(p.Popover).TextColor(p.PopoverForeground).Border(theme.HairlineWidth, p.Border)
			viewport := ui.Scroll(c).TrackScroll(&s.scroll).MinHeight(0).MaxHeight(min(theme.Scale(c, 240), max(0, h-80)))
			mygoui.ScrollEdges(c, viewport, &s.scroll, mygoui.ScrollEdgesOptions{Surface: p.Popover})
			viewport.Children(func() {
				// Prefix matches precede substring matches without constructing a filtered
				// slice. Source order within each rank remains stable.
				group := ""
				for rank := range 2 {
					for _, choice := range choices {
						lower := strings.ToLower(choice.Label)
						prefix := strings.HasPrefix(lower, query)
						if !strings.Contains(lower, query) || (rank == 0) != prefix {
							continue
						}
						if choice.Group != group {
							if count > 0 {
								ui.Box(c).Height(theme.HairlineWidth).Background(p.Border).MarginY(theme.Space4)
							}
							group = choice.Group
							if group != "" {
								ui.Text(c, group).Padding(theme.Space4, theme.Space8).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(p.MutedForeground)
							}
						}
						var item ui.Element
						if choice.Disabled {
							item = ui.Row(c.Key(choice.Value)).Disabled(true).Label(choice.Label).TextColor(p.DisabledForeground)
						} else {
							c.Key(choice.Value)
							item = cb.Item(choice.Label).Label(choice.Label)
							if item.Highlighted() {
								item.Background(p.Accent).TextColor(p.AccentForeground)
							}
							count++
						}
						item.Padding(theme.Space6, theme.Space8).Gap(theme.Space8).Radius(theme.Radii().XL).Children(func() {
							ui.Text(c, choice.Label).Grow(1).SingleLine()
							if slices.Contains(*values, choice.Value) {
								ui.Icon(c, lucide.Check()).Size(theme.IconMedium, theme.IconMedium)
							}
						})
					}
				}
			})
		})
		// ComboboxBase intentionally hides an empty option collection. This native
		// passive overlay supplies empty feedback without a phantom selectable row.
		if cb.Open() && count == 0 && !o.Disabled {
			ui.Overlay(c, func() {
				empty := o.Empty
				if empty == "" {
					empty = locales.ComponentNoMatches()
				}
				ui.Text(c, empty).AttachTo(in, ui.AnchorBottomStart, ui.AnchorTopStart).Margin(theme.Space4, 0, 0, 0).
					Padding(theme.Space12).Radius(theme.Radii().XXL).Background(p.Popover).TextColor(p.MutedForeground).Border(theme.HairlineWidth, p.Border)
			})
		}
		if chosen, ok := cb.Chosen(); ok && !o.Disabled {
			for _, choice := range choices {
				if choice.Label == chosen && !choice.Disabled {
					if o.Multiple {
						if slices.Contains(*values, choice.Value) {
							remove(choice.Value)
						} else {
							*values = append(*values, choice.Value)
							notify()
						}
						s.Query = ""
					} else {
						*values = append((*values)[:0], choice.Value)
						s.Query = choice.Label
						notify()
					}
					break
				}
			}
		}
	})
	return root
}
