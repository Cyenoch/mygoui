# Design

## Context

The library stops at components. `internal/gallery` holds the only
application shell: a custom title bar reading `c.TitleBar()`, a full-height
sidebar sharing the title band, and a compact navigation sheet
(`gallery.go`, `docs_layout.go`). `mygoui.ResizablePanel` already gives a
splitter with pointer and keyboard resizing, clamping and collapse;
`overlay.Sheet` gives a modal edge panel; `navigation.StatusBar`, `Toolbar`,
`Menubar`, `NavigationStack` and `DocumentTabs` fill a shell but do not form
one. MyGo supplies `Context.TitleBar()` (insets in zoomed DIPs, zero in full
screen and in windows with a system title bar), `Element.DragWindow()` and
`Tester.SetTitleBar`.

Surface brief: `.impeccable/surfaces/internal-gallery.md` (Operate). The
design language is root `DESIGN.md`; no new visual direction is chosen here.

## Survey

| Toolkit | Title bar | Panes | Narrow window |
| --- | --- | --- | --- |
| AppKit | Unified title bar and toolbar; `fullSizeContentView` lets a sidebar run under the traffic lights; toolbar items sit beside the controls | `NSSplitViewController`: sidebar, content, inspector items with minimum and maximum thickness, `canCollapse`, user-resizable dividers | The sidebar collapses on its own and is shown over the content when toggled |
| WinUI | `TitleBar` (icon, title, subtitle, content, right header) with explicit drag regions and caption-button insets; grows to 48 with content | `NavigationView` (fixed-width pane, left or top) and a two-pane view; no user resizing | `NavigationView` turns Left → LeftCompact → LeftMinimal: the pane becomes an overlay opened by a button |
| libadwaita | `HeaderBar` (start, title, end) inside `ToolbarView` with top and bottom bars | `OverlaySplitView` (sidebar beside or over content, `collapsed` and `show-sidebar` set by the application, usually from a breakpoint) and `NavigationSplitView` (sidebar | content turning into a navigation stack); neither is user-resizable | `collapsed` set by a breakpoint; the sidebar overlays or the view becomes a stack |

What they share, and what an application needs from this library:

1. A band with start, middle and end content that keeps clear of the system
   controls and moves the window.
2. Up to three columns: a navigation or list pane, the content, an inspector.
   Each side pane has a width (sometimes user-resizable), bounds, and can be
   hidden.
3. One adaptation for a narrow window: a side pane leaves the row and is
   shown over the content on demand.
4. A shell that puts these together, with a strip at the bottom.

## Goals / Non-Goals

**Goals:** the four items above, as the smallest set of builders; the
gallery as their first consumer without losing behavior; caller-owned state;
native focus, keys and accessibility; reduced motion; right-to-left.

**Non-Goals:** `NavigationSplitView` (list | detail collapsing into a stack):
`navigation.NavigationStack` exists and an application chooses between it
and `SplitView` with `layout.Compact`. A scaffold variant with the title bar
above the sidebar: a Column of `TitleBar`, `SplitView` and a status bar is
that layout already. Vertical splits: `ResizablePanel` with `PanelTop` or
`PanelBottom`. Persisting sizes. Menu bars, toolbars and tabs in the shell:
they are content of the slots.

## Brief (Impeccable shape)

- **Job and audience.** A Go application author wants the standard desktop
  shell in a few lines and expects it to behave as the platform's does.
  Operate mode: precision and known patterns, no expression.
- **Outcome.** The window's chrome never overlaps the system controls at any
  zoom, text scale or width; side panes resize, hide and come back; a narrow
  window keeps every pane reachable.
- **Direction.** The gallery's existing shell, generalized. Neutral
  surfaces, hairline dividers, the sidebar fill for the start pane, ghost
  small buttons in the bar. Nothing moves at rest.
- **States.** Bar: in band, stacked below the band, no insets (full screen,
  browser). Pane: docked, resizing, at a bound, collapsed, fixed, shown as a
  sheet, sheet closing. Scaffold: with and without each pane and the strip.
- **Constraints.** Existing tokens, `theme.Motion()` among them; no new duration
  literal; the toolkit owns the window move and maximize.

## Decisions

### D1. A new `layout` package

The compact pane is a modal sheet, and `overlay` owns modal sheets (scrim,
focus trap and return, Escape, slide motion). The root package cannot import
`overlay` (`overlay` imports the root), and re-implementing the sheet in the
root would need its own motion durations, which this change must not add. So
the family lives in `layout`, importing the root (`ResizablePanel`, `Button`)
and `overlay` (`Sheet`). Nothing imports `layout`; the dependency rule in
`docs/design.md` gains one line. `layout` does not import `navigation`: the
status strip and the bar's contents are caller callbacks.

Alternative: put `TitleBar` in the root and only `SplitView` in a family.
Rejected: the scaffold needs both, and three related builders in two packages
is harder to find than one family.

### D2. The set: TitleBar, SplitView, Scaffold, PaneToggle, Compact

`TitleBar`, `SplitView` and `Scaffold` are the three things the survey
found. `PaneToggle` exists because a pane has two flags (D5) and the button
that flips the right one is the same in every application; it also is the
keyboard path to a pane in compact mode. `Compact` exposes the rule the view
uses, so a caller's own content can follow it.

### D3. TitleBar: three slots, measured stacking

`Leading` and `Trailing` are fixed-width groups; `Center` fills the room
between and is the caller's to align, with `Title` as the common case (a
truncating label at the start, as the unified AppKit bar and WinUI place
it). A truly window-centred title would need an absolute layer that can
collide with either group; not built.

The gallery decided stacking from a hard-coded control count. A library
cannot know the count, so the bar reads its own and the leading and trailing
groups' laid-out widths from the previous frame at build time
(`LayoutBounds`, the pattern `navigation.DocumentTabs` uses), stores them,
and invalidates when one differs from what the bar was resolved with. Paint
callbacks write nothing: the bar's inner rows (never the returned element,
which stays the caller's to decorate) compare the width they were built
with to the one they are painted at and ask for another build.
Both groups never shrink, so their widths do not depend on the decision and
it cannot oscillate. When the reserved insets are zero, stacking gains no
room and is skipped.

`Span` says which window edges the bar reaches (`SpanWindow`, `SpanStart`,
`SpanEnd`), mapped to the physical insets by layout direction. The scaffold
needs it; so does anyone composing a full-height pane by hand.
`MinHeight` is the native band the window was created with
(`mygo.WindowOptions.TitleBarHeight`), in window DIPs: it matters where the
window reports no band (the browser build) and is divided by the zoom as the
gallery did. `Borderless` drops the bottom hairline for a pane's band.

### D4. SplitView is built on ResizablePanel and Sheet

A resizable pane is `mygoui.ResizablePanel` unchanged: its splitter role,
arrow, Home, End and Enter keys, drag clamping, drag-to-collapse and double
click are already tested. The view adds a default width, a dynamic maximum
that leaves `ContentMin` for the content (from the view's last laid-out
width), a group named by the pane's label around the body, and the pane's
label as the splitter's name, as the APG window splitter is named after its
primary pane.

`Fixed` panes exist because two of the three surveyed toolkits do not let
the user resize a sidebar, and the gallery's is fixed. A fixed pane has a
hairline and no splitter; it collapses over `theme.Motion().Reveal`.

Alternative: a new N-pane splitter with proportional sizes. Rejected: no
application in the survey needs more than start | content | end, and nesting
covers list | detail inside a sidebar layout.

### D5. Two flags per pane, three compact modes

Docked visibility (`Collapsed`) and sheet visibility (`Open`) are different
facts: narrowing the window must not open a sheet, and widening it must not
forget that the user had hidden the sidebar. The gallery already kept two
(`SidebarClosed`, `NavOpen`). They are pointers, like `ResizablePanel` and
`Sidebar` take, so existing application state plugs in. The view clears
`Open` when it is not compact, as the gallery did.

`CompactMode` is `CompactAuto` (window width under 760 DIP scaled by text,
the threshold `DESIGN.md` records for the gallery), `CompactNever` or
`CompactAlways`. libadwaita makes `collapsed` a property the application
sets; the two forced values give that control (the gallery's landing page
always navigates through the sheet) without a threshold option. A caller
with another threshold compares the width itself.

The rule reads the window width, not the view's: the sheet is a window-level
modal, dialogs already switch on window width, and a view-relative rule
would lag a frame.

### D6. The compact pane is a flush Sheet showing the same body

AppKit and libadwaita show the same sidebar over the content. `overlay.Sheet`
pads its body by 24 DIP and scrolls it, which suits a dialog's text, not a
pane that fills its height and scrolls itself. `SheetOptions.Flush` builds
the body alone in the panel. The sheet keeps its scrim, slide, focus trap
and return, and Escape. The pane's label names the panel. `Flush` is for a
left or right sheet without `Inset`. The edge motion of sheets and drawers
used a physical offset while the panel rests by layout direction; it now
mirrors right to left. A sheet's `Side`
already follows the layout direction, so the start pane is `SideLeft` and
the end pane `SideRight` in both directions.

Alternative: a second `Body` for the compact presentation. Rejected: two
builders for one pane, and the gallery's two navigation lists were a
consequence of the padded sheet, not a wish.

### D7. Scaffold: the full-height start pane

One layout: start pane from the top of the window to the strip, sharing the
title band (AppKit's full-size sidebar, the gallery's shell); title bar,
content and end pane beside it; strip across the bottom. It is
`SplitView{Start, Content: Column{TitleBar, SplitView{End, Content}}}` over
the strip, so it adds no pane logic. The start pane's band is a `TitleBar`
(`StartBar`) forced to `SpanStart`, borderless and to the content bar's
height. The content bar's height and stacking follow where the pane will
be: `SpanEnd` while it is docked and shown, `SpanWindow` otherwise. Its
start padding follows where the pane is in this build: it reserves
`max(0, start inset - shown pane width)`, so its controls slide with the
pane's edge and never pass under the window controls. A fixed pane knows
its shown width from its own fold; a resizable one reads
`ResizablePanelOptions.Shown`, added for this. Both bars read one
measurement, kept on the scaffold, so they agree in the same build.

The outer view reserves `ContentMin` plus the end pane's docked width, since
its content column holds the end pane too; the scaffold keeps a pane's size
itself when the caller passes none, so it knows that width. The band of a
start pane shown as a sheet also reserves the end controls by as much as the
sheet reaches them (a sheet as wide as a narrow window). The start pane is filled with the Sidebar
colour; other surfaces are the caller's.

### D8. Motion

Fixed-pane collapse uses `theme.Motion().Reveal` and `Ease`, as
`ResizablePanel` does. Divider
drags follow the pointer at once. The sheet's motion is `overlay.Sheet`'s.
Stacking and span changes are instant: they follow a window resize.
`AnimateTo` settles at once under reduced motion.

### D9. Accessibility and focus

- Bar: a group named by `Label`, else `Title`, else the locale's `TitleBar`.
  Controls are ordinary Tab stops in build order: leading, centre, trailing.
- Pane: a group named by `Label`; its divider a splitter with the same name,
  range and step, and Expanded when collapsible.
- Tab order of a scaffold: start band, start body, start divider, title bar,
  content, end divider, end body, strip. MyGo orders absolute children after
  those in the flow, so `ResizablePanel` now builds a start- or top-edge
  handle as a sticky child before the content: first for the focus, above
  the content for painting and the pointer, taking no room.
- Content keeps its identity across modes (keyed), so focus and widget state
  survive a resize. A focused control inside a pane that leaves loses focus,
  as any removed control does.
- `PaneToggle`: toggle button, checked while the pane shows, named by
  `Label` or the locale's `ToggleSidebar` / `ToggleInspector`.

### D10. Gallery adoption

The gallery's `View` becomes a `Scaffold`: start pane fixed at 260 DIP with
the brand in `StartBar`; search and the page list as its body; the page
header as `TitleBar` (`MinHeight` 52). `CompactAlways` on the landing page
keeps "no sidebar, navigation through the sheet" there. Kept private:
landing-page title-bar content, breadcrumb, search field, virtualized page
list, appearance flyout, zoom, the docs column and page index: they are the
gallery's content, not shell.

Behavior that changes, on purpose: the sidebar collapse animates over
200 ms (it was instant; `DESIGN.md` already prescribes the motion); the
compact sheet is the docked pane (sidebar fill, 260 DIP, 12 DIP gutters)
instead of a 300 DIP padded sheet with a second, unvirtualized list; the
compact toggle keeps the panel icon instead of switching to a menu icon.

## Risks / Trade-offs

- Stacking is one frame late when a control's width changes → the bar
  invalidates on a mismatch; the Tester settles it in the same call.
- The dynamic maximum uses the previous frame's width → a pane can overshoot
  for one frame during a fast window resize; the stored size is not changed.
- Window move, maximize on double click, caption-button hit testing and
  traffic-light placement are the toolkit's and the OS's → headless tests
  assert geometry, input targets and state only; platform behavior is listed
  as unverified in the report and `VALIDATION.md` is not given new claims.
- `Flush` widens `overlay.Sheet`'s surface by one flag → it is the only way
  to reuse the sheet for a pane without a second modal implementation.

## Open Questions

- Whether a published `DESIGN.md` layout rule should fix the bar's control
  size (small, 28 DIP) or leave it to the caller. The text handed over
  states it as guidance.
