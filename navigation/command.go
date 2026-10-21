package navigation

import (
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type CommandID string
type CommandItem struct {
	ID                                       CommandID
	Label, Detail, Group, Keywords, Shortcut string
	Badge, DisabledReason                    string
	Icon                                     *ui.SVG
	Media                                    func()
	Disabled                                 bool
	Mods                                     ui.Modifiers
	Key                                      ui.Key
	OnSelect                                 func()
}
type CommandState struct {
	Query     string
	active    CommandID
	first     ui.Handle
	search    ui.Handle
	lastQuery string
	scroll    ui.ScrollState
}
type CommandOptions struct {
	Label, Placeholder, Empty string
	Height                    float32
	AutoFocus                 bool
	OnBack                    func()
}

// commandRank ranks exact labels, prefixes, substring labels and keyword
// matches. All query terms must match; an empty query retains source order.
func commandRank(item CommandItem, query string) int {
	label := strings.ToLower(item.Label)
	if query == "" || label == query {
		return 0
	}
	if strings.HasPrefix(label, query) {
		return 1
	}
	if strings.Contains(label, query) {
		return 2
	}
	text := label + " " + strings.ToLower(item.Detail+" "+item.Keywords)
	for _, term := range strings.Fields(query) {
		if !strings.Contains(text, term) {
			return -1
		}
	}
	return 3
}

// Command provides inline search and actions with native editing and a native
// roving focus group. Shortcut labels describe actual scoped shortcuts when Key
// and Mods are supplied. Disabled commands remain readable and cannot execute.
func Command(c *ui.Context, s *CommandState, o CommandOptions, items []CommandItem) ui.Element {
	p := theme.Colors(c)
	query := strings.ToLower(strings.TrimSpace(s.Query))
	queryChanged := s.lastQuery != query
	s.lastQuery = query
	outerRadius, inset := theme.Radii().XXL, float32(theme.Space4)
	innerRadius := max(float32(0), outerRadius-inset)
	root := ui.Column(c).Gap(theme.Space4).Padding(inset).Radius(outerRadius).Clip().Background(p.Popover).TextColor(p.PopoverForeground)
	root.Children(func() {
		var first *CommandItem
		activePresent := false
		for rank := range 4 {
			for i := range items {
				item := &items[i]
				if commandRank(*item, query) != rank {
					continue
				}
				if !item.Disabled && item.OnSelect != nil {
					if first == nil {
						first = item
					}
					if item.ID == s.active {
						activePresent = true
					}
				}
			}
		}
		if (!activePresent || queryChanged) && first != nil {
			s.active = first.ID
		}
		ui.Row(c).Label(locales.ComponentCommandSearchSurface()).Padding(theme.Space6, theme.Space12).Gap(theme.Space8).Radius(innerRadius).Background(p.Input.Alpha(0.5)).Children(func() {
			if o.OnBack != nil {
				mygoui.Button(c.Key("back"), mygoui.ButtonOptions{Label: locales.ComponentBack(), Icon: lucide.ArrowLeft(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconSM,
					OnClick: func() { o.OnBack(); s.search.Focus(c) }})
			} else {
				ui.Icon(c, lucide.Search()).Size(theme.IconMedium, theme.IconMedium).TextColor(p.MutedForeground)
			}
			in := ui.TextInputBase(c.Key("search"), &s.Query).Bind(&s.search).Label(o.Label).Placeholder(o.Placeholder).Grow(1).MinWidth(0).MinHeight(theme.HeightInput)
			if o.AutoFocus {
				in.AutoFocus()
			}
			composing := in.Composing()
			in.OnTextKey(func(k ui.TextKey) bool {
				if composing || k.Mods != 0 {
					return false
				}
				switch k.Key {
				case ui.KeyDown:
					if first != nil {
						s.first.Focus(c)
						return true
					}
				case ui.KeyEnter:
					for _, item := range items {
						if item.ID == s.active && !item.Disabled && item.OnSelect != nil && commandRank(item, strings.ToLower(strings.TrimSpace(k.Text()))) >= 0 {
							item.OnSelect()
							return true
						}
					}
				case ui.KeyBackspace:
					if k.Text() == "" && o.OnBack != nil {
						o.OnBack()
						return true
					}
				}
				return false
			})
		})
		height := o.Height
		if height <= 0 {
			height = theme.Scale(c, 240)
		}
		viewport := ui.Scroll(c.Key("commands")).TrackScroll(&s.scroll).Height(height).MinHeight(0).FocusGroup(ui.Vertical)
		mygoui.ScrollEdges(c, viewport, &s.scroll, mygoui.ScrollEdgesOptions{Surface: p.Popover})
		viewport.Children(func() {
			group := ""
			count := 0
			bound := false
			for rank := range 4 {
				for _, item := range items {
					if commandRank(item, query) != rank {
						continue
					}
					if group != item.Group {
						if count > 0 {
							ui.Box(c).Height(theme.HairlineWidth).MarginY(theme.Space4).Background(p.Border)
						}
						group = item.Group
						if group != "" {
							ui.Text(c, group).Padding(theme.Space6, theme.Space8).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(p.MutedForeground)
						}
					}
					disabled := item.Disabled || item.OnSelect == nil
					b := ui.ButtonBase(c.Key(item.ID)).Label(item.Label).Disabled(disabled).Role(ui.RoleMenuItem).Justify(ui.Start).
						Padding(theme.Space6, theme.Space8).MinHeight(theme.Scale(c, 28)).Gap(theme.Space8).Radius(innerRadius)
					b.MinWidth(0).FillWidth().Description(item.Detail)
					if item.DisabledReason != "" {
						b.Description(item.DisabledReason).Tooltip(item.DisabledReason)
					}
					if !disabled && !bound {
						b.Bind(&s.first)
						bound = true
					}
					if b.Focused() {
						s.active = item.ID
						b.ScrollIntoView()
					}
					if b.Focused() || b.Hovered() {
						b.Background(p.Muted)
					}
					if disabled {
						b.TextColor(p.DisabledForeground)
					}
					b.OnClick(item.OnSelect).Children(func() {
						if item.Media != nil {
							item.Media()
						} else if item.Icon != nil {
							ui.Icon(c, item.Icon).Size(theme.Scale(c, theme.IconMedium), theme.Scale(c, theme.IconMedium)).Shrink(0)
						}
						ui.Column(c).Grow(1).MinWidth(0).Gap(theme.Space4).Children(func() {
							ui.Text(c, item.Label).SingleLine()
							if item.Detail != "" {
								ui.Text(c, item.Detail).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(p.MutedForeground)
							}
						})
						if item.Badge != "" {
							mygoui.Badge(c, mygoui.BadgeOptions{Label: item.Badge, Variant: mygoui.BadgeOutline})
						}
						if item.Shortcut != "" {
							mygoui.Kbd(c, item.Shortcut)
						}
					})
					if item.Key != ui.KeyUnknown && !disabled {
						root.OnShortcut(item.Mods, item.Key, item.OnSelect)
					}
					count++
				}
			}
			if count == 0 {
				empty := o.Empty
				if empty == "" {
					empty = locales.ComponentNoCommandsFound()
				}
				ui.Text(c, empty).Padding(theme.Space24).TextColor(p.MutedForeground)
			}
		})
	})
	return root
}

// CommandDialog composes the same command view with the shared modal lifecycle.
func CommandDialog(c *ui.Context, open *bool, s *CommandState, dialog overlay.DialogOptions, o CommandOptions, items []CommandItem) ui.Element {
	return overlay.Dialog(c, open, dialog, func() { Command(c, s, o, items) })
}
