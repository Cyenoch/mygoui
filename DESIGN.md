---
name: mygoui
description: Native coss components and a compact MyGo workbench.
colors:
  light-background: "#ffffff"
  light-foreground: "oklch(26.9% 0 0)"
  light-primary: "oklch(26.9% 0 0)"
  light-primary-foreground: "oklch(98.5% 0 0)"
  light-accent: "rgb(0 0 0 / 4%)"
  light-border: "rgb(0 0 0 / 8%)"
  light-input: "rgb(0 0 0 / 10%)"
  light-ring: "oklch(70.8% 0 0)"
  dark-foreground: "oklch(97% 0 0)"
  dark-primary: "oklch(97% 0 0)"
  dark-primary-foreground: "oklch(26.9% 0 0)"
  dark-accent: "rgb(255 255 255 / 4%)"
  dark-border: "rgb(255 255 255 / 6%)"
  dark-input: "rgb(255 255 255 / 8%)"
  dark-ring: "oklch(55.6% 0 0)"
typography:
  body:
    fontFamily: "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"
    fontSize: "16px"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontSize: "14px"
    fontWeight: 500
    lineHeight: 1.428571
  help:
    fontSize: "12px"
    fontWeight: 400
    lineHeight: 1.333333
  heading:
    fontSize: "36px"
    fontWeight: 600
    lineHeight: 1.111111
    letterSpacing: "-0.025em"
  code:
    fontFamily: "ui-monospace, SF Mono, SFMono-Regular, Menlo, Consolas, Liberation Mono, monospace"
    fontSize: "14px"
    fontWeight: 400
    lineHeight: 1.428571
rounded:
  xs: "4px"
  sm: "6px"
  md: "8px"
  lg: "10px"
  xl: "14px"
  xxl: "16px"
  xxxl: "24px"
  xxxxl: "32px"
spacing:
  unit: "4px"
  compact: "8px"
  control: "12px"
  section: "16px"
  panel: "24px"
  page: "32px"
components:
  button-primary-light:
    backgroundColor: "{colors.light-primary}"
    textColor: "{colors.light-primary-foreground}"
    height: "32px"
    rounded: "{rounded.lg}"
  button-primary-dark:
    backgroundColor: "{colors.dark-primary}"
    textColor: "{colors.dark-primary-foreground}"
    height: "32px"
    rounded: "{rounded.lg}"
  input:
    height: "32px"
    rounded: "{rounded.md}"
---

# Design System: mygoui

## Overview

Compact neutral surfaces, clear text hierarchy and restrained state feedback
translate the pinned coss component system into native MyGo. This records the
existing implementation; it does not establish a new visual identity.

The source of executable values is `theme/`, with component-specific geometry
in the corresponding Go files. Frontmatter records their default values in
portable CSS notation: one `px` represents one unscaled native DIP, not a
physical screen pixel. Native text scaling and configured preferences still
apply. Refresh this record when those defaults change.

The pinned reference and licensing boundary live in `docs/references.md`.
`docs/design.md` defines ownership and native behavior; `docs/components.md`
defines supported component subsets. This file describes visual treatment.

## Colors

Use semantic roles through `theme.Colors(c)` rather than local palettes.
Primary identifies actions; Accent identifies subtle interaction surfaces.
MutedBackground and MutedForeground remain distinct. Error, warning, success
and informational colors express state rather than decoration.

The frontmatter contains directly represented defaults from `theme/palette.go`.
Dark background/card/popover and muted text use Go's compositing/mix functions;
their exact result belongs to that implementation and `docs/theme.md`.
Do not substitute CSS color-mix or rounded hex approximations for them.
Preserve translucent borders and fills until they composite on the real surface.
Host overrides are separate for each appearance.

## Typography

The library defaults to native sans and mono fallbacks. Heading, body and
mono family preferences are independent; heading falls back to body.
Gallery startup explicitly registers Cal Sans. Configuring a family name
does not load its font, and library imports must not register fonts.

Frontmatter records the base, label, help, H1 and code roles from
`theme/tokens.go`; that file also defines H2-H6 and smaller body variants.
Semantic Typography has its own prose line heights through `SemanticType.Style`.
Use those APIs, `theme.Scale` and OS text scaling rather than treating the
portable token record as a complete layout engine.

## Layout

Use the existing four-DIP rhythm and content-owned composition. Set alignment
explicitly; MyGo row defaults and zero-valued alignment differ from CSS.
Controls keep native keyboard, selection and focus behavior.

The gallery owns its application layout: a full-height sidebar whose top row
shares the native title bar, a page title bar of the same height, and a
left-aligned content column where title, sections and examples share one
edge, with a page index beside it on wide windows. The landing page drops the
sidebar for a centered column. Transient controls (appearance, compact navigation) open as overlays,
never as inserted rows. Its scaled compact threshold, content cap
and gutters come from `internal/gallery/docs_layout.go`; its surface brief
records the workflow. Enlarged text and short windows must leave actions
reachable. A web breakpoint does not redefine native control sizing.

## Elevation & Depth

Use surface tone, thin borders and existing native shadow treatments.
Buttons change fill, border and shadow without shrinking on press.
Dialogs use the existing dim/blur backdrop. Popups, menus and dialogs share
`theme.OverlayShadows` (coss shadow-lg/5 over a 1px edge, dark lower edge in
light and light upper edge in dark); tooltips use `theme.TooltipShadows`
(shadow-md/5). Anchored popups share one popover surface. Overlays retain
inert exit paint rather than rebuilding consumer content after closure.

Popovers, date and color pickers fade in from 98% around the edge facing
their anchor and fade out; tooltips also shrink on exit. Menus, selects and
combo boxes open and close at once, as in coss.

Motion follows the existing Go durations/easings and honors
`theme.ReduceMotion(c)`. `.impeccable/design.json` records motion extensions.

### State changes move

A control that changes state shows the change, wherever a change of state
is what the user just caused or needs to notice. `internal/motion` holds the
timing so every control moves alike.

| Change | Motion |
| --- | --- |
| A fill, border or ink changing with the state (checked, selected, pressed on, invalid, status) | Colour fade, 150 ms ease-out |
| A mark appearing (check, radio dot, dash, star fill, step tick) | Draws or scales in over 150 ms; leaves by fading, never by drawing backwards |
| Content opening or closing (accordion, disclosure, spoiler, expander row) | Height and fade, 200 ms ease-in-out |
| A tree branch opening | Its rows fade in over 200 ms while the rows below slide down; closing drops them at once and turns the arrow |
| An indicator moving to another item (tabs, segmented control, pager) | Slide, 200 ms ease-in-out |
| A quantity changing (progress, meter) | Eases to the new value over 500 ms ease-out, as coss does |
| A slider fill moving to a value the application set | Eases over 200 ms ease-out; a held pointer is followed at once |
| A number changing (`NumberFlip`, so Statistic) | Each digit that differs rolls over 300 ms ease-out, up when the value rose and down when it fell; other characters and reduced motion crossfade |
| An icon standing for the state swapping (copy to check, eye to eye-off) | Crossfade, 150 ms |

Left at once on purpose: hover and focus rings (they track the pointer and
the keyboard), menu, select and combo box popups and their row highlight
(coss), text the user is typing, and anything the user drags. A value set
by the application on the first frame does not animate in. Reduced motion
keeps the fades and settles sizes and positions at once. Nothing requests
frames at rest.

A state never changes font weight. Selected, active, hover and current text
keeps its weight: heavier text is wider and shifts the layout around it. A
state shows in fill, ink, an indicator or an icon.

### Loading is one state

A component that can wait on the application takes a `Loading` flag the
caller owns; it starts no work. `internal/busy` holds the treatment.

| Kind | While loading |
| --- | --- |
| Action controls (Button, ToggleButton, SplitButton, Switch) | Keeps its size; content crossfades to a centred spinner at once; takes no press but does not look disabled |
| Fields (Input, SearchField, Select, ComboBox, Autocomplete) | A spinner takes the trailing affordance's box; the field stays editable; an empty popup shows one Loading row |
| Regions of data, first load (Table, ListBox, Tree, Command, Descriptions, Statistic, SettingsGroup) | Skeletons in the content's shape inside the real chrome, one shimmer per region, no spinner |
| Regions of data, refresh | Content stays at 55%, takes no input, under a 2 DIP indeterminate line on the top edge; no scrim |

Region indicators wait 150 ms, so a fast load shows nothing, and fade out
over 150 ms without waiting. The region is busy to assistive technology. `LoadingOverlay` covers
composed content the same way when a component has no flag of its own.
Headless rounded clips do not prove arbitrary GPU clipping correctness.

## Shapes

Radii derive from the configurable base; the frontmatter shows defaults at
base ten. Fields use the MD radius; cards use the larger surface treatment.
Joined groups receive explicit segment corners. The default ToggleGroup
keeps separated rounded items; joined corners belong to its outline form.
Use per-corner native painting where a component owns only part of a boundary.

## Components

Buttons and fields use the desktop size scale. The primary button's default
fill/text pair is shown in frontmatter; its height is scaled by the native
theme. Disabled controls keep MyGo's existing opacity and remain inert.
Read-only editors retain selection.

Cards and frames compose header, body and footer spacing explicitly.
Overlay footers own their muted full-width face and top rule. Native system
context menus remain the default; the themed in-window menu is opt-in.
Icon-only actions still have accessible names, and field help/errors retain
their native associations.

Apply the same semantic vocabulary to navigation, forms, feedback, dates and
tables. Component-specific states and native adaptations are recorded in
`docs/components.md`; a visual resemblance is not a parity guarantee.

## Do's and Don'ts

- Do preserve the pinned coss visual authority and native MyGo semantics.
- Do reuse theme roles and coherent public components before extracting more.
- Do inspect relevant light/dark, compact/wide and stateful native output.
- Do keep UI copy concise and reveal secondary help only when needed.
- Don't replace native editing/focus with pointer-only drawings.
- Don't invent fonts, palettes, mobile targets or new product claims.
- Don't treat HTML/CSS detector output as evidence about Go controls.
- Don't claim GPU, OS gesture or screen-reader success from headless images.
