# MyGo and CSS translation

coss ui expresses appearance through Tailwind classes and Base UI parts.
The native implementation uses MyGo layout, widget bases, effects and services.
These are the translation rules, with source provenance in [references.md](references.md).

## Layout and identity

| CSS/Tailwind concept | Native mapping |
| --- | --- |
| Flex row | `ui.Row`; cross-axis default is Center, so set `AlignItems(ui.Start)` for top alignment |
| Flex column stretch | `ui.Column` automatic cross-axis alignment resolves to Stretch |
| Fit-content child in a stretched column | `AlignSelf(ui.Start)` and appropriate `Shrink(0)` |
| `w-full` | `FillWidth()` |
| `flex: 1 1 0` | `Grow(1).Basis(0).MinWidth(0)`; content basis otherwise affects allocation |
| `flex-shrink: 1` | Opt into `Shrink(1)` and allow a suitable zero minimum |
| `margin: auto` | `ui.Auto`, including `MarginStart(ui.Auto)` for a trailing row action |
| Gap | `Gap` spaces children; padding supplies edge insets |
| Absolute inset | `Absolute` and Top/Right/Bottom/Left; anchored overlays use native attachment/placement |
| CSS z-index | Creation/paint order and native overlay layer; no arbitrary stacking-context system |
| DOM key/ref | Stable pre-construction `c.Key` and persistent `ui.Handle` |

`ui.Start` is the real zero alignment value. Do not treat it as unset.
A public option whose zero value should center content needs its own enum.
Native collections own virtualization, scroll, selection and column state;
applications retain their state rather than maintaining parallel row caches.

## Paint, clipping and effects

The element background paints before its Draw callback, which paints before
children. The element's own clip shapes children, not an arbitrary fill
already painted in Draw. Per-corner joined faces therefore need the actual
rounded element background or a correctly clipped child layer.

Capture the chosen hover/pressed/focus face while the element is valid in its
build, or resolve current geometry through supported handle-based facilities.
Never close over an expired element in a retained transition's paint callback.
A draw may run again after consumer content has left the live tree.

Native keyboard focus rings can extend outside a control. Ancestor clipping
or the next joined segment can cover them; joined controls need intentional
inset-ring geometry instead of CSS z-index tricks. Disabled painting is native
50% opacity, compared with the coss button reference's 64%.

`Element.Effect` renders an element and its native children to an intermediate
layer. `ui/shader` supplies `ColorEffect`, `DistortionEffect` and `LayerEffect`
entry points, with GPU execution and a CPU fallback. A ColorEffect can multiply
subtree alpha for a real edge mask; painting an opaque surface gradient over
content is a different treatment. ScrollFade's mask mode uses that native
facility. See [scrollfade.md](scrollfade.md).

The pinned `ui/effect.go` limits each effect-layer side to 8192 **pixels**,
including scale and effect margins. Above that size MyGo paints the subtree
without the effect. Embedded webviews do not show through an effect layer.
Headless fallback output does not establish actual GPU presentation.

Button feedback uses coss color/border/shadow states, without a pressed scale.
MyGo has pose animation APIs for suitable native transitions; their existence
does not require applying scale feedback to coss buttons.

## Text, controls and services

- Native type metrics are DIP values. `theme.Scale(c, dip)` incorporates OS
  text scale and the configured base size. Use explicit line boxes where a
  control needs stable vertical geometry.
- Family names select registered/installed fonts. Body, heading and mono are
  independent preferences; system fallbacks are not equivalent to the site's
  Cal Sans/Paper Mono. See [fonts.md](fonts.md).
- CSS descendant selectors and data attributes become explicit options or
  native state. Do not create an ambient provider to imitate selectors.
- Native text editors preserve selection, caret and IME contracts. HTML input
  types, browser autofill and automatic form validation have no blanket native
  equivalent.
- Native context menus use OS chrome. Ordinary Menu/Dropdown use native
  elements and typed menu data; they are separate presentation paths.
- Tooltips retain MyGo's pointer-rest and warm-reuse timing. Clipboard and URL
  opening use captured `ui.Services`, not browser hooks.

## Validation boundary

`ui.NewTester` runs native layout, event routing, semantic metadata and
headless painting. Inspecting pixels there can establish mask transparency or
state paint under that renderer; it cannot prove GPU behavior, OS menu chrome,
platform IME/gestures or screen-reader output. Record target-specific evidence
separately in [VALIDATION.md](../VALIDATION.md).
