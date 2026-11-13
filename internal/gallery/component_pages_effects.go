package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/effect"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// effectsState is the caller-owned state behind the effect demos.
type effectsState struct {
	maskFallback bool

	loading, loadingSet bool
	pulse               bool

	revealRun                     int
	revealForward, revealFallback bool
}

func effectsDemo(body func(c *ui.Context, st *effectsState)) func(*ui.Context, *State) {
	return func(c *ui.Context, s *State) {
		st := StateOf[effectsState](s, "effects")
		if !st.loadingSet {
			st.loading, st.loadingSet = true, true
		}
		ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() { body(c, st) })
	}
}

func effectSwitches(c *ui.Context, body func()) {
	ui.Row(c).FillWidth().MinWidth(0).Wrap().AlignItems(ui.Center).Gap(theme.Space16).Children(body)
}

// effectBackdrop is a varied background, which a mask lets through and an
// overpainted surface colour cannot match.
func effectBackdrop(c *ui.Context, body func()) {
	t := theme.Colors(c)
	radius := theme.RadiusScale().LG
	from, to := t.InfoSoft.Over(t.Background), t.SuccessSoft.Over(t.Background)
	ui.Column(c).FillWidth().MinWidth(0).Radius(radius).
		Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillGradient(r, ui.LinearGradient{From: from, To: to, Angle: 90}, radius)
		}).Children(body)
}

func maskListBody(c *ui.Context, st *effectsState) {
	effectSwitches(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.maskFallback, Label: "Fallback"})
	})
	effectBackdrop(c, func() {
		radius := theme.RadiusScale().LG
		// Clip cuts the overflowing rows in the fallback, as the effect's
		// layer does on the shader path.
		list := ui.Column(c).FillWidth().Height(theme.Scale(c, 184)).Padding(theme.Space12).Gap(theme.Space8).
			Radius(radius).Clip()
		list.Children(func() {
			for i := range 5 {
				mygoui.Button(c.Key(i), mygoui.ButtonOptions{
					Label: fmt.Sprintf("Release 2.%d", 6-i), Variant: mygoui.ButtonOutline, FullWidth: true,
				})
			}
		})
		effect.Mask(c, list, effect.MaskOptions{
			Edges: effect.EdgeBottom, Size: 48, Radius: radius, Fallback: st.maskFallback,
		})
	})
}

func maskSidesBody(c *ui.Context, st *effectsState) {
	effectBackdrop(c, func() {
		strip := ui.Row(c).FillWidth().MinWidth(0).Padding(theme.Space12).Gap(theme.Space8).Clip()
		strip.Children(func() {
			for _, tag := range []string{"macOS", "Windows", "Linux", "WebAssembly", "arm64", "amd64", "Metal", "Direct3D", "OpenGL"} {
				mygoui.Badge(c.Key(tag), mygoui.BadgeOptions{Label: tag, Variant: mygoui.BadgeOutline})
			}
		})
		effect.Mask(c, strip, effect.MaskOptions{
			Edges: effect.EdgeLeft | effect.EdgeRight, Fallback: st.maskFallback,
		})
	})
}

// effectSurface is the card a loading or arriving block sits on.
func effectSurface(c *ui.Context, body func()) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Space16).Radius(theme.RadiusScale().LG).
		Background(t.Surface).Border(theme.Hairline, t.Border).Children(body)
}

// effectPlaceholder is a loading card in the application's own layout.
func effectPlaceholder(c *ui.Context) ui.Element {
	row := ui.Row(c).FillWidth().Gap(theme.Space12).AlignItems(ui.Start)
	row.Children(func() {
		mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 40, Height: 40, Circle: true, Animation: mygoui.SkeletonNone})
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space8).Children(func() {
			for _, width := range []float32{180, 0, 120} {
				mygoui.Skeleton(c.Key(width), mygoui.SkeletonOptions{Width: width, Height: 12, Animation: mygoui.SkeletonNone})
			}
		})
	})
	return row
}

func shimmerBody(c *ui.Context, st *effectsState) {
	effectSwitches(c, func() {
		form.Switch(c, form.SwitchOptions{On: &st.loading, Label: "Loading"})
		form.Switch(c, form.SwitchOptions{On: &st.pulse, Label: "Pulse"})
	})
	effectSurface(c, func() {
		if !st.loading {
			ui.Column(c).Gap(theme.Space4).Children(func() {
				mygoui.Typography(c, "Ada Lovelace", mygoui.TypographyOptions{Type: theme.TypeH6})
				mygoui.Paragraph(c, "Wrote the first published program.", true)
			})
			return
		}
		placeholder := effectPlaceholder(c).Label("Loading profile").Busy(true)
		effect.Shimmer(c, placeholder, effect.ShimmerOptions{Pulse: st.pulse})
	})
}

func revealBody(c *ui.Context, st *effectsState) {
	effectSwitches(c, func() {
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Generate", Variant: mygoui.ButtonSecondary, Size: theme.ButtonSmall,
			OnClick: func() { st.revealRun++ },
		})
		form.Switch(c, form.SwitchOptions{On: &st.revealForward, Label: "Forward"})
		form.Switch(c, form.SwitchOptions{On: &st.revealFallback, Label: "Fallback"})
	})
	if st.revealRun == 0 {
		return
	}
	direction := effect.RevealDown
	if st.revealForward {
		direction = effect.RevealForward
	}
	effectSurface(c, func() {
		// A new key is a new element, so each run reveals again.
		summary := ui.Column(c.Key(st.revealRun)).FillWidth().Gap(theme.Space4)
		summary.Children(func() {
			mygoui.Typography(c, fmt.Sprintf("Summary %d", st.revealRun), mygoui.TypographyOptions{Type: theme.TypeH6})
			mygoui.Paragraph(c, "Build finished in 42 seconds. 128 tests passed, none skipped. Artifacts were uploaded to the release bucket and signed.", true)
		})
		effect.Reveal(c, summary, effect.RevealOptions{Direction: direction, Fallback: st.revealFallback})
	})
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "mask", Title: "Mask", Group: "Foundations", Order: 103,
		Description: "Fades the edges of any element to transparency where its content is cut.",
		Examples: []ComponentExample{
			{Title: "Bottom edge over a gradient", Demo: effectsDemo(maskListBody)},
			{Title: "Both sides", Demo: effectsDemo(maskSidesBody)},
		},
		Code: `preview := ui.Column(c).Height(184).Radius(10).Clip()
preview.Children(func() { /* rows that continue past the bottom */ })
effect.Mask(c, preview, effect.MaskOptions{
    Edges:  effect.EdgeBottom,
    Size:   48,
    Radius: 10,
})`,
		Params: []ComponentParam{
			{"e", "ui.Element", "required", "The element to fade, after its children and styling. Use ScrollFade for a scroll viewport."},
			{"Edges", "effect.Edges", "all", "Physical EdgeTop, EdgeRight, EdgeBottom and EdgeLeft flags."},
			{"Size", "float32", "Effects.FadeEdge", "Fade depth in DIPs."},
			{"Surface", "ui.Color", "shadow", "Opaque color behind the element; the fallback blends into it."},
			{"Radius", "float32", "0", "The element's corner radius, which the fallback follows."},
			{"Fallback", "bool", "false", "Shades the edges instead of masking them, as without a GPU."},
		},
		Tokens: []string{"Foreground"},
	})

	RegisterComponent(ComponentPage{
		ID: "shimmer", Title: "Shimmer", Group: "Foundations", Order: 104,
		Description: "One loading highlight across placeholders in your own layout.",
		Examples: []ComponentExample{
			{Title: "Loading card", Demo: effectsDemo(shimmerBody)},
		},
		Code: `var loading bool

line := mygoui.SkeletonOptions{Height: 12, Animation: mygoui.SkeletonNone}
placeholder := ui.Column(c).Gap(8).Label("Loading profile").Busy(true)
placeholder.Children(func() {
    mygoui.Skeleton(c, line)
    mygoui.Skeleton(c, line)
})
if loading {
    effect.Shimmer(c, placeholder, effect.ShimmerOptions{})
}`,
		Params: []ComponentParam{
			{"e", "ui.Element", "required", "The element holding the placeholders; apply it only while loading."},
			{"Pulse", "bool", "false", "Fades the content to half and back instead of sweeping a band."},
		},
		Tokens: []string{"MutedBackground"},
	})

	RegisterComponent(ComponentPage{
		ID: "reveal", Title: "Reveal", Group: "Foundations", Order: 105,
		Description: "Wipes arriving content in once, behind a soft front.",
		Examples: []ComponentExample{
			{Title: "Arriving summary", Demo: effectsDemo(revealBody)},
		},
		Code: `var arrived bool

if arrived {
    summary := ui.Column(c).Gap(4)
    summary.Children(func() { /* the content that arrived */ })
    effect.Reveal(c, summary, effect.RevealOptions{})
}`,
		Params: []ComponentParam{
			{"e", "ui.Element", "required", "Content that arrives; it reveals when first built."},
			{"Direction", "effect.RevealDirection", "RevealDown", "RevealDown, RevealUp, RevealForward or RevealBackward."},
			{"Fallback", "bool", "false", "Fades in without a travelling front, as under reduced motion."},
		},
	})
}
