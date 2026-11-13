---
name: mygoui-motion
description: "State motion in mygoui UI. Use when adding, reviewing or testing animation in a mygoui application."
---

# mygoui motion

For work on the mygoui repository itself, also read [MAINTAINING.md](MAINTAINING.md).

mygoui components already move as these rules say; build them as they are and keep transitions for custom controls and layout. `DESIGN.md` owns the rules in its state motion section; read it when a case is missing below. Durations and curves are theme tokens: read them with `theme.Motion()` so custom controls keep pace with the library under any theme (`docs/theme.md` lists them). It ships with the module version your application uses (`go list -m -f '{{.Dir}}' github.com/Cyenoch/mygoui`; when vendored, the repository at that version). MyGo API details: `go doc github.com/egoist/mygo/ui Animation` (and `ElementTransition`, `Element.AnimateTo`).

## What moves

A change of state the user caused, or needs to notice, shows as motion:

| Change | Motion | Token of `theme.Motion()` |
| --- | --- | --- |
| Fill, border or ink following a state (checked, selected, invalid, status) | Color fade | `State`, `Ease` |
| A mark appearing (check, dot, tick) | Draws or scales in; leaves by fading | `State`, `Ease` |
| An icon standing for the state swapping | Crossfade | `State` |
| Content opening or closing | Height and fade | `Reveal`, `Ease` |
| An indicator moving to another item | Slide | `Reveal`, `EaseTravel` |
| A slider fill moving to a value the application set | Eases there; a held pointer is followed at once | `Reveal`, `Ease` |
| A number changing | Digits roll (`mygoui.NumberFlip`) | `Roll`, `Ease` |
| Neighbours making room during a drag | Slide | `Shift`, `Ease` |
| A page replacing another | Slide | `Page`, `Ease` |
| A quantity changing (progress, meter) | Eases to the value | none; the progress bar's own (`DESIGN.md`) |

Overlay surfaces have `Popup`, `Dialog` and `Drawer`; mygoui's overlays apply them. `docs/theme.md` lists every token's default.

## What stays instant

Hover and focus rings, menu, select and combo box popups and their row highlight, text being typed, anything being dragged, and a value present on the first frame. Text never changes weight with state.

## Reduced motion and rest

- Under the OS reduced-motion request (`theme.ReduceMotion(c)`), fades and color changes still play; sizes and positions settle at once. MyGo applies this for you when each value says what it is (below).
- At rest, nothing requests frames. A decorative loop (shimmer) runs only while its state lasts and stays still under reduced motion, as mygoui's skeletons do.

## MyGo APIs

- Read `m := theme.Motion()` in the build that animates; a package-level variable would miss a theme published later.
- **Color**: `e.AnimateColor(key, target, ui.Animation{Duration: m.State, Ease: m.Ease})`. Colors change smoothly under reduced motion.
- **Opacity, mark drawing, crossfade**: `e.AnimateTo(key, target, ui.Animation{Duration: ..., NonMotion: true})`. `NonMotion` keeps it playing under reduced motion.
- **Position, size, scale, scroll**: `e.AnimateTo` without `NonMotion`; under reduced motion it jumps to the target.
- **Layout change, enter and exit**: `.Transition(ui.ElementTransition{...})` on an element with a stable `c.Key`. Elements the change moves need a transition too. Name what animates (`Position`, `Size` or `Width`/`Height`, `Colors`, `Shadows`); naming none animates all of them, including colors you already animate.
- `ui.Animation` and `ElementTransition` without an `Ease` use `ui.EaseOut`, whatever the theme says; pass `m.Ease`, or `m.EaseTravel` for an indicator moving between items.
- A value starts at its first target and an `Enter` plays only in a parent the previous frame showed, so first-frame values arrive without motion. For a change that should not show as motion (restoring a window, jumping to a result), call `c.DisableAnimations()` before building it.

```go
// pin is a custom control: its fill follows the state, its mark fades in.
func pin(c *ui.Context, pinned bool) ui.Element {
	t := theme.Colors(c)
	size := theme.Scale(c, theme.IconMD)
	e := ui.Box(c).Width(size).Height(size).Radius(theme.RadiusScale().SM)
	m := theme.Motion()
	shown, fill := float32(0), t.MutedBackground
	if pinned {
		shown, fill = 1, t.Primary
	}
	mark := e.AnimateTo(markKey{}, shown, ui.Animation{Duration: m.State, Ease: m.Ease, NonMotion: true})
	e.Background(e.AnimateColor(fillKey{}, fill, ui.Animation{Duration: m.State, Ease: m.Ease}))
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X + r.W/4, Y: r.Y + r.H/4, W: r.W / 2, H: r.H / 2}, t.PrimaryForeground.Alpha(mark), 2)
	})
	return e
}

// reveal opens and closes content: height and fade together. It names
// Height alone, so position and colors change at once.
func reveal() ui.ElementTransition {
	m := theme.Motion()
	return ui.ElementTransition{
		Duration: m.Reveal, Ease: m.Ease, Height: true,
		Enter: &ui.Motion{Collapse: true}, Exit: &ui.Motion{Collapse: true},
	}
}

func details(c *ui.Context, open bool, body func()) {
	if open {
		ui.Column(c.Key("details")).Transition(reveal()).Children(body)
	}
}
```

(`markKey` and `fillKey` are empty struct types; any comparable value works as a key.)

## Test it

With `ui.NewTester`: change the state, `Frame()`, and `FramePending()` is true while it moves; `Advance(time.Second)` and it is false at rest. Repeat after `SetPreferences(ui.Preferences{ReduceMotion: true})`: positions and sizes land at once, fades still play.
