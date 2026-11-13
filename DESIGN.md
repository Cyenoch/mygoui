---
name: mygoui
description: Compact neutral native components and a MyGo workbench.
colors:
  light-background: "#ffffff"
  light-foreground: "oklch(26.9% 0 0)"
  light-primary: "oklch(26.9% 0 0)"
  light-primary-foreground: "oklch(98.5% 0 0)"
  light-accent: "rgb(0 0 0 / 4%)"
  light-border: "rgb(0 0 0 / 8%)"
  light-input: "rgb(0 0 0 / 10%)"
  light-ring: "oklch(70.8% 0 0)"
  light-sidebar: "oklch(98.5% 0 0)"
  light-sidebar-border: "rgb(0 0 0 / 6%)"
  light-destructive: "oklch(63.7% 0.237 25.331)"
  light-destructive-foreground: "oklch(50.5% 0.213 27.518)"
  light-info: "oklch(62.3% 0.214 259.815)"
  light-info-foreground: "oklch(48.8% 0.243 264.376)"
  light-success: "oklch(69.6% 0.17 162.48)"
  light-success-foreground: "oklch(50.8% 0.118 165.612)"
  light-warning: "oklch(76.9% 0.188 70.08)"
  light-warning-foreground: "oklch(55.5% 0.163 48.998)"
  light-chart-1: "oklch(64.6% 0.222 41.116)"
  light-chart-2: "oklch(60% 0.118 184.704)"
  light-chart-3: "oklch(39.8% 0.07 227.392)"
  light-chart-4: "oklch(82.8% 0.189 84.429)"
  light-chart-5: "oklch(76.9% 0.188 70.08)"
  dark-foreground: "oklch(97% 0 0)"
  dark-primary: "oklch(97% 0 0)"
  dark-primary-foreground: "oklch(26.9% 0 0)"
  dark-accent: "rgb(255 255 255 / 4%)"
  dark-border: "rgb(255 255 255 / 6%)"
  dark-input: "rgb(255 255 255 / 8%)"
  dark-ring: "oklch(55.6% 0 0)"
  dark-sidebar-border: "rgb(255 255 255 / 5%)"
  dark-destructive-foreground: "oklch(70.4% 0.191 22.216)"
  dark-info-foreground: "oklch(70.7% 0.165 254.624)"
  dark-success-foreground: "oklch(76.5% 0.177 163.223)"
  dark-warning-foreground: "oklch(82.8% 0.189 84.429)"
  dark-chart-1: "oklch(48.8% 0.243 264.376)"
  dark-chart-4: "oklch(62.7% 0.265 303.9)"
  dark-chart-5: "oklch(64.5% 0.246 16.439)"
typography:
  body:
    fontFamily: "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"
    fontSize: "16px"
    fontWeight: 400
    lineHeight: 1.5
  body-sm:
    fontSize: "14px"
    fontWeight: 400
    lineHeight: 1.428571
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
  heading-2:
    fontSize: "30px"
    fontWeight: 600
    lineHeight: 1.2
    letterSpacing: "-0.025em"
  heading-3:
    fontSize: "24px"
    fontWeight: 600
    lineHeight: 1.333333
    letterSpacing: "-0.025em"
  heading-4:
    fontSize: "20px"
    fontWeight: 600
    lineHeight: 1.4
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
  space-2: "2px"
  space-4: "4px"
  space-6: "6px"
  space-8: "8px"
  space-12: "12px"
  space-16: "16px"
  space-24: "24px"
  space-32: "32px"
  space-48: "48px"
components:
  button-primary-light:
    backgroundColor: "{colors.light-primary}"
    textColor: "{colors.light-primary-foreground}"
    typography: "{typography.label}"
    height: "32px"
    padding: "0 11px"
    rounded: "{rounded.lg}"
  button-primary-dark:
    backgroundColor: "{colors.dark-primary}"
    textColor: "{colors.dark-primary-foreground}"
    typography: "{typography.label}"
    height: "32px"
    padding: "0 11px"
    rounded: "{rounded.lg}"
  button-ghost-light:
    textColor: "{colors.light-foreground}"
    typography: "{typography.label}"
    height: "32px"
    padding: "0 11px"
    rounded: "{rounded.lg}"
  input-light:
    backgroundColor: "{colors.light-background}"
    textColor: "{colors.light-foreground}"
    typography: "{typography.body-sm}"
    height: "32px"
    rounded: "{rounded.md}"
  tooltip-light:
    backgroundColor: "{colors.light-background}"
    textColor: "{colors.light-foreground}"
    typography: "{typography.help}"
    padding: "4px 8px"
    rounded: "{rounded.md}"
  popover-light:
    backgroundColor: "{colors.light-background}"
    textColor: "{colors.light-foreground}"
    rounded: "{rounded.lg}"
  card-light:
    backgroundColor: "{colors.light-background}"
    textColor: "{colors.light-foreground}"
    padding: "24px"
    rounded: "{rounded.xxl}"
  dialog-light:
    backgroundColor: "{colors.light-background}"
    textColor: "{colors.light-foreground}"
    padding: "24px"
    width: "512px"
    rounded: "{rounded.xxl}"
---

# Design System: mygoui

## Overview

Compact neutral surfaces, clear text hierarchy and restrained state feedback
for native MyGo desktop applications. This file is the authority for how
mygoui looks and moves. A rule that is missing is added here through an
OpenSpec change. It was extracted from the implementation; it does not
establish a new identity.

Values are written so a reader needs nothing else. They are executed in Go:
colours, type, spacing, radii, sizes, shadows, motion and effect tokens in
`theme/` (one `theme.Theme` value, `theme.Default()` being this record), and
component geometry in each component's file. When Go and this record
differ, this record wins: the mismatch is a defect, and a change settles it
by correcting the code, or the record if the record was wrong.

One `px` in the frontmatter is one unscaled native DIP, not a physical
pixel. Text scaling (`theme.Scale`) and host preferences still apply.
`docs/design.md` owns engineering contracts, `docs/components.md` owns the
supported subset of each component, and `docs/references.md` records where
these rules came from and the licensing boundary.

**Key Characteristics:**

- Neutral ink on white or near-black; colour is reserved for state and data.
- Hairline translucent borders and very soft shadows instead of heavy chrome.
- Desktop density: 32 DIP controls, 14 DIP control text, a 4 DIP rhythm.
- Short, ease-out state motion; nothing moves at rest.

## Colors

A neutral ink-on-paper palette in OKLCH, with four status hues and five
chart hues. Components read roles through `theme.Colors(c)`; they never
keep local palettes.

### Primary

- **Ink** (`oklch(26.9% 0 0)` light, `oklch(97% 0 0)` dark): Primary
  action fill and body text. Its foreground is `oklch(98.5% 0 0)` in light
  appearance and ink in dark. A host accent (`Prefs.Accent`) replaces
  Primary only; its foreground becomes `rgb(38 38 38)` on a light accent
  and white on a dark one.

Primary identifies actions. Accent identifies subtle interaction surfaces
(hover, selection). MutedBackground and MutedForeground stay distinct roles:
one is a fill, the other is text.

### Neutral

- **Background**: white in light appearance. In dark appearance
  `oklch(14.5% 0 0)` mixed 4% towards white.
- **Card / Surface**: white in light; dark background mixed a further 2%
  towards white. **Popover**: white in light; dark background mixed 4%
  towards white. Code blocks use the card value.
- **Accent, Secondary, Muted background**: black at 4% (light), white at 4%
  (dark). They stay translucent and composite on the surface below.
- **Muted foreground**: `oklch(55.6% 0 0)` mixed 10% towards black (light)
  or white (dark).
- **Border**: black 8% / white 6%. **Input**: black 10% / white 8%.
  **Ring**: `oklch(70.8% 0 0)` / `oklch(55.6% 0 0)`.
- **Sidebar**: `oklch(98.5% 0 0)` with a black 6% border (light);
  `oklch(14.5% 0 0)` mixed 3% towards white with a white 5% border (dark).
  Sidebar text is the foreground mixed 36% towards the sidebar.

### Status

Each status has a fill and a foreground for text on soft fills.

| Status | Fill | Foreground (light / dark) |
| --- | --- | --- |
| Destructive | `oklch(63.7% 0.237 25.331)`; dark mixes it 10% towards white | `oklch(50.5% 0.213 27.518)` / `oklch(70.4% 0.191 22.216)` |
| Info | `oklch(62.3% 0.214 259.815)` | `oklch(48.8% 0.243 264.376)` / `oklch(70.7% 0.165 254.624)` |
| Success | `oklch(69.6% 0.17 162.48)` | `oklch(50.8% 0.118 165.612)` / `oklch(76.5% 0.177 163.223)` |
| Warning | `oklch(76.9% 0.188 70.08)` | `oklch(55.5% 0.163 48.998)` / `oklch(82.8% 0.189 84.429)` |

Text on a solid destructive fill is white.

### Charts

Light: `oklch(64.6% 0.222 41.116)`, `oklch(60% 0.118 184.704)`,
`oklch(39.8% 0.07 227.392)`, `oklch(82.8% 0.189 84.429)`,
`oklch(76.9% 0.188 70.08)`. Dark: `oklch(48.8% 0.243 264.376)`, the success
fill, the warning fill, `oklch(62.7% 0.265 303.9)`,
`oklch(64.5% 0.246 16.439)`.

### Derived colours

Derived roles are computed, never picked:

- Mixing happens in premultiplied sRGB (`theme/palette.go` `mix`). Do not
  substitute CSS `color-mix` or rounded hex approximations.
- Hover of a solid fill is the fill at 90% of its alpha
  (`Effects.HoverOpacity`); pressed is the same.
- Soft status fills are the status fill at 8% alpha (light) or 16% (dark);
  their hover adds 4%.
- Secondary and tertiary backgrounds are MutedBackground and Input
  composited over the background. Secondary and tertiary borders multiply the border's
  alpha by 1.5 and 2.
- Text selection is Primary at 20%. The scrollbar is the foreground at 20%,
  6 DIP wide. A modal backdrop is black at 32%.
- A field fills with the background (light) or the input colour at 32% of
  its alpha (dark). A segmented track is the background (light) or Input
  over the background (dark).
- The native control face (MyGo `Surface`) is the popover colour; its hover
  and pressed faces composite Accent over it.

**The State Colour Rule.** Hue appears for state and data only:
destructive, info, success, warning and chart series. Navigation,
selection and emphasis use ink and neutral fills, unless a host sets an
accent for Primary.

**The Translucent Line Rule.** Borders, inputs and accent fills stay
translucent so they read the same on every surface they sit on. Host
overrides are set per appearance.

## Typography

Native sans and mono stacks; one family per role, independent heading,
body and mono preferences (heading falls back to body). Configuring a family
name does not load the font: the application registers fonts, and a
library import never does. The gallery registers Cal Sans at startup.

| Role | Size / line height (DIP) | Weight |
| --- | --- | --- |
| XS | 12 / 16 | 400 |
| SM (control text, body small) | 14 / 20 | 400 |
| Base (body) | 16 / 24 | 400 |
| LG | 18 / 28 | 400 |
| XL | 20 / 28 | 400 |
| H1 / H2 / H3 | 36 / 40, 30 / 36, 24 / 32 | 600, tracking −0.025 em |
| H4 / H5 / H6 | 20 / 28, 18 / 28, 16 / 24 | 600, tracking −0.025 em |
| Label | 14 / 20 | 500 |
| Help text | 12 / 16 | 400 |
| Code | 14 / 20 mono | 400 |
| Card title | 18 / 20 | 600 |
| Dialog title (H4) | 20 / 28 | 600 |
| Prose body / small / extra small | 16 / 28, 14 / 24, 12 / 20 | 400 |

Weights in use: 400, 500, 600, 700. Every size passes through
`theme.Scale` and the OS text scale.

**The Steady Weight Rule.** A state never changes font weight. Selected,
active, hover and current text keeps its weight, because heavier text is
wider and shifts the layout. A state shows in fill, ink, an indicator or an
icon.

## Layout

A 4 DIP rhythm (`Spacing.Unit`); components space by steps of it,
`theme.Space(n)`, and the built-in steps are named by their DIP value
(`Space2` … `Space64`). The values in use, most frequent first: 8,
12, 4, 16, 24, 6, 2, 20, 48, 32, 40, 10, 64, 14.

- Inside a control: 8 DIP between icon and label.
- Between related controls and rows: 8–12 DIP. Between sections: 16–24.
- Cards and dialogs pad their content by 24 DIP; a dialog footer pads by
  16 × 24 DIP.
- Alignment is explicit; MyGo row defaults and zero-valued alignment differ
  from CSS. Controls keep native keyboard, selection and focus behavior.

Windows narrower than 640 DIP change dialogs: a centred dialog becomes full
width with square corners and a top hairline, an inset sheet loses its
inset, and footer actions stack in reverse. Enlarged text and short
windows leave every action reachable. A web breakpoint never redefines
native control sizing.

The gallery owns its application layout: a full-height sidebar whose top row
shares the native title bar, a page title bar of the same height, and a
left-aligned content column where title, sections and examples share one
edge, with a page index beside it on wide windows. The landing page drops the
sidebar for a centered column. Transient controls (appearance, compact
navigation) open as overlays, never as inserted rows. Its compact threshold
(760 DIP, scaled), content cap and gutters come from
`internal/gallery/docs_layout.go`; its surface brief records the workflow.

**Window layout.** A custom title bar is at least the native band tall and
never draws window buttons. It keeps 12 DIP gutters beyond the window-control
insets, 8 DIP gaps, small (28 DIP) ghost controls and a hairline underneath;
when its controls do not fit between the window controls they move below the
band, at once. Side panes default to 256 DIP, at least 160, and leave the
content at least 240. Dividers are hairlines. A full-height start pane uses
the Sidebar surface and shares the title band, with no hairline under its
band. A fixed pane folds its width over 200 ms ease-out; reduced motion
settles at once. Below 760 DIP (scaled with text) panes leave the row and
show as a flush sheet from their edge, at the pane's width. The gallery shell
is this layout: a 52 DIP band and a 260 DIP fixed sidebar.

**The Density Rule.** Default controls are 32 DIP tall with 14 DIP text.
Larger sizes exist for emphasis, not for touch.

## Elevation & Depth

Depth is tonal first: a card is a step lighter than the background in dark
appearance, separated by a hairline border. Shadows are soft (5% black)
and paired with a 1 DIP edge; they signal that a surface floats, not
importance.

### Shadow Vocabulary

The 1 DIP edge is an inset line inside the surface, lit from above: black
along the bottom inside edge in light appearance, white along the top inside
edge in dark.

- **Surface** (light: `0 1px 2px` black 5% plus an inset `0 -1px` black 4%;
  dark: an inset `0 1px` white 6%): cards and resting surfaces.
- **Field** (light: `0 1px 2px` black 5% plus an inset `0 -1px` black 6%;
  dark: an inset `0 1px` white 6%): outline buttons. Fields at rest carry
  only `0 1px 2px` black 5%, dropped when focused, invalid or disabled.
- **Overlay** (`0 10px 15px -3px` and `0 4px 6px -4px`, black 5% each,
  plus the surface's inset edge: black 4% at the bottom in light, white 6%
  at the top in dark): popovers, menus, pickers and dialogs.
- **Tooltip** (`0 4px 6px -1px` and `0 2px 4px -2px`, black 5%, over the
  same edge): tooltips.
- **Solid button** (`0 1px 2px` of the fill at 24% plus an inset white 16%
  top edge): primary and destructive buttons; removed while pressed.

### Overlays

Anchored popups (popover, menu, select, combo box, pickers, teaching tip,
tour) share one surface: the popover fill, a hairline border, the overlay
shadow, the LG radius. Tooltips use the popover fill, the MD radius, XS
text, 4 × 8 DIP padding, at most 320 DIP wide.

Modal surfaces sit over a backdrop of black 32% with a 4 DIP blur (`Effects.BackdropBlur`) and use
the XXL radius where they float: a centred dialog (square and full width
below 640 DIP), an inset sheet, a bottom sheet and a floating panel. A
sheet anchored to a window edge is square unless inset. A drawer rounds
only the corners facing into the window, and all four when inset. A
dialog is 512 DIP wide by default and a side drawer or sheet 448; a bottom
sheet is capped at 85% of the window. Overlay footers own a full-width
face of MutedBackground at 72% with a top hairline. Dialog titles use H4. Overlays keep inert
exit paint instead of rebuilding consumer content after closing.

### Overlay motion

Tokens: anchored surfaces and the floating panel take `Motion.Popup`
(150 ms). Each kind of modal surface moves its scrim with it: dialogs and
sheets take `Motion.Dialog` (200 ms) on `Motion.EaseTravel`; drawers, bottom
and floating sheets `Motion.Drawer` (450 ms) on `Motion.EaseDrawer`. Ease-out is
`Motion.Ease`, ease-in-out is `Motion.EaseTravel`, and the 98% scale is
`Effects.PopupScale`. The tour reads `Motion.Reveal` and `Motion.Roll`; the
toast spring and its content fade are constants of the toast.

| Surface | Enter | Exit |
| --- | --- | --- |
| Popover, date and colour pickers, teaching tip | Fades in from 98% scale, origin at the edge facing the anchor; 150 ms ease-out | Fades out at full size; 150 ms |
| Tooltip, preview card | As a popover | Fades while shrinking to 98% |
| Menu, select, combo box | At once | At once |
| Dialog (centred) | Fades and scales from 98% on windows ≥ 640 DIP; rises 16 DIP on narrower windows; 200 ms ease-in-out | Reverse |
| Dialog anchored to an edge | Slides 32 DIP and fades; 200 ms ease-in-out | Reverse |
| Modal, ImageViewer, side dialog without slide | Fade, 200 ms ease-in-out; the backdrop fades with the surface | Fade |
| Dialog and AlertDialog scrim | Fades with its dialog, 200 ms ease-in-out | Fades |
| Drawer, bottom sheet, floating sheet | Slides from its edge over 450 ms on cubic-bezier(0.32, 0.72, 0, 1); a dragged sheet settles on the same curve | Reverse |
| Drawer scrim, bottom sheet backdrop | Fades over 450 ms on the drawer curve | Fades |
| ActionBar | Rises 8 DIP over 200 ms ease-out | Reverse |
| Tour | Scrim fades 200 ms ease-out; the highlight travels 300 ms ease-in-out while steps crossfade 300 ms linear; the panel height changes over 200 ms ease-in-out | Scrim fades |
| ImageViewer stage | Pan and zoom settle over 200 ms ease-in-out; toolbar colours fade 150 ms ease-out | — |
| Floating panel | Fade, 150 ms ease-out | Fade |
| Toast deck | Position and size follow a 320 ms spring; content fades over 200 ms | Spring |

**The Still Overlay Rule.** Under reduced motion every overlay, popup,
dialog, sheet, drawer, floating panel and scrim appears and disappears at
once, fades included.

**The Quiet List Rule.** Menus, selects and combo boxes open, close and
move their highlight at once: a list the user is scanning never animates.

**The Hairline Edge Rule.** Every floating surface carries a 1 DIP edge
with its shadow, so it stays legible when the shadow is invisible on a
dark or busy backdrop.

## Shapes

Radii derive from a configurable base (default 10 DIP): XS base−6 (4),
SM base−4 (6), MD base−2 (8), LG base (10), XL base+4 (14), XXL base+6 (16),
XXXL base+14 (24), XXXXL base+22 (32). Fields use base−2 (8).

| Radius | Used for |
| --- | --- |
| SM | Menu items |
| MD | Fields, tooltips, extra-small buttons |
| LG | Buttons, popovers, menus, selects, pickers, teaching tips, tours |
| XL | Alert, Frame, SettingsGroup card, Tag, Upload, close button and link focus area, action bar, image viewer toolbar, pressable feedback |
| XXL | Cards, dialogs, sheets, floating panels |

Joined groups receive explicit segment corners. The default ToggleGroup
keeps separated rounded items; joined corners belong to its outline form.
Use per-corner native painting where a component owns only part of a
boundary. Borders are 1 DIP hairlines (`theme.Hairline`).

## Components

### Sizes

| Control | Heights (DIP) |
| --- | --- |
| Button xs / sm / md / lg / xl | 24 / 28 / 32 / 36 / 40; horizontal padding 7 / 9 / 11 / 13 / 15; icon-only is square |
| Input | 32; textarea minimum 38 |
| Chip | 20 / 24 / 32 |
| Badge | 16 / 18 / 22 |
| Avatar | 24 / 32 / 40 |
| Kbd | 20 |
| List row | 40 |
| Icons | 12 / 14 / 16 / 20 |
| Status dot | 8 |
| Spinner | 12 / 16 / 24 / 32 |

### Buttons

Text is SM medium (XS for extra small, Base for extra large). Icons are
16 DIP at 80% opacity with an 8 DIP gap. Faces by variant:

- **Primary, destructive**: solid fill and the solid button shadow; hover
  and pressed at 90% alpha.
- **Secondary**: Secondary fill; hover 90%, pressed 80%.
- **Outline**: popover face in light appearance and the input colour at
  32% in dark, with the field shadow; hover and pressed move to Accent at
  50% (light) or input at 64% (dark).
- **Ghost**: transparent; hover and pressed use Accent.
- **Destructive outline**: hover fills destructive at 4%.

Buttons do not shrink when pressed; the shadow is removed instead.
`PressableFeedback`, for custom pressable surfaces, scales to 98.5% over
250 ms ease-out and shows its highlight at once.

### Fields

A field fills with the field background, a hairline Input border, the MD
field radius and SM text. Hover changes nothing. Focus turns the border to
Ring and adds a 3 DIP ring outside at Ring 24%. Invalid turns the border to
destructive at 36%; invalid and focused, the border is destructive at 64%
and the ring destructive at 16% (24% in dark). The invalid border fades
over 150 ms; the focus border changes at once. Read-only editors keep
selection. Field help and errors keep their native associations.

### Focus and disabled

Keyboard focus shows MyGo's native ring: 2 DIP of Ring, 1 DIP outside the
control. Joined segments draw the same ring 1 DIP inside instead, so it is
not clipped by a neighbour. Disabled controls are drawn at MyGo's 50%
opacity and are inert; loading controls are not disabled.

### Cards and surfaces

The same semantic vocabulary (roles, radii, sizes, motion) applies across
navigation, forms, feedback, dates and tables. A visual resemblance to
another component system is not a parity guarantee; `docs/components.md`
records each supported subset.

Cards use the Surface fill, a hairline border, the surface shadow and the
XXL radius; header, content and footer pad by 24 DIP with 6 DIP between
title and description. Native system context menus are the default; the
themed in-window menu is opt-in. Icon-only actions keep accessible names.
Component-specific states and native adaptations are recorded in
`docs/components.md`.

### State motion

A control that changes state shows the change, wherever the change is what
the user just caused or needs to notice. `theme.Motion()` holds the shared
timing: `State` 150 ms for a change in place, `Reveal` 200 ms for content
opening or a mark sliding, `Roll` 300 ms for a value rolling, `Shift`
150 ms for neighbours making room, `Page` 250 ms for a page changing.
`Ease` is ease-out and `EaseTravel` is ease-in-out. A duration below that
is none of these belongs to its component (`docs/theme.md` lists them).

| Change | Motion |
| --- | --- |
| A fill, border or ink changing with the state (checked, selected, pressed on, invalid, status) | Colour fade, 150 ms ease-out |
| A mark appearing (check, radio dot, dash, star fill, step tick) | Draws or scales in over 150 ms ease-out; leaves by fading, never by drawing backwards |
| Content opening or closing (accordion, disclosure, settings expander row, timeline, about sections) | Height and fade, 200 ms ease-out |
| A spoiler opening or closing | Height, 200 ms ease-out; its cut fade clears as it opens |
| A tree branch opening | Its rows fade in over 200 ms while the rows below slide down; the arrow turns over 200 ms ease-out; closing drops the rows at once |
| An indicator moving to another item (tabs, segmented control, tour dots) | Slide, 200 ms ease-in-out |
| A pager's current pip | Grows over 150 ms ease-out |
| A quantity changing (progress, meter) | Eases to the new value over 500 ms ease-out |
| A slider fill moving to a value the application set | 200 ms ease-out; a held pointer is followed at once |
| A number changing (`NumberFlip`, so Statistic) | Each digit that differs rolls over 300 ms ease-out, up when the value rose and down when it fell; other characters and reduced motion crossfade |
| An icon standing for the state swapping (copy to check, eye to eye-off) | Crossfade, 150 ms |
| A switch toggling | Thumb slides and track colour fades together, 150 ms ease-out; a held thumb widens 1.2× (`Effects.ThumbPressScale`) over 150 ms without delay, as a slider thumb grows |
| A slider or colour-area thumb held | Grows 1.2× (`Effects.ThumbPressScale`) over 150 ms |
| A page changing (carousel, navigation stack) | Slides over 250 ms ease-out |
| A row or tab moving aside during a drag | 150 ms ease-out |
| A sidebar or resizable panel collapsing | Width, 200 ms ease-out |
| A `layout` pane folding | Width, 200 ms ease-out; at once under reduced motion |
| A document tab opening, closing or reflowing | Tabs and the strip move and resize over 200 ms ease-out; a tab fades over 150 ms; a selected tab scrolls into view over 200 ms ease-out; wheel scrolling glides 150 ms ease-out |
| A carousel released short of the next page | Pulls back over 200 ms ease-out |
| A calendar view swapping | The year grid fades in, the day grid fades in and out, the header arrows fade and the chevron turns, all over 150 ms ease-out (`Motion.State`) |
| A banner dismissed | Leaves over 200 ms ease-in-out; its colours change over 150 ms ease-out |
| A link hovered | Its trailing icon fades from 60% to full over 150 ms ease-out |

Left at once on purpose: hover and focus rings (they track the pointer and
the keyboard), menu, select and combo box popups and their row highlight,
text the user is typing, and anything the user drags. Hover and press fills show
at once (close button, `PressableFeedback`, a resizable panel's handle under
the pointer; its keyboard focus fades over 150 ms). Current exceptions that
animate hover: the link icon (150 ms). A palette swatch grows 1.1× on hover
at once, since a tint would change the colour it shows. A value set by the
application on the first frame does not animate in. Nothing requests
frames at rest.

**The Reduced Motion Rule.** Under `theme.ReduceMotion(c)` control state
motion keeps its fades while sizes, positions, scales and loops settle at
once. Overlays keep nothing: see the Still Overlay Rule.

### Loading

A component that can wait on the application takes a `Loading` flag the
caller owns; it starts no work. `internal/busy` holds the treatment.

| Kind | While loading |
| --- | --- |
| Action controls (Button, ToggleButton, SplitButton, Switch) | Keeps its size; content crossfades to a centred spinner over 150 ms with no wait; takes no press but does not look disabled |
| Fields (Input, SearchField, Select, ComboBox, Autocomplete) | A spinner takes the trailing affordance's box; the field stays editable; an empty popup shows one Loading row |
| Regions of data, first load (Table, ListBox, Tree, Command, Descriptions, Statistic, SettingsGroup, charts) | Skeletons in the content's shape inside the real chrome, one shimmer per region, no spinner |
| Regions of data, refresh | Content stays at 55% opacity, takes no input, under a 2 DIP indeterminate line on the top edge; no scrim |

Region indicators wait 150 ms, so a fast load shows nothing, and fade out
over 150 ms without waiting. The region is busy to assistive technology.
`LoadingOverlay` covers composed content the same way when a component has
no flag of its own. Tokens: the wait is `Effects.LoadingDelay`, the fade
`Motion.State`, the refresh opacity `Effects.RefreshOpacity`, the spinner
turn `Effects.SpinnerPeriod`, and the sweep `Effects.SweepPeriod` and
`Effects.SweepBand`.

- **Spinner**: one turn per second. A spinner painted inside a field
  draws a 1.8 rad arc in twelve steps, stroke size × 2.5 / 24, and rests
  still under reduced motion.
- **Indeterminate sweep**: a band 40% of the track enters from wholly
  before it and leaves wholly past it every 1.5 s ease-in-out, travelling
  at least 4.5 bands, at most 60 frames per second; a progress bar and a
  region's refresh line share it. Under reduced motion the
  band holds still at the start of the track.
- **Progress value**: 500 ms ease-out.

**The One Loading State Rule.** A component shows waiting one way for its
kind (table above); it never stacks a spinner, a skeleton and a scrim.

### Effects

Shader effects run on an effect layer limited to 8192 pixels a side, scale
included, and paint a designed fallback where no GPU is available or the
layer is too large. Animated effects rest still under reduced motion;
skeletons and shimmer also rest under high contrast. An effect makes one
state legible and is never its only signal: the label, role or text says it
too. An element carries one effect.

The public family (`effect`, [docs/effects.md](docs/effects.md)) has three:

| Effect | Use it when | Tokens | Fallback |
| --- | --- | --- | --- |
| Mask | Content continues past an edge and is cut there over a background that is not one flat colour | `Effects.FadeEdge` | Content whole; the edge shaded with the foreground at 18%, or blended into a named opaque surface |
| Shimmer | Placeholders in an application's own layout wait as one region | `Effects.ShimmerPeriod`, `ShimmerPeak`, `ShimmerPeakDark` | The placeholders rest |
| Reveal | Content arrives after the user asked for it; never what a window opens with | `Motion.Reveal`, `Motion.Ease`, `Effects.FadeEdge` | A uniform fade; at once above the layer limit |

Mask is static. Shimmer loops only while applied. Reveal plays once and
leaves no layer; under reduced motion it is a fade with no travelling front.

Not in the language: a glow or spotlight shader (theme shadows and the
tour's path fill do it), grain or noise, a shader material (the backdrop is
MyGo's glass blur), an animated or gradient border (it would compete with
the focus ring), a press ripple, a success sheen, a noise dissolve.

- **Skeleton**: a 2 s cycle (`Effects.ShimmerPeriod`). Shimmer is a 120°
  band (stops 40 / 50 / 60%) peaking at white 64% in light and 4% in dark
  (`Effects.ShimmerPeak`, `ShimmerPeakDark`). Pulse lowers a
  skeleton's alpha to half and back, easing in and out, the same in a
  single skeleton and in `SkeletonGroup`, which runs one shader (the
  family's Shimmer) across its rows. Without
  a GPU the group paints still. Rows in a group are 8 DIP apart.
- **Scroll fade**: an edge that has more to scroll fades over 24 DIP
  (`Effects.FadeEdge`), in
  proportion to the distance left. Three modes: the default shadow (the
  foreground at 18%), a surface fade that blends towards a given opaque
  surface colour, and an opt-in alpha mask on the content (a shader; it
  falls back to the shadow without a GPU or above the layer limit). `ScrollShadow` fades an opaque
  surface over the same 24 DIP. The document tab strip masks 24 DIP; a
  closed spoiler fades its cut over 24 DIP.
- **Status dot**: pulse grows a ring to twice the dot over 2 s while its
  alpha falls from 40% to 0; breathe varies the dot between 65% and 100%
  over 1.6 s. Both redraw at 30 frames per second.
- **Caret**: the OTP caret blinks on a 500 ms cycle, visible for the first
  half.

## Do's and Don'ts

- Do read roles through `theme.Colors(c)` and sizes through `theme/`; keep this record and the Go values in step in the same change.
- Do reuse theme roles and coherent public components before extracting more.
- Do keep borders, inputs and accent fills translucent, and composite them on the real surface.
- Do pair every floating surface's shadow with its 1 DIP edge.
- Do show a state in fill, ink, an indicator or an icon, with the motion above.
- Do inspect light and dark, compact and wide, and stateful native output.
- Do keep UI copy concise and reveal secondary help only when needed.
- Do add a missing rule to this file through a change instead of borrowing one from another system.
- Don't change font weight with state.
- Don't animate menus, list highlights, typing, drags or anything at rest.
- Don't substitute CSS colour mixing or rounded hex values for derived colours.
- Don't replace native editing or focus with pointer-only drawings.
- Don't invent fonts, palettes, mobile targets or new product claims.
- Don't treat HTML/CSS detector output as evidence about Go controls.
- Don't claim GPU, OS gesture or screen-reader success from headless images.
