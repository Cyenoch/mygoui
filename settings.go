package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/busy"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SettingsGroupOptions configures a titled group of settings rows. Action
// builds a control at the end of the header, as a button acting on the
// whole group.
type SettingsGroupOptions struct {
	Title, Description string
	Action             func()
	Rows               []SettingsRow
	// Loading shows LoadingRows skeleton rows while the group has no rows
	// yet. With rows it is a refresh: they dim, ignore input and a line runs
	// along the top of the card.
	Loading     bool
	LoadingRows int
}

// SettingsRow is one row of a SettingsGroup: a leading icon, a title over a
// subtitle, and Trailing, which builds the row's control or value.
//
// Control is the handle the caller bound to the control Trailing builds, so
// that the title names it for assistive technology; on a plain row a click
// on the title also focuses or toggles it.
//
// OnClick makes the whole row a button with a chevron; its Trailing then
// holds passive content, as the current value. Rows makes the row an
// expander revealing them, open while *Expanded, which the row keeps per
// window when nil, under the row's Title, which then must differ among the
// expanders of a group; an expander's own OnClick is not used.
//
// The trailing control sits beside the text and goes below it where they
// do not fit in one line. Destructive inks the icon and the title with the
// danger color, for a row that deletes or resets. The rows of an expander
// start at the edge of its title.
type SettingsRow struct {
	Icon            *ui.SVG
	Title, Subtitle string
	Trailing        func()
	Control         *ui.Handle
	OnClick         func()
	Rows            []SettingsRow
	Expanded        *bool
	Disabled        bool
	Destructive     bool
}

type settingsExpandedKey struct{}

// SettingsGroup is a boxed list of settings: a header over one rounded card
// whose rows are separated by hairlines.
func SettingsGroup(c *ui.Context, o SettingsGroupOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space8)).Role(ui.RoleGroup)
	if o.Title != "" {
		root.Label(o.Title)
	}
	if o.Description != "" {
		root.Description(o.Description)
	}
	root.Busy(o.Loading)
	shown := busy.Show(root, o.Loading)
	return root.Children(func() {
		if o.Title != "" || o.Description != "" || o.Action != nil {
			settingsHeader(c, t, &o)
		}
		ui.Column(c).FillWidth().MinWidth(0).Children(func() {
			card := ui.Column(c).FillWidth().MinWidth(0).Radius(theme.RadiusScale().XL).Clip().
				Background(t.Card).TextColor(t.CardForeground).
				Border(theme.Hairline, t.Border)
			refresh := o.Loading && len(o.Rows) > 0
			busy.Dim(card, refresh, shown)
			card.Children(func() {
				if o.Loading && len(o.Rows) == 0 {
					settingsSkeletonRows(c, &o)
					return
				}
				for i := range o.Rows {
					if i > 0 {
						Separator(c, SeparatorOptions{})
					}
					settingsRow(c, t, &o.Rows[i], 0)
				}
			})
			if refresh {
				busy.Line(c, shown, theme.RadiusScale().XL)
			}
		})
	})
}

// settingsSkeletonRows stands in for the rows of a group still waiting for
// them: a title bar over a shorter subtitle bar and a control placeholder.
func settingsSkeletonRows(c *ui.Context, o *SettingsGroupOptions) {
	rows := o.LoadingRows
	if rows <= 0 {
		rows = 3
	}
	SkeletonGroup(c, SkeletonGroupOptions{}, func() {
		for i := range rows {
			if i > 0 {
				Separator(c, SeparatorOptions{})
			}
			bar := func(width, height float32) {
				Skeleton(c, SkeletonOptions{
					Width: theme.Scale(c, width), Height: theme.Scale(c, height), Animation: SkeletonNone,
				}).Shrink(1).MinWidth(0)
			}
			settingsRowFrame(c, ui.Row(c.Key(i)).FillWidth(), 0).Children(func() {
				ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Scale(c, theme.Space8)).
					Justify(ui.Center).AlignItems(ui.Start).Children(func() {
					bar(loadingBarWidths[i%len(loadingBarWidths)]+32, 14)
					bar(loadingBarWidths[(i+2)%len(loadingBarWidths)]+48, 12)
				})
				bar(36, 20)
			})
		}
	}).Gap(0)
}

func settingsHeader(c *ui.Context, t theme.Palette, o *SettingsGroupOptions) {
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.End).Gap(theme.Scale(c, theme.Space12)).
		Padding(0, theme.Scale(c, theme.Space4)).Children(func() {
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Scale(c, theme.Space2)).Children(func() {
			if o.Title != "" {
				style.Text(ui.Text(c, o.Title), c, theme.SM()).TextColor(t.Foreground).
					FontWeight(theme.WeightMedium)
			}
			if o.Description != "" {
				style.Text(ui.Text(c, o.Description), c, theme.SM()).TextColor(t.Muted)
			}
		})
		if o.Action != nil {
			ui.Row(c).Shrink(0).AlignItems(ui.Center).Children(o.Action)
		}
	})
}

// settingsRow builds o with its content start inset DIPs further in than a
// top-level row's.
func settingsRow(c *ui.Context, t theme.Palette, o *SettingsRow, inset float32) {
	switch {
	case len(o.Rows) > 0:
		settingsExpanderRow(c, t, o, inset)
	case o.OnClick != nil:
		settingsPressableRow(c, t, o, inset)
	default:
		settingsPlainRow(c, t, o, inset)
	}
}

// settingsRowFrame applies the row's box to e, its content inset further
// from the start edge.
func settingsRowFrame(c *ui.Context, e ui.Element, inset float32) ui.Element {
	e = e.MinWidth(0).MinHeight(theme.Scale(c, 52)).AlignItems(ui.Center).
		Gap(theme.Scale(c, theme.Space12)).
		Padding(theme.Scale(c, theme.Space8), theme.Scale(c, theme.Space16))
	if inset > 0 {
		e.PaddingStart(theme.Scale(c, theme.Space16) + inset)
	}
	return e
}

// settingsTextInset is how far a row's title starts from the start of its
// content: the icon and the gap after it.
func settingsTextInset(c *ui.Context, o *SettingsRow) float32 {
	if o.Icon == nil {
		return 0
	}
	return theme.Scale(c, theme.IconLG) + theme.Scale(c, theme.Space12)
}

// settingsRowBody builds the leading icon, then the texts and, if trailing,
// the control after them. The texts ask for a readable measure and the
// control goes below them where the two do not fit side by side. It returns
// the title.
func settingsRowBody(c *ui.Context, t theme.Palette, o *SettingsRow, trailing bool) ui.Element {
	ink := t.CardForeground
	iconInk := t.Muted
	if o.Destructive {
		ink, iconInk = t.Destructive, t.Destructive
	}
	if o.Icon != nil {
		size := theme.Scale(c, theme.IconLG)
		ui.Icon(c, o.Icon).Size(size, size).Shrink(0).TextColor(iconInk)
	}
	var title ui.Element
	ui.Row(c).Grow(1).Basis(0).MinWidth(0).Wrap().AlignItems(ui.Center).
		GapX(theme.Scale(c, theme.Space12)).GapY(theme.Scale(c, theme.Space8)).Children(func() {
		ui.Column(c).Grow(1).Basis(theme.Scale(c, 120)).MinWidth(0).Children(func() {
			title = style.Text(ui.Text(c, o.Title), c, theme.SM()).FontWeight(theme.WeightMedium).TextColor(ink)
			if o.Subtitle != "" {
				style.Text(ui.Text(c, o.Subtitle), c, theme.SM()).TextColor(t.Muted)
			}
		})
		if trailing {
			settingsTrailing(c, t, o)
		}
	})
	return title
}

func settingsTrailing(c *ui.Context, t theme.Palette, o *SettingsRow) {
	if o.Trailing == nil {
		return
	}
	ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space8)).
		TextColor(t.Muted).Children(o.Trailing)
}

// settingsControl returns the control the caller bound in Trailing, or a
// zero element.
func settingsControl(c *ui.Context, o *SettingsRow) ui.Element {
	if o.Control == nil {
		return ui.Element{}
	}
	return c.Resolve(*o.Control)
}

func settingsPlainRow(c *ui.Context, t theme.Palette, o *SettingsRow, inset float32) {
	e := style.Disabled(settingsRowFrame(c, ui.Row(c).FillWidth(), inset), o.Disabled)
	e.Children(func() {
		title := settingsRowBody(c, t, o, true)
		if control := settingsControl(c, o); control.Valid() {
			control.LabelledBy(title)
			if o.Subtitle != "" {
				control.Description(o.Subtitle)
			}
		}
	})
}

// settingsActivatable gives a row that takes a click its hover face and an
// inset focus ring, which the card's clip would otherwise cut.
func settingsActivatable(c *ui.Context, t theme.Palette, e ui.Element, o *SettingsRow) ui.Element {
	e = style.Disabled(e.Label(o.Title).Description(o.Subtitle).FocusRing(false).Cursor(ui.CursorPointer), o.Disabled)
	e = style.Background(e, 0, ui.Transparent, t.SurfaceHover, t.SurfaceHover)
	if e.FocusVisible() {
		inset := theme.Scale(c, theme.Space4)
		radius := theme.RadiusScale().MD
		e.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.FocusRing(ui.Rect{X: r.X + inset, Y: r.Y + inset, W: r.W - 2*inset, H: r.H - 2*inset},
				[4]float32{radius, radius, radius, radius})
		})
	}
	return e
}

func settingsPressableRow(c *ui.Context, t theme.Palette, o *SettingsRow, inset float32) {
	e := settingsRowFrame(c, ui.ButtonBase(c).FillWidth().Justify(ui.Start), inset)
	e = settingsActivatable(c, t, e, o)
	if !o.Disabled {
		e.OnClick(o.OnClick)
	}
	e.Children(func() {
		settingsRowBody(c, t, o, true)
		chevron := lucide.ChevronRight()
		if c.Direction() == ui.RightToLeft {
			chevron = lucide.ChevronLeft()
		}
		size := theme.Scale(c, theme.IconMD)
		ui.Icon(c, chevron).Size(size, size).Shrink(0).TextColor(t.Muted)
	})
}

func settingsExpanderRow(c *ui.Context, t theme.Palette, o *SettingsRow, inset float32) {
	// Disabling the whole row also stops its trailing control and the
	// nested rows of an expander left open.
	// The open state belongs to the row's title, not its position, so rows
	// added above it do not hand it on.
	if o.Title != "" {
		c.Key(o.Title)
	}
	root := style.Disabled(ui.Column(c).FillWidth().MinWidth(0), o.Disabled)
	open := o.Expanded
	if open == nil {
		open = ui.Local(root, settingsExpandedKey{}, func() bool { return false })
	}
	// The rows below start where this row's title does.
	childInset := inset + settingsTextInset(c, o)
	if o.Icon == nil {
		childInset = inset + theme.Scale(c, theme.Space16)
	}
	root.Children(func() {
		var parts ui.CollapsibleParts
		ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Stretch).Children(func() {
			parts = ui.CollapsibleBase(c, open)
			trigger := settingsRowFrame(c, parts.Trigger.Grow(1).Basis(0), inset)
			trigger = settingsActivatable(c, t, trigger, o)
			progress := parts.Progress()
			chevron := func() {
				size := theme.Scale(c, theme.IconMD)
				ui.Icon(c, lucide.ChevronDown()).Size(size, size).Shrink(0).
					TextColor(t.Muted).Rotate(180 * progress)
			}
			trigger.Children(func() {
				settingsRowBody(c, t, o, false)
				if o.Trailing == nil {
					chevron()
				}
			})
			// Clicked applies the collapsible's input, so the panel below
			// reads the new value in this build.
			_ = trigger.Clicked()
			if o.Trailing != nil {
				// The control cannot sit inside the trigger button. It and the
				// chevron, which ends the row, take the room after it, with the
				// row's own padding at the end.
				end := theme.Scale(c, theme.Space16)
				slot := ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space12)).
					Padding(0, end, 0, 0)
				if c.Direction() == ui.RightToLeft {
					slot.Padding(0, 0, 0, end)
				}
				slot.Children(func() {
					settingsTrailing(c, t, o)
					// A press on the chevron toggles like the trigger; the
					// keyboard has the trigger.
					size := theme.Scale(c, theme.IconMD)
					ui.Box(c).Size(size, size).Shrink(0).Center().Cursor(ui.CursorPointer).
						Role(ui.RoleNone).OnClick(func() { *open = !*open }).Children(chevron)
				})
				if control := settingsControl(c, o); control.Valid() {
					control.Label(o.Title)
				}
			}
		})
		panel := parts.Panel(func() {
			for i := range o.Rows {
				Separator(c, SeparatorOptions{})
				settingsRow(c, t, &o.Rows[i], childInset)
			}
		})
		if panel.Valid() {
			panel.FillWidth().MinWidth(0).Background(t.MutedBackground.Alpha(0.5))
			// The panel fades as it grows; open, it keeps its own opacity.
			if progress := parts.Progress(); progress < 1 {
				panel.Opacity(max(progress, 0))
			}
		}
	})
}
