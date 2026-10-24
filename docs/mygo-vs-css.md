# MyGo vs CSS

Concrete differences between CSS/Tailwind (HeroUI's source) and MyGo's native
layout and paint model, and the rules this library follows because of them.
Every item below changed a component's implementation.

## Layout

| CSS / Tailwind | MyGo |
| --- | --- |
| `display:flex` row stretches items | `ui.Row` sets `AlignItems(Center)` unless an alignment is set. Rows that must top-align call `AlignItems(ui.Start)`. |
| `align-items: stretch` | `ui.Column` keeps the automatic cross-axis alignment, which layout resolves to `Stretch`. |
| `align-items: stretch` on one child | `AlignSelf(ui.Start)` (or `Center`/`End`) hugs the child's own size across the axis. |
| `width:auto` / `width:fit-content` | There is no such width token. An element hugs its content by default; use `Shrink(0)` on the main axis and `AlignSelf(ui.Start)` across it. `FillWidth()` is `width:100%`. |
| `margin:auto` | `Margin(0, ui.Auto)` (`ui.Auto` is `+Inf`) takes the free space; `MarginStart(ui.Auto)` pushes a row child to the end. |
| `flex-shrink:1` (the CSS default) | Items do not shrink unless `Shrink(f)` is set; `Shrink(1)` opts in, and `MinWidth(0)`/`MinHeight(0)` lets an item go below its content. |
| `flex-shrink:0` | `Shrink(0)` is the default and affects only the main axis. It never stops a column from stretching a child's width. |
| `flex:1 1 0` | `Grow(1).Basis(0).MinWidth(0)`; `Basis(0)` is required or content width wins. |
| `gap` | `Gap` applies between items only, not at the ends; padding handles the ends. |
| `position:absolute; inset:0` | `Absolute()` plus `Top/Right/Bottom/Left(0)` fills the parent's padding box, and absolute children paint after flow children. Anchored UI uses `Attach`/`AttachTo` instead: MyGo flips and clamps it into the window. |
| `z-index` | There is no z-index or stacking context. Children paint in creation order, absolute children last, and `DrawOver` after children. The overlay layer (`ui.Overlay`) paints above the page. |

`Start` is the zero `ui.Align`, so `Align == 0` is a real value (left/top),
never a sentinel for "unset". Where the zero value must mean centered,
components define their own enum (for example `ButtonAlignCenter`) instead of
testing `Align != 0`.

## Paint

- `Element.Background(color)` paints the element's rounded fill, with the
  element's per-corner radius, before anything else. `Element.Draw(fn)` paints
  after that background and before children — and **before the element's clip
  is pushed**. `Clip()` shapes children, not the element's own background or
  `Draw`.
- `Painter.Fill(rect, color, radius)` takes a single radius; `Radius(r0, r1,
  r2, r3)` and `Clip()` are the only per-corner shaping APIs. A per-corner
  state face therefore cannot be a `Draw` fill: `p.Fill(r, fill, 0)` squares
  off the rounded ends of a joined group, and no radius argument can express
  `{8, 0, 0, 8}`. Paint it in a child layer (whose own clip shapes it) or as
  the element's background.
- Hover and pressed are paint-time states: read `e.Hovered()` and `e.Pressed()`
  inside `Draw`/`DrawOver`, as the theme's state fills do. The read also marks
  the element as watching hover, which is how the toolkit knows to schedule
  the repaint when the pointer moves.
- There is no alpha mask. `ScrollShadow` fades edges by painting an opaque
  surface color, not by masking content (see `docs/components.md`).
- There is no CSS pressed scale or ripple: layouts do not transform live
  element boxes. `Transition`, `AnimatePose` and `AnimateColor` animate
  committed boxes and colors; press feedback is color-only.
- Focus is `FocusRing(true)` on focusable elements: a 2 DIP ring starting 1 DIP
  outside the box, not an outline. An ancestor clip cuts it, so clipping a
  group to round its children also clips their rings. Joined controls turn the
  default ring off and draw an inset ring instead: the next segment of the
  group paints over an outward ring, while there is no z-index to raise the
  focused segment as HeroUI's `z-10` does.
- Gradients are `LinearGradient`/`ConicGradient` per element; there are no
  layered CSS backgrounds.

## Text and content

- Type sizes are absolute: a child inherits the parent's set size, but there is
  no `em`/`rem` relative scaling. `theme.Scale(c, dip)` multiplies the HeroUI
  DIP metric by the host text scale and base size; every size, line height,
  icon metric and spacing that tracks text goes through it. Controls use
  `FixedLineHeight` so their rows do not shift with content.
- There are no selectors or `group-hover:` / `data-[selected]:` variants. A
  parent cannot style descendants, so state is passed explicitly to the
  builder or read from the owning element (`Hovered`, `Pressed`, `Checked`).
- Lists and tables are native and own virtualization, scroll, selection, row
  building and column state (`ui.List`, `ui.Table`). Rows are built inside the
  native callback; rebuilding them in Go state fights the toolkit.

## Translating a Tailwind class set

1. Map `flex`/`grid` to `ui.Row`/`ui.Column`/`ui.Grid` and set alignment
   explicitly; never rely on CSS defaults.
2. Map `w-full` to `FillWidth`, `flex-1` to `Grow(1).Basis(0).MinWidth(0)`,
   and a hugged box to `Shrink(0)` plus `AlignSelf(ui.Start)` when the parent
   stretches.
3. Map `rounded-*` to `Radius`, `overflow-hidden` to `Clip()`; remember the
   clip applies to children only.
4. Map `bg-*` plus `hover:bg-*`/`active:bg-*` to the element background plus a
   state fill (`internal/style.Background`, or `CornerFill` when the corners
   differ), reading `Hovered()` and `Pressed()` at paint time.
5. Map `outline`/`ring` to `FocusRing(true)`, not to a border.
6. Keep layout out of paint: a box that depends on the container comes from
   layout. Inside `Draw`, `e.LayoutBounds()` is in the same layout coordinate
   space as the rect MyGo hands you; `e.Bounds()` is transformed to the window.
