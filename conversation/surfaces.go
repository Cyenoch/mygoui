// Package conversation provides native message surfaces without agent-specific
// state. Callers own data, interaction and loading; views perform no I/O.
package conversation

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type Align uint8

const (
	Start Align = iota
	End
)

type MessageOptions struct {
	Align                 Align
	Label, Header, Footer string
	Avatar                func()
	// Content callbacks replace the corresponding metadata string.
	HeaderContent, FooterContent func()
}

// Message lays out optional metadata and an avatar around caller-owned content.
// The footer occupies its own measured row, so an avatar stays on the final
// bubble's baseline even when delivery information wraps or contains actions.
// Alignment is logical and follows the window's reading direction.
func Message(c *ui.Context, o MessageOptions, content func()) ui.Element {
	colors := theme.Colors(c)
	alignment, bodyColumn, avatarColumn := ui.Start, 2, 1
	tracks := [2]ui.Track{ui.FitContent(), ui.Fr(1)}
	if o.Align == End {
		alignment, bodyColumn, avatarColumn = ui.End, 1, 2
		tracks[0], tracks[1] = tracks[1], tracks[0]
	}
	trackCount := 2
	if o.Avatar == nil {
		bodyColumn, trackCount = 1, 1
		tracks[0] = ui.Fr(1)
	}
	message := ui.Grid(c).FillWidth().MinWidth(0).ColumnTracks(tracks[:trackCount]...).GapX(8).GapY(10)
	if o.Label != "" {
		message.Label(o.Label)
	}
	message.Children(func() {
		if o.Avatar != nil {
			ui.Column(c).ColumnStart(avatarColumn).RowStart(1).AlignSelf(ui.End).
				MinWidth(32).AlignItems(ui.Center).Children(o.Avatar)
		}
		ui.Column(c).ColumnStart(bodyColumn).RowStart(1).MinWidth(0).
			AlignItems(alignment).Gap(10).Children(func() {
			if o.Header != "" || o.HeaderContent != nil {
				ui.Column(c).MinWidth(0).MaxWidthPercent(100).Padding(0, 12).
					AlignItems(alignment).FontSize(theme.Scale(c, 12)).FontWeight(500).
					TextColor(colors.MutedForeground).Children(func() {
					if o.HeaderContent != nil {
						o.HeaderContent()
					} else {
						ui.Text(c, o.Header)
					}
				})
			}
			content()
		})
		if o.Footer != "" || o.FooterContent != nil {
			ui.Column(c).ColumnStart(bodyColumn).RowStart(2).MinWidth(0).Padding(0, 12).
				AlignItems(alignment).FontSize(theme.Scale(c, 12)).FontWeight(500).
				TextColor(colors.MutedForeground).Children(func() {
				if o.FooterContent != nil {
					o.FooterContent()
				} else {
					ui.Text(c, o.Footer)
				}
			})
		}
	})
	return message
}

func MessageGroup(c *ui.Context, content func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Gap(8).Children(content)
}

type BubbleVariant uint8

const (
	BubbleDefault BubbleVariant = iota
	BubbleSecondary
	BubbleMuted
	BubbleTinted
	BubbleOutline
	BubbleGhost
	BubbleDestructive
)

type Reaction struct {
	Key, Label         string
	AccessibilityLabel string
	Count              int
	Selected, Disabled bool
	OnClick            func()
}

type BubbleOptions struct {
	Variant        BubbleVariant
	Align          Align
	Label          string
	OnClick        func()
	Reactions      []Reaction
	ReactionsAbove bool
	// ReserveReactionSpace retains the badge's protruding footprint even when
	// no reactions are present. Use it in tightly spaced conversation stacks.
	ReserveReactionSpace bool
}

func Bubble(c *ui.Context, o BubbleOptions, content func()) ui.Element {
	colors := theme.Colors(c)
	foreground, background := colors.PrimaryForeground, colors.Primary
	hover, pressed, border := colors.PrimaryHover, colors.PrimaryPressed, ui.Color{}
	switch o.Variant {
	case BubbleSecondary:
		foreground, background = colors.SecondaryForeground, colors.Secondary
		hover, pressed = colors.SecondaryHover, colors.SecondaryPressed
	case BubbleMuted:
		foreground, background = colors.Foreground, colors.Muted
		hover = theme.Over(colors.Foreground, background.Over(colors.Background), 0.05)
		pressed = theme.Over(colors.Foreground, background.Over(colors.Background), 0.09)
	case BubbleTinted:
		foreground, background = colors.SelectionForeground, colors.Selection
		// Keep tinted content on the same contrast-resolved semantic pair.
		hover, pressed = background, background
	case BubbleOutline:
		foreground, background, border = colors.Foreground, colors.Background, colors.Border
		hover, pressed = colors.Muted, colors.AccentPressed
	case BubbleGhost:
		foreground, background = colors.Foreground, ui.Color{}
		hover, pressed = colors.Muted, colors.AccentPressed
	case BubbleDestructive:
		foreground, background = colors.Destructive, colors.ErrorSurface
		// The subdued destructive surface uses readable destructive ink.
		hover, pressed = background, background
	}
	outer := ui.Column(c).MinWidth(0).MaxWidthPercent(80)
	if o.Variant == BubbleGhost {
		outer.MaxWidthPercent(100)
	}
	if o.Align == End {
		outer.AlignSelf(ui.End)
	} else {
		outer.AlignSelf(ui.Start)
	}
	// Badge text has a fixed scaled line box; count/selection changes cannot
	// change the reserved height or the surface's bounds.
	badgeHeight := theme.Scale(c, 20) + 4 + 6
	overlap := min(float32(10), badgeHeight*0.25)
	if o.ReserveReactionSpace {
		if o.ReactionsAbove {
			outer.Margin(badgeHeight-overlap, 0, 0, 0)
		} else {
			outer.Margin(0, 0, badgeHeight-overlap, 0)
		}
	}
	outer.Children(func() {
		var surface ui.Element
		if o.OnClick != nil {
			surface = ui.ButtonBase(c).Column().Justify(ui.Start).AlignItems(ui.Start).OnClick(o.OnClick)
		} else {
			surface = ui.Column(c)
		}
		surface.MinWidth(0).MaxWidthPercent(100).Padding(10, 12).Radius(24).
			Border(1, border).Background(background).TextColor(foreground).
			FontSize(theme.Scale(c, 14)).LineHeight(1.5)
		if o.Variant == BubbleGhost {
			surface.Padding(0).Radius(0)
		}
		if o.Label != "" {
			surface.Label(o.Label)
		}
		if o.OnClick != nil {
			if surface.Pressed() {
				surface.Background(pressed)
			} else if surface.Hovered() {
				surface.Background(hover)
			}
			if surface.Hovered() {
				if o.Variant == BubbleTinted {
					surface.BorderColor(colors.Ring)
				}
				if o.Variant == BubbleDestructive {
					surface.BorderColor(colors.Destructive)
				}
			}
		}
		surface.Children(content)
		if len(o.Reactions) == 0 {
			return
		}
		at, self, offset := ui.AnchorBottomEnd, ui.AnchorTopEnd, -overlap
		if o.ReactionsAbove {
			at, self, offset = ui.AnchorTopEnd, ui.AnchorBottomEnd, overlap
		}
		// Attach to the bubble's box, not the text. Overlap is capped by the
		// surface padding even at enlarged text sizes.
		ui.Row(c).Attach(at, self).Top(offset).MaxWidthPercent(100).MinWidth(0).
			Padding(2, 3).Gap(2).Radius(badgeHeight/2).Border(3, colors.Card).
			Background(colors.Muted.Over(colors.Card)).Children(func() {
			for _, reaction := range o.Reactions {
				text, label := reaction.Label, reaction.AccessibilityLabel
				if label == "" {
					label = reaction.Label
				}
				if reaction.Count > 0 {
					text = locales.ComponentReactionCount(text, reaction.Count)
					label = locales.ComponentReactionAccessibleCount(label, reaction.Count)
				}
				control := ui.ButtonBase(c.Key(reaction.Key)).MinWidth(0).Shrink(1).
					Padding(0, 3).Radius(theme.Scale(c, 10)).Label(label).
					Role(ui.RoleToggleButton).Checked(reaction.Selected).
					Disabled(reaction.Disabled).OnClick(reaction.OnClick).
					TextColor(colors.Foreground)
				if reaction.Selected {
					control.Background(colors.Selection).TextColor(colors.SelectionForeground)
				} else if control.Hovered() && !reaction.Disabled {
					control.Background(colors.AccentHover)
				}
				if reaction.Disabled {
					control.TextColor(colors.DisabledForeground)
				}
				control.Children(func() {
					ui.Text(c, text).MinWidth(0).FontSize(theme.Scale(c, 14)).
						FixedLineHeight(theme.Scale(c, 20)).SingleLine()
				})
			}
		})
	})
	return outer
}

type MarkerVariant uint8

const (
	MarkerDefault MarkerVariant = iota
	MarkerSeparator
	MarkerBorder
)

type MarkerOptions struct {
	Text    string
	Icon    *ui.SVG
	Variant MarkerVariant
	OnClick func()
	Label   string
}

func Marker(c *ui.Context, o MarkerOptions) ui.Element {
	colors := theme.Colors(c)
	row := ui.Row(c).FillWidth().MinWidth(0).MinHeight(16).Gap(8)
	if o.Variant == MarkerBorder {
		row.Padding(0, 0, 8, 0).BorderWidth(0, 0, 1, 0).BorderColor(colors.Border)
	}
	row.Children(func() {
		rule := func() { ui.Box(c).Grow(1).MinWidth(0).Height(1).Background(colors.Border) }
		if o.Variant == MarkerSeparator {
			rule()
		}
		build := func() {
			if o.Icon != nil {
				ui.Icon(c, o.Icon).Size(16, 16).TextColor(colors.MutedForeground)
			}
			ui.Text(c, o.Text).FontSize(theme.Scale(c, 14)).TextColor(colors.MutedForeground)
		}
		if o.OnClick != nil {
			ui.ButtonBase(c).MinWidth(0).Gap(8).Label(o.Text).OnClick(o.OnClick).Children(build)
		} else {
			ui.Row(c).MinWidth(0).Gap(8).Children(build)
		}
		if o.Variant == MarkerSeparator {
			rule()
		}
	})
	return row
}

type AttachmentState uint8

const (
	AttachmentDone AttachmentState = iota
	AttachmentIdle
	AttachmentUploading
	AttachmentProcessing
	AttachmentError
)

type AttachmentSize uint8

const (
	AttachmentDefault AttachmentSize = iota
	AttachmentSM
	AttachmentXS
)

type AttachmentOptions struct {
	Title, Description        string
	State                     AttachmentState
	Size                      AttachmentSize
	Vertical                  bool
	Image                     ui.ImageSource
	Icon                      *ui.SVG
	OnOpen, OnRemove, OnRetry func()
	Disabled                  bool
	Progress                  *float64
}

func Attachment(c *ui.Context, o AttachmentOptions) ui.Element {
	colors := theme.Colors(c)
	padding, gap, media, size := float32(8), float32(8), float32(40), float32(14)
	if o.Size == AttachmentSM {
		padding, gap, media, size = 6, 10, 32, 12
	}
	if o.Size == AttachmentXS {
		padding, gap, media, size = 4, 6, 28, 12
	}
	var outer ui.Element
	if o.Vertical {
		outer = ui.Column(c).Width(120)
	} else {
		outer = ui.Row(c).MinWidth(160)
	}
	outer.MinWidth(0).MaxWidthPercent(100).Padding(padding).Gap(gap).Radius(16).Border(1, colors.Border).Background(colors.Card).TextColor(colors.CardForeground)
	if o.State == AttachmentError {
		outer.BorderColor(colors.Destructive)
	}
	outer.Children(func() {
		var trigger ui.Element
		if o.OnOpen != nil {
			trigger = ui.ButtonBase(c).Justify(ui.Start).OnClick(o.OnOpen).Disabled(o.Disabled).Label(locales.ComponentOpenAttachment(o.Title))
		} else {
			trigger = ui.Row(c)
		}
		if o.Vertical {
			trigger.Column().FillWidth()
		}
		trigger.MinWidth(0).Gap(gap).Grow(1)
		trigger.Children(func() {
			mediaBox := ui.Box(c).Size(media, media).Radius(8).Clip().Background(colors.Muted).Center().Shrink(0)
			if o.Vertical {
				mediaBox.Size(104, 104)
			}
			mediaBox.Children(func() {
				switch {
				case o.Image != nil:
					ui.Image(c, o.Image).Fill().Fit(ui.Cover).Label(o.Title)
				case o.State == AttachmentUploading || o.State == AttachmentProcessing:
					mygoui.Spinner(c, mygoui.SpinnerOptions{Size: 16, Label: attachmentStatus(o.State)})
				default:
					icon := o.Icon
					if icon == nil {
						icon = lucide.FileText()
					}
					if o.State == AttachmentError {
						icon = lucide.CircleAlert()
					}
					ui.Icon(c, icon).Size(16, 16).TextColor(colors.Foreground)
				}
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
				ui.Text(c, o.Title).FontSize(theme.Scale(c, size)).FontWeight(500).SingleLine().Tooltip(o.Title)
				description := o.Description
				if description == "" {
					description = attachmentStatus(o.State)
				}
				ink := colors.MutedForeground
				if o.State == AttachmentError {
					ink = colors.Destructive
				}
				ui.Text(c, description).FontSize(theme.Scale(c, 12)).TextColor(ink).SingleLine().Tooltip(description)
				if o.Progress != nil {
					ui.Progress(c, *o.Progress).FillWidth().Label(locales.ComponentAttachmentUploadProgress())
				}
			})
		})
		if o.OnRetry != nil && o.State == AttachmentError {
			mygoui.Button(c, mygoui.ButtonOptions{Label: locales.ComponentRetryAttachment(o.Title), Icon: lucide.RotateCcw(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconXS, Disabled: o.Disabled, OnClick: o.OnRetry})
		}
		if o.OnRemove != nil {
			mygoui.Button(c, mygoui.ButtonOptions{Label: locales.ComponentRemoveAttachment(o.Title), Icon: lucide.X(), Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconXS, Disabled: o.Disabled, OnClick: o.OnRemove})
		}
	})
	return outer
}
func AttachmentGroup(c *ui.Context, content func()) ui.Element {
	return ui.ScrollHorizontal(c).FillWidth().MinWidth(0).Padding(4, 0).Children(func() { ui.Row(c).Gap(12).Children(content) })
}
func attachmentStatus(state AttachmentState) string {
	switch state {
	case AttachmentIdle:
		return locales.ComponentNotUploaded()
	case AttachmentUploading:
		return locales.ComponentUploading()
	case AttachmentProcessing:
		return locales.ComponentProcessing()
	case AttachmentError:
		return locales.ComponentUploadFailed()
	default:
		return locales.ComponentReady()
	}
}
