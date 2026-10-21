// Package trace presents compact activity lines and disclosures with native input
// and accessibility. Appearance follows shadcn's Base Trace and base-rhea style;
// persistent state, composition and focus follow MyGo's native model.
package trace

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/theme"
)

// Status affects presentation only. Titles must also describe the state in words.
type Status uint8

const (
	Default Status = iota
	Running
	Error
	Pending
)

// State belongs to one disclosure in one window. Keep it with the caller's data
// when a row can leave a virtual list or its parent can close.
type State struct {
	Open        bool
	initialized bool
	wasOpen     bool
	trigger     ui.Handle
	content     ui.Handle
}

// OpenChange observes a disclosure's committed expansion value.
type OpenChange func(bool)

// Options describes a line and its disclosure. Open binds controlled expansion;
// otherwise DefaultOpen initializes State.Open once. OnOpenChange observes the
// committed value. Body content is supplied as a native build callback.
type Options struct {
	Title        string
	Description  string
	Icon         *ui.SVG
	Status       Status
	DefaultOpen  bool
	Open         *bool
	Disabled     bool
	OnOpenChange OpenChange
}

// Trace groups related lines without deriving or inventing activity membership.
func Trace(c *ui.Context) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space8))
}

// Line is static text, without button semantics or a keyboard tab stop.
func Line(c *ui.Context, status Status) ui.Element {
	e := ui.Row(c).MinWidth(0).MinHeight(theme.Scale(c, theme.Space20)).Gap(theme.Scale(c, theme.Space8)).
		Font(theme.FontFamily(c)).FontSize(theme.Scale(c, theme.Space14)).FixedLineHeight(theme.Scale(c, theme.Space20)).
		TextColor(lineColor(c, status))
	if status == Pending {
		e.Opacity(0.6)
	}
	return e
}

// Icon is decorative; the title carries the accessible meaning.
func Icon(c *ui.Context, icon *ui.SVG) ui.Element {
	return ui.Icon(c, icon).Size(theme.Scale(c, theme.IconMedium), theme.Scale(c, theme.IconMedium)).Shrink(0)
}

// Title truncates on one line. The glyph-only shimmer preserves native text
// layout, truncation and AX, and is omitted entirely for reduced motion.
func Title(c *ui.Context, text string, status Status) ui.Element {
	e := ui.Text(c, text).MinWidth(0).Shrink(1).SingleLine().
		Font(theme.FontFamily(c)).FontSize(theme.Scale(c, theme.Space14)).FixedLineHeight(theme.Scale(c, theme.Space20))
	if status == Running && !c.Preferences().ReduceMotion {
		fx := ui.Local(e, "trace-shimmer", func() shimmer { return shimmer{} })
		fx.phase = e.Loop("trace-shimmer", 2*time.Second, ui.Linear)
		e.Effect(fx)
	}
	return e
}

// Description reserves its natural width and never ellipsizes counts or timing.
func Description(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).Shrink(0).NoWrap().FontFeatures("tnum").
		FontSize(theme.Scale(c, theme.Space14)).FixedLineHeight(theme.Scale(c, theme.Space20)).
		TextColor(theme.Colors(c).MutedForeground)
}

// Collapsible builds a native disclosure around caller-provided content. A nil
// state uses keyed native local state; use caller-owned State for virtualization.
func Collapsible(c *ui.Context, state *State, options Options, content func()) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0)
	if state == nil {
		state = ui.Local(root, "trace-state", func() State { return State{} })
	}
	root.Children(func() {
		Trigger(c.Key("trigger"), state, options)
		Content(c.Key("content"), state, content)
	})
	return root
}

// Trigger composes the native button, state binding and rotating chevron. It may
// also be used with Content inside a custom native layout and persistent State.
func Trigger(c *ui.Context, state *State, options Options) ui.Element {
	b := ui.ButtonBase(c).MinWidth(0).MinHeight(theme.Scale(c, theme.HeightWorkRow)).
		Justify(ui.Start).Gap(theme.Scale(c, theme.Space4)).Radius(theme.RadiusControl).
		Label(options.Title).Role(ui.RoleDisclosure).Disabled(options.Disabled)
	if state == nil {
		state = ui.Local(b, "trace-state", func() State { return State{} })
	}
	if !state.initialized {
		state.Open = state.Open || options.DefaultOpen
		state.initialized = true
	}
	if options.Open != nil {
		state.Open = *options.Open
	}
	b.Bind(&state.trigger).Expanded(state.Open)
	if state.wasOpen && !state.Open && state.content.FocusWithin(c) {
		state.trigger.Focus(c)
	}
	state.wasOpen = state.Open
	color := interactionColor(c, b, options.Status)
	b.TextColor(color)
	angle := float32(0)
	if state.Open {
		angle = 90
	}
	angle = b.Animate("trace-chevron", angle, 200*time.Millisecond)
	b.Children(func() {
		lineContents(c, options, color)
		Icon(c, lucide.ChevronRight()).Rotate(angle)
	})
	b.OnClick(func() {
		state.Open = !state.Open
		if options.Open != nil {
			*options.Open = state.Open
		}
		if !state.Open {
			state.trigger.Focus(c)
		}
		if options.OnOpenChange != nil {
			options.OnOpenChange(state.Open)
		}
	})
	return b
}

// Content builds only visible content. Nested caller state survives closing;
// closed disclosures create no scroll surface, empty output or extra spacer.
func Content(c *ui.Context, state *State, content func()) ui.Element {
	if state == nil || !state.Open || content == nil {
		return ui.Element{}
	}
	return ui.Column(c).FillWidth().MinWidth(0).Bind(&state.content).
		Margin(theme.Scale(c, theme.Space8), 0, theme.Scale(c, theme.Space4), theme.Scale(c, theme.Space8-theme.HairlineWidth)).
		BorderWidth(0, 0, 0, theme.HairlineWidth).BorderColor(theme.Colors(c).Border).
		Padding(0, theme.Scale(c, theme.Space4), 0, theme.Scale(c, theme.Space16)).
		TextColor(theme.Colors(c).MutedForeground).Children(content)
}

// Action supplies actual button semantics for an activity that opens another
// surface. A missing action disables the button.
func Action(c *ui.Context, options Options, action func()) ui.Element {
	return interactiveLine(c, options, ui.RoleButton, action)
}

// Link opens the actual URL and exposes native link semantics. Services is safe
// to retain in an action; no build Context or Element escapes the view.
func Link(c *ui.Context, options Options, url string) ui.Element {
	services := c.Services()
	if url == "" {
		options.Disabled = true
	}
	return interactiveLine(c, options, ui.RoleLink, func() { services.OpenURL(url) })
}

func interactiveLine(c *ui.Context, options Options, role ui.Role, action func()) ui.Element {
	b := ui.ButtonBase(c).Justify(ui.Start).MinWidth(0).MinHeight(theme.Scale(c, theme.HeightWorkRow)).
		Radius(theme.RadiusControl).Role(role).Label(options.Title).Disabled(options.Disabled || action == nil)
	color := interactionColor(c, b, options.Status)
	b.TextColor(color).Children(func() { lineContents(c, options, color) })
	if action != nil {
		b.OnClick(action)
	}
	return b
}

func lineContents(c *ui.Context, options Options, color ui.Color) {
	Line(c, options.Status).MinWidth(0).Shrink(1).TextColor(color).Children(func() {
		if options.Icon != nil {
			Icon(c, options.Icon)
		}
		Title(c, options.Title, options.Status)
		if options.Description != "" {
			Description(c, options.Description)
		}
	})
}

func lineColor(c *ui.Context, status Status) ui.Color {
	colors := theme.Colors(c)
	if status == Error {
		return colors.Destructive
	}
	return colors.MutedForeground
}

func interactionColor(c *ui.Context, e ui.Element, status Status) ui.Color {
	colors := theme.Colors(c)
	if e.IsDisabled() {
		return colors.DisabledForeground
	}
	if status == Error {
		return colors.Destructive
	}
	if e.Hovered() || e.Focused() {
		return colors.Foreground
	}
	return colors.MutedForeground
}

// A moving alpha band on the rendered glyphs gives shimmer without splitting
// Unicode text, reshaping source text or creating per-frame text copies.
type shimmer struct{ phase float32 }

func (*shimmer) EffectMargin() ui.Point { return ui.Point{} }

func (s *shimmer) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	scale := p.Scale()
	width, height := content.Size()
	source := ui.Rect{X: inset.X, Y: inset.Y, W: box.W * scale, H: box.H * scale}
	if width == 0 || height == 0 || box.W <= 0 {
		return
	}
	// Each band samples the same cached glyph bitmap, never paints its background.
	const bands = 24
	for i := range bands {
		left := box.W * float32(i) / bands
		right := box.W * float32(i+1) / bands
		distance := float32(math.Abs(float64(float32(i)/bands - (s.phase*1.6 - 0.3))))
		alpha := float32(0.7) + 0.3*max(0, 1-distance/0.2)
		part := ui.Rect{X: source.X + left*scale, Y: source.Y, W: (right - left) * scale, H: source.H}
		dest := ui.Rect{X: box.X + left, Y: box.Y, W: right - left, H: box.H}
		p.ImageRegion(content, part, dest, ui.ImageOptions{Tint: ui.RGBA(255, 255, 255, alpha)})
	}
}
